package user

// AuthPayload ...
type AuthPayload struct {
	Token *string `gqlgen:"token"`
	User  *User   `gqlgen:"user"`
}

// RegisterInput ...
type RegisterInput struct {
	Name      *string `gqlgen:"name" validate:"required,min=3,max=64"`
	Email     *string `gqlgen:"email" validate:"required,email,max=64"`
	StudentID *string `gqlgen:"studentID" validate:"required,numeric,len=9"`
	Password  *string `gqlgen:"password" validate:"required,min=6,max=64"`
}

// ToUser ...
func (r *RegisterInput) ToUser() *User {
	return &User{
		Name:      r.Name,
		Email:     r.Email,
		StudentID: r.StudentID,
		Password:  r.Password,
	}
}

// LoginInput ...
type LoginInput struct {
	Email     *string `db:"email" gqlgen:"email" validate:"omitempty,email,required_without=StudentID,max=64"`
	StudentID *string `db:"student_id" gqlgen:"studentID" validate:"omitempty,numeric,required_without=Email,len=9"`
	Password  *string `db:"password" validate:"required,min=6,max=64"`
}

// ToUser ...
func (r *LoginInput) ToUser() *User {
	return &User{
		Email:     r.Email,
		StudentID: r.StudentID,
		Password:  r.Password,
	}
}
