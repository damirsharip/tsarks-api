package substr

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type SubstrRequest struct {
	Word string `json:"word"`
}

func (r *SubstrRequest) Validate() error {
	if r.Word == "" {
		return fmt.Errorf("word required")
	}
	return nil
}

type SubstrResponse struct {
	Answer string `json:"answer"`
}

func FindSubstr() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SubstrRequest

		if err := decodeJSONBody(w, r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request payload")
			return
		}

		err := req.Validate()
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("bad request: %s", err))
			return
		}

		ans := findLongestSubstring(req.Word)

		respondJSON(w, http.StatusOK, SubstrResponse{Answer: ans})
	})
}

func findLongestSubstring(s string) string {
	n := len(s)
	charIndexMap := make(map[uint8]int)
	var result int
	var start int

	l := 0
	r := n

	for end := 0; end < n; end++ {
		duplicateIndex, isDuplicate := charIndexMap[s[end]]
		if isDuplicate {
			if end-start >= result {
				l = start
				r = end
				result = end - start
			}

			for i := start; i <= duplicateIndex; i++ {
				delete(charIndexMap, s[i])
			}

			start = duplicateIndex + 1
		}

		charIndexMap[s[end]] = end
	}

	result = int(math.Max(float64(result), float64(n-start)))
	if n-start >= result {
		l = start
		r = n
		result = n - start
	}

	return s[l:r]
}

func decodeJSONBody(_ http.ResponseWriter, r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
