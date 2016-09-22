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

func main() {
	config := config.Get()
	report := report.Librato(config.User, config.Token)
	defer report.Close()
	for _, cluster := range config.Clusters {
		go read(config.URL, cluster, report)
	}
	// sleep forever
	for {
		time.Sleep(10 * time.Second)
		log.Println(runtime.NumGoroutine(), "goroutines running")
	}
}

func restart(url, cluster string, report report.Report, err error) {
	log.Println(err)
	time.Sleep(10 * time.Second)
	read(url, cluster, report)
}

func read(url, cluster string, report report.Report) {
	log.Println("Starting", cluster)
	resp, err := http.Get(url + "?cluster=" + cluster)
	if err != nil {
		restart(url, cluster, report, err)
		return
	}
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			restart(url, cluster, report, err)
			return
		}
		sline := string(line)
		if !isData(sline) || !isCircuitReport(sline) {
			continue
		}
		doReport(report, cluster, sline)
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
