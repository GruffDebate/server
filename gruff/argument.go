package gruff

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

/*
 * Arguments are described in detail in the Canonical Debate White Paper: https://github.com/canonical-debate-lab/paper#312_Argument

  An Argument connects a Claim to another Claim or Argument

  That is:
     a Claim can be used as an ARGUMENT to either prove or disprove the truth of a claim,
     or to modify the relevance or impact of another argument.

  The TYPE of the argument indicates how the claim (or CLAIM) is being used:
    PRO TRUTH: The Claim is a claim that is being used to prove the truth of another claim
      Ex: "The defendant was in Cincinatti on the date of the murder"
    CON TRUTH: The Claim is used as evidence against another claim
      Ex: "The defendant was hospitalized on the date of the murder"
    PRO RELEVANCE: The Claim is being used to show that another Argument is relevant and/or important
      Ex: "The murder occurred in Cincinatti"
      Ex: "This argument clearly shows that the defendant has no alibi"
    CON RELEVANCE: The Claim is being used to show that another Argument is irrelevant and/or unimportant
      Ex: "The murder occurred in the same hospital in which the defendant was hospitalized"
      Ex: "There is no evidence that the defendant ever left their room"

  A quick explanation of the fields:
    Claim: The Debate (or claim) that is being used as the basis of the argument
    Target Claim: The "parent" Claim against which a pro/con truth argument is being made
    Target Argument: In the case of a relevance or impact argument, the argument to which it refers
    Strength: The strength of an Argument is a combination of the Truth of its underlying Claim, and the Relevance Score.
      It is a cached value derived from the Flat popular votes, as described here: https://github.com/canonical-debate-lab/paper#33311_Flat_Scores
      and here: https://github.com/canonical-debate-lab/paper#33323_Popular_Vote
    StrengthRU: The roll-up score, similar to Strength, but rolled up to Level 1 as described here: https://github.com/canonical-debate-lab/paper#33312_Rollup_Scores

  To help understand the difference between relevance and impact arguments, imagine an argument is a bullet:
    Impact is the size of your bullet
    Relevance is how well you hit your target
  (note that because this difference is subtle enough to be difficult to separate one from the other,
   the two concepts are reflected together in a single score called Relevance)

  Scoring:
    Truth: 1.0 = definitely true; 0.5 = equal chance true or false; 0.0 = definitely false. "The world is flat" should have a 0.000000000000000001 truth score.
    Relevance: 1.0 = Completely on-topic and important; 0.5 = Circumstantial or somewhat relevant; 0.01 = Totally off-point, should be ignored
    Strength: 1.0 = This argument is definitely the most important argument for this side - no need to read any others; 0.5 = This is one more argument to consider; 0.01 = Probably not even worth including in the discussion
*/

type Argument struct {
	VersionedModel
	TargetClaimID    *string    `json:"targetClaimId,omitempty"`
	TargetClaim      *Claim     `json:"targetClaim,omitempty"`
	TargetArgumentID *string    `json:"targetArgId,omitempty"`
	TargetArgument   *Argument  `json:"targetArg,omitempty"`
	ClaimID          string     `json:"claimId"`
	Claim            *Claim     `json:"claim,omitempty"`
	Title            string     `json:"title" valid:"length(3|1000),required"`
	Negation         string     `json:"negation"`
	Question         string     `json:"question"`
	Description      string     `json:"desc" valid:"length(3|4000)"`
	Note             string     `json:"note"`
	Pro              bool       `json:"pro"`
	Strength         float64    `json:"strength"`
	StrengthRU       float64    `json:"strengthRU"`
	ProArgs          []Argument `json:"proargs"`
	ConArgs          []Argument `json:"conargs"`
}

// ArangoObject interface

func (a Argument) CollectionName() string {
	return "arguments"
}

func (a Argument) ArangoKey() string {
	return a.Key
}

func (a Argument) ArangoID() string {
	return fmt.Sprintf("%s/%s", a.CollectionName(), a.ArangoKey())
}

func (a Argument) DefaultQueryParameters() ArangoQueryParameters {
	return DEFAULT_QUERY_PARAMETERS
}

