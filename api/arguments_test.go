package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/GruffDebate/server/gruff"
	"github.com/stretchr/testify/assert"
)

func TestGetArgument(t *testing.T) {
	setup()
	defer teardown()

	claim := gruff.Claim{
		Title:       "This is the API Get Argument test claim",
		Description: "Load all the things!",
		Negation:    "Don't load all the things.",
		Question:    "Load all the THINGS? Load ALL the things? LOAD all the things?",
		Note:        "This Claim needs to be all loaded.",
		Image:       "https://i.chzbgr.com/full/6434679808/h4ADBDEA5/",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	targetClaim := gruff.Claim{
		Title: "I'm the target of an API Get Argument",
	}
	err = targetClaim.Create(CTX)
	assert.NoError(t, err)

	arg := gruff.Argument{
		Title:         "The API needs to Get an Argument",
		Description:   "And that's me!",
		TargetClaimID: &targetClaim.ID,
		ClaimID:       claim.ID,
		Pro:           true,
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)

	arg1 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "Get Arg All?",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)

	arg2 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "GET Arg Load ALL!",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)

	arg3 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "Do it Get Arg Load ALL!",
		Pro:           true,
	}
	err = arg3.Create(CTX)
	assert.NoError(t, err)

	context := gruff.Context{
		ShortName: "Get Arg context",
		Title:     "Get Arg context",
		URL:       "https://en.wikipedia.org/wiki/Idealism",
	}
	err = context.Create(CTX)
	assert.NoError(t, err)

	err = claim.AddContext(CTX, context)
	assert.NoError(t, err)

	err = arg.LoadFull(CTX)
	assert.NoError(t, err)

	expected, _ := json.Marshal(arg)

	r := New(tokenForTestUser(DEFAULT_USER))
	r.GET(fmt.Sprintf("/api/arguments/%s", arg.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, string(expected), res.Body.String())
}

func TestMoveArgument(t *testing.T) {
	setup()
	defer teardown()

	claim := gruff.Claim{
		Title:       "This is the API Move Argument test claim",
		Description: "Load all the things!",
		Negation:    "Don't load all the things.",
		Question:    "Load all the THINGS? Load ALL the things? LOAD all the things?",
		Note:        "This Claim needs to be all loaded.",
		Image:       "https://i.chzbgr.com/full/6434679808/h4ADBDEA5/",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	targetClaim := gruff.Claim{
		Title: "I'm the original target of an API Move Argument",
	}
	err = targetClaim.Create(CTX)
	assert.NoError(t, err)

	targetArg := gruff.Argument{
		TargetClaimID: &targetClaim.ID,
		Title:         "I'm the final target of an API Move Argument",
		Pro:           true,
	}
	err = targetArg.Create(CTX)
	assert.NoError(t, err)

	arg := gruff.Argument{
		Title:         "The API needs to Move an Argument",
		Description:   "And that's me!",
		TargetClaimID: &targetClaim.ID,
		ClaimID:       claim.ID,
		Pro:           true,
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)

	arg1 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "Move Arg All?",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)

	arg2 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "GET Arg Load ALL!",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)

	arg3 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "Do it Move Arg Load ALL!",
		Pro:           true,
	}
	err = arg3.Create(CTX)
	assert.NoError(t, err)

	context := gruff.Context{
		ShortName: "Move Arg context",
		Title:     "Move Arg context",
		URL:       "https://en.wikipedia.org/wiki/Idealism",
	}
	err = context.Create(CTX)
	assert.NoError(t, err)

	err = claim.AddContext(CTX, context)
	assert.NoError(t, err)

	body := map[string]interface{}{
		"_key": arg.ArangoKey(),
		"pro":  false,
	}

	r := New(tokenForTestUser(DEFAULT_USER))
	r.PUT(fmt.Sprintf("/api/arguments/%s/move/arguments/%s", arg.ID, targetArg.ID))
	r.SetBody(body)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	// TODO: Make a Reload/ReloadFull method(s)
	arg.Key = ""
	err = arg.LoadFull(CTX)
	assert.NoError(t, err)
	expected, _ := json.Marshal(arg)
	assert.JSONEq(t, string(expected), res.Body.String())

	assert.Equal(t, targetArg.ID, *arg.TargetArgumentID)
	assert.Nil(t, arg.TargetClaimID)
	assert.Equal(t, claim.ID, arg.ClaimID)
	assert.False(t, arg.Pro)

	inference, err := arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), inference.To)
	assert.Equal(t, targetArg.ArangoID(), inference.From)
}

