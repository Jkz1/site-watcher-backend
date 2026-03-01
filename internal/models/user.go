package models

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username" validate:"required,min=4,max=20"`
	Password string `db:"password" json:"-" validate:"required,min=8"` // "-" means never show password in JSON
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
