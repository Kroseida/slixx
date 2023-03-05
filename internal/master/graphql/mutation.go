package graphql

import (
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"kroseida.org/slixx/internal/master/graphql/controller"
)

func registerMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()
	obj.FieldFunc("createStorage", controller.CreateStorage)
	obj.FieldFunc("deleteStorage", controller.DeleteStorage)
	obj.FieldFunc("createUser", controller.CreateUser)
	obj.FieldFunc("deleteUser", controller.DeleteUser)
	obj.FieldFunc("updateUser", controller.UpdateUser)
	obj.FieldFunc("addUserPermission", controller.AddUserPermission)
	obj.FieldFunc("removeUserPermission", controller.RemoveUserPermission)
	obj.FieldFunc("createPasswordAuthentication", controller.CreatePasswordAuthentication)
	obj.FieldFunc("authenticate", controller.Authenticate)
}
