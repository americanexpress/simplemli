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
	"testing"
)

func BenchmarkEncoding(b *testing.B) {
	mliTypes := []string{
		"2I",
		"2E",
		"4I",
		"4E",
		"2EE",
		"2BCD2",
		"A4E",
	}
	for _, k := range mliTypes {
		b.Run("Encoding "+k, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = Encode(k, 1500)
			}
		})

		x, _ := Encode(k, 1500)
		b.Run("Decoding "+k, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = Decode(k, &x)
			}
		})
	}
}
