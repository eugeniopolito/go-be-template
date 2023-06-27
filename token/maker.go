package token

import "time"

type Maker interface {
	CreateToken(username string, role int, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
