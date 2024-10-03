package User

import (
	"github.com/go-webauthn/webauthn/main/helper/db"
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	Name        string
	DisplayName string
	ID          string
}

func (u *User) WebAuthnID() []byte {
	return []byte(u.ID)
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnCredentials provides the list of Credential objects owned by the user.
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	var creds []webauthn.Credential
	creds = append(creds, db.CdpDB[u.ID])
	return creds
}

func NewUser(name string, ID string) *User {
	return &User{
		Name:        name,
		DisplayName: name,
		ID:          ID,
	}
}
