package user

type Request struct {
	DisplayName string `json:"displayName" binding:"required"`
}

type EventMessage struct {
	ID string `json:"id"`
}

type Status int

const (
	PENDING Status = iota
	COMPLETE
	FAILED
)
