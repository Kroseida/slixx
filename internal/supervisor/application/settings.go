package application

type Settings struct {
	Authentication Authentication `json:"authentication" graphql:"authentication"`
	Http           Http           `json:"http" graphql:"http"`
	Database       Database       `json:"database" graphql:"database"`
	Logger         Log            `json:"logger" graphql:"logger"`
	LogSync        LogSync        `json:"logSync" graphql:"logSync"`
}

type Authentication struct {
	HashCost        int `json:"hashCost" graphql:"hashCost"`
	TokenSize       int `json:"tokenSize" graphql:"tokenSize"`
	SessionDuration int `json:"sessionDuration" graphql:"sessionDuration"`
}

type Http struct {
	AllowedOrigin  string `json:"AllowedOrigin" graphql:"AllowedOrigin"`
	BindAddress    string `json:"bindAddress" graphql:"bindAddress"`
	EnableGraphiql bool   `json:"enableGraphiql" graphql:"enableGraphiql"`
}

type Database struct {
	Kind          string            `json:"kind" graphql:"kind"`
	Configuration map[string]string `json:"configuration" graphql:"configuration"`
}

type Log struct {
	Mode string `json:"mode" graphql:"mode"`
}

type LogSync struct {
	Active        bool `json:"active" graphql:"active"`
	RunCleanUp    bool `json:"runCleanUp" graphql:"runCleanUp"`
	LogRetention  int  `json:"logRetention" graphql:"logRetention"`
	CheckInterval int  `json:"checkInterval" graphql:"checkInterval"`
}

var DefaultSettings = Settings{
	Authentication: Authentication{
		HashCost:        16,
		TokenSize:       128,
		SessionDuration: 12,
	},
	Http: Http{
		AllowedOrigin:  "*",
		BindAddress:    ":3030",
		EnableGraphiql: false,
	},
	Database: Database{
		Kind: "sqlite",
		Configuration: map[string]string{
			"file": "data/slixx.db",
		},
	},
	Logger: Log{
		Mode: "info",
	},
	LogSync: LogSync{
		Active:        true,
		RunCleanUp:    true,
		LogRetention:  7 * 24,
		CheckInterval: 6,
	},
}
var CurrentSettings = DefaultSettings
