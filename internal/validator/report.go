package validator

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintReport(results []*RuleResult, format string) {
	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(results)
		return
	}

	// Text format
	passed := true
	for _, res := range results {
		if res.Status != "PASS" {
			passed = false
		}
		icon := "✅"
		if res.Status == "FAIL" {
			icon = "❌"
		} else if res.Status == "WARN" {
			icon = "⚠️"
		}

		fmt.Printf("%s Rule: %s\n", icon, res.RuleName)
		for _, err := range res.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	if !passed {
		fmt.Println("\nValidation FAILED.")
		os.Exit(1)
	}
	fmt.Println("\nValidation PASSED.")
}
