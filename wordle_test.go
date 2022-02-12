package wordle

import "testing"

func TestPatternParse(t *testing.T) {
	type stage struct {
		pattern  string
		negative string
		expect   string
		err      error
	}
	var stages = []stage{
		{
			pattern: "*****",
			expect:  "[ABCDEFGHIJKLMNOPQRSTUVWXYZ][ABCDEFGHIJKLMNOPQRSTUVWXYZ][ABCDEFGHIJKLMNOPQRSTUVWXYZ][ABCDEFGHIJKLMNOPQRSTUVWXYZ][ABCDEFGHIJKLMNOPQRSTUVWXYZ]",
		},
		{
			pattern: "H*^E^R*",
			negative: "UIOBNM",
			expect:  "H[ACDEFGHJKLPQRSTVWXYZ][ACDFGHJKLPQRSTVWXYZ][ACDEFGHJKLPQSTVWXYZ][ACDEFGHJKLPQRSTVWXYZ]",
		},
		{
			pattern: "*^R^U**",
			negative: "BSH",
			expect:  "[ACDEFGIJKLMNOPQRTUVWXYZ][ACDEFGIJKLMNOPQTUVWXYZ][ACDEFGIJKLMNOPQRTVWXYZ][ACDEFGIJKLMNOPQRTUVWXYZ][ACDEFGIJKLMNOPQRTUVWXYZ]",
		},
	}
	conv := func(rules [5]rule) (r string) {
		for i := 0; i < 5; i++ {
			r += rules[i].String()
		}
		return
	}
	for _, st := range stages {
		t.Run(st.pattern, func(t *testing.T) {
			var db DB
			pos, err := db.parseRules(st.pattern, st.negative)
			if err != nil {
				if err != st.err {
					t.Error(err)
				}
				return
			}
			s := conv(pos)
			if s != st.expect {
				t.Errorf("pattern parse fail: need '%s', got '%s'", st.expect, s)
			}
		})
	}
}
