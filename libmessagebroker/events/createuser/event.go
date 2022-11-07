package createuser

const QUEUE_NAME = "user.create"

type Model struct {
	DisplayName string `json:"displayName" binding:"required" bson:"displayName"`
}
