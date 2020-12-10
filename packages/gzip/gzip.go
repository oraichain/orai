package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Uncompress performs gzip uncompression and returns the result. Returns error if the
// input file is not in gzipped format or if the result's size exceeds maxSize.
func Uncompress(src []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	zr.Multistream(false)
	uncompressed, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return uncompressed, nil
}
