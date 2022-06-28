package ciscoiosxr_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/transport"
	"github.com/scrapli/scrapligocfg"
)

var update = flag.Bool( //nolint
	"update",
	false,
	"update the golden files",
)

func resolveFile(t *testing.T, f string) string {
	f, err := filepath.Abs(fmt.Sprintf("./test-fixtures/%s", f))
	if err != nil {
		t.Fatal(err)
	}

	return f
}

func readFile(t *testing.T, f string) []byte {
	b, err := os.ReadFile(fmt.Sprintf("./test-fixtures/%s", f))
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func writeGolden(t *testing.T, testName string, actualIn []byte, actualOut string) {
	goldenOut := filepath.Join("test-fixtures", "golden", testName+"-out.txt")
	goldenIn := filepath.Join("test-fixtures", "golden", testName+"-in.txt")

	err := os.WriteFile(goldenOut, []byte(actualOut), 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(goldenIn, actualIn, 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
}

func preparePlatform(
	t *testing.T,
	testName,
	payloadFile string,
) (*scrapligocfg.Cfg, *transport.File) {
	p, err := platform.NewPlatform(
		"cisco_iosxr",
		"dummy",
		options.WithTransportType(transport.FileTransport),
		options.WithFileTransportFile(resolveFile(t, payloadFile)),
		options.WithTransportReadSize(1),
		options.WithReadDelay(0),
	)
	if err != nil {
		t.Fatalf(
			"%s: encountered error creating network Platform instance, error: %s",
			testName,
			err,
		)
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
		t.Fatalf("%s: encountered error creating network Driver, error: %s", testName, err)
	}

	err = d.Channel.Open()
	if err != nil {
		t.Fatalf("%s: encountered error opening Channel, error: %s", testName, err)
	}

	fileTransportObj, ok := d.Transport.Impl.(*transport.File)
	if !ok {
		t.Fatalf("transport implementation is not Transport File")
	}

	c, err := scrapligocfg.NewCfg(
		d,
		"cisco_iosxr",
	)
	if err != nil {
		t.Fatalf(
			"%s: encountered error fetching Cfg instance from factory, error: %s",
			testName,
			err,
		)
	}

	err = c.Prepare()
	if err != nil {
		t.Fatalf("%s: encountered error preparing connection, error: %s", testName, err)
	}

	return c, fileTransportObj
}

func loadConfig(
	t *testing.T,
	testName string,
	c *scrapligocfg.Cfg,
	configF string,
	replace bool,
) {
	config := string(readFile(t, configF))

	r, err := c.LoadConfig(config, replace)
	if err != nil {
		t.Fatalf("%s: failed loading config, error: %s", testName, err)
	}

	if r.Failed != nil {
		t.Fatalf("%s: response object indicates failure",
			testName)
	}
}
