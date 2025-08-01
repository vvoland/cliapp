package cliapp

import (
	"log/slog"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := Init()
	slog.InfoContext(ctx, "works")
}
