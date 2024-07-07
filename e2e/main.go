package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/glethuillier/mps/e2e/internal/test"
)

const (
	clientBaseUrl = "http://localhost:3001"
	testDirectory = "tmp"
)

func main() {
	var err error

	gofakeit.Seed(0)
	test.CreateTestDirectory(testDirectory)
	defer func() {
		test.DeleteTestDirectory(testDirectory)
	}()

	filesCount := 100
	if len(os.Args) > 1 && os.Args[1] != "" {
		if filesCount, err = strconv.Atoi(os.Args[1]); err != nil {
			panic(err)
		}
	}

	filenames := test.GenerateFiles(testDirectory, filesCount)

	fmt.Printf("Generated %d test files\n", len(filenames))

	t1 := time.Now()
	receiptId, err := test.UploadFiles(
		fmt.Sprintf("%s/%s", clientBaseUrl, "upload"),
		testDirectory,
		filenames,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Files successfully uploaded (%s)\n", time.Since(t1))

	fmt.Println("Press the Enter Key to proceed to the verification")
	fmt.Println("(you can manually corrupt a file to ensure that the solution detects it)")
	fmt.Scanln()

	t2 := time.Now()
	err = test.DownloadFiles(
		fmt.Sprintf("%s/%s", clientBaseUrl, "download"),
		receiptId,
		filenames,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Files successfully verified (%s)\n", time.Since(t2))
}
