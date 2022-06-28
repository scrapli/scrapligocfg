package cisconxos_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type getConfigTestCase struct {
	description string
	source      string
	payloadFile string
}

func testGetConfig(testName string, testCase *getConfigTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		c, fileTransportObj := preparePlatform(t, testName, testCase.payloadFile)

		r, err := c.GetConfig(testCase.source)
		if err != nil {
			t.Errorf(
				"%s: encountered error running network GetConfig, error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		actualOut := r.Result
		actualIn := bytes.Join(fileTransportObj.Writes, []byte("\n"))

		if *update {
			writeGolden(t, testName, actualIn, actualOut)
		}

		expectedIn := readFile(t, fmt.Sprintf("golden/%s-in.txt", testName))
		expectedOut := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

		if !cmp.Equal(actualIn, expectedIn) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualIn,
				expectedIn,
			)
		}

		if !cmp.Equal(actualOut, string(expectedOut)) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestGetConfig(t *testing.T) {
	cases := map[string]*getConfigTestCase{
		"get-config-running": {
			description: "getconfig running",
			source:      "running",
			payloadFile: "getconfig-running.txt",
		},
		"get-config-startup": {
			description: "getconfig startup",
			source:      "startup",
			payloadFile: "getconfig-startup.txt",
		},
	}

	for testName, testCase := range cases {
		f := testGetConfig(testName, testCase)
		t.Run(testName, f)
	}
}
