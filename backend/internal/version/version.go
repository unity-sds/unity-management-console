package version

// Version is set during build time via -ldflags
var Version = "dev"

// GetVersion returns the application version
func GetVersion() string {
	return Version
}