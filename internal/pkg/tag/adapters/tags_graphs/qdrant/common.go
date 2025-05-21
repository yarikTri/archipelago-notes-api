package qdrant

import "io"

func StringFromReaderUnfallible(r io.Reader) string {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return ""
	}

	return string(bytes)
}
