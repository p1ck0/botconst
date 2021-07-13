package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bot struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID              string             `json:"userId" bson:"userId,omitempty"`
	Name                string             `json:"name" bson:"name,omitempty"`
	Telegram            string             `json:"telegram" bson:"telegram,omitempty"`
	WhatsAppID          string             `json:"whatsAppId" bson:"whatsAppId,omitempty"`
	WhatsAppToken       string             `json:"whatsAppToken" bson:"whatsAppToken,omitempty"`
	FacebookAppID       string             `json:"facebookAppId" bson:"facebookAppId,omitempty"`
	FacebookAppSecret   string             `json:"facebookAppSecret" bson:"facebookAppSecret,omitempty"`
	FacebookPageID      string             `json:"facebookPageId" bson:"facebookPageId,omitempty"`
	FacebookAccessToken string             `json:"facebookAccessToken" bson:"facebookAccessToken,omitempty"`
}
