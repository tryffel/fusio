package pipeline

import (
	"testing"
)

func TestSubscriptionStore(t *testing.T) {
	store := NewSubscriptionStore()

	err := store.subscribePipeline("g1", "d1", "m1", "p1")
	err = store.subscribePipeline("g2", "d1", "m1", "p1")
	err = store.subscribePipeline("g2", "d2", "m1", "p1")
	err = store.subscribePipeline("g2", "d2", "m2", "p1")

	if err != nil {
		t.Error(err)
	}

	pipelines, err := store.getSubscriptedPipelines(&[]string{"g1"}, "d1", "m1")
	if len(*pipelines) == 0 {
		t.Error("didn't get any pipelines, expected 1")
	}

	if (*pipelines)[0] != "p1" {
		t.Error("Received invalid pipeline id")
	}
}

func TestSubscriptionStoreMultipleGroups(t *testing.T) {
	store := NewSubscriptionStore()

	err := store.subscribePipeline("g1", "d1", "m1", "p1")
	_ = store.subscribePipeline("g1", "d1", "m1", "p2")
	_ = store.subscribePipeline("g2", "", "m1", "p3")
	_ = store.subscribePipeline("g2", "", "m2", "p4")
	_ = store.subscribePipeline("g2", "d2", "m1", "p5")

	if err != nil {
		t.Error(err)
	}

	pipelines, err := store.getSubscriptedPipelines(&[]string{"g1"}, "d1", "m1")
	if len(*pipelines) == 0 {
		t.Error("didn't get any pipelines, expected 1")
	}

	if (*pipelines)[0] != "p1" {
		t.Error("Received invalid pipeline id")
	}

	pipelines, err = store.getSubscriptedPipelines(&[]string{"g1", "g2"}, "d1", "m1")
	checkIdsMatches([]string{"p1", "p2", "p3"}, *pipelines, t, "checking multiple groups sub")

	pipelines, err = store.getSubscriptedPipelines(&[]string{"g1", "g2"}, "d2", "m1")
	checkIdsMatches([]string{"p3", "p5"}, *pipelines, t, "checking multiple groups sub")

	pipelines, err = store.getSubscriptedPipelines(&[]string{"g1", "g2"}, "d2", "m2")
	checkIdsMatches([]string{"p4"}, *pipelines, t, "checking multiple groups sub")
}

func BenchmarkGetSubsctions10Subs(b *testing.B) {
	store := NewSubscriptionStore()

	_ = store.subscribePipeline("g1", "d1", "m1", "p1")
	_ = store.subscribePipeline("g2", "d1", "m1", "p2")
	_ = store.subscribePipeline("g2", "d2", "m1", "p3")
	_ = store.subscribePipeline("g2", "d2", "m2", "p4")
	_ = store.subscribePipeline("g2", "", "m2", "p5")

	_ = store.subscribePipeline("g1", "", "m1", "p10")
	_ = store.subscribePipeline("g2", "", "m1", "p20")
	_ = store.subscribePipeline("g3", "", "m1", "p30")
	_ = store.subscribePipeline("g4", "", "m2", "p40")
	_ = store.subscribePipeline("g5", "", "m2", "p50")

	for i := 0; i < b.N; i++ {
		_, _ = store.getSubscriptedPipelines(&[]string{"g1"}, "", "m1")
	}
}

func checkIdsMatches(target []string, actual []string, t *testing.T, testCase string) {
	if len(target) != len(actual) {
		t.Errorf("%s, pipeline id array length doesn't match, got %d, expected %d", testCase, len(actual), len(target))
		return
	}

	for i := range target {
		if target[i] != actual[i] {
			t.Errorf("%s, pipeline id [%d] doesn't match, got %s, expected %s", testCase, i, actual[i], target[i])
		}
	}
}
