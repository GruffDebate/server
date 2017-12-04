package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bigokro/gruff-server/gruff"
	"github.com/stretchr/testify/assert"
)

func TestListContexts(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createContext()
	u2 := createContext()
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	expectedResults, _ := json.Marshal([]gruff.Context{u1, u2})

	r.GET("/api/contexts")
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListContextsPagination(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createContext()
	u2 := createContext()
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	r.GET("/api/contexts?start=0&limit=25")
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetContexts(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createContext()
	TESTDB.Create(&u1)

	expectedResults, _ := json.Marshal(u1)

	r.GET(fmt.Sprintf("/api/contexts/%d", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestCreateContexts(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createContext()

	r.POST("/api/contexts")
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestUpdateContexts(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createContext()
	TESTDB.Create(&u1)

	r.PUT(fmt.Sprintf("/api/contexts/%d", u1.ID))
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)
}

func TestDeleteContexts(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	u1 := createContext()
	TESTDB.Create(&u1)

	r.DELETE(fmt.Sprintf("/api/contexts/%d", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestAddContextToClaim(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	c1 := createContext()
	TESTDB.Create(&c1)

	cl1 := gruff.Claim{Title: "We should troll the trolls"}
	TESTDB.Create(&cl1)

	r.POST(fmt.Sprintf("/api/claims/%s/contexts/%d", cl1.ID, c1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	TESTDB.Preload("Contexts").Where("id = ?", cl1.ID).First(&cl1)
	assert.Equal(t, 1, len(cl1.Contexts))
	assert.Equal(t, c1.ID, cl1.Contexts[0].ID)
}

func TestRemoveContextFromClaim(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	c1 := createContext()
	TESTDB.Create(&c1)

	cl1 := gruff.Claim{Title: "We should troll the trolls"}
	TESTDB.Create(&cl1)

	TESTDB.Model(&cl1).Association("Contexts").Append(&c1)

	r.DELETE(fmt.Sprintf("/api/claims/%s/contexts/%d", cl1.ID, c1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	cl1.Contexts = []gruff.Context{}
	TESTDB.Preload("Contexts").Where("id = ?", cl1.ID).First(&cl1)
	assert.Equal(t, 0, len(cl1.Contexts))
}

func createContext() gruff.Context {
	c := gruff.Context{
		Title:       "contexts",
		Description: "contexts",
		Url:         "https://en.wikipedia.org/wiki/Peter_Christen_Asbj%C3%B8rnsen",
	}

	return c
}
