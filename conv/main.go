package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HaesungSeo/yamlconv"
	"gopkg.in/yaml.v2"
)

type SearchKey []string

func (m *SearchKey) String() string {
	return fmt.Sprint(*m)
}

func (m *SearchKey) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func main() {
	var searchKeys SearchKey
	yamlpath := flag.String("f", "/dev/stdin", "yaml file")
	ofmt := flag.String("o", "text", "output format, one of text, json")
	flag.Var(&searchKeys, "s", "define search keys multiple times, e.g. -s sriov -s [0] -s ip")
	flag.Parse()

	// read yaml file into buffer
	var filebuf []byte
	if yamlpath != nil {
		filename, err := filepath.Abs(*yamlpath)
		if err != nil {
			panic(err.Error())
		}
		filebuf, err = os.ReadFile(filename)
		if err != nil {
			panic(err.Error())
		}
	} else {
		filename, err := filepath.Abs("/dev/stdin")
		if err != nil {
			panic(err.Error())
		}
		filebuf, err = os.ReadFile(filename)
		if err != nil {
			panic(err.Error())
		}
	}

	// remove comments and empty lines, etc
	lines := strings.Split(string(filebuf), "\n")
	rline := make([]string, 0)
	for _, line := range lines {
		nline := strings.TrimRight(line, "\r\n")
		if len(nline) == 0 {
			continue
		} else if nline[0] == '#' {
			continue
		}
		rline = append(rline, nline)
	}

	// '\n' to \n, if input is single line of string
	if len(rline) <= 1 {
		nbuf := strings.ReplaceAll(string(filebuf), `\n`, "\n")
		filebuf = []byte(nbuf)
	}

	// parse it
	var data interface{}
	err := yaml.Unmarshal(filebuf, &data)
	if err != nil {
		panic(fmt.Sprintf("ERROR: %s\n", err.Error()))
	}

	switch *ofmt {
	case "text":
		data, err := yamlconv.Search(data, searchKeys)
		if err != nil {
			panic(fmt.Sprintf("ERROR: %s\n", err.Error()))
		}
		yamlconv.Print(data, "  ")
	case "json":
		ret, err := yamlconv.MarshalJson(data, searchKeys)
		if err != nil {
			panic(fmt.Sprintf("ERROR: %s\n", err.Error()))
		}
		fmt.Printf("%s\n", ret)
	}
}
