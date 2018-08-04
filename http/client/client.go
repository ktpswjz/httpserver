package client

import (
	"net/http"
	"crypto/tls"
	"io"
	"bytes"
	"encoding/json"
	"time"
	"io/ioutil"
	"github.com/ktpswjz/httpserver/security/rsakey"
	"crypto/rsa"
)

type Client struct {
	Transport	*http.Transport		// usually for https request
	Timeout		int64				// timeout in seconds unit, zero meas not timeout
}

func (s *Client) PostJson(url string, argument interface{}) ([]byte, *tls.ConnectionState, error)  {
	var body io.Reader = nil
	if argument != nil {
		switch argument.(type) {
		case []byte:
			body = bytes.NewBuffer(argument.([]byte))
		default:
			bodyData, err := json.Marshal(argument)
			if err != nil {
				return nil, nil, err
			}
			body = bytes.NewBuffer([]byte(bodyData))
		}
	}

	client := &http.Client{}
	if s.Transport != nil {
		client.Transport = s.Transport
	}
	if s.Timeout > 0 {
		timeout := s.Timeout * time.Second.Nanoseconds()
		client.Timeout = time.Duration(timeout)
	}

	resp, err := client.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return bodyData, resp.TLS, nil
}

func (s *Client) Download(url string, argument interface{}) ([]byte, *tls.ConnectionState, error) {
	client := &http.Client{}
	if s.Transport != nil {
		client.Transport = s.Transport
	}
	if s.Timeout > 0 {
		timeout := s.Timeout * time.Second.Nanoseconds()
		client.Timeout = time.Duration(timeout)
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return bodyData, resp.TLS, nil
}

func (s *Client) PrivateKey() *rsakey.Private {
	if s.Transport == nil {
		return nil
	}

	cfg := s.Transport.TLSClientConfig
	if cfg == nil {
		return nil
	}

	if len(cfg.Certificates) == 0 {
		return nil
	}

	cert := cfg.Certificates[0]
	key, ok := cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil
	}

	return &rsakey.Private{Key: key}
}