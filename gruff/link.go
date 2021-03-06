package gruff

import (
	"github.com/google/uuid"
)

type Link struct {
	VersionedModel
	Title       string    `json:"title" sql:"not null" valid:"length(3|1000)"`
	Description string    `json:"desc" valid:"length(3|4000)"`
	Url         string    `json:"url" valid:"length(3|4000)"`
	ClaimID     uuid.UUID `json:"claimId" sql:"type:uuid;not null"`
	Claim       *Claim    `json:"claim,omitempty"`
}

func (l Link) ValidateForCreate() Error {
	return ValidateStruct(l)
}

func (l Link) ValidateForUpdate(updates Updates) Error {
	if err := SetJsonValuesOnStruct(&l, updates, false); err != nil {
		return err
	}
	return l.ValidateForCreate()
}

func (l Link) ValidateForDelete() Error {
	return nil
}

func (l Link) ValidateField(f string) Error {
	return ValidateStructField(l, f)
}
