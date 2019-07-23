package api

import (
	_ "encoding/json"
	_ "fmt"
	"net/http"
	"testing"

	"github.com/GruffDebate/server/gruff"
	"github.com/stretchr/testify/assert"
)

func createUser(name string, username string, email string) gruff.User {
	u := gruff.User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: "123456",
	}

	u.Create(CTX)

	return u
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

func TestLogin(t *testing.T) {
	setup()
	defer teardown()

	r := New(nil)

	createUser("test1", "test1", "test1@test1.com")

	u := map[string]interface{}{
		"email": "test1@test1.com",
		"password": "123456",
	}

	r.POST("/api/auth")
	r.SetBody(u)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}


/*
func TestListUsers(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createUser("test1", "test1", "test1@test1.com")
	u2 := createUser("test2", "test2", "test2@test2.com")

	expectedResults, _ := json.Marshal([]gruff.User{u1, u2})

	r.GET("/api/users")
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListClaimsByUser(t *testing.T) {
	setup()
	defer teardown()

	u1 := createUser("test1", "test1", "test1@test1.com")

	r := New(tokenForTestUser(u1))

	c1 := gruff.Claim{
		Identifier:  gruff.Identifier{CreatedByID: u1.ArangoID()},
		Title:       "Claim 1",
		Description: "Claim 1",
		Truth:       0.23094,
	}
	c2 := gruff.Claim{
		Identifier:  gruff.Identifier{CreatedByID: u1.ArangoID()},
		Title:       "Claim 2",
		Description: "Claim 2",
		Truth:       0.23094,
	}
	c3 := gruff.Claim{
		Title:       "Claim 3",
		Description: "Claim 3",
		Truth:       0.25094,
	}
	c4 := gruff.Claim{
		Title:       "Claim 4",
		Description: "Claim 4",
		Truth:       0.26094,
	}
	c1.Create(CTX)
	c2.Create(CTX)
	c3.Create(CTX)
	c4.Create(CTX)

	expectedResults, _ := json.Marshal([]gruff.Claim{c1, c2})

	r.GET("/api/users/claims")
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListUsersPagination(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	createUser("test1", "test1", "test1@test1.com")
	createUser("test2", "test2", "test2@test2.com")

	r.GET("/api/users?start=0&limit=25")
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetUsers(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createUser("test1", "test1", "test1@test1.com")

	expectedResults, _ := json.Marshal(u1)

	r.GET(fmt.Sprintf("/api/users/%s", u1.ArangoKey()))
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetUserMe(t *testing.T) {
	setup()
	defer teardown()

	u1 := createUser("test1", "test1", "test1@test1.com")

	r := New(tokenForTestUser(u1))

	createUser("test2", "test2", "test2@test2.com")

	expectedResults, _ := json.Marshal(u1)

	r.GET("/api/users/me")
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

	r.PUT(fmt.Sprintf("/api/users/%s", u1.ArangoKey()))
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)
}

func TestUpdateUserMe(t *testing.T) {
	setup()
	defer teardown()

	u1 := createUser("test1", "test1", "test1@test1.com")

	r := New(tokenForTestUser(u1))

	createUser("test2", "test2", "test2@test2.com")

	u1.Email = "test1010@test1010.com"
	expectedResults, _ := json.Marshal(u1)

	r.PUT("/api/users/me")
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestDeleteUsers(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	u1 := createUser("test1", "test1", "test1@test1.com")

	r.DELETE(fmt.Sprintf("/api/users/%s", u1.ArangoKey()))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

*/
