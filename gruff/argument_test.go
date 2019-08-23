package gruff

import (
	"fmt"
	"testing"
	"time"

	"github.com/GruffDebate/server/support"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestArgumentValidateForCreate(t *testing.T) {
	a := Argument{}

	assert.Equal(t, "claimId: non zero value required;", a.ValidateForCreate().Error())

	a.Title = "A"
	assert.Equal(t, "title: A does not validate as length(3|1000);", a.ValidateForCreate().Error())

	a.Title = "This is a real argument"
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForCreate().Error())

	a.Description = "D"
	assert.Equal(t, "desc: D does not validate as length(3|4000);", a.ValidateForCreate().Error())

	a.Description = "This is a real description"
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForCreate().Error())

	a.ClaimID = ""
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForCreate().Error())

	a.ClaimID = uuid.New().String()
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForCreate().Error())

	a.TargetClaimID = support.StringPtr(uuid.New().String())
	assert.NoError(t, a.ValidateForCreate())

	a.TargetClaimID = nil
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForCreate().Error())

	a.TargetClaimID = support.StringPtr(uuid.New().String())
	assert.Nil(t, a.ValidateForCreate())
}

func TestArgumentValidateForUpdate(t *testing.T) {
	a := Argument{}
	updates := map[string]interface{}{}

	assert.Equal(t, "claimId: non zero value required;", a.ValidateForUpdate(updates).Error())

	updates["title"] = "A"
	assert.Equal(t, "title: A does not validate as length(3|1000);", a.ValidateForUpdate(updates).Error())

	updates["title"] = "This is a real argument"
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForUpdate(updates).Error())

	updates["desc"] = "D"
	assert.Equal(t, "desc: D does not validate as length(3|4000);", a.ValidateForUpdate(updates).Error())

	updates["desc"] = "This is a real description"
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForUpdate(updates).Error())

	updates["claimId"] = ""
	assert.Equal(t, "claimId: non zero value required;", a.ValidateForUpdate(updates).Error())

	updates["claimId"] = uuid.New().String()
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForUpdate(updates).Error())

	updates["targetClaimId"] = uuid.New().String()
	assert.NoError(t, a.ValidateForUpdate(updates))

	updates["targetClaimId"] = nil
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForUpdate(updates).Error())
}

func TestCreateArgumentForClaimNoBase(t *testing.T) {
	setupDB()
	defer teardownDB()

	u := User{}
	u.Key = "testuser"
	CTX.UserContext = u

	claim := Claim{
		Title:       "Let's create a new claim",
		Description: "Claims in general should be true or false",
		Negation:    "Let's not...",
		Question:    "Should we create a new Claim?",
		Note:        fmt.Sprintf("A Sprintf Note just to import fmt"),
		Image:       "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	arg := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Let's create a new argument",
		Description:   "Arguments are all about connecting things",
		Negation:      "Lettuce not...",
		Question:      "Should we create a new Argument?",
		Note:          "I'm not sure that there should be notes for this",
		Pro:           true,
	}
	arg.ProArgs = []Argument{Argument{Title: "I should not be saved"}}
	arg.ConArgs = []Argument{Argument{Title: "Me neither, bro."}}
	err = arg.Create(CTX)
	assert.NoError(t, err)

	saved := Argument{}
	saved.ID = arg.ID
	err = saved.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Equal(t, u.ArangoID(), saved.CreatedByID)
	assert.Nil(t, saved.DeletedAt)
	assert.Equal(t, arg.Title, saved.Title)
	assert.Equal(t, arg.Description, saved.Description)
	assert.Equal(t, arg.Negation, saved.Negation)
	assert.Equal(t, arg.Question, saved.Question)
	assert.Equal(t, arg.Note, saved.Note)
	assert.Equal(t, arg.TargetClaimID, saved.TargetClaimID)
	assert.Equal(t, arg.Pro, saved.Pro)
	assert.Equal(t, []Argument(nil), saved.ProArgs)
	assert.Equal(t, []Argument(nil), saved.ConArgs)

	// Make sure a base claim was created
	bc := Claim{}
	bc.ID = arg.ClaimID
	err = bc.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, bc.Key)
	assert.NotEmpty(t, bc.CreatedAt)
	assert.NotEmpty(t, bc.UpdatedAt)
	assert.Equal(t, u.ArangoID(), bc.CreatedByID)
	assert.Nil(t, bc.DeletedAt)
	assert.Equal(t, arg.Title, bc.Title)
	assert.Equal(t, arg.Description, bc.Description)
	assert.Equal(t, arg.Negation, bc.Negation)
	assert.Equal(t, arg.Question, bc.Question)
	assert.Equal(t, arg.Note, bc.Note)
	assert.False(t, bc.MultiPremise)

	// Check edges
	bce, err := arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, bce.Key)
	assert.NotEmpty(t, bce.CreatedAt)
	assert.Nil(t, bce.DeletedAt)
	assert.Equal(t, u.ArangoID(), bce.CreatedByID)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, bc.ArangoID(), bce.To)

	inf, err := arg.Inference(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, inf.Key)
	assert.NotEmpty(t, inf.CreatedAt)
	assert.Equal(t, u.ArangoID(), inf.CreatedByID)
	assert.Nil(t, inf.DeletedAt)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)
}

