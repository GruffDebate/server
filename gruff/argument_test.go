package gruff

import (
	"testing"

	"github.com/GruffDebate/server/support"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestArgumentValidateForCreate(t *testing.T) {
	a := Argument{}

	// assert.Equal(t, "Title: non zero value required;", a.ValidateForCreate().Error())

	// a.Title = "A"
	// assert.Equal(t, "Title: A does not validate as length(3|1000);", a.ValidateForCreate().Error())

	// a.Title = "This is a real argument"
	// assert.Equal(t, "ClaimID: non zero value required;", a.ValidateForCreate().Error())

	// a.Description = "D"
	// assert.Equal(t, "Description: D does not validate as length(3|4000);", a.ValidateForCreate().Error())

	a.Description = "This is a real description"
	assert.Equal(t, "ClaimID: non zero value required;", a.ValidateForCreate().Error())

	a.ClaimID = uuid.Nil
	assert.Equal(t, "ClaimID: non zero value required;", a.ValidateForCreate().Error())

	a.ClaimID = uuid.New()
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForCreate().Error())

	a.TargetClaimID = &NullableUUID{UUID: uuid.New()}
	assert.Equal(t, "Type: invalid;", a.ValidateForCreate().Error())

	a.Type = 7
	assert.Equal(t, "Type: invalid;", a.ValidateForCreate().Error())

	a.Type = ARGUMENT_FOR
	assert.Nil(t, a.ValidateForCreate())

	a.Type = ARGUMENT_AGAINST
	assert.Nil(t, a.ValidateForCreate())

	a.TargetClaimID = nil
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForCreate().Error())

	a.TargetClaimID = &NullableUUID{UUID: uuid.New()}
	assert.Nil(t, a.ValidateForCreate())
}

func TestArgumentValidateForUpdate(t *testing.T) {
	a := Argument{}

	// assert.Equal(t, "Title: non zero value required;", a.ValidateForUpdate().Error())

	// a.Title = "A"
	// assert.Equal(t, "Title: A does not validate as length(3|1000);", a.ValidateForUpdate().Error())

	// a.Title = "This is a real argument"
	// assert.Equal(t, "ClaimID: non zero value required;", a.ValidateForUpdate().Error())

	// a.Description = "D"
	// assert.Equal(t, "Description: D does not validate as length(3|4000);", a.ValidateForUpdate().Error())

	a.Description = "This is a real description"
	assert.Equal(t, "ClaimID: non zero value required;", a.ValidateForUpdate().Error())

	a.ClaimID = uuid.Nil
	assert.Equal(t, "ClaimID: non zero value required;", a.ValidateForUpdate().Error())

	a.ClaimID = uuid.New()
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForUpdate().Error())

	a.TargetClaimID = &NullableUUID{UUID: uuid.New()}
	assert.Equal(t, "Type: invalid;", a.ValidateForUpdate().Error())

	a.Type = 7
	assert.Equal(t, "Type: invalid;", a.ValidateForUpdate().Error())

	a.Type = ARGUMENT_FOR
	assert.Nil(t, a.ValidateForUpdate())

	a.Type = ARGUMENT_AGAINST
	assert.Nil(t, a.ValidateForUpdate())

	a.TargetClaimID = nil
	assert.Equal(t, "An Argument must have a target Claim or target Argument ID", a.ValidateForUpdate().Error())

	a.TargetArgumentID = &NullableUUID{UUID: uuid.New()}
	assert.Nil(t, a.ValidateForUpdate())
}

func TestOrderByBestArgument(t *testing.T) {
	setupDB()
	defer teardownDB()

	c1 := Claim{Title: "Claim 1", Truth: 0.5}
	c2 := Claim{Title: "Claim 2", Truth: 0.4}
	c3 := Claim{Title: "Claim 3", Truth: 0.6}
	c4 := Claim{Title: "Claim 4", Truth: 0.3}
	c5 := Claim{Title: "Claim 5", Truth: 0.7}
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)
	TESTDB.Create(&c3)
	TESTDB.Create(&c4)
	TESTDB.Create(&c5)

	a1 := Argument{Title: "Argument 1", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Strength: 0.4}
	a2 := Argument{Title: "Argument 1", TargetClaimID: NUUID(c1.ID), ClaimID: c3.ID, Strength: 0.2}
	a3 := Argument{Title: "Argument 1", TargetClaimID: NUUID(c1.ID), ClaimID: c4.ID, Strength: 0.6}
	a4 := Argument{Title: "Argument 1", TargetClaimID: NUUID(c1.ID), ClaimID: c5.ID, Strength: 0.7}
	TESTDB.Create(&a1)
	TESTDB.Create(&a2)
	TESTDB.Create(&a3)
	TESTDB.Create(&a4)

	args := []Argument{}
	TESTDB.Preload("Claim").Scopes(OrderByBestArgument).Find(&args)
	assert.Equal(t, 4, len(args))
	assert.Equal(t, a4.ID, args[0].ID)
	assert.Equal(t, a3.ID, args[1].ID)
	assert.Equal(t, a1.ID, args[2].ID)
	assert.Equal(t, a2.ID, args[3].ID)
}

