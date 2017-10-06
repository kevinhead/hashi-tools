package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/shirou/gopsutil/mem"
)

const (
	ver           = "0.0.1"
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
			ID:       "mem-check",
			Name:     "Memory Check",
			Notes:    "Checks percent used.",
			Script:   fmt.Sprintf("/opt/consul/bin/mem-check -critical %.0f -warn %.0f", max, warn),
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

	u, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("mem.VirtualMemory failed with - %v", err)
		os.Exit(-1)
	}

	if u.UsedPercent > max {
		fmt.Printf("MEMORY CRITICAL - %.2f%% used", u.UsedPercent)
		os.Exit(2)
	}

	if u.UsedPercent > warn {
		fmt.Printf("MEMORY WARNING - %.2f%% used", u.UsedPercent)
		os.Exit(1)
	}

	fmt.Printf("MEMORY OK - %.2f%% used", u.UsedPercent)

}
