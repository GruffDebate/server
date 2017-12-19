package gruff

import (
	"strings"
	"time"
)

type User struct {
	Model
	Name            string     `json:"name" sql:"not null" valid:"length(3|50)"`
	Username        string     `json:"username" settable:"false" sql:"unique_index;not null" valid:"length(3|50),matches(^[a-zA-Z0-9][a-zA-Z0-9-_]+$),required"`
	Email           string     `json:"email" sql:"not null" valid:"email"`
	Password        string     `json:"password,omitempty" sql:"-" valid:"length(5|64)"`
	HashedPassword  []byte     `json:"-" sql:"hashed_password;not null" gorm:"size:32"`
	Image           string     `json:"img,omitempty"`
	Curator         bool       `json:"curator"`
	Admin           bool       `json:"admin"`
	URL             string     `json:"url,omitempty"`
	EmailVerifiedAt *time.Time `json:"-" settable:"false"`
}

func (u User) ValidateForCreate() GruffError {
	err := u.ValidateField("Name")
	if err != nil {
		return err
	}
	err = u.ValidateField("Email")
	if err != nil {
		return err
	}
	err = u.ValidateField("Username")
	if err != nil {
		return err
	}
	err = u.ValidateField("Password")
	if err != nil {
		return err
	}
	return nil
}

func (u User) ValidateForUpdate() GruffError {
	return nil
}

func (u User) ValidateField(f string) GruffError {
	data := map[string]interface{}{"field": f}

	err := ValidateStructField(u, f)
	if err != nil {
		switch f {
		case "Email":
			if strings.Contains(err.Error(), "validate as email") {
				err = NewBusinessError(err.Error(), ERROR_SUBCODE_EMAIL_FORMAT, data)
			}
		case "Username":
			if strings.Contains(err.Error(), "validate as matches") {
				err = NewBusinessError(err.Error(), ERROR_SUBCODE_USERNAME_FORMAT, data)
			} else if strings.Contains(err.Error(), "validate as length") {
				err = NewBusinessError(err.Error(), ERROR_SUBCODE_USERNAME_LENGTH, data)
			}
		case "Password":
			if strings.Contains(err.Error(), "validate as password") {
				err = NewBusinessError(err.Error(), ERROR_SUBCODE_PASSWORD_FORMAT, data)
			} else if strings.Contains(err.Error(), "validate as length") {
				err = NewBusinessError(err.Error(), ERROR_SUBCODE_PASSWORD_LENGTH, data)
			}
		}
	}
	return err
}
