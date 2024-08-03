package execution

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/application"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/internal/supervisor/slixxreactive"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/statustype"
	"time"
)

func StartTimeoutDetector() {
	if !application.CurrentSettings.LogSync.Active {
		return
	}
	for {
		application.Logger.Info("Setting Timeout to executions older than 3 days")
		err := datasource.ExecutionProvider.UpdateWithStatusAndOlderThan(statustype.Info, time.Now().Add(-time.Hour*3*24), statustype.Timeout)
		if err != nil {
			application.Logger.Error("Failed to update old executions: ", err)
		}

		time.Sleep(time.Hour * time.Duration(application.CurrentSettings.LogSync.CheckInterval))
	}
}

func ApplyExecutionToIndex(
	id uuid.UUID,
	kind string,
	jobId uuid.UUID,
	percentage float64,
	statusType string,
	message string,
) error {
	execution, err := datasource.ExecutionProvider.Get(id)
	if err != nil {
		return err
	}
	if execution == nil {
		_, err := datasource.ExecutionProvider.Create(
			id,
			kind,
			jobId,
			statusType,
			nil,
		)
		if err != nil {
			return err
		}
	}

	var finishedAt time.Time

	if statusType == statustype.Finished || statusType == statustype.Error {
		finishedAt = time.Now()
	}

	_, err = datasource.ExecutionProvider.Update(id, nil, &statusType, &finishedAt)
	if err != nil {
		return err
	}

	_, err = datasource.ExecutionProvider.CreateExecutionHistory(
		id,
		percentage,
		statusType,
		message,
	)
	if err != nil {
		return err
	}
	slixxreactive.Event("execution.update." + id.String())
	slixxreactive.Event("execution.update.*")
	return nil
}

func Get(id uuid.UUID) (*model.Execution, error) {
	return datasource.ExecutionProvider.Get(id)
}

func GetHistory(id uuid.UUID) ([]*model.ExecutionHistory, error) {
	return datasource.ExecutionProvider.ListHistory(id)
}

func GetPaged(pagination *provider.Pagination[model.Execution], jobId *uuid.UUID) (*provider.Pagination[model.Execution], error) {
	return datasource.ExecutionProvider.ListPaged(pagination, jobId)
}
