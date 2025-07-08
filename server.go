package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/wahyujatirestu/sahabat-kurban/config"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
	"github.com/wahyujatirestu/sahabat-kurban/routes"
	"github.com/wahyujatirestu/sahabat-kurban/service"
	utilsrepo "github.com/wahyujatirestu/sahabat-kurban/utils/repository"
	utilsservice "github.com/wahyujatirestu/sahabat-kurban/utils/service"
)

type Server struct {
	userRepo 	repository.UserRepository
	authService service.AuthService
	jwtService	utilsservice.JWTService
	rtRepo 		utilsrepo.RefreshTokenRepository
	db 			*sql.DB
	engine 		*gin.Engine
	host		string
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	rtRepo := utilsrepo.NewRefreshTokenRepository(db)
	jwtService := utilsservice.NewJWTServie(cfg, rtRepo)
	authService := service.NewAuthService(cfg, userRepo, rtRepo, jwtService)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		userRepo: userRepo,
		rtRepo: rtRepo,
		db: db,
		authService: authService,
		jwtService: jwtService,
		engine: engine,
		host: host,
	}
}

func (s *Server) SetupRoutes() {
	apiV1 := s.engine.Group("/api/v1")

	authController := controller.NewAuthController(s.authService)

	routes.AuthRoute(apiV1, authController)
}

func (s *Server) Run() {
	s.SetupRoutes()
	if err := s.engine.Run(s.host); err != nil {
		log.Fatalf("failed to run server on %s: %v", s.host, err)
	}
}

func (s *Server) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
