package instagramvideodownloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func DownloadProxy(fileURL, filename string) error {
	if filename == "" {
		filename = "gram-grabberz-video.mp4"
	}

	if fileURL == "" {
		return errors.New("missingUrl")
	}

	if !strings.HasPrefix(fileURL, "https://") {
		return errors.New("invalidUrl")
	}

	err := saveVideoToDisk(fileURL, filename)
	if err != nil {
		return fmt.Errorf("failed to save video: %v", err)
	}
	return nil
}

func saveVideoToDisk(fileURL, filename string) error {
	resp, err := http.Get(fileURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch video: %v", err)
	}
	defer resp.Body.Close()

	if err := os.MkdirAll("./downloads", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create downloads directory: %v", err)
	}

	out, err := os.Create("./downloads/" + filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save video: %v", err)
	}
	return nil
}
