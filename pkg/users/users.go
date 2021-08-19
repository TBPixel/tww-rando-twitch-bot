package users

import (
	"fmt"
	"time"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/twitch"
)

// User
type User struct {
	TwitchID          string `badgerhold:"unique"`
	TwitchName        string
	TwitchDisplayName string
	FollowedByBot     bool
	JoinedAt          time.Time
}

// TwitchURL returns the absolute url to the users twitch channel
func (u User) TwitchURL() string {
	return fmt.Sprintf("%s/%s", twitch.URL, u.TwitchName)
}
