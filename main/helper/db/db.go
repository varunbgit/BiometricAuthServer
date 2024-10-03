package db

import "github.com/go-webauthn/webauthn/webauthn"

var (
	Redis map[string]webauthn.SessionData
	CdpDB map[string]webauthn.Credential
)
