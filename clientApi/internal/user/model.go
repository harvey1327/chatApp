package user

type Request struct {
	DisplayName string `json:"displayName" binding:"required"`
}
