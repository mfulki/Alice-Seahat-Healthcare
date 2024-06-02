package apperror

import (
	"net/http"
)

type AppError struct {
	StatusCode int
	Err        error
}

func New(code int, err error) *AppError {
	return &AppError{
		StatusCode: code,
		Err:        err,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	Unauthorized                      = New(http.StatusUnauthorized, ErrUnauthorized)
	InvalidCredential                 = New(http.StatusBadRequest, ErrInvalidCredential)
	InvalidToken                      = New(http.StatusBadRequest, ErrInvalidToken)
	InvalidEmail                      = New(http.StatusBadRequest, ErrInvalidEmail)
	InvalidPassword                   = New(http.StatusBadRequest, ErrInvalidPassword)
	InvalidPassToken                  = New(http.StatusBadRequest, ErrInvalidPassToken)
	InvalidParam                      = New(http.StatusBadRequest, ErrInvalidParam)
	ResourceNotFound                  = New(http.StatusBadRequest, ErrResourceNotFound)
	DrugExist                         = New(http.StatusBadRequest, ErrDrugExist)
	PharmacyDrugNotExist              = New(http.StatusBadRequest, ErrPharmacyDrugNotExist)
	PharmacyDrugCantEdit              = New(http.StatusBadRequest, ErrPharmacyDrugCantEdit)
	EmailExist                        = New(http.StatusBadRequest, ErrEmailExist)
	EmailCantResetPassword            = New(http.StatusBadRequest, ErrEmailCantResetPassword)
	FileTooLarge                      = New(http.StatusBadRequest, ErrFileTooLarge)
	FileInvalidType                   = New(http.StatusBadRequest, ErrFileInvalidType)
	MustMainAddress                   = New(http.StatusBadRequest, ErrMustMainAddress)
	SubdistrictNotAvail               = New(http.StatusBadRequest, ErrSubdistrictNotAvail)
	MainMustActiveAddress             = New(http.StatusBadRequest, ErrMainMustActiveAddress)
	CantDeleteMainAddress             = New(http.StatusBadRequest, ErrCantDeleteMainAddress)
	HasBeenVerified                   = New(http.StatusBadRequest, ErrHasBeenVerified)
	PrescriptedExist                  = New(http.StatusBadRequest, ErrPrescriptedExist)
	InsufficientStock                 = New(http.StatusBadRequest, ErrInsufficientStock)
	InsufficientStockMutation         = New(http.StatusBadRequest, ErrInsufficientStockMutation)
	NoValidCartOrder                  = New(http.StatusBadRequest, ErrNoValidCartInsideOrder)
	InvalidIdParams                   = New(http.StatusBadRequest, ErrInvalidIdParams)
	DoctorNotExist                    = New(http.StatusBadRequest, ErrDoctorNotExist)
	NoValidOrderPayment               = New(http.StatusBadRequest, ErrNoValidOrder)
	NoValidPayment                    = New(http.StatusBadRequest, ErrNoValidPayment)
	AddressNotExist                   = New(http.StatusBadRequest, ErrAddressNotExist)
	InvalidShipmentMethods            = New(http.StatusBadRequest, ErrShipmentMethodInvalid)
	TelemedicineHasBeenEnded          = New(http.StatusBadRequest, ErrTelemedicineHasBeenEnded)
	TelemedicineOngoingWithDoctor     = New(http.StatusBadRequest, ErrTelemedicineOngoingWithDoctor)
	CategoryNameExist                 = New(http.StatusBadRequest, ErrCategoryNameExist)
	CantCancelOrder                   = New(http.StatusBadRequest, ErrCantCancelOrder)
	CantRequestToSamePharmacy         = New(http.StatusBadRequest, ErrCantRequestToSamePharmacy)
	CantRequestToPharmaciesNotPartner = New(http.StatusBadRequest, ErrCantRequestToPharmaciesNotPartner)
	PharmacyDrugExist                 = New(http.StatusBadRequest, ErrPharmacyDrugExist)
	DrugNotExist                      = New(http.StatusBadRequest, ErrDrugNotExist)
	PharmacyNotExist                  = New(http.StatusBadRequest, ErrPharmacyNotExist)
	CategoryNotExist                  = New(http.StatusBadRequest, ErrCategoryNotExist)
	EmailOAuthNotFound                = New(http.StatusBadRequest, ErrEmailOAuthNotFound)
	SenderAndReceiverCantSame         = New(http.StatusBadRequest, ErrSenderAndReceiverCantSame)
	CantCancelStockMutation           = New(http.StatusBadRequest, ErrStockMutationCancel)
	CantApproveStockMutation          = New(http.StatusBadRequest, ErrStockMutationApprove)
	CantDeleteCategory                = New(http.StatusBadRequest, ErrCantDeleteCategory)
)

var (
	NoRoute = New(http.StatusNotFound, ErrNoRoute)
)
