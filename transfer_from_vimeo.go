package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PandaPayload struct {
	FolderID    string `json:"folder_id"`
	VideoID     string `json:"video_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type PandaResponse struct {
	WebsocketUrl string `json:"websocket_url"`
}

type APIError struct {
	Code    int    `json:"code"`
	ErrCode string `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

func Transfer_from_vimeo (download_url string) (string, error) {

	var pandaPayload PandaPayload = PandaPayload {
		FolderID: "1655686b-20dd-4fa3-a3cd-67bd33fe072b",
		Title: "Title teste",
		Description: "description teste",
		URL: download_url,
	}

	jsonPandaData, err := json.Marshal(pandaPayload)

	if err != nil {
		return "", fmt.Errorf("Erro ao serializar: %v", err)
	}


	url := "https://import.pandavideo.com:9443/videos/"

	// payload := strings.NewReader("{\"folder_id\":\"string\",\"video_id\":\"string\",\"title\":\"string\",\"description\":\"string\",\"url\":\"string\",\"size\":\"string\"}")
	payload := strings.NewReader(string(jsonPandaData))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Set("Authorization", "panda-06510f055590d747d8089df80adfcd7de4a098550ba44ab2e4a5f04cb78aa989") // Substitua pelo token real

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao executar requisição: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler corpo da resposta: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return "", fmt.Errorf("erro ao decodificar erro da API (status %d): %w", res.StatusCode, err)
		}
		return "", fmt.Errorf("erro da API (status %d): %s - %s", apiErr.Code, apiErr.ErrCode, apiErr.ErrMsg)
	}

	var response PandaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	return response.WebsocketUrl, nil

}
