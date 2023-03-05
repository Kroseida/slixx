package application

type Settings struct {
	Authentication Authentication `json:"authentication" graphql:"authentication"`
	Http           Http           `json:"http" graphql:"http"`
	Database       Database       `json:"database" graphql:"database"`
	Logger         Log            `json:"logger" graphql:"logger"`
}

type Authentication struct {
	HashCost        int `json:"hashCost" graphql:"hashCost"`
	TokenSize       int `json:"tokenSize" graphql:"tokenSize"`
	SessionDuration int `json:"sessionDuration" graphql:"sessionDuration"`
}

type Http struct {
	AllowedOrigin string `json:"AllowedOrigin" graphql:"AllowedOrigin"`
}

type Database struct {
	Kind          string            `json:"kind" graphql:"kind"`
	Configuration map[string]string `json:"configuration" graphql:"configuration"`
}

type Log struct {
	Mode string `json:"mode" graphql:"mode"`
}

var DefaultSettings = Settings{
	Authentication: Authentication{
		HashCost:        16,
		TokenSize:       128,
		SessionDuration: 12,
	},
	Http: Http{
		AllowedOrigin: "*",
	},
	Database: Database{
		Kind: "sqlite",
		Configuration: map[string]string{
			"file": "slixx.db",
		},
	},
	Logger: Log{
		Mode: "info",
	},
}
var CurrentSettings = DefaultSettings
