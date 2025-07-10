package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/wahyujatirestu/sahabat-kurban/config"
	"github.com/wahyujatirestu/sahabat-kurban/controller"
	"github.com/wahyujatirestu/sahabat-kurban/middleware"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
	"github.com/wahyujatirestu/sahabat-kurban/routes"
	"github.com/wahyujatirestu/sahabat-kurban/service"
	utilsrepo "github.com/wahyujatirestu/sahabat-kurban/utils/repository"
	utilsservice "github.com/wahyujatirestu/sahabat-kurban/utils/service"
)

type Server struct {
	userRepo 				repository.UserRepository
	pekurbanRepo 			repository.PekurbanRepository
	hewanKurbanRepo			repository.HewanKurbanRepository
	pekurbanHewanRepo		repository.PekurbanHewanRepository
	userService 			service.UserService
	authService 			service.AuthService
	jwtService				utilsservice.JWTService
	pekurbanService 		service.PekurbanService
	hewanKurbanService 		service.HewanKurbanService
	pekurbanHewanService 	service.PekurbanHewanService
	rtRepo 					utilsrepo.RefreshTokenRepository
	db 						*sql.DB
	engine 					*gin.Engine
	host					string
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
	pekurbanRepo := repository.NewPekurbanRepository(db)
	hewanKurbanRepo := repository.NewHewanKurbanRepository(db)
	pekurbanHewanRepo := repository.NewPekurbanHewanRepository(db)
	jwtService := utilsservice.NewJWTServie(cfg, rtRepo)
	authService := service.NewAuthService(cfg, userRepo, rtRepo, jwtService)
	userService := service.NewUserService(userRepo)
	pekurbanService := service.NewPekurbanService(pekurbanRepo, userRepo)
	hewanKurbanService := service.NewHewanKurbanService(hewanKurbanRepo)
	pekurbanHewanService := service.NewPekurbanHewanService(pekurbanHewanRepo)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	seedInitialAdmin(userRepo)

	return &Server{
		userRepo: userRepo,
		rtRepo: rtRepo,
		pekurbanRepo: pekurbanRepo,
		hewanKurbanRepo: hewanKurbanRepo,
		pekurbanHewanRepo: pekurbanHewanRepo,
		db: db,
		authService: authService,
		userService: userService,
		jwtService: jwtService,
		pekurbanService: pekurbanService,
		hewanKurbanService: hewanKurbanService,
		pekurbanHewanService: pekurbanHewanService,
		engine: engine,
		host: host,
	}
}

func (s *Server) SetupRoutes() {
	apiV1 := s.engine.Group("/api/v1")
	authMw := middleware.NewAuthMiddleware(s.jwtService)

	authController := controller.NewAuthController(s.authService)
	userController := controller.NewUserController(s.userService)
	pekurbanController := controller.NewPekurbanController(s.pekurbanService)
	hewanKurbanController := controller.NewHewanKurbanController(s.hewanKurbanService)
	pekurbanHewanController := controller.NewPekurbanHewanController(s.pekurbanHewanService, s.pekurbanService)

	routes.AuthRoute(apiV1, authController)
	routes.UserRoute(apiV1, userController, authMw)
	routes.PekurbanRoute(apiV1, pekurbanController, authMw)
	routes.HewanKurbanRoute(apiV1, hewanKurbanController, authMw)
	routes.PekurbanHewanRoute(apiV1, pekurbanHewanController, authMw)
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
