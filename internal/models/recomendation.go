package models

type Recom struct {
	ID   string `bson:"_id,omitempty"`
	Text string `bson:"text"`
}
