package analizer

import (
	"strings"
	"testing"
)

func TestEscaapeTargets(t *testing.T) {

	targets := []string{
		"@test Hello, world.",
		"https://lgb.com/analizer Hello, world.",
		"Hello, world. @test",
		"Hello, world. https://lgb.com/analizer",
		"@test Hello, world https://lgb.com/analizer",
		"@test Hello, http://lgb.com/analizer world.",
		"http://lgb.com/analizer Hello, world. @test",
		"http://lgb.com/analizer @test @test https://lgb.com/analizer",
	}

	for _, v := range targets {
		a := NewAnalizer([]string{v})
		a.EscapeTargets()

		if strings.Contains(a.GetTarget(0), "@test") {
			t.Errorf("Target %v contains @test\n", v)
		}

		if strings.Contains(a.GetTarget(0), "http://lgb.com/analizer") {
			t.Errorf("Target %v contains http://lgb.com/analizer\n", v)
		}

		if strings.Contains(a.GetTarget(0), "https://lgb.com/analizer") {
			t.Errorf("Target %v contains https://lgb.com/analizer\n", v)
		}

	}

}