func TestCreateArgumentForClaimWithBase(t *testing.T) {
	setupDB()
	defer teardownDB()

	u := User{}
	u.Key = "testuser"
	CTX.UserContext = u

	claim := Claim{
		Title:       "Let's create a new claim",
		Description: "Claims in general should be true or false",
		Negation:    "Let's not...",
		Question:    "Should we create a new Claim?",
		Note:        "He who notes is a note taker",
		Image:       "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	baseClaim := Claim{
		Title:       "Bass!",
		Question:    "How low can you go?",
		Negation:    "Death row",
		Description: "What a brotha knows",
		Image:       "http://straightfromthea.com/wp-content/uploads/2017/08/Flavour-and-Chuck-D-520x397.jpg",
	}
	err = baseClaim.Create(CTX)
	assert.NoError(t, err)

	arg := Argument{
		TargetClaimID: &claim.ID,
		ClaimID:       baseClaim.ID,
		Title:         "Let's create a new argument",
		Description:   "Arguments are all about connecting things",
		Negation:      "Lettuce not...",
		Question:      "Should we create a new Argument?",
		Note:          "I'm not sure that there should be notes for this",
		Pro:           true,
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)

	saved := Argument{}
	saved.ID = arg.ID
	err = saved.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Equal(t, u.ArangoID(), saved.CreatedByID)
	assert.Nil(t, saved.DeletedAt)
	assert.Equal(t, arg.Title, saved.Title)
	assert.Equal(t, arg.Description, saved.Description)
	assert.Equal(t, arg.Negation, saved.Negation)
	assert.Equal(t, arg.Question, saved.Question)
	assert.Equal(t, arg.Note, saved.Note)
	assert.Equal(t, arg.TargetClaimID, saved.TargetClaimID)
	assert.Equal(t, arg.ClaimID, saved.ClaimID)
	assert.Equal(t, arg.Pro, saved.Pro)

	// Check edges
	bce, err := arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, bce.Key)
	assert.NotEmpty(t, bce.CreatedAt)
	assert.Equal(t, u.ArangoID(), bce.CreatedByID)
	assert.Nil(t, bce.DeletedAt)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err := arg.Inference(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, inf.Key)
	assert.NotEmpty(t, inf.CreatedAt)
	assert.Equal(t, u.ArangoID(), inf.CreatedByID)
	assert.Nil(t, inf.DeletedAt)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)
}

