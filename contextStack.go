package main

type contextStack struct {
	head *node
}

//*************************************************************************

func (cs *contextStack) addStatement(varName string, typeNode *node)  {
	//create a node that will be inserted in the context
	n := newNode(varName,VARIABLE)
	//add n on top of the current head, add type on the right of the node
	n.linkNodes(cs.head,typeNode)
	//change pointer to head
	cs.head=n
}

//*************************************************************************

// check for the presence of a statement
func (cs contextStack) findStatement(varName string, typeNode *node) bool{
	//using indirect as an index, traverse the stack until the end is reached
	for indirect := cs.head; indirect!=nil; indirect=indirect.left {
		//if we find a statement with the same variable name we check for the type
		if indirect.value == varName{
			return compareSubtrees(indirect.right,typeNode)
		}
	}
	return false
}

//*************************************************************************

func compareSubtrees(NodeA *node,NodeB *node) bool{
	//Base case. If here, all the comparisons went well
	if NodeA == nil && NodeB == nil{
		return true
	}
	//recursion
	if NodeA.token==NodeB.token && NodeA.value==NodeB.value{
		return compareSubtrees(NodeA.left,NodeB.left) && compareSubtrees(NodeA.right,NodeB.right)
	} else{
		return false
	}
}