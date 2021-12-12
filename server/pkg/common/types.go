package common

import (
	"net/url"
	"time"
)

type (
	IsProduction      bool
	ClientID          string
	TraQBaseURL       *url.URL
	SessionSecret     string
	SessionKey        string
	SwiftAuthURL      *url.URL
	SwiftUserName     string
	SwiftPassword     string
	SwiftTenantID     string
	SwiftTenantName   string
	SwiftContainer    string
	FilePath          string
	AccessToken       string
	VerificationToken string
	DefaultChannels   []string
	UpdatedAt         time.Time
)
