package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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

	var (
		swiftAuthURL    common.SwiftAuthURL
		swiftUserName   common.SwiftUserName
		swiftPassword   common.SwiftPassword
		swiftTenantID   common.SwiftTenantID
		swiftTenantName common.SwiftTenantName
		swiftContainer  common.SwiftContainer
	)
	if isProduction {
		strSwiftAuthURL, ok := os.LookupEnv("OS_AUTH_URL")
		if !ok {
			panic("ENV OS_AUTH_URL is not set")
		}
		swiftAuthURL, err = url.Parse(strSwiftAuthURL)
		if err != nil {
			panic(fmt.Errorf("failed to parse swiftAuthURL: %w", err))
		}

		strSwiftUserName, ok := os.LookupEnv("OS_USERNAME")
		if !ok {
			panic("ENV OS_USERNAME is not set")
		}
		swiftUserName = common.SwiftUserName(strSwiftUserName)

		strSwiftPassword, ok := os.LookupEnv("OS_PASSWORD")
		if !ok {
			panic("ENV OS_PASSWORD is not set")
		}
		swiftPassword = common.SwiftPassword(strSwiftPassword)

		strSwiftTenantID, ok := os.LookupEnv("OS_TENANT_ID")
		if !ok {
			panic("ENV OS_TENANT_ID is not set")
		}
		swiftTenantID = common.SwiftTenantID(strSwiftTenantID)

		strSwiftTenantName, ok := os.LookupEnv("OS_TENANT_NAME")
		if !ok {
			panic("ENV OS_TENANT_NAME is not set")
		}
		swiftTenantName = common.SwiftTenantName(strSwiftTenantName)

		strSwiftContainer, ok := os.LookupEnv("OS_CONTAINER")
		if !ok {
			panic("ENV OS_CONTAINER is not set")
		}
		swiftContainer = common.SwiftContainer(strSwiftContainer)
	}

	accessToken, ok := os.LookupEnv("ACCESS_TOKEN")
	if !ok {
		panic("ENV ACCESS_TOKEN is not set")
	}

	verificationToken, ok := os.LookupEnv("VERIFICATION_TOKEN")
	if !ok {
		panic("ENV VERIFICATION_TOKEN is not set")
	}

	defaultChannels, ok := os.LookupEnv("DEFAULT_CHANNELS")
	if !ok {
		panic("ENV DEFAULT_CHANNELS is not set")
	}

	service, err := InjectService(&Config{
		IsProduction:      common.IsProduction(isProduction),
		SessionKey:        "sessions",
		SessionSecret:     common.SessionSecret(secret),
		TraQBaseURL:       common.TraQBaseURL(traQBaseURL),
		OAuthClientID:     common.ClientID(clientID),
		SwiftAuthURL:      swiftAuthURL,
		SwiftUserName:     swiftUserName,
		SwiftPassword:     swiftPassword,
		SwiftTenantID:     swiftTenantID,
		SwiftTenantName:   swiftTenantName,
		SwiftContainer:    swiftContainer,
		FilePath:          common.FilePath(filePath),
		HttpClient:        http.DefaultClient,
		AccessToken:       common.AccessToken(accessToken),
		VerificationToken: common.VerificationToken(verificationToken),
		DefaultChannels:   common.DefaultChannels(strings.Split(defaultChannels, ",")),
		UpdatedAt:         common.UpdatedAt(time.Now().Add(-time.Hour * 24 * 365 * 3)),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to inject API: %v", err))
	}

	api := service.API

	addr, ok := os.LookupEnv("ADDR")
	if !ok {
		panic("ADDR is not set")
	}

	err = api.Start(addr)
	if err != nil {
		panic(fmt.Sprintf("failed to start API: %v", err))
	}
}
