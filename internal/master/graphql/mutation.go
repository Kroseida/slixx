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
	obj.FieldFunc("updateUserName", controller.UpdateUserName)
	obj.FieldFunc("updateUserFirstName", controller.UpdateUserFirstName)
	obj.FieldFunc("updateUserLastName", controller.UpdateUserLastName)
	obj.FieldFunc("updateUserEmail", controller.UpdateUserEmail)
	obj.FieldFunc("updateUserDescription", controller.UpdateUserDescription)
	obj.FieldFunc("updateUserActive", controller.UpdateUserActive)
	obj.FieldFunc("addUserPermission", controller.AddUserPermission)
	obj.FieldFunc("removeUserPermission", controller.RemoveUserPermission)
	obj.FieldFunc("createPasswordAuthentication", controller.CreatePasswordAuthentication)
	obj.FieldFunc("authenticate", controller.Authenticate)
}
