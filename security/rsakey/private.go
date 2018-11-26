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
	"fmt"
	"github.com/ktpswjz/httpserver/security/hash"
	"io/ioutil"
	"os"
)

type Private struct {
	Key *rsa.PrivateKey
}

// size can be 1024, 2048, ...
func (s *Private) Create(size int) error {
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return err
	}

	s.Key = key

	return nil
}

// openssl genrsa -out private.pem 2048
func (s *Private) LoadFromFile(filePath, password string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("私钥文件内容无效")
	}

	blockData := block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		blockData, err = x509.DecryptPEMBlock(block, []byte(password))
		if err != err {
			return errors.New(fmt.Sprint("解密错误:", err))
		}
	}

	key, err := x509.ParsePKCS1PrivateKey(blockData)
	if err != nil {
		return err
	}

	s.Key = key

	return nil
}

func (s *Private) LoadFromMemory(data []byte, password string) error {

	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("Invalid key value")
	}

	blockData := block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		var err error
		blockData, err = x509.DecryptPEMBlock(block, []byte(password))
		if err != err {
			return errors.New(fmt.Sprint("解密错误:", err))
		}
	}

	key, err := x509.ParsePKCS1PrivateKey(blockData)
	if err != nil {
		return err
	}

	s.Key = key

	return nil
}

func (s *Private) SaveToFile(filePath string, password string) error {
	if s.Key == nil {
		return errors.New("Invalid key")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var block *pem.Block
	if len(password) > 0 {
		block, err = x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(s.Key), []byte(password), x509.PEMCipherAES256)
		if err != nil {
			return err
		}
	} else {
		block = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(s.Key),
		}
	}

	return pem.Encode(file, block)
}

func (s *Private) SaveToMemory(password string) ([]byte, error) {
	if s.Key == nil {
		return nil, errors.New("Invalid key")
	}

	var block *pem.Block
	var err error
	if len(password) > 0 {
		block, err = x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(s.Key), []byte(password), x509.PEMCipherAES256)
		if err != nil {
			return nil, err
		}
	} else {
		block = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(s.Key),
		}
	}

	return pem.EncodeToMemory(block), nil
}

func (s *Private) GetSize() int {
	if s.Key == nil {
		return 0
	}

	return s.Key.N.BitLen()
}

func (s *Private) PublicKey() (*Public, error) {
	if s.Key == nil {
		return nil, errors.New("Invalid private key")
	}

	data, err := x509.MarshalPKIXPublicKey(&s.Key.PublicKey)
	if err != err {
		return nil, err
	}
	key, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}

	return &Public{Key: key.(*rsa.PublicKey)}, nil
}

// 解密
// data：经RSA公钥加密的数据
func (s *Private) Decrypt(data []byte) ([]byte, error) {
	if s.Key == nil {
		return nil, errors.New("Invalid key")
	}

	var buf bytes.Buffer
	maxSize := s.Key.N.BitLen() / 8
	dataLength := len(data)
	count := dataLength / maxSize
	offset := 0
	for index := 1; index <= count; index++ {
		vav, err := rsa.DecryptPKCS1v15(rand.Reader, s.Key, data[offset:offset+maxSize])
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
		vav, err := rsa.DecryptPKCS1v15(rand.Reader, s.Key, data[offset:dataLength])
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

// 解密
// data：经RSA公钥加密的数据取base64
func (s *Private) DecryptData(data string) ([]byte, error) {
	if s.Key == nil {
		return nil, errors.New("Invalid key")
	}

	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return s.Decrypt(buf)
}

// 签名
// hashed：MD5哈希值
func (s *Private) Sign(hashed []byte) ([]byte, error) {
	if s.Key == nil {
		return nil, errors.New("Invalid key")
	}

	return rsa.SignPKCS1v15(rand.Reader, s.Key, crypto.MD5, hashed)
}

// 签名数据
// data: 数据
// 返回：对数据取MD5，然后对MD5进行签名，返回签名的base64
func (s *Private) SignData(data []byte) (string, error) {
	if s.Key == nil {
		return "", errors.New("Invalid key")
	}

	hash, err := hash.MD5Hash(data)
	if err != nil {
		return "", err
	}

	sig, err := s.Sign(hash)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

// 签名文件
// fileName: 文件名
// 返回：对文件数据取MD5，然后对MD5进行签名，返回签名的base64
func (s *Private) SignFile(filePath string) (string, error) {
	if s.Key == nil {
		return "", errors.New("Invalid key")
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return s.SignData(data)
}
