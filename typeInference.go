package main

//Context must be a COPY when called from typecheck
func typeInference(context contextStack, expressionTree *node) *node {
	switch expressionTree.token {
	case VARIABLE:
		// if top node is a variable just return the type in the context (or nil)
		return context.getType(expressionTree.value)
	case LAMBDA:
		//add statement x:T in the context
		context.addStatement(expressionTree.value, expressionTree.right)
		//get predicted type for the rest of the expression
		right := typeInference(context, expressionTree.left)

		//returning nil if there was no prediction for the rest
		//not the most elegant solution, but probably saves a lot of headaches
		if right == nil {
			return nil
		} else {
			//create the predicted type
			top := newNode("", ARROW)
			//TODO: FIX THIS
			left := getCopySubtree(expressionTree.right) //possibly wrong - creating a copy without references to children
			//make a copy of the expressionTree.right (which is its type) and then link it to top
			top.linkNodes(left, right)
			return top
		}
	case APPLICATION:
		//get types for left and right of the application
		leftType := typeInference(context.getCopy(), expressionTree.left)
		rightType := typeInference(context.getCopy(), expressionTree.right)
		//if the conditions are right, return the right of the arrow type
		// (T in the rule)
		if leftType == nil || rightType == nil {
			return nil
		}
		if leftType.token == ARROW &&
			leftType.left.compareSubtrees(rightType) {
			return leftType.right
		}
	}
	return nil
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
	t := typeInference(cs, rootExpressionNode)
	println("EXPRESSION: ", rootExpressionNode.toString())
	println("TYPE PREDICTED: ", t.toString())
	println()
	println()
}
