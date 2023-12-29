package utils

import "os"

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
