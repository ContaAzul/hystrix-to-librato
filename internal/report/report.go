package report

import (
	"github.com/caarlos0/hystrix-to-librato/internal/models"
	"github.com/caarlos0/hystrix-to-librato/internal/report/librato"
)

// Report metrics somehow
type Report interface {
	Report(data models.Data, cluster string)
	Close()
}

// Librato report type
func Librato(user, token string) Report {
	return librato.New(user, token)
}
