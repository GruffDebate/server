package gruff

import (
	"testing"
	"time"

	"github.com/GruffDebate/server/support"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestArgumentValidateForCreate(t *testing.T) {
	a := Argument{}

	assert.Equal(t, "title: non zero value required;", a.ValidateForCreate().Error())

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

	assert.Equal(t, "title: non zero value required;", a.ValidateForUpdate(updates).Error())

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
	arg.DeletedAt = support.TimePtr(time.Now().Add(-24 * time.Hour))

	err = arg.Create(CTX)
	assert.NoError(t, err)
	patch := map[string]interface{}{"start": time.Now().Add(-25 * time.Hour)}
	col, _ := CTX.Arango.CollectionFor(arg)
	col.UpdateDocument(CTX.Context, arg.ArangoKey(), patch)

	firstKey := arg.ArangoKey()

	arg.DeletedAt = support.TimePtr(time.Now().Add(-1 * time.Hour))
	err = arg.Create(CTX)
	assert.NoError(t, err)
	patch["start"] = time.Now().Add(-24 * time.Hour)
	col.UpdateDocument(CTX.Context, arg.ArangoKey(), patch)

	secondKey := arg.ArangoKey()

	arg.DeletedAt = nil
	err = arg.Create(CTX)
	assert.NoError(t, err)
	patch["start"] = time.Now().Add(-1 * time.Hour)
	col.UpdateDocument(CTX.Context, arg.ArangoKey(), patch)

	thirdKey := arg.ArangoKey()

	lookup := Argument{}
	lookup.ID = arg.ID
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())

	lookup.CreatedAt = time.Now().Add(-1 * time.Minute)
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())
	thirdCreatedAt := lookup.CreatedAt

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.CreatedAt = time.Now().Add(-2 * time.Hour)
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, secondKey, lookup.ArangoKey())
	secondCreatedAt := lookup.CreatedAt

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.CreatedAt = time.Now().Add(-25 * time.Hour)
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, firstKey, lookup.ArangoKey())
	firstCreatedAt := lookup.CreatedAt

	// TODO: Throw a NotFoundError?
	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.CreatedAt = time.Now().Add(-48 * time.Hour)
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Equal(t, "", lookup.ArangoKey())

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.CreatedAt = firstCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, firstKey, lookup.ArangoKey())

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.CreatedAt = secondCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.NotNil(t, lookup.DeletedAt)
	assert.Equal(t, secondKey, lookup.ArangoKey())

	lookup = Argument{}
	lookup.ID = arg.ID
	lookup.CreatedAt = thirdCreatedAt
	err = lookup.Load(CTX)
	assert.NoError(t, err)
	assert.Nil(t, lookup.DeletedAt)
	assert.Equal(t, thirdKey, lookup.ArangoKey())
}