/*
func TestCreateArgumentForClaim(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	d1 := gruff.Claim{
		Title:       "Claim",
		Description: "This is a test Claim",
		Truth:       0.001,
	}
	d2 := gruff.Claim{
		Title:       "Another Claim",
		Description: "This a target claim",
		Truth:       1.000,
	}
	TESTDB.Create(&d1)
	TESTDB.Create(&d2)

	a1 := gruff.Argument{
		ClaimID:       d1.ID,
		TargetClaimID: &gruff.NullableUUID{UUID: d2.ID},
		Type:          gruff.ARGUMENT_AGAINST,
		Title:         "This is an argument",
		Description:   "This is an arguous description",
	}

	r.POST("/api/arguments")
	r.SetBody(a1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	TESTDB.Where("title = ?", a1.Title).First(&a1)
	expectedResults, _ := json.Marshal(a1)

	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, gruff.ARGUMENT_AGAINST, a1.Type)
	assert.Equal(t, d1.ID, a1.ClaimID)
	assert.Equal(t, d2.ID, a1.TargetClaimID.UUID)
}

func TestCreateArgumentForArgument(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	d1 := gruff.Claim{
		Title:       "Claim",
		Description: "This is a test Claim",
		Truth:       0.001,
	}
	d2 := gruff.Claim{
		Title:       "Another Claim",
		Description: "This a target claim",
		Truth:       1.000,
	}
	d3 := gruff.Claim{
		Title:       "The new arg's Claim",
		Description: "This a claim claim",
		Truth:       1.000,
	}
	TESTDB.Create(&d1)
	TESTDB.Create(&d2)
	TESTDB.Create(&d3)

	a1 := gruff.Argument{
		ClaimID:       d1.ID,
		TargetClaimID: &gruff.NullableUUID{UUID: d2.ID},
		Type:          gruff.ARGUMENT_AGAINST,
		Title:         "This is an argument",
		Description:   "This is an arguous description",
	}
	TESTDB.Create(&a1)

	a2 := gruff.Argument{
		ClaimID:          d3.ID,
		TargetArgumentID: &gruff.NullableUUID{UUID: a1.ID},
		Type:             gruff.ARGUMENT_FOR,
		Title:            "This is an argument to an argument",
		Description:      "This is an argumentous description",
	}

	r.POST("/api/arguments")
	r.SetBody(a2)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	TESTDB.Where("title = ?", a2.Title).First(&a2)
	expectedResults, _ := json.Marshal(a2)

	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, gruff.ARGUMENT_FOR, a2.Type)
	assert.Equal(t, d3.ID, a2.ClaimID)
	assert.Equal(t, a1.ID, a2.TargetArgumentID.UUID)
}

func TestCreateArgumentNoClaim(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	a1 := createArgument()

	r.POST("/api/arguments")
	r.SetBody(a1)
	res, _ := r.Run(Router())
	//assert.Equal(t, 400, res.Code)
	assert.Equal(t, 500, res.Code)
}

func TestCreateArgumentWithNewClaim(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	d1 := gruff.Claim{
		Title:       "Claim",
		Description: "This is a test Claim",
	}
	d2 := gruff.Claim{
		Title:       "Another Claim",
		Description: "This a target claim",
		Truth:       1.000,
	}
	TESTDB.Create(&d2)

	a1 := gruff.Argument{
		Claim:         &d1,
		TargetClaimID: &gruff.NullableUUID{UUID: d2.ID},
		Type:          gruff.ARGUMENT_AGAINST,
		Title:         "This is an argument",
		Description:   "This is an arguous description",
	}

	r.POST("/api/arguments")
	r.SetBody(a1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	TESTDB.Preload("Claim").Where("title = ?", a1.Title).First(&a1)
	expectedResults, _ := json.Marshal(a1)

	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, gruff.ARGUMENT_AGAINST, a1.Type)
	assert.Equal(t, d2.ID, a1.TargetClaimID.UUID)

	assert.Equal(t, d1.Title, a1.Claim.Title)
	assert.Equal(t, d1.Description, a1.Claim.Description)
}

func TestCreateArgumentWithNewClaimAndContexts(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	d1 := gruff.Claim{
		Title:       "Claim",
		Description: "This is a test Claim",
	}
	d2 := gruff.Claim{
		Title:       "Another Claim",
		Description: "This a target claim",
		Truth:       1.000,
	}
	TESTDB.Create(&d2)

	c1 := gruff.Context{Title: "Taylor Swift", URL: "http://en.wikipedia.com/Taylor_Swift"}
	c2 := gruff.Context{Title: "Donald Trump", URL: "http://en.wikipedia.com/Donald_Trump"}
	c3 := gruff.Context{Title: "Bozo the Clown", URL: "http://en.wikipedia.com/Bozo_the_Clown"}
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)
	TESTDB.Create(&c3)

	d1.ContextIDs = []uint64{c2.ID, c3.ID, c1.ID}

	a1 := gruff.Argument{
		Claim:         &d1,
		TargetClaimID: &gruff.NullableUUID{UUID: d2.ID},
		Type:          gruff.ARGUMENT_AGAINST,
		Title:         "This is an argument",
		Description:   "This is an arguous description",
	}

	startDBLog()
	defer stopDBLog()

	r.POST("/api/arguments")
	r.SetBody(a1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	TESTDB.Preload("Claim").Preload("Claim.Contexts").Where("title = ?", a1.Title).First(&a1)
	assert.Equal(t, gruff.ARGUMENT_AGAINST, a1.Type)
	assert.Equal(t, d2.ID, a1.TargetClaimID.UUID)

	assert.Equal(t, d1.Title, a1.Claim.Title)
	assert.Equal(t, d1.Description, a1.Claim.Description)
	assert.Equal(t, 3, len(a1.Claim.Contexts))

	a1.Claim.ContextIDs = []uint64{}
	a1.Claim.Contexts = []gruff.Context{}
	expectedResults, _ := json.Marshal(a1)
	assert.Equal(t, string(expectedResults), res.Body.String())

}

func TestCreateArgumentWithNewClaimNoTitle(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	d1 := gruff.Claim{
		Title:       "Claim",
		Description: "This is a test Claim",
	}
	d2 := gruff.Claim{
		Title:       "Another Claim",
		Description: "This a target claim",
		Truth:       1.000,
	}
	TESTDB.Create(&d2)

	a1 := gruff.Argument{
		Claim:         &d1,
		TargetClaimID: &gruff.NullableUUID{UUID: d2.ID},
		Type:          gruff.ARGUMENT_AGAINST,
	}

	r.POST("/api/arguments")
	r.SetBody(a1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	TESTDB.Preload("Claim").Where("title = ?", d1.Title).First(&a1)
	expectedResults, _ := json.Marshal(a1)

	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, gruff.ARGUMENT_AGAINST, a1.Type)
	assert.Equal(t, d2.ID, a1.TargetClaimID.UUID)

	assert.Equal(t, d1.Title, a1.Claim.Title)
	assert.Equal(t, d1.Description, a1.Claim.Description)
}

func TestCreateArgumentForArgumentWithNewClaim(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	d1 := gruff.Claim{
		Title:       "Claim",
		Description: "This is a test Claim",
		Truth:       0.001,
	}
	d2 := gruff.Claim{
		Title:       "Another Claim",
		Description: "This a target claim",
		Truth:       1.000,
	}
	d3 := gruff.Claim{
		Title:       "The new arg's Claim",
		Description: "This a claim claim",
		Truth:       1.000,
	}
	TESTDB.Create(&d1)
	TESTDB.Create(&d2)

	a1 := gruff.Argument{
		ClaimID:       d1.ID,
		TargetClaimID: &gruff.NullableUUID{UUID: d2.ID},
		Type:          gruff.ARGUMENT_AGAINST,
		Title:         "This is an argument",
		Description:   "This is an arguous description",
	}
	TESTDB.Create(&a1)

	a2 := gruff.Argument{
		Claim:            &d3,
		TargetArgumentID: &gruff.NullableUUID{UUID: a1.ID},
		Type:             gruff.ARGUMENT_FOR,
		Title:            "This is an argument to an argument",
		Description:      "This is an argumentous description",
	}

	r.POST("/api/arguments")
	r.SetBody(a2)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	TESTDB.Preload("Claim").Where("title = ?", a2.Title).First(&a2)
	expectedResults, _ := json.Marshal(a2)

	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, gruff.ARGUMENT_FOR, a2.Type)
	assert.Equal(t, a1.ID, a2.TargetArgumentID.UUID)

	assert.Equal(t, d3.Title, a2.Claim.Title)
	assert.Equal(t, d3.Description, a2.Claim.Description)
}

func TestCreateArgumentWithNewClaimNoArgData(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	d1 := gruff.Claim{
		Title:       "Claim",
		Description: "This is a test Claim",
	}
	d2 := gruff.Claim{
		Title:       "Another Claim",
		Description: "This a target claim",
		Truth:       1.000,
	}
	TESTDB.Create(&d2)

	a1 := gruff.Argument{
		Claim:         &d1,
		TargetClaimID: &gruff.NullableUUID{UUID: d2.ID},
		Type:          gruff.ARGUMENT_AGAINST,
	}

	r.POST("/api/arguments")
	r.SetBody(a1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	TESTDB.Preload("Claim").Where("title = ?", d1.Title).First(&a1)
	expectedResults, _ := json.Marshal(a1)

	assert.Equal(t, string(expectedResults), res.Body.String())
	assert.Equal(t, gruff.ARGUMENT_AGAINST, a1.Type)
	assert.Equal(t, d2.ID, a1.TargetClaimID.UUID)

	assert.Equal(t, d1.Title, a1.Claim.Title)
	assert.Equal(t, d1.Description, a1.Claim.Description)
	assert.Equal(t, d1.Title, a1.Title)
	assert.Equal(t, "", a1.Description)
}

func TestUpdateArgument(t *testing.T) {
	setup()
	defer teardown()

	r := New(Token)

	d1 := gruff.Claim{
		Title:       "Claim",
		Description: "This is a test Claim",
		Truth:       0.001,
	}
	d2 := gruff.Claim{
		Title:       "Another Claim",
		Description: "This a target claim",
		Truth:       1.000,
	}
	TESTDB.Create(&d1)
	TESTDB.Create(&d2)

	a1 := gruff.Argument{
		ClaimID:       d1.ID,
		TargetClaimID: &gruff.NullableUUID{UUID: d2.ID},
		Type:          gruff.ARGUMENT_AGAINST,
		Title:         "This is an argument",
		Description:   "This is an arguous description",
	}
	TESTDB.Create(&a1)

	r.PUT(fmt.Sprintf("/api/arguments/%s", a1.ID))
	r.SetBody(a1)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)
}

func TestDeleteArgument(t *testing.T) {
	setup()
	defer teardown()
	r := New(Token)

	a1 := createArgument()
	TESTDB.Create(&a1)

	r.DELETE(fmt.Sprintf("/api/arguments/%s", a1.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
}
*/
