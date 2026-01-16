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

/*
Package simplemli is a Message Length Indicator Encoder/Decoder.

Message Length Indicators (MLI) are commonly used in communications over raw TCP/IP sockets. This method of denoting
message length is especially popular with users of ISO 8583 messages, a common communication protocol for financial
transactions.

This package provides an easy-to-use Encoder and Decoder for Message Length Indicators.

Usage:

	import (
		"github.com/americanexpress/simplemli"
	)

	func main() {
		msg := []byte("This is a message")

		// Encoding Example
		mli, err := simplemli.Encode(simplemli.MLI2I, len(msg))
		if err != nil {
			// Do something
		}

		// Append the MLI to the message
		msg := append(mli, msg)

		// Write to TCP Connection
		_, err = conn.Write(msg)


		// Reading MLI from TCP Connection
		b := make([]byte, simplemli.Size2I)
		_, err = conn.Read(b) // only read the MLI from buffer
		if err != nil {
			// Do something
		}

		// Decoding Example
		length, err := simplemli.Decode(simplemli.MLI2I, b)
		if err != nil {
			// Do something
		}

		// Reading Message from TCP Connection
		msg := make([]byte, length)
		_, err = conn.Read(msg)
	}

There are many common ways to encode message lengths and this library attempts to provide the most common MLI types.
*/
package simplemli

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
)

// MLI Size in bytes
const (
	Size2I    = 2
	Size2E    = 2
	Size4I    = 4
	Size4E    = 4
	Size2EE   = 2
	Size2BCD2 = 4
	SizeA4E   = 4
)

// Encoding/Decoding argument keys
const (
	// 2-byte network byte order with MLI included
	MLI2I = "2I"

	// 2-byte network byte order with MLI excluded
	MLI2E = "2E"

	// 4-byte network byte order with MLI included
	MLI4I = "4I"

	// 4-byte network byte order with MLI excluded
	MLI4E = "4E"

	// 2-byte network byte order with MLI excluded, additional 2-byte header is included with message
	MLI2EE = "2EE"

	// 2-byte header with a 2-byte binary-coded decimal with MLI excluded
	MLI2BCD2 = "2BCD2"

	// 4-byte ASCII string with MLI excluded
	MLIA4E = "A4E"
)

// ErrByteSize reports an attempt to decode byte data that does not match the expected size for the desired MLI type.
var ErrByteSize = fmt.Errorf("input bytes does not match expected size for selected mli type")

// ErrLength reports an attempt to decode or encode data with an invalid length (e.g., negative numbers).
var ErrLength = fmt.Errorf("invalid mli length provided")

// ErrOverflow reports an attempt to encode a length that is too large for the selected MLI type.
var ErrOverflow = fmt.Errorf("length exceeds maximum allowed for selected mli type")

