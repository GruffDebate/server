package gruff

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

const ARGUMENT_TYPE_PRO_TRUTH int = 1
const ARGUMENT_TYPE_CON_TRUTH int = 2
const ARGUMENT_TYPE_PRO_STRENGTH int = 3
const ARGUMENT_TYPE_CON_STRENGTH int = 4

/*
  An Argument connects a Claim to another Claim or Argument
  That is:
     a Claim can be used as an ARGUMENT to either prove or disprove the truth of a claim,
     or to modify the relevance or impact of another argument.

  The TYPE of the argument indicates how the claim (or CLAIM) is being used:
    PRO TRUTH: The Claim is a claim that is being used to prove the truth of another claim
      Ex: "The defendant was in Cincinatti on the date of the murder"
    CON TRUTH: The Claim is used as evidence against another claim
      Ex: "The defendant was hospitalized on the date of the murder"
    PRO RELEVANCE: The Claim is being used to show that another Argument is relevant
      Ex: "The murder occurred in Cincinatti"
    CON RELEVANCE: The Claim is being used to show that another Argument is irrelevant
      Ex: "The murder occurred in the same hospital in which the defendant was hospitalized"
    PRO IMPACT: The Claim is being used to show the importance of another Argument
      Ex: "This argument clearly shows that the defendant has no alibi"
    CON IMPACT: The Claim is being used to diminish the importance of another argument
      Ex: "There is no evidence that the defendant ever left their room"

  A quick explanation of the fields:
    Claim: The Debate (or claim) that is being used as an argument
    Target Claim: The "parent" Claim against which a pro/con truth argument is being made
    Target Argument: In the case of a relevance or impact argument, the argument to which it refers

  To help understand the difference between relevance and impact arguments, imagine an argument is a bullet:
    Impact is the size of your bullet
    Relevance is how well you hit your target

  Scoring:
    Truth: 1.0 = definitely true; 0.5 = equal chance true or false; 0.0 = definitely false. "The world is flat" should have a 0.000000000000000001 truth score.
    Impact: 1.0 = This argument is definitely the most important argument for this side - no need to read any others; 0.5 = This is one more argument to consider; 0.01 = Probably not even worth including in the discussion
    Relevance: 1.0 = Completely germaine and on-topic; 0.5 = Circumstantial or somewhat relevant; 0.01 = Totally off-point, should be ignored

 *
 * Topoi for Resolutions of Definition (for scoring Relevance/Impact):
 * - Is the interpretation relevant? (relevance)
 * - Is the interpretation fair?
 * - How should we choose among competing interpretations? (impact)
 *
 * Topoi for Resolutions of Value (for scoring Relevance/Impact):
 * - Is the condition truly good or bad as alleged? (i.e. which values are impacted, and is it positive or negative?)
 * - Has the value been properly applied? (relevance)
 * - How should we choose among competing values? (impact)
 *
 * Topoi for Resolutions of Policy (this would look differently in our model - one Issue with multiple claims as solutions?):
 * - Is there a problem? (could be represented by a "Do nothing" claim)
 * - Where is the credit or blame due?
 * - Will the proposal solve the problem?
 * - On balance, will things be better off? (trade offs - need to measure each proposal against multiple values)
 *

 * Types of evidence (Pro/Con-Truth arguments) (not implemented in Gruff):
 * - Examples
 * - Statistics
 * - Tangible objects
 * - Testimony
 * - Social consensus

 * Fallacies: accusations of standard fallacies can be used as arguments against relevance, impact, or truth
 * - Fallacies of Inference (con-impact? or con-relevance?):
 *   - Hasty generalizations
 *   - Unrepresentative samples
 *   - Fallacy of composition (if one is, then all are)
 *   - Fallacy of division (if most are, then this subgroup must be)
 *   - Errors in inference from sign (correlation vs. causation)
 *
 * - Fallacies of Relevance:
 *   - Ad Hominem
 *   - Appeal to Unreasonable Emotion
 * ... more

 --> True definition of fallacy: an argument that subverts the purpose of resolving a disagreement

*/
type Argument struct {
	Identifier
	TargetClaimID    *NullableUUID `json:"targetClaimId,omitempty" sql:"type:uuid"`
	TargetClaim      *Claim        `json:"targetClaim,omitempty"`
	TargetArgumentID *NullableUUID `json:"targetArgId,omitempty" sql:"type:uuid"`
	TargetArgument   *Argument     `json:"targetArg,omitempty"`
	ClaimID          uuid.UUID     `json:"claimId" sql:"type:uuid;not null"`
	Claim            *Claim        `json:"claim,omitempty"`
	Title            string        `json:"title" sql:"not null" valid:"length(3|1000),required"`
	Description      string        `json:"desc" valid:"length(3|4000)"`
	Type             int           `json:"type" sql:"not null"`
	Strength         float64       `json:"strength"`
	StrengthRU       float64       `json:"strengthRU"`
	ProStrength      []Argument    `json:"prostr,omitempty"`
	ConStrength      []Argument    `json:"constr,omitempty"`
}

