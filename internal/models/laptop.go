package models

type Laptop struct {
	ID          string `bson:"_id"`
	Brand       string `bson:"brand"`
	Model       string `bson:"model"`
	Description string `bson:"description"`
	Price       int    `bson:"price"`
	Id          string
}
