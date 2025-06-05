package pkg

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host         string
		Port         int
		Name         string
		Username     string
		Password     string
		SSL          string
		PoolMaxConns int
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

	Logging struct {
		Level  string
		Format string
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
	allowedSSL := map[string]bool{
		"":            true,
		"disable":     true,
		"allow":       true,
		"prefer":      true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if !allowedSSL[cfg.Database.SSL] {
		return fmt.Errorf("invalid database SSL mode (GOKB_DATABASE_SSL)")
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

	loggingLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !loggingLevels[cfg.Logging.Level] {
		return fmt.Errorf("invalid logging level (GOKB_LOGGING_LEVEL)")
	}
	if cfg.Logging.Format != "text" && cfg.Logging.Format != "json" {
		return fmt.Errorf("invalid logging format (GOKB_LOGGING_FORMAT)")
	}

	return nil
}

func InitConfig() (*Config, error) {
	// Read environment
	viper.AutomaticEnv()
	// In case if configuration file is needed
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

	cfg := &Config{}

	// Database configuration
	cfg.Database.Host = viper.GetString("GOKB_DATABASE_HOST")
	cfg.Database.Port = viper.GetInt("GOKB_DATABASE_PORT")
	cfg.Database.Name = viper.GetString("GOKB_DATABASE_NAME")
	cfg.Database.Username = viper.GetString("GOKB_DATABASE_USERNAME")
	cfg.Database.Password = viper.GetString("GOKB_DATABASE_PASSWORD")
	cfg.Database.SSL = viper.GetString("GOKB_DATABASE_SSL")
	if cfg.Database.SSL == "" {
		cfg.Database.SSL = "disable"
	}
	cfg.Database.PoolMaxConns = viper.GetInt("GOKB_DATABASE_POOL_MAX_CONNS")
	if cfg.Database.PoolMaxConns == 0 {
		cfg.Database.PoolMaxConns = 10
	}

	// Root user configuration
	cfg.Root.Email = viper.GetString("GOKB_ROOT_EMAIL")
	cfg.Root.Username = viper.GetString("GOKB_ROOT_USERNAME")
	if cfg.Root.Username == "" {
		cfg.Root.Username = "gokb"
	}

	// Cache server configuration
	cfg.Cache.URL = viper.GetString("GOKB_CACHE_URL")
	cfg.Cache.Port = viper.GetInt("GOKB_CACHE_PORT")
	cfg.Cache.Secure = viper.GetBool("GOKB_CACHE_SECURE")
	if cfg.Cache.Secure {
		cfg.Cache.Password = viper.GetString("GOKB_CACHE_PASSWORD")
	}

	// Logging configuration
	cfg.Logging.Level = viper.GetString("GOKB_LOGGING_LEVEL")
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	cfg.Logging.Format = viper.GetString("GOKB_LOGGING_FORMAT")
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "text"
	}

	// Configuration validation
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
