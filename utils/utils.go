package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func ToBytes(i interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	HandleError(encoder.Encode(i))
	return buffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleError(decoder.Decode(i))
}

func Hash(i interface{}) string {
	s := fmt.Sprint(i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}
