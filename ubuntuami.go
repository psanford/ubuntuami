package ubuntuami

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	// URL to fetch AMI data from
	ImagesURL = "http://cloud-images.ubuntu.com/locator/ec2/releasesTable"
)

type AMI struct {
	ID             string    // AMI ID (e.g. ami-05fbbeaba8fdee29a)
	Region         string    // AWS Region
	ReleaseName    string    // Ubuntu Release Name
	ReleaseVersion string    // Ubuntu Release Version
	Arch           string    // Hardware Architecture (amd64, arm64, i386)
	InstanceType   string    // (hvm:instance-store, hvm:ebs-ssd,
	ReleaseTS      string    // AMI Release Timstamp String
	ReleaseTime    time.Time // AMI Release Build Timestamp
	Link           string    // HTML link to AMI
	AKI            string    // Kernel ID (or hvm)
}

// Fetch the list of Ubuntu AMIs from cloud-images.ubuntu.com. The result
// of Fetch should be agressively cached.
func Fetch() ([]AMI, error) {
	resp, err := http.Get(ImagesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected http status: %d", resp.StatusCode)
	}

	return Parse(resp.Body)
}

// Parse the response from cloud-images releaseTable api.
func Parse(r io.Reader) ([]AMI, error) {
	var out []AMI
	bufReader := bufio.NewReader(r)
	for i := 0; ; i++ {
		line, err := bufReader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		line = bytes.TrimSpace(line)

		if !bytes.HasPrefix(line, []byte("[")) || !bytes.HasSuffix(line, []byte("],")) {
			continue
		}

		line = bytes.TrimSuffix(line, []byte(","))

		var amiData []string
		err = json.Unmarshal(line, &amiData)
		if err != nil {
			return nil, fmt.Errorf("decode json error on line %d: %s", i, err)
		}

		if got, expect := len(amiData), 8; got != expect {
			return nil, fmt.Errorf("AMI data with unexpected number of fields %d, expected %d", got, expect)
		}

		ami := AMI{
			Region:         amiData[0],
			ReleaseName:    amiData[1],
			ReleaseVersion: amiData[2],
			Arch:           amiData[3],
			InstanceType:   amiData[4],
			ReleaseTS:      amiData[5],
			Link:           amiData[6],
			AKI:            amiData[7],
		}

		id := ami.Link[strings.Index(ami.Link, ">")+1:]
		ami.ID = strings.TrimSuffix(id, "</a>")
		ami.ReleaseTime, _ = time.Parse("20060102", ami.ReleaseTS)

		out = append(out, ami)
	}

	return out, nil
}
