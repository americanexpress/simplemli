# SimpleMLI

A Message Length Indicator Encoder/Decoder

[![PkgGoDev](https://pkg.go.dev/badge/github.com/americanexpress/simplemli)](https://pkg.go.dev/github.com/americanexpress/simplemli)
[![Build](https://github.com/americanexpress/simplemli/actions/workflows/tests.yml/badge.svg)](https://github.com/americanexpress/simplemli/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/americanexpress/simplemli/badge.svg?branch=main)](https://coveralls.io/github/americanexpress/simplemli?branch=main)

Message Length Indicators (MLI) are commonly used in communications over raw TCP/IP sockets. This method of denoting 
message length is especially popular with users of [ISO 8583](https://en.wikipedia.org/wiki/ISO_8583) messages, a 
common communication protocol for financial transactions.

This package provides an easy-to-use Encoder and Decoder for Message Length Indicators. 

## Usage

```go
package main

import (
	"bytes"
	"io"

	"github.com/americanexpress/simplemli"
)

func writeMessage(w io.Writer, message []byte) error {
	mli, err := simplemli.Encode(simplemli.MLI2I, len(message))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, io.MultiReader(bytes.NewReader(mli), bytes.NewReader(message)))
	return err
}

func readMessage(r io.Reader) ([]byte, error) {
	mli := make([]byte, simplemli.Size2I)
	if _, err := io.ReadFull(r, mli); err != nil {
		return nil, err
	}

	length, err := simplemli.Decode(simplemli.MLI2I, &mli)
	if err != nil {
		return nil, err
	}

	message := make([]byte, length)
	if _, err := io.ReadFull(r, message); err != nil {
		return nil, err
	}
	return message, nil
}
```

## Message Length Types

There are many common ways to encode message lengths and this library attempts to provide the most common MLI types.

Valid Options are listed in the table below.

| Name | Description | `Encode` input range |
| ---- | ----------- | -------------------- |
| 2I | 2-byte network byte order with MLI included | 0–65,533 |
| 2E | 2-byte network byte order with MLI excluded | 0–65,535 |
| 4I | 4-byte network byte order with MLI included | 0–4,294,967,291 on 64-bit platforms; 0–2,147,483,647 on 32-bit platforms |
| 4E | 4-byte network byte order with MLI excluded | 0–4,294,967,295 on 64-bit platforms; 0–2,147,483,647 on 32-bit platforms |
| 2EE | 2-byte network byte order with MLI excluded; message contains an additional 2-byte embedded header | 2–65,537, including the embedded header |
| 2BCD2 | 2-byte header followed by a 2-byte packed-decimal MLI with its length included | 0–9,995 |
| A4E | 4-byte ASCII decimal string with MLI excluded | 0–9,999 |

### Inclusive vs. Exclusive MLI

An inclusive MLI is an MLI type where the length of the Message Length Indicator itself is included in the MLI value.

For example, if a message is 1500 bytes and 2I encoded, the resulting MLI will have a value of 1502. The reverse is 
true for Exclusive MLI types, a 1500-byte message with a 2E encoded MLI will have a value of 1500.

`Decode` handles inclusive and exclusive lengths. A 2I MLI value of 1502 decodes to a 1500-byte message.

2EE is different: its two-byte MLI excludes the MLI itself but describes a message containing an additional two-byte
embedded header. Both `Encode` input and `Decode` output include that embedded header. For example, a 1500-byte body
plus its two-byte embedded header is passed to `Encode` as 1502 and decodes as 1502.

Because the API returns `int`, decoding a 4I or 4E value larger than the platform's maximum `int` returns
`simplemli.ErrLength`.

## Contributing

We welcome Your interest in the American Express Open Source Community on Github. Any Contributor to
any Open Source Project managed by the American Express Open Source Community must accept and sign
an Agreement indicating agreement to the terms below. Except for the rights granted in this 
Agreement to American Express and to recipients of software distributed by American Express, You
reserve all right, title, and interest, if any, in and to Your Contributions. Please
[fill out the Agreement](https://cla-assistant.io/americanexpress/simplemli).

## License

Any contributions made under this project will be governed by the
[Apache License 2.0](./LICENSE).

## Code of Conduct

This project adheres to the [American Express Community Guidelines](./CODE_OF_CONDUCT.md). By
participating, you are expected to honor these guidelines.
