package lex

import (
	"fmt"
	"log"
	"regexp"
)

var ErrLexerEOF = fmt.Errorf("No more source code")

const (
	SELECT = 256 + iota
	FROM
	WHERE
	AS
	BY
	GROUP
	ORDER
	ASC
	DESC
	AND
	OR
	XOR
	NOT
	AMAS
	AMPQ
	AMCP
	BRC
	COMMA
	SEMCO
	QUOT
	DOT
	STAR
	NUMBER
	NAME
	STRING
)

var SQLPatterns = []LexPattern{
	LexPattern{false, "SELECT", SELECT},
	LexPattern{false, "FROM", FROM},
	LexPattern{false, "WHERE", WHERE},
	LexPattern{false, "AS", AS},
	LexPattern{false, "BY", BY},
	LexPattern{false, "GROUP", GROUP},
	LexPattern{false, "ORDER", ORDER},
	LexPattern{false, "ASC", ASC},
	LexPattern{false, "DESC", DESC},
	LexPattern{false, "AND", AND},
	LexPattern{false, "OR", OR},
	LexPattern{false, "XOR", XOR},
	LexPattern{false, "NOT", NOT},
	LexPattern{false, "+", AMAS},
	LexPattern{false, "-", AMAS},
	LexPattern{false, "*", AMPQ},
	LexPattern{false, "/", AMPQ},
	LexPattern{false, ">", AMCP},
	LexPattern{false, ">=", AMCP},
	LexPattern{false, "<", AMCP},
	LexPattern{false, "<=", AMCP},
	LexPattern{false, "<>", AMCP},
	LexPattern{false, "=", AMCP},
	LexPattern{false, ")", BRC},
	LexPattern{false, "(", BRC},
	LexPattern{false, ".", DOT},
	LexPattern{false, ",", COMMA},
	LexPattern{false, ";", SEMCO},
	LexPattern{true, "[A-Za-z][A-Za-z0-9_]*", NAME},
	LexPattern{true, "\"[^\"]*\"", STRING},
	LexPattern{true, "-?\\d+(\\.\\d+)?([eE][+-]?[0-9]+)?", NUMBER},
}

type RegexDetector struct {
	Tag      int32
	Detector *regexp.Regexp
}

func NewRegexDetector(tag int32, exp string) (*RegexDetector, error) {
	if exp[0] != '^' {
		exp = "^" + exp
	}

	e, err := regexp.Compile(exp)
	if err != nil {
		return nil, err
	}

	return &RegexDetector{Tag: tag, Detector: e}, nil
}

func (r *RegexDetector) Detect(s string) (rl int32, suc bool) {
	loc := r.Detector.FindStringIndex(s)
	if len(loc) == 0 {
		suc = false
		return
	}

	if loc[0] != 0 {
		suc = false
		return
	}

	suc = true
	rl = int32(loc[1])

	return
}

type LexPattern struct {
	IsRegex bool
	Pattern string
	Tag     int32
}

type Lexer struct {
	tokenTree *FlexTree
	regDets   []*RegexDetector
}

func NewLexer(patterns []LexPattern) (ret *Lexer, reterr error) {
	ret = &Lexer{}

	treePts := make(map[string]NodeTag)
	for _, pt := range patterns {
		if pt.IsRegex {
			if nd, err := NewRegexDetector(pt.Tag, pt.Pattern); err != nil {
				reterr = err
				break
			} else {
				ret.regDets = append(ret.regDets, nd)
			}
		} else {
			treePts[pt.Pattern] = NodeTag{Tag: pt.Tag}
		}
	}

	if ttk, err := BuildFlexTree(treePts); err != nil {
		reterr = err
	} else {
		ret.tokenTree = ttk
	}

	//log.Printf("Size regdets = %d", len(ret.regDets))

	return
}

func (l *Lexer) Detect(s string) (pattern string, tag int32, reterr error) {
	pt, tg, err := l.tokenTree.Recognize(s)
	if err != nil && err != ErrNothingMatch {
		reterr = fmt.Errorf("Error in tree Recognize, %v", err)
	}

	ln := int32(len(pt))

	ml := int32(-1)
	mi := -1
	for i, d := range l.regDets {
		l, suc := d.Detect(s)
		if suc {
			if l > ml {
				ml = l
				mi = i
			}
		}
	}

	if ml > ln {
		pattern = s[0:ml]
		tag = l.regDets[mi].Tag
	} else {
		pattern = s[0:ln]
		tag = tg.Tag
	}

	if ml < 0 && err == ErrNothingMatch {
		reterr = ErrNothingMatch
	}

	return
}

type SQLLexer struct {
	Lexer
	src     string
	cur     int32
	verbose bool
}

func NewSQLLexer(src string, verbose bool) (ret *SQLLexer, reterr error) {
	l, err := NewLexer(SQLPatterns)
	if err != nil {
		reterr = err
		return
	}

	ret = &SQLLexer{Lexer: *l, src: src, cur: 0, verbose: verbose}
	return
}

func (l *SQLLexer) Skip() {
	for ; l.cur < int32(len(l.src)) && (l.src[l.cur] == ' ' || l.src[l.cur] == '\t'); l.cur++ {
	}
}

func (l *SQLLexer) Next() (ptn string, tag int32, reterr error) {
	defer func() {
		if l.verbose {
			if reterr == nil {
				log.Printf("Match: %s", ptn)
			} else {
				log.Printf("Match EOF")
			}
		}
	}()

	l.Skip()
	if l.cur == int32(len(l.src)) {
		reterr = ErrLexerEOF
		return
	}

	ptn, tag, reterr = l.Detect(l.src[l.cur:])
	if reterr == nil {
		l.cur += int32(len(ptn))
	}

	return
}
