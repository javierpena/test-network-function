package hostname

import (
	"github.com/redhat-nfvpe/test-network-function/internal/reel"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"time"
)

// Hostname provides a hostname test implemented using command line tool "hostname".
type Hostname struct {
	result  int
	timeout time.Duration
	args    []string
	// The hostname
	hostname string
}

const (
	// Command is the command name for the unix "hostname" command.
	Command = "hostname"
	// SuccessfulOutputRegex is the regular expression match for hostname output.
	SuccessfulOutputRegex = `.+`
)

// Args returns the command line args for the test.
func (h *Hostname) Args() []string {
	return h.args
}

// Timeout return the timeout for the test.
func (h *Hostname) Timeout() time.Duration {
	return h.timeout
}

// Result returns the test result.
func (h *Hostname) Result() int {
	return h.result
}

// ReelFirst returns a step which expects an hostname summary for the given device.
func (h *Hostname) ReelFirst() *reel.Step {
	return &reel.Step{
		Expect:  []string{SuccessfulOutputRegex},
		Timeout: h.timeout,
	}
}

// ReelMatch parses the hostname output and set the test result on match.
// Returns no step; the test is complete.
func (h *Hostname) ReelMatch(_ string, _ string, match string) *reel.Step {
	h.hostname = match
	h.result = tnf.SUCCESS
	return nil
}

// ReelTimeout does nothing;  hostname requires no explicit intervention for a timeout.
func (h *Hostname) ReelTimeout() *reel.Step {
	return nil
}

// ReelEOF does nothing;  hostname requires no explicit intervention for EOF.
func (h *Hostname) ReelEOF() {
}

// GetHostname returns the extracted hostname, if one is extracted.
func (h *Hostname) GetHostname() string {
	return h.hostname
}

// NewHostname creates a new `Hostname` test which runs the "hostname" command.
func NewHostname(timeout time.Duration) *Hostname {
	return &Hostname{
		result:  tnf.ERROR,
		timeout: timeout,
		args:    []string{Command},
	}
}