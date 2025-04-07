package constant

type ProjectStatus string

const (
	ProjectPendingStatus  ItemStatus = "PENDING"
	ProjectApprovedStatus ItemStatus = "APPROVED"
	ProjectRejectedStatus ItemStatus = "REJECTED"
)
