// Copyright 2015 The go-DEWH Authors
// This file is part of the go-DEWH library.
//
// The go-DEWH library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-DEWH library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-DEWH library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/DEWH/go-DEWH/common"
	"github.com/DEWH/go-DEWH/common/hexutil"
	"github.com/DEWH/go-DEWH/common/math"
	"github.com/DEWH/go-DEWH/core/types"
	"github.com/DEWH/go-DEWH/params"
	"github.com/DEWH/go-DEWH/rlp"
)

// TransactionTest checks RLP DEWHoding and sender derivation of transactions.
type TransactionTest struct {
	json ttJSON
}

type ttJSON struct {
	BlockNumber math.HexOrDEWHimal64 `json:"blockNumber"`
	RLP         hexutil.Bytes       `json:"rlp"`
	Sender      hexutil.Bytes       `json:"sender"`
	Transaction *ttTransaction      `json:"transaction"`
}

//go:generate gencoDEWH -type ttTransaction -field-override ttTransactionMarshaling -out gen_tttransaction.go

type ttTransaction struct {
	Data     []byte         `gencoDEWH:"required"`
	GasLimit uint64         `gencoDEWH:"required"`
	GasPrice *big.Int       `gencoDEWH:"required"`
	Nonce    uint64         `gencoDEWH:"required"`
	Value    *big.Int       `gencoDEWH:"required"`
	R        *big.Int       `gencoDEWH:"required"`
	S        *big.Int       `gencoDEWH:"required"`
	V        *big.Int       `gencoDEWH:"required"`
	To       common.Address `gencoDEWH:"required"`
}

type ttTransactionMarshaling struct {
	Data     hexutil.Bytes
	GasLimit math.HexOrDEWHimal64
	GasPrice *math.HexOrDEWHimal256
	Nonce    math.HexOrDEWHimal64
	Value    *math.HexOrDEWHimal256
	R        *math.HexOrDEWHimal256
	S        *math.HexOrDEWHimal256
	V        *math.HexOrDEWHimal256
}

func (tt *TransactionTest) Run(config *params.ChainConfig) error {
	tx := new(types.Transaction)
	if err := rlp.DEWHodeBytes(tt.json.RLP, tx); err != nil {
		if tt.json.Transaction == nil {
			return nil
		}
		return fmt.Errorf("RLP DEWHoding failed: %v", err)
	}
	// Check sender derivation.
	signer := types.MakeSigner(config, new(big.Int).SetUint64(uint64(tt.json.BlockNumber)))
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return err
	}
	if sender != common.BytesToAddress(tt.json.Sender) {
		return fmt.Errorf("Sender mismatch: got %x, want %x", sender, tt.json.Sender)
	}
	// Check DEWHoded fields.
	err = tt.json.Transaction.verify(signer, tx)
	if tt.json.Sender == nil && err == nil {
		return errors.New("field validations succeeded but should fail")
	}
	if tt.json.Sender != nil && err != nil {
		return fmt.Errorf("field validations failed after RLP DEWHoding: %s", err)
	}
	return nil
}

func (tt *ttTransaction) verify(signer types.Signer, tx *types.Transaction) error {
	if !bytes.Equal(tx.Data(), tt.Data) {
		return fmt.Errorf("Tx input data mismatch: got %x want %x", tx.Data(), tt.Data)
	}
	if tx.Gas() != tt.GasLimit {
		return fmt.Errorf("GasLimit mismatch: got %d, want %d", tx.Gas(), tt.GasLimit)
	}
	if tx.GasPrice().Cmp(tt.GasPrice) != 0 {
		return fmt.Errorf("GasPrice mismatch: got %v, want %v", tx.GasPrice(), tt.GasPrice)
	}
	if tx.Nonce() != tt.Nonce {
		return fmt.Errorf("Nonce mismatch: got %v, want %v", tx.Nonce(), tt.Nonce)
	}
	v, r, s := tx.RawSignatureValues()
	if r.Cmp(tt.R) != 0 {
		return fmt.Errorf("R mismatch: got %v, want %v", r, tt.R)
	}
	if s.Cmp(tt.S) != 0 {
		return fmt.Errorf("S mismatch: got %v, want %v", s, tt.S)
	}
	if v.Cmp(tt.V) != 0 {
		return fmt.Errorf("V mismatch: got %v, want %v", v, tt.V)
	}
	if tx.To() == nil {
		if tt.To != (common.Address{}) {
			return fmt.Errorf("To mismatch when recipient is nil (contract creation): %x", tt.To)
		}
	} else if *tx.To() != tt.To {
		return fmt.Errorf("To mismatch: got %x, want %x", *tx.To(), tt.To)
	}
	if tx.Value().Cmp(tt.Value) != 0 {
		return fmt.Errorf("Value mismatch: got %x, want %x", tx.Value(), tt.Value)
	}
	return nil
}
