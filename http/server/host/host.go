package host

import (
	"fmt"
	"net/http"
	"crypto/tls"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/httpserver/http/server/configure"
	"errors"
	"github.com/ktpswjz/httpserver/security/certificate"
	"github.com/ktpswjz/httpserver/http/server/handler"
)

type Host interface {
	Run() error
}

func NewHost(config *configure.Server, handle handler.Handle, log types.Log) Host  {
	instance := &innerHost {config: config, handle: handle}
	instance.SetLog(log)

	return instance
}

type innerHost struct {
	types.Base
	config *configure.Server
	handle handler.Handle
}

func (s *innerHost) Run() error  {
	if s.config == nil {
		return errors.New(s.LogError("invalid configure"))
	}

	router, err := handler.NewHandler(s.handle, s.GetLog())
	if err != nil {
		return err
	}

	ch := make(chan string)
	srvCount := 0
	if s.config.Http.Enabled {
		srvCount++
		go func() {
			addr := fmt.Sprintf(":%s", s.config.Http.Port)
			s.LogInfo("http listening on \"", addr, "\"")
			err := http.ListenAndServe(addr, router)
			if err != nil {
				s.LogError("http", err)
			}
			ch <- "http stopped"
		}()
	}

	if s.config.Https.Enabled {
		srvCount++
		go func() {
			certFilePath := s.config.Https.Cert.File
			s.LogInfo("https cert file: ", certFilePath)
			pfx := &certificate.Pfx{}
			err := pfx.LoadFromFile(certFilePath, s.config.Https.Cert.Password)
			if nil == err {
				addr := fmt.Sprintf(":%s", s.config.Https.Port)
				s.LogInfo("https listening on \"", addr, "\"")

				srv := &http.Server{
					Addr:    addr,
					Handler: router,
					TLSConfig: &tls.Config{
						Certificates: pfx.TlsCertificates(),
						ClientAuth:   tls.NoClientCert,
					},
				}

				err = srv.ListenAndServeTLS("", "")
				if err != nil {
					s.LogError("https", err)
				}
			} else {
				s.LogError("https cert invalid: ", err)
			}

			addr := fmt.Sprintf(":%s", s.config.Https.Port)
			s.LogInfo("https listening on \"", addr, "\"")


			ch <- "https stopped"
		}()
	}

	for srvIndex := 0; srvIndex < srvCount; srvIndex++ {
		s.LogInfo(<-ch)
	}

	s.LogInfo("exited", fmt.Sprintf("server count: %d", srvCount))

	return nil
}