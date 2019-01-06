package journald

import (
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/pkg/errors"
)

type Storage interface {
	Save(position string) error
	Last() (string, error)
}

type Reader struct {
	src     *sdjournal.Journal
	storage Storage
}

func NewReader(src *sdjournal.Journal, storage Storage) *Reader {
	return &Reader{
		src:     src,
		storage: storage,
	}
}

// Seek seeks to last position. Either it is the position saved in storage or the tail of the journal
func (r *Reader) Seek() error {
	last, err := r.storage.Last()
	if err != nil {
		errjd := r.src.SeekTail()
		if errjd != nil {
			return errors.Wrapf(errjd, "could not seek to journal tail after storage error: %s", err.Error())
		}
	}
	err = r.src.SeekCursor(last)
	return errors.Wrapf(err, "could not seek to cursor %s", last)
}

// Next blocks until the next journal entry is available
func (r *Reader) Next() (*sdjournal.JournalEntry, error) {
	advanced, err := r.advance()
	if err != nil {
		return nil, errors.Wrap(err, "failed to advance")
	}
	if !advanced {
		r.src.Wait(sdjournal.IndefiniteWait)
		advanced, err = r.advance()
		if advanced != true {
			return nil, errors.New("finished wait but could not advance")
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to advance after wait")
		}
	}
	entry, err := r.src.GetEntry()
	if err == nil && entry != nil {
		err = r.storage.Save(entry.Cursor)
		return entry, errors.Wrap(err, "could not save cursor position to storage")
	}
	return entry, errors.Wrap(err, "could not get next entry")
}

// advance tries to jump to next journal entry.
func (r *Reader) advance() (bool, error) {
	rc, err := r.src.Next()
	if err != nil {
		return false, errors.Wrap(err, "could not get next. ")
	}
	return rc != 0, nil
}

func ToGolangTime(sdTime uint64) time.Time {
	seconds := sdTime / 1000000
	reminderMicroseconds := sdTime % 1000000
	return time.Unix(int64(seconds), int64(reminderMicroseconds*1000))
}
