package config

type Config struct {
	Server   Server   `mapstructure:"server"`
	MongoDB  Mongo    `mapstructure:"mongodb"`
	Redis    Redis    `mapstructure:"redis"`
	Firebase Firebase `mapstructure:"firebase"`
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
