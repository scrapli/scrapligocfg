package juniperjunos_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type loadConfigTestCase struct {
	description string
	payloadFile string
	configF     string
	replace     bool
}

func testCommitConfig(testName string, testCase *loadConfigTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		c, fileTransportObj := preparePlatform(t, testName, testCase.payloadFile)

		loadConfig(t, testName, c, testCase.configF, testCase.replace)

		r, err := c.CommitConfig()
		if err != nil {
			t.Errorf(
				"%s: encountered error running CommitConfig, error: %s",
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

func TestCommitConfig(t *testing.T) {
	cases := map[string]*loadConfigTestCase{
		"commit-config-replace": {
			description: "commit config replace",
			payloadFile: "commitconfig-replace.txt",
			configF:     "config-replace.txt",
			replace:     true,
		},
		"commit-config-merge": {
			description: "commit config merge",
			payloadFile: "commitconfig-merge.txt",
			configF:     "config-merge.txt",
			replace:     false,
		},
	}

	for testName, testCase := range cases {
		f := testCommitConfig(testName, testCase)
		t.Run(testName, f)
	}
}
