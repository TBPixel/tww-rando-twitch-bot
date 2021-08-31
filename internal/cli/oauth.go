package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

type tokenParserFunc func(config config.App, reader io.ReadCloser) (interface{}, error)

// authorizeUser implements the OAuth2 flow.
func authorizeUser(ctx *cli.Context, config config.App, authUrl, tokenUrl, clientID, clientSecret, redirectURL string, scopes []string, parserFunc tokenParserFunc) {

	// construct the authorization URL
	authorizationURL, _ := url.Parse(authUrl)
	q := authorizationURL.Query()
	q.Set("scope", strings.Join(scopes, " "))
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURL)
	authorizationURL.RawQuery = q.Encode()

	// start a web server to listen on a callback URL
	server := &http.Server{Addr: redirectURL}

	// define a handler that will get the authorization code, call the token endpoint, and close the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get the authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Println("oauth: Url Param 'code' is missing")
			io.WriteString(w, "Error: could not find 'code' URL parameter\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// trade the authorization code and the code verifier for an access token
		//codeVerifier := CodeVerifier.String()
		reader, err := getTokenResponse(tokenUrl, clientID, clientSecret, code, redirectURL)
		if err != nil {
			fmt.Printf("oauth: could not get access token: %s", err)
			io.WriteString(w, "Error: could not retrieve access token\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		contents, err := parserFunc(config, reader)
		if err != nil {
			fmt.Printf("oauth: could not parse access token: %s", err)
			io.WriteString(w, "error, failed to parse access token\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}
		ctx.Context = context.WithValue(ctx.Context, "token", contents)

		// return an indication of success to the caller
		io.WriteString(w, `
		<html>
			<body>
				<h1>Login successful!</h1>
				<h2>You can close this window and return to the twwr CLI.</h2>
			</body>
		</html>`)

		fmt.Println("Successfully logged into twitch API.")

		// close the HTTP server
		cleanup(server)
	})

	// parse the redirect URL for the port number
	u, err := url.Parse(redirectURL)
	if err != nil {
		fmt.Printf("twitch: bad redirect URL: %s\n", err)
		os.Exit(1)
	}

	// set up a listener on the redirect port
	port := fmt.Sprintf(":%s", u.Port())
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("twitch: can't listen to port %s: %s\n", port, err)
		os.Exit(1)
	}

	// open a browser window to the authorizationURL
	err = open.Run(authorizationURL.String())
	if err != nil {
		fmt.Printf("twitch: can't open browser to URL %s: %s\n", authorizationURL, err)
		os.Exit(1)
	}

	// start the blocking web server loop
	// this will exit when the handler gets fired and calls server.Close()
	server.Serve(l)
}

type TwitchAccessTokenContents struct {
	UserID            string  `json:"user_id"`
	ExpiresAt         float64 `json:"expires_in"`
	PreferredUsername string  `json:"preferred_username"`
}

type RacetimeAccessTokenContents struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Pronouns   string `json:"pronouns"`
	Flair      string `json:"flair"`
	TwitchName string `json:"twitch_name"`
}

// getTokenResponse trades the authorization code retrieved from the first OAuth2 leg for an access token
func getTokenResponse(tokenUrl, clientID, clientSecret, authorizationCode, callbackURL string) (io.ReadCloser, error) {
	// set the url and form-encoded data for the POST to the access token endpoint
	data := fmt.Sprintf(
		"grant_type=authorization_code"+
			"&client_id=%s"+
			"&client_secret=%s"+
			"&code=%s"+
			"&redirect_uri=%s",
		clientID, clientSecret, authorizationCode, callbackURL)
	payload := strings.NewReader(data)

	// create the request and execute it
	req, _ := http.NewRequest("POST", tokenUrl, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("oauth: HTTP error: %s", err)
		return nil, err
	}

	return res.Body, nil
}

// cleanup closes the HTTP server
func cleanup(server *http.Server) {
	// we run this as a goroutine so that this function falls through and
	// the socket to the browser gets flushed/closed before the server goes away
	go server.Close()
}

func twitchTokenParserFunc(_ config.App, reader io.ReadCloser) (interface{}, error) {
	defer reader.Close()
	var responseBody map[string]interface{}

	// unmarshal the json into a string map
	err := json.NewDecoder(reader).Decode(&responseBody)
	if err != nil {
		fmt.Printf("twitch: JSON error: %s", err)
		return nil, err
	}

	tkn, _, err := new(jwt.Parser).ParseUnverified(responseBody["id_token"].(string), jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	var contents TwitchAccessTokenContents
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
		contents = TwitchAccessTokenContents{
			PreferredUsername: claims["preferred_username"].(string),
			UserID:            claims["sub"].(string),
			ExpiresAt:         claims["exp"].(float64),
		}
	} else {
		log.Printf("failed to parse claims from token: %s", responseBody["access_token"])
	}

	// retrieve the access token out of the map, and return to caller
	return contents, nil
}

func racetimeTokenParserFunc(config config.App, reader io.ReadCloser) (interface{}, error) {
	defer reader.Close()
	var responseBody map[string]interface{}

	// unmarshal the json into a string map
	err := json.NewDecoder(reader).Decode(&responseBody)
	if err != nil {
		fmt.Printf("racetime: JSON error: %s", err)
		return nil, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/o/userinfo", config.Racetime.URL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", responseBody["access_token"]))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var contents RacetimeAccessTokenContents
	err = json.NewDecoder(res.Body).Decode(&contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
