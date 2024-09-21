package main

type Product struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	PricePerLb float64 `json:"price_per_lb"`
	PricePerOz float64 `json:"price_per_oz"`
	Store      string  `json:"store"`
	Volume     string  `json:"volume"`
	Weight     float64 `json:"weight"`
	OnSale     bool    `json:"on_sale"`
	DateTime   string  `json:"datetime"`
}
