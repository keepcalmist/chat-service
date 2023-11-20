package config

import "time"

type Config struct {
	Global   GlobalConfig  `toml:"global"`
	Log      LogConfig     `toml:"log"`
	Servers  ServersConfig `toml:"servers"`
	Sentry   Sentry        `toml:"sentry"`
	Clients  Clients       `toml:"clients"`
	Postgres Postgres      `toml:"postgres"`
	Services Services      `toml:"services"`
}

type Services struct {
	MsgProducer MsgProducer `toml:"msg_producer"`
	Outbox      Outbox      `toml:"outbox"`
}

type MsgProducer struct {
	Brokers    []string `toml:"brokers" validate:"required,dive,hostname_port"`
	Topic      string   `toml:"topic" validate:"required"`
	BatchSize  int      `toml:"batch_size" validate:"required"`
	EncryptKey string   `toml:"encrypt_key"`
}

type Outbox struct {
	Workers    int           `toml:"workers" validate:"required"`
	IdleTime   time.Duration `toml:"idle_time" validate:"required"`
	ReserveFor time.Duration `toml:"reserve_for" validate:"required"`
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
	Addr           string         `toml:"addr" validate:"required,hostname_port"`
	AllowOrigins   []string       `toml:"allow_origins" validate:"required,dive,url"`
	RequiredAccess RequiredAccess `toml:"required_access" validate:"required,dive"`
}

type RequiredAccess struct {
	Resource string `toml:"resource" validate:"required"`
	Role     string `toml:"role" validate:"required"`
}

type DebugServerConfig struct {
	Addr string `toml:"addr" validate:"required,hostname_port"`
}

type Sentry struct {
	DSN string `toml:"dsn" validate:"omitempty,url"`
}

type Clients struct {
	Keycloak Keycloak `toml:"keycloak"`
}

type Keycloak struct {
	BasePath     string `toml:"base_path" validate:"required,url"`
	Realm        string `toml:"realm" validate:"required"`
	ClientID     string `toml:"client_id" validate:"required"`
	ClientSecret string `toml:"client_secret" validate:"required"`
	DebugMode    bool   `toml:"debug_mode"`
}

type Postgres struct {
	Address  string `toml:"address" validate:"required,hostname_port"`
	Username string `toml:"username" validate:"required"`
	Password string `toml:"password" validate:"required"`
	Database string `toml:"database" validate:"required"`
	Debug    bool   `toml:"debug" validate:""`
}

func (c GlobalConfig) IsProduction() bool {
	return c.Env == "prod"
}
