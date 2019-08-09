package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/GruffDebate/server/gruff"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
)

func TestCreateClaim(t *testing.T) {
	setup()
	defer teardown()

	u := gruff.User{
		Name:     "claim user",
		Username: "APICreateClaimGuy",
		Email:    "claimguy@gruff.org",
		Password: "123456",
	}

	u.Create(CTX)

	claim := gruff.Claim{
		Title:        "API CreateClaim Let's create a new claim",
		Description:  "Claims in general should be true or false",
		Negation:     "Let's not...",
		Question:     "Should we create a new Claim?",
		Note:         "He who notes is a note taker",
		Image:        "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
		MultiPremise: true,
		PremiseRule:  gruff.PREMISE_RULE_ALL,
	}
	CTX.UserContext.Key = u.ArangoID()

	claim.VersionedModel.CreatedByID = u.ArangoID()

	r := New(tokenForTestUser(u))

	r.POST("/api/claims")
	r.SetBody(claim)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestGetClaim(t *testing.T) {
	setup()
	defer teardown()

	claim := gruff.Claim{
		Title:       "This is the API Get Claim test claim",
		Description: "Load all the things!",
		Negation:    "Don't load all the things.",
		Question:    "Load all the THINGS? Load ALL the things? LOAD all the things?",
		Note:        "This Claim needs to be all loaded.",
		Image:       "https://i.chzbgr.com/full/6434679808/h4ADBDEA5/",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	arg1 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "Get Claim All?",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)

	arg2 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "GET Claim Load ALL!",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)

	arg3 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "Do it Get Claim Load ALL!",
		Pro:           true,
	}
	err = arg3.Create(CTX)
	assert.NoError(t, err)

	context := gruff.Context{
		ShortName: "Get Claim context",
		Title:     "Get Claim context",
		URL:       "https://en.wikipedia.org/wiki/Claimed",
	}
	err = context.Create(CTX)
	assert.NoError(t, err)

	err = claim.AddContext(CTX, context)
	assert.NoError(t, err)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)

	expected, _ := json.Marshal(claim)

	r := New(tokenForTestUser(DEFAULT_USER))
	r.GET(fmt.Sprintf("/api/claims/%s", claim.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, string(expected), res.Body.String())
}

func TestConvertClaimToMultiPremise(t *testing.T) {
	setup()
	defer teardown()

	claim := gruff.Claim{
		Title:        "I'm just a premise, yes I'm only a premise. But I hope to be a multi-premise.",
		Description:  "ConvertToMultiPremise",
		Image:        "https://thesaurus.plus/img/synonyms/125/break_into_pieces.png",
		MultiPremise: false,
		PremiseRule:  gruff.PREMISE_RULE_NONE,
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "I don't belong on a multi-premise",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "What, you think I do?",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	context := gruff.Context{
		ShortName: "MPClaim Conversion",
		Title:     "Multi-Premise Claim Conversion",
		URL:       "https://en.wikipedia.org/wiki/Conversion",
	}
	err = context.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim.AddContext(CTX, context)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	anotherClaim := gruff.Claim{
		Title: "MPConversion other claim",
	}
	err = anotherClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argToAnotherClaim := gruff.Argument{
		TargetClaimID: &anotherClaim.ID,
		ClaimID:       claim.ID,
		Pro:           false,
	}
	err = argToAnotherClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	mpClaim := gruff.Claim{
		Title:        "I'm already an MPClaim. You wish you were!",
		MultiPremise: true,
		PremiseRule:  gruff.PREMISE_RULE_ALL,
	}
	err = mpClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = mpClaim.AddPremise(CTX, &claim)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	body := map[string]interface{}{
		"_key": claim.ArangoKey(),
	}

	r := New(tokenForTestUser(DEFAULT_USER))
	r.PUT(fmt.Sprintf("/api/claims/%s/convert", claim.ID))
	r.SetBody(body)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	expected, _ := json.Marshal(claim)
	assert.JSONEq(t, string(expected), res.Body.String())

	assert.True(t, claim.MultiPremise)
	assert.Equal(t, gruff.PREMISE_RULE_ALL, claim.PremiseRule)
}

func TestListTopClaims(t *testing.T) {
	setup()
	defer teardown()

	claim1 := gruff.Claim{
		Title: "I thought I was going to be clever, but this is Top Level Claim 1",
	}
	err := claim1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim2 := gruff.Claim{
		Title: "I thought I was going to be clever, but this is Top Level Claim 2",
	}
	err = claim2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim3 := gruff.Claim{
		Title: "I thought I was going to be clever, but this is Top Level Claim 3",
	}
	err = claim3.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim4 := gruff.Claim{
		Title: "I thought I was going to be clever, but this is Top Level Claim 4",
	}
	err = claim4.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim5 := gruff.Claim{
		Title: "I thought I was going to be clever, but this is Top Level Claim 5",
	}
	err = claim5.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg43 := gruff.Argument{
		TargetClaimID: &claim3.ID,
		ClaimID:       claim4.ID,
	}
	err = arg43.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claims := []gruff.Claim{claim5, claim3, claim2, claim1}
	expected, _ := json.Marshal(claims)

	r := New(tokenForTestUser(DEFAULT_USER))
	r.GET("/api/claims/top?start=0&limit=4")
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, string(expected), res.Body.String())

	claims = []gruff.Claim{claim2, claim1}
	expected, _ = json.Marshal(claims)

	r = New(tokenForTestUser(DEFAULT_USER))
	r.GET("/api/claims/top?start=2&limit=2")
	res, _ = r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, string(expected), res.Body.String())
}

