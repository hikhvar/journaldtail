package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/cortexproject/cortex/pkg/util/flagext"

	kitlog "github.com/go-kit/kit/log"

	"github.com/hikhvar/journaldtail/pkg/storage"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/grafana/loki/pkg/promtail"
	"github.com/hikhvar/journaldtail/pkg/journald"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
)

var lokiHostURL = "http://localhost:3100/api/prom/push"

func main() {
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	log.SetOutput(kitlog.NewStdlibAdapter(logger))
	// TODO: Store state on disk
	memStorage := storage.Memory{}
	journal, err := sdjournal.NewJournal()
	if err != nil {
		log.Fatal(fmt.Sprintf("could not open journal: %s", err.Error()))
	}
	reader := journald.NewReader(journal, &memStorage)

	// TODO: Read from CLI
	if v, isSet := os.LookupEnv("LOKI_URL"); isSet {
		lokiHostURL = v
	}

	cfg := promtail.ClientConfig{
		URL: flagext.URLValue{
			URL: MustParseURL(lokiHostURL),
		},
	}
	lokiClient, err := promtail.NewClient(cfg, logger)
	if err != nil {
		log.Fatal(fmt.Sprintf("could not create loki client: %s", err.Error()))
	}
	err = TailLoop(reader, lokiClient)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to tail journald: %s", err.Error()))
	}
}

func TailLoop(reader *journald.Reader, writer *promtail.Client) error {
	var lastTS time.Time
	for {
		r, err := reader.Next()
		if err != nil {
			return errors.Wrap(err, "could not get next journal entry")
		}
		if r != nil {
			ls := ToLabelSet(r)
			ts := journald.ToGolangTime(r.RealtimeTimestamp)
			msg := r.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE]

			if ts.Before(lastTS) {
				log.Fatal(fmt.Sprintf("%s is before %s! Message: %s", ts, lastTS, msg))
			}
			lastTS = ts
			err = writer.Handle(ls, ts, msg)
			if err != nil {
				return errors.Wrap(err, "could not enque systemd logentry")
			}
		}

	}
}

func ToLabelSet(reader *sdjournal.JournalEntry) model.LabelSet {
	ret := make(model.LabelSet)
	for key, value := range reader.Fields {
		if key != sdjournal.SD_JOURNAL_FIELD_MESSAGE {
			ret[model.LabelName(key)] = model.LabelValue(value)
		}
	}
	return ret
}

func MustParseURL(input string) *url.URL {
	u, err := url.Parse(input)
	if err != nil {
		panic(fmt.Sprintf("could not parse static url: %s", input))
	}
	return u
}