func TestCreateArgumentForArgument(t *testing.T) {
	setupDB()
	defer teardownDB()

	u := User{}
	u.Key = "testuser"
	CTX.UserContext = u

	claim := Claim{
		Title:       "Let's create a new claim",
		Description: "Claims in general should be true or false",
		Negation:    "Let's not...",
		Question:    "Should we create a new Claim?",
		Note:        "He who notes is a note taker",
		Image:       "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	targarg := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Daenerys",
		Description:   "Queen of dragons",
		Negation:      "Not Queen",
		Question:      "Will she be queen?",
		Note:          "Dracarys",
		Pro:           true,
	}
	err = targarg.Create(CTX)
	assert.NoError(t, err)

	arg := Argument{
		TargetArgumentID: &targarg.ID,
		Title:            "Let's create a new argument",
		Description:      "Arguments are all about connecting things",
		Negation:         "Lettuce not...",
		Question:         "Should we create a new Argument?",
		Note:             "I'm not sure that there should be notes for this",
		Pro:              true,
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)

	saved := Argument{}
	saved.ID = arg.ID
	err = saved.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.Key)
	assert.NotEmpty(t, saved.ID)
	assert.NotEmpty(t, saved.CreatedAt)
	assert.NotEmpty(t, saved.UpdatedAt)
	assert.Equal(t, u.ArangoID(), saved.CreatedByID)
	assert.Nil(t, saved.DeletedAt)
	assert.Equal(t, arg.Title, saved.Title)
	assert.Equal(t, arg.Description, saved.Description)
	assert.Equal(t, arg.Negation, saved.Negation)
	assert.Equal(t, arg.Question, saved.Question)
	assert.Equal(t, arg.Note, saved.Note)
	assert.Equal(t, arg.TargetClaimID, saved.TargetClaimID)
	assert.Equal(t, arg.ClaimID, saved.ClaimID)
	assert.Equal(t, arg.Pro, saved.Pro)

	// Make sure a base claim was created
	bc := Claim{}
	bc.ID = arg.ClaimID
	err = bc.Load(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, bc.Key)
	assert.NotEmpty(t, bc.CreatedAt)
	assert.NotEmpty(t, bc.UpdatedAt)
	assert.Equal(t, u.ArangoID(), bc.CreatedByID)
	assert.Nil(t, bc.DeletedAt)
	assert.Equal(t, arg.Title, bc.Title)
	assert.Equal(t, arg.Description, bc.Description)
	assert.Equal(t, arg.Negation, bc.Negation)
	assert.Equal(t, arg.Question, bc.Question)
	assert.Equal(t, arg.Note, bc.Note)
	assert.False(t, bc.MultiPremise)

	// Check edges
	bce, err := arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, bce.Key)
	assert.NotEmpty(t, bce.CreatedAt)
	assert.Nil(t, bce.DeletedAt)
	assert.Equal(t, u.ArangoID(), bce.CreatedByID)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, bc.ArangoID(), bce.To)

	inf, err := arg.Inference(CTX)
	assert.NoError(t, err)
	assert.NotEmpty(t, inf.Key)
	assert.NotEmpty(t, inf.CreatedAt)
	assert.Nil(t, inf.DeletedAt)
	assert.Equal(t, u.ArangoID(), inf.CreatedByID)
	assert.Equal(t, targarg.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)
}

func TestLoadArgumentAtDate(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:       "Let's create a new claim",
		Description: "Claims in general should be true or false",
		Negation:    "Let's not...",
		Question:    "Should we create a new Claim?",
		Note:        "He who notes is a note taker",
		Image:       "https://slideplayer.com/slide/4862164/15/images/9/7.3+Creating+Claims+7-9.+The+Create+Claims+button+in+the+Claim+Management+dialog+box+opens+the+Create+Claims+dialog+box..jpg",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)

	arg := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Let's create a new argument",
		Description:   "Arguments are all about connecting things",
		Negation:      "Lettuce not...",
		Question:      "Should we create a new Argument?",
		Note:          "I'm not sure that there should be notes for this",
		Pro:           true,
	}

	err = arg.Create(CTX)
	assert.NoError(t, err)
	patch := Updates{
		"start": time.Now().Add(-25 * time.Hour),
		"end":   time.Now().Add(-24 * time.Hour),
	}
	col, _ := CTX.Arango.CollectionFor(&arg)
	col.UpdateDocument(CTX.Context, arg.ArangoKey(), patch)

	firstKey := arg.ArangoKey()

	err = arg.Create(CTX)
	assert.NoError(t, err)
	patch["start"] = time.Now().Add(-24 * time.Hour)
	patch["end"] = time.Now().Add(-1 * time.Hour)
	col.UpdateDocument(CTX.Context, arg.ArangoKey(), patch)

	secondKey := arg.ArangoKey()

	err = arg.Create(CTX)
	assert.NoError(t, err)
	patch["start"] = time.Now().Add(-1 * time.Hour)
	delete(patch, "end")
	col.UpdateDocument(CTX.Context, arg.ArangoKey(), patch)

	thirdKey := arg.ArangoKey()

	lookup := Argument{}
	lookup.ID = arg.ID
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())

	lookup.QueryAt = support.TimePtr(time.Now().Add(-1 * time.Minute))
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())
	thirdCreatedAt := lookup.CreatedAt

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.QueryAt = support.TimePtr(time.Now().Add(-2 * time.Hour))
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, secondKey, lookup.ArangoKey())
	secondCreatedAt := lookup.CreatedAt

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.QueryAt = support.TimePtr(time.Now().Add(-25 * time.Hour))
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, firstKey, lookup.ArangoKey())
	firstCreatedAt := lookup.CreatedAt

	// TODO: Throw a NotFoundError?
	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.QueryAt = support.TimePtr(time.Now().Add(-48 * time.Hour))
	err = lookup.Load(CTX)
	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.QueryAt = &firstCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, firstKey, lookup.ArangoKey())

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.QueryAt = &secondCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, secondKey, lookup.ArangoKey())

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.QueryAt = &thirdCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())
}

