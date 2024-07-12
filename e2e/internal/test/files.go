package test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
)

func CreateTestDirectory(testDirectory string) error {
	if _, err := os.Stat(testDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(testDirectory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create folder: %v", err)
		}
	}

	info, err := os.Stat(testDirectory)
	if err != nil {
		return fmt.Errorf("unable to stat folder: %v", err)
	}

	mode := info.Mode()
	if mode&0200 == 0 {
		return fmt.Errorf("folder is not writable")
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if stat.Uid != uint32(os.Geteuid()) {
			if mode&0020 == 0 && mode&0002 == 0 {
				return fmt.Errorf("folder is not writable by group or others")
			}
		}
	}

	return nil
}

func DeleteTestDirectory(testDirectory string) error {
	return os.RemoveAll(testDirectory)
}

func randomFileSize(min, max int64) (int64, error) {
	n := max - min

	randInt, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		return 0, err
	}

	return randInt.Int64() + min, nil
}

// GenerateFiles generates random files and
// returns their filenames
func GenerateFiles(testDirectory string, filesCount int) []string {
	var filenames []string

	for i := 0; i < filesCount; i++ {
		filename := fmt.Sprintf(
			"%s.%s",
			gofakeit.Regex("[a-zA-Z0-9_]{5,20}"),
			gofakeit.FileExtension(),
		)

		fileSize, err := randomFileSize(10, 1000)
		if err != nil {
			panic(err)
		}

		contents := make([]byte, fileSize)
		rand.Read(contents)

		filepath := path.Join(testDirectory, filename)

		err = os.WriteFile(filepath, contents, 0644)
		if err != nil {
			panic(err)
		}

		filenames = append(filenames, filename)
	}

	return filenames
}
