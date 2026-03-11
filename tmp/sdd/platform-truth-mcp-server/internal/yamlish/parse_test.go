package yamlish

import "testing"

func TestParseNestedMapsAndLists(t *testing.T) {
	input := []byte(`platform:
  version: "2026.03"
platform_refs:
  principles:
    - id: "principles.api-versioning"
      reason: "Versioned APIs"
  contracts:
    - id: "contracts.customer-profile.v2"
      reason: "Shared profile contract"
change:
  requires_platform_change: false
`)

	value, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if got := LookupString(value, "platform", "version"); got != "2026.03" {
		t.Fatalf("platform.version = %q, want 2026.03", got)
	}
	if LookupBool(value, "change", "requires_platform_change") {
		t.Fatalf("requires_platform_change = true, want false")
	}

	principles := AsSlice(Lookup(value, "platform_refs", "principles"))
	if len(principles) != 1 {
		t.Fatalf("principles length = %d, want 1", len(principles))
	}
	first := AsMap(principles[0])
	if first["id"] != "principles.api-versioning" {
		t.Fatalf("principle id = %v", first["id"])
	}
}
