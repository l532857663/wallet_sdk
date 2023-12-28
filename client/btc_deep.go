package client

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func (c *BtcClient) BuildPSBTransfer(ins []Input, outs []Output) (interface{}, error) {
	err := c.createPsbtTransaction(ins, outs)
	if err != nil {
		fmt.Printf("BuildPSBTransfer createPsbtTransaction error: %+v\n", err)
		return nil, err
	}
	fmt.Printf("wch------1 c.PBST: %+v\n", c.PsbtUpdater.Upsbt.Inputs[0].NonWitnessUtxo)
	err = c.updatePsbtTransaction(InUtxos)
	if err != nil {
		fmt.Printf("BuildPSBTransfer updatePsbtTransaction error: %+v\n", err)
		return nil, err
	}
	fmt.Printf("wch------2 c.PBST: %+v\n", c.PsbtUpdater.Upsbt.Inputs[0].NonWitnessUtxo.TxHash().String())
	fmt.Printf("wch------2 c.PBST: %+v\n", c.PsbtUpdater.Upsbt.UnsignedTx.TxIn[0].PreviousOutPoint.Hash.String())

	err = c.signPsbtTransaction(inSigners1)
	if err != nil {
		fmt.Printf("BuildPSBTransfer signPsbtTransaction error: %+v\n", err)
		return nil, err
	}
	fmt.Printf("wch------3 c.PBST: %+v\n", c.PsbtUpdater.Upsbt.Inputs[0].WitnessUtxo)

	err = c.signPsbtTransaction(inSigners2)
	if err != nil {
		fmt.Printf("BuildPSBTransfer signPsbtTransaction error: %+v\n", err)
		return nil, err
	}

	raw, err := c.extractPsbtTransaction()
	if err != nil {
		fmt.Printf("BuildPSBTransfer ExtractPsbtTransaction error: %+v\n", err)
		return nil, err
	}

	fmt.Printf("Raw:%s\n", raw)

	return nil, nil
}

func (s *BtcClient) createPsbtTransaction(ins []Input, outs []Output) error {
	var (
		txOuts     []*wire.TxOut    = make([]*wire.TxOut, 0)
		txIns      []*wire.OutPoint = make([]*wire.OutPoint, 0)
		nSequences []uint32         = make([]uint32, 0)
	)
	for _, in := range ins {
		txHash, err := chainhash.NewHashFromStr(in.OutTxId)
		if err != nil {
			return err
		}
		prevOut := wire.NewOutPoint(txHash, in.OutIndex)
		txIns = append(txIns, prevOut)
		nSequences = append(nSequences, wire.MaxTxInSequenceNum)
	}

	for _, out := range outs {
		address, err := btcutil.DecodeAddress(out.Address, s.Params)
		if err != nil {
			return err
		}

		pkScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return err
		}

		txOut := wire.NewTxOut(int64(out.Amount), pkScript)
		txOuts = append(txOuts, txOut)
	}

	cPsbt, err := psbt.New(txIns, txOuts, int32(2), uint32(0), nSequences)
	if err != nil {
		return err
	}
	s.PsbtUpdater, err = psbt.NewUpdater(cPsbt)
	if err != nil {
		return err
	}
	return nil
}

func (s *BtcClient) updatePsbtTransaction(inUtxos []InputUtxo) error {
	for _, v := range inUtxos {
		switch v.UtxoType {
		case NonWitness:
			tx := wire.NewMsgTx(2)
			nonWitnessUtxoHex, err := hex.DecodeString(v.NonWitnessUtxo)
			if err != nil {
				fmt.Printf("hex.DecodeString err: %+v\n", err)
				return err
			}
			err = tx.Deserialize(bytes.NewReader(nonWitnessUtxoHex))
			if err != nil {
				fmt.Printf("tx.Deserialize err: %+v\n", err)
				return err
			}
			err = s.PsbtUpdater.AddInNonWitnessUtxo(tx, v.Index)
			if err != nil {
				fmt.Printf("s.PsbtUpdater.AddInNonWitnessUtxo err: %+v\n", err)
				return err
			}
			err = s.PsbtUpdater.AddInSighashType(v.SighashType, v.Index)
			if err != nil {
				fmt.Printf("s.PsbtUpdater.AddInSighashType err: %+v\n", err)
				return err
			}
			break
		case Witness:
			witnessUtxoScriptHex, err := hex.DecodeString(v.WitnessUtxoPkScript)
			if err != nil {
				return err
			}
			txout := wire.TxOut{Value: int64(v.WitnessUtxoAmount), PkScript: witnessUtxoScriptHex}
			err = s.PsbtUpdater.AddInWitnessUtxo(&txout, v.Index)
			if err != nil {
				return err
			}
			err = s.PsbtUpdater.AddInSighashType(v.SighashType, v.Index)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (s *BtcClient) signPsbtTransaction(inSigners []InputSigner) error {
	for _, v := range inSigners {
		// privateKeyBytes, err := hex.DecodeString(v.Pri)
		// if err != nil {
		// 	return err
		// }
		// privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
		// wif, _ := btcutil.NewWIF(privateKey, s.Params, true)
		// fmt.Printf("wif: %+v\n", wif.String())
		privateKey, err := GetPrikeyByWIF(v.Pri, s.Params)

		sigScript := []byte{}
		switch v.UtxoType {
		case NonWitness:
			sigScript, err = txscript.RawTxInSignature(s.PsbtUpdater.Upsbt.UnsignedTx, v.Index, s.PsbtUpdater.Upsbt.Inputs[v.Index].NonWitnessUtxo.TxOut[s.PsbtUpdater.Upsbt.UnsignedTx.TxIn[v.Index].PreviousOutPoint.Index].PkScript, v.SighashType, privateKey)
			if err != nil {
				return err
			}
			break
		case Witness:
			prevOutputFetcher := NewPrevOutputFetcher(s.PsbtUpdater.Upsbt.Inputs[v.Index].WitnessUtxo.PkScript, s.PsbtUpdater.Upsbt.Inputs[v.Index].WitnessUtxo.Value)
			sigHashes := txscript.NewTxSigHashes(s.PsbtUpdater.Upsbt.UnsignedTx, prevOutputFetcher)
			sigScript, err = txscript.RawTxInWitnessSignature(s.PsbtUpdater.Upsbt.UnsignedTx, sigHashes, v.Index, s.PsbtUpdater.Upsbt.Inputs[v.Index].WitnessUtxo.Value, s.PsbtUpdater.Upsbt.Inputs[v.Index].WitnessUtxo.PkScript, v.SighashType, privateKey)
			if err != nil {
				return err
			}
			break
		}

		fmt.Printf("sigScript: %s\n", hex.EncodeToString(sigScript))
		pubByte, err := hex.DecodeString(v.Pub)
		res, err := s.PsbtUpdater.Sign(v.Index, sigScript, pubByte, nil, nil)
		if err != nil || res != 0 {
			fmt.Printf("wch------ PsbtUpdater.Sign: res: %+v, err: %+v\n", res, err)
			return err
		}
		_, err = psbt.MaybeFinalize(s.PsbtUpdater.Upsbt, v.Index)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BtcClient) extractPsbtTransaction() (string, error) {
	if !s.PsbtUpdater.Upsbt.IsComplete() {
		err := psbt.MaybeFinalizeAll(s.PsbtUpdater.Upsbt)
		if err != nil {
			return "", err
		}
	}

	tx, err := psbt.Extract(s.PsbtUpdater.Upsbt)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = tx.Serialize(&b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b.Bytes()), nil
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
