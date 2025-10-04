package tokengenerator

import "time"

type Generator interface {
	GenerateToken(subject string, now time.Time, duration time.Duration) string
	ParseToken(token string, now time.Time) (string, error)
}
