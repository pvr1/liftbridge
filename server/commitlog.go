package server

import (
	"github.com/liftbridge-io/liftbridge/server/commitlog"
	"github.com/liftbridge-io/liftbridge/server/proto"
)

// CommitLog is the durable write-ahead log interface used to back each stream.
type CommitLog interface {
	// Delete closes the log and removes all data associated with it from the
	// filesystem.
	Delete() error

	// NewReader creates a new Reader starting at the given offset. If
	// uncommitted is true, the Reader will read uncommitted messages from the
	// log. Otherwise, it will only return committed messages.
	NewReader(offset int64, uncommitted bool) (*commitlog.Reader, error)

	// Truncate removes all messages from the log starting at the given offset.
	Truncate(offset int64) error

	// NewestOffset returns the offset of the last message in the log.
	NewestOffset() int64

	// OldestOffset returns the offset of the first message in the log.
	OldestOffset() int64

	// OffsetForTimestamp returns the earliest offset whose timestamp is
	// greater than or equal to the given timestamp.
	OffsetForTimestamp(timestamp int64) (int64, error)

	// SetHighWatermark sets the high watermark on the log. All messages up to
	// and including the high watermark are considered committed.
	SetHighWatermark(hw int64)

	// HighWatermark returns the high watermark for the log.
	HighWatermark() int64

	// Append writes the given batch of messages to the log and returns their
	// corresponding offsets in the log.
	Append(msg []*proto.Message) ([]int64, error)

	// AppendMessageSet writes the given message set data to the log and
	// returns the corresponding offsets in the log.
	AppendMessageSet(ms []byte) ([]int64, error)

	// Clean applies retention and compaction rules against the log, if
	// applicable.
	Clean() error

	// Close closes each log segment file and stops the background goroutine
	// checkpointing the high watermark to disk.
	Close() error
}
