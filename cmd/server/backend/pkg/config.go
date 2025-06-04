package pkg

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		Name     string
		Username string
		Password string
	}

	Root struct {
		Email    string
		Username string
	}

	Cache struct {
		URL      string
		Port     int
		Secure   bool
		Password string
	}
}

func validateConfig(cfg *Config) error {
	// Database config validation
	if cfg.Database.Host == "" {
		return fmt.Errorf("missing required database config (GOKB_DATABASE_HOST)")
	}
	if cfg.Database.Port == 0 {
		return fmt.Errorf("missing required database config (GOKB_DATABASE_PORT)")
	}
	if cfg.Database.Name == "" {
		return fmt.Errorf("missing required database config (GOKB_DATABASE_NAME)")
	}
	if cfg.Database.Username == "" {
		return fmt.Errorf("missing required database config (GOKB_DATABASE_USERNAME)")
	}
	if cfg.Database.Password == "" {
		return fmt.Errorf("missing required database config (GOKB_DATABASE_PASSWORD)")
	}

	// Root user validation
	if cfg.Root.Email == "" {
		return fmt.Errorf("missing required root user config (GOKB_ROOT_EMAIL)")
	}
	if cfg.Root.Username == "" {
		return fmt.Errorf("missing required root user config (GOKB_ROOT_USERNAME)")
	}

	// Cache server config validation
	if cfg.Cache.URL == "" {
		return fmt.Errorf("missing required cache server config (GOKB_CACHE_URL)")
	}
	if cfg.Cache.Port == 0 {
		return fmt.Errorf("missing required cache server config (GOKB_CACHE_PORT)")
	}
	if cfg.Cache.Secure && cfg.Cache.Password == "" {
		return fmt.Errorf("missing required cache server config (GOKB_CACHE_PASSWORD)")
	}

	return nil
}

func InitConfig() (*Config, error) {
	viper.AutomaticEnv()

	cfg := &Config{}

	// Database configuration
	cfg.Database.Host = viper.GetString("GOKB_DATABASE_HOST")
	cfg.Database.Port = viper.GetInt("GOKB_DATABASE_PORT")
	cfg.Database.Name = viper.GetString("GOKB_DATABASE_NAME")
	cfg.Database.Username = viper.GetString("GOKB_DATABASE_USERNAME")
	cfg.Database.Password = viper.GetString("GOKB_DATABASE_PASSWORD")

	// Root user configuration
	cfg.Root.Email = viper.GetString("GOKB_ROOT_EMAIL")
	cfg.Root.Username = viper.GetString("GOKB_ROOT_USERNAME")

	// Cache server configuration
	cfg.Cache.URL = viper.GetString("GOKB_CACHE_URL")
	cfg.Cache.Port = viper.GetInt("GOKB_CACHE_PORT")
	cfg.Cache.Secure = viper.GetBool("GOKB_CACHE_SECURE")
	if cfg.Cache.Secure {
		cfg.Cache.Password = viper.GetString("GOKB_CACHE_PASSWORD")
	}

	// Configuration validation
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
