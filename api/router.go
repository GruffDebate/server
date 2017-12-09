package api

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetUpRouter(test bool, db *gorm.DB) *echo.Echo {
	root := echo.New()

	root.Use(middleware.Logger())
	root.Use(middleware.Recover())
	root.Use(middleware.CORS())
	root.Use(middleware.Gzip())
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

	root.Use(DBMiddleware(db))
	root.Use(DetermineType)
	root.Use(InitializePayload)
	root.Use(SettingHeaders(test))

	root.GET("/", Home)

	// Public Api
	public := root.Group("/api")

	public.POST("/auth", SignIn)
	public.POST("/users", SignUp)

	// Private Api
	private := root.Group("/api")
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	private.Use(middleware.JWTWithConfig(config))
	private.Use(SessionUser)

	private.GET("/users", List)
	private.GET("/users/:id", Get)
	private.GET("/users/me", GetMe)
	private.PUT("/users/me", UpdateMe)
	private.PUT("/users/:id", Update)
	private.PUT("/users/password", ChangePassword)
	public.PUT("/users/changePassword", ChangePassword)
	private.DELETE("/users/:id", Delete)

	private.GET("/users/claims", ListClaimsUser)

	public.GET("/arguments", List)
	public.GET("/arguments/:id", GetArgument)
	private.POST("/arguments", CreateArgument)
	private.PUT("/arguments/:id", Update)
	private.DELETE("/arguments/:id", Delete)
	private.PUT("/arguments/:id/move/:newId/type/:type", MoveArgument)

	private.POST("/arguments/:id/strength", SetScore)
	private.PUT("/arguments/:id/strength", SetScore)

	public.GET("/contexts", List)
	public.GET("/contexts/:id", Get)
	private.POST("/contexts", Create)
	private.PUT("/contexts/:id", Update)
	private.DELETE("/contexts/:id", Delete)

	private.POST("/claims/:parentId/contexts/:id", AddAssociation)
	private.DELETE("/claims/:parentId/contexts/:id", RemoveAssociation)

	private.POST("/claims/:parentId/tags/:id", AddAssociation)
	private.DELETE("/claims/:parentId/tags/:id", RemoveAssociation)

	private.PUT("/claims/:parentId/contexts", ReplaceAssociation)
	private.PUT("/claims/:parentId/tags", ReplaceAssociation)

	public.GET("/claims", List)
	public.GET("/claims/top", ListTopClaims)
	public.GET("/claims/:id", GetClaim)
	private.POST("/claims", Create)
	private.PUT("/claims/:id", Update)
	private.DELETE("/claims/:id", Delete)
	private.POST("/claims/:id/truth", SetScore)
	private.PUT("/claims/:id/truth", SetScore)

	public.GET("/links", List)
	public.GET("/links/:id", Get)
	private.POST("/links", Create)
	private.PUT("/links/:id", Update)
	private.DELETE("/links/:id", Delete)

	public.GET("/tags", List)
	public.GET("/tags/:id", Get)
	private.POST("/tags", Create)
	private.PUT("/tags/:id", Update)
	private.DELETE("/tags/:id", Delete)

	public.GET("/tags/:id/claims", ListClaimsByTag)

	public.GET("/values", List)
	public.GET("/values/:id", Get)
	private.POST("/values", Create)
	private.PUT("/values/:id", Update)
	private.DELETE("/values/:id", Delete)

	private.GET("/notifications", ListNotifications)
	private.POST("/notifications/:id", MarkNotificationViewed)
	private.PUT("/notifications/:id", MarkNotificationViewed)

	return root
}
