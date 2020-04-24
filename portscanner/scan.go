package portscanner

import (
	"encoding/xml"
	"log"
	"os/exec"
)

type Nmaprun struct {
	Host struct {
		Ports struct {
			Port []struct {
				Protocol string `xml:"protocol,attr"`
				Portid   string `xml:"portid,attr"`
				State    struct {
					State     string `xml:"state,attr"`
					Reason    string `xml:"reason,attr"`
					ReasonTtl string `xml:"reason_ttl,attr"`
				} `xml:"state"`
				Service struct {
					Name string `xml:"name,attr"`
				} `xml:"service"`
			} `xml:"port"`
		} `xml:"ports"`
	} `xml:"host"`
	Runstats struct {
		Hosts struct {
			Up   string `xml:"up,attr"`
			Down string `xml:"down,attr"`
		} `xml:"hosts"`
	} `xml:"runstats"`
}

func ScanHost(ipv4 string) (map[string]map[string]string, error) {
	cmd := exec.Command("nmap", "-oX", "-", "-p", "3-400", ipv4)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.CombinedOutput() failed with %s\n", err)
	}

	b := &Nmaprun{}

	xml.Unmarshal(out, b)

	m := make(map[string]map[string]string)
	m[ipv4] = make(map[string]string)
	for _, v := range b.Host.Ports.Port {
		m[ipv4][v.Portid] = v.State.State
	}
	return m, err
}
