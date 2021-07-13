package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Scenario struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BotID    string             `json:"botId" bson:"botId"`
	Name     string             `json:"name" bson:"name"`
	Triggers []string           `json:"triggers" bson:"triggers,omitempty"`
	Actions  []Action           `json:"actions" bson:"actions,omitempty"`
}
