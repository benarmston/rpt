package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var update = flag.Bool("update", false, "update golden files")
var entryPoint = "main.go"

func TestDelay(t *testing.T) {
	var delay int64 = 1
	args := []string{"--delay", "1s", "3", "date", "+%s"}
	output, err := runBinary(args)
	assertExitCode(t, output, 0, err)

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var priorTimestamp int64
	for _, line := range lines {
		timestamp, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			t.Fatalf("failed to parse timestamp %s, %s", line, err)
		}
		if priorTimestamp != 0 {
			if priorTimestamp+delay != timestamp {
				t.Fatalf("expected timestamp %d to be %d seconds after prior timestamp %d", timestamp, delay, priorTimestamp)
			}
		}
		priorTimestamp = timestamp
	}
}

// Run the binary with specified args and compare output to golden files.
func TestGolden(t *testing.T) {
	tests := []struct {
		testName         string
		optionsAndArgs   []string
		fixture          string
		expectedExitCode int
	}{
		{
			"--help outputs expected help",
			[]string{"--help"},
			"help.golden",
			0,
		},
		{
			"Runs given command the given number times",
			[]string{"3", "echo", "The command"},
			"run-3-times.golden",
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			output, err := runBinary(tt.optionsAndArgs)
			assertExitCode(t, output, tt.expectedExitCode, err)
			if *update {
				writeFixture(t, tt.fixture, output)
			}
			actual := string(output)
			expected := loadFixture(t, tt.fixture)
			if !reflect.DeepEqual(actual, expected) {
				t.Logf("expected\n%s\n  got\n%s", expected, actual)
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(expected, actual, false)
				t.Log(dmp.DiffPrettyText(diffs))
				t.FailNow()
			}
		})
	}
}

func fixturePath(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), "testdata", fixture)
}

func writeFixture(t *testing.T, fixture string, content []byte) {
	err := os.WriteFile(fixturePath(t, fixture), content, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func loadFixture(t *testing.T, fixture string) string {
	content, err := os.ReadFile(fixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func runBinary(args []string) ([]byte, error) {
	fullArgs := append([]string{"run", entryPoint}, args...)
	cmd := exec.Command("go", fullArgs...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
	return cmd.CombinedOutput()
}

func assertExitCode(t *testing.T, output []byte, expectedExitCode int, err error) {
	t.Helper()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != expectedExitCode {
				t.Fatalf("output:\n%s\nerror:\n%s\n", output, err)
			}
		} else {
			t.Fatalf("output:\n%s\nerror:\n%s\n", output, err)
		}
	}
}
