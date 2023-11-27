package execution

import (
	"github.com/google/uuid"
	"kroseida.org/slixx/internal/supervisor/datasource"
	"kroseida.org/slixx/internal/supervisor/datasource/provider"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/strategy/statustype"
	"time"
)

func ApplyExecutionToIndex(
	id uuid.UUID,
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