// Decode accepts a message length in bytes and decodes the value into an integer. The byte slice provided to Decode
// must be the message length indicator itself and not include message headers or body. If the provided byte size does
// not match the expected MLI size, Decode will return an error.
//
// The return value provided by Decode will exclude the length of the MLI and provide the length of the message itself.
// For example, a 2I MLI of 1502 will return 1500 when Decoded.
//
//	length, err := simplemli.Decode(simplemli.MLI2I, b)
//	if err != nil {
//		// Do something
//	}
//
// Note: 2EE Message Length Indicators are unique in that they contain a 2-byte header which is not accounted for in
// the message length. When decoding a 2EE MLI of 1500, the return value will include the header length, 1502.
func Decode(key string, b []byte) (int, error) {
	switch key {
	case MLI2I:
		if len(b) != Size2I {
			return 0, ErrByteSize
		}
		n := int(binary.BigEndian.Uint16(b))
		if n == 0 {
			return 0, nil
		}
		n -= Size2I
		if n < 0 {
			return 0, ErrLength
		}
		return n, nil

	case MLI2E:
		if len(b) != Size2E {
			return 0, ErrByteSize
		}
		return int(binary.BigEndian.Uint16(b)), nil

	case MLI4I:
		if len(b) != Size4I {
			return 0, ErrByteSize
		}
		n := int(binary.BigEndian.Uint32(b))
		if n == 0 {
			return 0, nil
		}
		n -= Size4I
		if n < 0 {
			return 0, ErrLength
		}
		return n, nil

	case MLI4E:
		if len(b) != Size4E {
			return 0, ErrByteSize
		}
		return int(binary.BigEndian.Uint32(b)), nil

	case MLI2EE:
		if len(b) != Size2EE {
			return 0, ErrByteSize
		}
		return int(binary.BigEndian.Uint16(b)) + 2, nil

	case MLI2BCD2:
		if len(b) != Size2BCD2 {
			return 0, ErrByteSize
		}
		// Convert from hex to integer using Binary-Coded Decimal (last 2 bytes)
		n, err := strconv.Atoi(hex.EncodeToString(b[2:4]))
		if err != nil {
			return 0, fmt.Errorf("could not convert bcd string to integer: %w", err)
		}
		if n == 0 {
			return 0, nil
		}
		n -= Size2BCD2
		if n < 0 {
			return 0, ErrLength
		}
		return n, nil

	case MLIA4E:
		if len(b) != SizeA4E {
			return 0, ErrByteSize
		}
		// Check for edge case of "0000"
		isZero := true
		for _, v := range b {
			if v != '0' {
				isZero = false
				break
			}
		}
		if isZero {
			return 0, nil
		}
		n, err := strconv.Atoi(string(b))
		if err != nil {
			return 0, fmt.Errorf("unable to convert ascii values to integer: %w", err)
		}
		return n, nil

	default:
		return 0, fmt.Errorf("invalid mli type provided: %s", key)
	}
}

// Encode will accept a message length type and message length value desired. Encode will return a byte slice which
// contains a MLI formatted for in the desired message length type.
//
// For inclusive MLI types, the Encode function will add the MLI length to the returned encoded MLI. In all cases,
// users should provide the message length without including MLI length.
//
//	b, err := simplemli.Encode(key, len(msg))
//	if err != nil {
//		// Do something
//	}
//
// Note: 2EE Message Length Indicators are unique in that the messages should include a 2-byte embedded header which is
// not accounted for in the MLI. When encoding a 2EE MLI, users should include the embedded header in the length value.
// For example, a message of 1500 bytes, with a 2-byte embedded header will have a 2EE MLI value of 1500.
func Encode(key string, length int) ([]byte, error) {
	if length < 0 {
		return nil, ErrLength
	}

	switch key {
	case MLI2I:
		if length > math.MaxUint16-Size2I {
			return nil, ErrOverflow
		}
		b := make([]byte, Size2I)
		binary.BigEndian.PutUint16(b, uint16(length+Size2I))
		return b, nil

	case MLI2E:
		if length > math.MaxUint16 {
			return nil, ErrOverflow
		}
		b := make([]byte, Size2E)
		binary.BigEndian.PutUint16(b, uint16(length))
		return b, nil

	case MLI4I:
		if length > math.MaxUint32-Size4I {
			return nil, ErrOverflow
		}
		b := make([]byte, Size4I)
		binary.BigEndian.PutUint32(b, uint32(length+Size4I))
		return b, nil

	case MLI4E:
		if length > math.MaxUint32 {
			return nil, ErrOverflow
		}
		b := make([]byte, Size4E)
		binary.BigEndian.PutUint32(b, uint32(length))
		return b, nil

	case MLI2EE:
		if length < Size2EE {
			return nil, ErrLength
		}
		if length-Size2EE > math.MaxUint16 {
			return nil, ErrOverflow
		}
		b := make([]byte, Size2EE)
		binary.BigEndian.PutUint16(b, uint16(length-Size2EE))
		return b, nil

	case MLI2BCD2:
		if length+Size2BCD2 > 9999 {
			return nil, ErrOverflow
		}
		h, err := hex.DecodeString(fmt.Sprintf("%04d", length+Size2BCD2))
		if err != nil {
			return nil, fmt.Errorf("unable to convert length to bcd: %w", err)
		}
		b := make([]byte, Size2BCD2)
		copy(b[2:], h)
		return b, nil

	case MLIA4E:
		if length > 9999 {
			return nil, ErrOverflow
		}
		return []byte(fmt.Sprintf("%04d", length)), nil

	default:
		return nil, fmt.Errorf("invalid mli type provided: %s", key)
	}
}
