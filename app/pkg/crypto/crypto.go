package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func Md5(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	return hex.EncodeToString(md5.Sum(nil))
}

func Hmac(key, data string) string {
	hmac := hmac.New(sha1.New, []byte(key))
	hmac.Write([]byte(data))
	return hex.EncodeToString(hmac.Sum(nil))
}

func Sha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum(nil))
}

func HMacSHA256(src, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(src))
	return hex.EncodeToString(m.Sum(nil))
}

// AesEncrypt cbc 128位 pkcs7 iv和key保持一致
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// PKCS7Padding pkcs7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 这里区别上面方法：key用的string，更便于使用
func GetAesDecryptString(encData, key string) string {
	baseDecodeStr, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return ""
	}

	data, err := AesDecrypt(baseDecodeStr, []byte(key))
	if err != nil {
		return ""
	}

	return string(data)
}

// AesDecrypt aes解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData, err = PKCS7UnPadding(origData, blockSize)
	return origData, err
}

// PKCS7UnPadding pkcs7
func PKCS7UnPadding(origData []byte, blockSize int) ([]byte, error) {
	var text []byte

	textLength := len(origData)
	if textLength%blockSize != 0 {
		return nil, errors.New("Padded text not multiple of blocksize")
	}

	lastByte := origData[textLength-1]
	if lastByte < 1 || lastByte > 16 {
		return nil, errors.New("Cannot unpad text")
	}

	textPadding := origData[textLength-int(lastByte):]
	validPadding := bytes.Repeat([]byte{lastByte}, int(lastByte))

	if !bytes.Equal(textPadding, validPadding) {
		return nil, errors.New("Invalid padding")
	}

	text = origData[:textLength-int(lastByte)]

	return text, nil
}
