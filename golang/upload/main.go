package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

const (
	batchLines = 100
)

func uploadFile(fileName string) {
	language := strings.TrimSuffix(fileName, ".csv")
	log.Printf("uploading %s", language)
	languagePath := fmt.Sprintf("/mnt/efs/language/%s", language)
	err := os.MkdirAll(languagePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	inFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	batchId := 0
	lineCount := 0
	fileName = fmt.Sprintf("%s/%d", languagePath, batchId)
	log.Printf("creating %s", fileName)
	outFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()
	datawriter := bufio.NewWriter(outFile)

	for scanner.Scan() {
		datawriter.WriteString(scanner.Text())
		datawriter.WriteString("\n")
		lineCount++
		if lineCount%batchLines == 0 {
			datawriter.Flush()
			outFile.Close()
			batchId++
			fileName = fmt.Sprintf("%s/%d", languagePath, batchId)
			log.Printf("creating %s", fileName)
			outFile, err = os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}
			datawriter = bufio.NewWriter(outFile)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	err := os.MkdirAll("/mnt/efs/poll", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll("/mnt/efs/chat", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".csv") {
			uploadFile(file.Name())
		}
	}

	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
