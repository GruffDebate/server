package gruff

import (
	"fmt"
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

	err = claim.Create(CTX)
	assert.Error(t, err)
	assert.Equal(t, "A claim with the same ID already exists", err.Error())

	claim = Claim{}
	err = claim.Create(CTX)
	assert.Error(t, err)
	assert.Equal(t, "Title: must be between 3 and 1000 characters;", err.Error())

	claim.Title = "Something more than 3 characters"
	claim.Description = "AB"
	err = claim.Create(CTX)
	assert.Error(t, err)
	assert.Equal(t, "Description: must be blank, or between 3 and 4000 characters;", err.Error())

	claim.Description = ""
	err = claim.Create(CTX)
	assert.NoError(t, err)
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

	topClaim.DeletedAt = support.TimePtr(time.Now())
	err = topClaim.AddPremise(CTX, &premiseClaim1)
	assert.Error(t, err)
	assert.Equal(t, "A claim that has already been deleted, or has a newer version, cannot be modified", err.Error())

	topClaim.DeletedAt = nil
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
	claim.Key = fmt.Sprintf("thisusessprintfjustsofmtisalwaysimported")
	assert.Equal(t, "claims/thisusessprintfjustsofmtisalwaysimported", claim.ArangoID())
}

func TestClaimLoadByID(t *testing.T) {
	setupDB()
	defer teardownDB()

}

func TestClaimUpdate(t *testing.T) {
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

	ctx1 := Context{ShortName: "UpdateClaim First One", Title: "First context for update claim", URL: "https://en.wikipedia.org/wiki/First_Ones"}
	err = ctx1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	ctx2 := Context{ShortName: "UpdateClaim Second", Title: "Being first isn't everything", URL: "https://en.wikipedia.org/wiki/Second"}
	err = ctx2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim.AddContext(CTX, ctx1)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim.AddContext(CTX, ctx2)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Next check edges, then version and recheck everything
	premiseEdges, err := claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(premiseEdges))

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

	ces, err := claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ces))
	assert.Equal(t, ctx1.ArangoID(), ces[0].From)
	assert.Equal(t, claim.ArangoID(), ces[0].To)
	assert.Equal(t, ctx2.ArangoID(), ces[1].From)
	assert.Equal(t, claim.ArangoID(), ces[1].To)

	ctxs, err := claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ctxs))
	assert.Equal(t, ctx1.ArangoID(), ctxs[0].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), ctxs[1].ArangoID())

	// Update the claim
	curator := User{Username: "curator", Curator: true}
	err = curator.Create(CTX)
	assert.NoError(t, err)
	CTX.UserContext = curator

	err = claim.Load(CTX)
	assert.NoError(t, err)
	origClaimKey := claim.ArangoKey()

	updates := map[string]interface{}{
		"title": "New Title",
		"desc":  "New Description",
	}
	err = claim.Update(CTX, updates)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "New Title", claim.Title)
	assert.Equal(t, "New Description", claim.Description)
	assert.NotEqual(t, origClaimKey, claim.ArangoKey())
	assert.Equal(t, DEFAULT_USER.ArangoID(), claim.CreatedByID)
	assert.Equal(t, CTX.UserContext.ArangoID(), claim.UpdatedByID)

	origClaim := Claim{}
	origClaim.Key = origClaimKey
	err = origClaim.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, origClaim.DeletedAt)
	assert.Equal(t, "I dare you to doubt me", origClaim.Title)
	assert.Equal(t, "I am true. Woe be the person that doubts my veracity", origClaim.Description)
	assert.Equal(t, origClaimKey, origClaim.ArangoKey())

	// Verify new edges were created
	premiseEdges, err = claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(premiseEdges))

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

	ces, err = claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ces))
	assert.Equal(t, ctx1.ArangoID(), ces[0].From)
	assert.Equal(t, claim.ArangoID(), ces[0].To)
	assert.Equal(t, ctx2.ArangoID(), ces[1].From)
	assert.Equal(t, claim.ArangoID(), ces[1].To)

	ctxs, err = claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ctxs))
	assert.Equal(t, ctx1.ArangoID(), ctxs[0].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), ctxs[1].ArangoID())

	// Verify that the old edges were deleted
	premiseEdges, err = origClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(premiseEdges))

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

	ces, err = origClaim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ces))
	assert.Equal(t, ctx1.ArangoID(), ces[0].From)
	assert.Equal(t, origClaim.ArangoID(), ces[0].To)
	assert.NotNil(t, ces[0].DeletedAt)
	assert.Equal(t, ctx2.ArangoID(), ces[1].From)
	assert.Equal(t, origClaim.ArangoID(), ces[1].To)
	assert.NotNil(t, ces[1].DeletedAt)
}

