package models

// Latencies of a circuit
type Latencies struct {
	L0   int64 `json:"0"`
	L25  int64 `json:"25"`
	L50  int64 `json:"50"`
	L75  int64 `json:"75"`
	L90  int64 `json:"90"`
	L95  int64 `json:"95"`
	L99  int64 `json:"99"`
	L995 int64 `json:"99.5"`
	L100 int64 `json:"100"`
}

// Data Hystrix main data type
type Data struct {
	Group          string      `json:"group"`
	Name           string      `json:"name"`
	Open           interface{} `json:"isCircuitBreakerOpen"`
	MeanLatency    int64       `json:"latencyExecute_mean"`
	LatencieTotals Latencies   `json:"latencyTotal"`
}
