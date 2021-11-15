// parsing.go
// Programming Language: GoLang
//
// Course: Concepts of Programming Language
// Assignment 3: Type checker
// Class 2, Group 11
// Author(s) :	Emanuele Greco (s3375951),
//				Irie Railton (s3292037),
//				Kah ming Wong (s2641976).
//
// Date: 13th November, 2021.
//
//*************************************************************************

package main

import (
	"fmt"
	"os"
	"unicode"
)

//*************************************************************************

// Parses each line in the text file, and outputs the parsed string
// given that no errors has been encountered.
func parse() {
	lex()

	if nextToken == EOF {
		return
	}
	rootExpressionNode = expr()
	if nextToken != COLON {
		_, _ = fmt.Fprintf(os.Stderr, "MISSING TYPE\n")
	} else {
		lex()
		rootTypeNode = ema_typeParse()
	}

	if nextToken != EOL && nextToken != EOF {
		_, _ = fmt.Fprintf(os.Stderr, "INPUT STRING NOT FULLY PARSED\n")
		os.Exit(1)
	}
}

//*************************************************************************
/************************************************
*		 Grammar for expression parsing			*
*	<expr> 	::= <lexpr> <expr'>					*
*	<expr'> ::= <lexpr> <expr'> | ε   	   		*
*	<lexpr>	::= <pexpr> | λ<var>^<type><expr>	*
*	<pexpr>	::= <var> | '('<expr>')'	   		*
************************************************/

func expr() *node {
	lexprNode := lexpr()
	exprPNodes := expr_p()
	if len(exprPNodes) == 0 {
		return lexprNode
	} else {
		return appTreeCreate(append([]*node{lexprNode}, exprPNodes...))
	}
}

//*************************************************************************

// Finds valid expr_p expressions. May be "empty"
func expr_p() []*node {
	if !(nextToken == EOF || nextToken == EOL || nextToken == RIGHT_P || nextToken == COLON) {
		lexprNode := lexpr()
		exprPNodes := expr_p()
		exprPNodes = append([]*node{lexprNode}, exprPNodes...)
		return exprPNodes
	}
	return []*node{}
}

//*************************************************************************

// Finds valid lambda abstractions. If there's no lambda abstractions,
// continue to pexpr.
func lexpr() *node {
	if nextToken == LAMBDA { //check if we have a lambda abstraction
		lex()
		//check if we have a variable with a valid name
		if nextToken == VARIABLE && !unicode.IsUpper(lexeme[0]) {
			lambdaNode := newNode(string(lexeme[:lexLen]), LAMBDA)
			lex()
			if nextToken == TYPE_ASS {
				lex()
				lambdaNode.right = ema_typeParse() //right == type of the lambda expression
			}
			if nextToken != EOL && nextToken != EOF {
				lambdaNode.linkNodes(lexpr())
				return lambdaNode
			} else {
				_, _ = fmt.Fprintf(os.Stderr,
					"MISSING EXPRESSION AFTER LAMBDA ABSTRACTION\n")
				os.Exit(1)
			}
		} else { // nextToken != VARIABLE ERROR
			_, _ = fmt.Fprintf(os.Stderr, "NO VALID VARIABLE AFTER LAMBDA TOKEN\n")
			os.Exit(1)
		}
		return nil
	} else {
		return pexpr()
	}
}

//*************************************************************************

// Looks for a valid pexpr expression.
func pexpr() *node {
	if nextToken == LEFT_P {
		lex()
		if nextToken == RIGHT_P {
			_, _ = fmt.Fprintf(os.Stderr,
				"MISSING EXPRESSION AFTER OPENING PARENTHESIS\n")
			os.Exit(1)
		}
		exprNode := expr()
		if nextToken != RIGHT_P {
			_, _ = fmt.Fprintf(os.Stderr, "MISSING CLOSING PARENTHESIS\n")
			os.Exit(1)
		} else {
			lex()
		}
		return exprNode
	} else if nextToken == VARIABLE && !unicode.IsUpper(lexeme[0]) { //var case
		varNode := newNode(string(lexeme[:lexLen]), VARIABLE)
		lex()
		return varNode
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "INVALID VARIABLE NAME\n")
		os.Exit(1)
		return nil
	}
}

