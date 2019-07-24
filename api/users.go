package api

import (
	"net/http"
	"fmt"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c echo.Context) error {
	ctx := ServerContext(c)

	u := gruff.User{}

	if err := c.Bind(&u); err != nil {
		fmt.Println(err)
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if err := u.Create(ctx); err != nil {
		return AddGruffError(ctx, c, err)
	}

	t, err := TokenForUser(u)
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
	}
	user := map[string]interface{}{"user": u, "token": t}

	return c.JSON(http.StatusCreated, user)
}

func SignIn(c echo.Context) error {
	ctx := ServerContext(c)

	u := gruff.User{}
	if err := c.Bind(&u); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	user := gruff.User{
		Email: u.Email,
	}

	if err := user.Load(ctx); err != nil {
		return AddGruffError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
	}
	
	if ok, _ := verifyPassword(user, u.Password); ok {
		t, err := TokenForUser(user)
		if err != nil {
			return AddGruffError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
		}

		u := map[string]interface{}{"user": user, "token": t}

		return c.JSON(http.StatusOK, u)
	}

	return AddGruffError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
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
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	user := gruff.User{}
	user.Key = ctx.UserContext.Key
	if err := user.Load(ctx); err != nil {
		return AddGruffError(ctx, c, err)
	}

	user.Password = cp.NewPassword
	if err := user.ChangePassword(ctx, cp.OldPassword); err != nil {
		return AddGruffError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func GetMe(c echo.Context) error {
	ctx := ServerContext(c)

	user := gruff.User{}
	user.Key = ctx.UserContext.Key
	if err := user.Load(ctx); err != nil {
		return AddGruffError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func UpdateMe(c echo.Context) error {
	ctx := ServerContext(c)

	user := gruff.User{}
	user.Key = ctx.UserContext.Key
	if err := user.Load(ctx); err != nil {
		return AddGruffError(ctx, c, err)
	}

	updates := map[string]interface{}{}
	if err := c.Bind(&updates); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if err := user.Update(ctx, updates); err != nil {
		return AddGruffError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, user)
}

/*
func ListClaimsUser(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	claims := []gruff.Claim{}

	db = BasicJoins(ctx, c, db)
	db = db.Where("created_by_id = ?", ctx.UserContext.ID)
	db = BasicPaging(ctx, c, db)

	err := db.Find(&claims).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if ctx.Payload["ct"] != nil {
		ctx.Payload["results"] = claims
		return c.JSON(http.StatusOK, ctx.Payload)
	}

	return c.JSON(http.StatusOK, claims)
}
*/
