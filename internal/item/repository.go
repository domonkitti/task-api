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

    // เพิ่มเงื่อนไขการค้นหา
    if query.Statuses != "" {
        db = db.Where("status = ?", query.Statuses)
    }

    if query.Title != "" {
        db = db.Where("title = ?", query.Title)
    }

    // เรียงลำดับตามฟิลด์ที่ต้องการ (ตัวอย่างเช่น id)
    db = db.Order("id ASC")

    // ดึงข้อมูล
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
    if err := repo.Database.Delete(&model.Item{}, id).Error; err != nil {
        return err
    }
    return nil
}
