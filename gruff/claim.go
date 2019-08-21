package gruff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/GruffDebate/server/support"
	arango "github.com/arangodb/go-driver"
)

/*
 * A Claim is a proposed statement of fact
 *
 * Claims are described in detail in the Canonical Debate White Paper: https://github.com/canonical-debate-lab/paper#311_Claim
 *
 * According to David Zarefsky (https://www.thegreatcoursesplus.com/argumentation/argument-analysis-and-diagramming) there are 4 types:
 * - Fact: Al Gore received more popular votes than George Bush in the 2000 election
 * - Definition: Capital execution is murder
 * - Value: Environmental protection is more important than economic growth
 * - Policy: Congress should pass the president's budget
 *
 * Complex Claims:
 * - Series: Because of X, Y happened, which caused Z --> Not modeled in Gruff
 * - Convergent: Airline travel is becoming more unpleasant because of X, Y, Z, P, D, and Q --> Supported by standard Gruff structure
 * - Parallel: Same as convergent, except that any one argument is enough --> Supported by standard Gruff structure
 *
 * Complex Claims in the Canonical Debate are known as Multi-Premise Claims: https://github.com/canonical-debate-lab/paper#3115_Multipremise_Claims
 * They are currently unsupported, but soon will be.
 *
 */
const PREMISE_RULE_NONE int = 0
const PREMISE_RULE_ALL int = 1
const PREMISE_RULE_ANY int = 2
const PREMISE_RULE_ANY_TWO int = 3

const DEFAULT_CLAIM_SCORE float32 = 0.50

type Claim struct {
	VersionedModel
	Title         string     `json:"title" valid:"length(3|1000)"`
	Negation      string     `json:"negation"`
	Question      string     `json:"question"`
	Description   string     `json:"desc" valid:"length(3|4000)"`
	Note          string     `json:"note"`
	Image         string     `json:"img,omitempty"`
	MultiPremise  bool       `json:"mp"`
	PremiseRule   int        `json:"mprule"`
	Truth         float32    `json:"truth"` // Average score from direct opinions
	PremiseClaims []Claim    `json:"premises,omitempty"`
	ProArgs       []Argument `json:"proargs"`
	ConArgs       []Argument `json:"conargs"`
	Links         []Link     `json:"links,omitempty"`
	ContextElems  []Context  `json:"contexts,omitempty" skip:"true"`
}

// ArangoObject interface

func (c Claim) CollectionName() string {
	return "claims"
}

func (c Claim) ArangoKey() string {
	return c.Key
}

func (c Claim) ArangoID() string {
	return fmt.Sprintf("%s/%s", c.CollectionName(), c.ArangoKey())
}

func (c Claim) DefaultQueryParameters() ArangoQueryParameters {
	return DEFAULT_QUERY_PARAMETERS
}

func (c *Claim) Create(ctx *ServerContext) Error {
	// Only allow one claim with the same ID that isn't deleted
	if c.ID != "" {
		oldClaim := Claim{}
		oldClaim.ID = c.ID
		err := oldClaim.Load(ctx)
		if err != nil || (oldClaim.Key != "" && oldClaim.DeletedAt == nil) {
			return NewBusinessError("A claim with the same ID already exists")
		}
	}

	c.Truth = DEFAULT_CLAIM_SCORE

	aerr := CreateArangoObject(ctx, c)
	if aerr != nil {
		ctx.Rollback()
		return aerr
	}

	// TODO TEST
	if len(c.ContextElems) > 0 {
		for _, item := range c.ContextElems {
			context := Context{}
			context.Key = item.ArangoKey()
			err := context.Load(ctx)
			if err != nil {
				ctx.Rollback()
				return err
			}

			err = c.AddContext(ctx, context)
			if err != nil {
				ctx.Rollback()
				return err
			}
		}
	}

	return nil
}

func (c *Claim) Update(ctx *ServerContext, updates Updates) Error {
	return UpdateArangoObject(ctx, c, updates)
}

