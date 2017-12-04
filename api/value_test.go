package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bigokro/gruff-server/gruff"
	"github.com/stretchr/testify/assert"
)

func TestListValues(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createValue()
	u2 := createValue()
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	expectedResults, _ := json.Marshal([]gruff.Value{u1, u2})

	r.GET("/api/values")
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListValuesPagination(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createValue()
	u2 := createValue()
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	r.GET("/api/values?start=0&limit=25")
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetValues(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createValue()
	TESTDB.Create(&u1)

	expectedResults, _ := json.Marshal(u1)

	r.GET(fmt.Sprintf("/api/values/%d", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestCreateValues(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createValue()

	r.POST("/api/values")
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestUpdateValues(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createValue()
	TESTDB.Create(&u1)

	r.PUT(fmt.Sprintf("/api/values/%d", u1.ID))
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)
}

func TestDeleteValues(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	u1 := createValue()
	TESTDB.Create(&u1)

	r.DELETE(fmt.Sprintf("/api/values/%d", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func createValue() gruff.Value {
	v := gruff.Value{
		Title:       "Value",
		Description: "Value",
	}

	return v
}
