package main

func typeCheck(context contextStack, expressionTree *node, typeTree *node) bool {
	inferredType := typeInference(context, expressionTree)
	return inferredType.compareSubtrees(typeTree)
}
