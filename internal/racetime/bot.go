package racetime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
)

const (
	msgChatHistory = "chat.history"
	msgChatMessage = "chat.message"
	msgChatDelete  = "chat.delete"
	msgChatPurge   = "chat.purge"
	msgError       = "error"
	msgPong        = "pong"
	msgRaceData    = "race.data"
)

type chatMsg struct {
	Date    time.Time `json:"date"`
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
	Type string `json:"type"`
}

type msg struct {
	Action string `json:"action"`
	Data   map[string]string
}

type Bot struct {
	token  TokenSet
	config config.Racetime
}

func NewBot(c config.Racetime) (*Bot, error) {
	resp, err := http.PostForm(fmt.Sprintf("%s/o/token", c.URL), url.Values{
		"client_id":     []string{c.ClientID},
		"client_secret": []string{c.ClientSecret},
		"grant_type":    []string{"client_credentials"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token TokenSet
	err = json.NewDecoder(resp.Body).Decode(&token)

	return &Bot{
		token:  token,
		config: c,
	}, nil
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
		return fmt.Errorf("ws dial: %s", err)
	}
	defer c.Close()
	log.Println("connected to ws")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					continue
				}

				var msg chatMsg
				err = json.Unmarshal(message, &msg)
				if err != nil {
					log.Println("read:", err)
					continue
				}

				err = processChatMessage(c, msg)
				if err != nil {
					log.Println("processing:", err)
					continue
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				continue
			}
			// shutdown gracefully, force if close request takes > 1 second
			select {
			case <-ctx.Done():
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func processChatMessage(c *websocket.Conn, msg chatMsg) error {
	if msg.Type != msgChatMessage {
		return nil
	}

	return sendMsg(c, msg.Message.MessagePlain)
	//return nil
}

func sendMsg(c *websocket.Conn, message string) error {
	m := msg{
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
