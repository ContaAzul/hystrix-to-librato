package librato

import (
	"log"
	"sync"

	"github.com/caarlos0/hystrix-to-librato/internal/models"
	librato "github.com/rcrowley/go-librato"
)

var lock = sync.RWMutex{}

// New report type
func New(user, token string) *Librato {
	return &Librato{
		user:    user,
		token:   token,
		clients: make(map[string]librato.Metrics),
	}
}

// Librato type
type Librato struct {
	user    string
	token   string
	clients map[string]librato.Metrics
}

// Report the given data to librato for the given cluster
func (r *Librato) Report(data models.Data, cluster string) {
	log.Println("Report", cluster)
	go circuitOpen(data, r.client(cluster+"."+data.Group))
	go latencies(data, r.client(cluster+"."+data.Name))
}

// Close all clients
func (r *Librato) Close() {
	log.Println("Closing all librato clients...")
	for _, client := range r.clients {
		defer client.Wait()
		defer client.Close()
	}
}

func (r *Librato) client(source string) librato.Metrics {
	client, ok := r.clients[source]
	if !ok {
		client = librato.NewSimpleMetrics(r.user, r.token, source)
		log.Println("Creating new librato client for", source)
		lock.Lock()
		defer lock.Unlock()
		r.clients[source] = client
	}
	return client
}

func latencies(data models.Data, m librato.Metrics) {
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

func circuitOpen(data models.Data, m librato.Metrics) {
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
