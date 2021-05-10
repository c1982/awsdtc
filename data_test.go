package main

import "testing"

func TestParseUsageType(t *testing.T) {

	testcases := []struct {
		In  string
		Src string
		Dst string
		Dir string
	}{
		{"USW2-EU-AWS-Out-Bytes", "USW2", "EU", "Out"},
		{"USW2-APS2-AWS-In-Bytes", "USW2", "APS2", "In"},
	}

	for _, v := range testcases {
		src, dst, dir, err := ParseUsageType(v.In)
		if err != nil {
			t.Error(err)
		}

		if src != v.Src {
			t.Errorf("invalid src have: %s, got: %s", v.Src, src)
		}

		if dst != v.Dst {
			t.Errorf("invalid dst have: %s, got: %s", v.Dst, dst)
		}

		if dir != v.Dir {
			t.Errorf("invalid direction have: %s, got: %s", v.Dir, dir)
		}

	}

}
