package application

import "os"

func Exit() {
	Logger.Info("Shutting down Slixx supervisor")
	os.Exit(0)
}
