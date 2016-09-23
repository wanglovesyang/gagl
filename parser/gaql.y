/* symbolic tokens */

%{

package parser
import (
	"math/big"
)

%}

%union {
	numval *big.Rat
	strval string;
	subtok string;
	boolval bool
	slct *SelectStatement
	exp Expression
	exps []Expression
	slf *SelectField
	slfs []*SelectField
	strlst []string
	fc *FromClause
	oc *OrderByClause
	ocl []*OrderByClause
	fnc *FuncStatement
	intf interface{}
}
	
%token <strval> STRING, NAME
%token <numval> NUMBER

%type <slct> select_statement table_exp sub_query
%type <slf> select_item
%type <slfs> select_item_commalist selection
%type <strlst> column_ref_commalist opt_group_by_clause
%type <exp> where_clause opt_where_clause scalar_exp scalar_exp_1 scalar_exp_2 scalar_exp_3 scalar_exp_4
%type <exps> scalar_exp_comma_list
%type <fc> from_clause
%type <oc> ordering_spec
%type <ocl> ordering_spec_commalist opt_order_by_clause
%type <boolval> opt_asc_desc
%type <fnc> function_ref
%type <intf> gaql

	/* operators */
%left <subtok> OR
%left <subtok> AND
%left <subtok> NOT
%left <subtok> AMCP /* = <> < > <= >= */
%left <subtok> AMAS
%left <subtok> '*'
%left <subtok> '/'
%left <subtok> ';'
%nonassoc UMINUS

	/* literal keyword tokens */

%token <subtok> ALL AMMSC ANY AS ASC AUTHORIZATION BETWEEN BY SFUNC
%token <subtok> CHARACTER CHECK CLOSE COMMIT CONTINUE CREATE CURRENT
%token <subtok> CURSOR DECIMAL DECLARE DEFAULT DELETE DESC DISTINCT DOUBLE
%token <subtok> ESCAPE EXISTS FETCH FLOAT FOR FOREIGN FOUND FROM GOTO
%token <subtok> GRANT GROUP HAVING IN INDICATOR INSERT INTEGER INTO
%token <subtok> IS KEY LANGUAGE LIKE NULLX NUMERIC OF ON OPEN OPTION
%token <subtok> ORDER PARAMETER PRECISION PRIMARY PRIVILEGES PROCEDURE
%token <subtok> PUBLIC REAL REFERENCES ROLLBACK SCHEMA SELECT SET
%token <subtok> SMALLINT SOME SQLCODE SQLERROR TABLE TO UNION
%token <subtok> UNIQUE UPDATE USER VALUES VIEW WHENEVER WHERE WITH WORK


%%

gaql_list:
		gaql ';'	{ EndGaql($1); }
	|	gaql_list ';' gaql { EndGaql($2); }
	;

gaql:
        select_statement {$$ = $1}
;

select_statement:
    SELECT selection  
	table_exp
	{$3.Selection = $2; $$ = $3}
;

selection:
		select_item_commalist 
	|	'*' {$$ = []*SelectField{&SelectField{Name:"*"}}}
	;

select_item_commalist:
		select_item {$$ = []*SelectField{$1}}
	|	select_item_commalist ',' select_item {$$ = append($1, $3)}
	;

select_item:
		scalar_exp AS NAME {$$ = &SelectField{Name: $3, Exp: $1}}
	| 	scalar_exp {$$ = &SelectField{Name: "unknown", Exp: $1}}
	|	function_ref AS NAME {$$ = &SelectField{Name: $3, FExp: $1}}
	|	function_ref {$$ = &SelectField{Name: "unknown", FExp: $1}}
	;

scalar_exp:
		scalar_exp_1
	|	scalar_exp_1 OR scalar_exp_1 {$$ = &ScalaExp{Val:[2]Expression{$1, $3}, Operator:$2}}
	|	scalar_exp_1 AND scalar_exp_1 {$$ = &ScalaExp{Val:[2]Expression{$1, $3}, Operator:$2}}
	|   NOT scalar_exp_1 {$$ = &UnaryScalaExp{Val:$2, Operator:$1}}
	;

scalar_exp_1:
		scalar_exp_2
	|	scalar_exp_2 AMCP scalar_exp_2 {$$ = &ScalaExp{Val:[2]Expression{$1, $3}, Operator:$2}}
	;

scalar_exp_2:
		scalar_exp_3
	|	scalar_exp_2 AMAS scalar_exp_3 {$$ = &ScalaExp{Val:[2]Expression{$1, $3}, Operator:$2}}
	;

scalar_exp_3:
		scalar_exp_4
	|	scalar_exp_3 '*' scalar_exp_4 {$$ = &ScalaExp{Val:[2]Expression{$1, $3}, Operator:$2}}
	|	scalar_exp_3 '/' scalar_exp_4 {$$ = &ScalaExp{Val:[2]Expression{$1, $3}, Operator:$2}}
	;

scalar_exp_4:
		NUMBER {$$ = (*Number)($1)}
	|	NAME {$$ = Name($1)}
	|	'(' scalar_exp_2 ')' {$$ = $2}
	;

scalar_exp_comma_list:
		scalar_exp {$$ = []Expression{$1}}
	|	scalar_exp_comma_list ',' scalar_exp {$$ = append($1, $3)}
	;

table_exp:
		from_clause 
		opt_where_clause
		opt_group_by_clause
		opt_order_by_clause
		{$$ = &SelectStatement{From:$1, Where:$2, Group:$3, Order:$4}}
	;

opt_where_clause:
		/* empty */ 
		{$$ = nil}
	|	where_clause {$$ = $1}
	;

where_clause:
		WHERE scalar_exp {$$ = $2}
	;

from_clause:
		FROM STRING {$$ = &FromClause{DataPath: $2}}
    |   FROM sub_query NAME {$$ = &FromClause{SubQuery: $2, SubQueryName: $3}}
	;

sub_query:
	'(' select_statement ')' {$$ = $2}
	;

opt_group_by_clause:
		/* empty */
		{$$ = nil}
	|	GROUP BY column_ref_commalist {$$ = $3}
	;

opt_order_by_clause:
		/* empty */
		{$$ = nil}
	|	ORDER BY ordering_spec_commalist {$$ = $3}
	;

ordering_spec_commalist:
		ordering_spec {$$ = []*OrderByClause{$1}}
	|	ordering_spec_commalist ',' ordering_spec {$$ = append($1, $3)}
	;

ordering_spec:
		NAME opt_asc_desc {$$ = &OrderByClause{Field: $1, ASC: $2}}
	;

opt_asc_desc:
		/* empty */ {$$ = true}
	|	ASC {$$ = true}
	|	DESC {$$ = false}
	;

column_ref_commalist:
		NAME {$$ = []string{$1}}
	|	column_ref_commalist ',' NAME {$$ = append($1, $3)}
	;

     /* changed by yangyang, removing DISTINCT and ALL */
function_ref:
		NAME '(' '*' ')' {$$ = &FuncStatement{FuncName: $1, UsingAll:true}}
	|	NAME '(' scalar_exp_comma_list ')' {$$ = &FuncStatement{FuncName: $1, UsingAll: false, Arguments: $3}}
	;

%%
