package main

import (
	"fmt"
	"os"
	"strconv"
)

//AST Node Struct
type node struct {
	parent			*node
	left				*node
	right				*node
	boundNodes	[]*node
	canSkip			bool
	value				string
	token				int //type of node
	bound				bool
}

/*type astTree struct {
	var rootNode *node
} */

//*************************************************************************

func newNode(value string, token int) *node {
	return &node{
		value: value,
		token: token,
		bound: false,
		canSkip: false,
	}
}

//*************************************************************************

func (n *node) linkNodes(child ...*node) {
	n.left = child[0]
	child[0].parent = n
	if len(child) > 1 {
		n.right = child[1]
		child[1].parent = n
	}
}

//*************************************************************************

func (n *node) toString() string {
	return "value = " + n.value + ", tokenID = " + strconv.Itoa(n.token) + ", bound = " + strconv.FormatBool(n.bound)
}

//*************************************************************************

func appTreeCreate(nodes []*node) *node {
	if len(nodes) == 1 {
		newNode("", APPLICATION).linkNodes(nodes[0])
	} else {
		newNode("", APPLICATION).linkNodes(nodes[0], nodes[1])
	}
	if len(nodes) > 2 {
		return appTreeCreate(append([]*node{nodes[0].parent}, nodes[2:]...))
	}
	return nodes[0].parent
}

//*************************************************************************

// Checks if a lambda abstraction has a bound variable to it.
// if not, then we can apply alpha conversion just in case.
/*func noBoundVars(theNode *node, theChar rune) bool {
	if (theNode == nil) {
		return false
	} else if (theNode.token == LAMBDA && theNode.value == theChar) {
		return false
	} 

	if (theNode.token == VARIABLE && theNode.value == theChar) {
		return true
	}

	if (noBoundVars(theNode.left, theChar) || noBoundVars(theNode.right, theChar)) {
		return true
	}

	return false
}*/

//*************************************************************************


// Check if newByte is already in theSlice
// return true if it is already in theSlice
// return false if it's new
func isCopy (theSlice []rune, newSlice []rune) (bool, rune) {
	for i:= range theSlice {
		for j := range newSlice {
			if (theSlice[i] == newSlice[j]) {
				return true, theSlice[i]
			}
		}
	}
	return false, '0'
}

//*************************************************************************

// Gives the bound variables in a branch
// You can use it for example giveBoundVars(root->left, &boundVars)
// to get the bound variables of the left branch, and compare it with the free variables
// of the right branch.
func giveBoundVars(theNode *node, boundVars *[]rune) {
	if (theNode == nil) {
		return
	}
	
	if (theNode.token == LAMBDA) {
		theVar := []rune(theNode.value)
		varIsCopy, _ := isCopy((*boundVars), theVar)
		if (!varIsCopy) {
			(*boundVars) = append((*boundVars), theVar...)
		}
	}	

	giveBoundVars(theNode.left, boundVars)
	giveBoundVars(theNode.right, boundVars)
}

//*************************************************************************

// Checks wether a variable (a leaf) is bound to a lambda abstraction
// returns true if it is bound
// returns false if it is free
// Goes up into the tree
func isBound(theNode *node, theChar rune) bool {
	if (theNode == nil) {
		return false
	}

	if (theNode.token == LAMBDA) {
		nodeVar := []rune(theNode.value)
		if (nodeVar[0] == theChar) {
			return true;
		}
	}

	return isBound(theNode.parent, theChar)
}

//*************************************************************************

// Gives the free variables in a branch in the slice <boundVars>
// Again, can be used by calling giveFreeVars(root->right, &freeVars)
// to get the free variables of the right branch, and to further compare it with
// the variables of the left branch (for example).
// Not the most efficient.
func giveFreeVars(theNode *node, freeVars *[]rune) {
	if (theNode == nil) {
		return
	}

	if (theNode.token == VARIABLE) {
		theVar := []rune(theNode.value)
		varIsCopy, _ := isCopy((*freeVars), theVar)
		if (!isBound(theNode, theVar[0]) && !varIsCopy) {
			(*freeVars) = append((*freeVars), theVar...)
		}
		return
	}

	giveFreeVars(theNode.left, freeVars)
	giveFreeVars(theNode.right, freeVars)
}

//*************************************************************************

func giveFresh(theSlice []rune, duplicateVar rune) rune {
	buffer := int(duplicateVar)
	var freshChar rune

	varIsCopy := true
	for (varIsCopy) {
		buffer++
		if (buffer > 122) {
			buffer = 65
		}
		freshChar = rune(buffer)
		freshArray := []rune{freshChar}
		varIsCopy, _ = isCopy(theSlice, freshArray)
	}

	return freshChar
}

func alphaConversion(theNode *node, duplicateVar rune, freshChar rune) {
	if (theNode == nil) {
		return
	}

	nodeVar := []rune(theNode.value)
	if ((theNode.token == LAMBDA || theNode.token == VARIABLE) && nodeVar[0] == duplicateVar) {
		theNode.value = string(freshChar)
	}

	alphaConversion(theNode.left, duplicateVar, freshChar)
	alphaConversion(theNode.right, duplicateVar, freshChar)
}

