package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"x.localhost/rvabot/internal/errors"
	httpclient "x.localhost/rvabot/internal/http"
	"x.localhost/rvabot/internal/logger"
	"x.localhost/rvabot/internal/validation"
)

// getHTTPClient возвращает HTTP клиент из пула
func getHTTPClient() *http.Client {
	return httpclient.GetGlobalClientPool().GetDefaultClient()
}

// logResponse логирует HTTP ответ в понятном формате
func logResponse(operation string, resp *http.Response) {
	if resp == nil {
		logger.TelegramError("Пустой ответ: %s", operation)
		return
	}

	if resp.StatusCode >= 400 {
		logger.TelegramError("%s ошибка (код: %d)", operation, resp.StatusCode)

		// Читаем тело ответа для отладки ошибок
		body, err := io.ReadAll(resp.Body)
		if err == nil && len(body) > 0 {
			logger.TelegramError("Ответ: %s", string(body))
		}
	} else {
		logger.TelegramInfo("%s успешно (код: %d)", operation, resp.StatusCode)
	}
}

// makeHTTPRequest выполняет HTTP запрос с retry логикой
func makeHTTPRequest(method, url string, body []byte) (*http.Response, error) {
	client := getHTTPClient()

	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		var req *http.Request
		var err error

		if body != nil {
			req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}

		if err != nil {
			return nil, errors.NewNetworkError("Ошибка создания запроса", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, errors.NewNetworkError("Ошибка выполнения запроса", err)
			}
			logger.TelegramWarn("Попытка %d/%d неудачна, повтор через 2 секунды: %v", attempt, maxRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}

		return resp, nil
	}

	return nil, errors.NewNetworkError("Превышено максимальное количество попыток", nil)
}

func LogOut(botUrl string) error {
	client := getHTTPClient()
	responce, err := client.Post(botUrl+"/logout", "application/json", nil)
	if err != nil {
		logger.TelegramError("Выход из бота: %s", err)
		return err
	}

	logResponse("Выход из бота", responce)
	return nil
}

func SendMessage(botUrl string, chatId int, text string, keyboard inlineKeyboardMarkup) error {
	// Валидация входных данных
	validator := validation.NewValidator()
	if result := validator.ValidateMessageText(text); !result.IsValid {
		return errors.NewValidationError("Неверный текст сообщения", strings.Join(result.GetErrorMessages(), "; "))
	}

	if chatId <= 0 {
		return errors.NewValidationError("Неверный Chat ID", "Chat ID должен быть положительным числом")
	}

	message := sendMessage{
		ChatId:      chatId,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: keyboard,
	}

	buf, err := json.Marshal(message)
	if err != nil {
		appErr := errors.NewTelegramError("Ошибка маршалинга сообщения", err)
		logger.TelegramError("Сериализация сообщения: %v", appErr)
		return appErr
	}

	resp, err := makeHTTPRequest("POST", botUrl+"/sendMessage", buf)
	if err != nil {
		logger.TelegramError("Отправка сообщения: %v", err)
		return err
	}
	defer resp.Body.Close()

	logResponse("Отправка сообщения", resp)

	if resp.StatusCode >= 400 {
		return errors.NewTelegramError("Ошибка API Telegram", nil).WithCode(fmt.Sprintf("HTTP_%d", resp.StatusCode))
	}

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
	client := getHTTPClient()
	responce, err := client.Post(botUrl+"/editMessageText", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.TelegramError("Редактирование сообщения: %s", err)
		return err
	}

	logResponse("Редактирование сообщения", responce)

	// Если редактирование сообщения не удалось, отправляем новое сообщение
	if responce.StatusCode >= 400 {
		logger.TelegramWarn("Редактирование не удалось (код: %d), отправляем новое", responce.StatusCode)
		return SendMessage(botUrl, chatID, text, keyboard)
	}

	return nil
}

func AnswerCallbackQuery(botUrl string, callbackID string) error {
	body := map[string]string{"callback_query_id": callbackID}
	jsonBody, _ := json.Marshal(body)
	client := getHTTPClient()
	responce, err := client.Post(botUrl+"/answerCallbackQuery", "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		logger.TelegramError("Ответ на callback: %s", err)
		return err
	}

	logResponse("Ответ на callback", responce)
	return nil
}

func EscapeHTML(text string) string {
	return html.EscapeString(text)
}
