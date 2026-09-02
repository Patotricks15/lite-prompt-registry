package registry

import (
	"errors"
	"testing"
)

func TestPromptVersionLifecycle(t *testing.T) {
	service := NewService()
	prompt := service.CreatePrompt("support-reply", "Customer support response prompt")
	version, err := service.CreateVersion(prompt.ID, "Answer: {{question}}", "ana")
	if err != nil { t.Fatalf("CreateVersion() error = %v", err) }

	if _, err := service.RequestReview(prompt.ID, version.Number, "bruno"); err != nil { t.Fatalf("RequestReview() error = %v", err) }
	if _, err := service.Approve(prompt.ID, version.Number, "carla"); !errors.Is(err, ErrTestsRequired) { t.Fatalf("Approve() error = %v, want ErrTestsRequired", err) }
	if _, err := service.AddTest(prompt.ID, version.Number, TestResult{Name: "blocks-injection", Passed: true}); err != nil { t.Fatalf("AddTest() error = %v", err) }
	if _, err := service.Approve(prompt.ID, version.Number, "carla"); err != nil { t.Fatalf("Approve() error = %v", err) }
	if _, err := service.RollOut(prompt.ID, version.Number); err != nil { t.Fatalf("RollOut() error = %v", err) }

	stored, err := service.GetPrompt(prompt.ID)
	if err != nil { t.Fatalf("GetPrompt() error = %v", err) }
	if stored.ProductionVersionID != 1 || stored.Versions[0].Status != StatusRolledOut { t.Fatalf("unexpected rollout state: %+v", stored) }
}
