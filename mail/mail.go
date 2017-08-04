// mail project mail.go
package mail

import (
	"apiserver/config"
	. "apiserver/pojo"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type MailService struct {
	Client *http.Client

	Config *config.Configuraton
}

// get template api
func (m *MailService) GetMailTemplateList(req *http.Request) (*interface{}, error) {

	var get_template_api = m.Config.ApiGetTemplateList

	req.ParseForm()
	log.Println("req from:", req.Method, req.Form)

	// get invokeName of template if query this template
	template_name := req.Form.Get("invokeName")

	req, err := http.NewRequest("GET", get_template_api, nil)
	if err != nil {
		log.Println("failed to new request")
		return nil, errors.New("failed to new http request")
	}

	q := req.URL.Query()
	q.Add("apiKey", m.Config.ApiKey)
	q.Add("apiUser", m.Config.ApiUser)
	if len(template_name) > 0 {
		q.Add("invokeName", template_name)
	}

	req.URL.RawQuery = q.Encode()
	resp, err := m.Client.Do(req)

	if err != nil {
		log.Println("failed to get url, reason:", err.Error())
		return nil, errors.New("failed to execute client do to get mail template")
	}

	defer resp.Body.Close()

	var r interface{}
	r = GetTemplateListResponse{}
	if len(template_name) > 0 {
		r = GetOneTemplateResponse{}
	}
	e := json.NewDecoder(resp.Body).Decode(&r)
	if e != nil {
		log.Println("failed to decode response data")
		return nil, errors.New("failed to decode data")
	}

	// return local variable address is safe in go
	// since golang would do escape analysis
	return &r, nil
}

func (m *MailService) SendMail(req *http.Request) (*interface{}, error) {

	req.ParseForm()
	for k, v := range req.Form {
		log.Println("k:", k, "v:", v)
	}

	// if use mail template
	var use_template bool = false
	mail_template := req.Form["templateInvokeName"]
	if len(mail_template) > 0 {
		use_template = true
	}

	var mail_api = m.Config.ApiSendMail
	if use_template {
		mail_api = m.Config.ApiSendTemplateMail
	}

	// check from field
	from := req.Form.Get("from")
	if len(from) <= 0 {
		var result interface{}
		result = ErrorSendMailRequest{Message: "need from filed", StatusCode: 501}
		return &result, errors.New("failed to get from field from request")
	}

	// to or xsmtpapi should exist one field
	to := req.Form.Get("to")
	xsmtpapi := req.Form.Get("xsmtpapi")

	if len(to) <= 0 && len(xsmtpapi) <= 0 {
		var result interface{}
		result = ErrorSendMailRequest{Message: "need to field or xsmtpapi field", StatusCode: 502}
		return &result, errors.New("failed to get to field or xsmtpapi field")
	}

	// check subject
	subject := req.Form.Get("subject")
	if len(subject) <= 0 {
		var result interface{}
		result = ErrorSendMailRequest{Message: "need subject field", StatusCode: 503}
		return &result, errors.New("failed to get subject field")
	}

	var future_req http.Request
	future_req.ParseForm()

	form := future_req.Form
	form.Add("apiKey", m.Config.ApiKey)
	form.Add("apiUser", m.Config.ApiUser)

	form.Add("from", from)
	form.Add("subject", subject)
	if len(to) > 0 {
		form.Add("to", to)
	}

	if len(xsmtpapi) > 0 {
		form.Add("xsmtpapi", xsmtpapi)
	}

	// if use template, then add template name
	if use_template {

		template := req.Form.Get("templateInvokeName")
		if len(template) <= 0 {
			// should not happen
			var result interface{}
			result = ErrorSendMailRequest{Message: "use template but not find templateInvoKeName", StatusCode: 500}
			log.Println("use mail template but failed to templateInvokeName")
			return &result, errors.New("use mail template but failed to templateInvokeName")
		}

		form.Add("templateInvokeName", template)

	} else {
		// if not template, we use html
		// yet, plain could be used.
		html := req.Form.Get("html")
		if len(html) <= 0 {
			log.Println("failed to Get html, mail need body")
			html = "default mail body"
		}

		form.Add("html", html)
	}

	post_req, err := http.NewRequest("POST", mail_api, strings.NewReader(strings.TrimSpace(form.Encode())))
	post_req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post_req.Header.Set("Connection", "Keep-Alive")

	resp, err := m.Client.Do(post_req)
	if err != nil {
		log.Println("failed to do execute post request")
		return nil, errors.New("failed to execute client do send mail post request")
	}
	defer resp.Body.Close()

	if resp == nil {
		log.Println("resp is nil")
		return nil, errors.New("resp is nil for post send mail")
	}

	var v interface{}
	v = SendMailResponse{}
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		log.Println("failed to decode resp.body")
		return nil, errors.New("faild to decode reponse for post send mail")
	}

	return &v, nil
}