func TestClaimUpdateMP(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "I dare you to doubt me because I am Updated MP",
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
	assert.Error(t, err)
	assert.Equal(t, "You must convert this claim to be a multi-premise claim before adding new premises", err.Error())

	err = claim.ConvertToMultiPremise(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	autoPremises, err := claim.Premises(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(autoPremises))
	autoPremise := autoPremises[0]

	err = claim.AddPremise(CTX, &premiseClaim1)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premiseClaim2 := Claim{
		Title:        "Since it is I that am daring you, you therefore must not doubt",
		Description:  "I am undoubtable",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
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
		TargetArgumentID: &argPC1.ID,
		Title:            "Let's create a new argument argument",
		Pro:              false,
	}
	err = arg1arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Next check edges, then version and recheck everything
	premiseEdges, err := claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(premiseEdges))
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[1].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[2].From)
	assert.Equal(t, autoPremise.ArangoID(), premiseEdges[0].To)
	assert.Equal(t, premiseClaim1.ArangoID(), premiseEdges[1].To)
	assert.Equal(t, premiseClaim2.ArangoID(), premiseEdges[2].To)

	inferences, err := claim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(inferences))

	bces, err := claim.BaseClaimEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bces))
	assert.Equal(t, argDC1.ArangoID(), bces[0].From)
	assert.Equal(t, claim.ArangoID(), bces[0].To)

	args, err := claim.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(args))

	args, err = premiseClaim1.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, argPC1.ArangoID(), args[0].ArangoID())

	premiseEdges, err = mpClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premiseEdges))
	assert.Equal(t, mpClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].To)

	// Update the claim
	err = claim.Load(CTX)
	assert.NoError(t, err)
	origClaimKey := claim.ArangoKey()

	updates := map[string]interface{}{
		"title": "New Title",
		"desc":  "New Description",
	}
	err = claim.Update(CTX, updates)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "New Title", claim.Title)
	assert.Equal(t, "New Description", claim.Description)
	assert.NotEqual(t, origClaimKey, claim.ArangoKey())

	origClaim := Claim{}
	origClaim.Key = origClaimKey
	err = origClaim.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, origClaim.DeletedAt)
	assert.Equal(t, "I dare you to doubt me because I am Updated MP", origClaim.Title)
	assert.Equal(t, origClaimKey, origClaim.ArangoKey())

	// Verify new edges were created
	premiseEdges, err = claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(premiseEdges))
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[1].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[2].From)
	assert.Equal(t, autoPremise.ArangoID(), premiseEdges[0].To)
	assert.Equal(t, premiseClaim1.ArangoID(), premiseEdges[1].To)
	assert.Equal(t, premiseClaim2.ArangoID(), premiseEdges[2].To)

	inferences, err = claim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(inferences))

	bces, err = claim.BaseClaimEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bces))
	assert.Equal(t, argDC1.ArangoID(), bces[0].From)
	assert.Equal(t, claim.ArangoID(), bces[0].To)
	assert.Nil(t, bces[0].DeletedAt)

	args, err = claim.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(args))

	args, err = premiseClaim1.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, argPC1.ArangoID(), args[0].ArangoID())

	premiseEdges, err = mpClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premiseEdges))
	assert.Equal(t, mpClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].To)
	assert.Nil(t, premiseEdges[0].DeletedAt)

	// Verify that the old edges were deleted
	premiseEdges, err = origClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(premiseEdges))
	assert.Equal(t, origClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, origClaim.ArangoID(), premiseEdges[1].From)
	assert.Equal(t, origClaim.ArangoID(), premiseEdges[2].From)
	assert.Equal(t, autoPremise.ArangoID(), premiseEdges[0].To)
	assert.Equal(t, premiseClaim1.ArangoID(), premiseEdges[1].To)
	assert.Equal(t, premiseClaim2.ArangoID(), premiseEdges[2].To)
	assert.NotNil(t, premiseEdges[0].DeletedAt)
	assert.NotNil(t, premiseEdges[1].DeletedAt)
	assert.NotNil(t, premiseEdges[2].DeletedAt)

	inferences, err = origClaim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(inferences))

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

	err := claim.Create(CTX)
	assert.NoError(t, err)
	patch := map[string]interface{}{
		"start": time.Now().Add(-25 * time.Hour),
		"end":   time.Now().Add(-24 * time.Hour),
	}
	col, _ := CTX.Arango.CollectionFor(&claim)
	col.UpdateDocument(CTX.Context, claim.ArangoKey(), patch)

	firstKey := claim.ArangoKey()

	err = claim.version(CTX)
	assert.NoError(t, err)
	patch["start"] = time.Now().Add(-24 * time.Hour)
	patch["end"] = time.Now().Add(-1 * time.Hour)
	col.UpdateDocument(CTX.Context, claim.ArangoKey(), patch)

	secondKey := claim.ArangoKey()

	claim.DeletedAt = nil
	err = claim.version(CTX)
	assert.NoError(t, err)
	patch["start"] = time.Now().Add(-1 * time.Hour)
	delete(patch, "end")
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

	lookup = Claim{}
	lookup.ID = claim.ID
	lookup.QueryAt = support.TimePtr(time.Now().Add(-48 * time.Hour))
	err = lookup.Load(CTX)
	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())

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
	setupDB()
	defer teardownDB()

	CTX.RequestAt = nil

	claim1 := Claim{
		Title:        "A Top everything Claim for queries for top level claims",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}
	err := claim1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim2 := Claim{
		Title: "C Top everything Claim for queries for top level claims",
	}
	err = claim2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim3 := Claim{
		Title: "B Top everything Claim for queries for top level claims",
	}
	err = claim3.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premise := Claim{
		Title: "A premise Claim for queries for top level claims",
	}
	err = premise.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim1.AddPremise(CTX, &premise)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetClaimID: &premise.ID,
		Title:         "I need to argue for the top level claims",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := Argument{
		TargetClaimID: &claim2.ID,
		Title:         "I need to argue against the top level claims",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg3 := Argument{
		TargetArgumentID: &arg1.ID,
		Title:            "I might as well argue against the arguments",
		Pro:              false,
	}
	err = arg3.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claims := []Claim{}
	params := ArangoQueryParameters{}
	bindVars := BindVars{}
	query := Claim{}.QueryForTopLevelClaims(params)

	err = FindArangoObjects(CTX, query, bindVars, &claims)
	assert.NoError(t, err)
	// This will fail when run alone, since it depends on claims created in other tests
	assert.Equal(t, 20, len(claims))
	for _, claim := range claims {
		assert.NotEqual(t, premise.ID, claim.ID)
		assert.NotEqual(t, arg1.ClaimID, claim.ID)
		assert.NotEqual(t, arg2.ClaimID, claim.ID)
		assert.NotEqual(t, arg3.ClaimID, claim.ID)
	}

	params.Limit = support.IntPtr(3)
	query = Claim{}.QueryForTopLevelClaims(params)
	claims = []Claim{}
	err = FindArangoObjects(CTX, query, bindVars, &claims)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(claims))
	assert.Equal(t, claim3.ID, claims[0].ID)
	assert.Equal(t, claim2.ID, claims[1].ID)
	assert.Equal(t, claim1.ID, claims[2].ID)
	assert.Equal(t, claim1.Title, claims[2].Title)

	params.Offset = support.IntPtr(1)
	params.Limit = support.IntPtr(2)
	query = Claim{}.QueryForTopLevelClaims(params)
	claims = []Claim{}
	err = FindArangoObjects(CTX, query, bindVars, &claims)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(claims))
	assert.Equal(t, claim2.ID, claims[0].ID)
	assert.Equal(t, claim1.ID, claims[1].ID)

	params.Offset = nil
	params.Limit = support.IntPtr(20)
	params.Sort = support.StringPtr("obj._id")
	query = Claim{}.QueryForTopLevelClaims(params)
	claims = []Claim{}
	err = FindArangoObjects(CTX, query, bindVars, &claims)
	assert.Equal(t, 20, len(claims))
	var prevClaim Claim
	for _, claim := range claims {
		if prevClaim.ID != "" {
			assert.True(t, prevClaim.ArangoID() < claim.ArangoID())
		}
		assert.NotEqual(t, premise.ID, claim.ID)
		assert.NotEqual(t, arg1.ClaimID, claim.ID)
		assert.NotEqual(t, arg2.ClaimID, claim.ID)
		assert.NotEqual(t, arg3.ClaimID, claim.ID)
		prevClaim = claim
	}

	params.Sort = nil
	params.Limit = support.IntPtr(4)
	query = Claim{}.QueryForTopLevelClaims(params)
	claims = []Claim{}
	err = claim1.RemovePremise(CTX, premise.ArangoID())
	assert.NoError(t, err)
	err = FindArangoObjects(CTX, query, bindVars, &claims)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(claims))
	assert.Equal(t, claim1.ID, claims[0].ID)
	assert.Equal(t, claim1.Title, claims[0].Title)
	assert.Equal(t, premise.ID, claims[1].ID)
	assert.Equal(t, premise.Title, claims[1].Title)
	assert.Equal(t, claim3.ID, claims[2].ID)
	assert.Equal(t, claim2.ID, claims[3].ID)

	params.Limit = support.IntPtr(4)
	query = Claim{}.QueryForTopLevelClaims(params)
	claims = []Claim{}
	err = claim1.ConvertToMultiPremise(CTX)
	assert.NoError(t, err)
	err = claim1.AddPremise(CTX, &claim2)
	assert.NoError(t, err)
	err = FindArangoObjects(CTX, query, bindVars, &claims)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(claims))
	assert.Equal(t, claim1.ID, claims[0].ID)
	assert.Equal(t, premise.ID, claims[1].ID)
	assert.Equal(t, claim3.ID, claims[2].ID)
	assert.NotEqual(t, claim2.ID, claims[3].ID)

	claims = []Claim{}
	err = arg3.Delete(CTX)
	assert.NoError(t, err)
	err = FindArangoObjects(CTX, query, bindVars, &claims)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(claims))
	assert.Equal(t, claim1.ID, claims[0].ID)
	assert.Equal(t, arg3.ClaimID, claims[1].ID)
	assert.Equal(t, premise.ID, claims[2].ID)
	assert.Equal(t, claim3.ID, claims[3].ID)

	claims = []Claim{}
	err = claim3.Delete(CTX)
	assert.NoError(t, err)
	err = FindArangoObjects(CTX, query, bindVars, &claims)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(claims))
	assert.Equal(t, claim1.ID, claims[0].ID)
	assert.Equal(t, claim1.Title, claims[0].Title)
	assert.Equal(t, arg3.ClaimID, claims[1].ID)
	assert.Equal(t, arg3.Title, claims[1].Title)
	assert.Equal(t, premise.ID, claims[2].ID)
	assert.NotEqual(t, claim2.ID, claims[3].ID)
	assert.NotEqual(t, claim3.ID, claims[3].ID)
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
		MultiPremise: false,
		PremiseRule:  PREMISE_RULE_NONE,
	}
	err := claim.Create(CTX)
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

	// TODO: Contexts

	// Simple Load
	err = claim.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))

	// Load All
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
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 2, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg1, claim.ProArgs[0])
	assert.Equal(t, arg3, claim.ProArgs[1])
	assert.Equal(t, arg2, claim.ConArgs[0])

	err = arg1.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	arg1.Load(CTX)

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg3, claim.ProArgs[0])
	assert.Equal(t, arg2, claim.ConArgs[0])

	// Previous points in time
	arg1.Claim = &carg1
	arg2.Claim = &carg2
	arg3.Claim = &carg3
	claim.QueryAt = &arg1arg.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 2, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg1, claim.ProArgs[0])
	assert.Equal(t, arg3, claim.ProArgs[1])
	assert.Equal(t, arg2, claim.ConArgs[0])

	claim.QueryAt = &arg1.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))
	assert.Equal(t, arg1, claim.ProArgs[0])

	claim.QueryAt = &claim.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))

	claim.QueryAt = nil
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg3, claim.ProArgs[0])
	assert.Equal(t, arg2, claim.ConArgs[0])
}

