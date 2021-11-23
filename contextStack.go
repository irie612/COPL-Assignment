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
func (cs *contextStack) findStatement(varName string, typeNode *node) (bool, *node) {
	//using indirect as an index, traverse the stack until the end is reached
	for indirect := cs.head; indirect != nil; indirect = indirect.left {
		//if we find a statement with the same variable name we check for the type
		if indirect.value == varName {
			if compareSubtrees(indirect.right, typeNode) {
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
	//using indirect as an index, traverse the stack until the end is reached
	for indirect := cs.head; indirect != nil; indirect = indirect.left {
		//if we find a statement with the same variable name we return the type
		if indirect.value == varName {
			return indirect.right
		}
	}

	return nil
}

//*************************************************************************

//get an indentical copy of the stack
func (cs *contextStack) getCopy() contextStack {
	//crete a contextStack with the same underlying tree as the one is called on
	//can't use getCopySubtree(cs) it requires a *node type as argument
	return contextStack{getCopySubtree(cs.head)}
}

//*************************************************************************

func compareSubtrees(NodeA *node, NodeB *node) bool {
	//Base case. If here, then all the comparisons went well
	if NodeA == nil && NodeB == nil {
		return true
	}
	//recursion
	if NodeA.token == NodeB.token && (NodeA.value == NodeB.value || 
		NodeA.value == "?" || NodeB.value == "?") && !(NodeA.value == "?" && NodeB.value == "?") {
		return compareSubtrees(NodeA.left, NodeB.left) && compareSubtrees(NodeA.right, NodeB.right)
	} 

	return false
}
