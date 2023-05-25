package controller_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/graphql/controller"
	"testing"
)

func Test_GetJob(t *testing.T) {
	teardownSuite := setupSuite()

	storage, err := datasource.StorageProvider.CreateStorage("Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	executor, err := datasource.SatelliteProvider.CreateSatellite(
		"Executor",
		"description",
		"address",
		"token",
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	created, err := controller.CreateJob(withPermissions([]string{"job.create"}), controller.CreateJobDto{
		Name:                 "Testaaaaaa",
		Description:          "description",
		Strategy:             "COPY",
		Configuration:        "{}",
		OriginStorageId:      storage.Id,
		DestinationStorageId: storage.Id,
		ExecutorSatelliteId:  executor.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	job, err := controller.GetJob(withPermissions([]string{"job.view"}), controller.GetJobDto{
		Id: created.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, created.Id, job.Id)
	assert.Equal(t, "Testaaaaaa", job.Name)
	assert.Equal(t, "description", job.Description)
	assert.Equal(t, "COPY", job.Strategy)
	assert.Equal(t, "{\"blockSize\":0}", job.Configuration)
	assert.Equal(t, storage.Id, job.OriginStorageId)
	assert.Equal(t, storage.Id, job.DestinationStorageId)

	teardownSuite()
}

func Test_GetJobs(t *testing.T) {
	teardownSuite := setupSuite()

	storage, err := datasource.StorageProvider.CreateStorage("Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	executor, err := datasource.SatelliteProvider.CreateSatellite(
		"Executor",
		"description",
		"address",
		"token",
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = controller.CreateJob(withPermissions([]string{"job.create"}), controller.CreateJobDto{
		Name:                 "Testaaaaaa",
		Description:          "description1",
		Strategy:             "COPY",
		Configuration:        "{}",
		OriginStorageId:      storage.Id,
		DestinationStorageId: storage.Id,
		ExecutorSatelliteId:  executor.Id,
	})
	_, err = controller.CreateJob(withPermissions([]string{"job.create"}), controller.CreateJobDto{
		Name:                 "Testaaaaaa2",
		Description:          "description2",
		Strategy:             "COPY",
		Configuration:        "{}",
		OriginStorageId:      storage.Id,
		DestinationStorageId: storage.Id,
		ExecutorSatelliteId:  executor.Id,
	})
	_, err = controller.CreateJob(withPermissions([]string{"job.create"}), controller.CreateJobDto{
		Name:                 "Testaaaaaa3",
		Description:          "description3",
		Strategy:             "COPY",
		Configuration:        "{}",
		OriginStorageId:      storage.Id,
		DestinationStorageId: storage.Id,
		ExecutorSatelliteId:  executor.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	jobs, err := controller.GetJobs(withPermissions([]string{"job.view"}), controller.PageArgs{})

	assert.Equal(t, 3, len(jobs.Rows))
	assert.Equal(t, "Testaaaaaa", jobs.Rows[0].Name)
	assert.Equal(t, "description1", jobs.Rows[0].Description)
	assert.Equal(t, "Testaaaaaa2", jobs.Rows[1].Name)
	assert.Equal(t, "description2", jobs.Rows[1].Description)
	assert.Equal(t, "Testaaaaaa3", jobs.Rows[2].Name)
	assert.Equal(t, "description3", jobs.Rows[2].Description)
	teardownSuite()
}

func Test_DeleteJob(t *testing.T) {
	teardownSuite := setupSuite()

	storage, err := datasource.StorageProvider.CreateStorage("Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	executor, err := datasource.SatelliteProvider.CreateSatellite(
		"Executor",
		"description",
		"address",
		"token",
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = controller.CreateJob(withPermissions([]string{"job.create"}), controller.CreateJobDto{
		Name:                 "Testaaaaaa",
		Description:          "description1",
		Strategy:             "COPY",
		Configuration:        "{}",
		OriginStorageId:      storage.Id,
		DestinationStorageId: storage.Id,
		ExecutorSatelliteId:  executor.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	jobs, err := controller.GetJobs(withPermissions([]string{"job.view"}), controller.PageArgs{})
	assert.Equal(t, 1, len(jobs.Rows))

	_, err = controller.DeleteJob(withPermissions([]string{"job.delete"}), controller.DeleteJobDto{
		Id: jobs.Rows[0].Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	jobs, err = controller.GetJobs(withPermissions([]string{"job.view"}), controller.PageArgs{})
	assert.Equal(t, 0, len(jobs.Rows))
	teardownSuite()
}

func Test_CreateJob(t *testing.T) {
	teardownSuite := setupSuite()

	storage, err := datasource.StorageProvider.CreateStorage("Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	executor, err := datasource.SatelliteProvider.CreateSatellite(
		"Executor",
		"description",
		"address",
		"token",
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	createdJob, err := controller.CreateJob(withPermissions([]string{"job.create"}), controller.CreateJobDto{
		Name:                 "Testaaaaaa",
		Description:          "description1",
		Strategy:             "COPY",
		Configuration:        "{}",
		OriginStorageId:      storage.Id,
		DestinationStorageId: storage.Id,
		ExecutorSatelliteId:  executor.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	job, err := controller.GetJob(withPermissions([]string{"job.view"}), controller.GetJobDto{
		Id: createdJob.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, createdJob.Id, job.Id)
	assert.Equal(t, "Testaaaaaa", job.Name)
	assert.Equal(t, "description1", job.Description)
	assert.Equal(t, "COPY", job.Strategy)
	assert.Equal(t, "{\"blockSize\":0}", job.Configuration)
	assert.Equal(t, storage.Id, job.OriginStorageId)
	assert.Equal(t, storage.Id, job.DestinationStorageId)
	teardownSuite()
}

func Test_UpdateJob(t *testing.T) {
	teardownSuite := setupSuite()

	storage, err := datasource.StorageProvider.CreateStorage("Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	executor, err := datasource.SatelliteProvider.CreateSatellite(
		"Executor",
		"description",
		"address",
		"token",
	)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	createdJob, err := controller.CreateJob(withPermissions([]string{"job.create"}), controller.CreateJobDto{
		Name:                 "Testaaaaaa",
		Description:          "description1",
		Strategy:             "COPY",
		Configuration:        "{}",
		OriginStorageId:      storage.Id,
		DestinationStorageId: storage.Id,
		ExecutorSatelliteId:  executor.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	newName := "Testaaaaaa2"
	_, err = controller.UpdateJob(withPermissions([]string{"job.update"}), controller.UpdateJobDto{
		Id:   createdJob.Id,
		Name: &newName,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	updatedJob, err := controller.GetJob(withPermissions([]string{"job.view"}), controller.GetJobDto{
		Id: createdJob.Id,
	})
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, createdJob.Id, updatedJob.Id)
	assert.Equal(t, "Testaaaaaa2", updatedJob.Name)
	teardownSuite()
}
