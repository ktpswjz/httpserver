package client

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"github.com/ktpswjz/httpserver/security/rsakey"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	Transport *http.Transport // usually for https request
	Timeout   int64           // timeout in seconds unit, zero meas not timeout
}

func (s *Client) Get(url string, argument interface{}) ([]byte, []byte, *tls.ConnectionState, int, error) {
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
		return nil, nil, nil, 0, err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, resp.StatusCode, err
	}

	return nil, bodyData, resp.TLS, resp.StatusCode, nil
}

func (s *Client) PostJson(url string, argument interface{}, headers ...Header) ([]byte, []byte, *tls.ConnectionState, int, error) {
	var input []byte = nil
	var body io.Reader = nil
	if argument != nil {
		switch argument.(type) {
		case []byte:
			body = bytes.NewBuffer(argument.([]byte))
			input = argument.([]byte)
		default:
			bodyData, err := json.Marshal(argument)
			if err != nil {
				return bodyData, nil, nil, 0, err
			}
			body = bytes.NewBuffer([]byte(bodyData))
			input = bodyData
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

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return input, nil, nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	headerCount := len(headers)
	for i := 0; i < headerCount; i++ {
		header := headers[i]
		req.Header.Add(header.Key, header.Value)
	}
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return input, nil, nil, 0, err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return input, nil, nil, resp.StatusCode, err
	}

	return input, bodyData, resp.TLS, resp.StatusCode, nil
}

func (s *Client) PostXml(url string, argument interface{}) ([]byte, []byte, *tls.ConnectionState, int, error) {
	var input []byte = nil
	var body io.Reader = nil
	if argument != nil {
		switch argument.(type) {
		case []byte:
			body = bytes.NewBuffer(argument.([]byte))
			input = argument.([]byte)
		case string:
			body = bytes.NewBufferString(argument.(string))
			input = []byte(argument.(string))
		default:
			bodyData, err := xml.MarshalIndent(argument, "", "	")
			if err != nil {
				return bodyData, nil, nil, 0, err
			}
			body = bytes.NewBuffer([]byte(bodyData))
			input = bodyData
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

	resp, err := client.Post(url, "application/xml;charset=utf-8", body)
	if err != nil {
		return input, nil, nil, 0, err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return input, nil, nil, resp.StatusCode, err
	}

	return input, bodyData, resp.TLS, resp.StatusCode, nil
}

func (s *Client) PostSoap(url string, argument interface{}) ([]byte, []byte, *tls.ConnectionState, int, error) {
	soap := &Soap{
		Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:    "http://www.w3.org/2001/XMLSchema",
		Soap12: "http://www.w3.org/2003/05/soap-envelope",
		Body: SoapBody{
			Data: argument,
		},
	}

	input, err := xml.MarshalIndent(soap, "", "	")
	if err != nil {
		return input, nil, nil, 0, err
	}
	body := bytes.NewBuffer([]byte(input))

	client := &http.Client{}
	if s.Transport != nil {
		client.Transport = s.Transport
	}
	if s.Timeout > 0 {
		timeout := s.Timeout * time.Second.Nanoseconds()
		client.Timeout = time.Duration(timeout)
	}

	resp, err := client.Post(url, "application/soap+xml;charset=utf-8", body)
	if err != nil {
		return input, nil, nil, 0, err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return input, nil, nil, resp.StatusCode, err
	}

	return input, bodyData, resp.TLS, resp.StatusCode, nil
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
