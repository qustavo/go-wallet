package script

import "bytes"

func NewBytes(bufs ...[]byte) []byte {
	buf := bytes.NewBuffer(nil)
	for _, bz := range bufs {
		if _, err := buf.Write(bz); err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}
