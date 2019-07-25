package gruff

import (
	"testing"
	"time"

	"github.com/GruffDebate/server/support"
	"github.com/stretchr/testify/assert"
)

func TestCreateClaim(t *testing.T) {
	setupDB()
	defer teardownDB()

	u := User{}
	u.Key = "testuser"
	CTX.UserContext = u

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

	saved := Claim{}
	saved.ID = claim.ID
	err := saved.Load(CTX)
	assert.Error(t, err)
	assert.Empty(t, saved.Key)

	err = claim.Create(CTX)
	assert.NoError(t, err)
	saved = Claim{}
	saved.ID = claim.ID
	err = saved.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Equal(t, u.ArangoID(), saved.CreatedByID)
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
	saved := Claim{}
	saved.ID = topClaim.ID
	err = saved.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Nil(t, saved.DeletedAt)

	err = topClaim.AddPremise(CTX, &premiseClaim1)
	assert.NoError(t, err)
	saved = Claim{}
	saved.ID = premiseClaim1.ID
	err = saved.Load(CTX)
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

	premises, err := topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premises))
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())

	err = topClaim.AddPremise(CTX, &premiseClaim2)
	assert.NoError(t, err)
	saved.ID = premiseClaim2.ID
	err = saved.Load(CTX)
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

	premises, err = topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(premises))
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[1].ArangoID())
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

