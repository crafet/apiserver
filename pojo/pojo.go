// mail info
package pojo

type CommonResponse struct {
	Result     bool   `json:"result"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type TemplateInfo struct {
	Name           string `json:"name"`
	InvokeName     string `json:"invokeName"`
	TemplateType   int    `json:"templateType"`
	TemplateStat   int    `json:"templateStat"`
	CreateTime     string `json:"gmtCreated"`
	UpdateTime     string `json:"gmtModified"`
	Subject        string `json:"subject"`
	ContentSummary string `json:"contentSummary"`
}

type DataList struct {
	DataList []TemplateInfo `json:"dataList"`
}

// info is a json object, not a string field
type GetTemplateListResponse struct {
	CommonResponse
	Info DataList `json:"info"`
}

type GetOneTemplateData struct {
	TemplateInfo
	Html string `json:"html"`
}

type GetOneTemplateInfo struct {
	Data GetOneTemplateData `json:"data"`
}

// get one template
type GetOneTemplateResponse struct {
	CommonResponse
	Info GetOneTemplateInfo `json:"info"`
}

// send one mail request
type SendMailRequest struct {
	ApiKey  string
	ApiUser string

	// 发件人地址
	// required
	From string

	// 收件人地址. 多个地址使用';'分隔, 如 ben@ifaxin.com;joe@ifaxin.com
	// 如果使用了to，那么收件人会显示所有人
	// 但是如果使用下面的Xsmtpapi字段，那么就会单独发送
	To string

	// 	标题. 不能为空
	Subject string

	// 邮件的内容. 邮件格式为 text/html
	Html string

	// 邮件摘要
	ContentSummary string

	// 发件人名称 显示如: ifaxin客服支持<support@ifaxin.com>
	// 其中"ifaxin客服支持"为FromName
	// optional
	FromName string

	// optional
	Cc string

	// optional
	Bcc string

	// optional
	// 设置用户默认的回复邮件地址.多个地址使用';'分隔
	Replyto string

	// 本次发送所使用的标签ID. 此标签需要事先创建
	// optional
	LabelId int

	// 邮件头部信息. JSON 格式, 比如:{"header1": "value1", "header2": "value2"}
	// optional
	Headers string

	// 邮件附件. 发送附件时, 必须使用 multipart/form-data 进行 post 提交 (表单提交)
	// optional
	Attachments []interface{}

	// SMTP 扩展字段
	Xsmtpapi string

	// optional 邮件的内容. 邮件格式为 text/plain
	Plain string

	// 默认值: true
	// optional
	RespEmailId string

	// 默认值: false. 是否使用回执
	UseNotification string

	// 默认值: false. 是否使用地址列表发送.
	// 比如: to=group1@maillist.sendcloud.org;group2@maillist.sendcloud.org
	// 如果此字段为true，那么上面的To字段包含的就是地址列表，不是普通的email邮箱
	// 此时Cc,Bcc字段失效

	UseAddressList string

	// 如果UseAddressList为false,同时指定了Xsmtpapi扩展，那么使用Xsmtpapi中的to
	// 此时上面的To CC Bcc三个字段也不再用
}

type SendMailResponseInfo struct {
	EmailIdList []string `json:"emailIdList"`
}

// extend smtp valud
type Xsmtpapi struct {
	To  []string            `json:"to"`
	Sub map[string][]string `json:"sub"`
}

type SendMailResponse struct {
	CommonResponse
	Info SendMailResponseInfo `json:"info"`
}

// if error request
type ErrorSendMailRequest struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
