package scanner

import (
	"fmt"
	"io"
	"log"

	"github.com/dutchcoders/go-clamd"
)

// Scanner defines the interface for virus scanning.
type Scanner interface {
	// Scan checks the content for viruses.
	// Returns safe (bool), threat name (string), and error.
	Scan(r io.Reader) (bool, string, error)
}

// ClamAVScanner implements Scanner using ClamAV.
type ClamAVScanner struct {
	clam *clamd.Clamd
}

// NewClamAVScanner creates a new ClamAV scanner.
// address should be like "tcp://localhost:3310" or "unix:///tmp/clamd.socket"
func NewClamAVScanner(address string) (*ClamAVScanner, error) {
	c := clamd.NewClamd(address)

	// Optional: Test connection.
	// We might not want to fail hard if ClamAV is temporarily down,
	// but for initialization it's good to know.
	if err := c.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to clamav at %s: %w", address, err)
	}

	return &ClamAVScanner{clam: c}, nil
}

func (s *ClamAVScanner) Scan(r io.Reader) (bool, string, error) {
	abort := make(chan bool)
	defer close(abort)

	response, err := s.clam.ScanStream(r, abort)
	if err != nil {
		return false, "", fmt.Errorf("failed to start scan: %w", err)
	}

	for result := range response {
		switch result.Status {
		case clamd.RES_FOUND:
			return false, result.Description, nil
		case clamd.RES_ERROR:
			return false, "", fmt.Errorf("scan error: %s", result.Description)
		case clamd.RES_OK:
			// Continue waiting for other results or completion
		}
	}

	return true, "", nil
}

// NoOpScanner is a dummy scanner that always reports safe.
type NoOpScanner struct{}

func (s *NoOpScanner) Scan(r io.Reader) (bool, string, error) {
	log.Println("Warning: Virus scanning is disabled or using NoOpScanner.")
	return true, "", nil
}
