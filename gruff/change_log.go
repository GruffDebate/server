package gruff

const CHANGE_TYPE_CREATED_CLAIM int = 1
const CHANGE_TYPE_CREATED_ARGUMENT int = 2
const CHANGE_TYPE_CREATED_CLAIM_AND_ARGUMENT int = 3
const CHANGE_TYPE_MOVED_ARGUMENT int = 11
const CHANGE_TYPE_CLONE_CLAIM int = 21
const CHANGE_TYPE_MERGE_CLAIMS int = 31
const CHANGE_TYPE_MERGE_ARGUMENTS int = 32

/*
Types of Changes, and fields used:
- Created Claim: ClaimID
- Created Argument: ArgumentID, NewArgType, NewClaimID or NewArgID (parent)
- Created Claim and Argument: ArgumentID, NewArgType, NewClaimID or NewArgID (parent)
- Moved Argument: ArgumentID, OldClaimID or OldArgID, NewClaimID or NewArgID (parent), OldArgType, NewArgType
- Clone Claim:
  - One claim stays
  - New claim created, with same values, context, title and description (must be changed before saving)
  - Arguments stay with main claim
  --> Need Change Type add/remove values and contexts
  --> ClaimID, NewClaimID
  --> What about opinions?? I guess it would make a copy
  --> Would there be arguments between old and new claim(s)?
  ------- E.g. Fidel Castro is nice, and ended Apartheid --> Fidel Castro is nice, Fidel Castro ended Apartheid
- Merge Claims:
  - One claim becomes defunct
  - Must have "compatible" values/context (TBD)
  - All arguments attach to "winning" claim
  - Title, description stick with "winning" claim
  - All arguments with "losing" claim as base reattach to "winning" claim
  - Do we need a change log for each of these? Probably...
  - What about opinions? Should also merge...

*/
type ChangeLog struct {
	Model
	UserID     uint64        `json:"userId" sql:"not null"`
	User       *User         `json:"user,omitempty"`
	Type       int           `json:"type" sql:"not null"`
	ArgumentID *NullableUUID `json:"argumentId,omitempty" sql:"type:uuid"`
	Argument   *Argument     `json:"argument,omitempty"`
	ClaimID    *NullableUUID `json:"claimId,omitempty" sql:"type:uuid"`
	Claim      Claim         `json:"claim"`
	OldClaimID *NullableUUID `json:"oldClaimId,omitempty" sql:"type:uuid"`
	OldClaim   *Claim        `json:"oldClaim,omitempty"`
	OldArgID   *NullableUUID `json:"oldArgId,omitempty" sql:"type:uuid"`
	OldArg     *Argument     `json:"oldArg,omitempty"`
	NewClaimID *NullableUUID `json:"newClaimId,omitempty" sql:"type:uuid"`
	NewClaim   *Claim        `json:"newClaim,omitempty"`
	NewArgID   *NullableUUID `json:"newArgId,omitempty" sql:"type:uuid"`
	NewArg     *Argument     `json:"newArg,omitempty"`
	OldArgType *int          `json:"oldArgType,omitempty"`
	NewArgType *int          `json:"newArgType,omitempty"`
}

func (cl ChangeLog) ValidateForCreate() GruffError {
	return ValidateStruct(cl)
}

func (cl ChangeLog) ValidateForUpdate() GruffError {
	return cl.ValidateForCreate()
}

func (cl ChangeLog) ValidateField(f string) GruffError {
	return ValidateStructField(cl, f)
}
