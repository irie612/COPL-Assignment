// reduction.go
// Programming Language: GoLang
//
// Course: Concepts of Programming Language
// Assignment 2: Interpreter
// Class 2, Group 11
// Author(s) :	Emanuele Greco (s3375951),
//				Irie Railton (s3292037),
//				Kah ming Wong (s2641976).
//
// Date: 3rd November, 2021.
//

//*************************************************************************

package main

import (
	"fmt"
	"os"
)

//*************************************************************************

// Checks if variable is already in theSlice.
// Return true if it is already in theSlice.
// Return false if it's new.
func isPresent(variable string, theSlice []string) bool {
	for i := range theSlice {
		if theSlice[i] == variable {
			return true
		}
	}
	return false
}

//*************************************************************************

// Function to get all the bound variables (in the branch) of a lambda node.
// Expectations: <theNode> is a lambda node, and the tree is a valid tree.
// Return-value:	The pointer slice of every node (the variables) that are
// 								bound to lambda <theNode>.
func getCapturedNodes(theNode *node) []*node {
	if theNode.token != LAMBDA {
		print("Inside getCapturedNodes(). Node should be a lambda\n")
		os.Exit(1)
	}
	capturedNodes := []*node{}
	_giveCapturedNodes(theNode.value, theNode.left, &capturedNodes)
	return capturedNodes
}

//*************************************************************************

// Helper-function of <getCapturedNodes()>
//
// Recursively gives all the (bound) variables which's value is equal
// to <variableName>. If a lambda node with the same variable name as the
// initial one (= the one that is called with <getCapturedNodes()>)
// then we cut-off that branch from further calls.
// Expectations:	The tree is valid and complete.
// Result: All the bound variables are stored in <boundNodes>.
func _giveCapturedNodes(variableName string, theNode *node,
	boundNodes *[]*node) {
	if theNode == nil {
		return
	} else if theNode.token == LAMBDA && theNode.value == variableName {
		return
	} else if theNode.token == VARIABLE && theNode.value == variableName {
		*boundNodes = append(*boundNodes, theNode)
	}

	_giveCapturedNodes(variableName, theNode.left, boundNodes)
	_giveCapturedNodes(variableName, theNode.right, boundNodes)
}

//*************************************************************************

// Function to get all the free variables in the branch with <theNode> as
// the "root" of that branch.
// Expectations: The tree is valid and complete.
// Return-value: The slice with all the free variables in the branch.
func getFreeVars(theNode *node) []string {
	freeVars := []string{}
	_giveFreeVars(theNode, &freeVars)
	return freeVars
}

//*************************************************************************

// Helper-function of <getFreeVars()>
//
// Recursively gets all the free variables in the branch. The free
// variables are stored in <freeVars()>.
func _giveFreeVars(theNode *node, freeVars *[]string) {
	if theNode == nil {
		return
	}

	if theNode.token == VARIABLE {
		theVar := theNode.value
		varIsCopy := isPresent(theVar, *freeVars)
		if !isBound(theNode, theVar) && !varIsCopy {
			*freeVars = append(*freeVars, theVar)
		}
		return
	}

	_giveFreeVars(theNode.left, freeVars)
	_giveFreeVars(theNode.right, freeVars)
}

//*************************************************************************

// Helper-function of <_giveFreeVars()>
//
// Goes up into the tree to check whether the leaf that is found in
// <_giveFreeVars()> is bound to a lambda expression.
// Result:	-True if the leaf is a free variable.
//					-False if the leaf is a bound variable.
func isBound(theNode *node, theVar string) bool {
	if theNode == nil {
		return false
	} else if theNode.token == LAMBDA && theNode.value == theVar {
		return true
	}

	return isBound(theNode.parent, theVar)
}

//*************************************************************************

// Function for the alpha-conversion(s).
func alphaDriver(capturedNodes []*node, freeVars []string) {
	for _, theNode := range capturedNodes {
		ogValue := theNode.value //Value of the variable from which...
		tmpNode := theNode       //we go up the tree.
		//for each captured node we go up the tree. we stop when
		//we either found nil or the lambda expression that captured the node.
		for tmpNode != nil && !(tmpNode.token == LAMBDA &&
			tmpNode.value == ogValue) {
			if tmpNode.token == LAMBDA && isPresent(tmpNode.value, freeVars) {
				// tmp_node is a lambda node and the value is different from
				// og and contained in freeVars
				alphaConversion(tmpNode, getFreshVariable(freeVars))
			}
			tmpNode = tmpNode.parent
		}
	}
}

//*************************************************************************

// Get variables bounded to node, change the value of them,
// then changes its own value.
func alphaConversion(theNode *node, freshVar string) {
	for _, capturedNode := range getCapturedNodes(theNode) {
		capturedNode.value = freshVar
	}
	theNode.value = freshVar
}

