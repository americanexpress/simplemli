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

type MLICase struct {
	Name           string
	Size           int
	Encoded        string
	Invalid        string
	Value          int
	EmbeddedHeader int
}

func TestMLIs(t *testing.T) {
	mc := []MLICase{
		{
			Name:    "2I",
			Size:    Size2I,
			Encoded: "002d",
			Invalid: "0001",
			Value:   43,
		},
		{
			Name:    "2E",
			Size:    Size2E,
			Encoded: "002b",
			Value:   43,
		},
		{
			Name:    "4I",
			Size:    Size4I,
			Encoded: "00000035",
			Invalid: "00000001",
			Value:   49,
		},
		{
			Name:    "4E",
			Size:    Size4E,
			Encoded: "00000022",
			Value:   34,
		},
		{
			Name:           "2EE",
			Size:           Size2EE,
			Encoded:        "0036",
			Value:          56,
			EmbeddedHeader: 2,
		},
		{
			Name:    "2BCD2",
			Size:    Size2BCD2,
			Encoded: "00000288",
			Invalid: "00000001",
			Value:   284,
		},
		{
			Name:    "A4E",
			Size:    SizeA4E,
			Encoded: "30303433",
			Value:   43,
		},
	}

	// Execute Various Test Cases
	for _, c := range mc {
		t.Run("Decode "+c.Name, func(t *testing.T) {
			b, err := hex.DecodeString(c.Encoded)
			if err != nil {
				t.Errorf("Unable to decode test case sample payload hex - %s", err)
				t.FailNow()
			}

			n, err := Decode(c.Name, &b)
			if err != nil {
				t.Errorf("Unexpected error decoding sample MLI - %s", err)
			}

			if n != c.Value {
				t.Errorf("Unexpected value returned from MLI %s, got %d expected %d", c.Encoded, n, c.Value)
			}
		})

		t.Run("Encode "+c.Name, func(t *testing.T) {
			b, err := Encode(c.Name, c.Value)
			if err != nil {
				t.Errorf("Unable to encode test case length - %s", err)
				t.FailNow()
			}

			if hex.EncodeToString(b) != c.Encoded {
				t.Errorf("Encoded value does not match expectations, got %s, expected %s", hex.EncodeToString(b), c.Encoded)
			}
		})

		t.Run("Encode & Decode "+c.Name, func(t *testing.T) {

			b, err := Encode(c.Name, c.Value)
			if err != nil {
				t.Errorf("Unable to encode test case length - %s", err)
				t.FailNow()
			}

			n, err := Decode(c.Name, &b)
			if err != nil {
				t.Errorf("Unexpected error decoding sample MLI - %s", err)
			}

			if n != c.Value {
				t.Errorf("Unexpected value returned from MLI %s, got %d expected %d", c.Encoded, n, c.Value)
			}

		})

		t.Run("Zero Byte Encode & Decode "+c.Name, func(t *testing.T) {

			b, err := Encode(c.Name, 0+c.EmbeddedHeader)
			if err != nil {
				t.Errorf("Unable to encode test case length - %s", err)
				t.FailNow()
			}

			n, err := Decode(c.Name, &b)
			if err != nil {
				t.Errorf("Unexpected error decoding sample MLI - %s", err)
			}

			if n != 0+c.EmbeddedHeader {
				t.Errorf("Unexpected value returned from MLI %s, got %d expected %d", hex.EncodeToString(b), n, 0)
			}

		})

		if c.Invalid != "" {
			t.Run("Invalid MLI value "+c.Name, func(t *testing.T) {
				b, err := hex.DecodeString(c.Invalid)
				if err != nil {
					t.Errorf("Unable to decode test case sample payload hex - %s", err)
					t.FailNow()
				}

				_, err = Decode(c.Name, &b)
				if !errors.Is(err, ErrLength) {
					t.Errorf("Expected error decoding invalid MLI got %s", err)
				}
			})
		}

	}
}

func TestInvalid(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		_, err := Encode("Invalid", 0)
		if err == nil {
			t.Errorf("Expected error when calling Encode with bad mli type - got nil")
		}
	})

	t.Run("Encode with negative number", func(t *testing.T) {
		_, err := Encode("2I", -1)
		if !errors.Is(err, ErrLength) {
			t.Errorf("Expected ErrLength when calling Encode with a negative number - got %s", err)
		}
	})

	t.Run("Decode", func(t *testing.T) {
		_, err := Decode("Invalid", &empty)
		if err == nil {
			t.Errorf("Expected error when calling Decode with bad mli type - got nil")
		}
	})

	t.Run("A4E Random String", func(t *testing.T) {
		b := []byte("helo")
		_, err := Decode("A4E", &b)
		if err == nil {
			t.Errorf("Expected error when feeding decode a random string - got nil")
		}
	})
}

