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
		MLICase{
			Name:    "2I",
			Size:    Size2I,
			Encoded: "002d",
			Invalid: "0001",
			Value:   43,
		},
		MLICase{
			Name:    "2E",
			Size:    Size2E,
			Encoded: "002b",
			Value:   43,
		},
		MLICase{
			Name:    "4I",
			Size:    Size4I,
			Encoded: "00000035",
			Invalid: "00000001",
			Value:   49,
		},
		MLICase{
			Name:    "4E",
			Size:    Size4E,
			Encoded: "00000022",
			Value:   34,
		},
		MLICase{
			Name:           "2EE",
			Size:           Size2EE,
			Encoded:        "0036",
			Value:          56,
			EmbeddedHeader: 2,
		},
		MLICase{
			Name:    "2BCD2",
			Size:    Size2BCD2,
			Encoded: "00000288",
			Invalid: "00000001",
			Value:   284,
		},
		MLICase{
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
				if err == nil || err != ErrLength {
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
		if err == nil {
			t.Errorf("Expected error when calling Encode with a negative number - got nil")
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
	}
	for k, v := range tl {
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
