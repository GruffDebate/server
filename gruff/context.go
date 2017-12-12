package gruff

type Context struct {
	Model
	ParentID         *uint64   `json:"parentId"`
	Parent           *Context  `json:"parent,omitempty"`
	Title            string    `json:"title" sql:"not null" valid:"length(3|1000)"`
	Description      string    `json:"desc" valid:"length(3|4000)"`
	URL              string    `json:"url" valid:"url,required"`
	MID              string    `json:"mid"` // Google KG ID
	QID              string    `json:"qid"` // Wikidata ID
	MetaDataURL      *MetaData `json:"meta_url"`
	MetaDataGoogle   *MetaData `json:"meta_google"`
	MetaDataWikidata *MetaData `json:"meta_wikidata"`
}

type MetaData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
}

func (c Context) ValidateForCreate() GruffError {
	return ValidateStruct(c)
}

func (c Context) ValidateForUpdate() GruffError {
	return c.ValidateForCreate()
}

func (c Context) ValidateField(f string) GruffError {
	return ValidateStructField(c, f)
}
