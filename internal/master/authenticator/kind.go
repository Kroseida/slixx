package authenticator

import (
	"kroseida.org/slixx/pkg/model"
)

type Kind interface {
	Validate(authentication *model.Authentication, request string) (bool, error)
	ParseConfiguration(configurationJson string) (any, error)
}

var PASSWORD = "password"

var kinds = map[string]Kind{
	PASSWORD: Password{},
}

func GetKind(kind string) Kind {
	return kinds[kind]
}
