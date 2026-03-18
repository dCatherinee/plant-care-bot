package config

import (
	"errors"
	"strings"
	"testing"
)

func validConfig() Config {
	return Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "plants",
	}
}

func TestLoad(t *testing.T) {
	t.Setenv("DB_HOST", " localhost ")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", " postgres ")
	t.Setenv("DB_PASSWORD", " postgres ")
	t.Setenv("DB_NAME", " plants ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DBHost != "localhost" {
		t.Fatalf("expected DBHost %q, got %q", "localhost", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Fatalf("expected DBPort %q, got %q", "5432", cfg.DBPort)
	}
	if cfg.DBUser != "postgres" {
		t.Fatalf("expected DBUser %q, got %q", "postgres", cfg.DBUser)
	}
	if cfg.DBPassword != "postgres" {
		t.Fatalf("expected DBPassword %q, got %q", "postgres", cfg.DBPassword)
	}
	if cfg.DBName != "plants" {
		t.Fatalf("expected DBName %q, got %q", "plants", cfg.DBName)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name         string
		cfg          Config
		wantErr      bool
		wantVariable string
		wantProblem  string
		wantParts    []string
		checkAs      bool
	}{
		{
			name:    "valid",
			cfg:     validConfig(),
			wantErr: false,
		},
		{
			name: "missing port",
			cfg: func() Config {
				cfg := validConfig()
				cfg.DBPort = ""
				return cfg
			}(),
			wantErr:      true,
			wantVariable: "DB_PORT",
			wantProblem:  "is required",
			wantParts:    []string{"DB_PORT", "is required"},
			checkAs:      true,
		},
		{
			name: "invalid port format",
			cfg: func() Config {
				cfg := validConfig()
				cfg.DBPort = "abc"
				return cfg
			}(),
			wantErr:      true,
			wantVariable: "DB_PORT",
			wantProblem:  "must be a valid integer",
			wantParts:    []string{"DB_PORT", "must be a valid integer"},
			checkAs:      true,
		},
		{
			name: "port out of range",
			cfg: func() Config {
				cfg := validConfig()
				cfg.DBPort = "70000"
				return cfg
			}(),
			wantErr:      true,
			wantVariable: "DB_PORT",
			wantProblem:  "must be between 1 and 65535",
			wantParts:    []string{"DB_PORT", "must be between 1 and 65535"},
			checkAs:      true,
		},
		{
			name: "multiple errors",
			cfg: func() Config {
				cfg := validConfig()
				cfg.DBHost = ""
				cfg.DBPort = "70000"
				cfg.DBUser = ""
				cfg.DBPassword = ""
				cfg.DBName = ""
				return cfg
			}(),
			wantErr: true,
			wantParts: []string{
				"DB_HOST",
				"DB_PORT",
				"DB_USER",
				"DB_PASSWORD",
				"DB_NAME",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()

			if !tc.wantErr {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("expected ErrInvalidConfig, got %v", err)
			}

			if tc.checkAs {
				var validationErr ValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("expected ValidationError, got %T: %v", err, err)
				}

				if validationErr.Variable != tc.wantVariable {
					t.Fatalf("expected variable %q, got %q", tc.wantVariable, validationErr.Variable)
				}
				if validationErr.Problem != tc.wantProblem {
					t.Fatalf("expected problem %q, got %q", tc.wantProblem, validationErr.Problem)
				}
			}

			for _, part := range tc.wantParts {
				if !strings.Contains(err.Error(), part) {
					t.Fatalf("expected error %q to contain %q", err.Error(), part)
				}
			}
		})
	}
}
