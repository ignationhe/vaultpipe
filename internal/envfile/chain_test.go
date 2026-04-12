package envfile

import (
	"errors"
	"testing"
)

func addKeyStep(k, v string) ChainStep {
	return ChainStep{
		Name: "add-" + k,
		Fn: func(m map[string]string) (map[string]string, error) {
			out := copyMap(m)
			out[k] = v
			return out, nil
		},
	}
}

func failStep(name string) ChainStep {
	return ChainStep{
		Name: name,
		Fn: func(m map[string]string) (map[string]string, error) {
			return nil, errors.New("step failed")
		},
	}
}

func TestChain_AppliesStepsInOrder(t *testing.T) {
	initial := map[string]string{"A": "1"}
	steps := []ChainStep{addKeyStep("B", "2"), addKeyStep("C", "3")}

	out, _, err := Chain(initial, steps, DefaultChainOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestChain_StopsOnErrorWhenConfigured(t *testing.T) {
	initial := map[string]string{}
	steps := []ChainStep{addKeyStep("X", "x"), failStep("boom"), addKeyStep("Y", "y")}

	out, results, err := Chain(initial, steps, DefaultChainOptions())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := out["Y"]; ok {
		t.Error("Y should not be set after failure")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result entry, got %d", len(results))
	}
}

func TestChain_ContinuesOnErrorWhenNotStopping(t *testing.T) {
	initial := map[string]string{}
	opts := ChainOptions{StopOnError: false, Trace: false}
	steps := []ChainStep{addKeyStep("X", "x"), failStep("boom"), addKeyStep("Y", "y")}

	out, _, err := Chain(initial, steps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "x" || out["Y"] != "y" {
		t.Errorf("expected X and Y in output, got %v", out)
	}
}

func TestChain_TraceRecordsIntermediateResults(t *testing.T) {
	initial := map[string]string{}
	opts := ChainOptions{StopOnError: true, Trace: true}
	steps := []ChainStep{addKeyStep("A", "1"), addKeyStep("B", "2")}

	_, results, err := Chain(initial, steps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 trace results, got %d", len(results))
	}
	if results[0].Output["A"] != "1" {
		t.Errorf("trace step 0 missing A=1")
	}
	if results[1].Output["B"] != "2" {
		t.Errorf("trace step 1 missing B=2")
	}
}

func TestChain_NilFnReturnsError(t *testing.T) {
	initial := map[string]string{}
	steps := []ChainStep{{Name: "nil-fn", Fn: nil}}

	_, _, err := Chain(initial, steps, DefaultChainOptions())
	if err == nil {
		t.Fatal("expected error for nil Fn")
	}
}

func TestChain_EmptyStepsReturnsInitial(t *testing.T) {
	initial := map[string]string{"K": "V"}
	out, _, err := Chain(initial, nil, DefaultChainOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "V" {
		t.Errorf("expected initial map unchanged")
	}
}
