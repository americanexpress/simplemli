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
Package simplemli encodes and decodes Message Length Indicators (MLIs), which
frame messages sent over stream transports such as TCP. MLIs are commonly used
with ISO 8583 financial messages.

[Encode] accepts a message length and returns its MLI. [Decode] accepts only the
MLI bytes and returns the message length. Inclusive MLI types include the MLI's
own size in their wire value; Encode and Decode add or remove that size for the
caller. For 2EE, the input and returned length include the message's additional
two-byte embedded header.
*/
package simplemli

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

// empty is used as a quick return during errors
var empty = make([]byte, 0)

// Internal wire-format limits use int64 so validation behaves consistently on 32-bit and 64-bit platforms.
const (
	// maxUint16Length is the largest value representable by a two-byte unsigned binary MLI.
	maxUint16Length = int64(1<<16 - 1)

	// maxUint32Length is the largest value representable by a four-byte unsigned binary MLI.
	maxUint32Length = int64(1<<32 - 1)

	// maxBCDLength is the largest four-digit value representable by the two-byte packed-decimal field in a 2BCD2 MLI.
	maxBCDLength = int64(9999)

	// maxA4ELength is the largest four-digit decimal value representable by a four-byte ASCII MLI.
	maxA4ELength = int64(9999)
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
var ErrByteSize = errors.New("input bytes does not match expected size for selected mli type")

// ErrLength reports an attempt to decode or encode data with an invalid length, such as a negative value or an
// encoded value that cannot fit in the selected MLI type.
var ErrLength = errors.New("invalid mli length provided")

// Decode converts an MLI into its message length. b must contain only the MLI,
// without the message header or body. Decode returns [ErrByteSize] when b is nil
// or its size does not match key.
//
// The return value provided by Decode will exclude the length of the MLI and provide the length of the message itself.
// For example, a 2I MLI of 1502 will return 1500 when Decoded.
//
//	length, err := simplemli.Decode(simplemli.MLI2I, &b)
//	if err != nil {
//		// Do something
//	}
//
// For 2EE, the returned length includes the message's additional two-byte
// embedded header. A 2EE MLI value of 1500 therefore returns 1502.
//
// On platforms with a 32-bit int, Decode returns [ErrLength] when a 4I or 4E
// result cannot fit in int.
func Decode(key string, b *[]byte) (int, error) {
	if b == nil {
		return 0, ErrByteSize
	}

	switch key {
	case MLI2I:
		// Validate length vs. expected length
		if len(*b) != Size2I {
			return 0, ErrByteSize
		}

		// Convert to integer using Network Byte Order
		n := int(binary.BigEndian.Uint16(*b))
		// If 0 return right away
		if n == 0 {
			return 0, nil
		}

		// Remove MLI length and validate message length is valid
		n = n - Size2I
		if n < 0 {
			return 0, ErrLength
		}
		return n, nil

	case MLI2E:
		// Validate length vs expected length
		if len(*b) != Size2E {
			return 0, ErrByteSize
		}

		// Convert to integer using Network Byte Order
		n := int(binary.BigEndian.Uint16(*b))
		return n, nil

	case MLI4I:
		// Validate length vs expected length
		if len(*b) != Size4I {
			return 0, ErrByteSize
		}

		// Convert to integer using Network Byte Order
		n := binary.BigEndian.Uint32(*b)
		// If 0 return right away
		if n == 0 {
			return 0, nil
		}

		// Remove MLI length and validate message length is valid
		if n < Size4I {
			return 0, ErrLength
		}
		return uint32ToInt(n - Size4I)

	case MLI4E:
		// Validate length vs expected length
		if len(*b) != Size4E {
			return 0, ErrByteSize
		}

		// Convert to integer using Network Byte Order
		return uint32ToInt(binary.BigEndian.Uint32(*b))

	case MLI2EE:
		// Validate length vs expected length
		if len(*b) != Size2EE {
			return 0, ErrByteSize
		}

		// Convert to integer using Network Byte Order
		n := int(binary.BigEndian.Uint16(*b)) + 2 // add 2-byte header length
		return n, nil

	case MLI2BCD2:
		// Validate length vs expected length
		if len(*b) != Size2BCD2 {
			return 0, ErrByteSize
		}

		// Convert from hex to integer using Binary-Coded Decimal
		n, err := strconv.Atoi(hex.EncodeToString((*b)[2:4]))
		if err != nil {
			return 0, fmt.Errorf("could not convert hex string to integer: %w", err)
		}
		// If 0 return right away
		if n == 0 {
			return 0, nil
		}

		// Remove MLI length and validate message length is valid
		n = n - Size2BCD2
		if n < 0 {
			return 0, ErrLength
		}
		return n, nil

	case MLIA4E:
		// Validate length vs expected length
		if len(*b) != SizeA4E {
			return 0, ErrByteSize
		}

		// Check for edge case of 0 in hex format
		if string(*b) == "0000" {
			return 0, nil
		}

		// Convert to integer from ASCII
		n, err := strconv.Atoi(string(*b))
		if err != nil {
			return 0, fmt.Errorf("unable to convert string values to integer: %w", err)
		}
		return n, nil

	default:
		// Preserve historical error text for callers that compare Error() output.
		return 0, fmt.Errorf("Invalid MLI type provided") //nolint:staticcheck
	}
}

// Encode converts length into the MLI selected by key. length must exclude the
// MLI itself. Encode returns [ErrLength] when length is negative or outside the
// selected format's range.
//
// For inclusive MLI types, the Encode function will add the MLI length to the returned encoded MLI. In all cases,
// users should provide the message length without including MLI length.
//
//	b, err := simplemli.Encode(key, len(msg))
//	if err != nil {
//		// Do something
//	}
//
// Supported length ranges are:
//
//	2I:    0 through 65,533
//	2E:    0 through 65,535
//	4I:    0 through 4,294,967,291, limited by the platform's maximum int
//	4E:    0 through 4,294,967,295, limited by the platform's maximum int
//	2EE:   2 through 65,537
//	2BCD2: 0 through 9,995
//	A4E:   0 through 9,999
//
// For 2EE, length includes the message's additional two-byte embedded header.
// A 1500-byte body plus its header is passed as 1502 and produces an MLI value
// of 1500.
func Encode(key string, length int) ([]byte, error) {
	// Reject negative values
	if length < 0 {
		return empty, ErrLength
	}

	switch key {
	case MLI2I:
		if err := validateLength(length, 0, maxUint16Length-Size2I); err != nil {
			return empty, err
		}

		// Create MLI in Network Byte Order
		b := make([]byte, Size2I)
		binary.BigEndian.PutUint16(b, uint16(length+Size2I)) // include mli size
		return b, nil

	case MLI2E:
		if err := validateLength(length, 0, maxUint16Length); err != nil {
			return empty, err
		}

		// Create MLI in Network Byte Order
		b := make([]byte, Size2E)
		binary.BigEndian.PutUint16(b, uint16(length))
		return b, nil

	case MLI4I:
		if err := validateLength(length, 0, maxUint32Length-Size4I); err != nil {
			return empty, err
		}

		// Create MLI in Network Byte Order
		b := make([]byte, Size4I)
		binary.BigEndian.PutUint32(b, uint32(length+Size4I)) // include mli size
		return b, nil

	case MLI4E:
		if err := validateLength(length, 0, maxUint32Length); err != nil {
			return empty, err
		}

		// Create MLI in Network Byte Order
		b := make([]byte, Size4E)
		binary.BigEndian.PutUint32(b, uint32(length))
		return b, nil

	case MLI2EE:
		if err := validateLength(length, Size2EE, maxUint16Length+Size2EE); err != nil {
			return empty, err
		}

		// Create MLI in Network Byte Order
		b := make([]byte, Size2EE)
		binary.BigEndian.PutUint16(b, uint16(length-Size2EE)) // remove embedded 2-byte header length
		return b, nil

	case MLI2BCD2:
		if err := validateLength(length, 0, maxBCDLength-Size2BCD2); err != nil {
			return empty, err
		}

		// Create MLI in Binary-Coded Decimal
		h, err := hex.DecodeString(fmt.Sprintf("%04d", length+Size2BCD2)) // %04d is binary-coded decimal format, wrap in hex
		if err != nil {
			return empty, fmt.Errorf("unable to convert length to hex binary-coded decimal: %w", err)
		}
		// Create empty 2-byte header
		b := make([]byte, 2)
		b = append(b, h...)
		return b, nil

	case MLIA4E:
		if err := validateLength(length, 0, maxA4ELength); err != nil {
			return empty, err
		}

		// Create MLI as a four-byte ASCII decimal string
		return []byte(fmt.Sprintf("%04d", length)), nil

	default:
		// Preserve historical error text for callers that compare Error() output.
		return empty, fmt.Errorf("Invalid MLI type provided") //nolint:staticcheck
	}
}

// uint32ToInt rejects wire values that cannot fit in the platform's int type.
func uint32ToInt(n uint32) (int, error) {
	if uint64(n) > uint64(^uint(0)>>1) {
		return 0, ErrLength
	}
	return int(n), nil
}

func validateLength(length int, minLength, maxLength int64) error {
	n := int64(length)
	if n < minLength || n > maxLength {
		return ErrLength
	}
	return nil
}
