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
	/* Case (x: T) */
	case VARIABLE:
		/* If (x: T): Return T present in the context. If not present: nil */
		contextType := context.getType(n.value)
		if contextType == nil {
			return nil,
				fmt.Errorf("variable %s does not have a type", n.value)
		} else {
			return contextType, nil
		}

	/* Case (\x^T E): T->T' */
	case LAMBDA:
		/* Add T of \x^T to the context */
		context.addStatement(n.value, n.right)
		
		/* Infer what T' should be, by passing E to further calls */
		right, err := typeInference(context, n.left)

		/* Return nil if T' cannot be inferred, otherwise return (T->T') */
		if right == nil {
			return nil, err
		} else {
			/* Create a tree in the structure of T->T' */
			top := newNode("", ARROW)	// root of the tree
			left := getCopySubtree(n.right)	// T (of \x^T E)
			top.linkNodes(left, right)
			return top, nil
		}
	
	/* Case (E1E2): T' */
	case APPLICATION:
		/* Infer T->T' of E1, and T of E2*/
		leftType, err := typeInference(context.getCopy(), n.left)
		rightType, _ := typeInference(context.getCopy(), n.right)
		
		/* Either T->T' of E1, or T of E2 cannot be inferred */
		if leftType == nil || rightType == nil {
			return nil, err
		}

		/* Compare T of E1 and T of E2, if equivalent return T' (of E1) */
		if leftType.token == ARROW &&
			leftType.left.compareSubtrees(rightType) {
			return leftType.right, nil
		}
	}

	/* Default case */
	return nil, fmt.Errorf("cannot infer type")
}

//*************************************************************************
