package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const (
	API_KEY   = "panda-0ea7c727dbffbde35fde73a8d57d087657327208ba49c789bd8a0951531653e0"
	FOLDER_ID = "" // preencher caso seja necessario com um uuid
	FILENAME  = "nome do video teste" //pegar nome do vídeo do panda
)

var VIDEO_ID = uuid.New().String() // Gera um UUID v4

// parseToBase64 converte uma string para base64
func parseToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// uploadVideo faz o upload do vídeo
func uploadVideo(filename string) error {
	// Lê o arquivo binário
	binaryFile, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("erro ao ler o arquivo: %v", err)
	}

	// Monta os metadados em base64
	metadata := fmt.Sprintf("authorization %s", parseToBase64(API_KEY))
	if FOLDER_ID != "" {
		metadata += fmt.Sprintf(", folder_id %s", parseToBase64(FOLDER_ID))
	}
	metadata += fmt.Sprintf(", filename %s", parseToBase64(FILENAME))
	metadata += fmt.Sprintf(", video_id %s", parseToBase64(VIDEO_ID))

	// Requisição para obter os servidores de upload
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api-v2.pandavideo.com.br/hosts/uploader", nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição GET: %v", err)
	}
	req.Header.Set("Authorization", API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro na requisição GET: %v", err)
	}
	defer resp.Body.Close()

	// Decodifica a resposta JSON
	var uploadServers struct {
		Hosts map[string][]string `json:"hosts"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler resposta: %v", err)
	}
	if err := json.Unmarshal(body, &uploadServers); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	// Junta todos os hosts em uma única lista
	var allHosts []string
	for _, hosts := range uploadServers.Hosts {
		allHosts = append(allHosts, hosts...)
	}
	fmt.Println("Hosts disponíveis:", allHosts)

	// Escolhe um host aleatoriamente
	host := allHosts[rand.Intn(len(allHosts))]
	fmt.Printf("Iniciando upload para %s\n", host)

	// Requisição POST para fazer o upload do vídeo
	url := fmt.Sprintf("https://%v.pandavideo.com.br/files", host)
	req, err = http.NewRequest("POST", url, bytes.NewReader(binaryFile))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição POST: %v", err)
	}

	req.Header.Set("Tus-Resumable", "1.0.0")
	req.Header.Set("Upload-Length", fmt.Sprintf("%d", len(binaryFile)))
	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Upload-Metadata", metadata)

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer upload: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		fmt.Println("Upload concluído com sucesso")
	} else {
		return fmt.Errorf("erro no upload, status: %d", resp.StatusCode)
	}

	return nil
}

func Upload_from_local() {
	err := uploadVideo("certo.mp4")
	if err != nil {
		fmt.Println("UPLOAD ERROR")
		fmt.Println(err)
	}
}
