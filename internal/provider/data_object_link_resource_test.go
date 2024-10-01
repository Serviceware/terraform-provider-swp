package provider

import (
	"fmt"
	"testing"
)

var idDiffTests = []struct {
	stateTargetIDs []string
	planTargetIDs  []string

	expectedAdd    []string
	expectedRemove []string
}{
	{
		stateTargetIDs: []string{"a", "b", "c"},
		planTargetIDs:  []string{"b", "c", "d"},
		expectedAdd:    []string{"d"},
		expectedRemove: []string{"a"},
	},
	{
		stateTargetIDs: []string{"a", "b", "c"},
		planTargetIDs:  []string{"a", "b", "c"},
		expectedAdd:    []string{},
		expectedRemove: []string{},
	},
	{
		stateTargetIDs: []string{"a", "b"},
		planTargetIDs:  []string{"a", "b", "c", "d"},
		expectedAdd:    []string{"c", "d"},
		expectedRemove: []string{},
	},
	{
		stateTargetIDs: []string{"a", "b", "c", "d"},
		planTargetIDs:  []string{"a"},
		expectedAdd:    []string{},
		expectedRemove: []string{"b", "c", "d"},
	},
	{
		stateTargetIDs: []string{},
		planTargetIDs:  []string{"a", "b"},
		expectedAdd:    []string{"a", "b"},
		expectedRemove: []string{},
	},
	{
		stateTargetIDs: []string{"a", "b"},
		planTargetIDs:  []string{},
		expectedAdd:    []string{},
		expectedRemove: []string{"a", "b"},
	},
	{
		stateTargetIDs: []string{"a", "b", "c"},
		planTargetIDs:  []string{"d", "e", "f"},
		expectedAdd:    []string{"d", "e", "f"},
		expectedRemove: []string{"a", "b", "c"},
	},
}

func TestDiffStateAndPlanIDs(t *testing.T) {
	for _, tt := range idDiffTests {
		t.Run(fmt.Sprintf("diff(%v, %v)", tt.stateTargetIDs, tt.planTargetIDs), func(t *testing.T) {
			add, remove := diffStateAndPlanIDs(tt.stateTargetIDs, tt.planTargetIDs)

			if len(add) != len(tt.expectedAdd) {
				t.Errorf("expected %d add ids, got %d", len(tt.expectedAdd), len(add))
			}

			if len(remove) != len(tt.expectedRemove) {
				t.Errorf("expected %d remove ids, got %d", len(tt.expectedRemove), len(remove))
			}

			var addAsMap = make(map[string]bool)
			for _, id := range add {
				addAsMap[id] = true
			}

			var removeAsMap = make(map[string]bool)
			for _, id := range remove {
				removeAsMap[id] = true
			}

			for _, id := range tt.expectedAdd {
				if !addAsMap[id] {
					t.Errorf("expected %s to be in add", id)
				}
			}

			for _, id := range tt.expectedRemove {
				if !removeAsMap[id] {
					t.Errorf("expected %s to be in remove", id)
				}
			}
		})
	}
}
