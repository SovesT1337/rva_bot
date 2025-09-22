package telegram

import (
	"bytes"
	"encoding/json"
	"html"
	"log"
	"net/http"
)

func SendMessage(botUrl string, chatId int, text string, keyboard inlineKeyboardMarkup) error {
	message := sendMessage{
		ChatId:      chatId,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: keyboard,
	}

	buf, err := json.Marshal(message)
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
		return err
	}

	responce, err := http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("SendMessage error: %s", err)
		return err
	}

	log.Println("Responce: ", responce)

	return nil

}

func EditMessage(botUrl string, chatID int, messageID int, text string, keyboard inlineKeyboardMarkup) error {

	body := map[string]interface{}{
		"chat_id":      chatID,
		"message_id":   messageID,
		"text":         text,
		"parse_mode":   "HTML",
		"reply_markup": keyboard,
	}
	jsonBody, _ := json.Marshal(body)
	responce, err := http.Post(botUrl+"/editMessageText", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("EditMessage error: ", err)
		return err
	}

	log.Println("Responce: ", responce)

	// Если редактирование сообщения не удалось, отправляем новое сообщение
	if responce.StatusCode >= 400 {
		log.Printf("EditMessage failed with status %d, sending new message instead", responce.StatusCode)
		return SendMessage(botUrl, chatID, text, keyboard)
	}

	return nil
}

func AnswerCallbackQuery(botUrl string, callbackID string) error {

	body := map[string]string{"callback_query_id": callbackID}
	jsonBody, _ := json.Marshal(body)
	responce, err := http.Post(botUrl+"/answerCallbackQuery", "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Println("AnswerCallbackQuery error: ", err)
		return err
	}

	log.Println("Responce: ", responce)
	return nil
}

func EscapeHTML(text string) string {
	return html.EscapeString(text)
}
