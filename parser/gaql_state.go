package parser

import (
	"encoding/json"
	"fmt"
	"math/big"
)

const (
	TypeNumber = iota
	TypeString
	TypeTBD
)

type SelectStatement struct {
	Selection []*SelectField
	From      *FromClause
	Where     Expression
	Group     []string
	Order     []*OrderByClause
}

type FromClause struct {
	DataPath     string
	SubQuery     *SelectStatement
	SubQueryName string
}

type OrderByClause struct {
	Field string
	ASC   bool
}

type SelectField struct {
	Exp  Expression
	FExp *FuncStatement
	Name string
}

type Expression interface {
	GetType() (int32, error)
	String() string
	Value() interface{}
}

type Number big.Rat

func (n *Number) GetType() (ret int32, reterr error) {
	return TypeNumber, nil
}

func (n *Number) String() (ret string) {
	nn := (*big.Rat)(n)
	if nn.IsInt() {
		ret = nn.Num().String()
	} else {
		ret = nn.String()
	}

	return
}

func (n *Number) Value() interface{} {
	return n
}

type Name string

func (n Name) GetType() (ret int32, reterr error) {
	return TypeTBD, nil
}

func (n Name) String() (ret string) {
	return string(n)
}

func (n Name) Value() interface{} {
	return nil
}

type ScalaExp struct {
	Val      [2]Expression
	Operator string
}

func (e *ScalaExp) GetType() (ret int32, reterr error) {
	t1, err1 := e.Val[0].GetType()
	if err1 != nil {
		reterr = err1
		return
	}

	t2, err2 := e.Val[1].GetType()
	if err2 != nil {
		reterr = err2
		return
	}

	if t1 != t2 {
		reterr = fmt.Errorf("\"%s\" with type %d does not match \"%s\" with type %d", e.Val[0].String(), t1, e.Val[1].String(), t2)
	}

	return
}

func (e *ScalaExp) String() (ret string) {
	return fmt.Sprintf("(%s) %s (%s)", e.Val[0].String(), e.Operator, e.Val[1].String())
}

func (e *ScalaExp) Value() interface{} {
	return nil
}

type UnaryScalaExp struct {
	Operator string
	Val      Expression
}

func (e *UnaryScalaExp) GetType() (ret int32, reterr error) {
	return e.Val.GetType()
}

func (e *UnaryScalaExp) Value() interface{} {
	return nil
}

func (e *UnaryScalaExp) String() string {
	return fmt.Sprintf("%s (%s)", e.Val.String(), e.Operator)
}

type FuncStatement struct {
	FuncName  string
	UsingAll  bool
	Arguments []Expression
}

func EndGaql(gaql interface{}) {
	data, _ := json.Marshal(gaql)
	fmt.Printf("gaql = %s", data)
}
