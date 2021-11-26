// contextStack.go
// Programming Language: GoLang
//
// Course: Concepts of Programming Language
// Assignment 3: Type Checking
// Class 2, Group 11
// Author(s) :	Emanuele Greco (s3375951),
//				Irie Railton (s3292037),
//				Kah ming Wong (s2641976).
//
// Date: 15th November, 2021.
//

//*************************************************************************

package main

type contextStack struct {
	head *node
}

//*************************************************************************

func (cs *contextStack) addStatement(varName string, typeNode *node) {
	//create a node that will be inserted in the context
	n := newNode(varName, VARIABLE)
	//add n on top of the current head, add type on the right of the node
	n.linkNodes(cs.head, getCopySubtree(typeNode))
	//change pointer to head
	cs.head = n
}

//*************************************************************************

// check for the presence of a statement
func (cs *contextStack) findStatement(value string, n *node) (bool, *node) {
	//using indirect as an index, traverse the stack
	for indirect := cs.head; indirect != nil; indirect = indirect.left {
		//if statement with same variable name found check for its type
		if indirect.value == value {
			if indirect.right.compareSubtrees(n) {
				return true, indirect.right
			}
			return false, nil
		}
	}

	return false, nil
}

//*************************************************************************

// Given a variable name, gives back a type corresponding to the one
//in the first statement in the stack. (Used in type inference)
func (cs *contextStack) getType(varName string) *node {
	//using indirect as an index, traverse the stack
	for indirect := cs.head; indirect != nil; indirect = indirect.left {
		//if statement with same variable name found return its type
		if indirect.value == varName {
			return indirect.right
		}
	}

	return nil
}

//*************************************************************************

//get an identical copy of the stack
func (cs *contextStack) getCopy() contextStack {
	//can't use getCopySubtree(cs) it requires a *node type as argument
	return contextStack{getCopySubtree(cs.head)}
}

//*************************************************************************
