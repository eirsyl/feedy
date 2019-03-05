package pocket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eirsyl/flexit/log"
	"github.com/pkg/errors"

	"github.com/eirsyl/feedy/pkg/utils"
)

func (p *basePocket) Login(ctx context.Context, consumerKey string) (string, error) {

	logger := log.NewLogrusLogger(false)

	requestTokenBody, err := p.c.NewRequest("POST", "https://getpocket.com/v3/oauth/request", requestTokenRequest{
		ConsumerKey: consumerKey,
		RedirectURI: RedirectURL,
	})
	if err != nil {
		return "", errors.Wrap(err, "could not create request token body")
	}
	requestTokenBody.Header.Set("X-Accept", "application/json")

	var requestToken requestTokenResponse
	_, err = p.c.Do(ctx, requestTokenBody, &requestToken)
	if err != nil {
		return "", errors.Wrap(err, "could not retrieve request token")
	}

	// Open browser redirect
	loginURL := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", requestToken.Code, RedirectURL)
	utils.OpenURL(loginURL)
	logger.Infof("Open the following link in your browser: %s", loginURL)

	select {
	case cb := <-callbackServer():
		if cb.err != nil {
			return "", err
		}
	case <-interrupt():
		return "", ErrLoginCanceled
	}

	accessTokenBody, err := p.c.NewRequest("POST", "https://getpocket.com/v3/oauth/authorize", accessTokenRequest{
		ConsumerKey: consumerKey,
		Code:        requestToken.Code,
	})
	if err != nil {
		return "", errors.Wrap(err, "could not create access token request")
	}
	accessTokenBody.Header.Set("X-Accept", "application/json")

	var accessToken accessTokenResponse
	_, err = p.c.Do(ctx, accessTokenBody, &accessToken)
	if err != nil {
		return "", errors.Wrap(err, "could not retrieve access token")
	}

	return accessToken.AccessToken, nil
}

type requestTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectURI string `json:"redirect_uri"`
}

type requestTokenResponse struct {
	Code string `json:"code"`
}

type accessTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

type callback struct {
	err error
}

func callbackServer() chan *callback {
	c := make(chan *callback, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // nolint: errcheck
			return
		}

		w.Write([]byte("Go back to your terminal...")) // nolint: errcheck, gas

		c <- &callback{}
	})

	go func() {
		// Create listner
		serverListner, err := net.Listen("tcp", RedirectServerListen)
		if err != nil {
			c <- &callback{err: err}
			return
		}
		defer serverListner.Close() // nolint: errcheck

		// Create server
		server := http.Server{
			Handler:      mux,
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		}

		err = server.Serve(serverListner)
		c <- &callback{err: err}
	}()

	return c
}

func interrupt() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return c
}
