package global

import (
	"github.com/jinzhu/configor"
)

type config struct {
	Proxy struct {
		FallbackServer string `toml:"fallback_server" yaml:"fallback_server"  env:"MISSTODON_FALLBACK_SERVER"`
	} `toml:"proxy" yaml:"proxy"`
	Server struct {
		BindAddress string `toml:"bind_address" yaml:"bind_address" env:"MISSTODON_SERVER_BIND_ADDRESS"`
		AutoTLS     bool   `toml:"auto_tls" yaml:"auto_tls" env:"MISSTODON_SERVER_AUTO_TLS"`
		Domain      string `toml:"domain" yaml:"domain" env:"MISSTODON_SERVER_DOMAIN"`
		TlsCertFile string `toml:"tls_cert_file" yaml:"tls_cert_file" env:"MISSTODON_SERVER_TLS_CERT_FILE"`
		TlsKeyFile  string `toml:"tls_key_file" yaml:"tls_key_file" env:"MISSTODON_SERVER_TLS_KEY_FILE"`
	} `toml:"server" yaml:"server"`
	Logger struct {
		Level         int8   `toml:"level" yaml:"level" env:"MISSTODON_LOGGER_LEVEL"`
		ConsoleWriter bool   `toml:"console_writer" yaml:"console_writer" env:"MISSTODON_LOGGER_CONSOLE_WRITER"`
		RequestLogger bool   `toml:"request_logger" yaml:"request_logger" env:"MISSTODON_LOGGER_REQUEST_LOGGER"`
		Filename      string `toml:"filename" yaml:"filename" env:"MISSTODON_LOGGER_FILENAME"`
		MaxAge        int    `toml:"max_age" yaml:"max_age" env:"MISSTODON_LOGGER_MAX_AGE"`
		MaxBackups    int    `toml:"max_backups" yaml:"max_backups" env:"MISSTODON_LOGGER_MAX_BACKUPS"`
	} `toml:"logger" yaml:"logger"`
}

var Config config

func LoadConfig(filename string) error {
	return configor.
		New(&configor.Config{Environment: "production"}).
		Load(&Config, filename)
}
