package dn42

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Filter struct {
	Nr      uint32
	Action  string
	Prefix  *net.IPNet
	MinLen  byte
	MaxLen  byte
	Comment string
}
type RSPLObject map[string][]string
type Route struct {
	Prefix    string `json:"prefix"`
	MaxLength byte   `json:"maxLength"`
	Asn       string `json:"asn"`
	Ta        string `json:"ta"`
}

func ParseRoutes(dir string, filters []Filter) ([]Route, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	ret := make([]Route, 0)
	for _, f := range files {
		record, err := ParseObject(path.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		var prefixIP net.IP
		var prefixNet *net.IPNet
		if val, ok := record["route"]; ok {
			ip, net, err := net.ParseCIDR(val[0])
			if err != nil {
				return nil, err
			}
			prefixIP = ip
			prefixNet = net
		}
		if val, ok := record["route6"]; ok {
			ip, net, err := net.ParseCIDR(val[0])
			if err != nil {
				return nil, err
			}
			prefixIP = ip
			prefixNet = net
		}
		var max byte = 0
		permitted := false
		for _, filter := range filters {
			if filter.Prefix.Contains(prefixIP) {
				if filter.Action == "permit" {
					permitted = true
					break
				}
			}
		}
		if permitted != true {
			continue
		}
		for _, filter := range filters {
			if filter.Prefix.Contains(prefixIP) {
				if filter.Action == "permit" {
					max = filter.MaxLen
				}
			}
		}
		if val, ok := record["max-length"]; ok {
			maxLength, err := strconv.ParseInt(val[0], 10, 8)
			if err != nil {
				return nil, err
			}
			if byte(maxLength) < max {
				max = byte(maxLength)
			}
		}
		for _, origin := range record["origin"] {
			var route Route
			route.MaxLength = max
			route.Prefix = prefixNet.String()
			route.Asn = origin
			ret = append(ret, route)
		}
	}
	return ret, nil
}

func ParseObject(path string) (RSPLObject, error) {
	var lastAttr string
	var lastArr []string
	ret := make(RSPLObject)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		key := strings.TrimSpace(line[0:20])
		key = key[:len(key)-1]
		value := line[20:]
		switch key[0:1] {
		case "\t":
		case " ":
			lastArr[len(lastArr)-1] += value
			break
		case "+":
			lastArr[len(lastArr)-1] += "\n"
			break
		default:
			lastArr = make([]string, 0)
			lastAttr = key
			if ret[lastAttr] == nil {
				ret[lastAttr] = make([]string, 1)
				ret[lastAttr][0] = value
			} else {
				ret[lastAttr] = append(ret[lastAttr], value)
			}
		}
	}
	return ret, nil
}

func ParseFilter(path string) ([]Filter, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	filters := make([]Filter, 0, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}
		var filter Filter
		arr := regexp.MustCompile(" +").Split(line, 6)

		nr, err := strconv.ParseUint(strings.TrimSpace(arr[0]), 10, 32)
		if err != nil {
			return nil, err
		}
		filter.Nr = uint32(nr)

		filter.Action = strings.TrimSpace(arr[1])

		_, prefix, err := net.ParseCIDR(strings.TrimSpace(arr[2]))
		if err != nil {
			return nil, err
		}
		filter.Prefix = prefix

		minLen, err := strconv.ParseUint(strings.TrimSpace(arr[3]), 10, 8)
		if err != nil {
			return nil, err
		}
		filter.MinLen = byte(minLen)

		maxLen, err := strconv.ParseUint(strings.TrimSpace(arr[4]), 10, 8)
		if err != nil {
			return nil, err
		}
		filter.MaxLen = byte(maxLen)

		filter.Comment = arr[5]
		filters = append(filters, filter)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return filters, nil
}
