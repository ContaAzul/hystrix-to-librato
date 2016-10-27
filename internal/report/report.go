package report

import (
	"github.com/ContaAzul/hystrix-to-librato/internal/models"
	"github.com/ContaAzul/hystrix-to-librato/internal/report/librato"
)

// Report metrics somehow
type Report interface {
	Report(data models.Data, cluster string)
}

// Librato report type
func Librato(user, token string, metrics []string, interval int) Report {
	return librato.New(user, token, metrics, interval)
}