func TestClaimVersion(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "I dare you to doubt me",
		Description:  "I am true. Woe be the person that doubts my veracity",
		Negation:     "I dare you to accept me",
		Question:     "Do you dare to doubt me?",
		Note:         "This Claim is all about doubting. No links are going here.",
		Image:        "https://upload.wikimedia.org/wikipedia/en/thumb/7/7d/NoDoubtCover.png/220px-NoDoubtCover.png",
		MultiPremise: false,
		PremiseRule:  PREMISE_RULE_NONE,
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premiseClaim1 := Claim{
		Title:        "I am the one who is daring you to doubt mean",
		Description:  "The person that is daring you to doubt me being me",
		MultiPremise: false,
	}
	err = claim.AddPremise(CTX, &premiseClaim1)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premiseClaim2 := Claim{
		Title:        "Since it is I that am daring you, you therefore must not doubt",
		Description:  "I am undoubtable",
		MultiPremise: false,
	}
	err = claim.AddPremise(CTX, &premiseClaim2)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	distantClaim := Claim{
		Title:       "So very far away",
		Description: "So distant, you cannot see me.",
	}
	err = distantClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	mpClaim := Claim{
		Title:        "This is an MP that uses the main claim",
		Description:  "Not military police, mind you",
		Image:        "https://sgws3productimages.azureedge.net/sgwproductimages/images/33/4-5-2019/35ae88a198be52fc1.JPG",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}
	err = mpClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = mpClaim.AddPremise(CTX, &claim)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Let's create a new argument",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "I beg to differ",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argPC1 := Argument{
		TargetClaimID: &premiseClaim1.ID,
		Title:         "Let's create a new argument",
	}
	err = argPC1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argDC1 := Argument{
		TargetClaimID: &distantClaim.ID,
		ClaimID:       claim.ID,
		Title:         "I want to get away",
	}
	err = argDC1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1arg := Argument{
		TargetArgumentID: &arg1.ID,
		Title:            "Let's create a new argument argument",
		Pro:              false,
	}
	err = arg1arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Next check edges, then version and recheck everything
	premiseEdges, err := claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(premiseEdges))
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[1].From)
	assert.Equal(t, premiseClaim1.ArangoID(), premiseEdges[0].To)
	assert.Equal(t, premiseClaim2.ArangoID(), premiseEdges[1].To)

	inferences, err := claim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(inferences))
	assert.Equal(t, claim.ArangoID(), inferences[0].From)
	assert.Equal(t, claim.ArangoID(), inferences[1].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.Equal(t, arg2.ArangoID(), inferences[1].To)

	bces, err := claim.BaseClaimEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bces))
	assert.Equal(t, argDC1.ArangoID(), bces[0].From)
	assert.Equal(t, claim.ArangoID(), bces[0].To)

	args, err := claim.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())
	assert.Equal(t, arg2.ArangoID(), args[1].ArangoID())

	premiseEdges, err = mpClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premiseEdges))
	assert.Equal(t, mpClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].To)

	// Version the claim
	err = claim.Load(CTX)
	assert.NoError(t, err)
	origClaimKey := claim.ArangoKey()

	claim.Title = "New Title"
	err = claim.version(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "New Title", claim.Title)
	assert.NotEqual(t, origClaimKey, claim.ArangoKey())

	origClaim := Claim{}
	origClaim.Key = origClaimKey
	err = origClaim.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, origClaim.DeletedAt)
	assert.Equal(t, "I dare you to doubt me", origClaim.Title)
	assert.Equal(t, origClaimKey, origClaim.ArangoKey())

	// Verify new edges were created
	premiseEdges, err = claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(premiseEdges))
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[1].From)
	assert.Equal(t, premiseClaim1.ArangoID(), premiseEdges[0].To)
	assert.Equal(t, premiseClaim2.ArangoID(), premiseEdges[1].To)
	assert.Nil(t, premiseEdges[0].DeletedAt)
	assert.Nil(t, premiseEdges[1].DeletedAt)

	inferences, err = claim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(inferences))
	assert.Equal(t, claim.ArangoID(), inferences[0].From)
	assert.Equal(t, claim.ArangoID(), inferences[1].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.Equal(t, arg2.ArangoID(), inferences[1].To)
	assert.Nil(t, inferences[0].DeletedAt)
	assert.Nil(t, inferences[1].DeletedAt)

	bces, err = claim.BaseClaimEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bces))
	assert.Equal(t, argDC1.ArangoID(), bces[0].From)
	assert.Equal(t, claim.ArangoID(), bces[0].To)
	assert.Nil(t, bces[0].DeletedAt)

	args, err = claim.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())
	assert.Equal(t, arg2.ArangoID(), args[1].ArangoID())

	premiseEdges, err = mpClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premiseEdges))
	assert.Equal(t, mpClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].To)
	assert.Nil(t, premiseEdges[0].DeletedAt)

	// Verify that the old edges were deleted
	premiseEdges, err = origClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(premiseEdges))
	assert.Equal(t, origClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, origClaim.ArangoID(), premiseEdges[1].From)
	assert.Equal(t, premiseClaim1.ArangoID(), premiseEdges[0].To)
	assert.Equal(t, premiseClaim2.ArangoID(), premiseEdges[1].To)
	assert.NotNil(t, premiseEdges[0].DeletedAt)
	assert.NotNil(t, premiseEdges[1].DeletedAt)

	inferences, err = origClaim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(inferences))
	assert.Equal(t, origClaim.ArangoID(), inferences[0].From)
	assert.Equal(t, origClaim.ArangoID(), inferences[1].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.Equal(t, arg2.ArangoID(), inferences[1].To)
	assert.NotNil(t, inferences[0].DeletedAt)
	assert.NotNil(t, inferences[1].DeletedAt)

	bces, err = origClaim.BaseClaimEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bces))
	assert.Equal(t, argDC1.ArangoID(), bces[0].From)
	assert.Equal(t, origClaim.ArangoID(), bces[0].To)
	assert.NotNil(t, bces[0].DeletedAt)

	olderMpClaim := Claim{}
	olderMpClaim.Key = mpClaim.Key
	olderMpClaim.DeletedAt = origClaim.DeletedAt
	err = olderMpClaim.Load(CTX)
	assert.NoError(t, err)
	premiseEdges, err = olderMpClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premiseEdges))
	assert.Equal(t, olderMpClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, origClaim.ArangoID(), premiseEdges[0].To)
	assert.NotNil(t, premiseEdges[0].DeletedAt)
}

