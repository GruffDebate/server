package gruff

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

func (c Context) ValidateForUpdate(updates map[string]interface{}) GruffError {
	if err := SetJsonValuesOnStruct(&c, updates); err != nil {
		return err
	}
	return c.ValidateForCreate()
}

func (c Context) ValidateField(f string) GruffError {
	return ValidateStructField(c, f)
}
