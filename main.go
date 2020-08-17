package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Constants
const debug = true
const tgBaseURL = "https://api.telegram.org/bot"
const tgToken = "YOUR_TOKEN"

// Methods
const (
	getMe       = "getMe"
	getUpdates  = "getUpdates"
	sendMessage = "sendMessage"
)

// IntSlice ...
type IntSlice []int

// GetMeType ...
type GetMeType struct {
	Ok     bool            `json:"ok"`
	Result GetMeResultType `json:"result"`
}

// GetMeResultType ...
type GetMeResultType struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// SendMessageType ...
type SendMessageType struct {
	Ok     bool        `json:"ok"`
	Result MessageType `json:"result"`
}

// GetUpdatesType ...
type GetUpdatesType struct {
	Ok     bool                   `json:"ok"`
	Result []GetUpdatesResultType `json:"result"`
}

// GetUpdatesResultType ...
type GetUpdatesResultType struct {
	UpdateID int         `json:"update_id"`
	Message  MessageType `json:"message,omitempty"`
}

// MessageType ...
type MessageType struct {
	MessageID int                             `json:"message_id"`
	From      GetUpdatesResultMessageFromType `json:"from"`
	Chat      GetUpdatesResultMessageChatType `json:"chat"`
	Date      int64                           `json:"date"`
	Text      string                          `json:"text"`
}

// GetUpdatesResultMessageFromType ...
type GetUpdatesResultMessageFromType struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// GetUpdatesResultMessageChatType ...
type GetUpdatesResultMessageChatType struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

func main() {
	// start time
	now := time.Now()
	startTime := now.Unix()

	// urls
	getUpdatesURL := getURL(getUpdates)

	// handled messages

	var doneMessagesID = IntSlice{}
	var count int

	for {
		// get updates
		count = 0
		body := sendGet(getUpdatesURL)
		getUpdates := GetUpdatesType{}
		err := json.Unmarshal(body, &getUpdates)
		debugError(err)

		// return answers
		for _, item := range getUpdates.Result {
			if item.Message.Date > startTime && !doneMessagesID.Has(item.Message.MessageID) {
				handleMessage(item.Message)
				count++
				doneMessagesID = append(doneMessagesID, item.Message.MessageID)
			}
		}
		if count > 0 {
			fmt.Println("Done messages: " + strconv.Itoa(count))
		}
		time.Sleep(5)
	}
}

// Has ...
func (list IntSlice) Has(a int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// handleMessage ...
func handleMessage(msg MessageType) {
	sendMessageURL := getURL(sendMessage)
	chatID := strconv.Itoa(msg.Chat.ID)
	if msg.Text == "test" {
		values := map[string]string{"chat_id": chatID, "text": "Answer: " + msg.Text}
		jsonValue, _ := json.Marshal(values)
		sendPost(sendMessageURL, jsonValue)
	}
}

// DebugError print error
func debugError(err error) {
	if err != nil && debug {
		fmt.Println(err.Error())
	}
}

// GetURL conctruct full url by method
func getURL(methodName string) string {
	return tgBaseURL + tgToken + "/" + methodName
}

// SendGet send Get reauest and return Response body
func sendGet(url string) []byte {
	response, err := http.Get(url)
	debugError(err)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	debugError(err)

	return body
}

// SendPost send Post request and return Response body
func sendPost(url string, data []byte) []byte {
	r := bytes.NewReader(data)
	response, err := http.Post(url, "application/json", r)
	debugError(err)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	debugError(err)

	return body
}
