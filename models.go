package main

//Telegram structs
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}
type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

// Arshin structs
type RestResponseDevices struct {
	SearchResult Result `json:"result"`
}

type Result struct {
	Devices []Device `json:"items"`
}

type Device struct {
	Organization string `json:"org_title"`
	TypeOfDevice string `json:"mit_title"`
	DeviceNumber string `json:"mi_number"`
	ValidDate    string `json:"valid_date"`
}
