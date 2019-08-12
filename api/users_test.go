package api

import (
	"fmt"
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
	u.Load(CTX)
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
		"email":    "test1@test1.com",
		"password": "123456",
	}

	r.POST("/api/auth")
	r.SetBody(u)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestListClaimsByUser(t *testing.T) {
	setup()
	defer teardown()

	u1 := createUser("test1", "test1", "test1@test1.com")
	u1.Load(CTX)

	r := New(tokenForTestUser(u1))

	c1 := gruff.Claim{
		Title:        "Let's create a new claim 1",
		Description:  "Claims in general should be true or false",
		Negation:     "Let's not...",
		Question:     "Should we create a new Claim?",
		Note:         "He who notes is a note taker",
		Image:        "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
		MultiPremise: true,
		PremiseRule:  gruff.PREMISE_RULE_ALL,
	}
	c2 := gruff.Claim{
		Title:        "Let's create a new claim 2",
		Description:  "Claims in general should be true or false",
		Negation:     "Let's not...",
		Question:     "Should we create a new Claim?",
		Note:         "He who notes is a note taker",
		Image:        "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
		MultiPremise: true,
		PremiseRule:  gruff.PREMISE_RULE_ALL,
	}
	c3 := gruff.Claim{
		Title:        "Let's create a new claim 3",
		Description:  "Claims in general should be true or false",
		Negation:     "Let's not...",
		Question:     "Should we create a new Claim?",
		Note:         "He who notes is a note taker",
		Image:        "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
		MultiPremise: true,
		PremiseRule:  gruff.PREMISE_RULE_ALL,
	}
	c4 := gruff.Claim{
		Title:        "Let's create a new claim 4",
		Description:  "Claims in general should be true or false",
		Negation:     "Let's not...",
		Question:     "Should we create a new Claim?",
		Note:         "He who notes is a note taker",
		Image:        "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
		MultiPremise: true,
		PremiseRule:  gruff.PREMISE_RULE_ALL,
	}

	c4.Create(CTX)

	CTX.UserContext = u1
	c1.Create(CTX)
	c2.Create(CTX)
	c3.Create(CTX)

	c1.Load(CTX)
	c2.Load(CTX)
	c3.Load(CTX)
	c4.Load(CTX)

	// expectedResults, _ := json.Marshal([]gruff.Claim{c3, c2, c1})

	r.GET("/api/users/claims")
	res, _ := r.Run(Router())
	// assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestSetScore(t *testing.T) {
	setup()
	defer teardown()

	u := CTX.UserContext

	claim := gruff.Claim{
		Title: "Dude, I totally scored!",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	arg := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "Scored? Like, left scratch marks?",
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)

	body := map[string]interface{}{
		"score": 0.55,
	}

	r := New(tokenForTestUser(u))

	r.POST(fmt.Sprintf("/api/claims/%s/score", claim.ID))
	r.SetBody(body)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "0.55\n", res.Body.String())

	// TODO: Assert that total score was updated
	err = claim.Load(CTX)
	assert.NoError(t, err)

	score, err := u.ScoreFor(CTX, &claim)
	assert.NoError(t, err)
	assert.NotNil(t, score)
	assert.Equal(t, float32(0.55), score.Score)

	body = map[string]interface{}{
		"score": 0.22,
	}

	r.POST(fmt.Sprintf("/api/arguments/%s/score", arg.ID))
	r.SetBody(body)
	res, _ = r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "0.22\n", res.Body.String())

	// TODO: Assert that total score was updated
	err = arg.Load(CTX)
	assert.NoError(t, err)

	score, err = u.ScoreFor(CTX, &arg)
	assert.NoError(t, err)
	assert.NotNil(t, score)
	assert.Equal(t, float32(0.22), score.Score)
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
