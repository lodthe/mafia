package player

import "github.com/lodthe/mafia/pkg/mafiapb"

type Role string

const (
	RoleInnocent Role = "INNOCENT"
	RoleSheriff  Role = "SHERIFF"
	RoleMafiosi  Role = "MAFIOSI"
)

func (r Role) Proto() mafiapb.Role {
	switch r {
	case RoleInnocent:
		return mafiapb.Role_ROLE_INNOCENT

	case RoleSheriff:
		return mafiapb.Role_ROLE_SHERIFF

	case RoleMafiosi:
		return mafiapb.Role_ROLE_MAFIOSI

	default:
		return mafiapb.Role_ROLE_UNKNOWN
	}
}
