package commands

import "readygo/models"

// RunMigration run migration
func RunMigration() {
	models.DB.AutoMigrate(
		&models.Admin{},
		&models.Authorization{},
		&models.Permission{},
		&models.Role{},
		&models.Tag{},
	)
}
