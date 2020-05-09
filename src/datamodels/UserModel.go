package datamodels

type UserModel struct {
	ID           int    `sql:"id"`
	Username     string `sql:"user_name"`
	Userpwd     string `sql:"user_pwd"`
	Email        string `sql:"email"`
	Phone        string `sql:"phone"`
	SignupAt     string `sql:"signup_at"`
	LastActiveAt string `sql:"last_active"`
	Status       int    `sql:"status"`
}