func TestClaimLoadFullMP(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "This is the MP Claim LoadAll test claim",
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
		Title:        "First premise for the Claim MP LoadAll dude",
		Description:  "The person that is daring you to doubt me being me",
		MultiPremise: false,
	}
	err = claim.AddPremise(CTX, &premiseClaim1)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premiseClaim2 := Claim{
		Title:        "I am the second Claim MP LoadAll premise. I MUST be true.",
		Description:  "I am undoubtable",
		MultiPremise: false,
	}
	err = claim.AddPremise(CTX, &premiseClaim2)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	distantClaim := Claim{
		Title:       "So very far away from Claim MP LoadAll",
		Description: "So distant, you cannot see me.",
	}
	err = distantClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	mpClaim := Claim{
		Title:        "This is an MP that uses the main MP LoadAll claim",
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

	argPC11 := Argument{
		TargetClaimID: &premiseClaim1.ID,
		Title:         "MP All?",
		Pro:           true,
	}
	err = argPC11.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argPC12 := Argument{
		TargetClaimID: &premiseClaim1.ID,
		Title:         "MP Load ALL!",
		Pro:           false,
	}
	err = argPC12.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argPC13 := Argument{
		TargetClaimID: &premiseClaim1.ID,
		Title:         "MP Do it ALL!",
		Pro:           true,
	}
	err = argPC13.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argPC2 := Argument{
		TargetClaimID: &premiseClaim2.ID,
		Title:         "MP Let's create a new argument for the Premise of the claim LoadAll",
		Pro:           false,
	}
	err = argPC2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argDC1 := Argument{
		TargetClaimID: &distantClaim.ID,
		ClaimID:       claim.ID,
		Title:         "Distant MP LoadAll claim.",
	}
	err = argDC1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1arg := Argument{
		TargetArgumentID: &argPC11.ID,
		Title:            "MP Do it...ALLLLL!",
		Pro:              false,
	}
	err = arg1arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Simple Load
	err = claim.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the MP Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 0, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))

	// Load All
	premiseClaim1.Load(CTX)
	premiseClaim2.Load(CTX)
	argPC11.Load(CTX)
	argPC12.Load(CTX)
	argPC13.Load(CTX)
	argPC2.Load(CTX)
	var cargPC11, cargPC12, cargPC13, cargPC2 Claim
	cargPC11.ID = argPC11.ClaimID
	cargPC12.ID = argPC12.ClaimID
	cargPC13.ID = argPC13.ClaimID
	cargPC2.ID = argPC2.ClaimID
	cargPC11.Load(CTX)
	cargPC12.Load(CTX)
	cargPC13.Load(CTX)
	cargPC2.Load(CTX)
	argPC11.Claim = &cargPC11
	argPC12.Claim = &cargPC12
	argPC13.Claim = &cargPC13
	argPC2.Claim = &cargPC2
	premiseClaim1.ProArgs = []Argument{argPC11, argPC13}
	premiseClaim1.ConArgs = []Argument{argPC12}
	premiseClaim2.ConArgs = []Argument{argPC2}

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the MP Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 2, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))
	assert.Equal(t, premiseClaim1, claim.PremiseClaims[0])
	assert.Equal(t, premiseClaim2, claim.PremiseClaims[1])

	err = argPC11.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	err = premiseClaim2.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	argPC11.Load(CTX)
	premiseClaim2.Load(CTX)
	argPC2.Load(CTX)

	premiseClaim1.ProArgs = []Argument{argPC13}

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the MP Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 1, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))
	assert.Equal(t, premiseClaim1, claim.PremiseClaims[0])

	// Previous points in time
	premiseClaim1.QueryAt = &arg1arg.CreatedAt
	premiseClaim1.Load(CTX)
	premiseClaim2.QueryAt = &arg1arg.CreatedAt
	premiseClaim2.Load(CTX)
	argPC2.QueryAt = &arg1arg.CreatedAt
	argPC2.Load(CTX)
	argPC11.Claim = &cargPC11
	argPC12.Claim = &cargPC12
	argPC13.Claim = &cargPC13
	argPC2.Claim = &cargPC2
	premiseClaim1.ProArgs = []Argument{argPC11, argPC13}
	premiseClaim1.ConArgs = []Argument{argPC12}
	premiseClaim2.ConArgs = []Argument{argPC2}
	claim.QueryAt = &arg1arg.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the MP Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 2, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))
	assert.Equal(t, premiseClaim1, claim.PremiseClaims[0])
	assert.Equal(t, premiseClaim2, claim.PremiseClaims[1])

	premiseClaim1.ProArgs = []Argument{argPC11}
	premiseClaim1.ConArgs = nil
	premiseClaim2.ConArgs = nil
	claim.QueryAt = &argPC11.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the MP Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 2, len(claim.PremiseClaims))
	assert.Equal(t, premiseClaim1, claim.PremiseClaims[0])
	assert.Equal(t, premiseClaim2, claim.PremiseClaims[1])

	claim.QueryAt = &claim.CreatedAt
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the MP Claim LoadAll test claim", claim.Title)
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

	premiseClaim1.ProArgs = []Argument{argPC13}
	premiseClaim1.ConArgs = []Argument{argPC12}
	claim.QueryAt = nil
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)
	assert.Equal(t, "This is the MP Claim LoadAll test claim", claim.Title)
	assert.Equal(t, 2, len(claim.PremiseClaims))
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 0, len(claim.ConArgs))
	assert.Equal(t, premiseClaim1, claim.PremiseClaims[0])
	assert.Equal(t, premiseClaim3, claim.PremiseClaims[1])
}