func TestArgumentLoadFull(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title:        "This is the Argument LoadAll test claim",
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
		Title:       "So very far away from Argument LoadAll",
		Description: "So distant, you cannot see me.",
	}
	err = distantClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "ARG?",
		Pro:           true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Load ARG!",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg3 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "Do it ARG!",
		Pro:           true,
	}
	err = arg3.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	theArg := Argument{
		TargetClaimID: &distantClaim.ID,
		ClaimID:       claim.ID,
		Title:         "This is the Argument for LoadFull",
	}
	err = theArg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1arg := Argument{
		TargetArgumentID: &arg1.ID,
		Title:            "Do it...ARRRRRG!",
		Pro:              false,
	}
	err = arg1arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	theArgArg1 := Argument{
		TargetArgumentID: &theArg.ID,
		Title:            "Now we're just getting ridiculous... an argument for THE ARG?",
		Pro:              true,
	}
	err = theArgArg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	theArgArg2 := Argument{
		TargetArgumentID: &theArg.ID,
		Title:            "Now we're just getting ridiculous... another argument for THE ARG?",
		Pro:              false,
	}
	err = theArgArg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	theArgArg3 := Argument{
		TargetArgumentID: &theArg.ID,
		Title:            "Now we're just getting ridiculous... another nother argument for THE ARG?",
		Pro:              false,
	}
	err = theArgArg3.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	theArgArg3.Load(CTX)

	theArgArg3Base := Claim{}
	theArgArg3Base.ID = theArgArg3.ClaimID
	theArgArg3Base.Load(CTX)

	theArgArg3BaseArg := Argument{
		TargetClaimID: &theArgArg3Base.ID,
		Title:         "You think that's bad? How about one more for the base arg of the arg arg 3?",
		Pro:           true,
	}
	err = theArgArg3BaseArg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Simple Load
	err = theArg.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, theArg.DeletedAt)
	assert.Equal(t, "This is the Argument for LoadFull", theArg.Title)
	assert.Nil(t, theArg.Claim)
	assert.Equal(t, 0, len(theArg.ProArgs))
	assert.Equal(t, 0, len(theArg.ConArgs))

	// Load All
	arg1.Load(CTX)
	arg2.Load(CTX)
	arg3.Load(CTX)
	theArgArg1.Load(CTX)
	theArgArg2.Load(CTX)
	theArgArg3.Load(CTX)
	var carg1, carg2, carg3, ctaa1, ctaa2, ctaa3 Claim
	carg1.ID = arg1.ClaimID
	carg2.ID = arg2.ClaimID
	carg3.ID = arg3.ClaimID
	ctaa1.ID = theArgArg1.ClaimID
	ctaa2.ID = theArgArg2.ClaimID
	ctaa3.ID = theArgArg3.ClaimID
	carg1.Load(CTX)
	carg2.Load(CTX)
	carg3.Load(CTX)
	ctaa1.Load(CTX)
	ctaa2.Load(CTX)
	ctaa3.Load(CTX)
	arg1.Claim = &carg1
	arg2.Claim = &carg2
	arg3.Claim = &carg3
	theArgArg1.Claim = &ctaa1
	theArgArg2.Claim = &ctaa2
	theArgArg3.Claim = &ctaa3

	claim.Load(CTX)
	claim.ProArgs = []Argument{arg1, arg3}
	claim.ConArgs = []Argument{arg2}

	err = theArg.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, theArg.DeletedAt)
	assert.Equal(t, "This is the Argument for LoadFull", theArg.Title)
	assert.Equal(t, claim, *theArg.Claim)
	assert.Equal(t, 1, len(theArg.ProArgs))
	assert.Equal(t, 2, len(theArg.ConArgs))
	assert.Equal(t, theArgArg1, theArg.ProArgs[0])
	assert.Equal(t, theArgArg2, theArg.ConArgs[0])
	assert.Equal(t, theArgArg3, theArg.ConArgs[1])

	err = arg1.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	err = theArgArg2.Delete(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	arg1.Load(CTX)
	theArgArg2.Load(CTX)

	claim.ProArgs = []Argument{arg3}
	claim.ConArgs = []Argument{arg2}

	err = theArg.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, theArg.DeletedAt)
	assert.Equal(t, "This is the Argument for LoadFull", theArg.Title)
	assert.Equal(t, claim, *theArg.Claim)
	assert.Equal(t, 1, len(theArg.ProArgs))
	assert.Equal(t, 1, len(theArg.ConArgs))
	assert.Equal(t, theArgArg1, theArg.ProArgs[0])
	assert.Equal(t, theArgArg3, theArg.ConArgs[0])

	// Previous points in time
	arg1.Claim = &carg1
	arg2.Claim = &carg2
	arg3.Claim = &carg3
	theArgArg1.Claim = &ctaa1
	theArgArg2.Claim = &ctaa2
	theArgArg3.Claim = &ctaa3
	claim.ProArgs = []Argument{arg1, arg3}
	claim.ConArgs = []Argument{arg2}
	theArg.QueryAt = &theArgArg3BaseArg.CreatedAt
	err = theArg.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, theArg.DeletedAt)
	assert.Equal(t, "This is the Argument for LoadFull", theArg.Title)
	assert.Equal(t, claim, *theArg.Claim)
	assert.Equal(t, 1, len(theArg.ProArgs))
	assert.Equal(t, 2, len(theArg.ConArgs))
	assert.Equal(t, theArgArg1, theArg.ProArgs[0])
	assert.Equal(t, theArgArg2, theArg.ConArgs[0])
	assert.Equal(t, theArgArg3, theArg.ConArgs[1])

	claim.ProArgs = []Argument{arg1}
	claim.ConArgs = nil
	theArg.QueryAt = &arg1.CreatedAt
	err = theArg.LoadFull(CTX)
	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())

	claim.ProArgs = []Argument{arg1, arg3}
	claim.ConArgs = []Argument{arg2}
	theArg.QueryAt = &theArg.CreatedAt
	err = theArg.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, theArg.DeletedAt)
	assert.Equal(t, "This is the Argument for LoadFull", theArg.Title)
	assert.Equal(t, claim, *theArg.Claim)
	assert.Equal(t, 0, len(theArg.ProArgs))
	assert.Equal(t, 0, len(theArg.ConArgs))

	claim.ProArgs = []Argument{arg3}
	theArg.QueryAt = nil
	err = theArg.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Nil(t, theArg.DeletedAt)
	assert.Equal(t, "This is the Argument for LoadFull", theArg.Title)
	assert.Equal(t, claim, *theArg.Claim)
	assert.Equal(t, 1, len(theArg.ProArgs))
	assert.Equal(t, 1, len(theArg.ConArgs))
	assert.Equal(t, theArgArg1, theArg.ProArgs[0])
	assert.Equal(t, theArgArg3, theArg.ConArgs[0])
}

