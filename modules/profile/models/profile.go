package models

type Profile struct {
	ID             uint   `json:"id"`
	AccountID      uint64 `json:"account_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Mobile         string `json:"mobile"`
	ProfilePicture string `json:"profile_picture"`
}
