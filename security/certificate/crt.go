package certificate

import (
	"errors"
	"crypto/x509"
	"io/ioutil"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/ktpswjz/httpserver/security/rsakey"
	"crypto/rsa"
)

type Crt struct {
	cert *x509.Certificate
}

func (s *Crt) LoadFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return s.LoadFromMemory(data)
}

func (s *Crt) LoadFromBase64(text string) error {
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return err
	}

	return s.LoadFromMemory(data)
}

func (s *Crt) LoadFromMemory(data []byte) error {
	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("证书内容无效")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.New(fmt.Sprint("解析证书失败:", err))
	}
	s.cert = cert

	return nil
}

func (s *Crt) Certificate() *x509.Certificate {
	return s.cert
}

func (s *Crt) PublicKey() *rsakey.Public {
	if s.cert == nil {
		return nil
	}
	if s.cert.PublicKey == nil {
		return nil
	}

	key, ok := s.cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil
	}

	return &rsakey.Public{Key: key}
}

func (s *Crt) Pool() *x509.CertPool {
	if s.cert == nil {
		return nil
	}

	pool := x509.NewCertPool()
	pool.AddCert(s.cert)

	return pool
}
