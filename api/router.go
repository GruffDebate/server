package api

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var cors []string

type MiddlewareConfigurer interface {
	ConfigureDefaultApiMiddleware(*echo.Echo) *echo.Echo
	ConfigurePublicApiMiddleware(*echo.Echo) *echo.Group
	ConfigurePrivateApiMiddleware(*echo.Echo) *echo.Group
}

func SetUpRouter(mc MiddlewareConfigurer) *echo.Echo {
	root := echo.New()

	root.GET("/", Home)

	if os.Getenv("GRUFF_ENV") == "production" {
		cors = []string{"*"}
	} else {
		cors = []string{"*"}
	}

	//
	// PUBLIC ENDPOINTS
	//
	public := mc.ConfigurePublicApiMiddleware(root)

	public.POST("/auth", SignIn)
	public.POST("/users", SignUp)

	//
	// PRIVATE ENDPOINTS
	//
	private := mc.ConfigurePrivateApiMiddleware(root)

	private.GET("/users", List)
	private.GET("/users/:id", Get)
	private.GET("/users/me", GetMe)
	private.PUT("/users/me", UpdateMe)
	private.PUT("/users/:id", Update)
	private.PUT("/users/password", ChangePassword)
	public.PUT("/users/changePassword", ChangePassword)
	private.DELETE("/users/:id", Delete)

	private.GET("/users/claims", List)

	private.POST("/claims/:id/score", SetScore)
	private.PUT("/claims/:id/score", SetScore)
	private.POST("/arguments/:id/score", SetScore)
	private.PUT("/arguments/:id/score", SetScore)

	public.GET("/arguments/:id", Get)
	private.POST("/arguments", Create)
	private.PUT("/arguments/:id", Update)
	private.DELETE("/arguments/:id", Delete)
	private.PUT("/arguments/:id/move/:type/:targetId", MoveArgument)

	// TODO: Test all these
	public.GET("/contexts", List)
	public.GET("/contexts/search", SearchContext)
	// TODO: Should contexts be gotten using the short name instead of the id?
	public.GET("/contexts/:id", Get)
	private.POST("/contexts", Create)
	private.PUT("/contexts/:id", Update)
	private.DELETE("/contexts/:id", Delete)

	// TODO: Should contexts be added using the short name instead of the id?
	private.POST("/claims/:parentId/contexts/:id", AddContext)
	private.DELETE("/claims/:parentId/contexts/:id", RemoveContext)

	//private.POST("/claims/:parentId/tags/:id", AddAssociation)
	//private.DELETE("/claims/:parentId/tags/:id", RemoveAssociation)

	//private.PUT("/claims/:parentId/contexts", ReplaceAssociation)
	//private.PUT("/claims/:parentId/tags", ReplaceAssociation)

	public.GET("/claims", ListClaims("new"))
	public.GET("/claims/top", ListClaims("top"))
	public.GET("/claims/:id", Get)
	public.GET("/claims/:id/parents", ListParentArguments)
	private.POST("/claims", Create)
	private.PUT("/claims/:id", Update)
	private.DELETE("/claims/:id", Delete)
	private.PUT("/claims/:id/convert", ConvertClaimToMultiPremise)
	private.POST("/claims/:parentId/premises/:id", AddPremise)
	private.DELETE("/claims/:parentId/premises/:id", RemovePremise)
	//private.POST("/claims/:id/truth", SetScore)
	//private.PUT("/claims/:id/truth", SetScore)

	public.GET("/links", List)
	public.GET("/links/:id", Get)
	private.POST("/links", Create)
	private.PUT("/links/:id", Update)
	private.DELETE("/links/:id", Delete)

	//public.GET("/tags/:id/claims", ListClaimsByTag)

	//private.GET("/notifications", ListNotifications)
	//private.POST("/notifications/:id", MarkNotificationViewed)
	//private.PUT("/notifications/:id", MarkNotificationViewed)

	return root
}

type ProductionMiddlewareConfigurer struct{}

func (mc ProductionMiddlewareConfigurer) ConfigureDefaultApiMiddleware(root *echo.Echo) *echo.Echo {
	root.Use(middleware.Logger())
	root.Use(middleware.Recover())
	root.Use(middleware.CORS())
	root.Use(middleware.Secure())
	root.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		ContentSecurityPolicy: "default-src 'self'",
	}))
	root.Use(Secure(
		ReferrerPolicy("same-origin"),
	))

	root.Use(DBMiddleware(ARANGODB_POOL))
	root.Use(DetermineType)
	root.Use(InitializePayload)
	root.Use(SettingHeaders(false))

	return root
}

func (mc ProductionMiddlewareConfigurer) ConfigurePublicApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	public := api.Group("/api")
	public.Use(Session)

	return public
}

func (mc ProductionMiddlewareConfigurer) ConfigurePrivateApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	private := api.Group("/api")
	private.Use(Session)

	return private
}
