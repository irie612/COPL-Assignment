package main

func typeParse() *node {
	if nextToken == VARIABLE {
		tmp_node := newNode(string(lexeme[:lexLen]), VARIABLE)
		lex()
		if nextToken == ARROW {
			lex()
			arrow_node := newNode("☛", ARROW)
			arrow_node.linkNodes(tmp_node, typeParse())

			return arrow_node
		} else {
			//print("HELP!")
			if nextToken == RIGHT_P {
				lex()
			}
			return tmp_node
		}
	} else if nextToken == LEFT_P {
		lex()
		return typeParse()
	}
	return nil
}
