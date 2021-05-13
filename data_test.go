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

func TestParseRegionalType(t *testing.T) {
	testcases := []struct {
		In       string
		RegionID string
		Dir      string
	}{
		{"EUC1-DataTransfer-Out-Bytes", "EUC1", "Out"},
		{"EUC1-DataTransfer-In-Bytes", "EUC1", "In"},
		{"EUC1-DataTransfer-Regional-Bytes", "EUC1", "Regional"},
	}

	for _, v := range testcases {
		regionid, direction, err := ParseRegionalUsageType(v.In)
		if err != nil {
			t.Error(err)
		}
		if regionid != v.RegionID {
			t.Errorf("invalid src have: %s, got: %s", v.RegionID, regionid)
		}
		if direction != v.Dir {
			t.Errorf("invalid dst have: %s, got: %s", v.Dir, direction)
		}
	}
}
