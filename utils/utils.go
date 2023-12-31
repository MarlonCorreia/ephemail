package utils

import (
	"os"

	"golang.design/x/clipboard"
)

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func WriteFile(fileName string, content string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(fileName string) error {
	err := os.Remove(fileName)
	if err != nil {
		return err
	}
	return nil
}

func SendToClipBoard(msg string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(msg))
	return nil
}
