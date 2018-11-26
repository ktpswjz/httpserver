package rsakey

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/ktpswjz/httpserver/security/hash"
	"io/ioutil"
	"os"
)

type Public struct {
	Key *rsa.PublicKey
}

func (s Public) String() string {
	if s.Key == nil {
		return ""
	}

	data, err := x509.MarshalPKIXPublicKey(s.Key)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(data)
}

// openssl rsa -in private.pem -pubout -out public.pem
func (s *Public) LoadFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("Invalid key file")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	s.Key = key.(*rsa.PublicKey)

	return nil
}

func (s *Public) SaveToFile(filePath string) error {
	if s.Key == nil {
		return errors.New("Invalid key")
	}

	data, err := x509.MarshalPKIXPublicKey(s.Key)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: data,
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}

func (s *Public) SaveToMemory() ([]byte, error) {
	if s.Key == nil {
		return nil, errors.New("Invalid key")
	}

	data, err := x509.MarshalPKIXPublicKey(s.Key)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: data,
	}

	return pem.EncodeToMemory(block), nil
}

func (s *Public) GetSize() int {
	if s.Key == nil {
		return 0
	}

	return s.Key.N.BitLen()
}

func (s *Public) Encrypt(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	maxSize := s.Key.N.BitLen()/8 - 11
	dataLength := len(data)
	count := dataLength / maxSize
	offset := 0
	for index := 1; index <= count; index++ {
		vav, err := rsa.EncryptPKCS1v15(rand.Reader, s.Key, data[offset:offset+maxSize])
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(vav)
		if err != nil {
			return nil, err
		}

		offset += maxSize
	}

	if dataLength > offset {
		vav, err := rsa.EncryptPKCS1v15(rand.Reader, s.Key, data[offset:dataLength])
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(vav)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// 验证签名
// hash：签名时使用的哈希方法
// hashed：签名时的哈希值
// sig：经私钥签名的数据
func (s *Public) Verify(hashed []byte, sig []byte) error {
	return rsa.VerifyPKCS1v15(s.Key, crypto.MD5, hashed, sig)
}

// 验证签名
// data: 数据
// signature：签名(base64)
func (s *Public) VerifySign(data []byte, signature string) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	hash, err := hash.MD5Hash(data)
	if err != nil {
		return err
	}

	return s.Verify(hash, sig)
}

// 验证签名
// fileName: 文件名
// signature：签名(base64)
func (s *Public) VerifyFile(fileName string, signature string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	return s.VerifySign(data, signature)
}