func TestClaimDeleteLoop(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "I'm first, so any loop is your fault",
		Description:  "I am true. Woe be the person that doubts my veracity",
		Negation:     "I dare you to accept me",
		Question:     "Do you dare to doubt me?",
		Note:         "This Claim is for deleting in a loop",
		Image:        "https://static1.squarespace.com/static/58ed33aeb8a79b05bed202aa/t/5a1fed3a652dead776d6aaed/1512041798286/The+loop+logo+white+background.jpg?format=1000w",
		MultiPremise: false,
		PremiseRule:  PREMISE_RULE_NONE,
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	distantClaim := Claim{
		Title:       "So very far away, but the loop brings us closer",
		Description: "So distant, you cannot see me, unless you follow the loop.",
	}
	err = distantClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetClaimID: &claim.ID,
		ClaimID:       distantClaim.ID,
		Title:         "Let's create a loop argument",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argDC1 := Argument{
		TargetClaimID: &distantClaim.ID,
		ClaimID:       claim.ID,
		Title:         "I want to get away from the loop",
	}
	err = argDC1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1arg := Argument{
		TargetArgumentID: &arg1.ID,
		ClaimID:          distantClaim.ID,
		Title:            "Let's create a new argument argument related to the loop",
		Pro:              false,
	}
	err = arg1arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Next, delete the main claim
	err = claim.Delete(CTX)
	assert.Error(t, err)
	assert.Equal(t, "You cannot delete a claim that is being used as a base claim for other arguments", err.Error())

	// Delete the argument using it
	err = argDC1.Delete(CTX)
	assert.NoError(t, err)

	// Now delete the claim
	err = claim.Delete(CTX)
	assert.NoError(t, err)

	inferences, err := claim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inferences))
	assert.Equal(t, claim.ArangoID(), inferences[0].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.NotNil(t, inferences[0].DeletedAt)

	baseClaimEdges, err := claim.BaseClaimEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(baseClaimEdges))
	assert.Equal(t, argDC1.ArangoID(), baseClaimEdges[0].From)
	assert.Equal(t, claim.ArangoID(), baseClaimEdges[0].To)
	assert.NotNil(t, baseClaimEdges[0].DeletedAt)

	args, err := claim.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())
	assert.NotNil(t, args[0].DeletedAt)

	arg1arg.QueryAt = support.TimePtr(claim.DeletedAt.Add(-1 * time.Millisecond))
	err = arg1arg.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, arg1arg.DeletedAt)

	err = argDC1.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, argDC1.DeletedAt)

	err = distantClaim.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, distantClaim.DeletedAt)
}