func TestUpdateStrength(t *testing.T) {
	setupDB()
	defer teardownDB()

	c1 := Claim{Title: "C1"}
	c2 := Claim{Title: "C2"}
	c3 := Claim{Title: "C3"}
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)
	TESTDB.Create(&c3)

	a1 := Argument{Title: "A1", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID}
	a2 := Argument{Title: "Heinz 57", TargetClaimID: NUUID(c1.ID), ClaimID: c3.ID}
	TESTDB.Create(&a1)
	TESTDB.Create(&a2)

	a1.UpdateStrength(&CTX)
	a2.UpdateStrength(&CTX)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	assert.Equal(t, 0.0, a1.Strength)
	assert.Equal(t, 0.0, a2.Strength)

	ao1 := ArgumentOpinion{UserID: 1, ArgumentID: a1.ID, Strength: 0.5}
	TESTDB.Create(&ao1)

	a1.UpdateStrength(&CTX)
	a2.UpdateStrength(&CTX)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	assert.Equal(t, 0.5, a1.Strength)
	assert.Equal(t, 0.0, a2.Strength)

	ao2 := ArgumentOpinion{UserID: 2, ArgumentID: a2.ID, Strength: 0.9}
	TESTDB.Create(&ao2)

	a1.UpdateStrength(&CTX)
	a2.UpdateStrength(&CTX)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	assert.Equal(t, 0.5, a1.Strength)
	assert.Equal(t, 0.9, a2.Strength)

	ao3 := ArgumentOpinion{UserID: 3, ArgumentID: a1.ID, Strength: 0.7}
	TESTDB.Create(&ao3)

	a1.UpdateStrength(&CTX)
	a2.UpdateStrength(&CTX)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	assert.Equal(t, 0.6, a1.Strength)
	assert.Equal(t, 0.9, a2.Strength)

	ao4 := ArgumentOpinion{UserID: 4, ArgumentID: a1.ID, Strength: 0.3}
	ao5 := ArgumentOpinion{UserID: 5, ArgumentID: a2.ID, Strength: 0.6}
	ao6 := ArgumentOpinion{UserID: 6, ArgumentID: a2.ID, Strength: 0.2}
	ao7 := ArgumentOpinion{UserID: 7, ArgumentID: a2.ID, Strength: 0.8}
	ao8 := ArgumentOpinion{UserID: 8, ArgumentID: a2.ID, Strength: 0.8}
	TESTDB.Create(&ao4)
	TESTDB.Create(&ao5)
	TESTDB.Create(&ao6)
	TESTDB.Create(&ao7)
	TESTDB.Create(&ao8)

	a1.UpdateStrength(&CTX)
	a2.UpdateStrength(&CTX)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	assert.Equal(t, 0.5, a1.Strength)
	assert.Equal(t, 0.66, a2.Strength)

	ao7.Strength = 0.5
	TESTDB.Save(&ao7)

	a1.UpdateStrength(&CTX)
	a2.UpdateStrength(&CTX)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	assert.Equal(t, 0.5, a1.Strength)
	assert.Equal(t, 0.6, a2.Strength)
}