func (c *Claim) version(ctx *ServerContext, updates Updates) Error {
	c.QueryAt = nil
	oldVersion := *c

	// This should delete all the old edges, too,
	// except Inferences and BaseClaimEdges
	if err := oldVersion.performDelete(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	// Don't use Create method, since it performs unnecessary checks
	if err := CreateArangoObject(ctx, c); err != nil {
		ctx.Rollback()
		return err
	}

	// Find all edges going to old ver, make copy to new ver
	if c.MultiPremise {
		premiseEdges, err := oldVersion.PremiseEdges(ctx)
		if err != nil {
			ctx.Rollback()
			return err
		}
		for _, edge := range premiseEdges {
			newEdge := PremiseEdge{
				Edge: Edge{
					From: c.ArangoID(),
					To:   edge.To,
				},
				Order: edge.Order,
			}
			if err := newEdge.Create(ctx); err != nil {
				ctx.Rollback()
				return err
			}
		}
	}

	// Arguments
	inferences, err := oldVersion.Inferences(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range inferences {
		newEdge := Inference{Edge: Edge{
			From: c.ArangoID(),
			To:   edge.To,
		}}
		if err := newEdge.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
		if err := edge.Delete(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// Base Claim edges
	baseClaimEdges, err := oldVersion.BaseClaimEdges(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range baseClaimEdges {
		newEdge := BaseClaimEdge{Edge: Edge{
			From: edge.From,
			To:   c.ArangoID(),
		}}
		if err := newEdge.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
		if err := edge.Delete(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// Any edges using this Claim as a Premise
	premiseEdges, err := oldVersion.EdgesToThisPremise(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range premiseEdges {
		newEdge := PremiseEdge{
			Edge: Edge{
				From: edge.From,
				To:   c.ArangoID(),
			},
			Order: edge.Order,
		}
		if err := newEdge.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// Contexts
	contextEdges, err := oldVersion.ContextEdges(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range contextEdges {
		if len(c.ContextElems) > 0 {
			for _, item := range c.ContextElems {
				if item.ArangoID() == edge.From {
					newEdge := ContextEdge{Edge: Edge{
						From: edge.From,
						To:   c.ArangoID(),
					}}
					if err := newEdge.Create(ctx); err != nil {
						ctx.Rollback()
						return err
					}
					break
				}
			}
		} else {
			newEdge := ContextEdge{Edge: Edge{
				From: edge.From,
				To:   c.ArangoID(),
			}}
			if err := newEdge.Create(ctx); err != nil {
				ctx.Rollback()
				return err
			}
		}
	}

	// TODO: Links

	// UserScores
	// TODO: Do this as a bulk operation
	// TODO: Test
	userScores, err := oldVersion.UserScores(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range userScores {
		newEdge := UserScore{
			Edge: Edge{
				From: edge.From,
				To:   c.ArangoID(),
			},
			Score: edge.Score,
		}
		if err := newEdge.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	return nil
}

func (c *Claim) Delete(ctx *ServerContext) Error {
	if err := c.ValidateForDelete(); err != nil {
		return err
	}

	can, err := c.UserCanDelete(ctx)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to delete this item")
	}

	bces, err := c.BaseClaimEdges(ctx)
	if err != nil {
		return err
	}
	if len(bces) > 0 {
		return NewBusinessError("You cannot delete a claim that is being used as a base claim for other arguments")
	}

	// Arguments - only when really deleting the Claim, rather than versioning
	// WARNING: could create an infinite loop of deletions
	args, err := c.Arguments(ctx)
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

	// Base Claim edges
	args, err = c.ArgumentsBasedOnThisClaim(ctx)
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

	return c.performDelete(ctx)
}

// Execute the delete action without verifications
func (c *Claim) performDelete(ctx *ServerContext) Error {
	// Find all edges going to old ver, make copy to new ver
	if c.MultiPremise {
		premiseEdges, err := c.PremiseEdges(ctx)
		if err != nil {
			ctx.Rollback()
			return err
		}
		for _, edge := range premiseEdges {
			if err := edge.Delete(ctx); err != nil {
				ctx.Rollback()
				return err
			}
		}
	}

	// Premise edges to this premise
	premiseEdges, err := c.EdgesToThisPremise(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range premiseEdges {
		otherClaim := Claim{}
		otherClaim.Key = edge.From[len(otherClaim.CollectionName())+1:]
		if err := otherClaim.Load(ctx); err != nil {
			ctx.Rollback()
			return err
		}
		if err := otherClaim.RemovePremise(ctx, c.ArangoID()); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// Contexts
	contextEdges, err := c.ContextEdges(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range contextEdges {
		if err := edge.Delete(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// TODO: Links

	// UserScores
	// TODO: Test
	filter := "obj._to == @claim"
	bindVars := BindVars{
		"claim": c.ArangoID(),
	}
	if err := DeleteArangoObjects(ctx, UserScore{}.CollectionName(), filter, bindVars); err != nil {
		ctx.Rollback()
		return err
	}

	if err := DeleteArangoObject(ctx, c); err != nil {
		ctx.Rollback()
		return err
	}
	return nil
}

// Restrictor

func (c Claim) UserCanView(ctx *ServerContext) (bool, Error) {
	return true, nil
}

func (c Claim) UserCanCreate(ctx *ServerContext) (bool, Error) {
	return ctx.UserLoggedIn(), nil
}

func (c Claim) UserCanUpdate(ctx *ServerContext, updates Updates) (bool, Error) {
	return c.UserCanDelete(ctx)
}

func (c Claim) UserCanDelete(ctx *ServerContext) (bool, Error) {
	u := ctx.UserContext
	if u.Curator {
		return true, nil
	}
	return c.CreatedByID == u.ArangoID(), nil
}

// Validator

func (c Claim) ValidateForCreate() Error {
	if len(c.Title) < 3 || len(c.Title) > 1000 {
		return NewBusinessError("Title: must be between 3 and 1000 characters;")
	}
	if len(c.Description) > 0 && (len(c.Description) < 3 || len(c.Description) > 4000) {
		return NewBusinessError("Description: must be blank, or between 3 and 4000 characters;")
	}
	return ValidateStruct(c)
}

func (c Claim) ValidateForUpdate(updates Updates) Error {
	if c.DeletedAt != nil {
		return NewBusinessError("A claim that has already been deleted, or has a newer version, cannot be modified")
	}
	if err := SetJsonValuesOnStruct(&c, updates); err != nil {
		return err
	}
	return c.ValidateForCreate()
}

func (c Claim) ValidateForDelete() Error {
	if c.DeletedAt != nil {
		return NewBusinessError("This claim has already been deleted or versioned")
	}
	return nil
}

func (c Claim) ValidateField(f string) Error {
	return ValidateStructField(c, f)
}

// Loader

// If the Claim object has a key, that exact Claim will be loaded
// Otherwise, Load will look for Claims matching the ID
// If QueryAt is a non-nil value, it will load the Claim active at that time (if any)
// Otherwise, it will return the current active (undeleted) version.
func (c *Claim) Load(ctx *ServerContext) Error {
	var err Error
	if c.ID != "" {
		bindVars := BindVars{
			"id": c.ID,
		}
		query := fmt.Sprintf(`FOR obj IN %s 
                                       FILTER obj.id == @id
                                       %s
                                       SORT obj.start DESC
                                       LIMIT 1 
                                       RETURN obj`,
			c.CollectionName(),
			c.DateFilter(bindVars))
		err = FindArangoObject(ctx, query, bindVars, c)
	} else if c.ArangoKey() != "" {
		err = LoadArangoObject(ctx, c, c.ArangoKey())
	} else {
		err = NewBusinessError("There is no key or id for this Claim")
	}

	return err
}

func (c *Claim) LoadFull(ctx *ServerContext) Error {
	queryAt := c.QueryAt
	if err := c.Load(ctx); err != nil {
		return err
	}
	c.QueryAt = queryAt

	if c.MultiPremise {
		premises, err := c.Premises(ctx)
		if err != nil {
			return err
		}

		fullPremises := make([]Claim, len(premises))
		for i, premise := range premises {
			premise.QueryAt = c.QueryDate()
			if err := premise.LoadFull(ctx); err != nil {
				return err
			}
			premise.QueryAt = nil
			fullPremises[i] = premise
		}

		c.PremiseClaims = fullPremises
	} else {
		args, err := c.Arguments(ctx)
		if err != nil {
			return err
		}

		var proArgs, conArgs []Argument
		for _, arg := range args {
			bc := Claim{}
			bc.ID = arg.ClaimID
			bc.QueryAt = c.QueryDate()
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

		c.ProArgs = proArgs
		c.ConArgs = conArgs
	}

	return nil
}

// Arguments

func (c Claim) Arguments(ctx *ServerContext) ([]Argument, Error) {
	args := []Argument{}

	bindVars := BindVars{
		"claim": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s
                                 FOR a IN %s
                                   FILTER obj._to == a._id
                                      AND obj._from == @claim
                                   %s
                                   SORT a.start ASC
                                   RETURN a`,
		Inference{}.CollectionName(),
		Argument{}.CollectionName(),
		c.DateFilter(bindVars),
	)
	err := FindArangoObjects(ctx, query, bindVars, &args)
	return args, err
}

func (c Claim) ArgumentsBasedOnThisClaim(ctx *ServerContext) ([]Argument, Error) {
	args := []Argument{}

	bindVars := BindVars{
		"claim": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s
                                 FOR a IN %s
                                   FILTER obj._from == a._id
                                      AND obj._to == @claim
                                   %s
                                   RETURN a`,
		BaseClaimEdge{}.CollectionName(),
		Argument{}.CollectionName(),
		c.DateFilter(bindVars),
	)
	err := FindArangoObjects(ctx, query, bindVars, &args)
	return args, err
}

func (c Claim) Inferences(ctx *ServerContext) ([]Inference, Error) {
	edges := []Inference{}
	bindVars := BindVars{
		"from": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._from == @from
                                %s
                                RETURN obj`,
		Inference{}.CollectionName(),
		c.DateFilter(bindVars))
	err := FindArangoObjects(ctx, query, bindVars, &edges)
	return edges, err
}

// Premises

func (c *Claim) AddPremise(ctx *ServerContext, premise *Claim) Error {
	if premise == nil {
		ctx.Rollback()
		return NewServerError("Premise is nil")
	}

	// TODO: Test
	if !c.MultiPremise {
		ctx.Rollback()
		return NewBusinessError("You must convert this claim to be a multi-premise claim before adding new premises")
	}

	c.QueryAt = nil
	updates := Updates{}

	if err := c.ValidateForUpdate(updates); err != nil {
		ctx.Rollback()
		return err
	}

	can, err := c.UserCanUpdate(ctx, updates)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to modify this item")
	}

	// Check for premise loops
	if premise.ID == c.ID {
		ctx.Rollback()
		return NewBusinessError("A claim cannot be a premise of itself, nor one of its own premises. That's called \"Begging the Question\"")
	}

	hasPremise, err := premise.HasPremise(ctx, c.ArangoKey())
	if err != nil {
		ctx.Rollback()
		return err
	}
	if hasPremise {
		ctx.Rollback()
		return NewBusinessError("A claim cannot be a premise of itself, nor one of its own premises. That's called \"Begging the Question\"")
	}

	hasPremise, err = c.HasPremise(ctx, premise.ArangoKey())
	if err != nil {
		ctx.Rollback()
		return err
	}
	if hasPremise {
		ctx.Rollback()
		return NewBusinessError("This claim has already been added as a premise")
	}

	if premise.Key == "" {
		if err := premise.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// TODO: locking...
	max, err := c.NumberOfPremises(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}

	edge := PremiseEdge{
		Edge: Edge{
			From: c.ArangoID(),
			To:   premise.ArangoID(),
		},
		Order: int(max) + 1,
	}

	if err := edge.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}
	return nil
}

// TODO: Test
func (c *Claim) RemovePremise(ctx *ServerContext, premiseId string) Error {
	updates := Updates{}
	if err := c.ValidateForUpdate(updates); err != nil {
		return err
	}

	can, err := c.UserCanUpdate(ctx, updates)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to modify this item")
	}

	if !c.MultiPremise {
		return NewBusinessError("You cannot remove a premise from a Claim that isn't multi-premise")
	}

	premiseEdges, err := c.PremiseEdges(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}

	if !strings.HasPrefix(premiseId, c.CollectionName()) {
		// The ID is a generic ID, not an ArangoID
		premise := Claim{}
		premise.ID = premiseId
		if err := premise.Load(ctx); err != nil {
			ctx.Rollback()
			return err
		}
		premiseId = premise.ArangoID()
	}

	var removed bool
	for _, edge := range premiseEdges {
		if edge.To == premiseId {
			if err := edge.Delete(ctx); err != nil {
				ctx.Rollback()
				return err
			}
			removed = true
		} else if removed {
			if err := edge.UpdateOrder(ctx, edge.Order-1); err != nil {
				ctx.Rollback()
				return err
			}
		}
	}

	// TODO: make sure that a new version DOESN'T have the deleted edge
	if len(premiseEdges) == 1 {
		updates := Updates{
			"mp":     false,
			"mprule": PREMISE_RULE_NONE,
		}

		if err := c.Update(ctx, updates); err != nil {
			ctx.Rollback()
			return err
		}
		c.MultiPremise = false
		c.PremiseRule = PREMISE_RULE_NONE
	}

	return nil
}

func (c Claim) HasPremise(ctx *ServerContext, premiseArangoKey string) (bool, Error) {
	db := ctx.Arango.DB

	premise := Claim{}
	premise.Key = premiseArangoKey

	qctx := arango.WithQueryCount(ctx.Context)
	bindVars := BindVars{
		"rootc":   c.ArangoID(),
		"targetc": premise.ArangoID(),
	}
	query := `FOR v IN 1..5 OUTBOUND @rootc premises
                              FILTER v._id == @targetc
                              FILTER v.end == null
                              RETURN v._key`
	//                               PRUNE v._id == @targetc
	cursor, err := db.Query(qctx, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return false, NewServerError(err.Error())
	}
	n := cursor.Count()
	return n > 0, nil
}

func (c Claim) Premises(ctx *ServerContext) ([]Claim, Error) {
	premises := []Claim{}

	if c.MultiPremise {
		bindVars := BindVars{
			"claim": c.ArangoID(),
		}
		query := fmt.Sprintf(`FOR obj IN %s
                                 FOR c IN %s
                                   FILTER obj._to == c._id
                                      AND obj._from == @claim
                                   %s
                                   SORT obj.order
                                   RETURN c`,
			PremiseEdge{}.CollectionName(),
			Claim{}.CollectionName(),
			c.DateFilter(bindVars),
		)
		if err := FindArangoObjects(ctx, query, bindVars, &premises); err != nil {
			return premises, err
		}
	}

	return premises, nil
}

func (c Claim) ReorderPremise(ctx *ServerContext, premise Claim, new int) ([]Claim, Error) {
	premises := []Claim{}

	// TODO: test
	updates := Updates{}
	if err := c.ValidateForUpdate(updates); err != nil {
		return premises, err
	}

	// TODO: Test
	can, err := c.UserCanUpdate(ctx, updates)
	if err != nil {
		return premises, err
	}
	if !can {
		return premises, NewPermissionError("You do not have permission to modify this item")
	}

	if new <= 0 {
		return premises, NewBusinessError("Order: invalid value;")
	}

	num, err := c.NumberOfPremises(ctx)
	if err != nil {
		return premises, err
	}
	if new > int(num) {
		return premises, NewBusinessError("Order: the new order is higher than the number of premises;")
	}

	edges, err := c.PremiseEdges(ctx)
	if err != nil {
		return premises, err
	}

	var old int
	for i, edge := range edges {
		if edge.To == premise.ArangoID() {
			old = i + 1
		}
	}

	if old == 0 {
		return premises, NewNotFoundError("The premise you are trying to reorder was not found")
	}

	min := support.MinInt(new, old)
	max := support.MaxInt(new, old)
	var window bool
	for i, edge := range edges {
		curr := i + 1
		if curr >= min && curr <= max {
			window = true
		} else {
			window = false
		}
		if window {
			if curr == old {
				edge.Order = new
			} else if new < old {
				edge.Order = curr + 1
			} else {
				edge.Order = curr - 1
			}
		} else {
			edge.Order = curr
		}
		if err := edge.UpdateOrder(ctx, edge.Order); err != nil {
			ctx.Rollback()
			return premises, err
		}
	}

	premises, err = c.Premises(ctx)
	if err != nil {
		ctx.Rollback()
		return premises, err
	}
	return premises, nil
}

func (c Claim) PremiseEdges(ctx *ServerContext) ([]PremiseEdge, Error) {
	edges := []PremiseEdge{}

	bindVars := BindVars{
		"from": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._from == @from
                                %s
                                SORT obj.order
                                RETURN obj`,
		PremiseEdge{}.CollectionName(),
		c.DateFilter(bindVars))
	err := FindArangoObjects(ctx, query, bindVars, &edges)
	return edges, err
}

func (c Claim) NumberOfPremises(ctx *ServerContext) (int64, Error) {
	db := ctx.Arango.DB

	var n int64
	if c.MultiPremise {
		qctx := arango.WithQueryCount(ctx.Context)

		bindVars := BindVars{
			"from": c.ArangoID(),
		}
		query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._from == @from
                                %s
                                SORT obj.order
                                RETURN obj`,
			PremiseEdge{}.CollectionName(),
			c.DateFilter(bindVars))
		cursor, err := db.Query(qctx, query, bindVars)
		defer CloseCursor(cursor)
		if err != nil {
			return n, NewServerError(err.Error())
		}
		n = cursor.Count()
	}
	return n, nil
}

func (c Claim) EdgesToThisPremise(ctx *ServerContext) ([]PremiseEdge, Error) {
	edges := []PremiseEdge{}

	bindVars := BindVars{
		"to": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._to == @to 
                                %s
                                RETURN obj`,
		PremiseEdge{}.CollectionName(),
		c.DateFilter(bindVars))
	err := FindArangoObjects(ctx, query, bindVars, &edges)
	return edges, err
}

// Arguments that use this Claim

func (c Claim) BaseClaimEdges(ctx *ServerContext) ([]BaseClaimEdge, Error) {
	edges := []BaseClaimEdge{}

	bindVars := BindVars{
		"to": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._to == @to
                                %s
                                RETURN obj`,
		BaseClaimEdge{}.CollectionName(),
		c.DateFilter(bindVars))
	err := FindArangoObjects(ctx, query, bindVars, &edges)
	return edges, err
}

// Contexts
func (c *Claim) AddContext(ctx *ServerContext, context Context) Error {
	c.QueryAt = nil
	updates := Updates{}
	if err := c.ValidateForUpdate(updates); err != nil {
		ctx.Rollback()
		return err
	}

	can, err := c.UserCanUpdate(ctx, updates)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to modify this item")
	}

	if c.MultiPremise {
		ctx.Rollback()
		return NewBusinessError("Multi-premise claims inherit the union of contexts from all their premises")
	}

	// Check for duplicates
	contexts, err := c.Contexts(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, con := range contexts {
		if con.Key == context.Key {
			ctx.Rollback()
			return NewBusinessError("This context was already added to this claim")
		}
	}

	edge := ContextEdge{Edge: Edge{
		From: context.ArangoID(),
		To:   c.ArangoID(),
	}}

	if err := edge.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}
	return nil
}

func (c *Claim) RemoveContext(ctx *ServerContext, contextArangoKey string) Error {
	updates := Updates{}
	if err := c.ValidateForUpdate(updates); err != nil {
		return err
	}

	can, err := c.UserCanUpdate(ctx, updates)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to modify this item")
	}

	edge, err := FindContextEdge(ctx, contextArangoKey, c.ArangoKey())
	if err != nil {
		ctx.Rollback()
		return err
	}

	if err := edge.Delete(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	return nil
}

func (c Claim) Contexts(ctx *ServerContext) ([]Context, Error) {
	contexts := []Context{}

	if c.MultiPremise {
		// TODO: Do this in a single query
		m := map[string]Context{}
		premises, err := c.Premises(ctx)
		if err != nil {
			return contexts, err
		}
		for _, premise := range premises {
			premise.QueryAt = c.QueryAt
			pctx, err := premise.Contexts(ctx)
			if err != nil {
				return contexts, err
			}
			for _, context := range pctx {
				if _, ok := m[context.Key]; !ok {
					contexts = append(contexts, context)
					m[context.Key] = context
				}
			}
			sort.Slice(contexts, func(i, j int) bool {
				return contexts[i].ShortName < contexts[j].ShortName
			})
		}
	} else {
		bindVars := BindVars{
			"claim": c.ArangoID(),
		}
		query := fmt.Sprintf(`FOR obj IN %s
                                 FOR c IN %s
                                   FILTER obj._to == @claim
                                      AND obj._from == c._id
                                   %s
                                   SORT c.name
                                   RETURN c`,
			ContextEdge{}.CollectionName(),
			Context{}.CollectionName(),
			c.DateFilter(bindVars),
		)
		if err := FindArangoObjects(ctx, query, bindVars, &contexts); err != nil {
			return contexts, err
		}
	}

	return contexts, nil
}

func (c Claim) ContextEdges(ctx *ServerContext) ([]ContextEdge, Error) {
	edges := []ContextEdge{}

	bindVars := BindVars{
		"to": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._to == @to 
                                %s
                                SORT obj.start
                                RETURN obj`,
		ContextEdge{}.CollectionName(),
		c.DateFilter(bindVars))
	err := FindArangoObjects(ctx, query, bindVars, &edges)
	return edges, err
}

// Curation

// TODO: Test
func (c *Claim) ConvertToMultiPremise(ctx *ServerContext) Error {
	if c.MultiPremise {
		ctx.Rollback()
		return NewBusinessError("This claim is already a multi-premise claim")
	}
	c.QueryAt = nil

	updates := Updates{
		"mp":     true,
		"mprule": PREMISE_RULE_ALL,
	}

	// Make the current claim an MP claim (preserve the ID)
	if err := c.Update(ctx, updates); err != nil {
		ctx.Rollback()
		return err
	}
	if err := c.Load(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	// Create a new Claim with the old attributes as the first premise of this claim
	premise := Claim{}
	premise.Title = c.Title
	premise.Negation = c.Negation
	premise.Question = c.Question
	premise.Description = c.Description
	premise.Note = c.Note
	premise.Image = c.Image
	if err := premise.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}

	if err := c.AddPremise(ctx, &premise); err != nil {
		ctx.Rollback()
		return err
	}

	// Move arguments to current MP Claim down to new premise
	args, err := c.Arguments(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, arg := range args {
		if err := arg.MoveTo(ctx, &premise, arg.Pro); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// Move Contexts to current MP Claim down to new premise
	contextEdges, err := c.ContextEdges(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range contextEdges {
		newEdge := ContextEdge{Edge: Edge{
			From: edge.From,
			To:   premise.ArangoID(),
		}}
		if err := newEdge.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
		if err := edge.Delete(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	return nil
}

// TODO: Implement merge
// TODO: Implement search

// Scorer

func (c *Claim) Score(ctx *ServerContext) (float32, Error) {
	if c.QueryAt == nil {
		return c.Truth, nil
	}
	return c.scoreAt(ctx)
}

func (c *Claim) UpdateScore(ctx *ServerContext) Error {
	// TODO: not on deleted - validate for update
	c.QueryAt = nil
	score, err := c.scoreAt(ctx)
	if err != nil {
		return err
	}

	updates := Updates{
		"truth": score,
	}

	col, grr := ctx.Arango.CollectionFor(c)
	if grr != nil {
		return grr
	}
	if _, err := col.UpdateDocument(ctx.Context, c.ArangoKey(), updates); err != nil {
		return NewServerError(err.Error())
	}

	c.Truth = score
	return nil
}

func (c *Claim) scoreAt(ctx *ServerContext) (float32, Error) {
	var score float32
	results := map[string]interface{}{}

	bindVars := BindVars{
		"claim": c.ID,
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                 FOR c IN %s
                                   FILTER obj._to == c._id
                                      AND c.id == @claim
                                   %s
                                   COLLECT
                                   AGGREGATE 
                                     num = COUNT(obj),
                                     score = AVG(obj.score)
                                   RETURN { num, score }`,
		UserScore{}.CollectionName(),
		c.CollectionName(),
		c.DateFilter(bindVars))

	db := ctx.Arango.DB
	cursor, err := db.Query(ctx.Context, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return score, NewServerError(err.Error())
	}
	_, err = cursor.ReadDocument(ctx.Context, &results)
	if err != nil {
		return score, NewServerError(err.Error())
	}

	if val, ok := results["score"].(float64); ok {
		score = float32(val)
	}
	if score == 0.0 {
		if count, ok := results["num"].(float64); ok {
			if count == 0 {
				score = DEFAULT_CLAIM_SCORE
			}
		}
	}

	return score, nil
}

func (c Claim) UserScores(ctx *ServerContext) ([]UserScore, Error) {
	edges := []UserScore{}

	bindVars := BindVars{
		"to": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._to == @to 
                                %s
                                SORT obj.start
                                RETURN obj`,
		UserScore{}.CollectionName(),
		c.DateFilter(bindVars))
	err := FindArangoObjects(ctx, query, bindVars, &edges)
	return edges, err
}

// Graph methods

// TODO: THis could use the named graph debate_map
func (c Claim) HasCycle(ctx *ServerContext) (bool, Error) {
	db := ctx.Arango.DB

	qctx := arango.WithQueryCount(ctx.Context)
	bindVars := BindVars{
		"claim": c.ArangoID(),
	}
	// TODO: Add PRUNE?
	query := fmt.Sprintf(`FOR v, obj IN 1..10 OUTBOUND @claim 
                                inferences, base_claims, premises
                              FILTER v._id == @claim
                              %s
                              RETURN v._key`,
		c.DateFilter(bindVars))
	//                               PRUNE v._id == @targetc
	cursor, err := db.Query(qctx, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return false, NewServerError(err.Error())
	}
	n := cursor.Count()
	return n > 0, nil
}

// Queries

// TODO: Obviously, this is going to have to be denormalized at some point
func (c Claim) QueryForTopLevelClaims(params ArangoQueryParameters) string {
	params = c.DefaultQueryParameters().Merge(params)
	query := `FOR obj IN claims 
                    LET bcCount=(FOR bc IN base_claims 
                                   FILTER bc._to == obj._id 
                                      AND bc.end == null 
                                  COLLECT WITH COUNT INTO length 
                                   RETURN length) 
                    FILTER bcCount[0] == 0 
                    LET pCount=(FOR p IN premises
                                   FILTER p._to == obj._id 
                                      AND p.end == null 
                                  COLLECT WITH COUNT INTO length 
                                   RETURN length) 
                    FILTER bcCount[0] == 0 
                    FILTER pCount[0] == 0 
                    AND obj.end == null`
	return params.Apply(query)
}
