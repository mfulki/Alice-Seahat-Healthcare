package constant

import "time"

const (
	Production      = "production"
	DefaultPhotoURL = "https://res.cloudinary.com/aliceseahat/image/upload/v1713793225/static-assets/default-user.png"

	TimeoutShutdown = 5 * time.Second

	MailSubjectVerification = "Konfirmasi email aplikasi Seahat"
	MailSubjectReset        = "Reset password akun aplikasi Seahat"
	MailSubjectNewPartner   = "Selamat, akun anda terdaftar sebagai manager farmasi"

	StatusOnline  = "online"
	StatusOffline = "offline"

	User    = "user"
	Doctor  = "doctor"
	Manager = "manager"
	Admin   = "admin"

	LengthOfRequestID = 15
	MetreToKilometre  = 1000
)