func TestClaimAddPremiseLoop(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "I dare you to doubt me because I am MP Infinite Loop",
		Description:  "I am true. Woe be the person that doubts my veracity",
		Negation:     "I dare you to accept me",
		Question:     "Do you dare to doubt me?",
		Note:         "This Claim is all about infinite premise loops",
		Image:        "https://media.sanoma.fi/sites/default/files/styles/icon_lg/public/2018-03/Loop.png?itok=F630fzmT",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premiseClaim1 := Claim{
		Title:        "I am a normal premise. Use me!",
		Description:  "The person that is daring you to doubt me being me",
		MultiPremise: false,
	}
	err = claim.AddPremise(CTX, &premiseClaim1)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	premiseClaim2 := Claim{
		Title:        "I'm also an innocente premise... at first...",
		Description:  "I am undoubtable",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}
	err = claim.AddPremise(CTX, &premiseClaim2)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	mpClaim := Claim{
		Title:        "This is an MP that is going to be involved in a loop",
		Description:  "Not military police, mind you",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}
	err = mpClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = mpClaim.AddPremise(CTX, &claim)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Now try to create a loop and fail
	err = claim.AddPremise(CTX, &claim)
	assert.Error(t, err)
	assert.Equal(t, "A claim cannot be a premise of itself, nor one of its own premises. That's called \"Begging the Question\"", err.Error())

	err = claim.AddPremise(CTX, &mpClaim)
	assert.Error(t, err)
	assert.Equal(t, "A claim cannot be a premise of itself, nor one of its own premises. That's called \"Begging the Question\"", err.Error())

	err = premiseClaim2.AddPremise(CTX, &claim)
	assert.Error(t, err)
	assert.Equal(t, "A claim cannot be a premise of itself, nor one of its own premises. That's called \"Begging the Question\"", err.Error())

	err = claim.AddPremise(CTX, &premiseClaim1)
	assert.Error(t, err)
	assert.Equal(t, "This claim has already been added as a premise", err.Error())

	err = mpClaim.AddPremise(CTX, &premiseClaim1)
	assert.Error(t, err)
	assert.Equal(t, "This claim has already been added as a premise", err.Error())

}

