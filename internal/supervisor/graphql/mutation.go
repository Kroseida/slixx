package graphql

import (
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"kroseida.org/slixx/internal/supervisor/graphql/controller"
)

func registerMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()
	obj.FieldFunc("createJob", controller.CreateJob)
	obj.FieldFunc("updateJob", controller.UpdateJob)
	obj.FieldFunc("deleteJob", controller.DeleteJob)
	obj.FieldFunc("createStorage", controller.CreateStorage)
	obj.FieldFunc("updateStorage", controller.UpdateStorage)
	obj.FieldFunc("deleteStorage", controller.DeleteStorage)
	obj.FieldFunc("createUser", controller.CreateUser)
	obj.FieldFunc("deleteUser", controller.DeleteUser)
	obj.FieldFunc("updateUser", controller.UpdateUser)
	obj.FieldFunc("addUserPermission", controller.AddUserPermission)
	obj.FieldFunc("removeUserPermission", controller.RemoveUserPermission)
	obj.FieldFunc("createPasswordAuthentication", controller.CreatePasswordAuthentication)
	obj.FieldFunc("authenticate", controller.Authenticate)
	obj.FieldFunc("createSatellite", controller.CreateSatellite)
	obj.FieldFunc("updateSatellite", controller.UpdateSatellite)
	obj.FieldFunc("deleteSatellite", controller.DeleteSatellite)
	obj.FieldFunc("executeBackup", controller.ExecuteBackup)
}
