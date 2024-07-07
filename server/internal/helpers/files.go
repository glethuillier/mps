package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/glethuillier/fvs/server/internal/logger"
	"go.uber.org/zap"
)

// filesDir is the directory where files uploaded by the client
// are stored and retrieved
const filesDir = "downloads"

func Init() error {
	return ensureDirectory(filesDir)
}

// ensureDirectory ensures that the directory exists and is writable
func ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create folder: %s", err)
		}
	}

	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("unable to stat folder: %s", err)
	}

	mode := info.Mode()
	if mode&0200 == 0 {
		return fmt.Errorf("folder is not writable")
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if stat.Uid != uint32(os.Geteuid()) {
			if mode&0020 == 0 && mode&0002 == 0 {
				return fmt.Errorf("folder is not writable")
			}
		}
	}

	return nil
}

func SaveFile(id string, filename string, content []byte) {
	dir := filepath.Join(filesDir, id)
	err := ensureDirectory(dir)
	if err != nil {
		logger.Logger.Fatal("directory cannot be accessed", zap.Error(err))
		return
	}

	f := filepath.Join(dir, filename)
	logger.Logger.Debug("saving file", zap.String("filepath", f))

	err = os.WriteFile(f, content, 0644)
	if err != nil {
		logger.Logger.Fatal("cannot write file", zap.Error(err))
		return
	}
}

func GetFile(id string, filename string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(filesDir, id, filename))
	if err != nil {
		return nil, err
	}

	return data, nil
}
