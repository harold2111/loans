package postgres

import "github.com/lib/pq"

const (
	UniqueConstraintIdentification = "uix_clients_identification"
)

func IsUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" && pqErr.Constraint == constraintName
	}
	return false
}
