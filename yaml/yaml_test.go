package yaml

import (
	"testing"
)

func TestParse(t *testing.T) {
	got := Parse("test.yaml")
	if got.Cfgs[0].Name != "test" {
		t.Errorf("Parse(\"mts.yaml\") = %v; want something ELSE!", got)
	}
}
