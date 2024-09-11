package satusehat

import "time"

type Credential struct {
	AuthUrl      string `yaml:"auth_url" mapstructure:"auth_url" validate:"required,url"`
	BaseUrl      string `yaml:"base_url" mapstructure:"base_url" validate:"required,url"`
	ClientId     string `yaml:"client_id" mapstructure:"client_id" validate:"required"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret" validate:"required"`
}

type RestConfig struct {
	RetryCount       int           `yaml:"retry_count" mapstructure:"retry_count"`
	RetryWaitTime    time.Duration `yaml:"retry_wait_time" mapstructure:"retry_wait_time"`
	RetryMaxWaitTime time.Duration `yaml:"retry_max_wait_time" mapstructure:"retry_max_wait_time"`
	Timeout          time.Duration `yaml:"timeout" mapstructure:"timeout"`
}

func defaultRestConfig() *RestConfig {
	return &RestConfig{
		RetryCount:       1,
		RetryWaitTime:    2 * time.Second,
		RetryMaxWaitTime: 30 * time.Second,
		Timeout:          5 * time.Second,
	}
}

func defaultCredentials() *Credential {
	return &Credential{
		AuthUrl:      "https://api-satusehat.kemkes.go.id/oauth2/v1",
		BaseUrl:      "https://api-satusehat.kemkes.go.id/fhir-r4/v1",
		ClientId:     "",
		ClientSecret: "",
	}

}
