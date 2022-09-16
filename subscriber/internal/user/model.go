package model

type CreateUser struct {
	DisplayName string `json:"displayName" binding:"required"`
}