//*************************************************************************
/*******************************************
*		 Grammar for type parsing		       *
*	<type> 	::= <ptype> '->' <type'> | <ptype> * //Could we use just <type'>?
*	<type'> ::= <ptype> '->' <type'> | ε   	   *
*	<ptype>	::= <uvar> | '('<type>')'	       *
*******************************************/
// type "A ->" doesn't work T_T
//implementation of <type> line of grammar
func typeParse() *node {
	if !(nextToken == LEFT_P || nextToken == VARIABLE) {
		_, _ = fmt.Fprintf(os.Stderr, "MISSING EXPRESSION\n")
		os.Exit(1)
	}
	ptypeNode := ptypeParse() //get left side
	//if there is a right side get it else return left side
	if nextToken == ARROW {
		lex()
		if nextToken == VARIABLE || nextToken == LEFT_P {
			arrow := newNode("", 8)
			arrow.linkNodes(ptypeNode, type_pParse()) //get right side and link
			return arrow
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "MISSING EXPRESSION OR VARIABLE AFTER ARROW\n")
			os.Exit(1)
			return nil
		}
	} else if nextToken == VARIABLE || nextToken == LEFT_P {
		_, _ = fmt.Fprintf(os.Stderr, "MISSING ARROW IN BETWEEN EXPRESSIONS\n")
		os.Exit(1)
		return nil
	} else {
		return ptypeNode
	}
}

//*************************************************************************

//implementation of <type'> line of grammar
func type_pParse() *node {
	if !(nextToken == EOF || nextToken == EOL || nextToken == RIGHT_P) {
		ptypeNode := ptypeParse()
		if nextToken == ARROW || nextToken == LEFT_P {
			lex()
			if nextToken == VARIABLE || nextToken == LEFT_P {
				arrow := newNode("", 8)
				arrow.linkNodes(ptypeNode, type_pParse()) //get right side and link
				return arrow
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "MISSING EXPRESSION OR VARIABLE AFTER ARROW\n")
				os.Exit(1)
				return nil
			}
		} else {
			return ptypeNode
		}
	} else {
		return nil
	}
}

//*************************************************************************

//implementation of <ptype> line of grammar
func ptypeParse() *node {
	if nextToken == LEFT_P {
		lex()
		if nextToken == RIGHT_P {
			_, _ = fmt.Fprintf(os.Stderr,
				"MISSING EXPRESSION AFTER OPENING PARENTHESIS\n")
			os.Exit(1)
		}
		typeNode := typeParse()
		if nextToken != RIGHT_P {
			_, _ = fmt.Fprintf(os.Stderr, "MISSING CLOSING PARENTHESIS\n")
			os.Exit(1)
		} else {
			lex()
		}
		return typeNode
	} else {
		varNode := newNode(string(lexeme[:lexLen]), VARIABLE)
		lex()
		return varNode
	}
}

/********************************************************
*		 Grammar for ema_type parsing		    		*
*	<type> 	::= <uvar> <type'> | '(' <type> )' <type'> 	*
*	<type'> ::= '->' <type> | ε   	   			   	   	*
********************************************************/
func ema_typeParse() *node {
	var leftNode *node
	var arrowNode *node
	if nextToken == LEFT_P {
		lex()
		leftNode = ema_typeParse()
		if nextToken != RIGHT_P {
			_, _ = fmt.Fprintf(os.Stdout, "MISSING RIGHT PARENTHESIS\n")
		}
		lex()
	} else if nextToken == VARIABLE && unicode.IsUpper(lexeme[0]) {
		leftNode = newNode(string(lexeme[:lexLen]), VARIABLE)
		lex()
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "ILL FORMED TYPE EXPRESSION\n")
		os.Exit(1)
	}
	arrowNode = ema_typeParse_p()
	if arrowNode == nil {
		return leftNode
	} else {
		arrowNode.linkNodes(leftNode)
		return arrowNode
	}
}

//*************************************************************************

func ema_typeParse_p() *node {
	if nextToken == ARROW {
		lex()
		arrowNode := newNode("", ARROW)
		arrowNode.linkNodes(nil, ema_typeParse())
		return arrowNode
	} else {
		return nil
	}
}
