package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/explodes/greenhouse-pi/stats"
	"github.com/gorilla/mux"
)

func (api *Api) History(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// extract stat type
	// input
	statTypeRaw, ok := vars["stat"]
	if !ok {
		w.WriteHeader(400)
		w.Write([]byte("missing stat"))
		return
	}
	// parse
	statType := stats.StatType(statTypeRaw)
	switch statType {
	case stats.StatTypeTemp:
		break
	case stats.StatTypeHumidity:
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid stat type"))
		return
	}

	// extract start date
	// input
	startRaw, ok := vars["start"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing start time"))
		return
	}
	// parse
	start, err := parseTime(startRaw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid start time"))
		return
	}

	// extract end date
	// input
	endRaw, ok := vars["end"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing end time"))
		return
	}
	// parse
	end, err := parseTime(endRaw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid end time"))
		return
	}

	results, err := api.storage.Fetch(statType, start, end)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error fetching results: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(map[string]interface{}{
		"start": start,
		"end":   end,
		"stat":  statType,
		"items": convertStatsToResponse(results),
	})

}

func parseTime(s string) (time.Time, error) {
	var err error
	var result time.Time
	for _, format := range dateInputFormats {
		result, err = time.Parse(format, s)
		if err == nil {
			return result, nil
		}
	}
	return time.Time{}, err
}

func convertStatsToResponse(stats []stats.Stat) ([]statResponse) {
	results := make([]statResponse, 0, len(stats))
	for _, stat := range stats {
		results = append(results, statResponse{
			When:  stat.When,
			Value: stat.Value,
		})
	}
	return results
}
