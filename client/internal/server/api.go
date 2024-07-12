package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/glethuillier/mps/client/internal/common"
	"github.com/glethuillier/mps/client/internal/logger"
	"github.com/glethuillier/mps/client/internal/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type serverResponse struct {
	ReceiptId string `json:"receiptId,omitempty"`
	Error     string `json:"error,omitempty"`
}

type downloadRequest struct {
	ReceiptId string `json:"receipt_id"`
	Filename  string `json:"filename"`
}

// uploadFilesHandler handles requests to upload a batch of files
func uploadFilesHandler(
	ctx context.Context,
	service *middleware.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20) // limit to 32Mb
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			logger.Logger.Error(
				"unable to parse request form",
				zap.Error(err),
			)
			return
		}

		form := r.MultipartForm
		files := form.File["file"]

		var uploadedFiles []common.File

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Unable to open file", http.StatusInternalServerError)
				return
			}
			file.Close()

			// get contents
			contents, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Unable to get file contents", http.StatusInternalServerError)
			}

			uploadedFiles = append(uploadedFiles, common.File{
				Filename: fileHeader.Filename,
				Contents: contents,
			})
		}

		requestID := uuid.New()
		receiptId, err := service.ProcessUploadRequest(
			ctx,
			requestID,
			common.UploadRequest{
				Files: uploadedFiles,
			})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(serverResponse{Error: err.Error()})
			if err != nil {
				logger.Logger.Error(
					"cannot send error",
					zap.Error(err),
				)
			}
		} else {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(serverResponse{ReceiptId: receiptId})
			if err != nil {
				logger.Logger.Error(
					"cannot send response",
					zap.Error(err),
				)
			}
		}
	}
}

// downloadFilesHandler handles requests to download a given file
func downloadFilesHandler(ctx context.Context, service *middleware.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req downloadRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		requestID := uuid.New()

		file, err := service.ProcessDownloadRequest(
			ctx,
			requestID,
			common.DownloadRequest{
				ReceiptId: req.ReceiptId,
				Filename:  req.Filename,
			},
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(serverResponse{Error: file.Error.Error()})
			if err != nil {
				logger.Logger.Error(
					"cannot send error",
					zap.Error(err),
				)
			}
		} else {
			// return file
			w.Header().Set("Content-Disposition", "attachment; filename="+file.Filename)
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(file.Contents)))

			// custom proof-related headers
			for i, p := range file.Proof {
				w.Header().Set(
					fmt.Sprintf("Proof-Sibling-%d-%s", i, p.SiblingType.String()),
					p.SiblingHash,
				)
			}

			w.Header().Set("Proof-Root-Hash", req.ReceiptId)

			w.Write(file.Contents)
		}
	}
}

func Run(ctx context.Context, service *middleware.Service) {
	http.HandleFunc("/upload", uploadFilesHandler(ctx, service))
	http.HandleFunc("/download", downloadFilesHandler(ctx, service))

	logger.Logger.Info("API server started at :3001")
	http.ListenAndServe("0.0.0.0:3001", nil)
}