func TestLoadClaimAtDate(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:       "I dare you to doubt me",
		Description: "I am true. Woe be the person that doubts my veracity",
		Negation:    "I dare you to accept me",
		Question:    "Do you dare to doubt me?",
		Note:        "This Claim is all about doubting. No links are going here.",
		Image:       "https://upload.wikimedia.org/wikipedia/en/thumb/7/7d/NoDoubtCover.png/220px-NoDoubtCover.png",
	}
	claim.DeletedAt = support.TimePtr(time.Now().Add(-24 * time.Hour))

	err := claim.Create(CTX)
	assert.NoError(t, err)
	patch := map[string]interface{}{"start": time.Now().Add(-25 * time.Hour)}
	col, _ := CTX.Arango.CollectionFor(claim)
	col.UpdateDocument(CTX.Context, claim.ArangoKey(), patch)

	firstKey := claim.ArangoKey()

	claim.DeletedAt = support.TimePtr(time.Now().Add(-1 * time.Hour))
	err = claim.Create(CTX)
	assert.NoError(t, err)
	patch["start"] = time.Now().Add(-24 * time.Hour)
	col.UpdateDocument(CTX.Context, claim.ArangoKey(), patch)

	secondKey := claim.ArangoKey()

	claim.DeletedAt = nil
	err = claim.Create(CTX)
	assert.NoError(t, err)
	patch["start"] = time.Now().Add(-1 * time.Hour)
	col.UpdateDocument(CTX.Context, claim.ArangoKey(), patch)

	thirdKey := claim.ArangoKey()

	lookup := Claim{}
	lookup.ID = claim.ID
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())

	lookup = Claim{}
	lookup.ID = claim.ID
	lookup.CreatedAt = time.Now().Add(-1 * time.Minute)
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())
	thirdCreatedAt := lookup.CreatedAt

	lookup = Claim{}
	lookup.ID = claim.ID
	lookup.QueryAt = support.TimePtr(time.Now().Add(-2 * time.Hour))
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, secondKey, lookup.ArangoKey())
	secondCreatedAt := lookup.CreatedAt

	lookup = Claim{}
	lookup.ID = claim.ID
	lookup.QueryAt = support.TimePtr(time.Now().Add(-25 * time.Hour))
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, firstKey, lookup.ArangoKey())
	firstCreatedAt := lookup.CreatedAt

	// TODO: Throw a NotFoundError?
	lookup = Claim{}
	lookup.ID = claim.ID
	lookup.QueryAt = support.TimePtr(time.Now().Add(-48 * time.Hour))
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Equal(t, "", lookup.ArangoKey())

	lookup = Claim{}
	lookup.ID = claim.ID
	lookup.QueryAt = &firstCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, firstKey, lookup.ArangoKey())

	lookup = Claim{}
	lookup.ID = claim.ID
	lookup.QueryAt = &secondCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, secondKey, lookup.ArangoKey())

	lookup = Claim{}
	lookup.ID = claim.ID
	lookup.QueryAt = &thirdCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())
}