//*************************************************************************

// Get fresh variable, that is not contained in <freeVars>.
func getFreshVariable(freeVars []string) string {
	maxSize := 100
	charArray := []rune{}

	for i := 0; i < maxSize; i++ { //add a character to the string
		charArray = append(charArray, rune(int('a')-1))
		for j := 0; j < 26; j++ { //increment the last character
			charArray[i]++ //in the string (ex: a->b b->c ...)
			if !isPresent(string(charArray), freeVars) {
				return string(charArray)
			}
		}
	}
	return "#"
}

//*************************************************************************

// Main function for the beta-reductions.
// Exits the program if the maximum amount of reduction has been reached.
// Else applies beta-reduction, if possible, and return true.
func betaDriver(theNode *node) bool {
	counter := 0
	maxReduction := 100
	for betaReduction(theNode) {
		counter++
		if counter >= maxReduction {
			print("reached maximum number of reduction\n")
			os.Exit(2)
		}
	}
	return true
}

//*************************************************************************

// Applies beta-reduction once to the first applicable branch with
// preference to the left-hand side.
func betaReduction(theNode *node) bool {
	if theNode == nil || theNode.token == VARIABLE {
		return false
	}
	//Reduction possible
	if theNode.token == APPLICATION && theNode.left.token == LAMBDA &&
		theNode.right != nil {
		//check the need for alpha conversion
		//free variables on the right-hand side of the application
		freeVars := getFreeVars(theNode.right)
		//node captured by the lambda on the left-hand side of the application
		capturedNodes := getCapturedNodes(theNode.left)

		alphaDriver(capturedNodes, freeVars)

		//actual substitution
		targetVar := (theNode.left.value)
		substituteTree(theNode.left.left, theNode.right, targetVar)

		//link the sub tree of the lambda node to the its parent
		theNode.left.left.parent = theNode.parent
		//change the value stored at the address of the lambda node
		*theNode = *(theNode.left.left)

		return true
	} else {
		//Reduction is not possible. Check left branch. If no reduction found there, check right side
		if betaReduction(theNode.left) == false {
			return betaReduction(theNode.right)
		} else {
			return true
		}
	}

	return false
}

//*************************************************************************

// Substitutes the variables bound to the original lambda expression.
// <theNode>= nodes to apply substitution on.
// <subNode>= what the substitution must result into.
// targetVar= the variable of the original lambda expression.
func substituteTree(theNode *node, subNode *node, targetVar string) {
	if theNode == nil {
		return
	}

	if theNode.token == LAMBDA && theNode.value == targetVar {
		return
	}

	if theNode.token == VARIABLE && theNode.value == targetVar {
		newNode := getCopySubtree(subNode)
		parentNode := theNode.parent
		if parentNode.left == theNode {
			parentNode.left = newNode
		} else if parentNode.right == theNode {
			parentNode.right = newNode
		} else {
			print("SHOULDN'T BE HERE\n")
			os.Exit(1)
		}
		newNode.parent = parentNode
	}

	substituteTree(theNode.left, subNode, targetVar)
	substituteTree(theNode.right, subNode, targetVar)
}

//*************************************************************************

// Copies a branch.
// Return-value: node that contains a copy of a given branch.
func getCopySubtree(subtree *node) *node {
	//base case
	if subtree == nil {
		return nil
	}

	//Create a new identical node and then link a copy of the left
	//and the right of the subtree
	returnNode := newNode(subtree.value, subtree.token)
	returnNode.linkNodes(getCopySubtree(subtree.left),
		getCopySubtree(subtree.right))
	return returnNode
}

//*************************************************************************

// Prints the tree.
func printTree(theNode *node) {
	printPostOrder(theNode)
	fmt.Println()
}

//*************************************************************************

// Helper-function for <printTree()>.
func printPostOrder(theNode *node) {
	if theNode == nil {
		return
	}
	if theNode.token == LAMBDA {
		fmt.Fprintf(os.Stdout, "(")
		fmt.Fprintf(os.Stdout, "λ%s ", theNode.value)
	} else if theNode.token == VARIABLE {
		fmt.Fprintf(os.Stdout, "%s", theNode.value)
	}

	if theNode.left != nil && theNode.right != nil {
		fmt.Fprintf(os.Stdout, "(")
	}
	printPostOrder(theNode.left)

	if (theNode.left != nil && theNode.right != nil) &&
		(theNode.left.token == VARIABLE && theNode.right.token == VARIABLE) {
		fmt.Fprintf(os.Stdout, " ")
	}
	printPostOrder(theNode.right)
	if theNode.token == LAMBDA ||
		(theNode.left != nil && theNode.right != nil) {
		fmt.Fprintf(os.Stdout, ")")
	}
}

//*************************************************************************
