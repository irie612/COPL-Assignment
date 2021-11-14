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
//It has a bug Irie porcodio !!!!!!"
func (n *node) toString() string {
	var returnString string
	if n == nil {
		return ""
	}
	if n.token == LAMBDA {
		returnString = "λ" + n.value + "^" +
			n.right.toString() + " " +
			n.left.toString()
	} else if n.token == VARIABLE {
		returnString = n.value
	} else if n.token == APPLICATION {
		if n.left.token == VARIABLE {
			returnString = n.left.toString()
		} else {
			returnString = "(" + n.left.toString() + ")"
		}
		if n.left.token == VARIABLE &&
			n.right.token == VARIABLE {
			returnString += " " + n.right.toString()
		} else if n.right.token != VARIABLE {
			returnString += "(" + n.right.toString() + ")"
		} else {
			returnString += n.right.toString()
		}
	} else if n.token == ARROW {
		returnString = "(" + n.left.toString() + "→" +
			n.right.toString() + ")"
	}
	return returnString
}
