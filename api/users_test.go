package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bigokro/gruff-server/gruff"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestListUsers(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createUser("test1", "test1", "test1@test1.com")
	u2 := createUser("test2", "test2", "test2@test2.com")
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	expectedResults, _ := json.Marshal([]gruff.User{u1, u2})

	r.GET("/api/users")
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListUsersPagination(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createUser("test1", "test1", "test1@test1.com")
	u2 := createUser("test2", "test2", "test2@test2.com")
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	r.GET("/api/users?start=0&limit=25")
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetUsers(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createUser("test1", "test1", "test1@test1.com")
	TESTDB.Create(&u1)

	expectedResults, _ := json.Marshal(u1)

	r.GET(fmt.Sprintf("/api/users/%d", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestCreateUsers(t *testing.T) {
	setup()
	defer teardown()

	r := New(nil)

	u1 := createUser("test1", "test1", "test1@test1.com")

	r.POST("/api/users")
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestUpdateUsers(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createUser("test1", "test1", "test1@test1.com")
	TESTDB.Create(&u1)

	r.PUT(fmt.Sprintf("/api/users/%d", u1.ID))
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)
}

func TestDeleteUsers(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	u1 := createUser("test1", "test1", "test1@test1.com")
	TESTDB.Create(&u1)

	r.DELETE(fmt.Sprintf("/api/users/%d", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func createUser(name string, username string, email string) gruff.User {
	u := gruff.User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: "123456",
	}
	password := u.Password
	u.Password = ""
	u.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return u
}
