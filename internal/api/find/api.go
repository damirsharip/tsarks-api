package find

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type FindEmailRequest struct {
	Text string `json:"text"`
}

func (r *FindEmailRequest) Validate() error {
	if r.Text == "" {
		return fmt.Errorf("text required")
	}
	return nil
}

type FindEmailResponse struct {
	Emails []string `json:"emails"`
}

func FindEmail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req FindEmailRequest

		if err := decodeJSONBody(w, r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request payload")
			return
		}

		err := req.Validate()
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("bad request: %s", err))
			return
		}

		// Ищем все email адреса в ответе с помощью регулярного выражения
		emailRegex := regexp.MustCompile(`Email:\s*([^\s]+)`)
		matches := emailRegex.FindAllStringSubmatch(req.Text, -1)

		var ans []string
		// Выводим найденные email адреса
		for _, match := range matches {
			ans = append(ans, match[1])
		}

		respondJSON(w, http.StatusOK, FindEmailResponse{Emails: ans})
	})
}

type FindIINRequest struct {
	Text string `json:"text"`
}

func (r *FindIINRequest) Validate() error {
	if r.Text == "" {
		return fmt.Errorf("text required")
	}
	return nil
}

type FindIINResponse struct {
	IINs []string `json:"iins"`
}

func FindIIN() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req FindEmailRequest

		if err := decodeJSONBody(w, r, &req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request payload")
			return
		}

		err := req.Validate()
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("bad request: %s", err))
			return
		}

		iinRegex := regexp.MustCompile(`IIN:\s*(\d{12})`)
		matches := iinRegex.FindAllStringSubmatch(req.Text, -1)

		var ans []string
		for _, match := range matches {
			ans = append(ans, match[1])
		}

		respondJSON(w, http.StatusOK, FindIINResponse{IINs: ans})
	})
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
