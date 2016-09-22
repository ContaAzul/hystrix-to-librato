package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/caarlos0/hystrix-to-librato/internal/config"
	"github.com/caarlos0/hystrix-to-librato/internal/models"
	"github.com/caarlos0/hystrix-to-librato/internal/report"
)

var waitTime = 5 * time.Second

func main() {
	config := config.Get()
	report := report.Librato(config.User, config.Token)
	for _, cluster := range config.Clusters {
		go read(config.URL, cluster, report)
	}
	// sleep forever
	for {
		time.Sleep(waitTime)
		log.Println(runtime.NumGoroutine(), "goroutines running...")
	}
}

func read(url, cluster string, report report.Report) {
	time.Sleep(waitTime)
	log.Println("Starting", cluster)
	resp, err := http.Get(url + "?cluster=" + cluster)
	if err != nil {
		read(url, cluster, report)
		return
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !isData(line) || !isCircuitReport(line) {
			continue
		}
		doReport(report, cluster, line)
	}
}

func doReport(report report.Report, cluster, line string) {
	line = strings.TrimPrefix(line, "data:")
	var data models.Data
	if err := json.Unmarshal([]byte(line), &data); err != nil {
		return
	}
	go report.Report(data, cluster)
}

func isData(line string) bool {
	return strings.HasPrefix(line, "data:")
}

func isCircuitReport(line string) bool {
	return strings.Contains(line, "isCircuitBreakerOpen")
}
