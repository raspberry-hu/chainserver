package request

import (
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type UserClaims struct {
	UUID     uuid.UUID
	ID       uint
	TrueName string
	jwt.StandardClaims
}
