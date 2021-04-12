package auth

import "strings"

var ralphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_"
var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_"

func Encrypt(strtoencrypt string) string {
	var password = "kinge383e"
	pos_alpha_ary := []string{}
	for _, c := range password {
		pos := strings.Index(alphabet, string(c))
		if pos == -1 {
			return ""
		}
		pos_alpha_ary = append(pos_alpha_ary, alphabet[pos:])
	}
	var i, n int
	var encrypted_string string
	nn := len(password)
	c := len(strtoencrypt)
	for i < c {
		pos := strings.Index(alphabet, strtoencrypt[i:i+1])
		encrypted_string += pos_alpha_ary[n][pos : pos+1]
		n++
		if n == nn {
			n = 0
		}
		i++
	}
	return encrypted_string
}

func Decrypt(strtodecrypt string) string {
	var password = "kinge383e"
	pos_alpha_ary := []string{}
	for _, c := range password {
		pos := strings.Index(alphabet, string(c))
		if pos == -1 {
			return ""
		}
		pos_alpha_ary = append(pos_alpha_ary, alphabet[pos:])
	}
	var i, n int
	var decrypted_string string
	nn := len(password)
	c := len(strtodecrypt)
	for i < c {
		pos := strings.Index(pos_alpha_ary[n], strtodecrypt[i:i+1])
		decrypted_string += ralphabet[pos : pos+1]
		n++
		if n == nn {
			n = 0
		}
		i++
	}
	return decrypted_string
}
