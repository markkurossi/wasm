//
// binary_decode_test.go
//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package binary

import (
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	file, err := os.Open("testdata/cli.wasm")
	if err != nil {
		t.Fatalf("failed to open test data: %v", err)
	}
	defer file.Close()
	err = DecodeBinary(file)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}
}
