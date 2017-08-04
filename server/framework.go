package server

import (
	"apiserver/mail"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
)

var conf string

// load config
func init() {
	flag.StringVar(&conf, "c", "apiserver.conf", "-c apiserver.conf")
}

type HTTPFramework struct {
	Version string
	Server  *ApiServer

	MailService mail.MailService
}

type NoResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statucCode"`
}

type HandlerFunction func(w http.ResponseWriter, r *http.Request)

type HandlerFunctionMap map[string]HandlerFunction

func (f *HTTPFramework) RegisterHandler() HandlerFunctionMap {

	m := make(HandlerFunctionMap)

	m["/"] = Root
	m["/get_template"] = f.GetTemplate

	m["/send_mail"] = f.SendMail
	return m
}

func (f *HTTPFramework) Init() error {
	f.Version = "0.1"

	f.Server = &ApiServer{}

	// flags.Parse() before use the config
	flag.Parse()

	err := f.Server.Init(conf)
	if err != nil {
		log.Println("failed to init framework server")
		return err
	}

	// init services
	f.MailService = mail.MailService{Config: f.Server.Config, Client: f.Server.Client}

	return nil
}

func (f *HTTPFramework) Run() {

	mux := http.NewServeMux()

	m := f.RegisterHandler()
	for k, v := range m {
		mux.HandleFunc(k, v)
	}

	var port string = f.Server.Config.ServerPort
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Println("server run at port args:", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Println("failed to Run Server")
	}
}

// add more handler
func Root(w http.ResponseWriter, r *http.Request) {
	s := []byte("{\"message\":\"default router\"}")
	w.Write(s)
}

func (f *HTTPFramework) GetTemplate(w http.ResponseWriter, r *http.Request) {
	resp, err := f.MailService.GetMailTemplateList(r)
	if err != nil {
		log.Println("failed to get template list")
	}

	if resp == nil {
		nr, err := json.Marshal(NoResponse{Message: "no response", StatusCode: 500})
		if err != nil {
			log.Println("failed to marshal")
		}
		w.Write(nr)
		return

	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println("failed to marshal")
	}

	w.Write(b)
}

// send mail
func (f *HTTPFramework) SendMail(w http.ResponseWriter, r *http.Request) {

	resp, err := f.MailService.SendMail(r)
	if err != nil {
		log.Println("failed to get resp for SendMail")
	}

	if resp == nil {
		nr, err := json.Marshal(NoResponse{Message: "no response", StatusCode: 500})
		if err != nil {
			log.Println("failed to marshal resp")
		}

		w.Write(nr)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println("failed to marshal")
	}

	w.Write(b)
}
