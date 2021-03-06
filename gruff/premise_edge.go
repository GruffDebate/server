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

func (p *PremiseEdge) Create(ctx *ServerContext) Error {
	return CreateArangoObject(ctx, p)
}

func (p *PremiseEdge) Update(ctx *ServerContext, updates Updates) Error {
	return NewServerError("This item cannot be modified")
}

func (p *PremiseEdge) Delete(ctx *ServerContext) Error {
	return DeleteArangoObject(ctx, p)
}

// TODO: Preserve history...
func (p *PremiseEdge) UpdateOrder(ctx *ServerContext, order int) Error {
	col, err := ctx.Arango.CollectionFor(p)
	if err != nil {
		return err
	}

	patch := BindVars{"order": order}
	_, aerr := col.UpdateDocument(ctx.Context, p.ArangoKey(), patch)
	if aerr != nil {
		return NewServerError(aerr.Error())
	}
	return nil
}
