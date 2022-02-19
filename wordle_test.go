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
			expect:  "[abcdefghijklmnopqrstuvwxyz][abcdefghijklmnopqrstuvwxyz][abcdefghijklmnopqrstuvwxyz][abcdefghijklmnopqrstuvwxyz][abcdefghijklmnopqrstuvwxyz]",
		},
		{
			pattern:  "h*^e^r*",
			negative: "uiobnm",
			expect:   "h[acdefghjklpqrstvwxyz][acdfghjklpqrstvwxyz][acdefghjklpqstvwxyz][acdefghjklpqrstvwxyz]",
		},
		{
			pattern:  "*^r^u**",
			negative: "bsh",
			expect:   "[acdefgijklmnopqrtuvwxyz][acdefgijklmnopqtuvwxyz][acdefgijklmnopqrtvwxyz][acdefgijklmnopqrtuvwxyz][acdefgijklmnopqrtuvwxyz]",
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
			rules, _, err := db.parseRules(st.pattern, st.negative)
			if err != nil {
				if err != st.err {
					t.Error(err)
				}
				return
			}
			s := conv(rules)
			if s != st.expect {
				t.Errorf("pattern parse fail: need '%s', got '%s'", st.expect, s)
			}
		})
	}
}
