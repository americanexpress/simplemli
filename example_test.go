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
package simplemli_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/americanexpress/simplemli"
)

func Example() {
	message := []byte("This is a message")
	mli, err := simplemli.Encode(simplemli.MLI2I, len(message))
	if err != nil {
		panic(err)
	}

	stream := bytes.NewReader(append(mli, message...))
	header := make([]byte, simplemli.Size2I)
	if _, err := io.ReadFull(stream, header); err != nil {
		panic(err)
	}

	length, err := simplemli.Decode(simplemli.MLI2I, &header)
	if err != nil {
		panic(err)
	}

	decoded := make([]byte, length)
	if _, err := io.ReadFull(stream, decoded); err != nil {
		panic(err)
	}

	fmt.Println(string(decoded))
	// Output: This is a message
}
