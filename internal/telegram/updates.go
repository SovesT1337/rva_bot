package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func GetUpdates(botUrl string, offset int) ([]Update, error) {

	resp, err := http.Get(botUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse telegramResponse
	if err := json.Unmarshal(body, &restResponse); err != nil {
		return nil, err
	}

	return restResponse.Result, nil

}