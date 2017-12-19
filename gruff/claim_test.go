package gruff

import (
	"testing"

	"github.com/GruffDebate/server/support"
	"github.com/stretchr/testify/assert"
)

func TestCreateClaim(t *testing.T) {
	setupDB()
	defer teardownDB()

	d := Claim{
		Title:       "The first debate!",
		Description: "A description",
		Truth:       87.55,
	}
	TESTDB.Create(&d)

	assert.True(t, d.ID != ZERO_UUID)
}

func TestUpdateTruth(t *testing.T) {
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

	c1.UpdateTruth(&CTX)
	c2.UpdateTruth(&CTX)
	TESTDB.First(&c1)
	TESTDB.First(&c2)
	assert.Equal(t, 0.0, c1.Truth)
	assert.Equal(t, 0.0, c2.Truth)

	co1 := ClaimOpinion{UserID: 1, ClaimID: c1.ID, Truth: 0.5}
	TESTDB.Create(&co1)

	c1.UpdateTruth(&CTX)
	c2.UpdateTruth(&CTX)
	TESTDB.First(&c1)
	TESTDB.First(&c2)
	assert.Equal(t, 0.5, c1.Truth)
	assert.Equal(t, 0.0, c2.Truth)

	co2 := ClaimOpinion{UserID: 2, ClaimID: c2.ID, Truth: 0.3}
	TESTDB.Create(&co2)

	c1.UpdateTruth(&CTX)
	c2.UpdateTruth(&CTX)
	TESTDB.First(&c1)
	TESTDB.First(&c2)
	assert.Equal(t, 0.5, c1.Truth)
	assert.Equal(t, 0.3, c2.Truth)

	co3 := ClaimOpinion{UserID: 3, ClaimID: c2.ID, Truth: 0.5}
	TESTDB.Create(&co3)

	c1.UpdateTruth(&CTX)
	c2.UpdateTruth(&CTX)
	TESTDB.First(&c1)
	TESTDB.First(&c2)
	assert.Equal(t, 0.5, c1.Truth)
	assert.Equal(t, 0.4, c2.Truth)

	co4 := ClaimOpinion{UserID: 4, ClaimID: c1.ID, Truth: 0.9}
	co5 := ClaimOpinion{UserID: 5, ClaimID: c2.ID, Truth: 0.3}
	co6 := ClaimOpinion{UserID: 6, ClaimID: c2.ID, Truth: 0.3}
	co7 := ClaimOpinion{UserID: 7, ClaimID: c2.ID, Truth: 0.9}
	TESTDB.Create(&co4)
	TESTDB.Create(&co5)
	TESTDB.Create(&co6)
	TESTDB.Create(&co7)

	c1.UpdateTruth(&CTX)
	c2.UpdateTruth(&CTX)
	TESTDB.First(&c1)
	TESTDB.First(&c2)
	assert.Equal(t, 0.7, c1.Truth)
	assert.Equal(t, 0.46, c2.Truth)

	co6.Truth = 0.6
	TESTDB.Save(&co6)

	c1.UpdateTruth(&CTX)
	c2.UpdateTruth(&CTX)
	TESTDB.First(&c1)
	TESTDB.First(&c2)
	assert.Equal(t, 0.7, c1.Truth)
	assert.Equal(t, 0.52, c2.Truth)
}

