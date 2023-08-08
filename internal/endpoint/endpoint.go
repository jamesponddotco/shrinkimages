// Package endpoint contains the endpoint definitions for the service.
package endpoint

// CurrentAPIVersion is the current version of the API.
const CurrentAPIVersion string = "v1"

const (
	// Root is the endpoint for the root handler.
	Root string = "/"

	// Ping is the endpoint for the ping handler.
	Ping string = Root + CurrentAPIVersion + "/ping"

	// Shrink is the endpoint for the shrink handler.
	Shrink string = Root + CurrentAPIVersion + "/shrink"
)
