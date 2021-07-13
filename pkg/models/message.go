package models

type Message struct {
	BotID    string `json:"botId"`
	ChatID   string `json:"chatId"`
	Platform string `json:"platform"`
	Text     string `json:"text"`
}
