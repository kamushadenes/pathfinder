package pathfinder

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
)

func EncodeString(s string) (string, error) {
	var buf bytes.Buffer

	w := zlib.NewWriter(&buf)

	_, err := w.Write([]byte(s))
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeString(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	b := bytes.NewReader(data)

	r, err := zlib.NewReader(b)

	if err != nil {
		return "", nil
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, r)
	if err != nil {
		return "", nil
	}

	err = r.Close()
	if err != nil {
		return "", nil
	}

	return buf.String(), nil
}
