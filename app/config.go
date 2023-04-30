package app

import (
	"encoding/json"
	"os"
)

type Config struct {
	SMTPHost       string `json:"smtp_host"`
	SMTPPort       int    `json:"smtp_port"`
	SMTPUsername   string `json:"smtp_username"`
	SMTPSender     string `json:"smtp_sender"`
	SMTPPassword   string `json:"smtp_password"`
	AdminEmailAddr string `json:"admin_email_addr"`
}

func MustLoadConfig(fpath string) *Config {
	f, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	config := &Config{}
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		panic(err)
	}
	return config
}
