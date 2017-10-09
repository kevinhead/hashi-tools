package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

const (
	ver           = "0.0.2"
	checkTemplate = `
{
  "check": {
	"id": "{{.ID}}",
	"name": "{{.Name}}",
	"notes": "{{.Notes}}",
	"script": "{{.Script}}",
	"interval": "{{.Interval}}",
	"timeout": "{{.Timeout}}"
  }
}
	`
)

// healthCheck struct
type healthCheck struct {
	ID       string
	Name     string
	Notes    string
	Script   string
	Interval string
	Timeout  string
}

func main() {
	var version, check bool
	var max, warn float64

	flag.Float64Var(&max, "critical", 90, "Used percent critical threshold.")
	flag.Float64Var(&warn, "warn", 70, "Used percent warning threshold.")

	flag.BoolVar(&version, "version", false, "Prints version")
	flag.BoolVar(&check, "json", false, "Prints Consul Check definition and exits")

	flag.Parse()

	if version {
		fmt.Printf("%s v%s\n", os.Args[0], ver)
		return
	}

	if check {
		t := template.Must(template.New("check").Parse(checkTemplate))

		hc := healthCheck{
			ID:       "cpu-check",
			Name:     "CPU Check",
			Notes:    "Checks percent used. (all processors)",
			Script:   fmt.Sprintf("/opt/consul/bin/cpu-check -critical %.0f -warn %.0f", max, warn),
			Interval: "10s",
			Timeout:  "1s",
		}

		t.Execute(os.Stdout, hc)
		return
	}

	if flag.NFlag() != 2 {
		flag.Usage()
		os.Exit(-1)
	}

	u, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		log.Printf("cpu.Percent failed with - %v", err)
		os.Exit(-1)
	}

	if u[0] > max {
		fmt.Printf("CPU CRITICAL - %.2f%% used", u[0])
		os.Exit(2)
	}

	if u[0] > warn {
		fmt.Printf("CPU WARNING - %.2f%% used", u[0])
		os.Exit(1)
	}

	fmt.Printf("CPU OK - %.2f%% used", u[0])

}
