package gruff

/*
 * A Claim is a proposed statement of fact
 *
 * According to David Zarefsky (https://www.thegreatcoursesplus.com/argumentation/argument-analysis-and-diagramming) there are 4 types:
 * - Fact: Al Gore received more popular votes than George Bush in the 2000 election
 * - Definition: Capital execution is murder
 * - Value: Environmental protection is more important than economic growth
 * - Policy: Congress should pass the president's budget
 *
 * Also according to the professor, there are 4 parts to a claim/argument:
 * - Claim
 * - Evidence
 * - Inference
 * - Warrant
 *
 * In loose terms, a Claim here represents his Claim, and Evidence
 * An Argument of type 1 or 2 (truth) is an Inference
 * An Argument of type 3, 4, 5 or 6 is a Warrant
 *
 * Complex Claims:
 * - Series: Because of X, Y happened, which caused Z --> Not modeled in Gruff
 * - Convergent: Airline travel is becoming more unpleasant because of X, Y, Z, P, D, and Q --> Supported by standard Gruff structure
 * - Parallel: Same as convergent, except that any one argument is enough --> Supported by standard Gruff structure
 *
 * Topoi for Resolutions of Fact (for scoring Truth):
 * - What are the criteria (of truth)?
 * - Are the criteria satisfied?
 */
type Claim struct {
	Identifier
	Title       string     `json:"title" sql:"not null" valid:"length(3|1000)"`
	Description string     `json:"desc" valid:"length(3|4000)"`
	Image       string     `json:"img,omitempty"`
	Truth       float64    `json:"truth"`
	TruthRU     float64    `json:"truthRU"` // Average score rolled up from argument totals
	ProTruth    []Argument `json:"protruth,omitempty"`
	ConTruth    []Argument `json:"contruth,omitempty"`
	Links       []Link     `json:"links,omitempty"`
	Contexts    []Context  `json:"contexts,omitempty"  gorm:"many2many:claim_contexts;"`
	ContextIDs  []uint64   `json:"contextIds,omitempty" gorm:"-"`
	Values      []Value    `json:"values,omitempty"  gorm:"many2many:claim_values;"`
	Tags        []Tag      `json:"tags,omitempty"  gorm:"many2many:claim_tags;"`
}

func (c Claim) ValidateForCreate() GruffError {
	return ValidateStruct(c)
}

func (c Claim) ValidateForUpdate() GruffError {
	return c.ValidateForCreate()
}

func (c Claim) ValidateField(f string) GruffError {
	return ValidateStructField(c, f)
}

func (c Claim) UpdateTruth(ctx *ServerContext) {
	ctx.Database.Exec("UPDATE claims c SET truth = (SELECT AVG(truth) FROM claim_opinions WHERE claim_id = c.id) WHERE id = ?", c.ID)

	// TODO: test
	if c.TruthRU == 0.0 {
		// There's no roll up score yet, so the truth score itself is affecting related roll ups
		c.UpdateAncestorRUs(ctx)
	}
}

func (c *Claim) UpdateTruthRU(ctx *ServerContext) {
	// TODO: do it all in SQL?
	// TODO: should updates be recursive? (first, calculate sub-argument RUs)
	//       or, should it trigger an update of anyone that references it?
	proArgs, conArgs := c.Arguments(ctx)

	if len(proArgs) > 0 || len(conArgs) > 0 {
		proScore := 0.0
		for _, arg := range proArgs {
			remainder := 1.0 - proScore
			score := arg.ScoreRU(ctx)
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

	ctx.Database.Set("gorm:save_associations", false).Save(c)

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

func (c Claim) Arguments(ctx *ServerContext) (proArgs []Argument, conArgs []Argument) {
	proArgs = c.ProTruth
	conArgs = c.ConTruth

	if len(proArgs) == 0 {
		ctx.Database.
			Preload("Claim").
			Scopes(OrderByBestArgument).
			Where("type = ?", ARGUMENT_TYPE_PRO_TRUTH).
			Where("target_claim_id = ?", c.ID).
			Find(&proArgs)
	}

	if len(conArgs) == 0 {
		ctx.Database.
			Preload("Claim").
			Scopes(OrderByBestArgument).
			Where("type = ?", ARGUMENT_TYPE_CON_TRUTH).
			Where("target_claim_id = ?", c.ID).
			Find(&conArgs)
	}

	return
}