func (a Argument) ValidateForCreate() GruffError {
	err := a.ValidateField("Title")
	if err != nil {
		return err
	}
	err = a.ValidateField("Description")
	if err != nil {
		return err
	}
	err = a.ValidateField("Type")
	if err != nil {
		return err
	}
	err = a.ValidateIDs()
	if err != nil {
		return err
	}
	err = a.ValidateType()
	if err != nil {
		return err
	}
	return nil
}

func (a Argument) ValidateForUpdate() GruffError {
	return a.ValidateForCreate()
}

func (a Argument) ValidateField(f string) GruffError {
	err := ValidateStructField(a, f)
	return err
}

func (a Argument) ValidateIDs() GruffError {
	if a.ClaimID == uuid.Nil {
		return NewBusinessError("ClaimID: non zero value required;")
	}
	if (a.TargetClaimID == nil || a.TargetClaimID.UUID == uuid.Nil) &&
		(a.TargetArgumentID == nil || a.TargetArgumentID.UUID == uuid.Nil) {
		return NewBusinessError("An Argument must have a target Claim or target Argument ID")
	}
	if a.TargetClaimID != nil && a.TargetArgumentID != nil {
		return NewBusinessError("An Argument can have only one target Claim or target Argument ID")
	}
	return nil
}

func (a Argument) ValidateType() GruffError {
	switch a.Type {
	case ARGUMENT_TYPE_PRO_TRUTH, ARGUMENT_TYPE_CON_TRUTH:
		if a.TargetClaimID == nil || a.TargetClaimID.UUID == uuid.Nil {
			return NewBusinessError("A pro or con truth argument must refer to a target claim")
		}
	case ARGUMENT_TYPE_PRO_STRENGTH,
		ARGUMENT_TYPE_CON_STRENGTH:
		if a.TargetArgumentID == nil || a.TargetArgumentID.UUID == uuid.Nil {
			return NewBusinessError("An argument for or against argument strength must refer to a target argument")
		}
	default:
		return NewBusinessError("Type: invalid;")
	}
	return nil
}

// Business methods

func (a Argument) UpdateStrength(ctx *ServerContext) {
	ctx.Database.Exec("UPDATE arguments a SET strength = (SELECT AVG(strength) FROM argument_opinions WHERE argument_id = a.id) WHERE id = ?", a.ID)

	// TODO: test
	if a.StrengthRU == 0.0 {
		// There's no roll up score yet, so the strength score itself is affecting related roll ups
		a.UpdateAncestorRUs(ctx)
	}
}

func (a *Argument) UpdateStrengthRU(ctx *ServerContext) {
	// TODO: do it all in SQL?
	proArgs, conArgs := a.Arguments(ctx)

	if len(proArgs) > 0 || len(conArgs) > 0 {
		proScore := 0.0
		for _, arg := range proArgs {
			remainder := 1.0 - proScore
			score := arg.ScoreRU(ctx)
			addon := remainder * score
			proScore += addon
		}

		conScore := 0.0
		for _, arg := range conArgs {
			remainder := 1.0 - conScore
			score := arg.ScoreRU(ctx)
			addon := remainder * score
			conScore += addon
		}

		netScore := proScore - conScore
		netScore = 0.5 + 0.5*netScore

		a.StrengthRU = netScore
	} else {
		a.StrengthRU = 0.0
	}

	ctx.Database.Set("gorm:save_associations", false).Save(a)

	a.UpdateAncestorRUs(ctx)
}

func (a Argument) UpdateAncestorRUs(ctx *ServerContext) {
	if a.TargetClaimID != nil {
		claim := a.TargetClaim
		if claim == nil {
			claim = &Claim{}
			if err := ctx.Database.Where("id = ?", a.TargetClaimID).First(claim).Error; err != nil {
				return
			}
		}
		claim.UpdateTruthRU(ctx)
	} else {
		arg := a.TargetArgument
		if arg == nil {
			arg = &Argument{}
			if err := ctx.Database.Where("id = ?", a.TargetArgumentID).First(arg).Error; err != nil {
				fmt.Println("Error loading argument:", err.Error())
				return
			}
		}
		arg.UpdateStrengthRU(ctx)
	}
}

