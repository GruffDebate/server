package gruff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateClaim(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "Let's create a new claim",
		Description:  "Claims in general should be true or false",
		Negation:     "Let's not...",
		Question:     "Should we create a new Claim?",
		Note:         "He who notes is a note taker",
		Image:        "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}

	saved, err := claim.Load(CTX)
	assert.Error(t, err)
	assert.Empty(t, saved.Key)

	err = claim.Create(CTX)
	assert.NoError(t, err)
	saved, err = claim.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Nil(t, saved.DeletedAt)
	assert.Equal(t, claim.Title, saved.Title)
	assert.Equal(t, claim.Description, saved.Description)
	assert.Equal(t, claim.Negation, saved.Negation)
	assert.Equal(t, claim.Question, saved.Question)
	assert.Equal(t, claim.Note, saved.Note)
	assert.Equal(t, claim.Image, saved.Image)
	assert.True(t, saved.MultiPremise)
	assert.Equal(t, PREMISE_RULE_ALL, saved.PremiseRule)
}

func TestClaimAddPremise(t *testing.T) {
	setupDB()
	defer teardownDB()

	topClaim := Claim{
		Title:        "I dare you to doubt me",
		Description:  "I am true. Woe be the person that doubts my veracity",
		Negation:     "I dare you to accept me",
		Question:     "Do you dare to doubt me?",
		Note:         "This Claim is all about doubting. No links are going here.",
		Image:        "https://upload.wikimedia.org/wikipedia/en/thumb/7/7d/NoDoubtCover.png/220px-NoDoubtCover.png",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}

	premiseClaim1 := Claim{
		Title:        "I am the one who is daring you to doubt mean",
		Description:  "The person that is daring you to doubt me being me",
		MultiPremise: false,
	}

	premiseClaim2 := Claim{
		Title:        "Since it is I that am daring you, you therefore must not doubt",
		Description:  "I am undoubtable",
		MultiPremise: false,
	}

	err := topClaim.Create(CTX)
	assert.NoError(t, err)
	saved, err := topClaim.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Nil(t, saved.DeletedAt)

	err = topClaim.AddPremise(CTX, &premiseClaim1)
	assert.NoError(t, err)
	saved, err = premiseClaim1.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Nil(t, saved.DeletedAt)

	// Assert that the proper link has been created
	premiseEdges, err := topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premiseEdges))
	assert.Equal(t, topClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, premiseClaim1.ArangoID(), premiseEdges[0].To)
	assert.Equal(t, 1, premiseEdges[0].Order)

	n, err := topClaim.NumberOfPremises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), n)

	err = topClaim.AddPremise(CTX, &premiseClaim2)
	assert.NoError(t, err)
	saved, err = premiseClaim2.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Nil(t, saved.DeletedAt)

	// Assert that the proper link has been created
	premiseEdges, err = topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(premiseEdges))
	assert.Equal(t, topClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, premiseClaim1.ArangoID(), premiseEdges[0].To)
	assert.Equal(t, 1, premiseEdges[0].Order)
	assert.Equal(t, topClaim.ArangoID(), premiseEdges[1].From)
	assert.Equal(t, premiseClaim2.ArangoID(), premiseEdges[1].To)
	assert.Equal(t, 2, premiseEdges[1].Order)

	n, err = topClaim.NumberOfPremises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)

	n, err = premiseClaim1.NumberOfPremises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)

	n, err = premiseClaim2.NumberOfPremises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)

}

func TestClaimArangoID(t *testing.T) {
	claim := Claim{}
	claim.Key = "somethingpredictable"
	assert.Equal(t, "claims/somethingpredictable", claim.ArangoID())
}

func TestClaimLoadByID(t *testing.T) {
	setupDB()
	defer teardownDB()

}

// TODO: Test Adding an Argument
// TODO: Test getting an Argument
// TODO: Test Inferences
// TODO: Test BaseClaimEdges