func (a *Argument) Create(ctx *ServerContext) Error {
	// TODO: Test
	can, err := a.UserCanCreate(ctx)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You must be logged in to create this item")
	}

	col, err := ctx.Arango.CollectionFor(a)
	if err != nil {
		return err
	}

	var baseClaim Claim
	if a.ClaimID == "" {
		// Need to create a Base Claim for this Argument with the same title and description
		baseClaim = Claim{
			Title:       a.Title,
			Description: a.Description,
			Negation:    a.Negation,
			Question:    a.Question,
			Note:        a.Note,
		}
		if err := baseClaim.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
		a.ClaimID = baseClaim.ID
	} else {
		baseClaim.ID = a.ClaimID
		if err = baseClaim.Load(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// TODO: Test
	if err := a.ValidateForCreate(); err != nil {
		return err
	}

	a.PrepareForCreate(ctx)

	if _, dberr := col.CreateDocument(ctx.Context, a); dberr != nil {
		ctx.Rollback()
		return NewServerError(dberr.Error())
	}

	edge := BaseClaimEdge{Edge: Edge{
		From: a.ArangoID(),
		To:   baseClaim.ArangoID(),
	}}

	if err := edge.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	if a.TargetClaimID != nil {
		targetClaim := Claim{}
		targetClaim.ID = *a.TargetClaimID
		if err = (&targetClaim).Load(ctx); err != nil {
			ctx.Rollback()
			return err
		}
		if err = targetClaim.AddArgument(ctx, *a); err != nil {
			ctx.Rollback()
			return err
		}
	} else {
		targetArg := Argument{}
		targetArg.ID = *a.TargetArgumentID
		if err = targetArg.Load(ctx); err != nil {
			ctx.Rollback()
			return err
		}
		if err = targetArg.AddArgument(ctx, *a); err != nil {
			ctx.Rollback()
			return err
		}
	}

	return nil
}

func (a *Argument) Update(ctx *ServerContext, updates map[string]interface{}) Error {
	if err := a.ValidateForUpdate(updates); err != nil {
		return err
	}

	// TODO: Test
	can, err := a.UserCanUpdate(ctx, updates)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to modify this item")
	}

	col, err := ctx.Arango.CollectionFor(a)
	if err != nil {
		return err
	}

	// When an Argument is updated, it creates a new version
	if err := a.version(ctx); err != nil {
		return err
	}

	if _, err := col.UpdateDocument(ctx.Context, a.ArangoKey(), updates); err != nil {
		return NewServerError(err.Error())
	}

	return a.Load(ctx)
}

