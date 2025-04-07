package item

import (
	"errors"
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

func (repo Repository) Create(project *model.RequestProject) error {
    return repo.Database.Create(project).Error
}
// Get All Projects
func (repo Repository) GetAll() ([]model.RequestProject, error) {
    var projects []model.RequestProject
    err := repo.Database.Find(&projects).Error
    return projects, err
}

// Get Project By ID
func (repo Repository) GetByID(id uint) (model.RequestProject, error) {
    var project model.RequestProject
    err := repo.Database.First(&project, id).Error
    return project, err
}
func (repo Repository) GetSubTasksByProjectID(projectID uint) ([]model.SubTaskDraft, error) {
    var subTasks []model.SubTaskDraft
    err := repo.Database.Where("project_id = ?", projectID).Find(&subTasks).Error
    return subTasks, err
}

func (repo Repository) CreateSubTask(subtask *model.SubTaskDraft) error {
    return repo.Database.Create(subtask).Error
}
func (repo Repository) DeleteSubTask(projectID uint, draftID uint) error {
    return repo.Database.Where("project_id = ? AND draft_id = ?", projectID, draftID).Delete(&model.SubTaskDraft{}).Error
}


//หน้าลงทุนแต่ละปี
func (repo Repository) UpsertBudgetRequests(requests []model.BudgetRequestDraft) error {
	for _, req := range requests {
		var existing model.BudgetRequestDraft
		err := repo.Database.Where(
			"list = ? AND year = ? AND draft_subtask_id = ? AND request_pay = ? AND requested_year = ?",
			req.List, req.Year, req.DraftSubtaskID, req.RequestPay, req.RequestedYear,
		).First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if existing.ID != 0 {
			// Update
			existing.Value = req.Value
			existing.CommitInvest = req.CommitInvest
			existing.RequestedYear = req.RequestedYear // ✅ สำคัญเผื่อ update field requested_year ด้วย

			if err := repo.Database.Save(&existing).Error; err != nil {
				return err
			}
		} else {
			// Insert
			if err := repo.Database.Create(&req).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (repo Repository) FindBudgetRequestsBySubtaskId(subtaskID uint, requestedYear *int) ([]model.BudgetRequestDraft, error) {
	var requests []model.BudgetRequestDraft
	query := repo.Database.Where("draft_subtask_id = ?", subtaskID)

	// ✅ ถ้ามี requestedYear ให้เพิ่มเงื่อนไข
	if requestedYear != nil {
		query = query.Where("requested_year = ?", *requestedYear)
	}

	err := query.Find(&requests).Error
	return requests, err
}
//ทำหน้า sumary
func (repo Repository) GetBudgetRequestsByProjectID(projectID uint) ([]model.BudgetRequestWithSubtask, error) {
	var results []model.BudgetRequestWithSubtask

	err := repo.Database.Table("budgetrequestsdraft").
		Select(`budgetrequestsdraft.request_draft_id,
			budgetrequestsdraft.draft_subtask_id,
			budgetrequestsdraft.list,
			budgetrequestsdraft.year,
			budgetrequestsdraft.request_pay,
			budgetrequestsdraft.value,
			budgetrequestsdraft.commit_invest,
			subtasksdraft.subtask_name`).
		Joins("JOIN subtasksdraft ON subtasksdraft.draft_id = budgetrequestsdraft.draft_subtask_id").
		Where("subtasksdraft.project_id = ?", projectID).
		Find(&results).Error

	return results, err
}
