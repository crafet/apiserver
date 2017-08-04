// config project config.go
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Configuraton struct {
	ApiUser string `json:"api_user"`
	ApiKey  string `json:"api_key"`

	ApiGetTemplateList  string `json:"api_get_template_list"`
	ApiSendMail         string `json:"api_send_mail"`
	ApiSendTemplateMail string `json:"api_send_template_mail"`

	UseProxy      bool   `json:"use_proxy, omitempty"`
	ProxyAddr     string `json:"proxy_addr"`
	ProxyPort     string `json:"proxy_port"`
	ClientTimeout int32  `json:"client_timeout"`

	ServerPort string `json:"server_port"`
}

type Config struct {
	// config file
	File string
	C    Configuraton
}

func NewConfig(file string) *Config {

	return &Config{File: file}
}

// load the config file, which is json type.
// return the load item count and error if error
func (c *Config) LoadConfig() (err error) {

	f, err := os.Open(c.File)
	defer f.Close()

	if err != nil {
		log.Println("failed to open config file: ", c.File)
		return err

	}

	decoder := json.NewDecoder(f)

	if err := decoder.Decode(&(c.C)); err != nil {
		log.Println("decode failed")
		return err
	}

	log.Println("success to load config, value:", c.C)
	return nil
}

func Load() {
	fmt.Println("start to load config")

	DefaultConfig := NewConfig("E:\\go_workspace\\bin\\config.json")
	DefaultConfig.LoadConfig()

}