func TestClaimReorderPremise(t *testing.T) {
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

	premiseClaim3 := Claim{
		Title: "I don't even get a description",
	}

	premiseClaim4 := Claim{
		Title: "Talk to the hand",
	}

	err := topClaim.Create(CTX)
	assert.NoError(t, err)

	err = topClaim.AddPremise(CTX, &premiseClaim1)
	assert.NoError(t, err)

	err = topClaim.AddPremise(CTX, &premiseClaim2)
	assert.NoError(t, err)

	err = topClaim.AddPremise(CTX, &premiseClaim3)
	assert.NoError(t, err)

	err = topClaim.AddPremise(CTX, &premiseClaim4)
	assert.NoError(t, err)

	premises, err := topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(premises))
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[3].ArangoID())

	edges, err := topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(edges))
	assert.Equal(t, premiseClaim1.ArangoID(), edges[0].To)
	assert.Equal(t, premiseClaim2.ArangoID(), edges[1].To)
	assert.Equal(t, premiseClaim3.ArangoID(), edges[2].To)
	assert.Equal(t, premiseClaim4.ArangoID(), edges[3].To)
	assert.Equal(t, 1, edges[0].Order)
	assert.Equal(t, 2, edges[1].Order)
	assert.Equal(t, 3, edges[2].Order)
	assert.Equal(t, 4, edges[3].Order)

	premises, err = topClaim.ReorderPremise(CTX, premiseClaim1, 0)
	assert.Error(t, err)
	assert.Equal(t, "Order: invalid value;", err.Error())

	premises, err = topClaim.ReorderPremise(CTX, premiseClaim1, 5)
	assert.Error(t, err)
	assert.Equal(t, "Order: the new order is higher than the number of premises;", err.Error())

	// Move the first to the last
	premises, err = topClaim.ReorderPremise(CTX, premiseClaim1, 4)
	assert.NoError(t, err)
	assert.Equal(t, premiseClaim2.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim1.ArangoID(), premises[3].ArangoID())

	premises, err = topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(premises))
	assert.Equal(t, premiseClaim2.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim1.ArangoID(), premises[3].ArangoID())

	edges, err = topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(edges))
	assert.Equal(t, premiseClaim2.ArangoID(), edges[0].To)
	assert.Equal(t, premiseClaim3.ArangoID(), edges[1].To)
	assert.Equal(t, premiseClaim4.ArangoID(), edges[2].To)
	assert.Equal(t, premiseClaim1.ArangoID(), edges[3].To)
	assert.Equal(t, 1, edges[0].Order)
	assert.Equal(t, 2, edges[1].Order)
	assert.Equal(t, 3, edges[2].Order)
	assert.Equal(t, 4, edges[3].Order)

	// Move the last to the first
	premises, err = topClaim.ReorderPremise(CTX, premiseClaim1, 1)
	assert.NoError(t, err)
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[3].ArangoID())

	premises, err = topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(premises))
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[3].ArangoID())

	edges, err = topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(edges))
	assert.Equal(t, premiseClaim1.ArangoID(), edges[0].To)
	assert.Equal(t, premiseClaim2.ArangoID(), edges[1].To)
	assert.Equal(t, premiseClaim3.ArangoID(), edges[2].To)
	assert.Equal(t, premiseClaim4.ArangoID(), edges[3].To)
	assert.Equal(t, 1, edges[0].Order)
	assert.Equal(t, 2, edges[1].Order)
	assert.Equal(t, 3, edges[2].Order)
	assert.Equal(t, 4, edges[3].Order)

	// Move back
	premises, err = topClaim.ReorderPremise(CTX, premiseClaim2, 3)
	assert.NoError(t, err)
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[3].ArangoID())

	premises, err = topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(premises))
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[3].ArangoID())

	edges, err = topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(edges))
	assert.Equal(t, premiseClaim1.ArangoID(), edges[0].To)
	assert.Equal(t, premiseClaim3.ArangoID(), edges[1].To)
	assert.Equal(t, premiseClaim2.ArangoID(), edges[2].To)
	assert.Equal(t, premiseClaim4.ArangoID(), edges[3].To)
	assert.Equal(t, 1, edges[0].Order)
	assert.Equal(t, 2, edges[1].Order)
	assert.Equal(t, 3, edges[2].Order)
	assert.Equal(t, 4, edges[3].Order)

	// Move forward
	premises, err = topClaim.ReorderPremise(CTX, premiseClaim2, 2)
	assert.NoError(t, err)
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[3].ArangoID())

	premises, err = topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(premises))
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[3].ArangoID())

	edges, err = topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(edges))
	assert.Equal(t, premiseClaim1.ArangoID(), edges[0].To)
	assert.Equal(t, premiseClaim2.ArangoID(), edges[1].To)
	assert.Equal(t, premiseClaim3.ArangoID(), edges[2].To)
	assert.Equal(t, premiseClaim4.ArangoID(), edges[3].To)
	assert.Equal(t, 1, edges[0].Order)
	assert.Equal(t, 2, edges[1].Order)
	assert.Equal(t, 3, edges[2].Order)
	assert.Equal(t, 4, edges[3].Order)

	// Move second to last
	premises, err = topClaim.ReorderPremise(CTX, premiseClaim2, 4)
	assert.NoError(t, err)
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[3].ArangoID())

	premises, err = topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(premises))
	assert.Equal(t, premiseClaim1.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim4.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[3].ArangoID())

	edges, err = topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(edges))
	assert.Equal(t, premiseClaim1.ArangoID(), edges[0].To)
	assert.Equal(t, premiseClaim3.ArangoID(), edges[1].To)
	assert.Equal(t, premiseClaim4.ArangoID(), edges[2].To)
	assert.Equal(t, premiseClaim2.ArangoID(), edges[3].To)
	assert.Equal(t, 1, edges[0].Order)
	assert.Equal(t, 2, edges[1].Order)
	assert.Equal(t, 3, edges[2].Order)
	assert.Equal(t, 4, edges[3].Order)

	// Move third to first
	premises, err = topClaim.ReorderPremise(CTX, premiseClaim4, 1)
	assert.NoError(t, err)
	assert.Equal(t, premiseClaim4.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim1.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[3].ArangoID())

	premises, err = topClaim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(premises))
	assert.Equal(t, premiseClaim4.ArangoID(), premises[0].ArangoID())
	assert.Equal(t, premiseClaim1.ArangoID(), premises[1].ArangoID())
	assert.Equal(t, premiseClaim3.ArangoID(), premises[2].ArangoID())
	assert.Equal(t, premiseClaim2.ArangoID(), premises[3].ArangoID())

	edges, err = topClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(edges))
	assert.Equal(t, premiseClaim4.ArangoID(), edges[0].To)
	assert.Equal(t, premiseClaim1.ArangoID(), edges[1].To)
	assert.Equal(t, premiseClaim3.ArangoID(), edges[2].To)
	assert.Equal(t, premiseClaim2.ArangoID(), edges[3].To)
	assert.Equal(t, 1, edges[0].Order)
	assert.Equal(t, 2, edges[1].Order)
	assert.Equal(t, 3, edges[2].Order)
	assert.Equal(t, 4, edges[3].Order)
}

