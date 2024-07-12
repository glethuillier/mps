package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type downloadResponse struct {
	ReceiptID string `json:"receipt_id"`
	Filename  string `json:"filename"`
}

func DownloadFiles(url string, receiptId string, filenames []string) error {
	client := &http.Client{}

	for i, filename := range filenames {
		fmt.Printf("%d/%d — verifying %s: ", i+1, len(filenames), filename)

		payload := downloadResponse{
			ReceiptID: receiptId,
			Filename:  filename,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to perform request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("FAIL (file is corrupted)\n")
			return fmt.Errorf("received non-OK response: %s", resp.Status)
		} else {
			fmt.Printf("OK\n")
		}
		resp.Body.Close()
	}

	return nil
}
