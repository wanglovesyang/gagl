package lex

import (
	"testing"
)

func TestSQLLexer_source_lex(t *testing.T) {
	code := "SELECT ip AS port, 89.8E-5 AS num, \"str\" AS pool FROM \"http/share\" WHERE a < 5"
	l, err := NewSQLLexer(code)
	if err != nil {
		t.Fatalf("Fail to create sql, %v", err)
	}

	cnt := 0
	ptn, tg, err := l.Next()
	for ; err == nil; ptn, tg, err = l.Next() {
		t.Logf("Token[%d] = %s(%d)", cnt, ptn, tg)
		cnt++
	}

	if err != ErrLexerEOF {
		t.Fatalf("Error in lexer, %v", err)
	}

	return
}
