package main

import (
	"fmt"
	"os"
)
//Main function for type checking
func typeCheck(context contextStack, expressionTree *node, typeTree *node) (bool, *node) {
	switch expressionTree.token {
	case VARIABLE:
		/* Search in context for the 'Variable : Type' */
		isInContext, statementInContext := context.findStatement(expressionTree.value, typeTree)
		if (!isInContext) {
			print("Cannot find variable in context\n")
		}
		return isInContext, statementInContext

	case LAMBDA:
		/* Prerequisite: Check T of the lambda-expression, and T of the T->T' */
		if typeTree.token == ARROW && compareSubtrees(expressionTree.right, typeTree.left)  {
			/* Add statement (with the correct T) to the context */
			context.addStatement(expressionTree.value, expressionTree.right)

			/* Check rest of the expression: Advance to E (from \x^(T) E), and T' (from T->T') */
			typeCheckBool, statement := typeCheck(context, expressionTree.left, typeTree.right)

			/* Statement = T->T' */
			tNode := newNode(expressionTree.right.value, VARIABLE)
			arrow := newNode("", ARROW)
			arrow.linkNodes(tNode, statement)
			statement = arrow

			return typeCheckBool, statement
		}

	case APPLICATION:
		print("\n")
		E1Expression := expressionTree.left
		fmt.Fprintf(os.Stdout, "E1Expression: %s\n", E1Expression.toString())
		E2Expression := expressionTree.right
		fmt.Fprintf(os.Stdout, "E2Expression: %s\n", E2Expression.toString())
		E2TypeTree := newNode("?" , VARIABLE)	// T, may contain questionmarks
		n := newNode("?", VARIABLE)
		E1TypeTree := newNode("", ARROW)	// T->T' may contains questionmarks
		E1TypeTree.linkNodes(n, typeTree)
		fmt.Fprintf(os.Stdout, "E1TypeTree: %s\n", E1TypeTree.toString())

		/* E1Statement = Statement from new call (may contain additional statements apart from <T -> T'>) */
		E1Bool, E1Statement := typeCheck(context, E1Expression, E1TypeTree)
		fmt.Fprintf(os.Stdout, "E1Statement: %s\n", E1Statement.toString())

		/* Typecheck error found in further calls */
		if (!E1Bool) {
			print("E1Bool is false\n")
			return false, nil
		}
		
		/* E1TypeTreeNew = The actual T->T' statement, where the questionmarks has been determined */
		E1TypeTreeNew := getQuestionNode(E1TypeTree, E1Statement) // To compare with old typeTree
		fmt.Fprintf(os.Stdout, "E1TypeTreeNew: %s\n", E1TypeTreeNew.toString())
		if (!compareSubtrees(E1TypeTree, E1TypeTreeNew)) {
			print("whelp something went wrong again")
		}
		
		/* E2TypeTree = T of E1 (T -> T') */
		E2TypeTree = E1TypeTreeNew.left
		fmt.Fprintf(os.Stdout, "E2TypeTree: %s\n", E2TypeTree.toString())

		/* E2Statement = Statement from new call, may contain additoinal statements apart from <T> */
		E2Bool, _ := typeCheck(context, E2Expression, E2TypeTree)

		/* Apparently a typecheck error in the further calls from E2 */
		if (!E2Bool) {
			return false, nil
		}

		return true, E1TypeTreeNew
	}

	return false, nil
}

// checkTree: The tree with the questionmark
// answerTree: The tree with no questionmarks
func getQuestionNode(checkTree *node, answerTree *node) *node {
	checkTreeNrNodes := countNodes(checkTree)
	answerTreeNrNodes := countNodes(answerTree)
	if (answerTreeNrNodes < checkTreeNrNodes) {
		print("Whelp something went wrong")
		return nil
	}

	for answerTreeNrNodes > checkTreeNrNodes {
		answerTree = answerTree.right
		answerTreeNrNodes = countNodes(answerTree)
		print("AnswerTreeNrNodes: %d", answerTreeNrNodes)
	}
	fmt.Fprintf(os.Stdout, "QuestionNode:")
	return answerTree
}

func countNodes(theNode *node) int {
	if theNode == nil {
		return 0
	}

	return 1 + countNodes(theNode.left) + countNodes(theNode.right)
}