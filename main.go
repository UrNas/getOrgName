package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

type host struct {
	name    string
	ip      string
	orgName string
}

var (
	domains = flag.String("domains", "", "domains names")
)
var sliceDomains []string

func main() {
	flag.Parse()
	if *domains == "" {
		log.Fatal("should more than 0 domain")
	}
	sliceDomains = strings.Split(*domains, ",")
	switch len(sliceDomains) {
	case 0:
		log.Fatalln("len should be more than ", len(sliceDomains))
	default:
		data := getFullHost()
		for _, value := range data {
			fmt.Printf("[*] %s %s %s \n", value.name, value.ip, value.orgName)
		}

	}
}

func getHostAndIpMap(hosts []string) map[string][]string {
	hostsIpMap := make(map[string][]string)
	for _, host := range hosts {
		ips, err := net.LookupHost(host)
		if err != nil {
			hostsIpMap[host] = nil
		} else {
			hostsIpMap[host] = getIp(ips)
		}
	}
	return hostsIpMap
}
func getIp(ips []string) []string {
	var hostips []string
	for _, value := range ips {
		ip := net.ParseIP(value)
		if err := ip.To4(); err != nil {
			hostips = append(hostips, value)
		}
	}
	return hostips
}
func getFullHost() []host {
	var hosts []host
	data := getHostAndIpMap(sliceDomains)
	for key, value := range data {
		d := getOrgName(value)
		if len(d) < 1 {
			continue
		} else {
			for _, value := range d {
				for k, v := range value {
					hosts = append(hosts, host{
						name:    key,
						ip:      v,
						orgName: k,
					})
				}
			}
		}
	}
	return hosts
}
func getOrgName(ips []string) []map[string]string {
	var orgnames []map[string]string
	for _, hostIp := range ips {
		conn, err := net.Dial("tcp", "whois.arin.net:43")
		if err != nil {
			log.Println("could not connect to whois: ", err)
		}
		if _, err := conn.Write([]byte(hostIp + "\r\n")); err != nil {
			log.Printf("could not write for this ip %s\n", hostIp)
		}
		s := bufio.NewScanner(conn)
		for s.Scan() {
			data := s.Text()
			if strings.Contains(data, "OrgName:") {
				org := strings.Trim(strings.Split(data, " ")[8], ",!")
				orgnames = append(orgnames, map[string]string{org: hostIp})
			}
		}
	}
	return orgnames
}