func TestUpdateTruthRU(t *testing.T) {
	setupDB()
	defer teardownDB()

	c1 := Claim{Title: "C1", TruthRU: 0.5}
	c2 := Claim{Title: "C2", TruthRU: 1.0}
	c3 := Claim{Title: "C3", TruthRU: 0.5}
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)
	TESTDB.Create(&c3)

	(&c1).UpdateTruthRU(&CTX)
	assert.Equal(t, 0.0, c1.TruthRU)
	TESTDB.First(&c1)
	assert.Equal(t, 0.0, c1.TruthRU)

	a1 := Argument{Title: "A1", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Type: ARGUMENT_FOR, StrengthRU: 0.7}
	TESTDB.Create(&a1)

	(&c1).UpdateTruthRU(&CTX)
	assert.Equal(t, 0.85, c1.TruthRU)
	TESTDB.First(&c1)
	assert.Equal(t, 0.85, c1.TruthRU)

	a2 := Argument{Title: "Heinz 57", TargetClaimID: NUUID(c1.ID), ClaimID: c3.ID, Type: ARGUMENT_AGAINST, StrengthRU: 0.7}
	TESTDB.Create(&a2)

	(&c1).UpdateTruthRU(&CTX)
	assert.Equal(t, 0.675, c1.TruthRU)
	TESTDB.First(&c1)
	assert.Equal(t, 0.675, c1.TruthRU)

	a3 := Argument{Title: "Worcestershire", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Type: ARGUMENT_AGAINST, StrengthRU: 0.8}
	TESTDB.Create(&a3)

	(&c1).UpdateTruthRU(&CTX)
	assert.Equal(t, 0.415, c1.TruthRU)

	a4 := Argument{Title: "Miracle Whip", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Type: ARGUMENT_AGAINST, StrengthRU: 0.2}
	a5 := Argument{Title: "Grey Poupon", TargetClaimID: NUUID(c1.ID), ClaimID: c3.ID, Type: ARGUMENT_FOR, StrengthRU: 0.3}
	a6 := Argument{Title: "Dijonnaise", TargetClaimID: NUUID(c1.ID), ClaimID: c3.ID, Type: ARGUMENT_AGAINST, StrengthRU: 0.7}
	a7 := Argument{Title: "1000 Island", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Type: ARGUMENT_FOR, StrengthRU: 0.6}
	a8 := Argument{Title: "Tabasco", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Type: ARGUMENT_FOR, StrengthRU: 0.5}
	TESTDB.Create(&a4)
	TESTDB.Create(&a5)
	TESTDB.Create(&a6)
	TESTDB.Create(&a7)
	TESTDB.Create(&a8)

	(&c1).UpdateTruthRU(&CTX)
	assert.Equal(t, 0.5083, c1.TruthRU)

	a9 := Argument{Title: "Tabasco", TargetClaimID: NUUID(c2.ID), ClaimID: c3.ID, Type: ARGUMENT_FOR, StrengthRU: 1.0}
	TESTDB.Create(&a9)

	(&c1).UpdateTruthRU(&CTX)
	assert.Equal(t, 0.5083, c1.TruthRU)
}

