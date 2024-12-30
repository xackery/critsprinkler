package dps

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/xackery/critsprinkler/reporter"
	"github.com/xackery/critsprinkler/tracker"
)

var timeRegex = regexp.MustCompile(`\[(.*?)\]`)

func TestReporter(t *testing.T) {

	_, err := tracker.New("eqlog_Shin_thj.txt")
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	_, err = reporter.New()
	if err != nil {
		t.Fatalf("failed to create reporter: %v", err)
	}

	r, err := os.Open("eqlog_Shin_thj.txt")
	if err != nil {
		t.Fatalf("failed to open log file: %v", err)
	}
	defer r.Close()

	err = New()
	if err != nil {
		t.Fatalf("failed to create dps: %v", err)
	}

	// iterate each line to newline

	start := time.Now()
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("failed to read line: %v", err)
		}
		match := timeRegex.FindStringSubmatch(line)
		if len(match) < 2 {
			continue
		}
		event, err := time.Parse("Mon Jan 02 15:04:05 2006", match[1])
		if err != nil {
			continue
		}
		onLine(event, line)
	}

	t.Logf("processed log in %v", time.Since(start))

}