func (a *Argument) MoveTo(ctx *ServerContext, newId uuid.UUID, t int) GruffError {
	db := ctx.Database

	oldArg := Argument{TargetClaimID: a.TargetClaimID, TargetArgumentID: a.TargetArgumentID, Type: a.Type}
	oldTargetID := a.TargetArgumentID
	oldTargetType := OBJECT_TYPE_ARGUMENT
	if oldTargetID == nil {
		oldTargetID = a.TargetClaimID
		oldTargetType = OBJECT_TYPE_CLAIM
	}

	switch t {
	case ARGUMENT_TYPE_PRO_TRUTH, ARGUMENT_TYPE_CON_TRUTH:
		newClaim := Claim{}
		if err := db.Where("id = ?", newId).First(&newClaim).Error; err != nil {
			return NewNotFoundError(err.Error())
		}

		newIdN := NullableUUID{newId}
		a.TargetClaimID = &newIdN
		a.TargetClaim = &newClaim
		a.TargetArgumentID = nil

	case ARGUMENT_TYPE_PRO_STRENGTH, ARGUMENT_TYPE_CON_STRENGTH:
		newArg := Argument{}
		if err := db.Where("id = ?", newId).First(&newArg).Error; err != nil {
			return NewNotFoundError(err.Error())
		}

		newIdN := NullableUUID{newId}
		a.TargetArgumentID = &newIdN
		a.TargetArgument = &newArg
		a.TargetClaimID = nil

	default:
		return NewNotFoundError(fmt.Sprintf("Type unknown: %d", t))
	}
	a.Type = t

	if err := db.Set("gorm:save_associations", false).Save(a).Error; err != nil {
		return NewServerError(err.Error())
	}

	// TODO: Goroutine

	// Notify argument voters of move so they can vote again
	ops := []ArgumentOpinion{}
	if err := db.Where("argument_id = ?", a.ID).Find(&ops).Error; err != nil {
		db.Rollback()
		return NewServerError(err.Error())
	}

	for _, op := range ops {
		NotifyArgumentMoved(ctx, op.UserID, a.ID, oldTargetID.UUID, oldTargetType)
	}

	// Notify sub argument voters of move so they can double-check their vote
	uids := []uint64{}
	rows, dberr := db.Model(ArgumentOpinion{}).
		Select("DISTINCT argument_opinions.user_id").
		Where("argument_id IN (SELECT id FROM arguments WHERE target_argument_id = ?)", a.ID).
		Rows()
	defer rows.Close()

	if dberr == nil {
		for rows.Next() {
			var uid uint64
			err := rows.Scan(&uid)
			if err == nil {
				uids = append(uids, uid)
			}
		}
	}

	for _, uid := range uids {
		NotifyParentArgumentMoved(ctx, uid, a.ID, oldTargetID.UUID, oldTargetType)
	}

	// Clear opinions on the moved argument
	if err := db.Exec("DELETE FROM argument_opinions WHERE argument_id = ?", a.ID).Error; err != nil {
		db.Rollback()
		return NewServerError(err.Error())
	}

	a.UpdateAncestorRUs(ctx)
	oldArg.UpdateAncestorRUs(ctx)

	return nil
}

func (a Argument) Score(ctx *ServerContext) float64 {
	c := a.Claim
	if c == nil {
		c = &Claim{}
		ctx.Database.Where("id = ?", a.ClaimID).First(c)
	}

	return a.Strength * c.Truth
}

func (a Argument) ScoreRU(ctx *ServerContext) float64 {
	c := a.Claim
	if c == nil {
		c = &Claim{}
		ctx.Database.Where("id = ?", a.ClaimID).First(c)
	}

	truth := c.TruthRU
	if truth == 0.0 {
		truth = c.Truth
	}

	strength := a.StrengthRU
	if strength == 0.0 {
		strength = a.Strength
	}

	return strength * truth
}

func (a Argument) Arguments(ctx *ServerContext) (proArgs []Argument, conArgs []Argument) {
	proArgs = a.ProStrength
	conArgs = a.ConStrength

	if len(proArgs) == 0 {
		ctx.Database.
			Preload("Claim").
			Scopes(OrderByBestArgument).
			Where("type = ?", ARGUMENT_TYPE_PRO_STRENGTH).
			Where("target_argument_id = ?", a.ID).
			Find(&proArgs)
	}

	if len(conArgs) == 0 {
		ctx.Database.
			Preload("Claim").
			Scopes(OrderByBestArgument).
			Where("type = ?", ARGUMENT_TYPE_CON_STRENGTH).
			Where("target_argument_id = ?", a.ID).
			Find(&conArgs)
	}

	return
}

// Scopes

func OrderByBestArgument(db *gorm.DB) *gorm.DB {
	return db.Joins("LEFT JOIN claims c ON c.id = arguments.claim_id").
		Order("(arguments.strength * c.truth) DESC")
}
