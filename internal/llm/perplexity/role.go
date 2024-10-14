package perplexity

import (
	"fmt"
	"strings"
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=Role -trimprefix=Role
type Role int

const (
	RoleSystem Role = iota
	RoleUser
	RoleAssistant
)

func (r Role) MarshalText() ([]byte, error) {
	return []byte(strings.ToLower(r.String())), nil
}

func ParseRole(s string) (Role, error) {
	switch s {
	case strings.ToLower(RoleSystem.String()):
		return RoleSystem, nil
	case strings.ToLower(RoleUser.String()):
		return RoleUser, nil
	case strings.ToLower(RoleAssistant.String()):
		return RoleAssistant, nil
	}
	return RoleSystem, fmt.Errorf("perplexity: %q is not a valid role", s)
}

func (r *Role) UnmarshalText(text []byte) error {
	parsed, err := ParseRole(string(text))
	if err != nil {
		return err
	}

	*r = parsed
	return nil
}
