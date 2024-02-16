package graphql

import (
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"kroseida.org/slixx/internal/supervisor/graphql/controller"
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
	obj.FieldFunc("getBackups", controller.GetBackups)
	obj.FieldFunc("getSatelliteLogs", controller.GetSatelliteLogs)
	obj.FieldFunc("getExecutions", controller.GetExecutions)
	obj.FieldFunc("getExecution", controller.GetExecution)
	obj.FieldFunc("getExecutionHistory", controller.GetExecutionHistory)
	obj.FieldFunc("getBackup", controller.GetBackup)
	obj.FieldFunc("getJobSchedule", controller.GetJobSchedule)
	obj.FieldFunc("getJobSchedules", controller.GetJobSchedules)
	obj.FieldFunc("getJobScheduleKinds", controller.GetJobScheduleKinds)
}
