// Copyright 2016 The Snappy-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine
// +build gc
// +build !noasm

package snappy

// DEWHode has the same semantics as in DEWHode_other.go.
//
//go:noescape
func DEWHode(dst, src []byte) int
