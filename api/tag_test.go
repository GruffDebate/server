package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/GruffDebate/server/gruff"
	"github.com/stretchr/testify/assert"
)

func TestListTags(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createTag()
	u2 := createTag()
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	expectedResults, _ := json.Marshal([]gruff.Tag{u1, u2})

	r.GET("/api/tags")
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListTagsPagination(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createTag()
	u2 := createTag()
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)

	r.GET("/api/tags?start=0&limit=25")
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetTags(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createTag()
	TESTDB.Create(&u1)

	expectedResults, _ := json.Marshal(u1)

	r.GET(fmt.Sprintf("/api/tags/%d", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestCreateTags(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createTag()

	r.POST("/api/tags")
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestUpdateTags(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	u1 := createTag()
	TESTDB.Create(&u1)

	r.PUT(fmt.Sprintf("/api/tags/%d", u1.ID))
	r.SetBody(u1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)
}

func TestDeleteTags(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	u1 := createTag()
	TESTDB.Create(&u1)

	r.DELETE(fmt.Sprintf("/api/tags/%d", u1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListClaimsByTags(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	t1 := createTag()
	TESTDB.Create(&t1)

	c1 := gruff.Claim{Title: "Claim 1", Description: "Claim 1", Truth: 0.23094}
	c2 := gruff.Claim{Title: "Claim 2", Description: "Claim 2", Truth: 0.23094}
	c3 := gruff.Claim{Title: "Claim 3", Description: "Claim 3", Truth: 0.23094}
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)
	TESTDB.Create(&c3)

	TESTDB.Exec("INSERT INTO claim_tags (tag_id, claim_id) VALUES (?, ?)", t1.ID, c1.ID)
	TESTDB.Exec("INSERT INTO claim_tags (tag_id, claim_id) VALUES (?, ?)", t1.ID, c2.ID)

	expectedResults, _ := json.Marshal([]gruff.Claim{c1, c2})

	r.GET(fmt.Sprintf("/api/tags/%d/claims", t1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func createTag() gruff.Tag {
	t := gruff.Tag{
		Title: "Tag",
	}

	return t
}
