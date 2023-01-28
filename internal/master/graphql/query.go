package graphql

import (
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"kroseida.org/slixx/internal/master/graphql/controller"
)

func registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()

	obj.FieldFunc("getStorage", controller.GetStorage)
	obj.FieldFunc("getStorages", controller.GetStorages)
	obj.FieldFunc("getUsers", controller.GetUsers)
	obj.FieldFunc("getUser", controller.GetUser)
	obj.FieldFunc("getLocalUser", controller.GetLocalUser)
}
