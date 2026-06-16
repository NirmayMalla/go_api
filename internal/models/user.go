package models

type User struct {
	ID		int32			`json:"id"`
	Name 	string	`json:"name"`
	DOB		string	`json:"dob"`
	Age		int			`json:"age"`
}
