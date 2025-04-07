package item

import (
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
func (service Service) Create(req model.RequestProject) (model.RequestProject, error) {
    // กำหนดค่า Status เป็น "pending" เสมอ ไม่สนใจค่าที่ส่งเข้ามา
    req.Status = "pending"

    project := model.RequestProject{
        ProjectName:  req.ProjectName,
        Status:       req.Status, // จะเป็น "pending"
        StartDate:    req.StartDate,
        EndDate:      req.EndDate,
		Type : req.Type,
		SubType : req.SubType,
		Owner: req.Owner,
		Budget:req.Budget,
    }

    if err := service.Repository.Create(&project); err != nil {
        return model.RequestProject{}, err
    }

    return project, nil
}
// Get All Projects
func (service Service) GetAll() ([]model.RequestProject, error) {
    return service.Repository.GetAll()
}

// Get Project By ID
func (service Service) GetByID(id uint) (model.RequestProject, error) {
    return service.Repository.GetByID(id)
}
func (service Service) GetSubTasksByProjectID(projectID uint) ([]model.SubTaskDraft, error) {
    return service.Repository.GetSubTasksByProjectID(projectID)
}

func (s Service) CreateSubTask(projectID uint, name string, totalFunding float64) (model.SubTaskDraft, error) {
    subtask := model.SubTaskDraft{
        ProjectID:    projectID,
        SubTaskName:  name,
        TotalFunding: totalFunding,
    }

    if err := s.Repository.CreateSubTask(&subtask); err != nil {
        return model.SubTaskDraft{}, err
    }

    return subtask, nil
}
func (s Service) DeleteSubTask(projectID uint, draftID uint) error {
    return s.Repository.DeleteSubTask(projectID, draftID)
}

//หน้าลงทุนแต่ละปั
func (s Service) CreateOrUpdateBudgetRequests(subtaskID uint, inputs []model.BudgetRequestDraftInput) error {
	var requests []model.BudgetRequestDraft

	for _, input := range inputs {
		requests = append(requests, model.BudgetRequestDraft{
			DraftSubtaskID: subtaskID,
			List:           input.List,
			Year:           input.Year,
            RequestedYear:  input.RequestedYear,
			RequestPay:     input.RequestPay,  // ✅ map ใหม่
			Value:          input.Value,       // ✅ map ใหม่
			CommitInvest:   input.CommitInvest,
		})
	}

	return s.Repository.UpsertBudgetRequests(requests)
}
func (s Service) GetBudgetRequestsBySubtaskId(subtaskID uint, requestedYear *int) ([]model.BudgetRequestDraft, error) {
	return s.Repository.FindBudgetRequestsBySubtaskId(subtaskID, requestedYear)
}
//ทำ summary
func (s Service) GetBudgetRequestsByProjectID(projectID uint) ([]model.BudgetRequestWithSubtask, error) {
	return s.Repository.GetBudgetRequestsByProjectID(projectID)
}

