package gruff

import (
	"fmt"
)

// A PremiseEdge is an edge that goes from a Multi-premise Claim
// to one of the Claims that represents a specific premise
type PremiseEdge struct {
	Edge
	Order int `json:"order"`
}

// ArangoObject interface

func (p PremiseEdge) CollectionName() string {
	return "premises"
}

func (p PremiseEdge) ArangoKey() string {
	return p.Key
}

func (p PremiseEdge) ArangoID() string {
	return fmt.Sprintf("%s/%s", p.CollectionName(), p.ArangoKey())
}

func (p PremiseEdge) DefaultQueryParameters() ArangoQueryParameters {
	return DEFAULT_QUERY_PARAMETERS
}

func (p *PremiseEdge) Create(ctx *ServerContext) GruffError {
	col, err := ctx.Arango.CollectionFor(p)
	if err != nil {
		return err
	}

	p.PrepareForCreate(ctx)

	_, aerr := col.CreateDocument(ctx.Context, p)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}
	return nil
}

func (p *PremiseEdge) Update(ctx *ServerContext, updates map[string]interface{}) GruffError {
	return NewServerError("This item cannot be modified")
}

func (p *PremiseEdge) Delete(ctx *ServerContext) GruffError {
	p.PrepareForDelete(ctx)
	patch := map[string]interface{}{
		"end": p.DeletedAt,
	}
	col, err := ctx.Arango.CollectionFor(p)
	if err != nil {
		return err
	}
	_, aerr := col.UpdateDocument(ctx.Context, p.ArangoKey(), patch)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}

	return nil
}

// TODO: Preserve history...
func (p *PremiseEdge) UpdateOrder(ctx *ServerContext, order int) GruffError {
	col, err := ctx.Arango.CollectionFor(p)
	if err != nil {
		return err
	}

	patch := map[string]interface{}{"order": order}
	_, aerr := col.UpdateDocument(ctx.Context, p.ArangoKey(), patch)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}
	return nil
}