func (a *Argument) version(ctx *ServerContext) Error {
	oldVersion := *a

	// This should delete all the old edges, too
	if err := oldVersion.Delete(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	if err := a.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	// Find all edges going to old ver, make copy to new ver
	// Arguments
	inferences, err := oldVersion.Inferences(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range inferences {
		newEdge := Inference{Edge: Edge{
			From: a.ArangoID(),
			To:   edge.To,
		}}
		if err := newEdge.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// Inference edge
	inference, err := oldVersion.Inference(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	newInference := Inference{Edge: Edge{
		From: inference.From,
		To:   a.ArangoID(),
	}}
	if err := newInference.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	// Base Claim edge
	baseClaimEdge, err := oldVersion.BaseClaimEdge(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	newBaseClaimEdge := BaseClaimEdge{Edge: Edge{
		From: baseClaimEdge.From,
		To:   a.ArangoID(),
	}}
	if err := newBaseClaimEdge.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	// TODO: Contexts
	// TODO: Links

	return nil
}

func (a *Argument) Delete(ctx *ServerContext) Error {
	// TODO: test
	if err := DeleteArangoObject(ctx, a); err != nil {
		ctx.Rollback()
		return err
	}

	// Find all edges going to old ver, make copy to new ver
	inference, err := a.Inference(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	if err := inference.Delete(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	// Base Claim edge
	baseClaimEdge, err := a.BaseClaimEdge(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	if err := baseClaimEdge.Delete(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	// Arguments
	// WARNING: could create an infinite loop of deletions
	args, err := a.Arguments(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, arg := range args {
		if err := arg.Delete(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	return nil
}

// Restrictor
// TODO: Test
// TODO: Call in CRUD and other methods
func (a Argument) UserCanView(ctx *ServerContext) (bool, Error) {
	return true, nil
}

func (a Argument) UserCanCreate(ctx *ServerContext) (bool, Error) {
	return ctx.UserLoggedIn(), nil
}

func (a Argument) UserCanUpdate(ctx *ServerContext, updates map[string]interface{}) (bool, Error) {
	return a.UserCanDelete(ctx)
}

func (a Argument) UserCanDelete(ctx *ServerContext) (bool, Error) {
	u := ctx.UserContext
	if u.Curator {
		return true, nil
	}
	return a.CreatedByID == u.ArangoID(), nil
}

// Validator

func (a Argument) ValidateForCreate() Error {
	if err := a.ValidateField("title"); err != nil {
		return err
	}
	if err := a.ValidateField("desc"); err != nil {
		return err
	}
	if err := a.ValidateIDs(); err != nil {
		return err
	}
	return nil
}

func (a Argument) ValidateForUpdate(updates map[string]interface{}) Error {
	if a.DeletedAt != nil {
		return NewBusinessError("An argument that has already been deleted, or has a newer version, cannot be modified.")
	}
	if err := SetJsonValuesOnStruct(&a, updates); err != nil {
		return err
	}
	return a.ValidateForCreate()
}

func (a Argument) ValidateForDelete() Error {
	if a.DeletedAt != nil {
		return NewBusinessError("This argument has already been deleted or versioned.")
	}
	return nil
}

func (a Argument) ValidateField(f string) Error {
	err := ValidateStructField(a, f)
	return err
}

func (a Argument) ValidateIDs() Error {
	if a.ClaimID == "" {
		return NewBusinessError("claimId: non zero value required;")
	}
	if a.TargetClaimID == nil && a.TargetArgumentID == nil {
		return NewBusinessError("An Argument must have a target Claim or target Argument ID")
	}
	if a.TargetClaimID != nil && a.TargetArgumentID != nil {
		return NewBusinessError("An Argument can have only one target Claim or target Argument ID")
	}
	return nil
}

// Loader

func (a *Argument) Load(ctx *ServerContext) Error {
	db := ctx.Arango.DB

	col, err := ctx.Arango.CollectionFor(a)
	if err != nil {
		return err
	}

	if a.ArangoKey() != "" {
		_, dberr := col.ReadDocument(ctx.Context, a.ArangoKey(), a)
		if dberr != nil {
			return NewServerError(dberr.Error())
		}
	} else if a.ID != "" {
		var empty time.Time
		var query string
		bindVars := map[string]interface{}{
			"id": a.ID,
		}
		if a.CreatedAt == empty {
			query = fmt.Sprintf("FOR a IN %s FILTER a.id == @id AND a.end == null SORT a.start DESC LIMIT 1 RETURN a", a.CollectionName())
		} else {
			bindVars["start"] = a.CreatedAt
			query = fmt.Sprintf("FOR a IN %s FILTER a.id == @id AND a.start <= @start AND (a.end == null OR a.end > @start) SORT a.start DESC LIMIT 1 RETURN a", a.CollectionName())
		}
		cursor, err := db.Query(ctx.Context, query, bindVars)
		defer CloseCursor(cursor)
		if err != nil {
			return NewServerError(err.Error())
		}
		for cursor.HasMore() {
			_, err := cursor.ReadDocument(ctx.Context, a)
			if err != nil {
				return NewServerError(err.Error())
			}
		}
	} else {
		return NewBusinessError("There is no key or id for this Argument.")
	}

	return nil
}

// TODO
func (a *Argument) LoadFull(ctx *ServerContext) Error {
	if err := a.Load(ctx); err != nil {
		return err
	}

	args, err := a.Arguments(ctx)
	if err != nil {
		return err
	}

	var proArgs, conArgs []Argument
	for _, arg := range args {
		bc := Claim{}
		bc.ID = arg.ClaimID
		bc.QueryAt = a.QueryDate()
		if err := bc.Load(ctx); err != nil {
			return err
		}
		bc.QueryAt = nil
		arg.Claim = &bc

		if arg.Pro {
			proArgs = append(proArgs, arg)
		} else {
			conArgs = append(conArgs, arg)
		}
	}

	a.ProArgs = proArgs
	a.ConArgs = conArgs

	baseClaim := Claim{}
	baseClaim.ID = a.ClaimID
	baseClaim.QueryAt = a.QueryDate()
	if err = baseClaim.LoadFull(ctx); err != nil {
		return err
	}
	baseClaim.QueryAt = nil
	a.Claim = &baseClaim

	return nil
}

// Business methods

func (a Argument) AddArgument(ctx *ServerContext, arg Argument) Error {
	// TODO: test
	updates := map[string]interface{}{}
	if err := a.ValidateForUpdate(updates); err != nil {
		return err
	}

	// TODO: Test
	can, err := a.UserCanUpdate(ctx, updates)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to modify this item")
	}

	edge := Inference{Edge: Edge{
		From: a.ArangoID(),
		To:   arg.ArangoID(),
	}}

	if err := edge.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}
	return nil
}

func (a Argument) Arguments(ctx *ServerContext) ([]Argument, Error) {
	db := ctx.Arango.DB
	args := []Argument{}

	bindVars := map[string]interface{}{
		"arg": a.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s
                                 FOR a IN %s
                                   FILTER obj._to == a._id
                                      AND obj._from == @arg
                                   %s
                                   SORT a.start ASC
                                   RETURN a`,
		Inference{}.CollectionName(),
		Argument{}.CollectionName(),
		a.DateFilter(bindVars),
	)
	cursor, err := db.Query(ctx.Context, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return args, NewServerError(err.Error())
	}
	for cursor.HasMore() {
		arg := Argument{}
		_, err := cursor.ReadDocument(ctx.Context, &arg)
		if err != nil {
			return args, NewServerError(err.Error())
		}
		args = append(args, arg)
	}

	return args, nil
}

// TODO: Make generic by moving method to inference.go
func (a Argument) Inferences(ctx *ServerContext) ([]Inference, Error) {
	db := ctx.Arango.DB
	edges := []Inference{}

	bindVars := map[string]interface{}{
		"from": a.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._from == @from
                                %s
                                RETURN obj`,
		Inference{}.CollectionName(),
		a.DateFilter(bindVars))
	cursor, err := db.Query(ctx.Context, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return edges, NewServerError(err.Error())
	}
	for cursor.HasMore() {
		edge := Inference{}
		_, err := cursor.ReadDocument(ctx.Context, &edge)
		if err != nil {
			return edges, NewServerError(err.Error())
		}
		edges = append(edges, edge)
	}

	return edges, nil
}

func (a Argument) Inference(ctx *ServerContext) (Inference, Error) {
	db := ctx.Arango.DB
	edge := Inference{}

	query := fmt.Sprintf("FOR e IN %s FILTER e._to == @to LIMIT 1 RETURN e", edge.CollectionName())
	bindVars := map[string]interface{}{
		"to": a.ArangoID(),
	}
	cursor, err := db.Query(ctx.Context, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return edge, NewServerError(err.Error())
	}
	for cursor.HasMore() {
		_, err := cursor.ReadDocument(ctx.Context, &edge)
		if err != nil {
			return edge, NewServerError(err.Error())
		}
	}

	return edge, nil
}

func (a Argument) BaseClaimEdge(ctx *ServerContext) (BaseClaimEdge, Error) {
	db := ctx.Arango.DB
	edge := BaseClaimEdge{}

	query := fmt.Sprintf("FOR e IN %s FILTER e._from == @from LIMIT 1 RETURN e", edge.CollectionName())
	bindVars := map[string]interface{}{
		"from": a.ArangoID(),
	}
	cursor, err := db.Query(ctx.Context, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return edge, NewServerError(err.Error())
	}
	for cursor.HasMore() {
		_, err := cursor.ReadDocument(ctx.Context, &edge)
		if err != nil {
			return edge, NewServerError(err.Error())
		}
	}

	return edge, nil
}

// TODO: Create method should set default Strength to 0.5
// TODO: implement curator permissions

/*
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
	// TODO: use strategy pattern for different scoring mechanisms? Or leave external?
	// TODO: use latest algorithm
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

*/
/*
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

func (a *Argument) MoveTo(ctx *ServerContext, newId uuid.UUID, t, objType int) Error {
	db := ctx.Database

	oldArg := Argument{TargetClaimID: a.TargetClaimID, TargetArgumentID: a.TargetArgumentID, Type: a.Type}
	oldTargetID := a.TargetArgumentID
	oldTargetType := OBJECT_TYPE_ARGUMENT
	if oldTargetID == nil {
		oldTargetID = a.TargetClaimID
		oldTargetType = OBJECT_TYPE_CLAIM
	}

	switch objType {
	case OBJECT_TYPE_CLAIM:
		newClaim := Claim{}
		if err := db.Where("id = ?", newId).First(&newClaim).Error; err != nil {
			return NewNotFoundError(err.Error())
		}

		newIdN := NullableUUID{newId}
		a.TargetClaimID = &newIdN
		a.TargetClaim = &newClaim
		a.TargetArgumentID = nil

	case OBJECT_TYPE_ARGUMENT:
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
	if err := a.ValidateType(); err != nil {
		return err
	}

	if err := db.Set("gorm:save_associations", false).Save(a).Error; err != nil {
		return NewServerError(err.Error())
	}

	// TODO: Goroutine

	// TODO: More intelligent way to update scores?

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
*/

// Scopes

func OrderByBestArgument(db *gorm.DB) *gorm.DB {
	return db.Joins("LEFT JOIN claims c ON c.id = arguments.claim_id").
		Order("(arguments.strength * c.truth) DESC")
}
