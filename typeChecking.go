package main

import (
	"fmt"
)

//*************************************************************************

func typeCheck(context contextStack, expressionTree *node, typeTree *node) error {
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

func typeInference(context contextStack, expressionTree *node) (*node, error) {
	switch expressionTree.token {
	case VARIABLE:
		// if top node is a variable just return the type in the context (or nil)
		contextType := context.getType(expressionTree.value)
		if contextType == nil {
			return nil, fmt.Errorf("variable %s does not have a type", expressionTree.value)
		} else {
			return contextType, nil
		}
	case LAMBDA:
		//add statement x:T in the context
		context.addStatement(expressionTree.value, expressionTree.right)
		//get predicted type for the rest of the expression
		right, err := typeInference(context, expressionTree.left)

		//returning nil if there was no prediction for the rest
		//not the most elegant solution, but probably saves a lot of headaches
		if right == nil {
			return nil, err
		} else {
			//create the predicted type
			top := newNode("", ARROW)
			left := getCopySubtree(expressionTree.right)
			top.linkNodes(left, right)
			return top, nil
		}
	case APPLICATION:
		//get types for left and right of the application
		leftType, LErr := typeInference(context.getCopy(), expressionTree.left)
		rightType, _ := typeInference(context.getCopy(), expressionTree.right)
		//if the conditions are right, return the right of the arrow type
		// (T in the rule)
		if leftType == nil || rightType == nil {
			return nil, LErr
		}
		if leftType.token == ARROW &&
			leftType.left.compareSubtrees(rightType) {
			return leftType.right, nil
		}
	}
	return nil, fmt.Errorf("cannot infer type")
}

//*************************************************************************

//for testing purposes
func testTypeInference() {

	cs := contextStack{nil}
	/*
		{
			//block specific to the expression "a b"
			//created a context to make a valid prediction

			// a : A->B
			boh := newNode("->", ARROW)
			boh.linkNodes(newNode("A", VARIABLE), newNode("B", VARIABLE))
			cs.addStatement("a", boh)

			// b : a
			cs.addStatement("b", boh.left)
		}
	*/
	t, err := typeInference(cs, rootExpressionNode)
	println("EXPRESSION: ", rootExpressionNode.toString())
	checkError(err)
	println("TYPE PREDICTED: ", t.toString())
	println()
	println()
}
