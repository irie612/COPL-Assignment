// AST.go
// Programming Language: GoLang
//
// Course: Concepts of Programming Language
// Assignment 2: Interpreter
// Class 2, Group 11
// Author(s) :	Emanuele Greco (s3375951), 
//							Irie Railton (s3292037),
//							Kah ming Wong (s2641976).
//
// Date: 3rd November, 2021.
// 

//*************************************************************************

package main

//*************************************************************************

// AST Node Struct
type node struct {
	parent			*node
	left				*node
	right				*node
	value				string
	token				int //type of node
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