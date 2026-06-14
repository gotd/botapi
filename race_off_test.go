//go:build !race

package botapi

// raceDetectorEnabled reports whether the binary was built with -race.
const raceDetectorEnabled = false
