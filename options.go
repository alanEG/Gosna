package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

//declare flag.Usage values
func init() {
	flag.Usage = func() {
		h := "Usage: ./change [options]\n"
		h += "   -run          Run type [add,check]\n"
		h += "   -timeout      Requests timeout (default 5)\n"
		h += "   -thread       Requests thread\n"
		h += "   -header       Requests header\n"
		h += "   -dynamic      Check dynamic (default false)\n"
		h += "   -config       Config file   (default ~/.gosna_config.json)\n"
		h += "   -no-color     Disable color\n\n"

		fmt.Fprintf(os.Stderr, h)
	}
}

func (i *arrayFlags) String() string {

	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func options_parse() {

	//handling header array
	// ./change -headers value1 -headers value2
	flag.StringVar(&Run, "run", "", "")

	directoryHome, _ := os.UserHomeDir()

	flag.StringVar(&configFileFlag, "config", directoryHome+"/.gosna_config.json", "")
	flag.BoolVar(&flagDynmaic, "dynamic", false, "")
	flag.BoolVar(&flagNoColor, "no-color", false, "")
	flag.IntVar(&flagTimeout, "timeout", 5, "")
	flag.IntVar(&Thread, "thread", 5, "")
	flag.Var(&Headers, "header", "")
	flag.Var(&Headers, "H", "")
	flag.Parse()

	Header = parse_headers(Headers)

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}
}

//handling headers flags
func parse_headers(headers []string) map[string]string {

	headersv := make(map[string]string)

	for _, header := range headers {
		hv := strings.Split(header, `:`)
		headersv[hv[0]] = hv[1]

	}

	return headersv
}
