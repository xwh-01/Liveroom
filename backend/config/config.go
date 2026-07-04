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
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	Username     string `toml:"username"`
	Password     string `toml:"password"`
	Database     string `toml:"database"`
	Charset      string `toml:"charset"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxOpenConns int    `toml:"max_open_conns"`
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
