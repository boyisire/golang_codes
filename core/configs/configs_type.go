package configs

type KafkaMqConfig struct {
	Topics       []string `json:"topics"`
	Servers      []string `json:"servers"`
	Ak           string   `json:"ak"`
	Password     string   `json:"password"`
	ConsumerId   string   `json:"consumerId"`
	CertFile     string   `json:"cert_file"`
	FullCertFile string   `json:"full_cert_file"`
}
