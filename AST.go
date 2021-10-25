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
	if child[0] != nil {
		child[0].parent = n
	}
	if len(child) > 1 {
		n.right = child[1]
		if child[1] != nil {
			child[1].parent = n
		}
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


// Check if variable is already in theSlice
// return true if it is already in theSlice
// return false if it's new
func isPresent (variable string, theSlice []string ) (bool/*, string*/) {
	for i:= range theSlice {
		if (theSlice[i] == variable) {
			return true/*, theSlice[i]*/
		}
	}
	return false
}

//*************************************************************************

// Gives the bound variables in a branch
// You can use it for example giveBoundVars(root->left, &boundVars)
// to get the bound variables of the left branch, and compare it with the free variables
// of the right branch.
func giveBoundVars(theNode *node, boundVars *[]string) {
	if (theNode == nil) {
		return
	}
	
	if (theNode.token == LAMBDA) {
		theVar := theNode.value
		varIsCopy := isPresent(theVar,(*boundVars) )
		if (!varIsCopy) {
			(*boundVars) = append((*boundVars), theVar)
		}
	}	

	giveBoundVars(theNode.left, boundVars)
	giveBoundVars(theNode.right, boundVars)
}

//*************************************************************************


//given a variable name it searches for variable nodes with that same name starting from theNode downward
//it stops as soon as it reaches the end or a new lambda node with the same name. in the end we have variables captured
func _giveCapturedNodes(variableName string,theNode *node, boundNodes *[]*node) {
	if theNode == nil {
		return
	}else if (theNode.token == LAMBDA && theNode.value == variableName){
		return
	}else if (theNode.token == VARIABLE && theNode.value == variableName){
		*boundNodes = append(*boundNodes,theNode)
	}
	_giveCapturedNodes(variableName,theNode.left,boundNodes)
	_giveCapturedNodes(variableName,theNode.right,boundNodes)
}

//simpler interface. receive node with the lambda expression -> gives back the captured node
func getCapturedNodes(theNode *node) []*node{
	if theNode.token != LAMBDA{
		print ("Inside getCapturedNodes(). Node should be a lambda\n")
		os.Exit(1)
	}
	capturedNodes := []*node{}
	_giveCapturedNodes( theNode.value,theNode.left,&capturedNodes)
	return capturedNodes
}
//*************************************************************************

// Checks whether a variable (a leaf) is bound to a lambda abstraction
// returns true if it is bound
// returns false if it is free
// Goes up into the tree
func isBound(theNode *node, theVar string) bool {
	if (theNode == nil) {
		return false
	}

	if (theNode.token == LAMBDA) {
		if (theNode.value == theVar) {
			return true
		}
	}

	return isBound(theNode.parent, theVar)
}

//*************************************************************************

// Gives the free variables in a branch in the slice <boundVars>
// Again, can be used by calling giveFreeVars(root->right, &freeVars)
// to get the free variables of the right branch, and to further compare it with
// the variables of the left branch (for example).
// Not the most efficient.
func giveFreeVars(theNode *node, freeVars *[]string) {
	if (theNode == nil) {
		return
	}

	if (theNode.token == VARIABLE) {
		theVar := theNode.value
		varIsCopy := isPresent( theVar,(*freeVars))
		if (!isBound(theNode, theVar) && !varIsCopy) {
			(*freeVars) = append((*freeVars), theVar)
		}
		return
	}

	giveFreeVars(theNode.left, freeVars)
	giveFreeVars(theNode.right, freeVars)
}

//*************************************************************************

//driver for the alpha conversion(s)

func alphaDriver(/*theNode *node,*/ capturedNodes []*node,freeVars []string){
	for _,theNode:= range capturedNodes{
		og_value := theNode.value //value of the variable from which we go up the tree
		tmp_node := theNode
		//for each captured node we go up the tree. we stop when we either found nil
		//or the lambda expression that captured the node
		for tmp_node != nil && !(tmp_node.token == LAMBDA && tmp_node.value ==og_value){
			if tmp_node.token == LAMBDA && isPresent(tmp_node.value,freeVars){
				//tmp_node is lambda and the value is different than og and contained in freeVars
				alphaConversion(tmp_node,rune('a'),getFreshVariable(freeVars),false)
			}
			tmp_node = tmp_node.parent
		}
	}
}

//get variables bounded to node, change the value of them. then changes its own value
func alphaConversion(theNode *node, oldVar rune, freshVar string, flag bool) {
/*
	if (theNode == nil) {
		return
	}
	local_flag := flag
	nodeVar := []rune(theNode.value)
 //maybe can reuse get captured variables
	if nodeVar[0] == oldVar{
		switch theNode.token {
		case LAMBDA:
			if flag{
				return
			}
			//else
			local_flag = true
			fallthrough
		case VARIABLE:
			theNode.value = string(freshVar)
		}
	}

	alphaConversion(theNode.left, oldVar, freshVar,local_flag)
	alphaConversion(theNode.right, oldVar, freshVar,local_flag)
 */
	for _,capturedNode:= range getCapturedNodes(theNode){
		capturedNode.value = freshVar
	}
	theNode.value = freshVar
}

func getFreshVariable(freeVars []string) string{
	max_size := 100
	char_array := []rune{}

	i:=0
	for i < max_size{
		char_array = append(char_array,rune(int('a')-1) )
		j:=0
		for j < 26{
			char_array[i]++
			if !isPresent(string(char_array),freeVars){
				return string(char_array)
			}
		}
	}
	return "#"
}
//*************************************************************************

// substitutes the variables bound to the original lambda expression.
// newNode is new independent node
// not sure if it would work with only one right-hand side tree that's being pointed by the parents of every bound variable
func substituteTree(theNode *node, subNode *node, targetVar string) {
	if (theNode == nil) {
		return
	}

	//nodeVar := []rune(theNode.value)
	if (theNode.token == LAMBDA && theNode.value == targetVar) {
		return
	}

	if (theNode.token == VARIABLE && theNode.value == targetVar) {
		//var newNode *node
		//newNode = subNode
		newNode := getCopySubtree(subNode)
		parentNode := theNode.parent
		if (parentNode.left == theNode) {
			parentNode.left = newNode
		} else if (parentNode.right == theNode) {
			parentNode.right = newNode
		} else{
			print ("SHOULDN'T BE HERE\n")
			os.Exit(1)
		}
		newNode.parent = parentNode
		//theNode = newNode
	}

	substituteTree(theNode.left, subNode, targetVar)
	substituteTree(theNode.right, subNode, targetVar)
}

//get an exact copy of subtree. but every node is created again = indipendent
func getCopySubtree(subtree *node) *node{
	//base case
	if subtree == nil{
		return nil
	}
	//Create a new identical node and then link a copy of the left and the right of the subtree
	returnNode := newNode(subtree.value,subtree.token)
	returnNode.linkNodes(getCopySubtree(subtree.left) , getCopySubtree(subtree.right))
	return returnNode
}


// Applies beta-reduction once to the first applicable branch with preference
// to the left-hand side.
// might be working with single pointer
func applyBetaReduc(theNode **node) bool {
	if ((*theNode) == nil || (*theNode).token == VARIABLE) {
		return false
	}
	//Reduction possible
	if ((*theNode).token == APPLICATION && (*theNode).left.token == LAMBDA && (*theNode).right != nil) {
		//check the need for alpha conversion

		freeVars := []string{}
		capturedNodes := getCapturedNodes( (*theNode).left )
		giveFreeVars((*theNode).right, &freeVars)
		alphaDriver(capturedNodes,freeVars)

		//actual substitution
		targetVar := ((*theNode).left.value)
		//targetVar := targetSplice[0]
		substituteTree((*theNode).left.left, (*theNode).right, targetVar)	

		//link the sub tree of the lambda node to the its parent
		(*theNode).left.left.parent = (*theNode).parent
		//change the value stored at the address of the lambda node
		**theNode = *((*theNode).left.left)

		return true
	}else{
		//Reduction is not possible. Check left branch. If no reduction found there, check right side
		if applyBetaReduc( (&(*theNode).left) ) == false {
			return applyBetaReduc((&(*theNode).right))
		} else{
			return true
		}
	}

	return false
}

// this function doesnt really have a functionality yet.
// I am merely using it to test functions with.
func checkReduction (theNode **node) bool {
	/*
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
	*/
	if applyBetaReduc(&(*theNode)) {
		fmt.Fprintf(os.Stdout, "apply work\n")
	}
	return true
} 

//*************************************************************************


//*************************************************************************

func printPostOrder(theNode *node) {
	if (theNode == nil) {
		return
	}
	if (theNode.token == LAMBDA) {
		fmt.Fprintf(os.Stdout, "(")
		fmt.Fprintf(os.Stdout, "λ%s ", theNode.value)
	} else if (theNode.token == VARIABLE) {
		fmt.Fprintf(os.Stdout, "%s", theNode.value)
	}

	if (theNode.left != nil && theNode.right != nil) {
		fmt.Fprintf(os.Stdout, "(")
	}
	printPostOrder(theNode.left)

	if (theNode.left != nil && theNode.right != nil) &&
		(theNode.left.token == VARIABLE && theNode.right.token == VARIABLE) {
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