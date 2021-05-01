package ubuntuami

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	f, err := os.Open("test_data.jsonish")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	amis, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	if len(amis) != 775 {
		t.Fatalf("Expected 775 amis but got %d", len(amis))
	}

	expect := AMI{
		ID:             "ami-05fbbeaba8fdee29a",
		Region:         "af-south-1",
		ReleaseName:    "focal",
		ReleaseVersion: "20.04 LTS",
		Arch:           "arm64",
		InstanceType:   "hvm:ebs-ssd",
		ReleaseTS:      "20210429",
		ReleaseTime:    time.Date(2021, 04, 29, 0, 0, 0, 0, time.UTC),
		Link:           `<a href="https://console.aws.amazon.com/ec2/home?region=af-south-1#launchAmi=ami-05fbbeaba8fdee29a">ami-05fbbeaba8fdee29a</a>`,
		AKI:            "hvm",
	}

	if !cmp.Equal(amis[0], expect) {
		t.Fatal(cmp.Diff(amis[0], expect))
	}
}
