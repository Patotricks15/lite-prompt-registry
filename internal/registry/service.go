package registry

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type VersionStatus string

const (
	StatusDraft    VersionStatus = "draft"
	StatusReview   VersionStatus = "in_review"
	StatusApproved VersionStatus = "approved"
	StatusRolledOut VersionStatus = "rolled_out"
)

var (
	ErrPromptNotFound  = errors.New("prompt not found")
	ErrVersionNotFound = errors.New("version not found")
	ErrInvalidState    = errors.New("invalid workflow state")
	ErrTestsRequired   = errors.New("at least one passing test is required")
)

type TestResult struct {
	Name      string    `json:"name"`
	Passed    bool      `json:"passed"`
	Details   string    `json:"details,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Version struct {
	Number    int           `json:"number"`
	Template  string        `json:"template"`
	Status    VersionStatus `json:"status"`
	Author    string        `json:"author"`
	Reviewer  string        `json:"reviewer,omitempty"`
	Approver  string        `json:"approver,omitempty"`
	Tests     []TestResult  `json:"tests"`
	CreatedAt time.Time     `json:"created_at"`
}

type Prompt struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description,omitempty"`
	Versions            []*Version `json:"versions"`
	ProductionVersionID int        `json:"production_version,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

type Service struct {
	mu      sync.RWMutex
	prompts map[string]*Prompt
	nextID  int
}

func NewService() *Service {
	return &Service{prompts: make(map[string]*Prompt)}
}

func (s *Service) CreatePrompt(name, description string) *Prompt {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	prompt := &Prompt{
		ID:          fmt.Sprintf("prompt-%d", s.nextID),
		Name:        name,
		Description: description,
		Versions:    make([]*Version, 0),
		CreatedAt:   time.Now().UTC(),
	}
	s.prompts[prompt.ID] = prompt
	return clonePrompt(prompt)
}

func (s *Service) GetPrompt(id string) (*Prompt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	prompt, ok := s.prompts[id]
	if !ok {
		return nil, ErrPromptNotFound
	}
	return clonePrompt(prompt), nil
}

func (s *Service) CreateVersion(promptID, template, author string) (*Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	prompt, ok := s.prompts[promptID]
	if !ok {
		return nil, ErrPromptNotFound
	}
	version := &Version{Number: len(prompt.Versions) + 1, Template: template, Status: StatusDraft, Author: author, Tests: make([]TestResult, 0), CreatedAt: time.Now().UTC()}
	prompt.Versions = append(prompt.Versions, version)
	return cloneVersion(version), nil
}

func (s *Service) AddTest(promptID string, number int, result TestResult) (*Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	version, err := s.version(promptID, number)
	if err != nil {
		return nil, err
	}
	if version.Status != StatusDraft && version.Status != StatusReview {
		return nil, ErrInvalidState
	}
	result.CreatedAt = time.Now().UTC()
	version.Tests = append(version.Tests, result)
	return cloneVersion(version), nil
}

func (s *Service) RequestReview(promptID string, number int, reviewer string) (*Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	version, err := s.version(promptID, number)
	if err != nil {
		return nil, err
	}
	if version.Status != StatusDraft {
		return nil, ErrInvalidState
	}
	version.Status, version.Reviewer = StatusReview, reviewer
	return cloneVersion(version), nil
}

func (s *Service) Approve(promptID string, number int, approver string) (*Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	version, err := s.version(promptID, number)
	if err != nil {
		return nil, err
	}
	if version.Status != StatusReview {
		return nil, ErrInvalidState
	}
	if !allTestsPass(version.Tests) {
		return nil, ErrTestsRequired
	}
	version.Status, version.Approver = StatusApproved, approver
	return cloneVersion(version), nil
}

func (s *Service) RollOut(promptID string, number int) (*Version, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	prompt, ok := s.prompts[promptID]
	if !ok {
		return nil, ErrPromptNotFound
	}
	if number < 1 || number > len(prompt.Versions) {
		return nil, ErrVersionNotFound
	}
	version := prompt.Versions[number-1]
	if version.Status != StatusApproved {
		return nil, ErrInvalidState
	}
	version.Status = StatusRolledOut
	prompt.ProductionVersionID = number
	return cloneVersion(version), nil
}

func (s *Service) version(promptID string, number int) (*Version, error) {
	prompt, ok := s.prompts[promptID]
	if !ok {
		return nil, ErrPromptNotFound
	}
	if number < 1 || number > len(prompt.Versions) {
		return nil, ErrVersionNotFound
	}
	return prompt.Versions[number-1], nil
}

func allTestsPass(results []TestResult) bool {
	if len(results) == 0 {
		return false
	}
	for _, result := range results {
		if !result.Passed {
			return false
		}
	}
	return true
}

func cloneVersion(version *Version) *Version {
	copy := *version
	copy.Tests = append([]TestResult(nil), version.Tests...)
	return &copy
}

func clonePrompt(prompt *Prompt) *Prompt {
	copy := *prompt
	copy.Versions = make([]*Version, len(prompt.Versions))
	for index, version := range prompt.Versions {
		copy.Versions[index] = cloneVersion(version)
	}
	return &copy
}
