package provider

import (
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"gorm.io/gorm"
	"kroseida.org/slixx/pkg/model"
	"math"
)

// SatelliteProvider Satellite Provider
type SatelliteProvider struct {
	Database *gorm.DB
}

func (provider SatelliteProvider) Delete(id uuid.UUID) (*model.Satellite, error) {
	satellite, err := provider.Get(id)
	if satellite == nil {
		return nil, graphql.NewSafeError("satellite not found")
	}
	if err != nil {
		return nil, err
	}

	result := provider.Database.Delete(&satellite)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return satellite, nil
}

func (provider SatelliteProvider) Create(name string, description string, address string, token string) (*model.Satellite, error) {
	if name == "" {
		return nil, graphql.NewSafeError("name can not be empty")
	}
	if address == "" {
		return nil, graphql.NewSafeError("address can not be empty")
	}
	if token == "" {
		return nil, graphql.NewSafeError("token can not be empty")
	}

	satellite := model.Satellite{
		Id:          uuid.New(),
		Name:        name,
		Description: description,
		Address:     address,
		Token:       token,
	}

	result := provider.Database.Create(&satellite)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return &satellite, nil
}

func (provider SatelliteProvider) Update(id uuid.UUID, name *string, description *string, address *string, token *string) (*model.Satellite, error) {
	updateSatellite, err := provider.Get(id)
	if updateSatellite == nil {
		return nil, graphql.NewSafeError("satellite not found")
	}
	if err != nil {
		return nil, err
	}

	if name != nil {
		if *name == "" {
			return nil, graphql.NewSafeError("name can not be empty")
		}
		updateSatellite.Name = *name
	}

	if address != nil {
		updateSatellite.Address = *address
	}
	if token != nil {
		updateSatellite.Token = *token
	}
	if description != nil {
		updateSatellite.Description = *description
	}

	result := provider.Database.Save(&updateSatellite)
	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return updateSatellite, nil
}

func (provider SatelliteProvider) List() ([]*model.Satellite, error) {
	var satellites []*model.Satellite
	result := provider.Database.Find(&satellites)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return satellites, nil
}

func (provider SatelliteProvider) ListPaged(pagination *Pagination[model.Satellite]) (*Pagination[model.Satellite], error) {
	context := paginate(model.Satellite{}, "name", pagination, provider.Database)

	var satellites []model.Satellite
	result := context.Find(&satellites)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = satellites
	return pagination, nil
}

func (provider SatelliteProvider) Get(id uuid.UUID) (*model.Satellite, error) {
	var satellite *model.Satellite
	result := provider.Database.First(&satellite, "id = ?", id)

	if isSqlError(result.Error) {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return satellite, nil
}

func (provider SatelliteProvider) ApplyLogs(logs []*model.SatelliteLogEntry) error {
	result := provider.Database.Create(&logs)

	if isSqlError(result.Error) {
		return result.Error
	}

	return nil
}

func (provider SatelliteProvider) GetLogs(satelliteId uuid.UUID, pagination *Pagination[model.SatelliteLogEntry]) (*Pagination[model.SatelliteLogEntry], error) {
	var totalRows int64
	provider.Database.Model(model.SatelliteLogEntry{}).
		Where("satellite_id = ?", satelliteId).
		Order("logged_at DESC").
		Where("message like ?", "%"+pagination.Search+"%").
		Or("id like ?", "%"+pagination.Search+"%").
		Count(&totalRows)

	pagination.TotalRows = totalRows
	if pagination.GetLimit() > 0 {
		pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	} else {
		pagination.TotalPages = 1
	}

	context := provider.Database.Offset(pagination.GetOffset()).
		Where("satellite_id = ?", satelliteId).
		Order("logged_at DESC")

	if pagination.GetLimit() > 0 {
		context = context.Limit(pagination.GetLimit())
	}
	context.Where("message like ?", "%"+pagination.Search+"%").
		Or("id like ?", "%"+pagination.Search+"%")

	var logs []model.SatelliteLogEntry
	result := context.Find(&logs)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = logs
	return pagination, nil
}
