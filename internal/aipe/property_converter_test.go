package aipe

import "testing"

// AIPE can return non-string, non-bool property values (e.g. a nested
// object like "external-system-reference": {} on manually created
// objects). convertPropertiesToString must not panic on those; it should
// skip them instead, while still converting the properties it can handle.
func TestConvertPropertiesToString_SkipsNonStringNonBoolValues(t *testing.T) {
	input := map[string]interface{}{
		"name":                       "User Management",
		"external-system-reference":  map[string]interface{}{},
		"enabled":                    true,
	}

	result := convertPropertiesToString(input)

	if got, want := result["name"], "User Management"; got != want {
		t.Errorf("name = %q, want %q", got, want)
	}

	if got, want := result["enabled"], "true"; got != want {
		t.Errorf("enabled = %q, want %q", got, want)
	}

	if _, present := result["external-system-reference"]; present {
		t.Errorf("external-system-reference should have been skipped, got present: %v", result["external-system-reference"])
	}
}
