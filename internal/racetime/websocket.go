package racetime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	botPrefix      = "!twwr"
	msgChatHistory = "chat.history"
	msgChatMessage = "chat.message"
	msgChatDelete  = "chat.delete"
	msgChatPurge   = "chat.purge"
	msgError       = "error"
	msgPong        = "pong"
	msgRaceData    = "race.data"
)

type recv struct {
	Type string    `json:"type"`
	Date time.Time `json:"date"`
}

type chatMsg struct {
	recv
	Message struct {
		Bot          interface{} `json:"bot"`
		Delay        int         `json:"delay"`
		Highlight    bool        `json:"highlight"`
		ID           string      `json:"id"`
		IsBot        bool        `json:"is_bot"`
		IsMonitor    bool        `json:"is_monitor"`
		IsSystem     bool        `json:"is_system"`
		Message      string      `json:"message"`
		MessagePlain string      `json:"message_plain"`
		PostedAt     time.Time   `json:"posted_at"`
		User         UserData    `json:"user"`
	} `json:"message"`
}

type send struct {
	Action string            `json:"action"`
	Data   map[string]string `json:"data"`
}

// Connect to a raceroom's chat via name
func (b Bot) Connect(ctx context.Context, name string) error {
	u, err := url.Parse(b.config.URL)
	if err != nil {
		return err
	}
	u.Scheme = b.config.WSSchema
	u.Path = fmt.Sprintf("/ws/o/bot/%s", name)
	q := u.Query()
	q.Set("token", b.token.AccessToken)
	u.RawQuery = q.Encode()

	log.Printf("connecting to ws: %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			err = processChatMessage(c, message)
			if err != nil {
				log.Println("process: ", err)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("write close: %s", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func processChatMessage(c *websocket.Conn, msg []byte) error {
	var message recv
	err := json.Unmarshal(msg, &message)
	if err != nil {
		return err
	}

	if message.Type != msgChatMessage {
		return nil
	}

	var cm chatMsg
	err = json.Unmarshal(msg, &cm)
	if err != nil {
		return err
	}

	if cm.Message.IsBot {
		return nil
	}

	if !strings.HasPrefix(cm.Message.MessagePlain, botPrefix) || len(cm.Message.MessagePlain) <= len(botPrefix) {
		return nil
	}

	command := strings.TrimSpace(cm.Message.MessagePlain[len(botPrefix)+1:])
	log.Println(len(command))

	return sendMsg(c, cm.Message.MessagePlain)
}

func sendMsg(c *websocket.Conn, message string) error {
	m := send{
		Action: "message",
		Data: map[string]string{
			"message": message,
			"guid":    uuid.New().String(),
		},
	}

	d, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return c.WriteMessage(websocket.TextMessage, d)
}
