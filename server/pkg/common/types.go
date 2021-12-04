package common

import (
	"net/url"
)

type (
	IsProduction    bool
	ClientID        string
	TraQBaseURL     *url.URL
	SessionSecret   string
	SessionKey      string
	SwiftAuthURL    *url.URL
	SwiftUserName   string
	SwiftPassword   string
	SwiftTenantID   string
	SwiftTenantName string
	SwiftContainer  string
	FilePath        string
)
