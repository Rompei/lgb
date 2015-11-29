package analyzer

import (
	"strings"
	"testing"
)

func TestEscaapeTargets(t *testing.T) {

	targets := []string{
		"@test Hello, world.",
		"https://lgb.com/analyzer Hello, world.",
		"Hello, world. @test",
		"Hello, world. https://lgb.com/analyzer",
		"@test Hello, world https://lgb.com/analyzer",
		"@test Hello, http://lgb.com/analyzer world.",
		"http://lgb.com/analyzer Hello, world. @test",
		"http://lgb.com/analyzer @test @test https://lgb.com/analyzer",
	}

	for _, v := range targets {
		a := NewAnalyzer([]string{v})
		a.EscapeTargets()

		if strings.Contains(a.GetTarget(0), "@test") {
			t.Errorf("Target %v contains @test\n", v)
		}

		if strings.Contains(a.GetTarget(0), "http://lgb.com/analyzer") {
			t.Errorf("Target %v contains http://lgb.com/analyzer\n", v)
		}

		if strings.Contains(a.GetTarget(0), "https://lgb.com/analyzer") {
			t.Errorf("Target %v contains https://lgb.com/analyzer\n", v)
		}

	}

}
