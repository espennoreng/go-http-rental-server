package api

import (
	"log/slog"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/config"
	customMiddleware "github.com/espennoreng/go-http-rental-server/internal/middleware"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router chi.Router
}

func NewServer(
	cfg *config.AppConfig,
	verifier auth.TokenVerifier,
	log *slog.Logger,
	userService services.UserService,
	organizationService services.OrganizationService,
	organizationUserService services.OrganizationUserService,
	accessService services.AccessService,
) *Server {
	userHandler := NewUserHandler(userService, log)
	organizationHandler := NewOrganizationHandler(organizationService, log)
	organizationUserHandler := NewOrganizationUserHandler(organizationUserService, log)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.NewSlogMiddleware(log))

	setupRoutes(r, cfg, verifier, userService, userHandler, organizationHandler, organizationUserHandler, accessService)

	return &Server{
		router: r,
	}
}

func setupRoutes(
	r chi.Router,
	cfg *config.AppConfig,
	verifier auth.TokenVerifier,
	userService services.UserService,
	userHandler *userHandler,
	organizationHandler *organizationHandler,
	organizationUserHandler *organizationUserHandler,
	accessService services.AccessService,
) {

	authMiddleware := customMiddleware.NewAuthMiddleware(verifier, userService, cfg.GoogleOAuthClientID)
	accessMiddleware := customMiddleware.NewAccessMiddleware(accessService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Rental Server API"))
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(authMiddleware)
		
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			userHandler.CreateUser(w, r)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			userHandler.GetUserByID(w, r)
		})
	})

	r.Route("/organizations", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			organizationHandler.CreateOrganization(w, r)
		})

		r.With(accessMiddleware.RequireMember).Route("/{id}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				organizationHandler.GetOrganizationByID(w, r)
			})
		})

		r.With(accessMiddleware.RequireAdmin).Route("/{orgID}/users", func(r chi.Router) {

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				organizationUserHandler.AddUserToOrganization(w, r)
			})

			r.Put("/{userID}/role", func(w http.ResponseWriter, r *http.Request) {
				organizationUserHandler.UpdateUserRole(w, r)
			})

			r.Delete("/{userID}", func(w http.ResponseWriter, r *http.Request) {
				organizationUserHandler.DeleteUserFromOrganization(w, r)
			})
		})

		r.With(accessMiddleware.RequireMember).Route("/", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				organizationUserHandler.GetUsersByOrganizationID(w, r)
			})
		})
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