//*************************************************************************


// substitutes the variables bound to the original lambda expression.
// however, im unsure whether newNode is actually a new independent node
// from subNode (subnode= the "values" for substitution), or if it a some sort of pointer
// to the subNode. As in, every substituted node are all merely pointing to the single subNode.
func substituteTree(theNode *node, subNode *node, targetVar rune) {
	if (theNode == nil) {
		return
	}

	nodeVar := []rune(theNode.value)
	if (theNode.token == LAMBDA && nodeVar[0] == targetVar) {
		return
	}

	if (theNode.token == VARIABLE && nodeVar[0] == targetVar) {	
		var newNode *node
		newNode = subNode
		parentNode := theNode.parent
		if (parentNode.left == theNode) {
			parentNode.left = newNode
			newNode.parent = parentNode
			theNode = newNode
		} else if (parentNode.right == theNode) {
			parentNode.right = newNode
			newNode.parent = parentNode
			theNode = newNode
		}
	}

	substituteTree(theNode.left, subNode, targetVar)
	substituteTree(theNode.right, subNode, targetVar)
}


// Applies beta-reduction once to the first applicable branch with preference
// to the left hand side.
// might be working with single pointer
func applyBetaReduc(theNode **node) bool {
	if ((*theNode) == nil || (*theNode).token == VARIABLE) {
		return false
	}
	//Reduction possible
	if ((*theNode).token == APPLICATION && (*theNode).left.token == LAMBDA && (*theNode).right != nil) {
		varIsCopy:=true
		//check the need for alpha conversion
		for (varIsCopy) {
			boundVars := []rune{}
			freeVars := []rune{}
			giveBoundVars((*theNode).left, &boundVars)
			giveFreeVars((*theNode).right, &freeVars)
			varIsCopy, duplicateVar := isCopy(boundVars, freeVars)
			
			if (!varIsCopy) {
				break
			}
			
			usedVars := append(boundVars, freeVars...)
			freshVar := giveFresh(usedVars, duplicateVar)
			alphaConversion((*theNode).left, duplicateVar, freshVar)

			giveBoundVars((*theNode).left, &boundVars)
			giveFreeVars((*theNode).right, &freeVars)
			varIsCopy, _ = isCopy(boundVars, freeVars)
		}
		//actual substitution
		targetSplice := []rune((*theNode).left.value)
		targetVar := targetSplice[0]
		substituteTree((*theNode).left.left, (*theNode).right, targetVar)	

		//link the sub tree of the lambda node to the its parent
		(*theNode).left.left.parent = (*theNode).parent
		//change the value stored at the address of the lambda node
		**theNode = *((*theNode).left.left)

		//Reduction is not possible. Check left branch. If no reduction found there, check right side
	}else{
		if applyBetaReduc( (&(*theNode).left) ) == false {
			return applyBetaReduc((&(*theNode).right))
		}else{
			return true
		}
	}

	return false
} 



/*func betaReduc(theNode *node) {
	if (applyBetaReduc(theNode)) {
		
	}
}*/

// this function doesnt really have a functionality yet.
// I am merely using it to test functions with.
func checkReduction (theNode **node) bool {
	boundVars := []rune{}
	giveBoundVars((*theNode), &boundVars)
	for i := range boundVars {
		theByte := boundVars[i]
		fmt.Fprintf(os.Stdout,  "bound %c \n", theByte)
	}

	freeVars := []rune{}
	giveFreeVars((*theNode), &freeVars)
	for i := range freeVars {
		theByte := freeVars[i]
		fmt.Fprintf(os.Stdout,  "free %c \n", theByte)
	}

	if applyBetaReduc(&(*theNode)) {
		fmt.Fprintf(os.Stdout, "apply work")
	}
	
	
	return true
} 

//*************************************************************************

func testAha (theNode *node) {
	theNode.left.value = "AHA"
}

//*************************************************************************

func printPostOrder(theNode *node) {
	if (theNode == nil) {
		return
	}
	if (theNode.token == LAMBDA) {
		fmt.Fprintf(os.Stdout, "(")
	}

	if (theNode.token == LAMBDA) {
		fmt.Fprintf(os.Stdout, "λ%s ", theNode.value)
	} else if (theNode.token == VARIABLE) {
		fmt.Fprintf(os.Stdout, "%s", theNode.value)
	}

	if (theNode.left != nil && theNode.right != nil) {
		fmt.Fprintf(os.Stdout, "(")
	}
	printPostOrder(theNode.left)

	if (theNode.left != nil && theNode.right != nil) {
		fmt.Fprintf(os.Stdout, " ")
	}
	printPostOrder(theNode.right)
	if (theNode.token == LAMBDA || (theNode.left != nil && theNode.right != nil)) {
		fmt.Fprintf(os.Stdout, ")")
	}
}

//*************************************************************************

func printTree(theNode *node) {
	printPostOrder(theNode)
	fmt.Println()
}