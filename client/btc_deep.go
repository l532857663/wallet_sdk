package client

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts"
)

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
