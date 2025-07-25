// dns-query-function.go is a Cloud Function that performs DNS queries
package dnsqueryfunction

import (
	"fmt"
	"net/http"

	"github.com/miekg/dns"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("dnsQueryGet", dnsQueryGet)
}

// dnsQueryGet is an HTTP Cloud Function handler that performs DNS queries based on
// user-provided parameters. It expects two query parameters:
//   - "fqdn": the fully qualified domain name to query.
//   - "querytype": the DNS record type (e.g., "A", "AAAA", "MX", "TXT").
//
// The function validates the input parameters and constructs a DNS query using the
// github.com/miekg/dns library. It sends the query to Google's public DNS server (8.8.8.8)
// and returns the results in the HTTP response. If the parameters are missing or invalid,
// or if the DNS query fails, it responds with an appropriate HTTP error code and message.
//
// Example request:
//
//	GET /?fqdn=example.com&querytype=A
//
// Example response:
//
//	example.com.	3600	IN	A	93.184.216.34
func dnsQueryGet(w http.ResponseWriter, r *http.Request) {
	fqdn := r.URL.Query().Get("fqdn")
	qtype := r.URL.Query().Get("querytype")
	if fqdn == "" || qtype == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing fqdn or querytype parameter")
		return
	}

	dnsType, ok := dns.StringToType[qtype]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid querytype parameter")
		return
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fqdn), dnsType)

	c := new(dns.Client)
	dnsResp, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DNS query failed: %v", err)
		return
	}

	for _, ans := range dnsResp.Answer {
		fmt.Fprintf(w, "%v\n", ans)
	}
}
