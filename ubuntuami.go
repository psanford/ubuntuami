package ubuntuami

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	// URL to fetch AMI data from
	ImagesURL = "https://cloud-images.ubuntu.com/locator/ec2/releasesTable"
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

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var aa aaData
	err = json.Unmarshal(data, &aa)
	if err != nil {
		return nil, err
	}

	for _, amiData := range aa.AAData {
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

		ami.ReleaseName = strings.SplitN(strings.ToLower(ami.ReleaseName), " ", 2)[0]

		out = append(out, ami)
	}

	return out, nil
}

type aaData struct {
	AAData  [][]string `json:"aaData"`
	Updated string     `json:"updated"`
}
