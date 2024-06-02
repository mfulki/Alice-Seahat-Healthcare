package rajaongkir

type CostPayload struct {
	Origin      uint   `json:"origin"`
	Destination uint   `json:"destination"`
	Weight      uint   `json:"weight"`
	Courier     string `json:"courier"`
}

type CostResponse struct {
	RajaOngkir rajaOngkirResponse `json:"rajaongkir"`
}

type rajaOngkirResponse struct {
	Results []costCourier `json:"results"`
}

type costCourier struct {
	Code  string           `json:"code"`
	Name  string           `json:"name"`
	Costs []serviceCourier `json:"costs"`
}

type serviceCourier struct {
	Service     string               `json:"service"`
	Description string               `json:"description"`
	Cost        []costServiceCourier `json:"cost"`
}

type costServiceCourier struct {
	Value float64 `json:"value"`
	Etd   string  `json:"etd"`
	Note  string  `json:"note"`
}