func TestArgumentUpdate(t *testing.T) {
	setupDB()
	defer teardownDB()

	claim := Claim{
		Title: "You are going to update the Argument for this Claim, aren't you?",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	baseClaim := Claim{
		Title: "How low can you go? Update... what a Claim knows.",
	}
	err = baseClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg := Argument{
		TargetClaimID: &claim.ID,
		ClaimID:       baseClaim.ID,
		Title:         "This is the argument everyone wants to Update",
		Pro:           true,
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetArgumentID: &arg.ID,
		Title:            "I'm an argument to an updated argument",
		Pro:              true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := Argument{
		TargetArgumentID: &arg.ID,
		Title:            "No, I'm an argument to an updated argument",
		Pro:              false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Next check edges, then version and recheck everything
	inferences, err := arg.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(inferences))
	assert.Equal(t, arg.ArangoID(), inferences[0].From)
	assert.Equal(t, arg.ArangoID(), inferences[1].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.Equal(t, arg2.ArangoID(), inferences[1].To)

	bce, err := arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err := arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	args, err := arg.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())
	assert.Equal(t, arg2.ArangoID(), args[1].ArangoID())

	// Update the argument
	curator := User{Username: "another_curator", Curator: true}
	err = curator.Create(CTX)
	assert.NoError(t, err)
	CTX.UserContext = curator

	err = arg.Load(CTX)
	assert.NoError(t, err)
	origKey := arg.ArangoKey()

	updates := map[string]interface{}{
		"title": "Now this argument has gone and been updated",
		"desc":  "And so has its description",
	}
	err = arg.Update(CTX, updates)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, arg.DeletedAt)
	assert.Equal(t, "Now this argument has gone and been updated", arg.Title)
	assert.Equal(t, "And so has its description", arg.Description)
	assert.NotEqual(t, origKey, arg.ArangoKey())
	assert.Equal(t, DEFAULT_USER.ArangoID(), arg.CreatedByID)
	assert.Equal(t, CTX.UserContext.ArangoID(), arg.UpdatedByID)

	origArg := Argument{}
	origArg.Key = origKey
	err = origArg.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, origArg.DeletedAt)
	assert.Equal(t, "This is the argument everyone wants to Update", origArg.Title)
	assert.Equal(t, "", origArg.Description)
	assert.Equal(t, origKey, origArg.ArangoKey())

	// Verify new edges were created
	inferences, err = arg.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(inferences))
	assert.Equal(t, arg.ArangoID(), inferences[0].From)
	assert.Equal(t, arg.ArangoID(), inferences[1].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.Equal(t, arg2.ArangoID(), inferences[1].To)
	assert.Nil(t, inferences[0].DeletedAt)
	assert.Nil(t, inferences[1].DeletedAt)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)
	assert.Nil(t, bce.DeletedAt)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)
	assert.Nil(t, inf.DeletedAt)

	args, err = arg.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())
	assert.Equal(t, arg2.ArangoID(), args[1].ArangoID())

	// Verify that the old edges were deleted
	inferences, err = origArg.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(inferences))
	assert.Equal(t, origArg.ArangoID(), inferences[0].From)
	assert.Equal(t, origArg.ArangoID(), inferences[1].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)
	assert.Equal(t, arg2.ArangoID(), inferences[1].To)
	assert.NotNil(t, inferences[0].DeletedAt)
	assert.NotNil(t, inferences[1].DeletedAt)

	bce, err = origArg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, origArg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)
	assert.NotNil(t, bce.DeletedAt)

	inf, err = origArg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, origArg.ArangoID(), inf.To)
	assert.NotNil(t, inf.DeletedAt)

	args, err = origArg.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())
	assert.Equal(t, arg2.ArangoID(), args[1].ArangoID())
}

