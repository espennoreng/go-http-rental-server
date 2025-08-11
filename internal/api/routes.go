package api

import (
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


type Server struct {
	router chi.Router
}

func NewServer(
	userService services.UserService,
	organizationService services.OrganizationService,
) *Server {
	userHandler := NewUserHandler(userService)
	organizationHandler := NewOrganizationHandler(organizationService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	setupRoutes(r, userHandler, organizationHandler)

	return &Server{
		router: r,
	}
}

func setupRoutes(r chi.Router, userHandler *userHandler, organizationHandler *organizationHandler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Rental Server API"))
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			userHandler.CreateUser(w, r)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			userHandler.GetUserByID(w, r)
		})
	})

	r.Route("/organizations", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			organizationHandler.CreateOrganization(w, r)
		})

	
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
