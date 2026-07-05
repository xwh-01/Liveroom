package config

import "github.com/BurntSushi/toml"

type Config struct {
	Server ServerConfig `toml:"server"`
	Redis  RedisConfig  `toml:"redis"`
	MySQL  MySQLConfig  `toml:"mysql"`
	Room   RoomConfig   `toml:"room"`
}

type ServerConfig struct {
	Port int `toml:"port"`
}

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type MySQLConfig struct {
	DSN                    string `toml:"dsn"`
	MaxOpenConns           int    `toml:"max_open_conns"`
	MaxIdleConns           int    `toml:"max_idle_conns"`
	ConnMaxLifetimeSeconds int    `toml:"conn_max_lifetime_seconds"`
}

type RoomConfig struct {
	DefaultID string `toml:"default_id"`
}

var Cfg *Config

func Load(path string) error {
	Cfg = &Config{}
	_, err := toml.DecodeFile(path, Cfg)
	return err
}
