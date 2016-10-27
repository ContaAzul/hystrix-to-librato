package report

import (
	"time"

	"github.com/ContaAzul/hystrix-to-librato/internal/models"
	"github.com/ContaAzul/hystrix-to-librato/internal/report/librato"
)

// Report metrics somehow
type Report interface {
	Report(data models.Data, cluster string)
}

// Librato report type
func Librato(user, token string, metrics []string, interval time.Duration) Report {
	return librato.New(user, token, metrics, interval)
}
