package provider

import (
	"github.com/google/uuid"
	"github.com/samsarahq/thunder/graphql"
	"gorm.io/gorm"
	"kroseida.org/slixx/pkg/model"
)

// SatelliteProvider Satellite Provider
type SatelliteProvider struct {
	Database *gorm.DB
}

func (provider SatelliteProvider) DeleteSatellite(id uuid.UUID) (*model.Satellite, error) {
	satellite, err := provider.GetSatellite(id)
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

func (provider SatelliteProvider) CreateSatellite(name string, description string, address string, token string) (*model.Satellite, error) {
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

func (provider SatelliteProvider) UpdateSatellite(id uuid.UUID, name *string, description *string, address *string, token *string) (*model.Satellite, error) {
	updateSatellite, err := provider.GetSatellite(id)
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

func (provider SatelliteProvider) GetSatellites() ([]*model.Satellite, error) {
	var satellites []*model.Satellite
	result := provider.Database.Find(&satellites)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	return satellites, nil
}

func (provider SatelliteProvider) GetSatellitesPaged(pagination *Pagination[model.Satellite]) (*Pagination[model.Satellite], error) {
	context := paginate(model.Satellite{}, "name", pagination, provider.Database)

	var satellites []model.Satellite
	result := context.Find(&satellites)

	if isSqlError(result.Error) {
		return nil, result.Error
	}

	pagination.Rows = satellites
	return pagination, nil
}

func (provider SatelliteProvider) GetSatellite(id uuid.UUID) (*model.Satellite, error) {
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

func (provider SatelliteProvider) CreateSatelliteLogs(logs []*model.SatelliteLogEntry) error {
	result := provider.Database.Create(&logs)

	if isSqlError(result.Error) {
		return result.Error
	}

	return nil
}