func TestClaimHasCycle(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:       "The cycle of life continues",
		Description: "This is for the Claim HasCycle test",
		Image:       "https://cdn5.vectorstock.com/i/1000x1000/41/39/life-cycle-of-a-chicken-for-kids-vector-4924139.jpg",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "First has cycle argument",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Second has cycle argument",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	has, err := claim.HasCycle(CTX)
	assert.NoError(t, err)
	assert.False(t, has)

	arg := Argument{
		TargetArgumentID: &arg1.ID,
		ClaimID:          claim.ID,
		Title:            "This is the argument that's going to mess it up for everyone",
		Pro:              false,
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	has, err = claim.HasCycle(CTX)
	assert.NoError(t, err)
	assert.True(t, has)

	err = arg.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	has, err = claim.HasCycle(CTX)
	assert.NoError(t, err)
	assert.False(t, has)

	claim.QueryAt = support.TimePtr(arg.CreatedAt)
	has, err = claim.HasCycle(CTX)
	assert.NoError(t, err)
	assert.True(t, has)

	claim.QueryAt = nil

	arga := Argument{
		TargetArgumentID: &arg1.ID,
		Title:            "Now we're going to go out of our way to complete a cycle",
		Pro:              false,
	}
	err = arga.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argb := Argument{
		TargetArgumentID: &arga.ID,
		Title:            "And by that I mean way out of our way",
		Pro:              true,
	}
	err = argb.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argbClaimId := argb.ClaimID
	argc := Argument{
		TargetClaimID: &argbClaimId,
		Title:         "Yeah, pretty far out to create a cycle",
		Pro:           true,
	}
	err = argc.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argcClaimId := argc.ClaimID
	argd := Argument{
		TargetClaimID: &argcClaimId,
		ClaimID:       claim.ID,
		Title:         "This should be far enough for the cycle",
		Pro:           false,
	}
	err = argd.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	has, err = claim.HasCycle(CTX)
	assert.NoError(t, err)
	assert.True(t, has)

	err = argd.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	has, err = claim.HasCycle(CTX)
	assert.NoError(t, err)
	assert.False(t, has)

	argcClaim := Claim{}
	argcClaim.ID = argc.ClaimID
	err = argcClaim.Load(CTX)
	assert.NoError(t, err)
	err = argcClaim.AddPremise(CTX, &claim)
	assert.Error(t, err)
	assert.Equal(t, "You must convert this claim to be a multi-premise claim before adding new premises", err.Error())

	err = argcClaim.ConvertToMultiPremise(CTX)
	assert.NoError(t, err)

	err = argcClaim.AddPremise(CTX, &claim)
	assert.NoError(t, err)

	has, err = claim.HasCycle(CTX)
	assert.NoError(t, err)
	assert.True(t, has)
}

func TestClaimAddContext(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "Here we have a claim that is all about adding Context.",
		Description:  "Add Context. That's what it's about.",
		Image:        "http://danieleizans.com/wp-content/uploads/2011/01/ContextInCS.jpg",
		MultiPremise: false,
		PremiseRule:  PREMISE_RULE_NONE,
	}

	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	ctx1 := Context{ShortName: "AddContext Goodies", Title: "Goody Goody Yum Yum", URL: "https://en.wikipedia.org/wiki/The_Goodies_(TV_series)"}
	err = ctx1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim.AddContext(CTX, ctx1)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Assert that the proper link has been created
	contextEdges, err := claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(contextEdges))
	assert.Equal(t, claim.ArangoID(), contextEdges[0].To)
	assert.Equal(t, ctx1.ArangoID(), contextEdges[0].From)

	contexts, err := claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(contexts))
	assert.Equal(t, ctx1.ArangoID(), contexts[0].ArangoID())

	ctx2 := Context{ShortName: "AddContext Politics", Title: "Politics in general", URL: "https://en.wikipedia.org/wiki/Politics"}
	ctx3 := Context{ShortName: "AddContext Comedy", Title: "Politics is always funny", URL: "https://en.wikipedia.org/wiki/Comedy"}
	err = ctx2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	err = ctx3.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	u := User{Username: "AddContextNoPermissions", Email: "noperms@addcontext.com"}
	err = u.Create(CTX)
	assert.NoError(t, err)

	CTX.UserContext = u
	err = claim.AddContext(CTX, ctx2)
	assert.Error(t, err)
	assert.Equal(t, "You do not have permission to modify this item", err.Error())

	CTX.UserContext = DEFAULT_USER
	err = claim.AddContext(CTX, ctx2)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	err = claim.AddContext(CTX, ctx3)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	threeCtxTime := time.Now()

	contextEdges, err = claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contextEdges))
	assert.Equal(t, claim.ArangoID(), contextEdges[0].To)
	assert.Equal(t, ctx1.ArangoID(), contextEdges[0].From)
	assert.Equal(t, claim.ArangoID(), contextEdges[1].To)
	assert.Equal(t, ctx2.ArangoID(), contextEdges[1].From)
	assert.Equal(t, claim.ArangoID(), contextEdges[2].To)
	assert.Equal(t, ctx3.ArangoID(), contextEdges[2].From)

	contexts, err = claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx1.ArangoID(), contexts[1].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[2].ArangoID())

	// Remove a Context
	CTX.UserContext = u
	err = claim.RemoveContext(CTX, ctx2.ArangoKey())
	assert.Error(t, err)
	assert.Equal(t, "You do not have permission to modify this item", err.Error())
	CTX.UserContext = DEFAULT_USER

	err = claim.RemoveContext(CTX, ctx2.ArangoKey())
	assert.NoError(t, err)
	CTX.RequestAt = nil
	twoCtxTime := time.Now()

	ctx2, err = FindContext(CTX, ctx2.ArangoID())
	assert.NoError(t, err)
	assert.Nil(t, ctx2.DeletedAt)

	contextEdges, err = claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(contextEdges))
	assert.Equal(t, claim.ArangoID(), contextEdges[0].To)
	assert.Equal(t, ctx1.ArangoID(), contextEdges[0].From)
	assert.Equal(t, claim.ArangoID(), contextEdges[1].To)
	assert.Equal(t, ctx3.ArangoID(), contextEdges[1].From)

	contexts, err = claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx1.ArangoID(), contexts[1].ArangoID())

	// Query Contexts at a point in time
	claim.QueryAt = &threeCtxTime
	contextEdges, err = claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contextEdges))
	assert.Equal(t, claim.ArangoID(), contextEdges[0].To)
	assert.Equal(t, ctx1.ArangoID(), contextEdges[0].From)
	assert.Equal(t, claim.ArangoID(), contextEdges[1].To)
	assert.Equal(t, ctx2.ArangoID(), contextEdges[1].From)
	assert.Equal(t, claim.ArangoID(), contextEdges[2].To)
	assert.Equal(t, ctx3.ArangoID(), contextEdges[2].From)
	assert.NotNil(t, contextEdges[1].DeletedAt)

	contexts, err = claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx1.ArangoID(), contexts[1].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[2].ArangoID())

	// Add removed Context again
	claim.QueryAt = nil
	err = claim.AddContext(CTX, ctx2)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	contextEdges, err = claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contextEdges))
	assert.Equal(t, claim.ArangoID(), contextEdges[0].To)
	assert.Equal(t, ctx1.ArangoID(), contextEdges[0].From)
	assert.Equal(t, claim.ArangoID(), contextEdges[1].To)
	assert.Equal(t, ctx3.ArangoID(), contextEdges[1].From)
	assert.Equal(t, claim.ArangoID(), contextEdges[2].To)
	assert.Equal(t, ctx2.ArangoID(), contextEdges[2].From)

	contexts, err = claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx1.ArangoID(), contexts[1].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[2].ArangoID())

	// Query back in time
	claim.QueryAt = &twoCtxTime
	contextEdges, err = claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(contextEdges))
	assert.Equal(t, claim.ArangoID(), contextEdges[0].To)
	assert.Equal(t, ctx1.ArangoID(), contextEdges[0].From)
	assert.Equal(t, claim.ArangoID(), contextEdges[1].To)
	assert.Equal(t, ctx3.ArangoID(), contextEdges[1].From)

	contexts, err = claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx1.ArangoID(), contexts[1].ArangoID())

	// Try to add repeat Context
	claim.QueryAt = nil
	err = claim.AddContext(CTX, ctx3)
	assert.Error(t, err)
	assert.Equal(t, "This context was already added to this claim", err.Error())

	// Add to an MP Claim
	mpClaim := Claim{Title: "MP Claim in an AddContext world", MultiPremise: true, PremiseRule: PREMISE_RULE_ALL}
	err = mpClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = mpClaim.AddContext(CTX, ctx1)
	assert.Error(t, err)
	assert.Equal(t, "Multi-premise claims inherit the union of contexts from all their premises", err.Error())

	// Get Contexts from MP Claim
	claim.QueryAt = nil
	err = claim.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, claim.DeletedAt)

	err = mpClaim.AddPremise(CTX, &claim)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim2 := Claim{Title: "A MultiPremise Premise in an AddContext world", MultiPremise: true, PremiseRule: PREMISE_RULE_ALL}
	err = claim2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim3 := Claim{Title: "A regular Premise in an AddContext world"}
	err = claim3.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	claim4 := Claim{Title: "A redundant Premise in an AddContext world"}
	err = claim4.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim2.AddPremise(CTX, &claim3)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim2.AddPremise(CTX, &claim4)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = mpClaim.AddPremise(CTX, &claim2)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim3.AddContext(CTX, ctx2)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim3.AddContext(CTX, ctx3)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim4.AddContext(CTX, ctx1)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim4.AddContext(CTX, ctx3)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	ctx4 := Context{ShortName: "AddContext One Last Tim", Title: "One last new context", URL: "https://en.wikipedia.org/wiki/One_Last_Time"}
	err = ctx4.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = claim3.AddContext(CTX, ctx4)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Query MP Claim Contexts
	contextEdges, err = mpClaim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contextEdges))

	contexts, err = mpClaim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx1.ArangoID(), contexts[1].ArangoID())
	assert.Equal(t, ctx4.ArangoID(), contexts[2].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[3].ArangoID())

	// Delete and query at points in time
	err = claim.RemoveContext(CTX, ctx1.ArangoKey())
	assert.NoError(t, err)
	CTX.RequestAt = nil

	contextEdges, err = mpClaim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contextEdges))

	contexts, err = mpClaim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx1.ArangoID(), contexts[1].ArangoID())
	assert.Equal(t, ctx4.ArangoID(), contexts[2].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[3].ArangoID())

	err = claim4.RemoveContext(CTX, ctx1.ArangoKey())
	assert.NoError(t, err)
	CTX.RequestAt = nil

	threeCtxTime = time.Now()

	contextEdges, err = mpClaim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contextEdges))

	contexts, err = mpClaim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx4.ArangoID(), contexts[1].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[2].ArangoID())

	err = claim3.RemoveContext(CTX, ctx4.ArangoKey())
	assert.NoError(t, err)
	CTX.RequestAt = nil

	contextEdges, err = mpClaim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contextEdges))

	contexts, err = mpClaim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[1].ArangoID())

	mpClaim.QueryAt = &threeCtxTime
	contextEdges, err = mpClaim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contextEdges))

	contexts, err = mpClaim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx4.ArangoID(), contexts[1].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[2].ArangoID())

	err = claim4.AddContext(CTX, ctx4)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	mpClaim.QueryAt = nil
	contextEdges, err = mpClaim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contextEdges))

	contexts, err = mpClaim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx4.ArangoID(), contexts[1].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[2].ArangoID())

	// TODO: delete context edges
	err = claim4.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	contextEdges, err = mpClaim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contextEdges))

	contexts, err = mpClaim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(contexts))
	assert.Equal(t, ctx3.ArangoID(), contexts[0].ArangoID())
	assert.Equal(t, ctx2.ArangoID(), contexts[1].ArangoID())

	// Add to deleted Claim
	err = claim4.AddContext(CTX, ctx1)
	assert.Error(t, err)
	assert.Equal(t, "A claim that has already been deleted, or has a newer version, cannot be modified", err.Error())
}

