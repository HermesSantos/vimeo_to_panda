package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 55ee21e741241e541ade85500013e77d
type File struct {
	Link string `json:"link"`
}

type VimeoResponse struct {
	Files []File `json:"files"`
}

func Get_vimeo_video (video_id string) (string, error) {
	url := "https://api.vimeo.com/videos/" + video_id + "?time_links=false"
	url_token := "55ee21e741241e541ade85500013e77d"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Erro creating new request: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer " + url_token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error executing the request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Erro ao ler corpo: %v", err)
	}

	var data VimeoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("Erro ao deserializar JSON: %v", err)
	}

	return data.Files[0].Link, nil
}
