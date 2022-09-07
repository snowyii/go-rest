package model

type Product struct {
	Name        string  `json:"name,omitempty bson:"name",omitempty"`
	Price       float64 `json:"price,omitempty bson:"price",omitempty"`
	Description string  `json:"description,omitempty bson:"description,omitempty""`
}

type User struct {
	Name    string
	Surname string
	Age     int
	Balance float64
}