func TestClaimValidateForDelete(t *testing.T) {
	claim := Claim{}
	assert.Nil(t, claim.ValidateForDelete())

	claim.DeletedAt = support.TimePtr(time.Now())
	err := claim.ValidateForDelete()
	assert.NotNil(t, err)
	assert.Equal(t, "This claim has already been deleted or versioned", err.Error())
}

func TestClaimUserCanDelete(t *testing.T) {
	u := User{}
	CTX.UserContext = u

	claim := Claim{}
	can, err := claim.UserCanDelete(CTX)
	assert.Nil(t, err)
	assert.False(t, can)

	u.PrepareForCreate(CTX)
	CTX.UserContext = u
	claim = Claim{}
	can, err = claim.UserCanDelete(CTX)
	assert.Nil(t, err)
	assert.False(t, can)

	claim = Claim{}
	can, err = claim.UserCanDelete(CTX)
	assert.Nil(t, err)
	assert.False(t, can)

	claim.CreatedByID = u.ArangoID()
	can, err = claim.UserCanDelete(CTX)
	assert.Nil(t, err)
	assert.True(t, can)

	u = User{}
	u.PrepareForCreate(CTX)
	CTX.UserContext = u
	can, err = claim.UserCanDelete(CTX)
	assert.Nil(t, err)
	assert.False(t, can)

	u.Curator = true
	CTX.UserContext = u
	can, err = claim.UserCanDelete(CTX)
	assert.Nil(t, err)
	assert.True(t, can)
}

func TestClaimUserCanCreate(t *testing.T) {
	u := User{}
	CTX.UserContext = u

	claim := Claim{}
	can, err := claim.UserCanCreate(CTX)
	assert.Nil(t, err)
	assert.False(t, can)

	u.PrepareForCreate(CTX)
	CTX.UserContext = u

	can, err = claim.UserCanCreate(CTX)
	assert.Nil(t, err)
	assert.True(t, can)

	u.Curator = true
	CTX.UserContext = u
	can, err = claim.UserCanCreate(CTX)
	assert.Nil(t, err)
	assert.True(t, can)
}

