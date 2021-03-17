package user

import (
	"errors"
	"fmt"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/authentication"
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"

	"golang.org/x/crypto/bcrypt"
)

// Generate the JWT auth token
func storeAuth(uid int) (*authentication.Token, error) {
	token, err := authentication.CreateToken(uid)
	if err != nil {
		return nil, err
	}

	err = authentication.CreateAuth(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// Compare the password form and database if match
func compareHashPassword(hashed, password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(hashed)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

// Login ...
// The Email, (StudentID or Email) must be set
func (u User) Login() (*User, *authentication.Token, error) {
	var user User
	var err error

	if u.Email != nil {
		sql, _ := Sqls{}.Get(ByEmail)
		err = database.GetDB().Get(&user, sql, *u.Email)
	} else if u.StudentID != nil {
		sql, _ := Sqls{}.Get(ByStudentID)
		err = database.GetDB().Get(&user, sql, *u.StudentID)
	} else {
		return nil, nil, fmt.Errorf("The email or student ID must be send on login")
	}

	if err != nil {
		return nil, nil, fmt.Errorf("User not found")
	}

	if compareHashPassword(*user.Password, *u.Password) != nil {
		return nil, nil, errors.New("Invalid password")
	}

	token, err := storeAuth(*user.ID)

	return &user, token, nil
}

// Register ...
func (u User) Register() (*User, *authentication.Token, error) {
	user, err := u.AddAndCheck()
	if err != nil {
		return nil, nil, err
	}

	token, err := storeAuth(*user.ID)

	return user, token, nil
}

// Get ...
func Get(userID int) (user User, err error) {
	sql, _ := Sqls{}.Get(ByUserID)
	err = database.GetDB().Get(&user, sql, userID)
	return user, err
}

// AddAndCheck ...
func (u *User) AddAndCheck() (*User, error) {
	getDb := database.GetDB()

	user := new(User)
	sql, _ := Sqls{}.Get(ByEmail + ByStudentID)
	err := getDb.Get(user, sql, u.Email, u.StudentID)

	// Check if the user exists in database
	if err == nil {
		return nil, errors.New("User already exists")
	}

	bytePassword := []byte(*u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err) // Something really went wrong here...
	}

	sql, _ = Sqls{}.Add()

	// Create the user and return back the user ID
	err = getDb.Get(user, sql, u.Name, u.Email, u.Role, u.StudentID, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddOrGet ...
func (u *User) AddOrGet() error {
	getDb := database.GetDB()

	sql, _ := Sqls{}.Get(ByEmail + ByStudentID)
	err := getDb.Get(u, sql, u.Email, u.StudentID)

	// Check if the user exists in database
	if err == nil {
		return nil
	}

	bytePassword := []byte(*u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err) // Something really went wrong here...
	}

	sql, _ = Sqls{}.Add()

	// Create the user and return back the user ID
	err = getDb.Get(u, sql, u.Name, u.Email, u.Role, u.StudentID, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