func TestUpdateAncestorRUs(t *testing.T) {
	setupDB()
	defer teardownDB()

	c1 := Claim{Title: "C1", Truth: 0.3, TruthRU: 0.5}
	c2 := Claim{Title: "C2", Truth: 0.8, TruthRU: 1.0}
	c3 := Claim{Title: "C3", Truth: 0.6, TruthRU: 0.5}
	c4 := Claim{Title: "C4", Truth: 0.5}
	c5 := Claim{Title: "C5", Truth: 0.3, TruthRU: 0.3}
	TESTDB.Create(&c1)
	TESTDB.Create(&c2)
	TESTDB.Create(&c3)
	TESTDB.Create(&c4)
	TESTDB.Create(&c5)

	a1 := Argument{Title: "Argument 1", TargetClaimID: NUUID(c1.ID), ClaimID: c2.ID, Strength: 0.1, StrengthRU: 0.15, Type: ARGUMENT_FOR}
	a2 := Argument{Title: "Argument 2", TargetClaimID: NUUID(c1.ID), ClaimID: c3.ID, Strength: 0.2, Type: ARGUMENT_AGAINST}
	a3 := Argument{Title: "Argument 3", TargetClaimID: NUUID(c1.ID), ClaimID: c4.ID, Strength: 0.6, StrengthRU: 0.7, Type: ARGUMENT_FOR}
	TESTDB.Create(&a1)
	TESTDB.Create(&a2)
	TESTDB.Create(&a3)

	a4 := Argument{Title: "Argument 4", TargetArgumentID: NUUID(a1.ID), ClaimID: c5.ID, Strength: 0.7, StrengthRU: 0.65, Type: ARGUMENT_FOR}
	TESTDB.Create(&a4)

	a1.UpdateAncestorRUs(&CTX)

	TESTDB.First(&c1)
	TESTDB.First(&c2)
	TESTDB.First(&c3)
	TESTDB.First(&c4)
	TESTDB.First(&c5)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	TESTDB.First(&a3)
	TESTDB.First(&a4)

	assert.Equal(t, 0.67375, support.RoundToDecimal(c1.TruthRU, 5))
	assert.Equal(t, 1.0, c2.TruthRU)
	assert.Equal(t, 0.5, c3.TruthRU)
	assert.Equal(t, 0.0, c4.TruthRU)
	assert.Equal(t, 0.3, c5.TruthRU)
	assert.Equal(t, 0.15, a1.StrengthRU)
	assert.Equal(t, 0.0, a2.StrengthRU)
	assert.Equal(t, 0.7, a3.StrengthRU)
	assert.Equal(t, 0.65, a4.StrengthRU)

	a2.UpdateAncestorRUs(&CTX)
	a3.UpdateAncestorRUs(&CTX)

	TESTDB.First(&c1)
	TESTDB.First(&c2)
	TESTDB.First(&c3)
	TESTDB.First(&c4)
	TESTDB.First(&c5)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	TESTDB.First(&a3)
	TESTDB.First(&a4)

	assert.Equal(t, 0.67375, support.RoundToDecimal(c1.TruthRU, 5))
	assert.Equal(t, 1.0, c2.TruthRU)
	assert.Equal(t, 0.5, c3.TruthRU)
	assert.Equal(t, 0.0, c4.TruthRU)
	assert.Equal(t, 0.3, c5.TruthRU)
	assert.Equal(t, 0.15, a1.StrengthRU)
	assert.Equal(t, 0.0, a2.StrengthRU)
	assert.Equal(t, 0.7, a3.StrengthRU)
	assert.Equal(t, 0.65, a4.StrengthRU)

	a4.UpdateAncestorRUs(&CTX)

	TESTDB.First(&c1)
	TESTDB.First(&c2)
	TESTDB.First(&c3)
	TESTDB.First(&c4)
	TESTDB.First(&c5)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	TESTDB.First(&a3)
	TESTDB.First(&a4)

	assert.Equal(t, 0.819188, support.RoundToDecimal(c1.TruthRU, 6))
	assert.Equal(t, 1.0, c2.TruthRU)
	assert.Equal(t, 0.5, c3.TruthRU)
	assert.Equal(t, 0.0, c4.TruthRU)
	assert.Equal(t, 0.3, c5.TruthRU)
	assert.Equal(t, 0.5975, a1.StrengthRU)
	assert.Equal(t, 0.0, a2.StrengthRU)
	assert.Equal(t, 0.7, a3.StrengthRU)
	assert.Equal(t, 0.65, a4.StrengthRU)

	c5.TruthRU = 0.4
	TESTDB.Save(&c5)
	c5.UpdateAncestorRUs(&CTX)

	TESTDB.First(&c1)
	TESTDB.First(&c2)
	TESTDB.First(&c3)
	TESTDB.First(&c4)
	TESTDB.First(&c5)
	TESTDB.First(&a1)
	TESTDB.First(&a2)
	TESTDB.First(&a3)
	TESTDB.First(&a4)

	assert.Equal(t, 0.83300, support.RoundToDecimal(c1.TruthRU, 6))
	assert.Equal(t, 1.0, c2.TruthRU)
	assert.Equal(t, 0.5, c3.TruthRU)
	assert.Equal(t, 0.0, c4.TruthRU)
	assert.Equal(t, 0.4, c5.TruthRU)
	assert.Equal(t, 0.64, a1.StrengthRU)
	assert.Equal(t, 0.0, a2.StrengthRU)
	assert.Equal(t, 0.7, a3.StrengthRU)
	assert.Equal(t, 0.0, a4.StrengthRU)
}
