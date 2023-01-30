// Bind Zonefile to PowerDNS converter
//
// Written in 2023 by Taras Kondratyuk (https://getpinga.com)
//
// @license MIT

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/miekg/dns"
)

type Record struct {
	Type     uint16
	TTL      uint32
	Data     string
	Priority uint16
}

func loadRecords(filename string) ([]Record, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config struct {
		Records []Record `json:"records"`
	}
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return config.Records, nil
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	records, err := loadRecords("records.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, record := range records {
		switch record.Type {
		case dns.TypeA:
			rr := &dns.A{
				Hdr: dns.RR_Header{
					Name:   r.Question[0].Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    record.TTL,
				},
				A: net.ParseIP(record.Data),
			}
			m.Answer = append(m.Answer, rr)
		case dns.TypeMX:
			rr := &dns.MX{
				Hdr: dns.RR_Header{
					Name:   r.Question[0].Name,
					Rrtype: dns.TypeMX,
					Class:  dns.ClassINET,
					Ttl:    record.TTL,
				},
				Mx:     record.Data,
				Priority: record.Priority,
			}
			m.Answer = append(m.Answer, rr)
		default:
			rr := &dns.RR{
				Name:   r.Question[0].Name,
				Rrtype: record.Type,
				Class:  dns.ClassINET,
				Ttl:    record.TTL,
				String: fmt.Sprintf("%d %s", record.TTL, record.Data),
			}
			m.Answer = append(m.Answer, rr)
		}
	}

	w.WriteMsg(m)
}

func main() {
	dns.HandleFunc(".", handleRequest)

	err := dns.ListenAndServe(":53", "udp", nil)
	if err != nil {
		fmt.Println(err)
	}
}
