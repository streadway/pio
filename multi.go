/*
Copyright (c) 2014, Sean Treadway, SoundCloud Ltd.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

Redistributions in binary form must reproduce the above copyright notice, this
list of conditions and the following disclaimer in the documentation and/or
other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// package pio implements io types for concurrent access
package pio

import (
	"io"
)

type multiReader struct {
	readers chan io.Reader
	current io.Reader
}

// MultiReader returns a reader that consumes all bytes from all readers on src
// until src is closed.
func MultiReader(src chan io.Reader) io.Reader {
	return &multiReader{
		readers: src,
	}
}

// Read implements io.Reader consuming all bytes from all readers.  Returns
// io.EOF after all bytes from all readers are consumed and readers is closed.
func (r *multiReader) Read(b []byte) (int, error) {
	for {
		// current reader until closed
		if r.current == nil {
			var more bool
			r.current, more = <-r.readers
			if !more {
				return 0, io.EOF
			}
		}
		n, err := r.current.Read(b)
		if n > 0 || err != io.EOF {
			if err == io.EOF {
				// we've finished this reader, there could be more on the next
				r.current = nil
				err = nil
			}
			return n, err
		}
		r.current = nil
	}
}
