package hash

import (
	"encoding/hex"
	"crypto/md5"
	"crypto/sha256"
	"hash"
	"crypto/sha1"
	"crypto/sha512"
	"errors"
)

const (
	MD5		= 11
	SHA1	= 21
	SHA256	= 22
	SHA384	= 23
	SHA512	= 24
)

// 将哈希值转换成字符串
// hash: 二进制哈希值
// 返回：小写的十六进制字符串
func ToString(hash []byte) string {
	return hex.EncodeToString(hash)
}

// 将十六进制字符串的哈希值转换成二进制
// hash: 十六进制字符串
// 返回：二进制哈希值
func ToHash(hash string) []byte {
	hashed, err := hex.DecodeString(hash)
	if err != nil {
		return nil
	}

	return hashed
}

func Hash(data string, format uint64) (string, error) {
	if len(data) < 1 {
		return data, nil
	}

	var hasher hash.Hash = nil
	if format == MD5 {
		hasher = md5.New()
	} else if format == SHA1 {
		hasher = sha1.New()
	} else if format == SHA256 {
		hasher = sha256.New()
	} else if format == SHA384 {
		hasher = sha512.New384()
	} else if format == SHA512 {
		hasher = sha512.New()
	}

	if hasher == nil {
		return "", errors.New("invalid format")
	}

	hashed, err := calcHash([]byte(data), hasher)
	if err != nil {
		return "", err
	}

	return ToString(hashed), nil
}

// 使用MD5计算哈希值
// data: 目标数据
// 返回：二进制哈希值，或者nil如果发生错误
func MD5Hash(data []byte) ([]byte, error) {
	return calcHash(data, md5.New())
}

// 使用SHA1计算哈希值
// data: 目标数据
// 返回：二进制哈希值，或者nil如果发生错误
func SHA1Hash(data []byte) ([]byte, error) {
	return calcHash(data, sha1.New())
}

// 使用SHA256计算哈希值
// data: 目标数据
// 返回：二进制哈希值，或者nil如果发生错误
func SHA256Hash(data []byte) ([]byte, error) {
	return calcHash(data, sha256.New())
}

// 使用SHA384计算哈希值
// data: 目标数据
// 返回：二进制哈希值，或者nil如果发生错误
func SHA384Hash(data []byte) ([]byte, error) {
	return calcHash(data, sha512.New384())
}

// 使用SHA512计算哈希值
// data: 目标数据
// 返回：二进制哈希值，或者nil如果发生错误
func SHA512Hash(data []byte) ([]byte, error) {
	return calcHash(data, sha512.New())
}

func calcHash(data []byte, hasher hash.Hash) ([]byte, error) {
	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
