package jpena

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/test-network-function/test-network-function/pkg/tnf"
	"github.com/test-network-function/test-network-function/pkg/tnf/dependencies"
	"github.com/test-network-function/test-network-function/pkg/tnf/identifier"
	"github.com/test-network-function/test-network-function/pkg/tnf/reel"
)

const (
	VersionRegex = `(?m)Red Hat Enterprise Linux release (\d+)\.(\d+).*$`
)

var (
	ReleaseCommand = fmt.Sprintf("%s '^Red Hat' /etc/redhat-release", dependencies.GrepBinaryName)
)

// Release is an implementation of tnf.Test used to determine whether a container is based on Red Hat UBI.
type Release struct {
	// result is the result of the test.
	result int
	// timeout is the timeout duration for the test.
	timeout time.Duration
	// args stores the command and arguments.
	args []string
	// major contains the Red Hat major release, 0 if non-existent
	major int
	// minor contains the Red Hat minor release, 0 if non-existent
	minor int
	// isRedHatBased contains whether the container is based on Red Hat technologies.
	isRedHatBased bool
}

// Args returns the command line arguments for the test.
func (r *Release) Args() []string {
	return r.args
}

// GetIdentifier returns the tnf.Test specific identifier.
func (r *Release) GetIdentifier() identifier.Identifier {
	return identifier.VersionIdentifier
}

// Timeout returns the timeout for the test.
func (r *Release) Timeout() time.Duration {
	return r.timeout
}

// Result returns the test result.
func (r *Release) Result() int {
	return r.result
}

// ReelFirst returns a reel.Step which expects output from running the Args command.
func (r *Release) ReelFirst() *reel.Step {
	return &reel.Step{
		Execute: ReleaseCommand,
		Expect:  []string{VersionRegex},
		Timeout: r.timeout,
	}
}

// ReelMatch determines whether the container is based on Red Hat technologies through pattern matching logic.
func (r *Release) ReelMatch(pattern, _, match string) *reel.Step {
	if pattern == VersionRegex {
		// If the above conditional is not triggered, it can be deduced that we have matched the VersionRegex.
		r.result = tnf.SUCCESS
		r.isRedHatBased = true
		// Find major and minor versions
		re := regexp.MustCompile(pattern)
		matched := re.FindStringSubmatch(match)
		r.major, _ = strconv.Atoi(matched[1])
		r.minor, _ = strconv.Atoi(matched[2])
		fmt.Printf("JPENA we got version %d.%d\n", r.major, r.minor)
	} else {
		r.result = tnf.ERROR
		r.isRedHatBased = false
	}
	return nil
}

// ReelTimeout does nothing;  no intervention is needed for a timeout.
func (r *Release) ReelTimeout() *reel.Step {
	return nil
}

// ReelEOF does nothing;  no intervention is needed for EOF.
func (r *Release) ReelEOF() {
}

// NewRelease create a new Release tnf.Test.
func NewRelease(timeout time.Duration) *Release {
	return &Release{result: tnf.ERROR, timeout: timeout, args: make([]string, 0)}
}
