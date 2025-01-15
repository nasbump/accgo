package encrypts

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
)

// 使用 HKDF 从 token 派生密钥
func deriveKey(token []byte) ([]byte, error) {
	hash := sha256.New
	hkdf := hkdf.New(hash, token, nil, nil) // 使用 HKDF 派生密钥
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := io.ReadFull(hkdf, key); err != nil {
		return nil, err
	}
	return key, nil
}

// 加密函数
func GoAhead(plaintext string, token string) (string, error) {
	// 派生密钥
	key, err := deriveKey([]byte(token))
	if err != nil {
		return "", err
	}

	// 创建 ChaCha20-Poly1305 加密器
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return "", err
	}

	// 生成随机 Nonce
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密明文
	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), nil)

	// 将 Nonce 和密文拼接在一起
	encryptedData := append(nonce, ciphertext...)

	// 返回 Base64 编码的加密数据
	return base64.URLEncoding.EncodeToString(encryptedData), nil
}

// 解密函数
func ComeBack(ciphertext string, token string) (string, error) {
	// 派生密钥
	key, err := deriveKey([]byte(token))
	if err != nil {
		return "", err
	}

	// 解码 Base64 加密数据
	encryptedData, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建 ChaCha20-Poly1305 加密器
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return "", err
	}

	// 检查加密数据长度是否合法
	if len(encryptedData) < aead.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	// 提取 Nonce 和实际密文
	nonce := encryptedData[:aead.NonceSize()]
	ciphertextBytes := encryptedData[aead.NonceSize():]

	// 解密密文
	plaintext, err := aead.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	// 返回解密后的明文
	return string(plaintext), nil
}

// func main() {
// 	token := "my-secret-token" // 不固定长度的 token
// 	plaintext := "Hello, this is a secret message!"

// 	// 加密
// 	encrypted, err := encrypt(plaintext, token)
// 	if err != nil {
// 		fmt.Println("Encryption error:", err)
// 		return
// 	}
// 	fmt.Println("Encrypted:", encrypted)

// 	// 解密
// 	decrypted, err := decrypt(encrypted, token)
// 	if err != nil {
// 		fmt.Println("Decryption error:", err)
// 		return
// 	}
// 	fmt.Println("Decrypted:", decrypted)
// }
