package gruff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateClaim(t *testing.T) {
	setupDB()
	defer teardownDB()

	assert.Equal(t, "A", "A")
}

func TestCreateFullGraph(t *testing.T) {
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

	saved, err := topClaim.Load(CTX)
	assert.Error(t, err, "some kind of error")
	assert.Empty(t, saved.Key)

	err = topClaim.Create(CTX)
	assert.NoError(t, err)
	saved, err = topClaim.Load(CTX)
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

	// TODO: assert versioning!
}
