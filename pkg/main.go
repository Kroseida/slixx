package main

import (
	"kroseida.org/slixx/pkg/storage"
	"kroseida.org/slixx/pkg/strategy"
)

// This is just a sample.
func main() {
	origin := &storage.FtpKind{}
	destination := []storage.Kind{&storage.FtpKind{}}

	origin.Initialize(&storage.FtpKindConfiguration{
		Host:     "10.30.100.102:21",
		Username: "test",
		Password: "123123123",
		File:     "/test",
	})
	destination[0].Initialize(&storage.FtpKindConfiguration{
		Host:     "10.30.100.102:21",
		Username: "test",
		Password: "123123123",
		File:     "/backup",
	})

	targetStrategy := strategy.COPY
	err := targetStrategy.Execute(origin, destination)
	if err != nil {
		panic(err)
	}
}
