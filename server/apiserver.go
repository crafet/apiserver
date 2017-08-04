// apiserver project apiserver.go
package server

import (
	"apiserver/config"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiServer struct {
	Config *config.Configuraton

	Client *http.Client
}

func (s *ApiServer) InitClient() error {

	sec := s.Config.ClientTimeout
	if sec < 0 || sec > 100 {
		// default timeout 5s
		sec = 5
	}

	if !s.Config.UseProxy {
		s.Client = &http.Client{Timeout: time.Second * time.Duration(sec)}
		return nil
	}

	// init client with proxy
	var fullUrl string

	if !strings.HasPrefix(s.Config.ProxyAddr, "http://") {
		fullUrl = "http://" + s.Config.ProxyAddr + ":" + s.Config.ProxyPort
	} else {
		fullUrl = s.Config.ProxyAddr + ":" + s.Config.ProxyPort
	}

	log.Println("fullUrl: ", fullUrl)

	proxyUrl, err := url.Parse(fullUrl)
	if err != nil {
		log.Println("parse proxy url failed with proxyUrl:", proxyUrl)
		return err
	}

	tr := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	s.Client = &http.Client{Transport: tr, Timeout: time.Second * time.Duration(sec)}
	if s.Client == nil {
		return errors.New("faild to build client with proxy")
	}

	// build client successfully
	return nil
}

func (s *ApiServer) InitConfig(conf string) error {

	c := config.NewConfig(conf)
	err := c.LoadConfig()
	if err != nil {
		return err
	}

	s.Config = &c.C
	return nil
}

func (s *ApiServer) Init(conf string) error {

	// init config
	err := s.InitConfig(conf)
	if err != nil {
		log.Println("failed to init conf")
		return err
	}

	// init client
	if err := s.InitClient(); err != nil {
		log.Println("failed to init client")
		return err
	}

	return nil
}
