package auth

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"testing"
)

//加解密测试
func TestEnDecrypt(t *testing.T) {
	username := "D01test123"
	password := "RXOZDqu"

	h := sha1.New()
	h.Write([]byte(username + ":" + password))
	key := h.Sum(nil)

	var crypt AuthCrypt
	var clinetCrypt ClientAuthCrypt
	crypt.Init(key)
	clinetCrypt.Init(key)

	header := new(bytes.Buffer)
	var mlen uint16 = 108
	var opcode uint16 = 9
	binary.Write(header, binary.LittleEndian, mlen)
	binary.Write(header, binary.LittleEndian, opcode)

	crypt.EncryptSend(header.Bytes())

	b := make([]byte, 4)
	copy(b, header.Bytes())

	clinetCrypt.DecryptRecv(b)
	msgLen := binary.LittleEndian.Uint16(b[:2])
	opCode := binary.LittleEndian.Uint16(b[2:])

	fmt.Printf("after DecryptRecv() msglen=%d, opCode=%d", msgLen, opCode)

}
