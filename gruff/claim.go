package gruff

import (
	"fmt"
	"time"

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

type Claim struct {
	Identifier
	Title        string     `json:"title" valid:"length(3|1000)"`
	Negation     string     `json:"negation"`
	Question     string     `json:"question"`
	Description  string     `json:"desc" valid:"length(3|4000)"`
	Note         string     `json:"note"`
	Image        string     `json:"img,omitempty"`
	MultiPremise bool       `json:"mp"`
	PremiseRule  int        `json:"mprule"`
	Truth        float64    `json:"truth"`   // Average score from direct opinions
	TruthRU      float64    `json:"truthRU"` // Average score rolled up from argument totals
	ProArgs      []Argument `json:"proargs"`
	ConArgs      []Argument `json:"conargs"`
	Links        []Link     `json:"links,omitempty"`
	Contexts     []Context  `json:"contexts,omitempty"`
	ContextIDs   []uint64   `json:"contextIds,omitempty"`
	Tags         []Tag      `json:"tags,omitempty"`
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

// Validator

func (c Claim) ValidateForCreate() GruffError {
	return ValidateStruct(c)
}

func (c Claim) ValidateForUpdate() GruffError {
	return c.ValidateForCreate()
}

func (c Claim) ValidateField(f string) GruffError {
	return ValidateStructField(c, f)
}

// Creator

func (c *Claim) Create(ctx *ServerContext) GruffError {
	col, err := ctx.Arango.CollectionFor(c)
	if err != nil {
		return err

	}

	// TODO: validate for create
	c.PrepareForCreate(ctx.UserContext)

	if _, dberr := col.CreateDocument(ctx.Context, c); dberr != nil {
		return NewServerError(dberr.Error())
	}
	return nil
}

// Loader

// If the Claim object has a key, that exact Claim will be loaded
// Otherwise, Load will look for Claims matching the ID
// If CreatedAt is a non-blank value, it will load the Claim active at that time (if any)
// Otherwise, it will return the current active (undeleted) version.
func (c *Claim) Load(ctx *ServerContext) GruffError {
	db := ctx.Arango.DB

	col, err := ctx.Arango.CollectionFor(c)
	if err != nil {
		return err
	}

	if c.ArangoKey() != "" {
		_, dberr := col.ReadDocument(ctx.Context, c.ArangoKey(), c)
		if dberr != nil {
			return NewServerError(dberr.Error())
		}
	} else if c.ID != "" {
		var empty time.Time
		var query string
		bindVars := map[string]interface{}{
			"id": c.ID,
		}
		if c.CreatedAt == empty {
			query = fmt.Sprintf("FOR c IN %s FILTER c.id == @id AND c.end == null SORT c.start DESC LIMIT 1 RETURN c", c.CollectionName())
		} else {
			bindVars["start"] = c.CreatedAt
			query = fmt.Sprintf("FOR c IN %s FILTER c.id == @id AND c.start <= @start AND (c.end == null OR c.end > @start) SORT c.start DESC LIMIT 1 RETURN c", c.CollectionName())
		}
		cursor, err := db.Query(ctx.Context, query, bindVars)
		if err != nil {
			return NewServerError(err.Error())
		}
		defer cursor.Close()
		for cursor.HasMore() {
			_, err := cursor.ReadDocument(ctx.Context, c)
			if err != nil {
				return NewServerError(err.Error())
			}
		}
	} else {
		return NewBusinessError("There is no key or id for this Claim.")
	}

	return nil
}

// Versioner

func (c *Claim) Version(ctx *ServerContext) (Claim, GruffError) {
	var newVersion Claim

	oldVersion := Claim{}
	oldVersion.ID = c.ID
	err := oldVersion.Load(ctx)
	if err != nil {
		return newVersion, err
	}

	// This should delete all the old edges, too
	if err := oldVersion.Delete(ctx); err != nil {
		ctx.Rollback()
		return newVersion, err
	}

	c.PrepareForCreate(ctx.UserContext)
	if err := c.Create(ctx); err != nil {
		ctx.Rollback()
		return newVersion, err
	}

	// Find all edges going to old ver, make copy to new ver
	// TODO: The method to get edges needs to only return the undeleted edges... but they were already deleted...
	if c.MultiPremise {
		premiseEdges, err := oldVersion.PremiseEdges(ctx)
		if err != nil {
			ctx.Rollback()
			return newVersion, err
		}
		for _, edge := range premiseEdges {
			newEdge := PremiseEdge{
				From:  c.ArangoID(),
				To:    edge.To,
				Order: edge.Order,
			}
			if err := newEdge.Create(ctx); err != nil {
				ctx.Rollback()
				return newVersion, err
			}
		}
	}

	// Arguments
	inferences, err := oldVersion.Inferences(ctx)
	if err != nil {
		ctx.Rollback()
		return newVersion, err
	}
	for _, edge := range inferences {
		newEdge := Inference{
			From: c.ArangoID(),
			To:   edge.To,
		}
		if err := newEdge.Create(ctx); err != nil {
			ctx.Rollback()
			return newVersion, err
		}
	}

	// Base Claim edges
	baseClaimEdges, err := oldVersion.BaseClaimEdges(ctx)
	if err != nil {
		ctx.Rollback()
		return newVersion, err
	}
	for _, edge := range baseClaimEdges {
		newEdge := BaseClaimEdge{
			From: edge.From,
			To:   c.ArangoID(),
		}
		if err := newEdge.Create(ctx); err != nil {
			ctx.Rollback()
			return newVersion, err
		}
	}

	// TODO: Contexts
	// TODO: References
	// TODO: Tags

	newVersion = *c
	return newVersion, nil
}

func (c *Claim) Delete(ctx *ServerContext) GruffError {
	c.PrepareForDelete()
	patch := map[string]interface{}{
		"end": c.DeletedAt,
	}
	col, err := ctx.Arango.CollectionFor(c)
	if err != nil {
		return err
	}
	_, dberr := col.UpdateDocument(ctx.Context, c.ArangoKey(), patch)
	if dberr != nil {
		return NewServerError(dberr.Error())
	}

	// Delete any edges to or from this Claim
	// TODO: How to handle deleting a Claim that is used as a BaseClaim
	// Note that Delete is also used when versioning
	// TODO: This would probably be faster just to execute a singe update query per edge type

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

	// Arguments
	inferences, err := c.Inferences(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range inferences {
		if err := edge.Delete(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// Base Claim edges
	baseClaimEdges, err := c.BaseClaimEdges(ctx)
	if err != nil {
		ctx.Rollback()
		return err
	}
	for _, edge := range baseClaimEdges {
		if err := edge.Delete(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	// TODO: Contexts
	// TODO: References
	// TODO: Tags

	return nil
}

// Business methods

func (c Claim) AddArgument(ctx *ServerContext, a Argument) GruffError {
	edge := Inference{
		From: c.ArangoID(),
		To:   a.ArangoID(),
	}

	if err := edge.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}
	return nil
}

func (c *Claim) AddPremise(ctx *ServerContext, premise *Claim) GruffError {
	if premise == nil {
		ctx.Rollback()
		return NewServerError("Premise is nil")
	}

	if premise.Key == "" {
		if err := premise.Create(ctx); err != nil {
			ctx.Rollback()
			return err
		}
	}

	if !c.MultiPremise {
		c.MultiPremise = true
		c.PremiseRule = PREMISE_RULE_ALL

		if _, err := c.Version(ctx); err != nil {
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
		From:  c.ArangoID(),
		To:    premise.ArangoID(),
		Order: int(max) + 1,
	}

	if err := edge.Create(ctx); err != nil {
		ctx.Rollback()
		return err
	}
	return nil
}

func (c Claim) Arguments(ctx *ServerContext) ([]Argument, GruffError) {
	db := ctx.Arango.DB
	args := []Argument{}

	// TODO: not deleted, or deleted > c.Deleted
	query := fmt.Sprintf(`FOR i IN %s
                                 FOR a IN %s
                                   FILTER i._to == a._id
                                      AND i._from == @claim
                                   RETURN a`,
		Inference{}.CollectionName(),
		Argument{}.CollectionName(),
	)
	bindVars := map[string]interface{}{
		"claim": c.ArangoID(),
	}
	cursor, err := db.Query(ctx.Context, query, bindVars)
	if err != nil {
		return args, NewServerError(err.Error())
	}
	defer cursor.Close()
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

func (c Claim) Premises(ctx *ServerContext) ([]Claim, GruffError) {
	db := ctx.Arango.DB
	premises := []Claim{}

	// TODO: not deleted, or deleted > c.Deleted
	query := fmt.Sprintf(`FOR p IN %s
                                 FOR c IN %s
                                   FILTER p._to == c._id
                                      AND p._from == @claim
                                   SORT p.order
                                   RETURN c`,
		PremiseEdge{}.CollectionName(),
		Claim{}.CollectionName(),
	)
	bindVars := map[string]interface{}{
		"claim": c.ArangoID(),
	}
	cursor, err := db.Query(ctx.Context, query, bindVars)
	if err != nil {
		return premises, NewServerError(err.Error())
	}
	defer cursor.Close()
	for cursor.HasMore() {
		claim := Claim{}
		_, err := cursor.ReadDocument(ctx.Context, &claim)
		if err != nil {
			return premises, NewServerError(err.Error())
		}
		premises = append(premises, claim)
	}

	return premises, nil
}

func (c Claim) ReorderPremise(ctx *ServerContext, premise Claim, new int) ([]Claim, GruffError) {
	premises := []Claim{}

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
		return premises, NewNotFoundError("The premise you are trying to reorder was not found.")
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

// Edges

func (c Claim) PremiseEdges(ctx *ServerContext) ([]PremiseEdge, GruffError) {
	db := ctx.Arango.DB
	edges := []PremiseEdge{}

	// TODO: order by, not deleted, or deleted > c.Deleted
	query := fmt.Sprintf("FOR e IN %s FILTER e._from == @from SORT e.order RETURN e", PremiseEdge{}.CollectionName())
	bindVars := map[string]interface{}{
		"from": c.ArangoID(),
	}
	cursor, err := db.Query(ctx.Context, query, bindVars)
	if err != nil {
		return edges, NewServerError(err.Error())
	}
	defer cursor.Close()
	for cursor.HasMore() {
		edge := PremiseEdge{}
		_, err := cursor.ReadDocument(ctx.Context, &edge)
		if err != nil {
			return edges, NewServerError(err.Error())
		}
		edges = append(edges, edge)
	}

	return edges, nil
}

func (c Claim) NumberOfPremises(ctx *ServerContext) (int64, GruffError) {
	db := ctx.Arango.DB
	qctx := arango.WithQueryCount(ctx.Context)
	var n int64

	// TODO: not deleted, or deleted > c.Deleted
	query := fmt.Sprintf("FOR e IN %s FILTER e._from == @from RETURN e", PremiseEdge{}.CollectionName())
	bindVars := map[string]interface{}{
		"from": c.ArangoID(),
	}
	cursor, err := db.Query(qctx, query, bindVars)
	if err != nil {
		return n, NewServerError(err.Error())
	}
	defer cursor.Close()
	n = cursor.Count()
	return n, nil
}

// TODO: this could most definitely be made more generic...
func (c Claim) Inferences(ctx *ServerContext) ([]Inference, GruffError) {
	db := ctx.Arango.DB
	edges := []Inference{}

	// TODO: not deleted, or deleted > c.Deleted
	query := fmt.Sprintf("FOR e IN %s FILTER e._from == @from RETURN e", Inference{}.CollectionName())
	bindVars := map[string]interface{}{
		"from": c.ArangoID(),
	}
	cursor, err := db.Query(ctx.Context, query, bindVars)
	if err != nil {
		return edges, NewServerError(err.Error())
	}
	defer cursor.Close()
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

// TODO: this could most definitely be made more generic...
func (c Claim) BaseClaimEdges(ctx *ServerContext) ([]BaseClaimEdge, GruffError) {
	db := ctx.Arango.DB
	edges := []BaseClaimEdge{}

	// TODO: not deleted, or deleted > c.Deleted
	query := fmt.Sprintf("FOR e IN %s FILTER e._to == @to RETURN e", BaseClaimEdge{}.CollectionName())
	bindVars := map[string]interface{}{
		"to": c.ArangoID(),
	}
	cursor, err := db.Query(ctx.Context, query, bindVars)
	if err != nil {
		return edges, NewServerError(err.Error())
	}
	defer cursor.Close()
	for cursor.HasMore() {
		edge := BaseClaimEdge{}
		_, err := cursor.ReadDocument(ctx.Context, &edge)
		if err != nil {
			return edges, NewServerError(err.Error())
		}
		edges = append(edges, edge)
	}

	return edges, nil
}

// TODO: Create method should set default Truth to 0.5
// TODO: Implement merge
// TODO: Implement search

func (c Claim) UpdateTruth(ctx *ServerContext) {
	//ctx.Database.Exec("UPDATE claims c SET truth = (SELECT AVG(truth) FROM claim_opinions WHERE claim_id = c.id) WHERE id = ?", c.ID)

	// TODO: test
	if c.TruthRU == 0.0 {
		// There's no roll up score yet, so the truth score itself is affecting related roll ups
		//c.UpdateAncestorRUs(ctx)
	}
}

/*
func (c *Claim) UpdateTruthRU(ctx *ServerContext) {
	// TODO: do it all in SQL?
	// TODO: should updates be recursive? (first, calculate sub-argument RUs)
	//       or, should it trigger an update of anyone that references it?
	proArgs, conArgs := c.Arguments(ctx)

	if len(proArgs) > 0 || len(conArgs) > 0 {
		proScore := 0.0
		for _, arg := range proArgs {
			remainder := 1.0 - proScore
			score := 0 //arg.ScoreRU(ctx)
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

		c.TruthRU = netScore
	} else {
		c.TruthRU = 0.0
	}

	//ctx.Database.Set("gorm:save_associations", false).Save(c)

	c.UpdateAncestorRUs(ctx)
}

func (c Claim) UpdateAncestorRUs(ctx *ServerContext) {
	args := []Argument{}
	ctx.Database.Where("claim_id = ?", c.ID).Find(&args)
	for _, arg := range args {
		// TODO: instead, add to list of things to be updated in bg
		// TODO: what about cycles??
		// TODO: test
		arg.UpdateStrengthRU(ctx)
	}
}
*/

// Queries

func (c Claim) QueryForTopLevelClaims(params ArangoQueryParameters) string {
	params = c.DefaultQueryParameters().Merge(params)
	query := "FOR obj IN claims LET bcCount=(FOR bc IN base_claims FILTER bc._to == obj._id COLLECT WITH COUNT INTO length RETURN length) FILTER bcCount[0] == 0"
	return params.Apply(query)
}
