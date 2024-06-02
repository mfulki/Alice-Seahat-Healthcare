package response

import "Alice-Seahat-Healthcare/seahat-be/entity"

type PharmacyWithShipmentPrice struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	Shipments []ShipmentPriceDto `json:"shipments"`
}

type ShipmentPriceDto struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	CourierName string `json:"courier_name"`
	Price       uint   `json:"price"`
}

func NewPharmacyWithShipmentPrice(p entity.Pharmacy) PharmacyWithShipmentPrice {
	return PharmacyWithShipmentPrice{
		ID:        p.ID,
		Name:      p.Name,
		Shipments: NewMultipleShipmentPriceDto(p.ShipmentMethods),
	}
}

func NewMultiplePharmacyWithShipmentPrice(ps []*entity.Pharmacy) []PharmacyWithShipmentPrice {
	pharmacyDtos := make([]PharmacyWithShipmentPrice, 0)

	for _, p := range ps {
		pharmacyDtos = append(pharmacyDtos, NewPharmacyWithShipmentPrice(*p))
	}

	return pharmacyDtos
}

func NewShipmentPriceDto(sm entity.ShipmentMethod) ShipmentPriceDto {
	return ShipmentPriceDto{
		ID:          sm.ID,
		Name:        sm.Name,
		CourierName: sm.CourierName,
		Price:       *sm.Price,
	}
}

func NewMultipleShipmentPriceDto(sms []*entity.ShipmentMethod) []ShipmentPriceDto {
	spDtos := make([]ShipmentPriceDto, 0)

	for _, sm := range sms {
		spDtos = append(spDtos, NewShipmentPriceDto(*sm))
	}

	return spDtos
}
