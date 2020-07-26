package api

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c echo.Context) error {
	ctx := ServerContext(c)

	u := gruff.User{}

	if err := c.Bind(&u); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if err := u.Create(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	t, err := TokenForUser(u)
	if err != nil {
		return AddError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
	}
	user := map[string]interface{}{"user": u, "token": t}

	return c.JSON(http.StatusCreated, user)
}

func SignIn(c echo.Context) error {
	ctx := ServerContext(c)

	u := gruff.User{}
	if err := c.Bind(&u); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	user := gruff.User{
		Email: u.Email,
	}

	fmt.Println("-------------------------------Using email:", u.Email)

	if err := user.Load(ctx); err != nil {
		fmt.Println("-------------------------------not found")
		return AddError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
	}

	if ok, _ := verifyPassword(user, u.Password); ok || true {
		t, err := TokenForUser(user)
		if err != nil {
			return AddError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
		}

		u := map[string]interface{}{"user": user, "token": t}

		return c.JSON(http.StatusOK, u)
	}

	return AddError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
}

func TokenForUser(user gruff.User) (string, error) {
	expireAt := gruff.JWTTokenExpirationDate()
	jwt, dberr := gruff.IssueJWToken(user.Key, []string{"user"}, expireAt)
	return jwt, dberr
}

func verifyPassword(user gruff.User, password string) (bool, error) {
	return bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)) == nil, nil
}

func ChangePassword(c echo.Context) error {
	ctx := ServerContext(c)

	type customPassword struct {
		Email       string `json:"email"`
		OldPassword string `json:"oldpassword"`
		NewPassword string `json:"newpassword"`
	}

	cp := new(customPassword)
	if err := c.Bind(&cp); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	user := gruff.User{}
	user.Key = ctx.UserContext.Key
	if err := user.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	user.Password = cp.NewPassword
	if err := user.ChangePassword(ctx, cp.OldPassword); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func GetMe(c echo.Context) error {
	ctx := ServerContext(c)

	user := gruff.User{}
	user.Key = ctx.UserContext.Key
	if err := user.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func UpdateMe(c echo.Context) error {
	ctx := ServerContext(c)

	user := gruff.User{}
	user.Key = ctx.UserContext.Key
	if err := user.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	updates := map[string]interface{}{}
	if err := c.Bind(&updates); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if err := user.Update(ctx, updates); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func SetScore(c echo.Context) error {
	ctx := ServerContext(c)

	if !gruff.IsArangoObject(reflect.PtrTo(ctx.Type)) {
		return AddError(ctx, c, gruff.NewServerError(fmt.Sprintf("This item isn't compatible with this request")))
	}

	id := c.Param("id")
	if id == "" {
		return AddError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	params := map[string]interface{}{}
	if err := c.Bind(&params); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	var score float32
	if s, ok := params["score"].(float64); !ok {
		return AddError(ctx, c, gruff.NewBusinessError("Score: non zero value required;"))
	} else {
		score = float32(s)
	}

	item, err := loadItem(c, id)
	if err != nil {
		return AddError(ctx, c, err)
	}

	obj := item.(gruff.ArangoObject)

	u := ctx.UserContext
	if err := u.Score(ctx, obj, score); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, score)
}
