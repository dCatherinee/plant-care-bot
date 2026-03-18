package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var ErrInvalidConfig = errors.New("invalid config")

type ValidationError struct {
	Variable string
	Problem  string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("config validation: %s: %s", e.Variable, e.Problem)
}

func (e ValidationError) Unwrap() error {
	return ErrInvalidConfig
}

func Load() (Config, error) {
	cfg := Config{
		DBHost:     strings.TrimSpace(os.Getenv("DB_HOST")),
		DBPort:     strings.TrimSpace(os.Getenv("DB_PORT")),
		DBUser:     strings.TrimSpace(os.Getenv("DB_USER")),
		DBPassword: strings.TrimSpace(os.Getenv("DB_PASSWORD")),
		DBName:     strings.TrimSpace(os.Getenv("DB_NAME")),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var errs []error

	if c.DBHost == "" {
		errs = append(errs, ValidationError{
			Variable: "DB_HOST",
			Problem:  "is required",
		})
	}

	if c.DBPort == "" {
		errs = append(errs, ValidationError{
			Variable: "DB_PORT",
			Problem:  "is required",
		})
	} else {
		port, err := strconv.Atoi(c.DBPort)
		if err != nil {
			errs = append(errs, ValidationError{
				Variable: "DB_PORT",
				Problem:  "must be a valid integer",
			})
		} else if port < 1 || port > 65535 {
			errs = append(errs, ValidationError{
				Variable: "DB_PORT",
				Problem:  "must be between 1 and 65535",
			})
		}
	}

	if c.DBUser == "" {
		errs = append(errs, ValidationError{
			Variable: "DB_USER",
			Problem:  "is required",
		})
	}

	if c.DBPassword == "" {
		errs = append(errs, ValidationError{
			Variable: "DB_PASSWORD",
			Problem:  "is required",
		})
	}

	if c.DBName == "" {
		errs = append(errs, ValidationError{
			Variable: "DB_NAME",
			Problem:  "is required",
		})
	}

	return errors.Join(errs...)
}
