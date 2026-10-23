package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

func writeString(w io.Writer, s string) error {
	if err := binary.Write(w, binary.BigEndian, int16(len(s))); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

func readString(r io.Reader) (string, error) {
	var length int16
	if err := binary.Read(r, binary.BigEndian, length); err != nil {
		return "", err
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func writeBytes(w io.Writer, b []byte) error {
	if err := binary.Write(w, binary.BigEndian, int32(len(b))); err != nil {
		return err
	}

	_, err := w.Write(b)
	return err
}

func readBytes(r io.Reader) ([]byte, error) {
	var length int32
	if err := binary.Read(r, binary.BigEndian, length); err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, nil
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func EncodeProduceRequest() {

}

func EncodeFetchRequest() {

}

// Helper for pre-fixed length framing
func ReadFrame(r io.Reader) ([]byte, error) {
	var lenght int32
	if err := binary.Read(r, binary.BigEndian, &lenght); err != nil {
		return nil, err
	}

	if lenght < 0 || lenght > 16<<20 { // 16 MB limit
		return nil, fmt.Errorf("Invalid frame limit %d", lenght)
	}

	buf := make([]byte, lenght)
	if n, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	} else {
		log.Printf("%d bytes copied\n", n)
	}
	return buf, nil
}

func WriteFrame(w io.Writer, payload []byte) error {
	if err := binary.Write(w, binary.BigEndian, int32(len(payload))); err != nil {
		return err
	}

	_, err := w.Write(payload)
	return err
}
