package item

import (
	"task-api/internal/constant"
	"task-api/internal/model"

	"gorm.io/gorm"
)

type Service struct {
	Repository Repository
}

func NewService(db *gorm.DB) Service {
	return Service{
		Repository: NewRepository(db),
	}
}
// สร้าง Item ใหม่
func (service Service) Create(req model.RequestItem) (model.Item, error) {
	item := model.Item{
		Title:    req.Title,
		Price:    req.Price,
		Quantity: req.Quantity,
		Status:   constant.ItemPendingStatus,
	}

	if err := service.Repository.Create(&item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}
// ค้นหา Item ทั้งหมด
func (service Service) Find(query model.RequestFindItem) ([]model.Item, error) {
	return service.Repository.Find(query)
}
//ค้นหา Item จากไอดีหนึาง
func (service Service) FindByID(id uint) (model.Item, error) {
    return service.Repository.FindByID(id)
}
// อัปเดตเฉพาะ Status ของ Item
func (service Service) UpdateStatus(id uint, status constant.ItemStatus) (model.Item, error) {
	// Find item
	item, err := service.Repository.FindByID(id)
	if err != nil {
		return model.Item{}, err
	}

	// Fill data
	item.Status = status

	// Replace
	if err := service.Repository.Replace(item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}
// อัปเดตข้อมูลทั่วไปของ Item
func (service Service) UpdateIteminfo(id uint, req model.RequestUpdateIteminfo) (model.Item, error) {
	// ค้นหา Item ตาม ID
	item, err := service.Repository.FindByID(id)
	if err != nil {
		return model.Item{}, err
	}

	// อัปเดตฟิลด์ต่างๆ
	item.Title = req.Title
	item.Price = req.Price
	item.Quantity = req.Quantity

	// แทนที่ข้อมูลในฐานข้อมูล
	if err := service.Repository.Replace(item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}

// ลบ Item ตาม ID
func (service Service) Delete(id uint) error {
	return service.Repository.DeleteByID(id)
}