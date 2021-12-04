package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/comail/colog"
	"github.com/mazrean/Quantainer/pkg/common"
)

func main() {
	env := os.Getenv("QUANTAINER_ENV")
	isProduction := env != "development"

	if !isProduction {
		colog.SetMinLevel(colog.LDebug)
		colog.SetFormatter(&colog.StdFormatter{
			Colors: true,
			Flag:   log.Ldate | log.Ltime | log.Lshortfile,
		})
	} else {
		colog.SetMinLevel(colog.LError)
		colog.SetFormatter(&colog.StdFormatter{
			Colors: false,
			Flag:   log.Ldate | log.Ltime | log.Lshortfile,
		})
	}

	colog.Register()

	secret, ok := os.LookupEnv("SESSION_SECRET")
	if !ok {
		panic("SESSION_SECRET is not set")
	}

	clientID := os.Getenv("CLIENT_ID")
	if len(clientID) == 0 {
		panic(errors.New("ENV CLIENT_ID IS NULL"))
	}
	clientSecret := os.Getenv("CLIENT_SECRET")
	if len(clientSecret) == 0 {
		panic(errors.New("ENV CLIENT_SECRET IS NULL"))
	}

	traQBaseURL, err := url.Parse("https://q.trap.jp/api/v3")
	if err != nil {
		panic(fmt.Sprintf("failed to parse traQBaseURL: %v", err))
	}

	filePath, ok := os.LookupEnv("FILE_PATH")
	if !ok {
		panic("ENV FILE_PATH is not set")
	}

	api, err := InjectAPI(&Config{
		IsProduction:  common.IsProduction(isProduction),
		SessionKey:    "sessions",
		SessionSecret: common.SessionSecret(secret),
		TraQBaseURL:   common.TraQBaseURL(traQBaseURL),
		OAuthClientID: common.ClientID(clientID),
		/* SwiftAuthURL:    swiftAuthURL,
		SwiftUserName:   swiftUserName,
		SwiftPassword:   swiftPassword,
		SwiftTenantID:   swiftTenantID,
		SwiftTenantName: swiftTenantName,
		SwiftContainer:  swiftContainer,*/
		FilePath:   common.FilePath(filePath),
		HttpClient: http.DefaultClient,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to inject API: %v", err))
	}

	addr, ok := os.LookupEnv("ADDR")
	if !ok {
		panic("ADDR is not set")
	}

	err = api.Start(addr)
	if err != nil {
		panic(fmt.Sprintf("failed to start API: %v", err))
	}
}
