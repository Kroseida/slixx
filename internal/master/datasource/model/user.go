package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id          uuid.UUID `sql:"default:uuid_generate_v4()",json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Active      bool      `json:"active"`
	Description string    `json:"description"`
	Permissions string    `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Authentication struct {
	Id            uuid.UUID `sql:"default:uuid_generate_v4()",json:"id"`
	Kind          string    `json:"kind"`
	Configuration string    `json:"-",graphql-secret:"true"`
	UserId        uuid.UUID `json:"userId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Session struct {
	Id        uuid.UUID `sql:"default:uuid_generate_v4()",json:"id",graphql:"id"`
	UserId    uuid.UUID `json:"userId",graphql:"userId"`
	ExpiresAt time.Time `json:"expiresAt",graphql:"expiresAt"`
	CreatedAt time.Time `json:"createdAt",graphql:"createdAt"`
	Token     string    `json:"token",graphql:"token",graphql-secret:"true"`
}

func (u *User) GetActivePermissions() ([]string, error) {
	var activePermissions []string
	err := json.Unmarshal([]byte(u.Permissions), &activePermissions)
	if err != nil {
		return nil, err
	}
	return activePermissions, nil
}

func (u *User) IsPermitted(permission []string) (bool, error) {
	activePermissions, err := u.GetActivePermissions()
	if err != nil {
		return false, err
	}
	for _, activePermission := range activePermissions {
		for _, permission := range permission {
			if activePermission == permission {
				return true, nil
			}
		}
	}
	return false, err
}

func (u *User) AddPermission(permissions []string) error {
	activePermissions, err := u.GetActivePermissions()
	if err != nil {
		return err
	}
	for _, permission := range permissions {
		activePermissions = append(activePermissions, permission)
	}
	jsonData, err := json.Marshal(activePermissions)
	if err != nil {
		return err
	}
	u.Permissions = string(jsonData)

	return nil
}

func (u *User) RemovePermission(permissions []string) error {
	activePermissions, err := u.GetActivePermissions()
	if err != nil {
		return err
	}
	for _, permission := range permissions {
		for _, activePermission := range activePermissions {
			if activePermission == permission {
				activePermissions = remove(activePermissions, permission)
			}
		}
	}
	jsonData, err := json.Marshal(activePermissions)
	if err != nil {
		return err
	}
	u.Permissions = string(jsonData)
	return nil
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
