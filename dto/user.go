package dto

type User struct {
	FirstName   string  `json:"fname"`
	LastName    string  `json:"lname"`
	Email       string  `json:"email"`
	PostAddress Address `json:"Address"`
}
type Address struct {
	HouseNo string
	City    string
}

// dto คือสิ่งที่อยากให้ user ดู
