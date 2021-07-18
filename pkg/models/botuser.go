package models

type BotUser struct {
	BotId     string `json:"bot_id" bson:"bot_id"`
	BotUserId string `json:"bot_user_id" bson:"bot_user_id"`
	Platform  string `json:"platform" bson:"platform"`
	State     string `json:"state" bson:"state"`
}
