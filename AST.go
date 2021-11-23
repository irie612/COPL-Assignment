// AST.go
// Programming Language: GoLang
//
// Course: Concepts of Programming Language
// Assignment 3: Type checker
// Class 2, Group 11
// Author(s) :	Emanuele Greco (s3375951),
//				Irie Railton (s3292037),
//				Kah ming Wong (s2641976).
//
// Date: 13th November, 2021.
//

//*************************************************************************

package main

import(
	"os"
	"fmt"
)

//*************************************************************************

// AST Node Struct
type node struct {
	parent *node
	left   *node
	right  *node  //used to store the type of the node for Lambda abstraction nodes and for variables in contextStack
	value  string
	token  int   //type of node
}

//*************************************************************************

func newNode(value string, token int) *node {
	return &node{
		value: value,
		token: token,
	}
}

//*************************************************************************

// Function to link a node to its children.
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

//function to convert a tree into a string recursively
func (n *node) toString() string {
	var returnString string
	if n == nil {
		return ""
	}
	if n.token == LAMBDA {
		returnString = "λ" + n.value + "^"
		if n.right.token == ARROW {
			returnString += bracket(n.right.toString())
		} else {
			returnString += n.right.toString()
		}
		if n.left.token == APPLICATION || n.left.token == LAMBDA {
			returnString += bracket(n.left.toString())
		} else {
			returnString += " " + n.left.toString()
		}
	} else if n.token == VARIABLE {
		returnString = n.value
	} else if n.token == APPLICATION {
		if n.left.token == VARIABLE {
			returnString = n.left.toString()
		} else {
			returnString = bracket(n.left.toString())
		}
		if n.left.token == VARIABLE &&
			n.right.token == VARIABLE {
			returnString += " " + n.right.toString()
		} else if n.right.token != VARIABLE {
			returnString += bracket(n.right.toString())
		} else {
			returnString += n.right.toString()
		}
	} else if n.token == ARROW {
		if n.parent != nil {
			returnString = bracket(n.left.toString() + "→" +
				n.right.toString())
		} else {
			returnString = n.left.toString() + "→" +
				n.right.toString()
		}

	}
	return returnString
}

func bracket(s string) string {
	return "(" + s + ")"
}

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
		(theNode.left.token == VARIABLE && theNode.right.token == VARIABLE) && theNode.token != ARROW {
		fmt.Fprintf(os.Stdout, " ")
	} 

	if theNode.token == ARROW {
		fmt.Fprintf(os.Stdout, "->")
	}

	printPostOrder(theNode.right)
	if theNode.token == LAMBDA || 
	(theNode.left != nil && theNode.right != nil) {
		fmt.Fprintf(os.Stdout, ")")
	}
}
