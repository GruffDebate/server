package gruff

import (
	"fmt"
	"strings"

	"github.com/GruffDebate/server/support"
	arango "github.com/arangodb/go-driver"
)

/*
 * Forgive us the confusion between ServerContext and Context!
 *
 * Context is a critical element to the debate graph, in that it removes ambiguity
 * from the Claims as much as can be possible.
 *
 * Each Context element represents an object or concept in the real world,
 * preferably one which has been represented within an external knowledge graph.
 * By linking Context elements to a Claim, it turns inexact statements into very
 * specific Claims.
 *
 * For example, the statement "Martin Luther was responsible for the revolution"
 * could potentially be referring to the German priest that began the 16th century Reformation,
 * or could be referring to the 20th century American minister who was a leader in the U.S. civil rights movement.
 * By attaching Context elements (akin to wiki pages) to the Claims, it can be very clear
 * which one is meant without long descriptions in the title of the Claim.
 *
 * By linking Contexts to knowledge graphs, it also becomes possible to perform enhanced
 * searches based on graph relationships (show me "other revolutions", "other leaders of the civil rights movement", "other ministers").
 * It also enables better automated de-duplication algorithms, since Claims with very
 * different titles, but the same or similar Contexts may be compared for duplication.
 */

type Context struct {
	Model
	ShortName        string    `json:"name" valid:"length(1|60),required"`
	Title            string    `json:"title" sql:"not null" valid:"length(3|1000),required"`
	Description      string    `json:"desc" valid:"length(3|4000)"`
	URL              string    `json:"url" valid:"url,required"`
	MID              string    `json:"mid,omitempty"` // Google KG ID
	QID              string    `json:"qid,omitempty"` // Wikidata ID
	MetaDataURL      *MetaData `json:"meta_url,omitempty"`
	MetaDataGoogle   *MetaData `json:"meta_google,omitempty"`
	MetaDataWikidata *MetaData `json:"meta_wikidata,omitempty"`
	// TODO: add other KBs
}

type MetaData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
}

// ArangoObject interface

func (c Context) CollectionName() string {
	return "contexts"
}

func (c Context) ArangoKey() string {
	return c.Key
}

func (c Context) ArangoID() string {
	return fmt.Sprintf("%s/%s", c.CollectionName(), c.ArangoKey())
}

func (c Context) DefaultQueryParameters() ArangoQueryParameters {
	params := ArangoQueryParameters{
		Sort: support.StringPtr("obj.name ASC"),
	}

	return DEFAULT_QUERY_PARAMETERS.Merge(params)
}

func (c *Context) Create(ctx *ServerContext) Error {
	// TODO: Unique indexes? Unique checks?
	return CreateArangoObject(ctx, c)
}

// TODO: Test
// TODO: generic...
func (c *Context) Update(ctx *ServerContext, updates Updates) Error {
	if err := c.ValidateForUpdate(updates); err != nil {
		return err
	}

	// TODO: Test
	can, err := c.UserCanUpdate(ctx, updates)
	if err != nil {
		return err
	}
	if !can {
		return NewPermissionError("You do not have permission to modify this item")
	}

	col, err := ctx.Arango.CollectionFor(c)
	if err != nil {
		return err

	}

	if _, err := col.UpdateDocument(ctx.Context, c.ArangoKey(), updates); err != nil {
		return NewServerError(err.Error())
	}

	return nil
}

// TODO: Test
func (c *Context) Delete(ctx *ServerContext) Error {
	n, err := c.NumberOfClaims(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return NewBusinessError("A context cannot be deleted if it's used by any claims")
	}

	return DeleteArangoObject(ctx, c)
}

// Restrictor
// TODO: Test
// TODO: Call in CRUD and other methods
func (c Context) UserCanView(ctx *ServerContext) (bool, Error) {
	return true, nil
}

func (c Context) UserCanCreate(ctx *ServerContext) (bool, Error) {
	return ctx.UserLoggedIn(), nil
}

func (c Context) UserCanUpdate(ctx *ServerContext, updates Updates) (bool, Error) {
	return c.UserCanDelete(ctx)
}

func (c Context) UserCanDelete(ctx *ServerContext) (bool, Error) {
	u := ctx.UserContext
	return u.Curator, nil
}

// Validator

func (c Context) ValidateForCreate() Error {
	return ValidateStruct(c)
}

func (c Context) ValidateForUpdate(updates Updates) Error {
	if err := SetJsonValuesOnStruct(&c, updates, false); err != nil {
		return err
	}
	return c.ValidateForCreate()
}

func (c Context) ValidateForDelete() Error {
	return nil
}

func (c Context) ValidateField(f string) Error {
	return ValidateStructField(c, f)
}

// Business methods

// TODO: TEst
func (c Context) NumberOfClaims(ctx *ServerContext) (int64, Error) {
	db := ctx.Arango.DB

	qctx := arango.WithQueryCount(ctx.Context)

	var n int64
	bindVars := BindVars{
		"from": c.ArangoID(),
	}
	query := fmt.Sprintf(`FOR obj IN %s 
                                FILTER obj._from == @from
                                RETURN obj`,
		ContextEdge{}.CollectionName())
	// TODO: make generic arango count items method
	cursor, err := db.Query(qctx, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return n, NewServerError(err.Error())
	}
	n = cursor.Count()

	return n, nil
}

func FindContext(ctx *ServerContext, contextArangoId string) (Context, Error) {
	context := Context{}
	bindVars := BindVars{
		"context": contextArangoId,
	}
	query := fmt.Sprintf(`FOR obj IN %s
                                      FILTER obj._id == @context
                                       LIMIT 1
                                      RETURN obj`,
		Context{}.CollectionName(),
	)
	err := FindArangoObject(ctx, query, bindVars, &context)
	return context, err
}

func SearchContext(ctx *ServerContext, contextText string) ([]Context, Error) {
	contexts := []Context{}
	bindVars := BindVars{
		"name":  "%" + contextText + "%",
		"title": "%" + contextText + "%",
	}
	query := fmt.Sprintf(`
		FOR obj IN %s
		FILTER obj.name like @name OR obj.title like @title
		RETURN obj
	`, Context{}.CollectionName())

	err := FindArangoObjects(ctx, query, bindVars, &contexts)
	return contexts, err
}

func (c *Context) Load(ctx *ServerContext) Error {
	var err Error
	if c.ArangoKey() != "" {
		err = LoadArangoObject(ctx, c, c.ArangoKey())
	} else {
		err = NewBusinessError("There is no key or id for this Context")
	}

	return err
}

func (c *Context) LoadFull(ctx *ServerContext) Error {
	if err := c.Load(ctx); err != nil {
		return err
	}
	return nil
}

func SearchContexts(ctx *ServerContext, term string) ([]Context, Error) {
	contexts := []Context{}
	if term == "" {
		return contexts, NewBusinessError("Query term required")
	}

	term = strings.ToLower(term)
	term = strings.ReplaceAll(term, "%", "\\%")
	bindVars := BindVars{
		"term": "%" + term + "%",
	}
	query := fmt.Sprintf(`FOR obj IN %s
                                      FILTER lower(obj.title) LIKE @term
                                          OR lower(obj.name) LIKE @term
                                      RETURN obj`,
		Context{}.CollectionName(),
	)
	err := FindArangoObjects(ctx, query, bindVars, &contexts)
	return contexts, err
}
