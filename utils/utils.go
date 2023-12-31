package utils

import (
	"fmt"
	"os"
	"strings"

	"golang.design/x/clipboard"
)

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func safeFileName(filePath string) string {
	var safeName strings.Builder
	for _, b := range filePath {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			safeName.WriteByte(byte(b))
		} else if b == ' ' {
			safeName.WriteByte('_')
		}
	}
	return safeName.String()
}

func fileNameExtFromPath(filePath string) (string, string) {
	extIdx := strings.LastIndex(filePath, ".")
	return filePath[0:extIdx], filePath[extIdx:]
}

func WriteFile(filePath string, content []byte) error {
	name, ext := fileNameExtFromPath(filePath)
	safeFilePath := fmt.Sprintf("%s%s", safeFileName(name), ext)

	f, err := os.Create(safeFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
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
