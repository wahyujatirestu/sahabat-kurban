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
	payserv "github.com/wahyujatirestu/sahabat-kurban/payments/service"
	utilsrepo "github.com/wahyujatirestu/sahabat-kurban/utils/repository"
	utilsservice "github.com/wahyujatirestu/sahabat-kurban/utils/service"
)

type Server struct {
	userRepo 				repository.UserRepository
	pekurbanRepo 			repository.PekurbanRepository
	hewanKurbanRepo			repository.HewanKurbanRepository
	pekurbanHewanRepo		repository.PekurbanHewanRepository
	penyembelihanRepo		repository.PenyembelihanRepository
	penerimaRepo			repository.PenerimaDagingRepository
	distribusiRepo			repository.DistribusiDagingRepository
	pembayaranRepo			repository.PembayaranKurbanRepository
	emailRepo				utilsrepo.EmailVerificationRepository
	resetRepo 				utilsrepo.ResetPasswordRepository
	userService 			service.UserService
	authService 			service.AuthService
	emailService			utilsservice.EmailService
	jwtService				utilsservice.JWTService
	pekurbanService 		service.PekurbanService
	hewanKurbanService 		service.HewanKurbanService
	pekurbanHewanService 	service.PekurbanHewanService
	penyembelihanService 	service.PenyembelihanService
	penerimaService 		service.PenerimaDagingService
	distribusiService 		service.DistribusiDagingService
	midtransService			payserv.MidtransService
	pembayaranService		service.PembayaranKurbanService
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
	emailRepo := utilsrepo.NewEmailVerificationRepository(db)
	resetRepo := utilsrepo.NewResetPasswordRepository(db)
	pekurbanRepo := repository.NewPekurbanRepository(db)
	hewanKurbanRepo := repository.NewHewanKurbanRepository(db)
	pekurbanHewanRepo := repository.NewPekurbanHewanRepository(db)
	penyembelihanRepo := repository.NewPenyembelihanRepository(db)
	penerimaRepo := repository.NewPenerimaDagingRepository(db)
	distribusiRepo := repository.NewDistribusiDagingRepository(db)
	pembayaranRepo := repository.NewPembayaranKurbanRepository(db)

	emailService := utilsservice.NewEmailService(
		cfg.SendgridAPIKey,
		cfg.EmailSender,
		cfg.EmailSenderName,
		cfg.AppBaseURL,
	)

	jwtService := utilsservice.NewJWTServie(cfg, rtRepo, userRepo)
	authService := service.NewAuthService(cfg, userRepo, rtRepo, emailRepo, resetRepo, jwtService, emailService)
	userService := service.NewUserService(userRepo)
	pekurbanService := service.NewPekurbanService(pekurbanRepo, userRepo)
	hewanKurbanService := service.NewHewanKurbanService(hewanKurbanRepo, penyembelihanRepo)
	pekurbanHewanService := service.NewPekurbanHewanService(pekurbanHewanRepo, hewanKurbanRepo)
	penyembelihanService := service.NewPenyembelihanService(penyembelihanRepo)
	penerimaService := service.NewPenerimaDagingService(penerimaRepo, pekurbanRepo)
	distribusiService := service.NewDistribusiDagingService(distribusiRepo, penerimaRepo, hewanKurbanRepo)
	midtransService := payserv.NewMidtransService()
	pembayaranService := service.NewPembayaranKurbanService(pembayaranRepo, midtransService, pekurbanHewanRepo, hewanKurbanRepo, pekurbanRepo)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	seedInitialAdmin(userRepo)

	return &Server{
		userRepo: userRepo,
		rtRepo: rtRepo,
		resetRepo: resetRepo,
		pekurbanRepo: pekurbanRepo,
		hewanKurbanRepo: hewanKurbanRepo,
		pekurbanHewanRepo: pekurbanHewanRepo,
		penyembelihanRepo: penyembelihanRepo,
		penerimaRepo: penerimaRepo,
		distribusiRepo: distribusiRepo,
		pembayaranRepo: pembayaranRepo,
		emailRepo: emailRepo,
		db: db,
		authService: authService,
		userService: userService,
		emailService: emailService,
		jwtService: jwtService,
		pekurbanService: pekurbanService,
		hewanKurbanService: hewanKurbanService,
		pekurbanHewanService: pekurbanHewanService,
		penyembelihanService: penyembelihanService,
		penerimaService: penerimaService,
		distribusiService: distribusiService,
		midtransService: midtransService,
		pembayaranService: pembayaranService,
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
	penyembelihanController := controller.NewPenyembelihanController(s.penyembelihanService)
	penerimaController := controller.NewPenerimaDagingController(s.penerimaService)
	distribusiController := controller.NewDistribusiDagingController(s.distribusiService)
	pembayaranController := controller.NewPembayaranController(s.pembayaranService, s.pekurbanService)

	routes.AuthRoute(apiV1, authController)
	routes.UserRoute(apiV1, userController, authMw)
	routes.PekurbanRoute(apiV1, pekurbanController, authMw)
	routes.HewanKurbanRoute(apiV1, hewanKurbanController, authMw)
	routes.PekurbanHewanRoute(apiV1, pekurbanHewanController, authMw)
	routes.PenyembelihanRoute(apiV1, penyembelihanController, authMw)
	routes.PenerimaDagingRoute(apiV1, penerimaController, authMw)
	routes.DistribusiDagingRoute(apiV1, distribusiController, authMw)
	routes.PembayaranRoute(apiV1, pembayaranController, authMw)
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
