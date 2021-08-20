package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

// authorizeUser implements the OAuth2 flow.
func authorizeUser(ctx *cli.Context, authUrl, tokenUrl, clientID, clientSecret, redirectURL string) {
	// construct the authorization URL
	authorizationURL, _ := url.Parse(authUrl)
	q := authorizationURL.Query()
	q.Set("scope", "openid")
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
			fmt.Println("twitch: Url Param 'code' is missing")
			io.WriteString(w, "Error: could not find 'code' URL parameter\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// trade the authorization code and the code verifier for an access token
		//codeVerifier := CodeVerifier.String()
		contents, err := getAccessTokenContents(tokenUrl, clientID, clientSecret, code, redirectURL)
		if err != nil {
			fmt.Printf("twitch: could not get access token: %s", err)
			io.WriteString(w, "Error: could not retrieve access token\n")

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

type twitchAccessTokenContents struct {
	UserID            string  `json:"user_id"`
	ExpiresAt         float64 `json:"expires_in"`
	PreferredUsername string  `json:"preferred_username"`
}

// getAccessTokenContents trades the authorization code retrieved from the first OAuth2 leg for an access token
func getAccessTokenContents(tokenUrl, clientID, clientSecret, authorizationCode, callbackURL string) (*twitchAccessTokenContents, error) {
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
		fmt.Printf("twitch: HTTP error: %s", err)
		return nil, err
	}

	// process the response
	defer res.Body.Close()
	var responseBody map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)

	// unmarshal the json into a string map
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		fmt.Printf("twitch: JSON error: %s", err)
		return nil, err
	}

	tkn, _, err := new(jwt.Parser).ParseUnverified(responseBody["id_token"].(string), jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	var contents *twitchAccessTokenContents
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
		contents = &twitchAccessTokenContents{
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

// cleanup closes the HTTP server
func cleanup(server *http.Server) {
	// we run this as a goroutine so that this function falls through and
	// the socket to the browser gets flushed/closed before the server goes away
	go server.Close()
}
