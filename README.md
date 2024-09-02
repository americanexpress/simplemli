# SimpleMLI

A Message Length Indicator Encoder/Decoder

[![PkgGoDev](https://pkg.go.dev/badge/github.com/americanexpress/simplemli)](https://pkg.go.dev/github.com/americanexpress/simplemli)
[![Build](https://github.com/americanexpress/simplemli/actions/workflows/tests.yml/badge.svg)](https://github.com/americanexpress/simplemli/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/americanexpress/simplemli/badge.svg?branch=main)](https://coveralls.io/github/americanexpress/simplemli?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/americanexpress/simplemli)](https://goreportcard.com/report/github.com/americanexpress/simplemli)

Message Length Indicators (MLI) are commonly used in communications over raw TCP/IP sockets. This method of denoting 
message length is especially popular with users of [ISO 8583](https://en.wikipedia.org/wiki/ISO_8583) messages, a 
common communication protocol for financial transactions.

This package provides an easy-to-use Encoder and Decoder for Message Length Indicators. 

## Usage:

```golang
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
	_, err = conn.Read(&b) // only read the MLI from buffer
	if err != nil {
		// Do something
	}
	
	// Decoding Example
	length, err := simplemli.Decode(simplemli.MLI2I, &b)
	if err != nil {
		// Do something
	}
	
	// Reading Message from TCP Connection
	msg := make([]byte, length)
	_, err = conn.Read(&msg)
}
```

## Message Length Types

There are many common ways to encode message lengths and this library attempts to provide the most common MLI types.

Valid Options are listed in the table below.

| Name | Description |
| ---- | -------- |
| 2I | 2-byte network byte order with MLI included|
| 2E | 2-byte network byte order with MLI excluded |
| 4I | 4-byte network byte order with MLI included |
| 4E | 4-byte network byte order with MLI excluded |
| 2EE | 2-byte network byte order with MLI excluded, additional 2-byte header is included with message |
| 2BCD2 | 2-byte Header with a 2-byte binary-coded decimal with MLI excluded |
| A4E | 4-byte ASCII string with MLI excluded |

### Inclusive vs. Exclusive MLI

An inclusive MLI is an MLI type where the length of the Message Length Indicator itself is included in the MLI value.

For example, if a message is 1500 bytes and 2I encoded, the resulting MLI will have a value of 1502. The reverse is 
true for Exclusive MLI types, a 1500-byte message with a 2E encoded MLI will have a value of 1500.

When calling the Decoder, the MLI inclusive/exclusive nature is already taken care of. If you pass an MLI with a value 
of 1502 and decode it with 2I encoding. The resulting integer will be 1500.

## Contributing

We welcome Your interest in the American Express Open Source Community on Github. Any Contributor to
any Open Source Project managed by the American Express Open Source Community must accept and sign
an Agreement indicating agreement to the terms below. Except for the rights granted in this 
Agreement to American Express and to recipients of software distributed by American Express, You
reserve all right, title, and interest, if any, in and to Your Contributions. Please
[fill out the Agreement](https://cla-assistant.io/americanexpress/simplemli).

## License

Any contributions made under this project will be governed by the
[Apache License 2.0](./LICENSE.txt).

## Code of Conduct

This project adheres to the [American Express Community Guidelines](./CODE_OF_CONDUCT.md). By
participating, you are expected to honor these guidelines.
