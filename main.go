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
	sf := flag.String("s", "subdomains.txt", "Subdomains to check against each resolver")
	domain := flag.String("d", GetRandomString()+".com", "Non-existent domain to check")
	concurrency := flag.Int("c", 20, "Number of concurrent checks to run")
	details := flag.Bool("details", false, "List servers and non-existent domain responses")
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
			defer wg.Done()
			for r := range resolvers {
				sub.Process(r, *details)
			}
		}()
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
	s.Subs = append(s.Subs, GetRandomString())
}

func (s Subdomain) Process(resolver string, details bool) {
	// a server ends up here if it resolves ANY record we provide
	badNS := make(map[string]bool)
	for _, ss := range s.Subs {
		record := ss + "." + s.Domain
		resp, err := LookupHost(resolver, record)
		// the only error returned is one for no answer
		if err != nil {
			continue
		}
		// the server returned something, this is bad
		badNS[resolver] = true
		for _, r := range resp {
			// only print this if details are wanted
			if details {
				fmt.Printf("%s ::: %s => %s\n", resolver, record, r)
			}
		}
	}
	if !details && !badNS[resolver] {
		fmt.Println(resolver)
	}
}

func GetRandomString() string {
	return strings.ToLower(uniuri.New())
}

func LookupHost(resolver, record string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(record), dns.TypeA)
	m.RecursionDesired = true
	result, _, err := c.Exchange(m, resolver+":53")
	if err != nil {
		return nil, nil
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
