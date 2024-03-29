// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package file

import (
	"bufio"
	"io"
	"os"
)

// https://github.com/tailscale/tailscale/blob/main/util/lineread/lineread.go

// LineReadFile opens name and calls fn for each line. It returns an error if the Open failed
// or once fn returns an error.
func LineReadFile(name string, fn func(line []byte) error) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return LineReadReader(f, fn)
}

// LineReadReader calls fn for each line.
// If fn returns an error, Reader stops reading and returns that error.
// LineReadReader may also return errors encountered reading and parsing from r.
// To stop reading early, use a sentinel "stop" error value and ignore
// it when returned from Reader.
func LineReadReader(r io.Reader, fn func(line []byte) error) error {
	bs := bufio.NewScanner(r)
	for bs.Scan() {
		if err := fn(bs.Bytes()); err != nil {
			return err
		}
	}
	return bs.Err()
}