func TestArgumentMoveTo(t *testing.T) {
	setupDB()
	defer teardownDB()

	// TODO: can't move to MP Claim

	claim := Claim{
		Title: "This Claim is here for testing MoveTo for Arguments (what else?)",
	}
	err := claim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	baseClaim := Claim{
		Title: "How low can you go? MoveTo... what a Claim knows.",
	}
	err = baseClaim.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg := Argument{
		TargetClaimID: &claim.ID,
		ClaimID:       baseClaim.ID,
		Title:         "This is the argument everyone wants to Move",
		Pro:           true,
	}
	err = arg.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg1 := Argument{
		TargetArgumentID: &arg.ID,
		Title:            "I'm an argument to an argument that gets moved",
		Pro:              true,
	}
	err = arg1.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	arg2 := Argument{
		TargetClaimID: &claim.ID,
		Title:         "I'm an argument that you can move to",
		Pro:           false,
	}
	err = arg2.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	argBC := Argument{
		TargetClaimID: &baseClaim.ID,
		Title:         "Move the Argument to me and you'll be sorry!",
		Pro:           true,
	}
	err = argBC.Create(CTX)
	assert.NoError(t, err)
	CTX.RequestAt = nil

	// Next check starting edges
	inferences, err := arg.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inferences))
	assert.Equal(t, arg.ArangoID(), inferences[0].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)

	bce, err := arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err := arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	args, err := arg.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, arg.ArangoID(), claim.ProArgs[0].ArangoID())
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg2.ArangoID(), claim.ConArgs[0].ArangoID())

	// Move it to the same claim, no change expected
	arangoKey := arg.Key
	err = arg.MoveTo(CTX, &claim, arg.Pro)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arangoKey, arg.Key)
	assert.True(t, arg.Pro)

	inferences, err = arg.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inferences))
	assert.Equal(t, arg.ArangoID(), inferences[0].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	args, err = arg.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, arg.ArangoID(), claim.ProArgs[0].ArangoID())
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg2.ArangoID(), claim.ConArgs[0].ArangoID())

	// Move to the same claim, but change the polarity
	arangoKey = arg.Key
	err = arg.MoveTo(CTX, &claim, false)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	arg.DeletedAt = nil
	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.NotEqual(t, arangoKey, arg.Key)
	assert.False(t, arg.Pro)

	inferences, err = arg.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inferences))
	assert.Equal(t, arg.ArangoID(), inferences[0].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	args, err = arg.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 2, len(claim.ConArgs))
	assert.Equal(t, arg2.ArangoID(), claim.ConArgs[0].ArangoID())
	assert.Equal(t, arg.ArangoID(), claim.ConArgs[1].ArangoID())

	// Move to another claim
	c2 := Claim{}
	c2.ID = arg2.ClaimID
	err = c2.Load(CTX)
	assert.NoError(t, err)

	arangoKey = arg.Key
	err = arg.MoveTo(CTX, &c2, true)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	arg.DeletedAt = nil
	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.NotEqual(t, arangoKey, arg.Key)
	assert.True(t, arg.Pro)
	assert.Equal(t, c2.ID, *arg.TargetClaimID)

	c2Time := time.Now()

	inferences, err = arg.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inferences))
	assert.Equal(t, arg.ArangoID(), inferences[0].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, c2.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	args, err = arg.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg2.ArangoID(), claim.ConArgs[0].ArangoID())

	err = c2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c2.ProArgs))
	assert.Equal(t, 0, len(c2.ConArgs))
	assert.Equal(t, arg.ArangoID(), c2.ProArgs[0].ArangoID())

	// Move to an argument
	arangoKey = arg.Key
	err = arg.MoveTo(CTX, &arg2, false)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	arg.DeletedAt = nil
	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.NotEqual(t, arangoKey, arg.Key)
	assert.False(t, arg.Pro)
	assert.Nil(t, arg.TargetClaimID)
	assert.Equal(t, arg2.ID, *arg.TargetArgumentID)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg2.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	err = c2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(c2.ProArgs))
	assert.Equal(t, 0, len(c2.ConArgs))

	err = arg2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(arg2.ProArgs))
	assert.Equal(t, 1, len(arg2.ConArgs))
	assert.Equal(t, arg.ArangoID(), arg2.ConArgs[0].ArangoID())

	// Swap polarity on the argument
	arangoKey = arg.Key
	err = arg.MoveTo(CTX, &arg2, true)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	arg.DeletedAt = nil
	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.NotEqual(t, arangoKey, arg.Key)
	assert.True(t, arg.Pro)
	assert.Equal(t, arg2.ID, *arg.TargetArgumentID)

	arg2SwapTime := time.Now()

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg2.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	err = arg2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(arg2.ProArgs))
	assert.Equal(t, 0, len(arg2.ConArgs))
	assert.Equal(t, arg.ArangoID(), arg2.ProArgs[0].ArangoID())

	// Try to move it to its own base claim
	arangoKey = arg.Key
	err = arg.MoveTo(CTX, &baseClaim, true)
	assert.Error(t, err)
	assert.Equal(t, "An argument cannot be moved to its own base claim", err.Error())
	CTX.RequestAt = nil
	arg.DeletedAt = nil
	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arangoKey, arg.Key)
	assert.True(t, arg.Pro)
	assert.Equal(t, arg2.ID, *arg.TargetArgumentID)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg2.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	err = arg2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(arg2.ProArgs))
	assert.Equal(t, 0, len(arg2.ConArgs))
	assert.Equal(t, arg.ArangoID(), arg2.ProArgs[0].ArangoID())

	// Try to move it to one of its own arguments
	arangoKey = arg.Key
	err = arg.MoveTo(CTX, &arg1, false)
	assert.Error(t, err)
	assert.Equal(t, "An argument cannot be moved to one of its own arguments", err.Error())
	CTX.RequestAt = nil
	arg.DeletedAt = nil
	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arangoKey, arg.Key)
	assert.True(t, arg.Pro)
	assert.Equal(t, arg2.ID, *arg.TargetArgumentID)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg2.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	err = arg2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(arg2.ProArgs))
	assert.Equal(t, 0, len(arg2.ConArgs))
	assert.Equal(t, arg.ArangoID(), arg2.ProArgs[0].ArangoID())

	// Move it back to its original claim
	arangoKey = arg.Key
	err = arg.MoveTo(CTX, &claim, true)
	assert.NoError(t, err)
	CTX.RequestAt = nil
	err = arg.Load(CTX)
	assert.NoError(t, err)
	assert.NotEqual(t, arangoKey, arg.Key)
	assert.True(t, arg.Pro)
	assert.Equal(t, claim.ID, *arg.TargetClaimID)
	assert.Nil(t, arg.TargetArgumentID)

	inferences, err = arg.Inferences(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inferences))
	assert.Equal(t, arg.ArangoID(), inferences[0].From)
	assert.Equal(t, arg1.ArangoID(), inferences[0].To)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, claim.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	args, err = arg.Arguments(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, arg1.ArangoID(), args[0].ArangoID())

	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(claim.ProArgs))
	assert.Equal(t, arg.ArangoID(), claim.ProArgs[0].ArangoID())
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg2.ArangoID(), claim.ConArgs[0].ArangoID())

	err = arg2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(arg2.ProArgs))
	assert.Equal(t, 0, len(arg2.ConArgs))

	// Load back to after swapping polarity on the argument
	arg.QueryAt = &arg2SwapTime
	err = arg.LoadFull(CTX)
	assert.NoError(t, err)
	assert.True(t, arg.Pro)
	assert.Nil(t, arg.TargetClaimID)
	assert.Equal(t, arg2.ID, *arg.TargetArgumentID)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg2.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	arg2.QueryAt = &arg2SwapTime
	err = arg2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(arg2.ProArgs))
	assert.Equal(t, 0, len(arg2.ConArgs))
	assert.Equal(t, arg.ArangoID(), arg2.ProArgs[0].ArangoID())

	c2.QueryAt = &arg2SwapTime
	err = c2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(c2.ProArgs))
	assert.Equal(t, 0, len(c2.ConArgs))

	claim.QueryAt = &arg2SwapTime
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg2.ArangoID(), claim.ConArgs[0].ArangoID())

	// Load back to after moving to c2
	arg.QueryAt = &c2Time
	err = arg.LoadFull(CTX)
	assert.NoError(t, err)
	assert.True(t, arg.Pro)
	assert.Equal(t, c2.ID, *arg.TargetClaimID)
	assert.Nil(t, arg.TargetArgumentID)

	bce, err = arg.BaseClaimEdge(CTX)
	assert.NoError(t, err)
	assert.Equal(t, arg.ArangoID(), bce.From)
	assert.Equal(t, baseClaim.ArangoID(), bce.To)

	inf, err = arg.Inference(CTX)
	assert.NoError(t, err)
	assert.Equal(t, c2.ArangoID(), inf.From)
	assert.Equal(t, arg.ArangoID(), inf.To)

	arg2.QueryAt = &c2Time
	err = arg2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(arg2.ProArgs))
	assert.Equal(t, 0, len(arg2.ConArgs))

	c2.QueryAt = &c2Time
	err = c2.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c2.ProArgs))
	assert.Equal(t, 0, len(c2.ConArgs))
	assert.Equal(t, arg.ArangoID(), c2.ProArgs[0].ArangoID())

	claim.QueryAt = &c2Time
	err = claim.LoadFull(CTX)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(claim.ProArgs))
	assert.Equal(t, 1, len(claim.ConArgs))
	assert.Equal(t, arg2.ArangoID(), claim.ConArgs[0].ArangoID())
}
