package http

import (
	"net/http"
	"os"
	"time"

	"dekamond/internal/config"
	userdomain "dekamond/internal/domain/user"
	"dekamond/internal/http/handlers"
	"dekamond/internal/http/middleware"
	"dekamond/internal/infra/cache"
	postgresrepositories "dekamond/internal/infra/db/postgres/repositories"
	authusecase "dekamond/internal/usecase/auth"
	userusecase "dekamond/internal/usecase/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(conf config.Config, pg *pgxpool.Pool, redis *redis.Client) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.AllowAll().Handler)

	userRepo := postgresrepositories.NewPostgresUserRepository(pg)
	var _ userdomain.Repository = userRepo

	cacheStore := cache.NewRedisStore(redis)
	var _ userdomain.CacheStore = cacheStore

	authUsecase := authusecase.New(userRepo, cacheStore, conf)
	authHandler := handlers.NewAuthHandler(authUsecase)

	userUsecase := userusecase.New(userRepo)
	userHandler := handlers.NewUserHandler(userUsecase)

	r.Route("/api", func(api chi.Router) {
		api.Route("/auth", func(auth chi.Router) {
			auth.With(middleware.OTPRateLimit(cacheStore, 3, 10*time.Minute)).Post("/request-otp", authHandler.RequestOTP)
			auth.Post("/verify-otp", authHandler.VerifyOTP)
		})
		api.Route("/users", func(users chi.Router) {
			users.Use(middleware.JwtAuth(conf.JWTSecret))
			users.Get("/", userHandler.List)
			users.Get("/{id}", userHandler.GetByID)
		})
	})

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		openapiPath := "/openapi.yaml"
		if _, err := os.Stat(openapiPath); os.IsNotExist(err) {
			openapiPath = "openapi.yaml"
		}
		
		data, err := os.ReadFile(openapiPath)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/openapi.yaml"),
	))

	return r
}