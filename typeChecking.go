package main

//Main function for type checking
func typeCheck(context contextStack, expressionTree *node, typeTree *node) bool {

	switch expressionTree.token {
	case VARIABLE:
		//Search for the statement x : T
		return context.findStatement(expressionTree.value, typeTree)
	case LAMBDA:
		//prerequisites for the application of the rule
		if typeTree.token == ARROW && compareSubtrees(typeTree.left, expressionTree.right) {
			//add statement to gamma
			context.addStatement(expressionTree.value, getCopySubtree(expressionTree.right))
			//recursively calls itself to check the rest of the expression
			return typeCheck(context, expressionTree.left, typeTree.right)
		} else {
			return false
		}
	case APPLICATION:
		//Coming soon....
	default:
		return false
	}
	print("SHOULDN'T BE HERE LOL \n")
	return false
}