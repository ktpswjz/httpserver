package hash

import (
	"encoding/hex"
	"crypto/md5"
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

// 使用MD5计算哈希值
// data: 目标数据
// 返回：二进制哈希值，或者nil如果发生错误
func MD5Hash(data []byte) ([]byte, error) {
	hasher := md5.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
