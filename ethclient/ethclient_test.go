// Copyright 2016 The go-DEC Authors
// This file is part of the go-DEC library.
//
// The go-DEC library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-DEC library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-DEC library. If not, see <http://www.gnu.org/licenses/>.

package ethclient

import "github.com/DEC/go-DEC"

// Verify that Client implements the DEC interfaces.
var (
	_ = DEC.ChainReader(&Client{})
	_ = DEC.TransactionReader(&Client{})
	_ = DEC.ChainStateReader(&Client{})
	_ = DEC.ChainSyncReader(&Client{})
	_ = DEC.ContractCaller(&Client{})
	_ = DEC.GasEstimator(&Client{})
	_ = DEC.GasPricer(&Client{})
	_ = DEC.LogFilterer(&Client{})
	_ = DEC.PendingStateReader(&Client{})
	// _ = DEC.PendingStateEventer(&Client{})
	_ = DEC.PendingContractCaller(&Client{})
)
