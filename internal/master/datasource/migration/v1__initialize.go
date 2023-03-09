package migration

import (
	"gorm.io/gorm"
)

type V1Initialize struct {
}

func (V1Initialize) GetName() string {
	return "v1__initialize"
}

func (V1Initialize) Migrate(database *gorm.DB) error {
	err := database.Exec(`CREATE TABLE users (id char(36) NOT NULL,
 												  name text NOT NULL,
 												  email text NOT NULL,
 												  first_name text NOT NULL,
 												  last_name text NOT NULL,
 												  active boolean NOT NULL,
 												  permissions text NOT NULL,
 												  description text NOT NULL,
 												  created_at DATETIME, 
 												  updated_at DATETIME, 
 												  deleted_at DATETIME, 
 												  PRIMARY KEY (id)
                   								 )`).Error
	if err != nil {
		return err
	}

	err = database.Exec(`CREATE TABLE authentications (id char(36) NOT NULL, 
														   kind text NOT NULL, 
														   configuration text NOT NULL, 
														   user_id char(36) NOT NULL, 
														   created_at DATETIME, 
														   updated_at DATETIME, 
														   deleted_at DATETIME, 
														   PRIMARY KEY (id),
														   FOREIGN KEY (user_id) REFERENCES users(id)
                             							  )`).Error
	if err != nil {
		return err
	}

	err = database.Exec(`CREATE TABLE sessions (id char(36) NOT NULL,
													token text NOT NULL,
	 												user_id char(36) NOT NULL,
	 												expires_at DATETIME NOT NULL,
	 												created_at DATETIME,
	 												updated_at DATETIME,
	 												deleted_at DATETIME,
	 												PRIMARY KEY (id),
	 												FOREIGN KEY (user_id) REFERENCES users(id)
                             					   )`).Error
	if err != nil {
		return err
	}

	err = database.Exec(`CREATE TABLE storages (id char(36) NOT NULL, 
                                                    name text NOT NULL, 
                                                    kind text NOT NULL,
                                                    configuration text NOT NULL,
                                                    created_at DATETIME, 
                                                    updated_at DATETIME, 
                                                    deleted_at DATETIME,
                                                    PRIMARY KEY (id)
                                                   )`).Error
	if err != nil {
		return err
	}

	return nil
}
