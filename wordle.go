package wordle

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

const (
	WordSize = 5
)

type DB struct {
	buf []string
}

type rule struct {
	exact byte
	buf   []byte
}

var (
	patternRE = regexp.MustCompile(`([A-Z]|[*]|\^[A-Z])`)
)

func (db *DB) Load(path string) (err error) {
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return
	}
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	db.buf = db.buf[:0]
	scr := bufio.NewScanner(f)
	for scr.Scan() {
		db.buf = append(db.buf, scr.Text())
	}
	err = scr.Err()
	return
}

func (db DB) Unwordle(dst []string, pattern, negative string) ([]string, error) {
	rules, req, err := db.parseRules(pattern, negative)
	if err != nil {
		return dst, err
	}
	dst = dst[:0]
	for i := 0; i < len(db.buf); i++ {
		if db.checkMatch(db.buf[i], rules, req) {
			dst = append(dst, db.buf[i])
		}
	}
	return dst, nil
}

func (db DB) checkMatch(x string, rules [5]rule, req []byte) bool {
	if len(x) != WordSize {
		return false
	}
	xl := strings.ToLower(x)
	for i := 0; i < WordSize; i++ {
		r := rules[i]
		c := xl[i]
		if r.exact > 0 {
			if r.exact != c {
				return false
			}
			continue
		}
		var m bool
		for j := 0; j < len(r.buf); j++ {
			if r.buf[j] == c {
				m = true
				break
			}
		}
		if !m {
			return false
		}
	}
	for i := 0; i < len(req); i++ {
		var m bool
		for j := 0; j < WordSize; j++ {
			if xl[j] == req[i] {
				m = true
				break
			}
		}
		if !m {
			return false
		}
	}
	return true
}

func (db DB) parseRules(pattern, negative string) (r [WordSize]rule, req []byte, err error) {
	m := patternRE.FindAllStringSubmatch(pattern, -1)
	if len(m) != WordSize {
		err = fmt.Errorf("pattern has %d rules, but need %d", len(m), WordSize)
		return
	}
	allow := map[byte]struct{}{
		'a': {}, 'b': {}, 'c': {}, 'd': {}, 'e': {}, 'f': {}, 'g': {}, 'h': {}, 'i': {}, 'j': {}, 'k': {}, 'l': {},
		'm': {}, 'n': {}, 'o': {}, 'p': {}, 'q': {}, 'r': {}, 's': {}, 't': {}, 'u': {}, 'v': {}, 'w': {}, 'x': {},
		'y': {}, 'z': {},
	}
	neg := map[byte]struct{}{}
	for i := 0; i < len(negative); i++ {
		neg[negative[i]] = struct{}{}
		delete(allow, negative[i])
	}
	exact := make(map[byte]struct{}, WordSize)
	for i := 0; i < WordSize; i++ {
		b := m[i][0][0]
		if b == '^' {
			b = m[i][0][1]
			if _, ok := neg[b]; ok {
				err = fmt.Errorf("char %s present in both pattern and negative string", string(b))
				return
			}
			exact[b] = struct{}{}
			req = append(req, b)
		} else if b != '*' {
			exact[b] = struct{}{}
			delete(allow, b)
		}
	}
	for i := 0; i < WordSize; i++ {
		switch {
		case m[i][0] == "*":
			for b := range allow {
				r[i].buf = append(r[i].buf, b)
			}
		case m[i][0][0] == '^':
			b := m[i][0][1]
			for b1 := range allow {
				if b == b1 {
					continue
				}
				r[i].buf = append(r[i].buf, b1)
			}
		default:
			r[i].exact = m[i][0][0]
		}
	}
	for k := 0; k < WordSize; k++ {
		sort.Slice(r[k].buf, func(i, j int) bool {
			return r[k].buf[i] < r[k].buf[j]
		})
	}
	return
}

func (p rule) String() string {
	if p.exact > 0 {
		return string(p.exact)
	}
	if len(p.buf) > 0 {
		return "[" + string(p.buf) + "]"
	}
	return ""
}
