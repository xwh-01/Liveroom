package config

import "github.com/BurntSushi/toml"

type Config struct {
	Server ServerConfig `toml:"server"`
	Redis  RedisConfig  `toml:"redis"`
}

type ServerConfig struct {
	Port int `toml:"port"`
}

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

var Cfg *Config

func Load(path string) error {
	Cfg = &Config{}
	_, err := toml.DecodeFile(path, Cfg)
	return err
}
