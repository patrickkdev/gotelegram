package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type chat struct {
	name   string
	chatID string
	apiKey string
}

func Chat(name string, id string, apiKey string) chat {
	return chat{name: name, chatID: id, apiKey: apiKey}
}

func (t chat) GetChatId() string {
	return t.chatID
}

func (t chat) GetApiKey() string {
	return t.apiKey
}

func (t chat) GetChatName() string {
	return t.name
}

func (t chat) String() string {
	return fmt.Sprintf("Telegram: %s (%s)", t.name, t.chatID)
}

func (t chat) SendAsync(message any) {
	go t.Send(message)
}

func (t chat) Send(message any) error {
	return t.send(message)	
}

func (t chat) send(message any) error {
	// if message is a string, just send it
	msgStr, ok := message.(string)
	if !ok {
		if jsonData, err := json.MarshalIndent(message, "", "  "); err == nil {
			// if message is not a string, try to convert it to a JSON indented string

			msgStr = string(jsonData)
		} else {
			// if it fails, just convert it to a string

			msgStr = fmt.Sprintf("%v", message)
		}
	}

	if len(msgStr) <= 4096 {
		// If the message is within the limit, send it as is
		return t.sendChunk(msgStr)
	} else {
		// If the message is longer than 4096 characters, split it and send in chunks
		chunks := splitMessageIntoChunks(msgStr, 4096)

		var err error
		for _, chunk := range chunks {
			err = t.sendChunk(chunk)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (t chat) sendChunk(message string) error {
	body := map[string]string{
		"chat_id":                  t.chatID,
		"text":                     message,
		"disable_web_page_preview": "true",
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.apiKey)

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil || resp.StatusCode > 299 {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func splitMessageIntoChunks(message string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(message); i += chunkSize {
		end := i + chunkSize
		if end > len(message) {
			end = len(message)
		}
		chunks = append(chunks, message[i:end])
	}
	return chunks
}
