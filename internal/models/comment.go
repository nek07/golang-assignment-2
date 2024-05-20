package models

type Comment struct {
	ID       string `bson:"_id,omitempty"`
	LaptopID string `bson:"laptopId"`
	UserName string `bson:"userName"`
	Text     string `bson:"text"`
	Time     string `bson:"time"`
}
