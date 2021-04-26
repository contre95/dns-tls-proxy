package resolver

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	"tls-dns-proxy/pkg/domain/proxy"
)

const cflRootCert = `-----BEGIN CERTIFICATE-----
MIIEQzCCAyugAwIBAgIQCidf5wTW7ssj1c1bSxpOBDANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD
QTAeFw0yMDA5MjMwMDAwMDBaFw0zMDA5MjIyMzU5NTlaMFYxCzAJBgNVBAYTAlVT
MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxMDAuBgNVBAMTJ0RpZ2lDZXJ0IFRMUyBI
eWJyaWQgRUNDIFNIQTM4NCAyMDIwIENBMTB2MBAGByqGSM49AgEGBSuBBAAiA2IA
BMEbxppbmNmkKaDp1AS12+umsmxVwP/tmMZJLwYnUcu/cMEFesOxnYeJuq20ExfJ
qLSDyLiQ0cx0NTY8g3KwtdD3ImnI8YDEe0CPz2iHJlw5ifFNkU3aiYvkA8ND5b8v
c6OCAa4wggGqMB0GA1UdDgQWBBQKvAgpF4ylOW16Ds4zxy6z7fvDejAfBgNVHSME
GDAWgBQD3lA1VtFMu2bwo+IbG8OXsj3RVTAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0l
BBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBIGA1UdEwEB/wQIMAYBAf8CAQAwdgYI
KwYBBQUHAQEEajBoMCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdpY2VydC5j
b20wQAYIKwYBBQUHMAKGNGh0dHA6Ly9jYWNlcnRzLmRpZ2ljZXJ0LmNvbS9EaWdp
Q2VydEdsb2JhbFJvb3RDQS5jcnQwewYDVR0fBHQwcjA3oDWgM4YxaHR0cDovL2Ny
bDMuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNybDA3oDWgM4Yx
aHR0cDovL2NybDQuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNy
bDAwBgNVHSAEKTAnMAcGBWeBDAEBMAgGBmeBDAECATAIBgZngQwBAgIwCAYGZ4EM
AQIDMA0GCSqGSIb3DQEBDAUAA4IBAQDeOpcbhb17jApY4+PwCwYAeq9EYyp/3YFt
ERim+vc4YLGwOWK9uHsu8AjJkltz32WQt960V6zALxyZZ02LXvIBoa33llPN1d9R
JzcGRvJvPDGJLEoWKRGC5+23QhST4Nlg+j8cZMsywzEXJNmvPlVv/w+AbxsBCMqk
BGPI2lNM8hkmxPad31z6n58SXqJdH/bYF462YvgdgbYKOytobPAyTgr3mYI5sUje
CzqJx1+NLyc8nAK8Ib2HxnC+IrrWzfRLvVNve8KaN9EtBH7TuMwNW4SpDCmGr6fY
1h3tDjHhkTb9PA36zoaJzu0cIw265vZt6hCmYWJC+/j+fgZwcPwL
-----END CERTIFICATE-----
`

type cloudFlare struct {
	ip          string
	port        int
	rootCert    string
	readTimeOut uint
}

func NewCloudFlareResolver(ip string, port int, rto uint) proxy.Resolver {
	return &cloudFlare{ip, port, cflRootCert, rto}
}

func (cfl *cloudFlare) GetTLSConnection() (*tls.Conn, error) {
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(cfl.rootCert)) {
		log.Println("Fail to parse rootCert")
		return nil, errors.New("Fail to parse rootCert")
	}
	dnsCloudFlareConn, err := tls.Dial("tcp", cfl.ip+":"+strconv.Itoa(cfl.port), &tls.Config{
		RootCAs: roots,
	})
	if err != nil {
		log.Println("Error connecting to CloudFlare")
		return nil, err
	}
	_ = dnsCloudFlareConn.SetReadDeadline(time.Now().Add(time.Duration(cfl.readTimeOut) * time.Millisecond))
	return dnsCloudFlareConn, nil
}

func (cfl *cloudFlare) Solve(um proxy.UnsolvedMsg) (proxy.SolvedMsg, error) {
	// Levanto una conexión con CF
	conn, err := cfl.GetTLSConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, e := conn.Write(um)
	if e != nil {
		fmt.Printf("%v", e)
	}
	var reply [2045]byte
	n, er := conn.Read(reply[:])
	if er != nil {
		fmt.Printf("Could read response from CloudFlare: %v \n", er)
	} else {
		log.Println("Succesfuly fullfiled the request")
	}
	return reply[:n], nil
}
