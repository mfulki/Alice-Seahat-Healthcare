package server

import (
	"database/sql"
	"io"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/handler"
	"Alice-Seahat-Healthcare/seahat-be/libs/firebase"
	"Alice-Seahat-Healthcare/seahat-be/libs/mail"
	"Alice-Seahat-Healthcare/seahat-be/libs/rajaongkir"
	"Alice-Seahat-Healthcare/seahat-be/libs/validator"
	"Alice-Seahat-Healthcare/seahat-be/middleware"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Handlers struct {
	CustomHandler          *handler.CustomHandler
	Middleware             *middleware.Middleware
	DrugHandler            *handler.DrugHandler
	PharmacyDrugHandler    *handler.PharmacyDrugHandler
	UserHandler            *handler.UserHandler
	DoctorHandler          *handler.DoctorHandler
	PharmacyManagerHandler *handler.PharmacyManagerHandler
	AdminHandler           *handler.AdminHandler
	UploadHandler          *handler.UploadHandler
	TelemedicineHandler    *handler.TelemedicineHandler
	MessageBubbleHandler   *handler.MessageBubbleHandler
	AdminReportHandler     *handler.AdminReportHandler
	SpecializationHandler  *handler.SpecializationHandler
	PartnerHandler         *handler.PartnerHandler
	AddressHandler         *handler.AddressHandler
	CartItemHandler        *handler.CartItemHandler
	OrderHandler           *handler.OrderHandler
	PaymentHandler         *handler.PaymentHandler
	PharmacyHandler        *handler.PharmacyHandler
	StockRequestHandler    *handler.StockRequestHandler
	ShipmentMethodHandler  *handler.ShipmentMethodHandler
	CategoryHandler        *handler.CategoryHandler
	StockJournalHandler    *handler.StockJournalHandler
	ManufacturerHandler    *handler.ManufacturerHandler
}

type Server struct {
	transactor transaction.Transactor
	db         transaction.DBTransaction
	mailDialer mail.MailDialer
	appLog     io.Writer
	rajaOngkir rajaongkir.RajaOngkir
	firebase   firebase.Firebase
}

func NewServer(
	db *sql.DB,
	dialer mail.MailDialer,
	appLog io.Writer,
	rajaOngkir rajaongkir.RajaOngkir,
	firebase firebase.Firebase,
) *Server {
	return &Server{
		transactor: transaction.NewTransactor(db),
		db:         transaction.NewDBTransaction(db),
		mailDialer: dialer,
		appLog:     appLog,
		rajaOngkir: rajaOngkir,
		firebase:   firebase,
	}
}

func (s Server) SetupServer() *gin.Engine {
	validator.SetCustom(binding.Validator.Engine())

	customHandler := handler.NewCustomHandler()
	middleware := middleware.NewMiddleware()

	drugRepository := repository.NewDrugRepository(s.db)
	pharmacyDrugRepository := repository.NewPharmacyDrugRepository(s.db)
	userRepository := repository.NewUserRepository(s.db)
	tokenRepository := repository.NewTokenRepository(s.db)
	doctorRepository := repository.NewDoctorRepository(s.db)
	pharmacyManagerRepository := repository.NewPharmacyManagerRepository(s.db)
	adminRepository := repository.NewAdminRepository(s.db)
	partnerRepository := repository.NewPartnerRepository(s.db)
	telemedicineRepository := repository.NewTelemedicineRepository(s.db)
	messageBubbleRepository := repository.NewMessageBubblesRepository(s.db)
	adminReportRepository := repository.NewAdminReportRepository(s.db)
	specializationRepository := repository.NewSpecializationRepository(s.db)
	addressRepository := repository.NewAddressRepository(s.db)
	cartItemRepository := repository.NewCartItemRepository(s.db)
	orderRepository := repository.NewOrderRepository(s.db)
	paymentRepository := repository.NewPaymentRepository(s.db)
	orderDetailRepository := repository.NewOrderDetailRepository(s.db)
	pharmacyRepository := repository.NewPharmacyRepository(s.db)
	stockJournalRepository := repository.NewStockJournalRepository(s.db)
	stockRequestRepository := repository.NewStockRequestRepository(s.db)
	stockRequestDrugRepository := repository.NewStockRequestDrugRepository(s.db)
	shipmentMethodRepository := repository.NewShipmentMethodRepository(s.db, s.rajaOngkir)
	categoryRepository := repository.NewCategoryRepository(s.db)
	prescriptionRepository := repository.NewPrescriptionRepository(s.db)
	manufacturerRepository := repository.NewManufacturerRepository(s.db)

	drugUsecase := usecase.NewDrugUsecase(drugRepository, s.transactor)
	userUsecase := usecase.NewUserUsecase(userRepository, doctorRepository, tokenRepository, s.transactor, s.mailDialer, s.firebase)
	pharmacyDrugUsecase := usecase.NewPharmacyDrugUsecase(pharmacyDrugRepository, addressRepository, drugRepository, pharmacyRepository, categoryRepository, s.transactor, stockJournalRepository)
	doctorUsecase := usecase.NewDoctorUsecase(doctorRepository, tokenRepository, s.transactor, s.mailDialer, s.firebase)
	pharmacyManagerUsecase := usecase.NewPharmacyManagerUsecase(pharmacyManagerRepository, partnerRepository, s.transactor)
	adminUsecase := usecase.NewAdminUsecase(userRepository, doctorRepository, pharmacyManagerRepository, adminRepository, s.transactor)
	uploadUsecase := usecase.NewUploadUsecase()
	telemedicineUsecase := usecase.NewTelemedicineUsecase(telemedicineRepository, userRepository, doctorRepository, prescriptionRepository)
	messageBubbleUsecase := usecase.NewMessageBubbleUsecase(messageBubbleRepository)
	adminReportUsecase := usecase.NewAdminReportUsecase(adminReportRepository)
	specializationUsecase := usecase.NewSpecializationUsecase(specializationRepository, s.transactor)
	partnerUsecase := usecase.NewPartnerUsecase(pharmacyManagerRepository, partnerRepository, s.transactor, s.mailDialer)
	addressUsecase := usecase.NewAddressUsecase(addressRepository, shipmentMethodRepository, s.transactor)
	cartItemUsecase := usecase.NewCartItemUsecase(cartItemRepository, pharmacyDrugRepository, s.transactor)
	orderUsecase := usecase.NewOrderUsecase(orderRepository, orderDetailRepository, s.transactor, paymentRepository, pharmacyDrugRepository, cartItemRepository, stockJournalRepository, stockRequestRepository, stockRequestDrugRepository, shipmentMethodRepository, addressRepository)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepository, orderRepository, s.transactor)
	pharmacyUsecase := usecase.NewPharmacyUsecase(pharmacyRepository, shipmentMethodRepository, s.transactor)
	stockRequestUsecase := usecase.NewStockRequestUsecase(pharmacyRepository, stockRequestRepository, stockRequestDrugRepository, s.transactor, pharmacyDrugRepository, stockJournalRepository)
	shipmentMethodUsecase := usecase.NewShipmentMethodUsecase(shipmentMethodRepository, s.transactor)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository, s.transactor, pharmacyDrugRepository)
	stockJournalUsecase := usecase.NewStockJournalUsecase(orderRepository, orderDetailRepository, s.transactor, paymentRepository, pharmacyDrugRepository, cartItemRepository, stockJournalRepository, stockRequestRepository, stockRequestDrugRepository)
	manufacturerUsecase := usecase.NewManufacturerUsecase(manufacturerRepository)

	drugHandler := handler.NewDrugHandler(drugUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	doctorHandler := handler.NewDoctorHandler(doctorUsecase)
	pharmacyDrugHandler := handler.NewPharmacyDrugHandler(pharmacyDrugUsecase)
	pharmacyManagerHandler := handler.NewPharmacyManagerHandler(pharmacyManagerUsecase)
	adminHandler := handler.NewAdminHandler(adminUsecase)
	uploadHandler := handler.NewUploadHandler(uploadUsecase)
	telemedicineHandler := handler.NewTelemedicineHandler(telemedicineUsecase)
	messageBubbleHandler := handler.NewMessageBubbleHandler(messageBubbleUsecase)
	adminReportHandler := handler.NewAdminReportHandler(adminReportUsecase)
	specializationHandler := handler.NewSpecializationHandler(specializationUsecase)
	partnerHandler := handler.NewPartnerHandler(partnerUsecase)
	addressHandler := handler.NewAddressHandler(addressUsecase)
	cartItemHandler := handler.NewCartItemHandler(cartItemUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)
	pharmacyHandler := handler.NewPharmacyHandler(pharmacyUsecase)
	stockRequestHandler := handler.NewStockRequestHandler(stockRequestUsecase)
	shipmentMethodHandler := handler.NewShipmentMethodHandler(shipmentMethodUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	stockJournalHandler := handler.NewStockJournalHandler(stockJournalUsecase)
	manufacturerHandler := handler.NewManufacturerHandler(manufacturerUsecase)

	return SetupRouter(&Handlers{
		CustomHandler:          customHandler,
		UserHandler:            userHandler,
		Middleware:             middleware,
		DrugHandler:            drugHandler,
		DoctorHandler:          doctorHandler,
		PharmacyDrugHandler:    pharmacyDrugHandler,
		PharmacyManagerHandler: pharmacyManagerHandler,
		AdminHandler:           adminHandler,
		UploadHandler:          uploadHandler,
		TelemedicineHandler:    telemedicineHandler,
		MessageBubbleHandler:   messageBubbleHandler,
		AdminReportHandler:     adminReportHandler,
		SpecializationHandler:  specializationHandler,
		PartnerHandler:         partnerHandler,
		AddressHandler:         addressHandler,
		CartItemHandler:        cartItemHandler,
		OrderHandler:           orderHandler,
		PaymentHandler:         paymentHandler,
		PharmacyHandler:        pharmacyHandler,
		StockRequestHandler:    stockRequestHandler,
		ShipmentMethodHandler:  shipmentMethodHandler,
		CategoryHandler:        categoryHandler,
		StockJournalHandler:    stockJournalHandler,
		ManufacturerHandler:    manufacturerHandler,
	}, s.appLog)
}
