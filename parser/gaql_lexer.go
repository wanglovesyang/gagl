package parser

import (
	"fmt"
	lx "gaql/parser/lexer"
	"log"
	"math/big"
)

type MGAQLLexer struct {
	lexer *lx.SQLLexer
}

func NewMGAQLLexer(src string, verbose bool) *MGAQLLexer {
	l, _ := lx.NewSQLLexer(src, verbose)
	return &MGAQLLexer{lexer: l}
}

func (e *MGAQLLexer) Error(s string) {
	log.Fatalf("Error in lexer, %s", s)
}

func (e *MGAQLLexer) Errorf(ft string, args ...interface{}) {
	e.Error(fmt.Sprintf(ft, args...))
}

func (e *MGAQLLexer) Lex(s *GAQLSymType) (ret int) {
	ptn, tag, err := e.lexer.Next()
	if err == lx.ErrLexerEOF {
		return 0
	}

	if err != nil {
		e.Errorf("Error in lexer next, %v", err)
	}

	switch tag {
	case lx.NUMBER:
		s.numval = &big.Rat{}
		s.numval.SetString(ptn)
	case lx.STRING, lx.NAME:
		s.strval = ptn
	case lx.AMPQ, lx.AMAS, lx.AMCP, lx.BRC, lx.AND, lx.OR, lx.NOT:
		s.subtok = ptn
	case lx.SELECT, lx.WHERE, lx.AS, lx.FROM, lx.GROUP, lx.BY, lx.ORDER, lx.ASC, lx.DESC:
	case lx.SEMCO, lx.COMMA:
	default:
		e.Errorf("Invalid tag in expression, %s, %d", ptn, tag)
	}

	switch tag {
	case lx.SELECT:
		ret = SELECT
	case lx.WHERE:
		ret = WHERE
	case lx.AS:
		ret = AS
	case lx.FROM:
		ret = FROM
	case lx.GROUP:
		ret = GROUP
	case lx.BY:
		ret = BY
	case lx.ORDER:
		ret = ORDER
	case lx.ASC:
		ret = ASC
	case lx.DESC:
		ret = DESC
	case lx.AND:
		ret = AND
	case lx.OR:
		ret = OR
	case lx.NOT:
		ret = NOT
	case lx.NUMBER:
		ret = NUMBER
	case lx.BRC:
		ret = int(ptn[0])
	case lx.STRING:
		ret = STRING
	case lx.NAME:
		ret = NAME
	case lx.AMPQ:
		ret = int(ptn[0])
	case lx.AMAS:
		ret = AMAS
	case lx.AMCP:
		ret = AMCP
	case lx.SEMCO:
		ret = int(';')
	case lx.COMMA:
		ret = ','
	default:
		e.Errorf("Invalid tag in expression, %s, %d", ptn, tag)
	}

	return
}
