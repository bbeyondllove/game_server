package auth

import (
	"crypto/hmac"
	"crypto/rc4"
	"crypto/sha1"
)

type AuthCrypt struct {
	clientDecrypt *rc4.Cipher
	serverEncrypt *rc4.Cipher
	initialized   bool
}

func (ac *AuthCrypt) Init(k []byte) {
	ServerEncryptionKey := []uint8{0xCC, 0x98, 0xAE, 0x04, 0xE8, 0x97, 0xEA, 0xCA, 0x12, 0xDD, 0xC0, 0x93, 0x42, 0x91, 0x53, 0x57}
	serverEncryptHmac := hmac.New(sha1.New, ServerEncryptionKey)
	serverEncryptHmac.Write(k)
	encryptHash := serverEncryptHmac.Sum(nil)

	ServerDecryptionKey := []uint8{0xC2, 0xB3, 0x72, 0x3C, 0xC6, 0xAE, 0xD9, 0xB5, 0x34, 0x3C, 0x53, 0xEE, 0x2F, 0x43, 0x67, 0xCE}
	clientDecryptHmac := hmac.New(sha1.New, ServerDecryptionKey)
	clientDecryptHmac.Write(k)
	decryptHash := clientDecryptHmac.Sum(nil)

	ac.clientDecrypt, _ = rc4.NewCipher(decryptHash)
	ac.serverEncrypt, _ = rc4.NewCipher(encryptHash)

	ac.initialized = true
}

func (ac *AuthCrypt) DecryptRecv(b []byte) {
	if !ac.initialized {
		return
	}
	ac.clientDecrypt.XORKeyStream(b, b)
}

func (ac *AuthCrypt) EncryptSend(b []byte) {
	if !ac.initialized {
		return
	}
	ac.serverEncrypt.XORKeyStream(b, b)
}
