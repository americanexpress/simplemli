//go:build amd64 || arm64 || arm64be || ppc64 || ppc64le || mips64 || mips64le || mips64p32 || mips64p32le || riscv64 || s390x || sparc64 || wasm || loong64

/*
 * Copyright 2020 American Express Travel Related Services Company, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */
package simplemli

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestEncodeLengthBoundaries64Bit(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		length    int
		wantHex   string
		wantError error
	}{
		{name: "4I maximum", key: MLI4I, length: 4294967291, wantHex: "ffffffff"},
		{name: "4I too large", key: MLI4I, length: 4294967292, wantError: ErrLength},
		{name: "4E maximum", key: MLI4E, length: 4294967295, wantHex: "ffffffff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.key, tt.length)
			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Fatalf("Encode(%q, %d) error = %v, want %v", tt.key, tt.length, err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("Encode(%q, %d) unexpected error: %v", tt.key, tt.length, err)
			}
			if gotHex := hex.EncodeToString(got); gotHex != tt.wantHex {
				t.Errorf("Encode(%q, %d) = %s, want %s", tt.key, tt.length, gotHex, tt.wantHex)
			}
		})
	}
}

func TestEncodeDecodeRoundTripBoundaries64Bit(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		length int
	}{
		{name: "4I maximum", key: MLI4I, length: 4294967291},
		{name: "4E maximum", key: MLI4E, length: 4294967295},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := Encode(tt.key, tt.length)
			if err != nil {
				t.Fatalf("Encode(%q, %d) unexpected error: %v", tt.key, tt.length, err)
			}

			got, err := Decode(tt.key, &b)
			if err != nil {
				t.Fatalf("Decode(%q) unexpected error: %v", tt.key, err)
			}
			if got != tt.length {
				t.Errorf("Decode(%q) = %d, want %d", tt.key, got, tt.length)
			}
		})
	}
}
