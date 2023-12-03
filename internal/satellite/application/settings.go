package application

import (
	"kroseida.org/slixx/pkg/utils"
)

type Settings struct {
	Logger    Log       `json:"logger" graphql:"logger"`
	Satellite Satellite `json:"satellite" graphql:"satellite"`
	LogSync   LogSync   `json:"logSync" graphql:"logSync"`
}

type LogSync struct {
	Active bool `json:"active" graphql:"active"`
}

type Log struct {
	Mode string `json:"mode" graphql:"mode"`
}

type Satellite struct {
	Network             Network `json:"network" graphql:"network"`
	AuthenticationToken string  `json:"authenticationToken" graphql:"authenticationToken"`
}

type Network struct {
	BindAddress string `json:"bindAddress" graphql:"bindAddress"`
}

var DefaultSettings = Settings{
	Logger: Log{
		Mode: "info",
	},
	LogSync: LogSync{
		Active: true,
	},
	Satellite: Satellite{
		Network: Network{
			BindAddress: ":9623",
		},
		AuthenticationToken: utils.GenerateSecureToken(32),
	},
}
var CurrentSettings = DefaultSettings
