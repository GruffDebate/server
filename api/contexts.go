package api

import (
	"net/http"

	"github.com/GruffDebate/server/gruff"
	"github.com/badoux/goscraper"
	"github.com/labstack/echo"
)

func ListContext(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	contexts := []gruff.Context{}

	db = BasicJoins(ctx, c, db)
	db = BasicPaging(ctx, c, db)

	err := db.Find(&contexts).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	for i := range contexts {
		s, err := goscraper.Scrape(contexts[i].URL, 1)
		if err != nil {
			return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
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
