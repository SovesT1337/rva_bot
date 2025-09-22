package telegram

import (
	"bytes"
	"encoding/json"
	"html"
	"io"
	"log"
	"net/http"
)

// logResponse логирует HTTP ответ в понятном формате
func logResponse(operation string, resp *http.Response) {
	if resp == nil {
		log.Printf("%s: получен пустой ответ", operation)
		return
	}

	statusText := "успешно"
	if resp.StatusCode >= 400 {
		statusText = "ошибка"
	}

	log.Printf("%s: %s (код: %d)", operation, statusText, resp.StatusCode)

	// Читаем тело ответа для отладки ошибок
	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err == nil && len(body) > 0 {
			log.Printf("%s: тело ответа: %s", operation, string(body))
		}
	}
}

func LogOut(botUrl string) error {
	responce, err := http.Post(botUrl+"/logout", "application/json", nil)
	if err != nil {
		log.Printf("Выход из бота: ошибка соединения - %s", err)
		return err
	}

	logResponse("Выход из бота", responce)
	return nil
}

func SendMessage(botUrl string, chatId int, text string, keyboard inlineKeyboardMarkup) error {
	message := sendMessage{
		ChatId:      chatId,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: keyboard,
	}

	buf, err := json.Marshal(message)
	if err != nil {
		log.Printf("Отправка сообщения: ошибка сериализации - %s", err)
		return err
	}

	responce, err := http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("Отправка сообщения: ошибка соединения - %s", err)
		return err
	}

	logResponse("Отправка сообщения", responce)
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
		log.Printf("Редактирование сообщения: ошибка соединения - %s", err)
		return err
	}

	logResponse("Редактирование сообщения", responce)

	// Если редактирование сообщения не удалось, отправляем новое сообщение
	if responce.StatusCode >= 400 {
		log.Printf("Редактирование сообщения не удалось (код: %d), отправляем новое сообщение", responce.StatusCode)
		return SendMessage(botUrl, chatID, text, keyboard)
	}

	return nil
}

func AnswerCallbackQuery(botUrl string, callbackID string) error {

	body := map[string]string{"callback_query_id": callbackID}
	jsonBody, _ := json.Marshal(body)
	responce, err := http.Post(botUrl+"/answerCallbackQuery", "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Printf("Ответ на callback: ошибка соединения - %s", err)
		return err
	}

	logResponse("Ответ на callback", responce)
	return nil
}

func EscapeHTML(text string) string {
	return html.EscapeString(text)
}
