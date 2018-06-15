package certificate

import (
	"math/big"
	"time"
	"crypto/x509/pkix"
	"io/ioutil"
	"crypto/x509"
	"errors"
)

type RevokedItem struct {
	SerialNumber       *big.Int
	RevocationTime     time.Time
	Organization       string
	OrganizationalUnit string
	Locality           string
	NotBefore          *time.Time
	NotAfter           *time.Time
}

type RevokedInfo struct {
	ThisUpdate *time.Time
	NextUpdate *time.Time
	Items      []RevokedItem
}

type Crl struct {
	crl *pkix.CertificateList
}

func (s *Crl) LoadFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return s.LoadFromMemory(data)
}

func (s *Crl) LoadFromMemory(data []byte) error {
	crl, err := x509.ParseCRL(data)
	if err != nil {
		return err
	}
	s.crl = crl

	return nil
}

func (s *Crl) VerifySignature(cert *Crt) error {
	if cert == nil {
		return errors.New("证书无效")
	}
	cert509 := cert.Certificate()
	if cert509 == nil {
		return errors.New("证书无效")
	}

	return cert509.CheckCRLSignature(s.crl)
}

func (s *Crl) List() *pkix.TBSCertificateList {
	if s.crl == nil {
		return nil
	}

	return &s.crl.TBSCertList
}

func (s *Crl) Info() *RevokedInfo {
	info := &RevokedInfo{
		Items: make([]RevokedItem, 0),
	}
	list := s.List()
	if list == nil {
		return info
	}
	info.ThisUpdate = &list.ThisUpdate
	info.NextUpdate = &list.NextUpdate

	lst := list.RevokedCertificates
	lstCount := len(lst)
	for lstIndex := 0; lstIndex < lstCount; lstIndex++ {
		item := RevokedItem{
			SerialNumber:   lst[lstIndex].SerialNumber,
			RevocationTime: lst[lstIndex].RevocationTime,
		}

		oid := &Oid{Extensions: lst[lstIndex].Extensions}
		oid.GetValue(OidOrganization, &item.Organization)
		oid.GetValue(OidOrganizationalUnit, &item.OrganizationalUnit)
		oid.GetValue(OidLocality, &item.Locality)
		oid.GetValue(OidNotBefore, &item.NotBefore)
		oid.GetValue(OidNotAfter, &item.NotAfter)

		info.Items = append(info.Items, item)
	}

	return info
}
