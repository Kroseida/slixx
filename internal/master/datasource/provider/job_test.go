package provider_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/internal/master/datasource"
	"testing"
)

func Test_CreateJob(t *testing.T) {
	teardownSuite := setupSuite()

	storage, err := datasource.StorageProvider.CreateStorage("Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.JobProvider.CreateJob("Test", "Test", "COPY", "{}", storage.Id, storage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	jobs, err := datasource.JobProvider.GetJobs()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 1, len(jobs))
	assert.Equal(t, "Test", jobs[0].Name)
	assert.Equal(t, "COPY", jobs[0].Strategy)
	assert.Equal(t, "{}", jobs[0].Configuration)
	teardownSuite()
}

func Test_CreateJob_MissingStorage(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.JobProvider.CreateJob("Test", "Test", "COPY", "{}", uuid.New(), uuid.New())
	if err == nil {
		t.Error("Expected error")
		teardownSuite()
		return
	}

	assert.Equal(t, "origin storage not found", err.Error())
	teardownSuite()
}

func Test_CreateJob_EmptyName(t *testing.T) {
	teardownSuite := setupSuite()

	originStorage, err := datasource.StorageProvider.CreateStorage("Origin Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.JobProvider.CreateJob("", "Test", "COPY", "{}", originStorage.Id, originStorage.Id)
	if err == nil {
		t.Error("Expected error")
		teardownSuite()
		return
	}

	assert.Equal(t, "name can not be empty", err.Error())
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

	job, err := datasource.JobProvider.CreateJob("Test", "Test", "COPY", "{}", storage.Id, storage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	updatedName := "Updated Name"
	datasource.JobProvider.UpdateJob(job.Id, &updatedName, nil, nil, nil, nil, nil)

	jobs, err := datasource.JobProvider.GetJobs()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, jobs[0].Name, updatedName)
	teardownSuite()
}

func Test_UpdateJob_MissingJob(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.JobProvider.UpdateJob(uuid.New(), nil, nil, nil, nil, nil, nil)
	if err == nil {
		t.Error("Expected error")
		teardownSuite()
		return
	}

	assert.Equal(t, "job not found", err.Error())
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

	job, err := datasource.JobProvider.CreateJob("Test", "Test", "COPY", "{}", storage.Id, storage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	datasource.JobProvider.DeleteJob(job.Id)

	jobs, err := datasource.JobProvider.GetJobs()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, 0, len(jobs))
	teardownSuite()
}

func Test_DeleteJob_MissingJob(t *testing.T) {
	teardownSuite := setupSuite()

	_, err := datasource.JobProvider.DeleteJob(uuid.New())
	if err == nil {
		t.Error("Expected error")
		teardownSuite()
		return
	}

	assert.Equal(t, "job not found", err.Error())
	teardownSuite()
}

func Test_GetJob(t *testing.T) {
	teardownSuite := setupSuite()

	storage, err := datasource.StorageProvider.CreateStorage("Storage", "", "FTP", "{}")
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	_, err = datasource.JobProvider.CreateJob("Test_Invalid", "Test", "COPY", "{}", storage.Id, storage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	job, err := datasource.JobProvider.CreateJob("Test", "Test", "COPY", "{}", storage.Id, storage.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	job, err = datasource.JobProvider.GetJob(job.Id)
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}

	assert.Equal(t, "Test", job.Name)
	teardownSuite()
}
