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

	if len(amis) != 573 {
		t.Fatalf("Expected 573 amis but got %d", len(amis))
	}

	expect := AMI{
		ID:             "ami-0bd6ed93d44bb24ae",
		Region:         "us-west-1",
		ReleaseName:    "lunar",
		ReleaseVersion: "23.04",
		Arch:           "amd64",
		InstanceType:   "hvm:ebs-ssd",
		ReleaseTS:      "20230420",
		ReleaseTime:    time.Date(2023, 04, 20, 0, 0, 0, 0, time.UTC),
		Link:           `<a href="https://console.aws.amazon.com/ec2/home?region=us-west-1#launchAmi=ami-0bd6ed93d44bb24ae">ami-0bd6ed93d44bb24ae</a>`,
		AKI:            "hvm",
	}

	if !cmp.Equal(amis[0], expect) {
		t.Fatal(cmp.Diff(amis[0], expect))
	}
}