func TestUpdateClaim(t *testing.T) {
	setup()
	defer teardown()

	claim := gruff.Claim{
		Title:       "This is the API Update Claim test claim",
		Description: "Update all the things!",
		Negation:    "Don't update all the things.",
		Question:    "Update all the THINGS? Update ALL the things? UPDATE all the things?",
		Note:        "This Claim needs to be all updated.",
		Image:       "https://i.chzbgr.com/full/6434679808/h4ADBDEA5/",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	arg1 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "Update Claim All?",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)

	arg2 := gruff.Argument{
		TargetClaimID: &claim.ID,
		Title:         "UPDATE Claim Update ALL!",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)

	body := map[string]interface{}{
		"_key":  claim.ArangoKey(),
		"title": "Just updating the title for this claim update",
	}

	r := New(tokenForTestUser(DEFAULT_USER))
	r.PUT(fmt.Sprintf("/api/claims/%s", claim.ID))
	r.SetBody(body)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	expected, _ := json.Marshal(claim)
	assert.JSONEq(t, string(expected), res.Body.String())

	assert.Equal(t, "Just updating the title for this claim update", claim.Title)
	assert.NotEqual(t, body["_key"], claim.ArangoKey())
}

func TestDeleteClaim(t *testing.T) {
	setup()
	defer teardown()

	claim := gruff.Claim{
		Title: "A simple claim to be deleted",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	r := New(tokenForTestUser(DEFAULT_USER))
	r.DELETE(fmt.Sprintf("/api/claims/%s", claim.ID))
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	claim.QueryAt = &claim.CreatedAt
	err = claim.Load(CTX)
	assert.NoError(t, err)
	expected, _ := json.Marshal(claim)
	assert.JSONEq(t, string(expected), res.Body.String())

	assert.NotNil(t, claim.DeletedAt)
}

func TestAddContextRemoveContext(t *testing.T) {
	setup()
	defer teardown()

	claim := gruff.Claim{
		Title: "getContext()",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	context := gruff.Context{
		ShortName: "Add Claim context",
		Title:     "Add Claim context",
		URL:       "https://golden.com/search/add%20context",
	}
	err = context.Create(CTX)
	assert.NoError(t, err)

	body := map[string]interface{}{
		"_key": claim.ArangoKey(),
	}

	r := New(tokenForTestUser(DEFAULT_USER))
	r.POST(fmt.Sprintf("/api/claims/%s/contexts/%s", claim.ID, context.Key))
	r.SetBody(body)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	expected, _ := json.Marshal(claim)
	assert.JSONEq(t, string(expected), res.Body.String())

	contexts, err := claim.Contexts(CTX)
	assert.Equal(t, 1, len(contexts))
	assert.Equal(t, context.ArangoKey(), contexts[0].ArangoKey())

	r = New(tokenForTestUser(DEFAULT_USER))
	r.DELETE(fmt.Sprintf("/api/claims/%s/contexts/%s", claim.ID, context.Key))
	res, _ = r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	expected, _ = json.Marshal(claim)
	assert.JSONEq(t, string(expected), res.Body.String())

	contexts, err = claim.Contexts(CTX)
	assert.Equal(t, 0, len(contexts))
}

func TestAddPremiseRemovePremise(t *testing.T) {
	setup()
	defer teardown()

	claim := gruff.Claim{
		Title:        "Add a Premise",
		MultiPremise: true,
		PremiseRule:  gruff.PREMISE_RULE_ALL,
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	premise1 := gruff.Claim{
		Title: "Already a premise",
	}
	err = premise1.Create(CTX)
	assert.NoError(t, err)

	err = claim.AddPremise(CTX, &premise1)
	assert.NoError(t, err)

	premise2 := gruff.Claim{
		Title: "Soon to be a premise",
	}
	err = premise2.Create(CTX)
	assert.NoError(t, err)

	body := map[string]interface{}{
		"_key": claim.ArangoKey(),
	}

	r := New(tokenForTestUser(DEFAULT_USER))
	r.POST(fmt.Sprintf("/api/claims/%s/premises/%s", claim.ID, premise2.ID))
	r.SetBody(body)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	expected, _ := json.Marshal(claim)
	assert.JSONEq(t, string(expected), res.Body.String())

	premises, err := claim.Premises(CTX)
	assert.Equal(t, 2, len(premises))
	assert.Equal(t, premise2.ArangoID(), premises[1].ArangoID())

	r = New(tokenForTestUser(DEFAULT_USER))
	r.DELETE(fmt.Sprintf("/api/claims/%s/premises/%s", claim.ID, premise1.ID))
	res, _ = r.Run(Router())
	assert.Equal(t, http.StatusOK, res.Code)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	expected, _ = json.Marshal(claim)
	assert.JSONEq(t, string(expected), res.Body.String())

	premises, err = claim.Premises(CTX)
	assert.Equal(t, 1, len(premises))
	assert.Equal(t, premise2.ArangoID(), premises[0].ArangoID())
}

/*
func TestSetTruthScore(t *testing.T) {
	setup()
	defer teardown()
	u := createTestUser()
	r := New(tokenForTestUser(u))

	c1 := createClaim()
	c2 := createClaim()
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)

	m := map[string]interface{}{
		"score": 0.2394,
	}

	r.POST(fmt.Sprintf("/api/claims/%s/truth", c1.ID))
	r.SetBody(m)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	co := gruff.ClaimOpinion{}
	err := TESTDB.Where("user_id = ?", u.ID).Where("claim_id = ?", c1.ID).First(&co).Error
	assert.Nil(t, err)
	assert.Equal(t, 0.2394, co.Truth)

	TESTDB.First(&c1)
	assert.Equal(t, 0.2394, c1.Truth)
}

func TestSetTruthScoreUpdate(t *testing.T) {
	setup()
	defer teardown()
	u := createTestUser()
	r := New(tokenForTestUser(u))

	c1 := createClaim()
	c2 := createClaim()
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)

	co := gruff.ClaimOpinion{
		UserID:  u.ID,
		ClaimID: c1.ID,
		Truth:   0.8239,
	}
	TESTDB.Create(&co)

	c1.UpdateTruth(CTX)
	TESTDB.First(&c1)
	assert.Equal(t, 0.8239, c1.Truth)

	m := map[string]interface{}{
		"score": 0.2394,
	}

	r.POST(fmt.Sprintf("/api/claims/%s/truth", c1.ID))
	r.SetBody(m)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)

	co = gruff.ClaimOpinion{}
	err := TESTDB.Where("user_id = ?", u.ID).Where("claim_id = ?", c1.ID).First(&co).Error
	assert.Nil(t, err)
	assert.Equal(t, 0.2394, co.Truth)

	TESTDB.First(&c1)
	assert.Equal(t, 0.2394, c1.Truth)
}

func TestSetScoreStrength(t *testing.T) {
	setup()
	defer teardown()
	u := createTestUser()
	r := New(tokenForTestUser(u))

	c1 := createClaim()
	c2 := createClaim()
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)

	a := gruff.Argument{TargetClaimID: gruff.NUUID(c1.ID), ClaimID: c2.ID, Title: "An arg"}
	TESTDB.Create(&a)

	m := map[string]interface{}{
		"score": 0.2394,
	}

	r.POST(fmt.Sprintf("/api/arguments/%s/strength", a.ID))
	r.SetBody(m)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusCreated, res.Code)

	co := gruff.ArgumentOpinion{}
	err := TESTDB.Where("user_id = ?", u.ID).Where("argument_id = ?", a.ID).First(&co).Error
	assert.Nil(t, err)
	assert.Equal(t, 0.2394, co.Strength)

	TESTDB.First(&a)
	assert.Equal(t, 0.2394, a.Strength)
}

func TestSetScoreStrengthUpdate(t *testing.T) {
	setup()
	defer teardown()
	u := createTestUser()
	r := New(tokenForTestUser(u))

	c1 := createClaim()
	c2 := createClaim()
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)

	a := gruff.Argument{TargetClaimID: gruff.NUUID(c1.ID), ClaimID: c2.ID, Title: "An arg"}
	TESTDB.Create(&a)

	co := gruff.ArgumentOpinion{
		UserID:     u.ID,
		ArgumentID: a.ID,
		Strength:   0.8239,
	}
	TESTDB.Create(&co)

	a.UpdateStrength(CTX)
	// a.UpdateRelevance(CTX.ServerContext())

	m := map[string]interface{}{
		"score": 0.2394,
	}

	r.POST(fmt.Sprintf("/api/arguments/%s/strength", a.ID))
	r.SetBody(m)
	res, _ := r.Run(Router())
	assert.Equal(t, http.StatusAccepted, res.Code)

	co = gruff.ArgumentOpinion{}
	err := TESTDB.Where("user_id = ?", u.ID).Where("argument_id = ?", a.ID).First(&co).Error
	assert.Nil(t, err)
	assert.Equal(t, 0.2394, co.Strength)

	TESTDB.First(&a)
	assert.Equal(t, 0.2394, a.Strength)
}

*/
