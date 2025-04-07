package model

type RequestProject struct {
	ProjectID    uint      `gorm:"primaryKey;column:project_id" json:"project_id"`
	ProjectName  string    `gorm:"column:project_name" json:"project_name"`
	Status       string    `gorm:"column:status" json:"status"`
	StartDate    string `gorm:"column:start_date" json:"start_date"`
	EndDate      string `gorm:"column:end_date" json:"end_date"`
	Type string			   `gorm:"column:type" json:"type"`
	SubType string		   `gorm:"column:subtype" json:"subtype"`
	Owner string 		   `gorm:"column:owner" json:"owner"`
	Budget int		   `gorm:"column:budget" json:"budget"`
	// คุณสามารถเพิ่มฟิลด์อื่น ๆ ได้ เช่น Type, Subtype เป็นต้น
}
// TableName กำหนดชื่อ table ให้ GORM
func (RequestProject) TableName() string {
    return "projects"
}
// ==== SubTaskDraft งานย่อย ====
type SubTaskDraft struct {
	DraftID       uint    `gorm:"primaryKey;column:draft_id" json:"draft_id"`
	ProjectID     uint    `gorm:"column:project_id" json:"project_id"`
	SubTaskName   string  `gorm:"column:subtask_name" json:"subtask_name"`
	TotalFunding  float64 `gorm:"column:total_funding" json:"total_funding"`
}

func (SubTaskDraft) TableName() string {
	return "subtasksdraft"
}

type CreateSubTaskRequest struct {
    SubTaskName  string  `json:"subtask_name"`
    TotalFunding float64 `json:"total_funding"`
}
//เอาไว้ทำเงินกู้เงินรายได  ดึงข้อมูล
type BudgetRequestDraft struct {
	ID             uint    `gorm:"primaryKey;column:request_draft_id"`
	DraftSubtaskID uint    `gorm:"column:draft_subtask_id"`
	List           string  `gorm:"column:list"`
	Year           int     `gorm:"column:year"`
	RequestedYear  int     `json:"requested_year"`
	RequestPay     string  `gorm:"column:request_pay"` // ✅ ใหม่!
	Value          float64 `gorm:"column:value"`      // ✅ ใหม่!
	CommitInvest   string  `gorm:"column:commit_invest"`
}

func (BudgetRequestDraft) TableName() string {
	return "budgetrequestsdraft"
}

// เขียน
type BudgetRequestDraftInput struct {
	List          string  `json:"list"`
	Year          int     `json:"year"`
	RequestedYear int 	  `json:"requestedYear" gorm:"column:requested_year"` // ✅ เพิ่มเข้า struct ด้วย!
	RequestPay    string  `json:"requestPay"`
	Value         float64 `json:"value"`
	CommitInvest  string  `json:"commitInvest"`
}
// ✅ Struct ตัวใหม่ สำหรับแสดงผลรวมชื่อ subtask ด้วย
type BudgetRequestWithSubtask struct {
	RequestDraftID uint    `json:"requestDraftId"` // ✅ เปลี่ยนให้ตรง
	DraftSubtaskID uint    `json:"draftSubtaskId"`
	List           string  `json:"list"`
	Year           int     `json:"year"`
	RequestPay     string  `json:"requestPay"`
	Value          float64 `json:"value"`
	CommitInvest   string  `json:"commitInvest"`
	SubtaskName    string  `json:"subtaskName"`
}
