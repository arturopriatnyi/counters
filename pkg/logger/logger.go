package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Mode string

const (
	ProdMode Mode = "prod"
	DevMode  Mode = "dev"
)

type Option func(cfg zap.Config) zap.Config

func WithMode(mode Mode) Option {
	return func(cfg zap.Config) zap.Config {
		switch mode {
		case ProdMode:
			cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
			cfg.Development = false
		case DevMode:
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			cfg.Development = true
		}

		return cfg
	}
}

func New(opts ...Option) (*zap.Logger, error) {
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "lvl",
			TimeKey:     "ts",
			EncodeTime:  zapcore.RFC3339NanoTimeEncoder,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: make(map[string]interface{}),
	}

	for _, opt := range opts {
		cfg = opt(cfg)
	}

	return cfg.Build()
}
