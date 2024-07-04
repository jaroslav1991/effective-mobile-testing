package validators

import (
	"errors"
	"github.com/lib/pq"
)

func CheckLenPassport(passport []string) bool {
	if len(passport) != 2 {
		return false
	}

	if len(passport[0]) != 4 || len(passport[1]) != 6 {
		return false
	}

	return true
}

func IsConstrainError(err error) (error, bool) {
	if err != nil {
		var errPq *pq.Error
		if errors.As(err, &errPq) {
			if errPq.Code == "23505" || errPq.Code == "23503" {
				return err, true
			}
		}
	}

	return err, false
}
