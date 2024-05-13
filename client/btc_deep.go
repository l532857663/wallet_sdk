package client

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts"
)

func (c *BtcClient) GenerateSignedListingPSBTBase64(in *Input, out *Output) (interface{}, error) {
	network := c.Params
	txHash, err := chainhash.NewHashFromStr(in.TxId)
	if err != nil {
		return "", err
	}
	prevOut := wire.NewOutPoint(txHash, in.VOut)
	inputs := []*wire.OutPoint{{Index: 0}, {Index: 1}, prevOut}

	pkScript, err := NewPubkeyHash(out.Address, network)
	if err != nil {
		return "", err
	}
	// placeholder
	dummyPkScript, err := NewPubkeyHash(c.Placeholder, network)
	if err != nil {
		return "", err
	}
	outputs := []*wire.TxOut{{PkScript: dummyPkScript}, {PkScript: dummyPkScript}, wire.NewTxOut(out.Amount, pkScript)}

	nSequences := []uint32{wire.MaxTxInSequenceNum, wire.MaxTxInSequenceNum, wire.MaxTxInSequenceNum}
	p, err := psbt.New(inputs, outputs, int32(2), uint32(0), nSequences)
	if err != nil {
		return "", err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return "", err
	}

	dummyWitnessUtxo := wire.NewTxOut(0, dummyPkScript)
	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 0)
	if err != nil {
		return "", err
	}
	err = updater.AddInWitnessUtxo(dummyWitnessUtxo, 1)
	if err != nil {
		return "", err
	}

	prevPkScript, err := NewPubkeyHash(in.Address, network)
	if err != nil {
		return "", err
	}
	witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
	prevOuts := map[wire.OutPoint]*wire.TxOut{
		wire.OutPoint{Index: 0}: dummyWitnessUtxo,
		wire.OutPoint{Index: 1}: dummyWitnessUtxo,
		*prevOut:                witnessUtxo,
	}
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	err = signInput(updater, SellerSignatureIndex, in, prevOutputFetcher, txscript.SigHashSingle|txscript.SigHashAnyOneCanPay, network)
	if err != nil {
		return "", err
	}

	return p.B64Encode()
}

func (c *BtcClient) BuildPSBTransfer(ins []Input, outs []Output) (interface{}, error) {
	var inputs []*wire.OutPoint
	var nSequences []uint32
	network := c.Params

	for _, in := range ins {
		txHash, err := chainhash.NewHashFromStr(in.TxId)
		if err != nil {
			return "", err
		}
		inputs = append(inputs, wire.NewOutPoint(txHash, in.VOut))

		nSequences = append(nSequences, in.Sequence|wire.SequenceLockTimeDisabled)
	}

	var outputs []*wire.TxOut
	for _, out := range outs {
		pkScript, err := NewPubkeyHash(out.Address, network)
		if err != nil {
			return "", err
		}
		outputs = append(outputs, wire.NewTxOut(out.Amount, pkScript))
	}

	p, err := psbt.New(inputs, outputs, int32(2), uint32(0), nSequences)
	if err != nil {
		return "", err
	}

	updater, err := psbt.NewUpdater(p)
	if err != nil {
		return "", err
	}

	for i, in := range ins {
		publicKeyBytes, err := hex.DecodeString(in.PublicKey)
		if err != nil {
			return "", err
		}
		prevPkScript, err := NewPubkeyHash(in.Address, network)
		if err != nil {
			return "", err
		}
		if txscript.IsPayToPubKeyHash(prevPkScript) {
			prevTx := wire.NewMsgTx(2)
			txBytes, err := hex.DecodeString(in.NonWitnessUtxo)
			if err != nil {
				return "", err
			}
			if err := prevTx.Deserialize(bytes.NewReader(txBytes)); err != nil {
				return "", err
			}
			if err := updater.AddInNonWitnessUtxo(prevTx, i); err != nil {
				return "", err
			}
		} else {
			witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
			if err := updater.AddInWitnessUtxo(witnessUtxo, i); err != nil {
				return "", err
			}
			if txscript.IsPayToScriptHash(prevPkScript) {
				redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(publicKeyBytes))
				if err != nil {
					return "", err
				}
				if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
					return "", err
				}
			}
		}

		derivationPath, err := accounts.ParseDerivationPath(in.DerivationPath)
		if err != nil {
			return "", err
		}
		if err := updater.AddInBip32Derivation(in.MasterFingerprint, derivationPath, publicKeyBytes, i); err != nil {
			return "", err
		}
	}

	return nil, nil
}

