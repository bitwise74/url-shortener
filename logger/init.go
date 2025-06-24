// Package logger includes stuff related to logging data
package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Init initializes the zap logger
func Init(dev bool) (*zap.Logger, error) {
	if dev {
		l, err := zap.NewDevelopment(zap.AddCaller())
		if err != nil {
			return nil, fmt.Errorf("failed to initialize development logger: %w", err)
		}
		defer l.Sync()
		return l, nil
	}

	l, err := zap.NewProduction(zap.AddCaller())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize production logger: %w", err)
	}

	return l, nil
}
