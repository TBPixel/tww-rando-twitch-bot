package racetime

const (
	TokenURL = "o/token"
	AuthURL  = "o/authorize"
)

type TokenSet struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
