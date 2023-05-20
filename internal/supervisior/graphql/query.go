package graphql

import (
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"kroseida.org/slixx/internal/supervisior/graphql/controller"
)

func registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()
	obj.FieldFunc("getJob", controller.GetJob)
	obj.FieldFunc("getJobs", controller.GetJobs)
	obj.FieldFunc("getStorage", controller.GetStorage)
	obj.FieldFunc("getStorages", controller.GetStorages)
	obj.FieldFunc("getUsers", controller.GetUsers)
	obj.FieldFunc("getUser", controller.GetUser)
	obj.FieldFunc("getLocalUser", controller.GetLocalUser)
	obj.FieldFunc("getStorageKinds", controller.GetStorageKinds)
	obj.FieldFunc("getJobStrategies", controller.GetJobStrategies)
	obj.FieldFunc("getPermissions", controller.GetPermissions)
	obj.FieldFunc("getSatellite", controller.GetSatellite)
	obj.FieldFunc("getSatellites", controller.GetSatellites)
}