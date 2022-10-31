package mpegts

// Timestamp is a 33-bit MPEG-2 timestamp for PTS/DTS
type Timestamp uint64

const (
	NonTimestamp Timestamp = 1 << 33
	MaxTimestamp Timestamp = NonTimestamp - 1
)

const (
	SystemClock = 90000 // 90kHz
)

// Scale returns the timestamp from the duration and timescale.
// For example, need to get number of 90kHz ticks in 250 milliseconds.
// Duration is 250 and the timescale is 1000 - number of milliseconds in second.
// `Scale(250, 1000)` return 22500 ticks
func Scale(duration, timescale int) Timestamp {
	return Timestamp(duration) * SystemClock / Timestamp(timescale)
}

// Delta returns the difference t-u considering value overflow
func (t Timestamp) Delta(u Timestamp) Timestamp {
	if t >= u {
		return t - u
	} else {
		return NonTimestamp - u + t
	}
}

// Add returns the timestamp t+u
func (t Timestamp) Add(u Timestamp) Timestamp {
	return (t + u) & MaxTimestamp
}
