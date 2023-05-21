package application

import (
	"os"
)

func Exit() {
	Logger.Info("Shutting down Slixx satellite")
	os.Exit(0)
}
