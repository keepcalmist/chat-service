package config

type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Log     LogConfig     `toml:"log"`
	Servers ServersConfig `toml:"servers"`
	Sentry  Sentry        `toml:"sentry"`
}

type GlobalConfig struct {
	Env string `toml:"env" validate:"oneof=dev stage prod"`
}

type LogConfig struct {
	Level string `toml:"level" validate:"oneof=debug info warn error"`
}

type ServersConfig struct {
	Debug  DebugServerConfig  `toml:"debug"`
	Client ClientServerConfig `toml:"client"`
}

type ClientServerConfig struct {
	Addr         string   `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins []string `toml:"allow_origins" validate:"required,dive,url"`
	SwaggerFile  string   `toml:"swagger_file" validate:"required,filepath"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type Sentry struct {
	DSN string `toml:"dsn" validate:"sentrydsn"`
}

func (c GlobalConfig) IsProduction() bool {
	return c.Env == "prod"
}
