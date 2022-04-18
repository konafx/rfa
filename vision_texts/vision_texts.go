package vision_texts

import (
	"context"
	"fmt"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"go.uber.org/zap"
)

func Detect(filename string) string {
	// Creates a client.
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		zap.L().Fatal("Failed to create client", zap.Error(err))
	}
	defer client.Close()

	file, err := os.Open(filename)
	if err != nil {
		zap.L().Fatal("Failed to read file", zap.Error(err))
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		zap.L().Fatal("Failed to create image", zap.Error(err))
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		zap.L().Fatal("Failed to detect labels", zap.Error(err))
	}

	if len(annotations) == 0 {
		zap.S().Warn("No text found.")
		return ""
	}

	result := annotations[0].Description
	if contains(strings.Split(result, "\n"), "本日の運動結果") ||
		contains(strings.Split(result, "\n"), "Today's Results") {

		zap.L().Info(fmt.Sprintf("Detect results: %s\n", filename), zap.String("result", result))
		return result
	}

	return ""
}

func contains(texts []string, key string) bool {
	for _, text := range texts {
		if strings.HasPrefix(text, key) {
			return true
		}
	}
	return false
}