type PrevOutputFetcher struct {
	pkScript []byte
	value    int64
}

func NewPrevOutputFetcher(pkScript []byte, value int64) *PrevOutputFetcher {
	return &PrevOutputFetcher{
		pkScript,
		value,
	}
}

func (d *PrevOutputFetcher) FetchPrevOutput(wire.OutPoint) *wire.TxOut {
	return &wire.TxOut{
		Value:    d.value,
		PkScript: d.pkScript,
	}
}

func signInput(updater *psbt.Updater, i int, in *Input, prevOutFetcher *txscript.MultiPrevOutFetcher, hashType txscript.SigHashType, network *chaincfg.Params) error {
	wif, err := btcutil.DecodeWIF(in.PrivateKey)
	if err != nil {
		return err
	}
	privKey := wif.PrivKey

	prevPkScript, err := NewPubkeyHash(in.Address, network)
	if err != nil {
		return err
	}
	if txscript.IsPayToPubKeyHash(prevPkScript) {
		prevTx := wire.NewMsgTx(2)
		txBytes, err := hex.DecodeString(in.NonWitnessUtxo)
		if err != nil {
			return err
		}
		if err = prevTx.Deserialize(bytes.NewReader(txBytes)); err != nil {
			return err
		}
		if err = updater.AddInNonWitnessUtxo(prevTx, i); err != nil {
			return err
		}
	} else {
		witnessUtxo := wire.NewTxOut(in.Amount, prevPkScript)
		if err = updater.AddInWitnessUtxo(witnessUtxo, i); err != nil {
			return err
		}
	}

	if err = updater.AddInSighashType(hashType, i); err != nil {
		return err
	}

	if txscript.IsPayToTaproot(prevPkScript) {
		internalPubKey := schnorr.SerializePubKey(privKey.PubKey())
		updater.Upsbt.Inputs[i].TaprootInternalKey = internalPubKey

		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)
		if hashType == txscript.SigHashAll {
			hashType = txscript.SigHashDefault
		}
		witness, err := txscript.TaprootWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes,
			i, in.Amount, prevPkScript, hashType, privKey)
		if err != nil {
			return err
		}

		updater.Upsbt.Inputs[i].TaprootKeySpendSig = witness[0]
	} else if txscript.IsPayToPubKeyHash(prevPkScript) {
		signature, err := txscript.RawTxInSignature(updater.Upsbt.UnsignedTx, i, prevPkScript, hashType, privKey)
		if err != nil {
			return err
		}
		signOutcome, err := updater.Sign(i, signature, privKey.PubKey().SerializeCompressed(), nil, nil)
		if err != nil {
			return err
		}
		if signOutcome != psbt.SignSuccesful {
			return err
		}
	} else {
		pubKeyBytes := privKey.PubKey().SerializeCompressed()
		sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)

		script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
		if err != nil {
			return err
		}
		signature, err := txscript.RawTxInWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes, i, in.Amount, script, hashType, privKey)
		if err != nil {
			return err
		}

		if txscript.IsPayToScriptHash(prevPkScript) {
			redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return err
			}
			err = updater.AddInRedeemScript(redeemScript, i)
			if err != nil {
				return err
			}
		}

		signOutcome, err := updater.Sign(i, signature, pubKeyBytes, nil, nil)
		if err != nil {
			return err
		}
		if signOutcome != psbt.SignSuccesful {
			return err
		}
	}
	return nil
}