func TestClaimQueryForTopLevelClaims(t *testing.T) {
	params := ArangoQueryParameters{}
	assert.Equal(t, "FOR obj IN claims LET bcCount=(FOR bc IN base_claims FILTER bc._to == obj._id COLLECT WITH COUNT INTO length RETURN length) FILTER bcCount[0] == 0 SORT obj.start DESC LIMIT 0, 20 RETURN obj", Claim{}.QueryForTopLevelClaims(params))

	params.Return = support.StringPtr("obj._id")
	assert.Equal(t, "FOR obj IN claims LET bcCount=(FOR bc IN base_claims FILTER bc._to == obj._id COLLECT WITH COUNT INTO length RETURN length) FILTER bcCount[0] == 0 SORT obj.start DESC LIMIT 0, 20 RETURN obj._id", Claim{}.QueryForTopLevelClaims(params))

	params.Return = support.StringPtr("{claim: obj, count: bcCount[0]}")
	assert.Equal(t, "FOR obj IN claims LET bcCount=(FOR bc IN base_claims FILTER bc._to == obj._id COLLECT WITH COUNT INTO length RETURN length) FILTER bcCount[0] == 0 SORT obj.start DESC LIMIT 0, 20 RETURN {claim: obj, count: bcCount[0]}", Claim{}.QueryForTopLevelClaims(params))

	params.Return = nil
	params.Sort = support.StringPtr("obj._id")
	assert.Equal(t, "FOR obj IN claims LET bcCount=(FOR bc IN base_claims FILTER bc._to == obj._id COLLECT WITH COUNT INTO length RETURN length) FILTER bcCount[0] == 0 SORT obj._id LIMIT 0, 20 RETURN obj", Claim{}.QueryForTopLevelClaims(params))

	params.Offset = support.IntPtr(5)
	assert.Equal(t, "FOR obj IN claims LET bcCount=(FOR bc IN base_claims FILTER bc._to == obj._id COLLECT WITH COUNT INTO length RETURN length) FILTER bcCount[0] == 0 SORT obj._id LIMIT 5, 20 RETURN obj", Claim{}.QueryForTopLevelClaims(params))

	params.Limit = support.IntPtr(10)
	assert.Equal(t, "FOR obj IN claims LET bcCount=(FOR bc IN base_claims FILTER bc._to == obj._id COLLECT WITH COUNT INTO length RETURN length) FILTER bcCount[0] == 0 SORT obj._id LIMIT 5, 10 RETURN obj", Claim{}.QueryForTopLevelClaims(params))

	params.Offset = nil
	assert.Equal(t, "FOR obj IN claims LET bcCount=(FOR bc IN base_claims FILTER bc._to == obj._id COLLECT WITH COUNT INTO length RETURN length) FILTER bcCount[0] == 0 SORT obj._id LIMIT 0, 10 RETURN obj", Claim{}.QueryForTopLevelClaims(params))
}

