// typeChecking.go
// Programming Language: GoLang
//
// Course: Concepts of Programming Language
// Assignment 3: Type Checking
// Class 2, Group 11
// Author(s) :	Emanuele Greco (s3375951),
//				Irie Railton (s3292037),
//				Kah ming Wong (s2641976).
//
// Date: 30th November, 2021.
//

//*************************************************************************

package main

import (
	"fmt"
)

//*************************************************************************

func typeCheck(expressionTree *node, typeTree *node) error {
	context := contextStack{nil}
	inferredType, err := typeInference(context, expressionTree)
	if inferredType == nil {
		return err
	} else {
		if inferredType.compareSubtrees(typeTree) {
			return nil
		} else {
			return fmt.Errorf("does not type check")
		}
	}
}

//*************************************************************************

func typeInference(context contextStack, n *node) (*node, error) {
	switch n.token {
	case VARIABLE:
		//if top node is a variable return the type in context (or nil)
		contextType := context.getType(n.value)
		if contextType == nil {
			return nil,
				fmt.Errorf("variable %s does not have a type", n.value)
		} else {
			return contextType, nil
		}
	case LAMBDA:
		//add statement x:T in the context
		context.addStatement(n.value, n.right)
		//get predicted type for the rest of the expression
		right, err := typeInference(context, n.left)

		//returning nil if there was no prediction for the rest
		if right == nil {
			return nil, err
		} else {
			//create the predicted type
			top := newNode("", ARROW)
			left := getCopySubtree(n.right)
			top.linkNodes(left, right)
			return top, nil
		}
	case APPLICATION:
		//get types for left and right of the application
		leftType, err := typeInference(context.getCopy(), n.left)
		rightType, _ := typeInference(context.getCopy(), n.right)
		//if conditions are right, return the right of the arrow type
		// (T' in the rule)
		if leftType == nil || rightType == nil {
			return nil, err
		}
		if leftType.token == ARROW &&
			leftType.left.compareSubtrees(rightType) {
			return leftType.right, nil
		}
	}
	return nil, fmt.Errorf("cannot infer type")
}

//*************************************************************************
