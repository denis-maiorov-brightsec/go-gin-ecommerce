package auth

import (
	"strings"
)

const (
	PermissionManagePromotions = "promotions:manage"

	StubPromotionsAdminToken = "stub-promotions-admin"
	StubReadonlyUserToken    = "stub-readonly-user"
)

type Identity struct {
	SubjectID   string
	Permissions map[string]struct{}
}

func (i Identity) HasPermission(permission string) bool {
	_, ok := i.Permissions[permission]
	return ok
}

type Authenticator interface {
	AuthenticateBearerToken(token string) (*Identity, bool)
}

type StubAuthenticator struct{}

func NewStubAuthenticator() StubAuthenticator {
	return StubAuthenticator{}
}

func (StubAuthenticator) AuthenticateBearerToken(token string) (*Identity, bool) {
	switch strings.TrimSpace(token) {
	case StubPromotionsAdminToken:
		return &Identity{
			SubjectID: "promotions-admin",
			Permissions: map[string]struct{}{
				PermissionManagePromotions: {},
			},
		}, true
	case StubReadonlyUserToken:
		return &Identity{
			SubjectID:   "readonly-user",
			Permissions: map[string]struct{}{},
		}, true
	default:
		return nil, false
	}
}
