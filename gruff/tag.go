package gruff

type Tag struct {
	Model
	Title  string  `json:"title" sql:"not null" valid:"length(3|50)"`
	Claims []Claim `json:"claims,omitempty"  gorm:"many2many:claim_tags;"`
}

func (t Tag) ValidateForCreate() GruffError {
	return ValidateStruct(t)
}

func (t Tag) ValidateForUpdate() GruffError {
	return t.ValidateForCreate()
}

func (t Tag) ValidateField(f string) GruffError {
	return ValidateStructField(t, f)
}
