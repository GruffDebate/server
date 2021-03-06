package gruff

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Name            string     `json:"name" sql:"not null" valid:"length(3|50)"`
	Username        string     `json:"username" settable:"false" sql:"unique_index;not null" valid:"length(3|50),matches(^[a-zA-Z0-9][a-zA-Z0-9-_]+$),required"`
	Email           string     `json:"email" sql:"not null" valid:"email"`
	Password        string     `json:"password,omitempty" sql:"-" valid:"length(5|64)"`
	HashedPassword  string     `json:"hashed_password"` // TODO: don't return this value via the API
	Image           string     `json:"img,omitempty"`
	Curator         bool       `json:"curator"`
	Admin           bool       `json:"admin"`
	URL             string     `json:"url,omitempty"`
	EmailVerifiedAt *time.Time `json:"-" settable:"false"`
}

// ArangoObject interface

func (u User) CollectionName() string {
	return "users"
}

func (u User) ArangoKey() string {
	return u.Key
}

func (u User) ArangoID() string {
	return fmt.Sprintf("%s/%s", u.CollectionName(), u.ArangoKey())
}

func (u User) DefaultQueryParameters() ArangoQueryParameters {
	return DEFAULT_QUERY_PARAMETERS
}

func (u *User) Create(ctx *ServerContext) Error {
	col, err := ctx.Arango.CollectionFor(u)
	if err != nil {
		return err

	}

	// TODO: Test
	can, err := u.UserCanCreate(ctx)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to create this item")
	}

	u.PrepareForCreate(ctx)

	if err := u.ValidateForCreate(); err != nil {
		return err
	}

	password := u.Password
	u.Password = ""
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	u.HashedPassword = string(hashedPassword[:])
	if _, dberr := col.CreateDocument(ctx.Context, u); dberr != nil {
		return NewServerError(dberr.Error())
	}
	return nil
}

// TODO: Test
func (u *User) Update(ctx *ServerContext, updates Updates) Error {
	return UpdateArangoObject(ctx, u, updates)
}

// TODO: Test
func (u *User) Delete(ctx *ServerContext) Error {
	return DeleteArangoObject(ctx, u)
}

// Restrictor
// TODO: Test
// TODO: Call in CRUD and other methods
func (u User) UserCanView(ctx *ServerContext) (bool, Error) {
	user := ctx.UserContext
	if user.Curator {
		return true, nil
	}
	return u.ArangoKey() == user.ArangoKey(), nil
}

func (u User) UserCanCreate(ctx *ServerContext) (bool, Error) {
	return true, nil
}

func (u User) UserCanUpdate(ctx *ServerContext, updates Updates) (bool, Error) {
	return u.UserCanView(ctx)
}

func (u User) UserCanDelete(ctx *ServerContext) (bool, Error) {
	user := ctx.UserContext
	return user.Curator, nil
}

// Validator

func (u User) ValidateForCreate() Error {
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

func (u User) ValidateForUpdate(updates Updates) Error {
	updated := User{
		Name:     updates["name"].(string),
		Email:    updates["email"].(string),
		Username: updates["username"].(string),
	}
	if updated.Name != "" {
		if err := updated.ValidateField("Name"); err != nil {
			return err
		}
	}
	if updated.Email != "" {
		if err := updated.ValidateField("Email"); err != nil {
			return err
		}
	}
	if updated.Username != "" {
		if err := updated.ValidateField("Username"); err != nil {
			return err
		}
	}
	return nil
}

func (u User) ValidateForDelete() Error {
	return nil
}

func (u User) ValidateField(f string) Error {
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

// Loader

func (u *User) Load(ctx *ServerContext) Error {
	var err Error
	if u.ArangoKey() != "" {
		err = LoadArangoObject(ctx, u, u.ArangoKey())
	} else {
		var query string
		bindVars := BindVars{}
		if u.Username != "" {
			bindVars["username"] = strings.ToLower(u.Username)
			// TODO: unique index on lower(username)
			query = fmt.Sprintf("FOR obj IN %s FILTER LOWER(obj.username) == @username LIMIT 1 RETURN obj", u.CollectionName())
		} else if u.Email != "" {
			bindVars["email"] = strings.ToLower(u.Email)
			// TODO: unique index on lower(email)
			query = fmt.Sprintf("FOR obj IN %s FILTER LOWER(obj.email) == @email LIMIT 1 RETURN obj", u.CollectionName())
		} else {
			return NewBusinessError("There is no value available to load this User.")
		}
		err = FindArangoObject(ctx, query, bindVars, u)
	}
	return err
}

func (u *User) LoadFull(ctx *ServerContext) Error {
	if err := u.Load(ctx); err != nil {
		return err
	}
	return nil
}

// Business methods

func (u *User) VerifyPassword(ctx *ServerContext, password string) (bool, Error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		return false, NewBusinessError(err.Error())
	}
	return true, nil
}

func (u *User) ChangePassword(ctx *ServerContext, oldPassword string) Error {
	col, err := ctx.Arango.CollectionFor(u)
	if err != nil {
		return err
	}

	if u.Password == "" {
		return NewBusinessError("New Password: non zero value required;")
	}
	newPassword := u.Password
	u.Password = ""
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	u.HashedPassword = string(hashedPassword[:])
	update := Updates{
		"hashed_password": u.HashedPassword,
	}

	if _, err := col.UpdateDocument(ctx.Context, u.ArangoKey(), update); err != nil {
		return NewServerError(err.Error())
	}

	return nil
}

// Scoring

func (u User) Score(ctx *ServerContext, target ArangoObject, score float32) Error {
	oldScore, err := u.ScoreFor(ctx, target)
	if err != nil {
		return err
	}

	if oldScore != nil {
		if err := oldScore.Delete(ctx); err != nil {
			return err
		}
	}

	newScore := UserScore{
		Edge: Edge{
			From: u.ArangoID(),
			To:   target.ArangoID(),
		},
		Score: score,
	}

	if err := newScore.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	//go target.UpdateScores(ctx)
	if scorer, ok := target.(Scorer); ok {
		scorer.UpdateScore(ctx)
	}

	return nil
}

func (u *User) ScoreFor(ctx *ServerContext, target ArangoObject) (*UserScore, Error) {
	score := UserScore{}
	bindVars := BindVars{
		"user": u.ArangoID(),
	}
	var dateFilter string
	if claim, ok := target.(*Claim); ok {
		dateFilter = claim.DateFilter(bindVars)
		bindVars["target"] = claim.ID
	} else if arg, ok := target.(*Argument); ok {
		dateFilter = arg.DateFilter(bindVars)
		bindVars["target"] = arg.ID
	}
	query := fmt.Sprintf(`FOR obj IN %s
                                 FOR targ IN %s
                                   FILTER obj._to == targ._id
                                      AND obj._from == @user
                                      AND targ.id == @target
                                   %s
                                   SORT obj.start ASC
                                   RETURN obj`,
		UserScore{}.CollectionName(),
		target.CollectionName(),
		dateFilter,
	)
	err := FindArangoObject(ctx, query, bindVars, &score)
	if err != nil {
		if err.Code() == ERROR_CODE_NOT_FOUND {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &score, nil
}