func TestBadSizedBytes(t *testing.T) {
	tl := map[string]int{
		"2I":    Size2I,
		"2E":    Size2E,
		"4I":    Size4I,
		"4E":    Size4E,
		"2EE":   Size2EE,
		"2BCD2": Size2BCD2,
		"A4E":   SizeA4E,
	}
	for k, v := range tl {
		t.Run(k+" Nil input test", func(t *testing.T) {
			_, err := Decode(k, nil)
			if !errors.Is(err, ErrByteSize) {
				t.Fatalf("Decode(%q, nil) error = %v, want %v", k, err, ErrByteSize)
			}
		})

		t.Run(k+" Bigger than expected test", func(t *testing.T) {
			b := make([]byte, v+10000)
			_, err := Decode(k, &b)
			if err == nil {
				t.Errorf("Expected error when sending too big byte slice to decode sent %d for mli type %s", len(b), k)
			}
		})

		t.Run(k+" Smaller than expected test", func(t *testing.T) {
			b := make([]byte, v-1)
			_, err := Decode(k, &b)
			if err == nil {
				t.Errorf("Expected error when sending too small byte slice to decode sent %d for mli type %s", len(b), k)
			}
		})
	}
}

func TestEncodeLengthBoundaries(t *testing.T) {
	type testCase struct {
		name      string
		key       string
		length    int
		wantHex   string
		wantError error
	}

	tests := []testCase{
		{name: "2I minimum", key: MLI2I, length: 0, wantHex: "0002"},
		{name: "2I maximum", key: MLI2I, length: 65533, wantHex: "ffff"},
		{name: "2I too large", key: MLI2I, length: 65534, wantError: ErrLength},
		{name: "2E minimum", key: MLI2E, length: 0, wantHex: "0000"},
		{name: "2E maximum", key: MLI2E, length: 65535, wantHex: "ffff"},
		{name: "2E too large", key: MLI2E, length: 65536, wantError: ErrLength},
		{name: "2EE minimum with embedded header", key: MLI2EE, length: 2, wantHex: "0000"},
		{name: "2EE maximum", key: MLI2EE, length: 65537, wantHex: "ffff"},
		{name: "2EE below embedded header", key: MLI2EE, length: 1, wantError: ErrLength},
		{name: "2EE too large", key: MLI2EE, length: 65538, wantError: ErrLength},
		{name: "2BCD2 minimum", key: MLI2BCD2, length: 0, wantHex: "00000004"},
		{name: "2BCD2 maximum", key: MLI2BCD2, length: 9995, wantHex: "00009999"},
		{name: "2BCD2 too large", key: MLI2BCD2, length: 9996, wantError: ErrLength},
		{name: "A4E minimum", key: MLIA4E, length: 0, wantHex: "30303030"},
		{name: "A4E maximum", key: MLIA4E, length: 9999, wantHex: "39393939"},
		{name: "A4E too large", key: MLIA4E, length: 10000, wantError: ErrLength},
	}
	tests = append(tests,
		testCase{name: "4I minimum", key: MLI4I, length: 0, wantHex: "00000004"},
		testCase{name: "4E minimum", key: MLI4E, length: 0, wantHex: "00000000"},
	)

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

func TestEncodeDecodeRoundTripBoundaries(t *testing.T) {
	type testCase struct {
		name   string
		key    string
		length int
	}

	tests := []testCase{
		{name: "2I maximum", key: MLI2I, length: 65533},
		{name: "2E maximum", key: MLI2E, length: 65535},
		{name: "2EE maximum", key: MLI2EE, length: 65537},
		{name: "2BCD2 maximum", key: MLI2BCD2, length: 9995},
		{name: "A4E maximum", key: MLIA4E, length: 9999},
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

func TestDecodeMalformedValues(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		input     []byte
		wantError bool
	}{
		{name: "2BCD2 invalid decimal nibble", key: MLI2BCD2, input: []byte{0x00, 0x00, 0x00, 0x0a}, wantError: true},
		{name: "A4E non-decimal ascii", key: MLIA4E, input: []byte("12a4"), wantError: true},
		{name: "A4E zero ascii", key: MLIA4E, input: []byte("0000")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decode(tt.key, &tt.input)

			if tt.wantError && err == nil {
				t.Fatal("Decode() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("Decode() unexpected error: %v", err)
			}
		})
	}
}

func FuzzDecode(f *testing.F) {
	keys := []string{MLI2I, MLI2E, MLI4I, MLI4E, MLI2EE, MLI2BCD2, MLIA4E}
	for _, key := range keys {
		b, err := Encode(key, 43)
		if err != nil {
			f.Fatalf("Encode(%q) seed error: %v", key, err)
		}
		f.Add(key, b)
	}

	f.Fuzz(func(t *testing.T, key string, input []byte) {
		switch key {
		case MLI2I, MLI2E, MLI4I, MLI4E, MLI2EE, MLI2BCD2, MLIA4E:
			_, _ = Decode(key, &input)
		default:
			t.Skip()
		}
	})
}
