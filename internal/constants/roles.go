package constants

type JWTRole uint8

var jwtRoles = [...]string{
	"free",
	"pro",
	"admin",
}

const (
	Free JWTRole = iota
	Pro
	Admin
)

func (rs JWTRole) String() string {
	return jwtRoles[rs]
}
