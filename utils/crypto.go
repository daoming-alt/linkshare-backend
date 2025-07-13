// utils/crypto.go
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Encrypt encrypts a string using AES
func Encrypt(plainText, key string) (string, error) {
    block, err := aes.NewCipher([]byte(key[:32])) // Use first 32 bytes of key
    if err != nil {
        return "", err
    }

    cipherText := make([]byte, aes.BlockSize+len(plainText))
    iv := cipherText[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))
    return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts a string using AES
func Decrypt(cipherText, key string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(cipherText)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher([]byte(key[:32]))
    if err != nil {
        return "", err
    }

    iv := data[:aes.BlockSize]
    data = data[aes.BlockSize:]
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(data, data)
    return string(data), nil
}