package model

type Attendee struct {
	Id       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Phone_no string `db:"phone_no" json:"phone_no"`
	Email    string `db:"email" json:"email"`
}
