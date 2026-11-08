package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"gogogo/internal/handler"
	"gogogo/internal/middleware"
	"gogogo/internal/model"
	"gogogo/internal/pkg/response"
)

type Handlers struct {
	Auth      *handler.AuthHandler
	User      *handler.UserHandler
	Plant     *handler.PlantHandler
	Area      *handler.AreaHandler
	Equipment *handler.EquipmentHandler
	Point     *handler.PointHandler
	Source    *handler.SourceHandler
	Binding   *handler.BindingHandler
	Sample    *handler.SampleHandler
	Health    *handler.HealthHandler
}

func New(handlers Handlers, authMiddleware gin.HandlerFunc, corsOrigins []string, logger *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Recovery(logger),
		middleware.Logger(logger),
		middleware.CORS(corsOrigins),
	)
	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "route not found")
	})

	r.GET("/healthz", handlers.Health.Healthz)
	r.GET("/readyz", handlers.Health.Readyz)

	admin := middleware.RequireRole(model.RoleAdmin)
	operator := middleware.RequireRole(model.RoleAdmin, model.RoleOperator)

	api := r.Group("/api/v1")
	api.GET("/healthz", handlers.Health.Healthz)
	api.GET("/readyz", handlers.Health.Readyz)

	auth := api.Group("/auth")
	auth.POST("/login", handlers.Auth.Login)
	auth.POST("/refresh", handlers.Auth.Refresh)
	auth.POST("/register", authMiddleware, admin, handlers.Auth.Register)
	auth.POST("/logout", authMiddleware, handlers.Auth.Logout)

	users := api.Group("/users", authMiddleware)
	users.GET("/me", handlers.User.Me)
	users.PUT("/me", handlers.User.UpdateProfile)
	users.PUT("/me/password", handlers.User.ChangePassword)

	plants := api.Group("/plants", authMiddleware)
	plants.GET("", handlers.Plant.List)
	plants.POST("", admin, handlers.Plant.Create)
	plants.GET("/:id", handlers.Plant.Get)
	plants.PUT("/:id", admin, handlers.Plant.Update)
	plants.DELETE("/:id", admin, handlers.Plant.Delete)
	plants.GET("/:id/areas", handlers.Area.ListByPlant)
	plants.GET("/:id/latest", handlers.Sample.LatestByPlant)

	areas := api.Group("/areas", authMiddleware)
	areas.POST("", admin, handlers.Area.Create)
	areas.GET("/:id", handlers.Area.Get)
	areas.PUT("/:id", admin, handlers.Area.Update)
	areas.DELETE("/:id", admin, handlers.Area.Delete)

	equipments := api.Group("/equipments", authMiddleware)
	equipments.GET("", handlers.Equipment.List)
	equipments.POST("", operator, handlers.Equipment.Create)
	equipments.GET("/:id", handlers.Equipment.Get)
	equipments.PUT("/:id", operator, handlers.Equipment.Update)
	equipments.DELETE("/:id", operator, handlers.Equipment.Delete)

	points := api.Group("/measurement-points", authMiddleware)
	points.GET("", handlers.Point.List)
	points.POST("", operator, handlers.Point.Create)
	points.GET("/:id", handlers.Point.Get)
	points.PUT("/:id", operator, handlers.Point.Update)
	points.DELETE("/:id", operator, handlers.Point.Delete)
	points.GET("/:id/samples", handlers.Sample.Query)
	points.GET("/:id/latest", handlers.Sample.Latest)

	sources := api.Group("/data-sources", authMiddleware)
	sources.GET("", handlers.Source.List)
	sources.POST("", admin, handlers.Source.Create)
	sources.GET("/:id", handlers.Source.Get)
	sources.PUT("/:id", admin, handlers.Source.Update)
	sources.DELETE("/:id", admin, handlers.Source.Delete)
	sources.POST("/:id/test", admin, handlers.Source.Test)

	bindings := api.Group("/point-bindings", authMiddleware)
	bindings.GET("", handlers.Binding.List)
	bindings.POST("", operator, handlers.Binding.Create)
	bindings.GET("/:id", handlers.Binding.Get)
	bindings.PUT("/:id", operator, handlers.Binding.Update)
	bindings.DELETE("/:id", operator, handlers.Binding.Delete)

	return r
}
