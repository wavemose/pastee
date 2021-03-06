package pastee

import (
	"math"
	"testing"
)

type IntToMBase31Mapping struct {
	Value   int64  // MBase31.Value
	Encoded string // MBase31.ToString()
}

var testCases = []IntToMBase31Mapping{
	// Boundary conditions.
	{0, "a"},
	{1, "b"},
	{30, "9"},
	{31, "ba"},
	{int64(math.Pow(31, 2) - 1), "99"},
	{int64(math.Pow(31, 2)), "baa"},
	{int64(math.Pow(31, 3) - 1), "999"},
	{int64(math.Pow(31, 3)), "baaa"},
	{int64(math.Pow(31, 4) - 1), "9999"},
	{int64(math.Pow(31, 4)), "baaaa"},
	{int64(math.Pow(31, 5) - 1), "99999"},
	{int64(math.Pow(31, 5)), "baaaaa"},

	// Max int64.
	{9223372036854775807, "nzadvkfuvmedh"},
	// Min int64
	{-9223372036854775808, "-nzadvkfuvmedj"},

	// The following test cases were autogenerated from an independent reference
	// implementation.

	{13500, "rbs"},
	{13529, "rcq"},
	{13839, "rpq"},
	{14072, "rw8"},
	{14114, "ryk"},
	{14322, "r7a"},
	{14383, "r89"},
	{14889, "ssk"},
	{-14986, "-svq"},
	{15200, "s4m"},
	{15332, "s8v"},
	{15518, "tev"},
	{15702, "tmt"},
	{15708, "tmz"},
	{15890, "ttv"},
	{-16372, "-ube"},
	{16742, "uqc"},
	{16826, "us3"},
	{16834, "utb"},
	{17072, "u2z"},

	{106648320, "dzr6pp"},
	{119490945, "efn9gw"},
	{127464176, "erawaw"},
	{12971867, "rbqkm"},
	{132629055, "ewv9sw"},
	{135583968, "ez4fmp"},
	{150625090, "fjdb56"},
	{-154833879, "-fpxmrj"},
	{16585848, "u8z8n"},
	{179293615, "gjen5a"},
	{181225698, "gmhhmj"},
	{1865280, "cav9m"},
	{186590882, "gtbmhz"},
	{186668988, "gtd8tj"},
	{191564400, "gyqjv3"},
	{-192477600, "-gzp755"},
	{209389488, "hkzwcw"},
	{213648633, "hrmvcp"},
	{219906560, "hydw9a"},
	{22149155, "29sa6"},
	{221604380, "hz8wyp"},
	{2221318, "cpurq"},

	{1009468124602, "bfyqewahp"},
	{1115341224630, "bktzhk7kw"},
	{11440461637600, "qp4v45cka"},
	{117205506000, "ejb7u59g"},
	{12306128168250, "rqkackp5f"},
	{12409466631048, "rubqwrygd"},
	{-1267981172800, "-bscy52hzq"},
	{12920565455798, "sewmawmya"},
	{1301919064380, "btk8kdzru"},
	{1337968884784, "buwusky58"},
	{14462269860480, "t8xq6j6kb"},
	{1527378050924, "b3s9q7fnh"},
	{15511370814416, "vf3tn4peb"},
	{1607626486935, "b6qpr7hyp"},
	{161250389356, "f5yn7v3d"},
	{1629352101888, "b7g6mx9mg"},
	{1863644319264, "cfz6bgdw8"},
	{186574363061, "g3g8bdny"},
	{1880822391140, "cgnhb4ws3"},
}

func TestToAndFromMBase31(t *testing.T) {
	// int64 -> MBase31 string.
	for i, test := range testCases {
		var expected string = test.Encoded
		var actual string = MBase31{Value: test.Value}.ToString()
		if expected != actual {
			t.Errorf("%d. MBase31{%d} => %s. Expected: %s", i, test.Value, actual, expected)
		}
	}

	// MBase31 string -> int64
	for i, test := range testCases {
		expected := MBase31{Value: test.Value}
		actual, err := MBase31FromString(test.Encoded)
		if err != nil {
			t.Errorf("%d. MBase31FromString(%s) => error: %+v", i, test.Encoded, err)
		} else if expected != actual {
			t.Errorf("%d. MBase31FromString(%s) => %d. Expected: %d", i,
				test.Encoded, actual, expected)
		}
	}
}

func TestParseErrors(t *testing.T) {
	invalidRunes := []rune{'i', 'l', 'o', '0', '1', 'A', 'B'}

	for i, r := range invalidRunes {
		mb31, err := MBase31FromString(string(r))
		if err == nil {
			t.Errorf("%d. MBase31FromString(%s) => %+v. Expected error.", i, string(r), mb31)
		}
	}
}

func TestOverflow(t *testing.T) {
	overflows := []string{
		"nzadvkfuvmedj",  // Max int64 + 1
		"-nzadvkfuvmedk", // Min int64 - 1
	}

	for i, s := range overflows {
		mb31, err := MBase31FromString(s)
		if err == nil {
			t.Errorf("%d. MBaseFromString(%s) => %+v. Expected error.", i, s, mb31)
		} else {
			t.Logf("Got (expected) error: %+v", err)
		}
	}
}