func TestArgumentMoveTo(t *testing.T) {
	setupDB()
	defer teardownDB()

	c1 := Claim{Title: "Claim 1", Truth: 0.5, TruthRU: 0.6}
	c2 := Claim{Title: "Claim 2", Truth: 0.4, TruthRU: 0.5}
	c3 := Claim{Title: "Claim 3", Truth: 0.6, TruthRU: 0.7}
	c4 := Claim{Title: "Claim 4", Truth: 0.3, TruthRU: 0.4}
	c5 := Claim{Title: "Claim 5", Truth: 0.7, TruthRU: 0.8}
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)
	TESTDB.Create(&c3)
	TESTDB.Create(&c4)
	TESTDB.Create(&c5)

	a1 := Argument{Title: "Argument 1", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Strength: 0.1, StrengthRU: 0.2, Type: ARGUMENT_FOR}
	a2 := Argument{Title: "Argument 2", TargetClaimID: NUUID(c1.ID), ClaimID: c3.ID, Strength: 0.2, StrengthRU: 0.3, Type: ARGUMENT_AGAINST}
	a3 := Argument{Title: "Argument 3", TargetClaimID: NUUID(c1.ID), ClaimID: c4.ID, Strength: 0.6, StrengthRU: 0.7, Type: ARGUMENT_FOR}
	a4 := Argument{Title: "Argument 4", TargetClaimID: NUUID(c1.ID), ClaimID: c5.ID, Strength: 0.7, StrengthRU: 0.6, Type: ARGUMENT_AGAINST}
	TESTDB.Create(&a1)
	TESTDB.Create(&a2)
	TESTDB.Create(&a3)
	TESTDB.Create(&a4)

	u1 := User{Name: "User 1", Username: "user1", Email: "email1@gruff.org"}
	u2 := User{Name: "User 2", Username: "user2", Email: "email2@gruff.org"}
	u3 := User{Name: "User 3", Username: "user3", Email: "email3@gruff.org"}
	u4 := User{Name: "User 4", Username: "user4", Email: "email4@gruff.org"}
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)
	TESTDB.Create(&u3)
	TESTDB.Create(&u4)

	co11 := ClaimOpinion{UserID: u1.ID, ClaimID: c1.ID, Truth: 0.2398}
	co12 := ClaimOpinion{UserID: u2.ID, ClaimID: c1.ID, Truth: 0.3290}
	co14 := ClaimOpinion{UserID: u4.ID, ClaimID: c1.ID, Truth: 0.0290}
	co21 := ClaimOpinion{UserID: u1.ID, ClaimID: c2.ID, Truth: 0.389}
	co32 := ClaimOpinion{UserID: u2.ID, ClaimID: c3.ID, Truth: 0.34985}
	co33 := ClaimOpinion{UserID: u3.ID, ClaimID: c3.ID, Truth: 0.5487}
	co44 := ClaimOpinion{UserID: u4.ID, ClaimID: c4.ID, Truth: 0.4839}
	TESTDB.Create(&co11)
	TESTDB.Create(&co12)
	TESTDB.Create(&co14)
	TESTDB.Create(&co21)
	TESTDB.Create(&co32)
	TESTDB.Create(&co33)
	TESTDB.Create(&co44)

	ao13 := ArgumentOpinion{UserID: u3.ID, ArgumentID: a1.ID, Strength: 0.23984}
	ao14 := ArgumentOpinion{UserID: u4.ID, ArgumentID: a1.ID, Strength: 0.923}
	ao21 := ArgumentOpinion{UserID: u1.ID, ArgumentID: a2.ID, Strength: 0.23984}
	ao23 := ArgumentOpinion{UserID: u3.ID, ArgumentID: a2.ID, Strength: 0.9832}
	ao32 := ArgumentOpinion{UserID: u2.ID, ArgumentID: a3.ID, Strength: 0.4839}
	ao42 := ArgumentOpinion{UserID: u2.ID, ArgumentID: a4.ID, Strength: 0.2983}
	ao44 := ArgumentOpinion{UserID: u4.ID, ArgumentID: a4.ID, Strength: 0.89384}
	TESTDB.Create(&ao13)
	TESTDB.Create(&ao14)
	TESTDB.Create(&ao21)
	TESTDB.Create(&ao23)
	TESTDB.Create(&ao32)
	TESTDB.Create(&ao42)
	TESTDB.Create(&ao44)

	err := a2.MoveTo(&CTX, a1.ID, ARGUMENT_FOR, OBJECT_TYPE_ARGUMENT)

	assert.Nil(t, err)
	// Test removal of opinions
	cos := []ClaimOpinion{}
	TESTDB.Find(&cos)
	assert.Equal(t, 7, len(cos))
	aos := []ArgumentOpinion{}
	TESTDB.Find(&aos)
	assert.Equal(t, 5, len(aos))
	TESTDB.Where("id = ?", a2.ID).First(&a2)
	assert.Nil(t, a2.TargetClaimID)
	assert.Equal(t, a1.ID, a2.TargetArgumentID.UUID)
	assert.Equal(t, ARGUMENT_FOR, a2.Type)
	// Test notification of voters
	ns := []Notification{}
	TESTDB.Scopes(FindArgumentMovedNotifications).Where("item_id = ?", a2.ID).Order("user_id ASC").Find(&ns)
	assert.Equal(t, 2, len(ns))
	assert.Equal(t, u1.ID, ns[0].UserID)
	assert.Equal(t, NOTIFICATION_TYPE_MOVED, ns[0].Type)
	assert.Equal(t, a2.ID, ns[0].ItemID.UUID)
	assert.Equal(t, OBJECT_TYPE_ARGUMENT, *ns[0].ItemType)
	assert.Equal(t, c1.ID, ns[0].OldID.UUID)
	assert.Equal(t, OBJECT_TYPE_CLAIM, *ns[0].OldType)
	assert.Equal(t, u3.ID, ns[1].UserID)
	assert.Equal(t, NOTIFICATION_TYPE_MOVED, ns[1].Type)
	assert.Equal(t, a2.ID, ns[1].ItemID.UUID)
	assert.Equal(t, OBJECT_TYPE_ARGUMENT, *ns[1].ItemType)
	assert.Equal(t, c1.ID, ns[1].OldID.UUID)
	assert.Equal(t, OBJECT_TYPE_CLAIM, *ns[1].OldType)
	// Test notification of sub-argument voters
	ns = []Notification{}
	TESTDB.Scopes(FindParentArgumentMovedNotifications).Where("item_id = ?", a2.ID).Order("user_id ASC").Find(&ns)
	assert.Equal(t, 0, len(ns))
	// Test recalc of rollup scores
	TESTDB.First(&c1)
	TESTDB.First(&c2)
	TESTDB.First(&c3)
	TESTDB.First(&c4)
	TESTDB.First(&c5)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	TESTDB.First(&a3)
	TESTDB.First(&a4)
	assert.Equal(t, 0.5089, support.RoundToDecimal(c1.TruthRU, 4))
	assert.Equal(t, 0.5, c2.TruthRU)
	assert.Equal(t, 0.7, c3.TruthRU)
	assert.Equal(t, 0.4, c4.TruthRU)
	assert.Equal(t, 0.8, c5.TruthRU)
	assert.Equal(t, 0.605, a1.StrengthRU)
	assert.Equal(t, 0.3, a2.StrengthRU)
	assert.Equal(t, 0.7, a3.StrengthRU)
	assert.Equal(t, 0.6, a4.StrengthRU)

	// Set up some sub-arguments
	a5 := Argument{Title: "Argument 5", TargetArgumentID: NUUID(a3.ID), ClaimID: c4.ID, Strength: 0.6, StrengthRU: 0.7, Type: ARGUMENT_FOR}

	a6 := Argument{Title: "Argument 6", TargetArgumentID: NUUID(a3.ID), ClaimID: c5.ID, Strength: 0.7, StrengthRU: 0.6, Type: ARGUMENT_AGAINST}
	TESTDB.Create(&a5)
	TESTDB.Create(&a6)

	ao52 := ArgumentOpinion{UserID: u2.ID, ArgumentID: a5.ID, Strength: 0.2983}
	ao53 := ArgumentOpinion{UserID: u3.ID, ArgumentID: a5.ID, Strength: 0.89384}
	ao61 := ArgumentOpinion{UserID: u1.ID, ArgumentID: a6.ID, Strength: 0.89384}
	TESTDB.Create(&ao52)
	TESTDB.Create(&ao53)
	TESTDB.Create(&ao61)

	err = a3.MoveTo(&CTX, c2.ID, ARGUMENT_AGAINST, OBJECT_TYPE_CLAIM)

	assert.Nil(t, err)
	// Test removal of opinions
	cos = []ClaimOpinion{}
	TESTDB.Find(&cos)
	assert.Equal(t, 7, len(cos))
	aos = []ArgumentOpinion{}
	TESTDB.Find(&aos)
	assert.Equal(t, 7, len(aos))
	TESTDB.Where("id = ?", a3.ID).First(&a3)
	assert.Nil(t, a3.TargetArgumentID)
	assert.Equal(t, c2.ID, a3.TargetClaimID.UUID)
	assert.Equal(t, ARGUMENT_AGAINST, a3.Type)
	ns = []Notification{}
	// Test notification of voters
	TESTDB.Scopes(FindArgumentMovedNotifications).Where("item_id = ?", a3.ID).Order("user_id ASC").Find(&ns)
	assert.Equal(t, 1, len(ns))
	assert.Equal(t, u2.ID, ns[0].UserID)
	assert.Equal(t, NOTIFICATION_TYPE_MOVED, ns[0].Type)
	assert.Equal(t, a3.ID, ns[0].ItemID.UUID)
	assert.Equal(t, OBJECT_TYPE_ARGUMENT, *ns[0].ItemType)
	assert.Equal(t, c1.ID, ns[0].OldID.UUID)
	assert.Equal(t, OBJECT_TYPE_CLAIM, *ns[0].OldType)
	// Test notification of sub-argument voters
	ns = []Notification{}
	TESTDB.Scopes(FindParentArgumentMovedNotifications).Where("item_id = ?", a3.ID).Order("user_id ASC").Find(&ns)
	assert.Equal(t, 3, len(ns))
	assert.Equal(t, u1.ID, ns[0].UserID)
	assert.Equal(t, NOTIFICATION_TYPE_PARENT_MOVED, ns[0].Type)
	assert.Equal(t, a3.ID, ns[0].ItemID.UUID)
	assert.Equal(t, OBJECT_TYPE_ARGUMENT, *ns[0].ItemType)
	assert.Equal(t, c1.ID, ns[0].OldID.UUID)
	assert.Equal(t, OBJECT_TYPE_CLAIM, *ns[0].OldType)
	assert.Equal(t, u2.ID, ns[1].UserID)
	assert.Equal(t, NOTIFICATION_TYPE_PARENT_MOVED, ns[1].Type)
	assert.Equal(t, a3.ID, ns[1].ItemID.UUID)
	assert.Equal(t, OBJECT_TYPE_ARGUMENT, *ns[1].ItemType)
	assert.Equal(t, c1.ID, ns[1].OldID.UUID)
	assert.Equal(t, OBJECT_TYPE_CLAIM, *ns[1].OldType)
	assert.Equal(t, u3.ID, ns[2].UserID)
	assert.Equal(t, NOTIFICATION_TYPE_PARENT_MOVED, ns[2].Type)
	assert.Equal(t, a3.ID, ns[2].ItemID.UUID)
	assert.Equal(t, OBJECT_TYPE_ARGUMENT, *ns[2].ItemType)
	assert.Equal(t, c1.ID, ns[2].OldID.UUID)
	assert.Equal(t, OBJECT_TYPE_CLAIM, *ns[2].OldType)
	// Test recalc of rollup scores
	TESTDB.First(&c1)
	TESTDB.First(&c2)
	TESTDB.First(&c3)
	TESTDB.First(&c4)
	TESTDB.First(&c5)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	TESTDB.First(&a3)
	TESTDB.First(&a4)
	assert.Equal(t, 0.3689, support.RoundToDecimal(c1.TruthRU, 4))
	assert.Equal(t, 0.36, c2.TruthRU)
	assert.Equal(t, 0.7, c3.TruthRU)
	assert.Equal(t, 0.4, c4.TruthRU)
	assert.Equal(t, 0.8, c5.TruthRU)
	assert.Equal(t, 0.605, a1.StrengthRU)
	assert.Equal(t, 0.3, a2.StrengthRU)
	assert.Equal(t, 0.7, a3.StrengthRU)
	assert.Equal(t, 0.6, a4.StrengthRU)
}

