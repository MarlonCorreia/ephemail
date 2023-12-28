package clipb

import "golang.design/x/clipboard"

func SendToClipBoard(msg string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(msg))
	return nil
}
