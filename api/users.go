package api

import (
	"net/http"
	"time"

	"github.com/GruffDebate/server/gruff"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    string `json:"img"`
	Curator  bool   `json:"curator"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

type customPassword struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

func SignUp(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	u := new(gruff.User)

	if err := c.Bind(u); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	password := u.Password
	u.Password = ""
	u.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err := db.Create(u).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	return c.JSON(http.StatusCreated, u)
}

func SignIn(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	u := gruff.User{}
	if err := c.Bind(&u); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if u.Email != "" {
		user := gruff.User{}
		if err := db.Where("email = ?", u.Email).Find(&user).Error; err != nil {
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
	}

	return AddGruffError(ctx, c, gruff.NewUnauthorizedError("Unauthorized"))
}

func TokenForUser(user gruff.User) (string, error) {
	claims := &jwtCustomClaims{
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.Image,
		user.Curator,
		user.Admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	return t, err
}

func verifyPassword(user gruff.User, password string) (bool, error) {
	return bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)) == nil, nil
}

func ChangePassword(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	u := new(customPassword)
	if err := c.Bind(&u); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	user := gruff.User{}
	err := db.Where("email = ?", u.Email).Find(&user).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if ok, _ := verifyPassword(user, u.OldPassword); ok {
		user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(u.NewPassword), bcrypt.DefaultCost)

		if err := db.Save(user).Error; err != nil {
			return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
		}

		return c.JSON(http.StatusOK, user)
	}

	return AddGruffError(ctx, c, gruff.NewNotFoundError("Not Found"))
}
