package item

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"task-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Controller struct {
	Service Service
}

func NewController(db *gorm.DB) Controller {
	return Controller{
		Service: NewService(db),
	}
}

type ApiError struct {
	Field  string
	Reason string
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "จำเป็นต้องกรอกข้อมูลนี้"
	case "email":
		return "Invalid email"
	case "gt":
		return fmt.Sprintf("Number must greater than %v", param)
	case "gte":
		return fmt.Sprintf("Number must greater than or equal %v", param)
	}
	return ""
}

func getValidationErrors(err error) []ApiError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Param())}
		}
		return out
	}
	return nil
}
//สร้างของใหม่
func (controller Controller) CreateProject(ctx *gin.Context) {
    var request model.RequestProject
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "ข้อมูลไม่ถูกต้อง",
        })
        return
    }

    project, err := controller.Service.Create(request)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "เกิดข้อผิดพลาดขณะสร้างข้อมูล",
        })
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "data": project,
    })
}
// Get All Projects
func (controller Controller) GetAllProjects(ctx *gin.Context) {
    projects, err := controller.Service.GetAll()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "เกิดข้อผิดพลาดขณะดึงข้อมูล",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data": projects,
    })
}

// Get Project By ID
func (controller Controller) GetProjectByID(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "ID ไม่ถูกต้อง",
        })
        return
    }

    project, err := controller.Service.GetByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "message": "ไม่พบข้อมูลโครงการ",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data": project,
    })
}
func (controller Controller) GetSubTasksByProjectID(ctx *gin.Context) {
    idParam := ctx.Param("project_id")
    projectID, err := strconv.Atoi(idParam)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": "project_id ไม่ถูกต้อง",
        })
        return
    }

    subTasks, err := controller.Service.GetSubTasksByProjectID(uint(projectID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "เกิดข้อผิดพลาดขณะดึงข้อมูลงานย่อย",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data": subTasks,
    })
}

func (c Controller) CreateSubTask(ctx *gin.Context) {
    var (
        request model.CreateSubTaskRequest
        projectID uint
    )

    if id, err := strconv.Atoi(ctx.Param("projectId")); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "projectId ไม่ถูกต้อง"})
        return
    } else {
        projectID = uint(id)
    }

    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
        return
    }

    subtask, err := c.Service.CreateSubTask(projectID, request.SubTaskName, request.TotalFunding)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการสร้างงานย่อย"})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"data": subtask})
}
func (c Controller) DeleteSubTask(ctx *gin.Context) {
    projectIDParam := ctx.Param("projectId")
    draftIDParam := ctx.Query("draftId") // รับจาก query param

    projectID, err := strconv.Atoi(projectIDParam)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "projectId ไม่ถูกต้อง"})
        return
    }

    draftID, err := strconv.Atoi(draftIDParam)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "draftId ไม่ถูกต้อง"})
        return
    }

    err = c.Service.DeleteSubTask(uint(projectID), uint(draftID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการลบงานย่อย"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "ลบงานย่อยสำเร็จ"})
}
//หน้าลงทุนแต่ละปี
func (c Controller) CreateBudgetRequests(ctx *gin.Context) {
    // ✅ ใช้ subProjectId ตาม route
    subProjectIDParam := ctx.Param("subProjectId")

    subProjectIDInt, err := strconv.Atoi(subProjectIDParam)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "SubProject ID ไม่ถูกต้อง"})
        return
    }

    var inputs []model.BudgetRequestDraftInput
    if err := ctx.ShouldBindJSON(&inputs); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
        return
    }

    if err := c.Service.CreateOrUpdateBudgetRequests(uint(subProjectIDInt), inputs); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถบันทึกข้อมูลได้"})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"message": "บันทึกข้อมูลสำเร็จ"})
}
func (c Controller) GetBudgetRequestsBySubtaskId(ctx *gin.Context) {
	subtaskIDParam := ctx.Param("subProjectId")
	subtaskID, err := strconv.Atoi(subtaskIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "SubProject ID ไม่ถูกต้อง"})
		return
	}

	// ✅ ดึง query param requestedYear (optional)
	requestedYearParam := ctx.Query("requestedYear")
	var requestedYear *int

	if requestedYearParam != "" {
		year, err := strconv.Atoi(requestedYearParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Requested Year ไม่ถูกต้อง"})
			return
		}
		requestedYear = &year // ✅ assign ค่าเข้า pointer
	}

	result, err := c.Service.GetBudgetRequestsBySubtaskId(uint(subtaskID), requestedYear)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลได้"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
//ทำsummary
func (c Controller) GetBudgetRequestsByProjectID(ctx *gin.Context) {
	projectIDParam := ctx.Param("projectId")

	// แปลง projectId จาก string -> uint
	projectID, err := strconv.Atoi(projectIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Project ID ไม่ถูกต้อง"})
		return
	}

	// เรียก Service
	data, err := c.Service.GetBudgetRequestsByProjectID(uint(projectID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลได้"})
		return
	}

	ctx.JSON(http.StatusOK, data)
}