func TestClaimLoadFull(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "This is the Claim LoadAll test claim",
		Description:  "Load all the things!",
		Negation:     "Don't load all the things.",
		Question:     "Load all the THINGS? Load ALL the things? LOAD all the things?",
		Note:         "This Claim needs to be all loaded.",
		Image:        "https://i.chzbgr.com/full/6434679808/h4ADBDEA5/",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premiseClaim1 := Claim{
		Title:        "First premise for the Claim LoadAll dude",
		Description:  "The person that is daring you to doubt me being me",
		MultiPremise: false,
	}
	err = claim.AddPremise(CTX, &premiseClaim1)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premiseClaim2 := Claim{
		Title:        "I am the second Claim LoadAll premise. I MUST be true.",
		Description:  "I am undoubtable",
		MultiPremise: false,
	}
	err = claim.AddPremise(CTX, &premiseClaim2)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	distantClaim := Claim{
		Title:       "So very far away from Claim LoadAll",
		Description: "So distant, you cannot see me.",
	}
	err = distantClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	mpClaim := Claim{
		Title:        "This is an MP that uses the main LoadAll claim",
		Description:  "Not military police, mind you",
		Image:        "https://sgws3productimages.azureedge.net/sgwproductimages/images/33/4-5-2019/35ae88a198be52fc1.JPG",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}
	err = mpClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = mpClaim.AddPremise(CTX, &claim)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "All?",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Load ALL!",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg3 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Do it ALL!",
		Pro:           true,
	}
	err = arg3.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argPC1 := Argument{
		TargetClaimID: &premiseClaim1.ID,
		Title:         "Let's create a new argument for the Premise of the claim LoadAll",
	}
	err = argPC1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argDC1 := Argument{
		TargetClaimID: &distantClaim.ID,
		ClaimID:       claim.ID,
		Title:         "Distant LoadAll claim.",
	}
	err = argDC1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1arg := Argument{
		TargetArgumentID: &arg1.ID,
		Title:            "Do it...ALLLLL!",
		Pro:              false,
	}
	err = arg1arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Simple Load
	err = claim.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))

	// Load All
	premiseClaim1.Load(CTX)
	premiseClaim2.Load(CTX)
	arg1.Load(CTX)
	arg2.Load(CTX)
	arg3.Load(CTX)
	var carg1, carg2, carg3 Claim
	carg1.ID = arg1.ClaimID
	carg2.ID = arg2.ClaimID
	carg3.ID = arg3.ClaimID
	carg1.Load(CTX)
	carg2.Load(CTX)
	carg3.Load(CTX)
	arg1.Claim = &carg1
	arg2.Claim = &carg2
	arg3.Claim = &carg3

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 2, len(claim.PremiseClaims))
	assert.Equal(t, 2, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, premiseClaim1, claim.PremiseClaims[0])
	assert.Equal(t, premiseClaim2, claim.PremiseClaims[1])
	assert.Equal(t, arg1, claim.ProArgs[0])
	assert.Equal(t, arg3, claim.ProArgs[1])
	assert.Equal(t, arg2, claim.ConArgs[0])

	err = arg1.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	err = premiseClaim1.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	arg1.Load(CTX)
	premiseClaim1.Load(CTX)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 1, len(claim.PremiseClaims))
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, premiseClaim2, claim.PremiseClaims[0])
	assert.Equal(t, arg3, claim.ProArgs[0])
	assert.Equal(t, arg2, claim.ConArgs[0])

	// Previous points in time
	claim.QueryAt = &arg1arg.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 2, len(claim.PremiseClaims))
	assert.Equal(t, 2, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, premiseClaim1, claim.PremiseClaims[0])
	assert.Equal(t, premiseClaim2, claim.PremiseClaims[1])
	assert.Equal(t, arg1, claim.ProArgs[0])
	assert.Equal(t, arg3, claim.ProArgs[1])
	assert.Equal(t, arg2, claim.ConArgs[0])

	claim.QueryAt = &arg1.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 2, len(claim.PremiseClaims))
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))
	assert.Equal(t, premiseClaim1, claim.PremiseClaims[0])
	assert.Equal(t, premiseClaim2, claim.PremiseClaims[1])
	assert.Equal(t, arg1, claim.ProArgs[0])

	claim.QueryAt = &claim.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))

	premiseClaim3 := Claim{
		Title:        "You can never have enough premises in a LoadAll test",
		Description:  "At least, that's my claim...",
		MultiPremise: false,
	}
	err = claim.AddPremise(CTX, &premiseClaim3)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	premiseClaim3.Load(CTX)

	claim.QueryAt = nil
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 2, len(claim.PremiseClaims))
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, premiseClaim2, claim.PremiseClaims[0])
	assert.Equal(t, premiseClaim3, claim.PremiseClaims[1])
	assert.Equal(t, arg3, claim.ProArgs[0])
	assert.Equal(t, arg2, claim.ConArgs[0])
}