func TestArgumentScoreRU(t *testing.T) {
	setupDB()
	defer teardownDB()

	c1 := Claim{Title: "Claim 1", TruthRU: 0.5}
	c2 := Claim{Title: "Claim 2", Truth: 0.4}
	c3 := Claim{Title: "Claim 3", TruthRU: 0.6}
	c4 := Claim{Title: "Claim 4", TruthRU: 0.3}
	c5 := Claim{Title: "Claim 5", TruthRU: 0.7}
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)
	TESTDB.Create(&c3)
	TESTDB.Create(&c4)
	TESTDB.Create(&c5)

	a1 := Argument{Title: "Argument 1", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Strength: 0.1, StrengthRU: 0.15}
	a2 := Argument{Title: "Argument 2", TargetClaimID: NUUID(c1.ID), ClaimID: c3.ID, Strength: 0.2}
	a3 := Argument{Title: "Argument 3", TargetClaimID: NUUID(c1.ID), ClaimID: c4.ID, Strength: 0.6, StrengthRU: 0.7}
	a4 := Argument{Title: "Argument 4", TargetClaimID: NUUID(c1.ID), ClaimID: c5.ID, Strength: 0.7, StrengthRU: 0.65}
	TESTDB.Create(&a1)
	TESTDB.Create(&a2)
	TESTDB.Create(&a3)
	TESTDB.Create(&a4)

	assert.Equal(t, 0.06, a1.ScoreRU(&CTX))
	assert.Equal(t, 0.12, a2.ScoreRU(&CTX))
	assert.Equal(t, 0.21, a3.ScoreRU(&CTX))
	assert.Equal(t, 0.455, support.RoundToDecimal(a4.ScoreRU(&CTX), 3))
}
