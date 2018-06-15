package certificate

import (
	"crypto/tls"
	"io/ioutil"
	"golang.org/x/crypto/pkcs12"
	"errors"
	"fmt"
	"encoding/pem"
	"crypto/x509"
	"github.com/ktpswjz/httpserver/security/rsakey"
	"crypto/rsa"
)

type Pfx struct {
	cert *tls.Certificate
}

func (s *Pfx) LoadFromFile(filePath, password string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return s.LoadFromMemory(data, password)
}

func (s *Pfx) LoadFromMemory(data []byte, password string) error {
	certBlocks, err := pkcs12.ToPEM(data, password)
	if err != nil {
		return errors.New(fmt.Sprint("解析证书失败:", err))
	}

	blockData := make([]byte, 0)
	for _, b := range certBlocks {
		blockData = append(blockData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(blockData, blockData)
	if err != nil {
		return errors.New(fmt.Sprint("构造证书失败:", err))
	}
	s.cert = &cert

	return nil
}

func (s *Pfx) TlsCertificates() []tls.Certificate {
	if s.cert == nil {
		return nil
	}

	return []tls.Certificate{*s.cert}
}

func (s *Pfx) TlsCertificate() *tls.Certificate {
	return s.cert
}

func (s *Pfx) Certificate() *x509.Certificate {
	if s.cert == nil {
		return nil
	}

	if len(s.cert.Certificate) == 0 {
		return nil
	}

	cert, err := x509.ParseCertificate(s.cert.Certificate[0])
	if err != nil {
		return nil
	}

	return cert
}

func (s *Pfx) PrivateKey() *rsakey.Private {
	if s.cert == nil {
		return nil
	}
	if s.cert.PrivateKey == nil {
		return nil
	}

	key, ok := s.cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil
	}

	return &rsakey.Private{Key: key}
}

func (s *Pfx) PublicKey() *rsakey.Public {
	privateKey := s.PrivateKey()
	if privateKey == nil {
		return nil
	}

	publicKey, err := privateKey.PublicKey()
	if err != nil {
		return nil
	}

	return publicKey
}
