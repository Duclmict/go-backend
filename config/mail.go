package config

type MailTemplate struct {
	From 		string					`json:"from"`
	Name	 	string					`json:"name"`
	Subject 	string 					`json:"subject"`
	Template 	string 					`json:"template"`
	Option []interface{}				`json:"option"`
}

type CCStruct struct {
	Address 	string					`json:"from"`
	Name	 	string					`json:"name"`
	Option []interface{}				`json:"option"`
}

type AttachStruct struct {
	FileName 	string					`json:"from"`
	Option []interface{}				`json:"option"`
}

// prefix MAIL
var (
	MAIL_DRIVER	string
	MAIL_HOST string
	MAIL_PORT int
	MAIL_USERNAME string
	MAIL_PASSWORD string
	MAIL_ENCRYPTION string
	MAIL_FROM_ADDRESS string
	MAIL_FROM_NAME string
	MAIL_TEMPLATE map[string]*MailTemplate
)

