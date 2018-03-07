package diagnose

import (
	"fmt"
	"time"

	"encoding/json"

	"sync"
)

// HealthStatus is type for health status
type HealthStatus string

const (
	// StatusOK means the component is OK
	StatusOK HealthStatus = "OK"
	// StatusError means there is an error with the component
	StatusError HealthStatus = "ERROR"
)

// ComponentReport define each component's report
type ComponentReport struct {
	Status     HealthStatus  `json:"status"`
	Name       string        `json:"name"`
	Message    string        `json:"message"`
	Suggestion string        `json:"suggestion"`
	Latency    time.Duration `json:"latency"`
}

type HealthReport struct {
	Status  HealthStatus      `json:"status"`
	Details []ComponentReport `json:"details"`
}

// Component is interface for a component health check
type Component interface {
	Diagnose() ComponentReport
}

// HealthChecker is main health diagnoser
type HealthChecker struct {
	Components []Component
}

// SimpleCheck simplified health check function for components where
// if the function returns an error it should be able to distinguish
// the health status of the component
type SimpleCheck func() error

// Add component for checking
func (h *HealthChecker) Add(com Component) *HealthChecker {
	if h.Components == nil {
		h.Components = make([]Component, 0, 2)
	}
	h.Components = append(h.Components, com)
	return h
}

// Add a new component report
func (h *HealthReport) Add(report ComponentReport) {
	if h.Details == nil {
		h.Details = make([]ComponentReport, 0, 3)
	}
	h.Details = append(h.Details, report)
}

// Check starts health check of components
func (h *HealthChecker) Check() HealthReport {
	report := &HealthReport{Status: StatusOK}
	if h.Components == nil {
		return *report
	}
	wait := sync.WaitGroup{}
	for _, c := range h.Components {
		wait.Add(1)
		go func(component Component) {
			diagnose := component.Diagnose()
			if diagnose.Status == StatusError && report.Status != StatusError {
				report.Status = StatusError
			}
			report.Add(diagnose)
			wait.Done()
		}(c)
	}
	wait.Wait()
	return *report
}

// NewReport constructor
func NewReport(component string) *ComponentReport {
	return &ComponentReport{
		Status:  StatusOK,
		Name:    component,
		Message: "ok",
	}
}

// SimpleDiagnose is a simple diagnose check based on a check function that returns an error
func SimpleDiagnose(componentName string, check SimpleCheck) ComponentReport {
	var (
		err   error
		start time.Time
	)
	report := NewReport(componentName)
	start = time.Now()
	err = check()
	report.Check(err, "OK", "Check configuration")
	report.AddLatency(start)
	return *report
}

// Check for error and add custom message
func (c *ComponentReport) Check(err error, message, suggestion string) {
	if err != nil {
		c.Status = StatusError
		c.Message = fmt.Sprintf("%s: \"%s\"", message, err.Error())
		c.Suggestion = suggestion
	}
}

// AddLatency add a latency for the start time
func (c *ComponentReport) AddLatency(start time.Time) {
	duration := time.Since(start)
	if c.Latency == time.Duration(0) || c.Latency < duration {
		c.Latency = duration
	}
}

// MarshalJSON is custom json formatter
func (c *ComponentReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status     HealthStatus `json:"status"`
		Name       string       `json:"name"`
		Message    string       `json:"message"`
		Suggestion string       `json:"suggestion"`
		Latency    string       `json:"latency"`
	}{
		c.Status,
		c.Name,
		c.Message,
		c.Suggestion,
		c.Latency.String(),
	})
}

// New constructor function for HealthChecker
func New() (*HealthChecker, error) {
	return &HealthChecker{}, nil
}
