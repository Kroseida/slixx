package main

import (
	_storage "kroseida.org/slixx/pkg/storage"
	_strategy "kroseida.org/slixx/pkg/strategy"
)

// This is just a sample.
func main() {
	origin := &_storage.FtpKind{}
	destination := &_storage.FtpKind{}

	origin.Initialize(&_storage.FtpKindConfiguration{
		Host:     "10.30.100.102:21",
		Username: "test",
		Password: "123123123",
		File:     "/test",
	})
	destination.Initialize(&_storage.FtpKindConfiguration{
		Host:     "10.30.100.102:21",
		Username: "test",
		Password: "123123123",
		File:     "/backup",
	})

	targetStrategy := _strategy.COPY
	err := targetStrategy.Execute(origin, destination)
	if err != nil {
		panic(err)
	}
}
