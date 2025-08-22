package formatter

import (
	"strings"

	"github.com/qlik-trial/usage-telemetry-publisher/internal/model"
)

func Write(events []*model.ScrubbedEvent) string {
	res := make([]string, 0, len(events))
	for _, e := range events {
		res = append(res, Flatten(e))
	}
	return strings.Join(res, "\n")
}
