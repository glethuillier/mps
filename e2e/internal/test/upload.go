package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type uploadResponse struct {
	ReceiptID string `json:"receiptId"`
}

// uploadFiles uploads the test files to the server and returns
// the receipt ID or an error
func UploadFiles(url, testDirectory string, filenames []string) (string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	for _, filename := range filenames {
		file, err := os.Open(path.Join(testDirectory, filename))
		if err != nil {
			return "", err
		}

		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return "", err
		}

		file.Close()
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"failed to upload files: %s",
			response.Status,
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var resp uploadResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		return "", err
	}

	return resp.ReceiptID, nil
}
