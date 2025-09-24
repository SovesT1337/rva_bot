package telegram

import (
	"encoding/json"
	"io"
	"strconv"

	"x.localhost/rvabot/internal/errors"
	httpclient "x.localhost/rvabot/internal/http"
	"x.localhost/rvabot/internal/logger"
)

func GetUpdates(botUrl string, offset int) ([]Update, error) {
	// Валидация входных данных
	if offset < 0 {
		return nil, errors.NewValidationError("Неверный offset", "offset должен быть неотрицательным числом")
	}

	client := httpclient.GetGlobalClientPool().GetDefaultClient()

	url := botUrl + "/getUpdates?offset=" + strconv.Itoa(offset)
	resp, err := client.Get(url)
	if err != nil {
		appErr := errors.NewNetworkError("Ошибка получения обновлений", err)
		logger.TelegramError("Ошибка HTTP запроса: %v", appErr)
		return nil, appErr
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		appErr := errors.NewTelegramError("Ошибка API Telegram", nil).WithCode(strconv.Itoa(resp.StatusCode))
		logger.TelegramError("Ошибка API (код %d)", resp.StatusCode)
		return nil, appErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		appErr := errors.NewNetworkError("Ошибка чтения ответа", err)
		logger.TelegramError("Ошибка чтения тела ответа: %v", appErr)
		return nil, appErr
	}

	var restResponse telegramResponse
	if err := json.Unmarshal(body, &restResponse); err != nil {
		appErr := errors.NewTelegramError("Ошибка парсинга JSON", err)
		logger.TelegramError("Ошибка парсинга ответа: %v", appErr)
		return nil, appErr
	}

	logger.TelegramInfo("Получено %d обновлений", len(restResponse.Result))
	return restResponse.Result, nil
}
