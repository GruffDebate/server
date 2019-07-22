package api

import (
	_ "encoding/json"
	_ "fmt"
	_ "net/http"
	_ "testing"

	_ "github.com/GruffDebate/server/gruff"
	_ "github.com/stretchr/testify/assert"
)

/*
func TestListLinks(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createLink()
	u2 := createLink()
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	expectedResults, _ := json.Marshal([]gruff.Link{u1, u2})

	r.GET("/api/links")
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListLinksPagination(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createLink()
	u2 := createLink()
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	r.GET("/api/links?start=0&limit=25")
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetLink(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createLink()
	TESTDB.Create(&u1)

	expectedResults, _ := json.Marshal(u1)

	fmt.Printf("/api/links/%s\n", u1.ID)
	r.GET(fmt.Sprintf("/api/links/%s", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestCreateLink(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createLink()

	r.POST("/api/links")
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestUpdateLink(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createLink()
	TESTDB.Create(&u1)

	r.PUT(fmt.Sprintf("/api/links/%s", u1.ID))
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)
}

func TestDeleteLink(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	u1 := createLink()
	TESTDB.Create(&u1)

	r.DELETE(fmt.Sprintf("/api/links/%s", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func createLink() gruff.Link {
	l := gruff.Link{
		Title:       "Links",
		Description: "Links",
		Url:         "www.gruff.org",
	}

	return l
}
*/
