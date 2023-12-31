package mail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MarlonCorreia/ephemail/utils"
	"github.com/k3a/html2text"
)

var (
	baseEmailAPI = "https://www.1secmail.com/api/v1/"
)

type Attachment struct {
	FileName    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        uint   `json:"size"`
}

type MessageContent struct {
	Body        string       `json:"body"`
	TextBody    string       `json:"textBody"`
	HtmlBody    string       `json:"htmlBody"`
	Attachments []Attachment `json:"attachments"`
}

type Message struct {
	Id      uint   `json:"id"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Date    string `json:"date"`
	Content MessageContent
}

type EmailModel struct {
	User     string
	Domain   string
	Messages []*Message
}

func (m *EmailModel) GetEmail() string {
	return m.User + "@" + m.Domain
}

func (m *EmailModel) ContainsMessage(id uint) bool {
	for _, ele := range m.Messages {
		if ele.Id == id {
			return true
		}
	}
	return false
}

func (m *EmailModel) AddMessage(msg *Message) {
	if !m.ContainsMessage(msg.Id) {
		m.GetMessageContent(msg)
		m.Messages = append(m.Messages, msg)
	}
}

func (m *EmailModel) UpdateEmailMessages() error {
	url := fmt.Sprintf("%s?action=getMessages&login=%s&domain=%s", baseEmailAPI, m.User, m.Domain)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	msgs := []*Message{}
	json.Unmarshal(byteBody, &msgs)

	for _, msg := range msgs {
		m.AddMessage(msg)
	}

	return nil
}

func (m *EmailModel) GetMessageContent(msg *Message) error {
	url := fmt.Sprintf("%s?action=readMessage&login=%s&domain=%s&id=%d", baseEmailAPI, m.User, m.Domain, msg.Id)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	bytebody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var body MessageContent

	err = json.Unmarshal(bytebody, &body)
	if err != nil {
		return err
	}
	strBody := html2text.HTML2TextWithOptions(body.Body, html2text.WithUnixLineBreaks())
	msg.Content = body
	msg.Content.TextBody = strBody

	return nil
}

func (m *EmailModel) DownloadAttachment(msg *Message, att Attachment) error {
	strId := fmt.Sprint(msg.Id)
	url := fmt.Sprintf(
		"%s?action=download&login=%s&domain=%s&id=%s&file=%s",
		baseEmailAPI,
		m.User,
		m.Domain,
		strId,
		att.FileName,
	)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("Unable to download file")
	}

	content, err := io.ReadAll(res.Body)
	utils.WriteFile(att.FileName, content)

	return nil
}

func (msg *Message) DisplayMessage() string {
	return fmt.Sprintf("%s - From: %s - At: %s\n", msg.Subject, msg.From, msg.Date)
}

func (msg *Message) DisplayCompleteEmail() string {
	emailBody := "\nSubject: %s\nFrom: %s\nAt: %s\n\n%s"
	return fmt.Sprintf(emailBody, msg.Subject, msg.From, msg.Date, msg.Content.TextBody)
}

func (m *EmailModel) BuildNewEmail() error {
	url := fmt.Sprintf("%s?action=genRandomMailbox&count=1", baseEmailAPI)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	bytesBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var body []string
	err = json.Unmarshal(bytesBody, &body)
	if err != nil {
		return err
	}

	user_domain := strings.Split(body[0], "@")
	m.User = user_domain[0]
	m.Domain = user_domain[1]

	return nil
}
