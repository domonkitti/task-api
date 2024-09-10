package item

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

func (repo Repository) Create(item *model.Item) error {
	return repo.Database.Create(item).Error
}

func (repo Repository) Find(query model.RequestFindItem) ([]model.Item, error) {
	var results []model.Item

	db := repo.Database

	if statuses := query.Statuses; len(statuses) > 0 {
		db = db.Where("status = ?", statuses)
	}

	if err := db.Find(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

func (repo Repository) FindByID(id uint) (model.Item, error) {
	var result model.Item

	if err := repo.Database.First(&result, id).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (repo Repository) Replace(item model.Item) error {
	return repo.Database.Model(&item).Updates(item).Error
}

func (repo Repository) DeleteByID(id uint) error {
    var item model.Item

    // ตรวจสอบว่าไอเท็มมีอยู่หรือไม่
    if err := repo.Database.First(&item, id).Error; err != nil {
        return err // ถ้าไม่พบไอเท็ม ให้ส่งคืนข้อผิดพลาด
    }

    // ลบไอเท็ม
    if err := repo.Database.Delete(&item).Error; err != nil {
        return err
    }

    return nil
}