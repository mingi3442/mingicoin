package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte { //interface는 base type이라 뭐든지 interface가 될 수 있다
	var aBuffer bytes.Buffer // Buffer는 bytes를 넣으며 read-write
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(i))

}
func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)        //interface를 string으로 변환
	hash := sha256.Sum256([]byte(s)) //hashing
	return fmt.Sprintf("%x", hash)   //16진수로 return

}

func Splitter(s string, sep string, i int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}

func ToJSON(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}
