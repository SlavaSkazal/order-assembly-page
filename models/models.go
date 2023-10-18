package models

type Rack struct {
	ID       int
	Name     string
	Products map[int]bool
}

type Order struct {
	ID           int
	Number       int
	ProdQauntyty map[int]int
}

type Product struct {
	ID              int
	Name            string
	MainRack        int
	AdditionalRacks map[int]bool
	Orders          map[int]bool
}
