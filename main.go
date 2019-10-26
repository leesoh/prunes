package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dchest/uniuri"
	"github.com/miekg/dns"
)

func main() {
	sf := flag.String("s", "", "Subdomains to check against each resolver")
	domain := flag.String("d", "", "Non-existent domain to check")
	concurrency := flag.Int("c", 10, "Number of concurrent checks to run")
	flag.Parse()
	var sub Subdomain
	sub.Domain = *domain
	sub.Subfile = *sf
	sub.Load()
	var wg sync.WaitGroup
	resolvers := make(chan string)
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			for r := range resolvers {
				sub.Process(r)
			}
		}()
		wg.Done()
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		res := strings.ToLower(strings.TrimSpace(sc.Text()))
		if res == "" {
			continue
		}
		resolvers <- res
	}
	close(resolvers)
	wg.Wait()
}

type Subdomain struct {
	Domain  string
	Subfile string
	Subs    []string
}

func (s *Subdomain) Load() {
	subfile, err := os.Open(s.Subfile)
	if err != nil {
		panic(err)
	}
	defer subfile.Close()
	sc := bufio.NewScanner(subfile)
	for sc.Scan() {
		s.Subs = append(s.Subs, sc.Text())
	}
	// add a totally random subdomain
	s.Subs = append(s.Subs, uniuri.New())
}

func (s Subdomain) Process(resolver string) {
	for _, ss := range s.Subs {
		record := ss + "." + s.Domain
		resp, err := LookupHost(resolver, record)
		if err != nil {
			continue
		}
		for _, r := range resp {
			fmt.Printf("%s: %s => %s\n", resolver, record, r)
		}
	}
}

func LookupHost(resolver, record string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(record), dns.TypeA)
	m.RecursionDesired = true
	result, _, err := c.Exchange(m, resolver+":53")
	if err != nil {
		return nil, err
	}
	if len(result.Answer) == 0 {
		return nil, fmt.Errorf("no answer for %s", record)
	}
	var records []string
	for _, a := range result.Answer {
		if r, ok := a.(*dns.A); ok {
			records = append(records, r.A.String())
		}
	}
	return records, nil
}
