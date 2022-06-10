package telegram

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

type FileReader interface {
	ReadFile(string) ([]byte, error)
}

func GetWords(
	fileReader FileReader,
	language string,
	batchId int,
) ([]string, int, error) {
	if batchId == -1 {
		batchId = rand.Intn(260)
		log.Printf("random batchId: %d", batchId)
	}

	key := fmt.Sprintf("/mnt/efs/language/%s/%d", language, batchId)
	content, err := fileReader.ReadFile(key)
	if err != nil {
		return nil, 0, err
	}

	return strings.Split(string(content), "\n"), batchId, nil
}
