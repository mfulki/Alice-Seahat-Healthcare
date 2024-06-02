package server

import (
	"io"
	"time"

	"Alice-Seahat-Healthcare/seahat-be/config"
	"Alice-Seahat-Healthcare/seahat-be/constant"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRouter(h *Handlers, fileLog io.Writer) *gin.Engine {
	router := gin.New()
	router.ContextWithFallback = true

	if config.App.Env == constant.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	log := logrus.New()
	log.SetOutput(fileLog)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	router.Use(h.Middleware.Cors())
	router.Use(h.Middleware.RequestID)
	router.Use(h.Middleware.Logger(log))
	router.Use(h.Middleware.ErrorHandling(logrus.New()))
	router.Use(gin.Recovery())

	mwUserDoctor := h.Middleware.AuthMulti([]string{constant.User, constant.Doctor})
	mwDoctorAdmin := h.Middleware.AuthMulti([]string{constant.Doctor, constant.Admin})
	mwManagerAdmin := h.Middleware.AuthMulti([]string{constant.Manager, constant.Admin})

	userRouter := router.Group("/")
	{
		userRouter.POST("/register", h.UserHandler.Register)
		userRouter.POST("/register/oauth", h.UserHandler.RegisterOAuth)
		userRouter.POST("/login", h.UserHandler.Login)
		userRouter.POST("/login/oauth", h.UserHandler.LoginOAuth)
		userRouter.POST("/logout", h.UserHandler.Logout)
		userRouter.POST("/verification", h.UserHandler.Verification)
		userRouter.POST("/forgot-password", h.UserHandler.ForgotPassword)
		userRouter.POST("/reset-password", h.UserHandler.ResetPassword)

		userRouter.POST("/upload", h.UploadHandler.UploadCloudinary)
		userRouter.GET("/pharmacy-drugs", h.PharmacyDrugHandler.GetAllWithinRadius)
		userRouter.GET("/pharmacy-drugs/:id", h.PharmacyDrugHandler.GetByID)

		userRouter.GET("/provinces", h.AddressHandler.GetAllProvinces)
		userRouter.GET("/cities", h.AddressHandler.GetAllCities)
		userRouter.GET("/subdistricts", h.AddressHandler.GetAllSubdistrict)

		userRouter.GET("/specializations", h.SpecializationHandler.GetAll)
		userRouter.GET("/drugs/:id", h.DrugHandler.GetOneById)
		userRouter.GET("/drugs/:id/nearest-pharmacy", h.PharmacyDrugHandler.GetNearestPharmacyByDrugID)

		userRouter.GET("/categories", h.CategoryHandler.GetAllCategory)
		userRouter.GET("/drugs", mwDoctorAdmin, h.DrugHandler.GetAll)
		userRouter.GET("/drugs-of-the-day", h.PharmacyDrugHandler.GetDrugOfTheDay)

		userRouter.GET("/manufacturers", mwManagerAdmin, h.ManufacturerHandler.GetAllManufacturers)
		userRouter.GET("/shipment-methods", mwManagerAdmin, h.ShipmentMethodHandler.GetAllShipmentMethods)

		privateUserRouter := userRouter.Group("/")
		{
			privateUserRouter.Use(h.Middleware.UserAuth())
			privateUserRouter.GET("/profile", h.UserHandler.GetProfile)
			privateUserRouter.PUT("/profile", h.UserHandler.UpdatePersonal)
			privateUserRouter.PUT("/update-password", h.UserHandler.UpdatePassword)

			privateUserRouter.GET("/addresses", h.AddressHandler.GetAllAddress)
			privateUserRouter.POST("/addresses", h.AddressHandler.AddAddress)
			privateUserRouter.GET("/addresses/:id", h.AddressHandler.GetAddressByID)
			privateUserRouter.PUT("/addresses/:id", h.AddressHandler.UpdateAddress)
			privateUserRouter.DELETE("/addresses/:id", h.AddressHandler.DeleteAddress)
			privateUserRouter.GET("/addresses/:id/shipment-price", h.AddressHandler.GetShipmentPriceByAddressID)

			privateUserRouter.GET("/cart-items", h.CartItemHandler.GetAllCartItem)
			privateUserRouter.POST("/cart-items", h.CartItemHandler.CreateCartItem)
			privateUserRouter.PUT("/cart-items/:id", h.CartItemHandler.UpdateQtyCartItem)
			privateUserRouter.DELETE("/cart-items", h.CartItemHandler.DeleteManyCartItem)

			privateUserRouter.POST("/resend-verification", h.UserHandler.ResendVerification)

			privateUserRouter.POST("/orders", h.OrderHandler.CreateOrder)
			privateUserRouter.PATCH("/orders/:id/confirm-order", h.OrderHandler.UpdateConfirmOrder)
			privateUserRouter.GET("/payments", h.PaymentHandler.GetAllPaymentByUserId)
			privateUserRouter.PATCH("/payments/:id/update-payment-proof", h.PaymentHandler.UpdatePaymentProof)
			privateUserRouter.PATCH("/payments/:id/cancel-payment", h.PaymentHandler.UserCancelPayment)
		}
	}

	chatRouter := router.Group("/telemedicines")
	{
		privateChatRouter := chatRouter.Group("")
		{
			privateChatRouter.Use(mwUserDoctor)
			privateChatRouter.GET("", h.TelemedicineHandler.GetAllTelemedicine)
			privateChatRouter.GET("/:id", h.TelemedicineHandler.GetTelemedicineByID)
			privateChatRouter.POST("/chat", h.MessageBubbleHandler.AddChat)
		}

		privateDoctorChatRouter := chatRouter.Group("")
		{
			privateDoctorChatRouter.Use(h.Middleware.DoctorAuth())
			privateDoctorChatRouter.PUT("/:id", h.TelemedicineHandler.UpdateTelemedicine)
			privateDoctorChatRouter.POST("/:id/certificate", h.TelemedicineHandler.GenerateMedicalCertificate)
			privateDoctorChatRouter.POST("/:id/prescriptions", h.TelemedicineHandler.AddPrescriptions)
		}

		privateUserChatRouter := chatRouter.Group("")
		{
			privateUserChatRouter.Use(h.Middleware.UserAuth())
			privateUserChatRouter.POST("", h.TelemedicineHandler.AddTelemedicine)

		}
	}

	doctorRouter := router.Group("/doctor")
	{
		doctorRouter.POST("/register", h.DoctorHandler.Register)
		doctorRouter.POST("/register/oauth", h.DoctorHandler.RegisterOAuth)
		doctorRouter.POST("/login", h.DoctorHandler.Login)
		doctorRouter.POST("/login/oauth", h.DoctorHandler.LoginOAuth)
		doctorRouter.POST("/forgot-password", h.DoctorHandler.ForgotPassword)
		doctorRouter.POST("/verification", h.DoctorHandler.Verification)
		doctorRouter.GET("", h.DoctorHandler.GetAllDoctors)
		doctorRouter.GET("/:id", h.DoctorHandler.GetDoctorByID)

		privateDoctorRouter := doctorRouter.Group("/")
		{
			privateDoctorRouter.Use(h.Middleware.DoctorAuth())
			privateDoctorRouter.GET("/profile", h.DoctorHandler.GetProfile)
			privateDoctorRouter.PUT("/profile", h.DoctorHandler.UpdatePersonal)
			privateDoctorRouter.PUT("/update-password", h.DoctorHandler.UpdatePassword)
			privateDoctorRouter.PUT("/update-status", h.DoctorHandler.UpdateStatus)
		}
	}

	pharmacyManagerRouter := router.Group("/pharmacy-manager")
	{
		pharmacyManagerRouter.POST("/login", h.PharmacyManagerHandler.Login)
		privateManagerRouter := pharmacyManagerRouter.Group("/")
		{
			privateManagerRouter.Use(h.Middleware.ManagerAuth())
			privateManagerRouter.GET("/profile", h.PharmacyManagerHandler.GetProfile)
			privateManagerRouter.GET("/drugs/:id", h.PharmacyDrugHandler.GetAllPharmacyDrugs)
			privateManagerRouter.PUT("/drugs/edit/:id", h.PharmacyDrugHandler.UpdatePharmacyDrug)
			privateManagerRouter.POST("/drugs/insert/:id", h.PharmacyDrugHandler.CreatePharmacyDrug)
			privateManagerRouter.GET("/orders", h.OrderHandler.GetOrderByPharmacyManagerId)
			privateManagerRouter.PATCH("/orders/:id/order-proceed", h.OrderHandler.OrderProceed)
			privateManagerRouter.PATCH("/orders/:id/sent", h.OrderHandler.OrderSent)
			privateManagerRouter.PATCH("/orders/:id/cancel", h.OrderHandler.OrderCancelByPM)

			privateManagerRouter.POST("/drugs/insert", h.PharmacyDrugHandler.CreatePharmacyDrug)
			privateManagerRouter.POST("/stock-mutation/request", h.StockRequestHandler.StockMutationManualRequest)
			privateManagerRouter.POST("/stock-mutation/:id/approve", h.StockRequestHandler.StockMutationManualApprove)
			privateManagerRouter.POST("/stock-mutation/:id/cancel", h.StockRequestHandler.StockMutationManualCancel)

			privateManagerRouter.GET("/shipment-methods", h.ShipmentMethodHandler.GetAllShipmentMethods)
			privateManagerRouter.GET("/pharmacies/:id/stock-journal", h.StockJournalHandler.GetAllStockJournalByPharmacyId)
			privateManagerRouter.GET("/stock-requests", h.StockRequestHandler.GetAllStockRequest)

			privateManagerRouter.POST("/stock-requests/drugs", h.StockRequestHandler.GetAvailableDrugFromSenderAndReceiverPharmacy)

			privateManagerRouter.GET("/reports/drugs", h.AdminReportHandler.GetManagerDrugReport)
			privateManagerRouter.GET("/reports/categories", h.AdminReportHandler.GetManagerCategoryReport)
		}
	}

	pharmacyRouter := router.Group("/pharmacies")
	{
		pharmacyRouter.GET("", mwManagerAdmin, h.PharmacyHandler.GetAllPharmacies)
		pharmacyRouter.POST("", h.Middleware.AdminAuth(), h.PharmacyHandler.AddPharmacy)
		pharmacyRouter.GET("/:id", mwManagerAdmin, h.PharmacyHandler.GetPharmacyByID)
		pharmacyRouter.PUT("/:id", h.Middleware.ManagerAuth(), h.PharmacyHandler.EditPharmacy)
	}

	adminRouter := router.Group("/admin")
	{
		adminRouter.POST("/login", h.AdminHandler.Login)

		privateAdminRouter := adminRouter.Group("/")
		{
			privateAdminRouter.Use(h.Middleware.AdminAuth())
			privateAdminRouter.GET("/profile", h.AdminHandler.GetProfile)
			privateAdminRouter.POST("/drugs", h.DrugHandler.InsertOne)
			privateAdminRouter.PUT("/drugs/:id", h.DrugHandler.UpdateOne)

			privateAdminRouter.GET("/partners", h.PartnerHandler.GetAll)
			privateAdminRouter.POST("/partners", h.PartnerHandler.CreatePartner)
			privateAdminRouter.GET("/partners/:id", h.PartnerHandler.GetPartnerByID)
			privateAdminRouter.PUT("/partners/:id", h.PartnerHandler.UpdatePartnerByID)

			privateAdminRouter.GET("/payments", h.PaymentHandler.GetAllPaymentToConfirm)
			privateAdminRouter.PATCH("/payments/:id/confirm", h.PaymentHandler.PaymentConfirmation)
			privateAdminRouter.PATCH("/payments/:id/cancel", h.PaymentHandler.AdminCancelPayment)
			privateAdminRouter.PATCH("/payments/:id/reject", h.PaymentHandler.AdminRejectPayment)

			privateAdminRouter.POST("/categories", h.CategoryHandler.CreateCategory)
			privateAdminRouter.GET("/categories/:id", h.CategoryHandler.GetCategoryByID)
			privateAdminRouter.PUT("/categories/:id", h.CategoryHandler.UpdateCategoryByID)
			privateAdminRouter.DELETE("/categories/:id", h.CategoryHandler.DeleteCategoryByID)

			privateAdminRouter.GET("/registrants/user", h.AdminHandler.GetAllUser)
			privateAdminRouter.GET("/registrants/doctor", h.AdminHandler.GetAllDoctor)
			privateAdminRouter.GET("/registrants/pharmacy-manager", h.AdminHandler.GetAllManager)

			privateAdminRouter.GET("pharmacy-managers/:id/pharmacies", h.PharmacyHandler.GetAllPharmaciesByManagerID)

			privateAdminRouter.GET("/reports/drugs", h.AdminReportHandler.GetAdminDrugReport)
			privateAdminRouter.GET("/reports/categories", h.AdminReportHandler.GetAdminCategoryReport)
		}
	}

	router.NoRoute(h.CustomHandler.NoRoute)
	return router
}
