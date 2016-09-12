/* symbolic tokens */

%{

package main
import (
	"math/big"
)

type SelectStatement struct {

}

type QueryStatement struct {
	
}

type WhereClause struct {

}

type FromClause struct {

}

type SelectField struct {
	Exp *ScalaExp
	Name string
}

type ScalaExp struct {

}

%}

%union {
	intval int64;
	double float64;
	strval string;
	subtok int64;
}
	
%token NAME
%token STRING
%token INTNUM APPROXNUM

	/* operators */

%left OR
%left AND
%left NOT
%left <subtok> COMPARISON /* = <> < > <= >= */
%left '+' '-'
%left '*' '/'
%nonassoc UMINUS

	/* literal keyword tokens */

%token ALL AMMSC ANY AS ASC AUTHORIZATION BETWEEN BY SFUNC
%token CHARACTER CHECK CLOSE COMMIT CONTINUE CREATE CURRENT
%token CURSOR DECIMAL DECLARE DEFAULT DELETE DESC DISTINCT DOUBLE
%token ESCAPE EXISTS FETCH FLOAT FOR FOREIGN FOUND FROM GOTO
%token GRANT GROUP HAVING IN INDICATOR INSERT INTEGER INTO
%token IS KEY LANGUAGE LIKE NULLX NUMERIC OF ON OPEN OPTION
%token ORDER PARAMETER PRECISION PRIMARY PRIVILEGES PROCEDURE
%token PUBLIC REAL REFERENCES ROLLBACK SCHEMA SELECT SET
%token SMALLINT SOME SQLCODE SQLERROR TABLE TO UNION
%token UNIQUE UPDATE USER VALUES VIEW WHENEVER WHERE WITH WORK


%%

gaql_list:
		gaql ';'	{ end_gaql(); }
	|	gaql_list gaql ';' { end_gaql(); }
	;

gaql:
        select_statement 
    |   write_statement
    |   map_statement
    |   reduce_statement
;

select_statement:
    SELECT selection  
	table_exp
;

selection:
		scalar_exp_commalist
	|	'*'
	;

scalar_exp:
		scalar_exp '+' scalar_exp
	|	scalar_exp '-' scalar_exp
	|	scalar_exp '*' scalar_exp
	|	scalar_exp '/' scalar_exp
	|	'+' scalar_exp %prec UMINUS
	|	'-' scalar_exp %prec UMINUS
	|	atom
	|	column_ref
	|	function_ref
	|	'(' scalar_exp ')'
	;

table_exp:
		from_clause
		opt_where_clause
		opt_group_by_clause
	;

from_clause:
		FROM table_ref_commalist
    |   sub_query NAME
	;

sub_query:
	'(' select_statement ')'
	;

table_commalist:
		table
	|	table_commalist ',' table
	;

opt_group_by_clause:
		/* empty */
	|	GROUP BY column_ref_commalist
	;

column_ref_commalist:
		column_ref
	|	column_ref_commalist ',' column_ref
	;

     /* changed by yangyang, removing DISTINCT and ALL */
function_ref:
		SFUNC '(' '*' ')'
	|	SFUNC '(' scalar_exp ')'
	;

column_ref:
		NAME
	|	NAME '.' NAME	/* needs semantics */
	;