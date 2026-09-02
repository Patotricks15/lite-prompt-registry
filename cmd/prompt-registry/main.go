package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-guardrail/open-prompt-registry/internal/registry"
)

type api struct {
	service *registry.Service
}

func main() {
	app := &api{service: registry.NewService()}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", app.health)
	mux.HandleFunc("POST /prompts", app.createPrompt)
	mux.HandleFunc("GET /prompts/{id}", app.getPrompt)
	mux.HandleFunc("POST /prompts/{id}/versions", app.createVersion)
	mux.HandleFunc("POST /prompts/{id}/versions/{number}/tests", app.addTest)
	mux.HandleFunc("POST /prompts/{id}/versions/{number}/review", app.requestReview)
	mux.HandleFunc("POST /prompts/{id}/versions/{number}/approve", app.approve)
	mux.HandleFunc("POST /prompts/{id}/versions/{number}/rollout", app.rollOut)

	log.Println("open-prompt-registry listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (a *api) health(w http.ResponseWriter, _ *http.Request) { writeJSON(w, http.StatusOK, map[string]string{"status": "ok"}) }

func (a *api) createPrompt(w http.ResponseWriter, r *http.Request) {
	var input struct { Name string `json:"name"`; Description string `json:"description"` }
	if !decode(w, r, &input) { return }
	if strings.TrimSpace(input.Name) == "" { writeError(w, http.StatusBadRequest, "name is required"); return }
	writeJSON(w, http.StatusCreated, a.service.CreatePrompt(input.Name, input.Description))
}

func (a *api) getPrompt(w http.ResponseWriter, r *http.Request) {
	prompt, err := a.service.GetPrompt(r.PathValue("id"))
	if err != nil { writeServiceError(w, err); return }
	writeJSON(w, http.StatusOK, prompt)
}

func (a *api) createVersion(w http.ResponseWriter, r *http.Request) {
	var input struct { Template string `json:"template"`; Author string `json:"author"` }
	if !decode(w, r, &input) { return }
	if strings.TrimSpace(input.Template) == "" || strings.TrimSpace(input.Author) == "" { writeError(w, http.StatusBadRequest, "template and author are required"); return }
	version, err := a.service.CreateVersion(r.PathValue("id"), input.Template, input.Author)
	if err != nil { writeServiceError(w, err); return }
	writeJSON(w, http.StatusCreated, version)
}

func (a *api) addTest(w http.ResponseWriter, r *http.Request) {
	var input struct { Name string `json:"name"`; Passed bool `json:"passed"`; Details string `json:"details"` }
	if !decode(w, r, &input) { return }
	if strings.TrimSpace(input.Name) == "" { writeError(w, http.StatusBadRequest, "name is required"); return }
	number, ok := versionNumber(w, r)
	if !ok { return }
	version, err := a.service.AddTest(r.PathValue("id"), number, registry.TestResult{Name: input.Name, Passed: input.Passed, Details: input.Details})
	if err != nil { writeServiceError(w, err); return }
	writeJSON(w, http.StatusCreated, version)
}

func (a *api) requestReview(w http.ResponseWriter, r *http.Request) { a.transition(w, r, "review") }
func (a *api) approve(w http.ResponseWriter, r *http.Request) { a.transition(w, r, "approve") }
func (a *api) rollOut(w http.ResponseWriter, r *http.Request) { a.transition(w, r, "rollout") }

func (a *api) transition(w http.ResponseWriter, r *http.Request, action string) {
	var input struct { Actor string `json:"actor"` }
	if action != "rollout" {
		if !decode(w, r, &input) { return }
		if strings.TrimSpace(input.Actor) == "" { writeError(w, http.StatusBadRequest, "actor is required"); return }
	}
	number, ok := versionNumber(w, r)
	if !ok { return }
	var version *registry.Version
	var err error
	switch action {
	case "review": version, err = a.service.RequestReview(r.PathValue("id"), number, input.Actor)
	case "approve": version, err = a.service.Approve(r.PathValue("id"), number, input.Actor)
	case "rollout": version, err = a.service.RollOut(r.PathValue("id"), number)
	}
	if err != nil { writeServiceError(w, err); return }
	writeJSON(w, http.StatusOK, version)
}

func versionNumber(w http.ResponseWriter, r *http.Request) (int, bool) {
	number, err := strconv.Atoi(r.PathValue("number"))
	if err != nil || number < 1 { writeError(w, http.StatusBadRequest, "version number must be a positive integer"); return 0, false }
	return number, true
}

func decode(w http.ResponseWriter, r *http.Request, destination any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(destination); err != nil { writeError(w, http.StatusBadRequest, "invalid JSON request"); return false }
	return true
}
func writeJSON(w http.ResponseWriter, status int, value any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(value) }
func writeError(w http.ResponseWriter, status int, message string) { writeJSON(w, status, map[string]string{"error": message}) }
func writeServiceError(w http.ResponseWriter, err error) { if errors.Is(err, registry.ErrPromptNotFound) || errors.Is(err, registry.ErrVersionNotFound) { writeError(w, http.StatusNotFound, err.Error()); return }; if errors.Is(err, registry.ErrInvalidState) || errors.Is(err, registry.ErrTestsRequired) { writeError(w, http.StatusConflict, err.Error()); return }; writeError(w, http.StatusInternalServerError, "internal server error") }
