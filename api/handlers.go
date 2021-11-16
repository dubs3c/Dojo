package api

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPad implements the Padding interface PKCS7UnPad method.
func PKCS7UnPad(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	char := src[srcLen-1]
	paddingLen := int(char)

	// Hmm, is it correct to check paddingLen == 0?
	// Padding byte should never be \00
	if paddingLen >= srcLen || paddingLen > blockSize || paddingLen == 0 {
		return nil, errors.New("padding error")
	}

	for i := 0; i < paddingLen; i++ {
		if src[srcLen-paddingLen+i] != char {
			return nil, errors.New("padding error")
		}
	}
	return src[:srcLen-paddingLen], nil
}

func responseWithStatusCode(status_code int, w http.ResponseWriter) {
	w.WriteHeader(status_code)
}

func responseWithMessage(message string, status_code int, w http.ResponseWriter) {
	w.WriteHeader(status_code)
	w.Write([]byte(message))
}

func aesEncrypt(key []byte, iv []byte, plaintext []byte) (string, error) {

	padded_plaintext := PKCS7Padding(plaintext, aes.BlockSize)
	log.Println(hex.EncodeToString(padded_plaintext))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("error setting key")
	}

	ciphertext := make([]byte, len(padded_plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, []byte(padded_plaintext))

	hex_ciphertext := hex.EncodeToString(ciphertext)

	return hex_ciphertext, nil
}

func aesDecrypt(key []byte, iv []byte, ciphertext []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		return []byte{}, errors.New("error setting key")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return PKCS7UnPad(ciphertext, aes.BlockSize)
}

// PaddingOracleV1Decrypt Decrypter
func PaddingOracleV1Decrypt(w http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		errString := fmt.Sprintf("'%s' http method is not supported. Please use POST.", req.Method)
		responseWithMessage(errString, 405, w)
		return
	}

	encodedCiphertext := req.PostFormValue("ciphertext")

	if len(encodedCiphertext) == 0 {
		responseWithMessage("please send ciphertext=[your-ciphertext]", 400, w)
		return
	}

	key, _ := b64.StdEncoding.DecodeString("VQk8k3tQ+PDN2F1Ymk4tDLuu5IRJA0WdnFoqwzWIoe4=")
	iv, _ := b64.StdEncoding.DecodeString("/c3GqMto6yAsPFafqJsVtA==")
	ciphertext, err := hex.DecodeString(encodedCiphertext)

	if err != nil {
		responseWithMessage(err.Error(), 400, w)
		return
	}

	_, err = aesDecrypt(key, iv, ciphertext)

	if err != nil {
		responseWithMessage(err.Error(), 400, w)
		return
	}

	responseWithStatusCode(200, w)
}

// PaddingOracleV1Encrypt Encrypter
func PaddingOracleV1Encrypt(w http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		errString := fmt.Sprintf("'%s' http method is not supported. Please use POST.", req.Method)
		responseWithMessage(errString, 405, w)
		return
	}

	inputPlaintext := req.PostFormValue("plaintext")

	if inputPlaintext == "" {
		responseWithMessage("please send plaintext=[your-plaintext]", 400, w)
		return
	}

	key, _ := b64.StdEncoding.DecodeString("VQk8k3tQ+PDN2F1Ymk4tDLuu5IRJA0WdnFoqwzWIoe4=")
	iv, _ := b64.StdEncoding.DecodeString("/c3GqMto6yAsPFafqJsVtA==")
	plaintext := []byte(inputPlaintext)

	log.Println(plaintext)
	ciphertext, err := aesEncrypt(key, iv, plaintext)

	if err != nil {
		responseWithMessage(err.Error(), 400, w)
		return
	}

	responseWithMessage(ciphertext, 200, w)
}
