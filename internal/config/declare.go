package config

import "time"

type Config struct {
	Server   Server   `mapstructure:"server"`
	MongoDB  Mongo    `mapstructure:"mongodb"`
	Redis    Redis    `mapstructure:"redis"`
	Firebase Firebase `mapstructure:"firebase"`
	JWT      JWT      `mapstructure:"jwt"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Mongo struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Firebase struct {
	ProjectID      string `mapstructure:"project_id"`
	CredentialFile string `mapstructure:"credential_file"`
}

type JWT struct {
	SecretKey     string        `mapstructure:"secret_key"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}
