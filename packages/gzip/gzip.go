package gzip

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
)

// Uncompress performs gzip uncompression and returns the result. Returns error if the
// input file is not in gzipped format or if the result's size exceeds maxSize.
func Uncompress(src []byte, maxSize int64) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	zr.Multistream(false)
	uncompressed, err := ioutil.ReadAll(io.LimitReader(zr, maxSize+1))
	if err != nil {
		return nil, err
	}
	if len(uncompressed) > int(maxSize) {
		return uncompressed, errors.New("uncompressed file exceeds maxSize")
	}
	return uncompressed, nil
}
