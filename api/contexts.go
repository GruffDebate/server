package api

import (
	"net/http"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

func SearchContexts(c echo.Context) error {
	ctx := ServerContext(c)

	query := c.QueryParam("query")

	result, err := gruff.SearchContexts(ctx, query)
	if err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, result)
}

/*
func ListContext(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	contexts := []gruff.Context{}

	db = DefaultJoins(ctx, c, db)
	db = DefaultPaging(ctx, c, db)

	if err := db.Find(&contexts).Error; err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	for i := range contexts {
		s, err := goscraper.Scrape(contexts[i].URL, 1)
		if err != nil {
			return AddError(ctx, c, gruff.NewServerError(err.Error()))
		}

		contexts[i].MetaDataURL = &gruff.MetaData{
			Title:       s.Preview.Title,
			Description: s.Preview.Description,
			Image:       s.Preview.Images[0],
			URL:         s.Preview.Link,
		}
	}

	if ctx.Payload["ct"] != nil {
		ctx.Payload["results"] = contexts
		return c.JSON(http.StatusOK, ctx.Payload)
	}

	return c.JSON(http.StatusOK, contexts)
}
*/
