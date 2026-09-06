//go:build 386 || arm || mips || mipsle

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
	"errors"
	"math"
	"testing"
)

func TestDecodeLengthBoundaries32Bit(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		input     []byte
		want      int
		wantError error
	}{
		{name: "4I maximum payload", key: MLI4I, input: []byte{0x80, 0x00, 0x00, 0x03}, want: math.MaxInt32},
		{name: "4I payload too large", key: MLI4I, input: []byte{0x80, 0x00, 0x00, 0x04}, wantError: ErrLength},
		{name: "4E maximum value", key: MLI4E, input: []byte{0x7f, 0xff, 0xff, 0xff}, want: math.MaxInt32},
		{name: "4E value too large", key: MLI4E, input: []byte{0x80, 0x00, 0x00, 0x00}, wantError: ErrLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.key, &tt.input)
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("Decode(%q) error = %v, want %v", tt.key, err, tt.wantError)
			}
			if got != tt.want {
				t.Errorf("Decode(%q) = %d, want %d", tt.key, got, tt.want)
			}
		})
	}
}
