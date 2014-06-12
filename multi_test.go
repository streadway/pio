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

package pio

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func TestIdentity(t *testing.T) {
	const (
		blockSize = 4 * 1024
		fileSize  = 12 * blockSize
	)

	buf := &bytes.Buffer{}
	io.CopyN(buf, rand.Reader, 12*blockSize)

	readers := make(chan io.Reader, 1)
	go func() {
		for n := 0; n < fileSize; n += blockSize {
			readers <- bytes.NewReader(buf.Bytes()[n : n+blockSize])
		}
		close(readers)
	}()

	res := &bytes.Buffer{}
	io.Copy(res, MultiReader(readers))

	if want, got := buf.Bytes(), res.Bytes(); bytes.Compare(want, got) != 0 {
		t.Fatalf("expected parallel read to be equal to original, len(want): %v, len(got): %v", len(want), len(got))
	}
}
