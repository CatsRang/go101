package model

// User represents the database entity with exported fields
type User struct {
	Id         int64
	Username   string
	Email      string
	Status     int
	CreateTime string
}
