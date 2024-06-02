package entity

type Collection struct {
	Filter       map[string]string
	Sort         string
	Search       string
	Page         uint
	Limit        uint
	TotalRecords uint
	Args         []any
}
