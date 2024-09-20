package user

import (
	"task-api/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	Database *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		Database: db,
	}
}

func (repo Repository) FindOneByUsername(username string) (model.User, error) {
	var result model.User

	db := repo.Database
	db = db.Where("username = ?", username)

	if err := db.Find(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}
func (repo Repository) Signup(user *model.User) error {
	return repo.Database.Create(user).Error
}