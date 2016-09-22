package librato

import (
	"log"
	"sync"
	"time"

	"github.com/caarlos0/hystrix-to-librato/internal/models"
	librato "github.com/rcrowley/go-librato"
)

// New report type
func New(user, token string) *Librato {
	return &Librato{
		user:    user,
		token:   token,
		reports: make(map[string]time.Time),
		lock:    sync.RWMutex{},
	}
}

// Librato type
type Librato struct {
	user    string
	token   string
	reports map[string]time.Time
	lock    sync.RWMutex
}

// Report the given data to librato for the given cluster
func (r *Librato) Report(data models.Data, cluster string) {
	r.circuitOpen(data, cluster+"."+data.Group)
	r.latencies(data, cluster+"."+data.Group+"."+data.Name)
}

func (r *Librato) shouldReport(source string) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	val, ok := r.reports[source]
	if ok && time.Since(val).Seconds() < 5 {
		return false
	}
	r.reports[source] = time.Now()
	return true
}
func (r *Librato) latencies(data models.Data, source string) {
	if !r.shouldReport(source) {
		return
	}
	log.Println("Report", source)
	m := librato.NewSimpleMetrics(r.user, r.token, source)
	defer m.Wait()
	defer m.Close()
	m.NewCounter("hystrix.latency.100th") <- data.LatencieTotals.L100
	m.NewCounter("hystrix.latency.99.5th") <- data.LatencieTotals.L99_5
	m.NewCounter("hystrix.latency.99th") <- data.LatencieTotals.L99
	m.NewCounter("hystrix.latency.95th") <- data.LatencieTotals.L95
	m.NewCounter("hystrix.latency.90th") <- data.LatencieTotals.L90
	m.NewCounter("hystrix.latency.75th") <- data.LatencieTotals.L75
	m.NewCounter("hystrix.latency.50th") <- data.LatencieTotals.L50
	m.NewCounter("hystrix.latency.25th") <- data.LatencieTotals.L25
	m.NewCounter("hystrix.latency.0th") <- data.LatencieTotals.L0
	m.NewCounter("hystrix.latency.mean") <- data.MeanLatency
}

func (r *Librato) circuitOpen(data models.Data, source string) {
	if !r.shouldReport(source) {
		return
	}
	log.Println("Report", source)
	m := librato.NewSimpleMetrics(r.user, r.token, source)
	defer m.Wait()
	defer m.Close()
	c := m.NewCounter("hystrix.circuit.open")
	if isOpen(data.Open) {
		c <- 1
	} else {
		c <- 0
	}
}
func isOpen(data interface{}) bool {
	if b, ok := data.(bool); ok {
		return b
	}
	return true
}
