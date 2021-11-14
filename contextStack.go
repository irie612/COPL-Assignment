package main

type contextStack struct {
	head *node
}

func ( cs contextStack ) addStatement(varName string, typeNode *node)  {
	n := newNode(varName,VARIABLE)
	//add n on top of the current head, add type on the right of the node
	n.linkNodes(cs.head,typeNode)
	//change pointer to head
	cs.head=n
}

func (cs contextStack) findVariable(varName string, typeNode *node) bool{


	for indirect := cs.head; indirect!=nil; indirect=indirect.left {
		if indirect.value == varName{
			return compareSubtrees(indirect.right,typeNode)
		}
	}
}


func compareSubtrees(NodeA *node,NodeB *node) bool{
	if NodeA == nil && NodeA == nil{
		return true
	}
	if NodeA.token==NodeB.token && NodeA.value==NodeB.value{
		return compareSubtrees(NodeA.left,NodeB.left) && compareSubtrees(NodeA.right,NodeB.right)
	} else{
		return false
	}
}