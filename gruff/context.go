package gruff

import (
	"fmt"
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
	return DEFAULT_QUERY_PARAMETERS
}

// Validator

func (c Context) ValidateForCreate() GruffError {
	return ValidateStruct(c)
}

func (c Context) ValidateForUpdate(updates map[string]interface{}) GruffError {
	if err := SetJsonValuesOnStruct(&c, updates); err != nil {
		return err
	}
	return c.ValidateForCreate()
}

func (c Context) ValidateField(f string) GruffError {
	return ValidateStructField(c, f)
}

// Creator
// TODO: Test validations, etc.
func (c *Context) Create(ctx *ServerContext) GruffError {
	if err := c.ValidateForCreate(); err != nil {
		return err
	}

	// TODO: Unique indexes? Unique checks?

	col, err := ctx.Arango.CollectionFor(c)
	if err != nil {
		return err
	}

	c.PrepareForCreate(ctx)

	if _, dberr := col.CreateDocument(ctx.Context, c); dberr != nil {
		ctx.Rollback()
		return NewServerError(dberr.Error())
	}

	return nil
}

// Business methods

func FindContext(ctx *ServerContext, contextArangoId string) (Context, GruffError) {
	db := ctx.Arango.DB

	context := Context{}
	bindVars := map[string]interface{}{
		"context": contextArangoId,
	}
	query := fmt.Sprintf(`FOR obj IN %s
                                      FILTER obj._id == @context
                                       LIMIT 1
                                      RETURN obj`,
		Context{}.CollectionName(),
	)
	cursor, err := db.Query(ctx.Context, query, bindVars)
	defer CloseCursor(cursor)
	if err != nil {
		return context, NewServerError(err.Error())
	}
	for cursor.HasMore() {
		_, err := cursor.ReadDocument(ctx.Context, &context)
		if err != nil {
			return context, NewServerError(err.Error())
		}
	}

	return context, nil
}
