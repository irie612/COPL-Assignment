YEAR: 2021-2022

COURSE: CONCEPTS OF PROGRAMMING LANGUAGES (ASSIGNMENT 3)

CLASS: 2

GROUP: 11

AUTHOR(S): Emanuele Greco (s3375951), Irie Railton (s3292037), & Kah Ming Wong (s2641976)

GoLang Compiler Version: go 1.13

Operating System (Emanuele): Ubuntu 20.04

Operating System (Irie): Ubuntu 20.04.1 LTS

Operating System (Kah Ming): Ubuntu 20.04.3 LTS



The program works correctly on all parts. ~~Aspects that may be called deficiencies are:~~
- ~~The program is not able to correctly parse an expression with dots as part of the lambda abstractions~~
    The program is now able to correctly parse an expression with dots.

______________________________________________________________________________________________________________________________________________
## Assignment 1

The program is a parser for LL grammar using recursive descent. It starts in the main function, which can be seen as the main driver of the whole program. Here, we first check whether an argument is given which contains an acceptable text file. If so, then the program will parse the text file until the EOF. The parsing begins in the parse function. This function first performs lexical analysis onto first character of the line, with the help of the lex() function. The lex() function essentially determines what the token is from the first non-blank character: variable, digit, or UNKNOWN. If it's unknown then another function in the lex() function is called, to resolve what token it is, based on its character. 

The parse function will call the expr() function after the program completes the lexical analysis for a specific character. The expr function calls both lexpr(), and the expr_p() function. In the lexpr() function we try to look for a valid lambda abstraction and add it to the output string. In the instance the lambda construction is invalid, the program exits with status 1 (syntax error). If a lambda abstraction was not found in lexpr(), then the program goes to the pexpr() function. Here, we determine whether the token is a left parenthesis. If it is a left parenthesis, then the next tokens should be a variable, and followed with a right parenthesis. If either does not hold, then it again exits with status 1. If at pexpr() we do not encounter a left parenthesis first, then we can determine that the token is a variable, and thus, add it to the output string.

After the variable has been added to the output string within the variable case of pexpr(), if the next token is a variable or a left parenthesis, it will append a right parenthesis to the output string. The expr() function also keeps track of the starting position of the current part of the expression and at the end of the function it calls a helper function that counts, from the given start position, the number of each type of parenthesis and adds the necessary number of left parenthesis.

The expr_p() function does this whole process again, in the instance that there are more valid expressions. However, the expr_p() function does not necessarily need to add anything, and can essentially be "empty". If called, the program would return to the previous operation in the stack.

After every token has been analyzed, and the line has been parsed, the parse() function outputs the string. If the token does not equal to either EOL or EOF after every valid expression has been analyzed and parsed, then the parse() function will exit with status 1 (string not fully parsed). 

If by then EOF has not been reached, then the main driver would parse the next line with the exact same steps.


______________________________________________________________________________________________________________________________________________
## Assignment 2

Assignment 2 builds onto our previous iteration of lexer and parser. That is, with the addition of creating nodes whilst traversing between the recursive descents of our lexer and parser. These nodes result into an AST tree which we can apply alpha-conversion and beta-reductions on (interpreter). To do that, we simply call the `betaDriver(rootNode)` function in our old `parse()` function. `betaDriver()` is a function that calls `betaReduction()` function. `betaReduction()` is a function that looks for a valid beta-reduction in the tree, with preference to the left-hand side; and applies beta-reduction once per call. Additionally, `betaReduction()` calls `alphaDriver()` which checks whether alpha-conversion is needed, and if so, applies alpha-conversion onto the specific branch. `betaReduction()` will return a true to `betaDriver()` if it found and applied a beta-reduction, else a false (no more valid beta-reductions in the tree). Thus, we can keep looping `betaReduction()` in `betaDriver()` until it gives a false, or, until the maximum amount of reductions has been reached. See below for our reduction strategy.

## Reduction strategy (in the perspective of `betaReduction()`)

1. Is the "current node" a valid node to apply beta-reduction on? Yes= go to step 2, No= go to step 3.
2. Is alpha-conversion required? If so, apply alpha-conversion and then beta-reduction, if not, immediately apply beta-reduction. Return true.
3. Move down the tree, starting with the left child. And check whether the current node is a valid node for beta-reduction (step 1). We keep moving down the tree through the left nodes until we hit a NULL (or there was a valid beta-reduction). If we hit a NULL `betaReudction()` will return a false (to the previous recursive calls). And move to step 4. If the call of `betaReduction()` with the left child resulted into a true, then return true. Otherwise, also move to step 4.
4. If the branch under the left child of the "current node" has no valid beta-reduction, then go down the tree through the right child of "current node". And go back to step 1. At this point, the "current node" will return the value of the `betaReduction()` call of the right child.


______________________________________________________________________________________________________________________________________________
## Assignment 3

Assigment 3 is the last addition to our existing program, which involves 'type checking'. For this, we slightly modified our existing parser, and furthermore, created new functions to parse types with. The end result of one call of parse on a correct judgement should give us a valid `rootExpressionNode` containing the expression as a tree, and a valid `rootTypeNode` containing the type of the judgement as a tree. If either of these nodes are invalid, from invalid input or grammar, then the program will terminate with exit status 1. Additionally, types in the expression part of the judgement are now to be stored in the right child of lambda nodes, which, previously, were empty.


Our type check strategy is by inference, as explained in the work groups. Hence, it is not necessarily a standard one-to-one type check strategy resulting from the three rules; but a derivation from it, such that we go in the _reverse_ of the rules. Thus:
- If the expression is a variable: Return the type in the context (since x: T). Otherwise, return nil if the variable is not present in the context.
- If the expression is a lambda (+body) (= λx^T E): Add the type of the lambda-variable (the T in λx^T E), and proceed the inference (of E) to get the T' of the expression (since E:T'). If the inference for T' results into a nil, then we stop the inference, and declare that the expression cannot be type checked; because the type of the expression HAS to be of the structure `T→T'` (since λx^T E: T→T'). If we do get a T', then we make a tree that is essentially of the structure `T→T'`. In other words, the head is an arrow, the left child is the T from the lambda-variable, and the right child is the T' that we just inferred. This will also be our return value.
- If the expression is an application expression (= E1E2): We pass both E1, and E2 through the inference function to get their respective types. If either happens to be nil, then we declare that the expression cannot be type checked. Otherwise, we will compare the T of E1 (since E1: T→T'), and the T of E2 (since E2: T). In the case when the T's are exactly the same, then we return T' (since (E1 E2): T').

In simply typed lambda calculus, type inference is decidable and if there is a type it is unique. This gives rise to the ability to either infer the type, or decide that the expression does not have a valid type. 
So after we have used the expression to infer the type we only have to compare this to the type that is in the judgment.
If these are not equivalent, then we can declare that the judgement does not type check. Otherwise, we can declare that the judgement correctly type checks.