func TestClaimUserCanUpdate(t *testing.T) {
	u := User{}
	CTX.UserContext = u

	updates := map[string]interface{}{}

	claim := Claim{}
	can, err := claim.UserCanUpdate(CTX, updates)
	assert.Nil(t, err)
	assert.False(t, can)

	u.PrepareForCreate(CTX)
	CTX.UserContext = u
	can, err = claim.UserCanUpdate(CTX, updates)
	assert.Nil(t, err)
	assert.False(t, can)

	claim.UpdatedByID = u.ArangoID()
	can, err = claim.UserCanUpdate(CTX, updates)
	assert.Nil(t, err)
	assert.False(t, can)

	claim.CreatedByID = u.ArangoID()
	can, err = claim.UserCanUpdate(CTX, updates)
	assert.Nil(t, err)
	assert.True(t, can)

	u = User{}
	u.PrepareForCreate(CTX)
	CTX.UserContext = u
	can, err = claim.UserCanUpdate(CTX, updates)
	assert.Nil(t, err)
	assert.False(t, can)

	u.Curator = true
	CTX.UserContext = u
	can, err = claim.UserCanUpdate(CTX, updates)
	assert.Nil(t, err)
	assert.True(t, can)
}

func TestClaimConvertToMultiPremise(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "I'm just a premise, yes I'm only a premise. But I hope to be a multi-premise.",
		Description:  "ConvertToMultiPremise",
		Image:        "https://thesaurus.plus/img/synonyms/125/break_into_pieces.png",
		MultiPremise: false,
		PremiseRule:  PREMISE_RULE_NONE,
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "I don't belong on a multi-premise",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "What, you think I do?",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	context := Context{
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

	anotherClaim := Claim{
		Title: "MPConversion other claim",
	}
	err = anotherClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argToAnotherClaim := Argument{
		TargetClaimID: &anotherClaim.ID,
		ClaimID:       claim.ID,
		Pro:           false,
	}
	err = argToAnotherClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	mpClaim := Claim{
		Title:        "I'm already an MPClaim. You wish you were!",
		MultiPremise: true,
		PremiseRule:  PREMISE_RULE_ALL,
	}
	err = mpClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = mpClaim.AddPremise(CTX, &claim)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	startTime := time.Now()

	// Check connections
	premiseEdges, err := claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(premiseEdges))

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
	assert.Equal(t, argToAnotherClaim.ArangoID(), bces[0].From)
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

	ces, err := claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ces))
	assert.Equal(t, context.ArangoID(), ces[0].From)
	assert.Equal(t, claim.ArangoID(), ces[0].To)

	ctxs, err := claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ctxs))
	assert.Equal(t, context.ArangoID(), ctxs[0].ArangoID())

	// Convert to Multi-premise
	err = claim.ConvertToMultiPremise(CTX)
	assert.NoError(t, err)

	premiseEdges, err = claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premiseEdges))

	premise := Claim{}
	premise.Key = premiseEdges[0].To[7:]
	err = premise.Load(CTX)
	assert.NoError(t, err)

	inferences, err = claim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(inferences))

	arg1.Load(CTX)
	arg2.Load(CTX)
	inferences, err = premise.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(inferences))
	assert.Equal(t, premise.ArangoID(), inferences[0].From)
	assert.Equal(t, premise.ArangoID(), inferences[1].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.Equal(t, arg2.ArangoID(), inferences[1].To)

	bces, err = claim.BaseClaimEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bces))
	assert.Equal(t, argToAnotherClaim.ArangoID(), bces[0].From)
	assert.Equal(t, claim.ArangoID(), bces[0].To)

	args, err = claim.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(args))

	args, err = premise.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())
	assert.Equal(t, arg2.ArangoID(), args[1].ArangoID())

	premiseEdges, err = mpClaim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premiseEdges))
	assert.Equal(t, mpClaim.ArangoID(), premiseEdges[0].From)
	assert.Equal(t, claim.ArangoID(), premiseEdges[0].To)

	ces, err = claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(ces))

	ces, err = premise.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ces))
	assert.Equal(t, context.ArangoID(), ces[0].From)
	assert.Equal(t, premise.ArangoID(), ces[0].To)

	// This method returns Contexts of the premises!
	ctxs, err = claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ctxs))
	assert.Equal(t, context.ArangoID(), ctxs[0].ArangoID())

	// TODO: actually, both should maintain contexts... for the purpose of search?
	ctxs, err = premise.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ctxs))
	assert.Equal(t, context.ArangoID(), ctxs[0].ArangoID())

	err = premise.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(premise.ProArgs))
	assert.Equal(t, 1, len(premise.ConArgs))
	assert.Equal(t, arg1.ArangoID(), premise.ProArgs[0].ArangoID())
	assert.Equal(t, arg2.ArangoID(), premise.ConArgs[0].ArangoID())

	// Check past connections
	claim.QueryAt = &startTime
	claim.Load(CTX)
	arg1.QueryAt = &startTime
	arg1.Load(CTX)
	arg2.QueryAt = &startTime
	arg2.Load(CTX)
	mpClaim.QueryAt = &startTime
	mpClaim.Load(CTX)

	premiseEdges, err = claim.PremiseEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(premiseEdges))

	inferences, err = claim.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(inferences))
	assert.Equal(t, claim.ArangoID(), inferences[0].From)
	assert.Equal(t, claim.ArangoID(), inferences[1].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.Equal(t, arg2.ArangoID(), inferences[1].To)

	bces, err = claim.BaseClaimEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bces))
	assert.Equal(t, argToAnotherClaim.ArangoID(), bces[0].From)
	assert.Equal(t, claim.ArangoID(), bces[0].To)

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

	ces, err = claim.ContextEdges(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ces))
	assert.Equal(t, context.ArangoID(), ces[0].From)
	assert.Equal(t, claim.ArangoID(), ces[0].To)

	ctxs, err = claim.Contexts(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ctxs))
	assert.Equal(t, context.ArangoID(), ctxs[0].ArangoID())

}
