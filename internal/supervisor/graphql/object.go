package graphql

import (
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"kroseida.org/slixx/internal/supervisor/graphql/controller"
)

func registerObject(schema *schemabuilder.Schema) {
	storage(schema)
}

func storage(schema *schemabuilder.Schema) {
	schema.Object("Storage", controller.Storage{}).Key("id")
	schema.Object("StoragePrototype", controller.StoragePrototype{}).Key("id")
	schema.Object("User", controller.User{}).Key("id")
	schema.Object("Session", controller.Session{}).Key("id")
	schema.Object("Authentication", controller.Authentication{}).Key("id")
}
