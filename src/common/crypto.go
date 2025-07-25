package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// 加密函数，使用AES-GCM模式
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// 解密函数，使用AES-GCM模式
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// 生成指定长度的密钥
func GenerateKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	return key, err
}

// 密钥处理，确保长度符合AES要求(16, 24, 32字节)
func ProcessKey(key string) []byte {
	keyBytes := []byte(key)
	keyLen := len(keyBytes)

	// 截断或填充密钥以满足AES要求
	switch {
	case keyLen >= 32:
		return keyBytes[:32]
	case keyLen >= 24:
		return keyBytes[:24]
	default:
		// 不足16字节则填充
		padded := make([]byte, 16)
		copy(padded, keyBytes)
		return padded
	}
}
