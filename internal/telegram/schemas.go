package telegram

type telegramResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId      int            `json:"update_id"`
	Message       Message        `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type Message struct {
	MessageId int     `json:"message_id"`
	Chat      Chat    `json:"chat"`
	Text      string  `json:"text"`
	Sticker   Sticker `json:"sticker"`
}

type Sticker struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type CallbackQuery struct {
	ID      string  `json:"id"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

type inlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}

type inlineKeyboardMarkup struct {
	InlineKeyboard [][]inlineKeyboardButton `json:"inline_keyboard"`
}

type sendMessage struct {
	ChatId      int                  `json:"chat_id"`
	Text        string               `json:"text"`
	ParseMode   string               `json:"parse_mode"`
	ReplyMarkup inlineKeyboardMarkup `json:"reply_markup,omitempty"`
}
