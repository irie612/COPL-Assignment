YEAR: 2021-2022

COURSE: CONCEPTS OF PROGRAMMING LANGUAGES (ASSIGNMENT 1)

CLASS: 2

GROUP: 11

AUTHOR(S): Emanuele Greco (s3375951), Irie Railton (s3292037) & Kah Ming wong (s2641976)



The program works correctly on most parts. Aspects that may be called deficiencies are:
- The output does not contain the least amount of brackets.

There are no deviations from the assignment.


______________________________________________________________________________________________________________________________________________


The program is a parser for LL grammar using recursive descent. It starts in the main function, which can be seen as the main driver of the whole program. Here, we first check whether an argument is given which contains an acceptable text file. If so, then the program will parse the text file until the EOF. The parsing begins in the parse function. This function first performs lexical analysis onto first character of the line, with the help of the lex() function. The lex() function essentially determines what the token is from the first non-blank character: variable, digit, or UNKNOWN. If it's unknown then another function in the lex() function is called, to resolve what token it is, based on its character. 

The parse function will call the expr() function after the program completes the lexical analysis for a specific character. The expr function calls both lexpr(), and the expr_p() function. In the lexpr() function we try to look for a valid lambda abstraction and add it to the output string. In the instance the lambda construction is invalid, the program exits with status 1 (syntax error). If a lambda abstraction was not found in lexpr(), then the program goes to the pexpr() function. Here, we determine whether the token is a left parenthesis. If it is a left parenthesis, then the next tokens should be a variable, and followed with a right parenthesis. If either does not hold, then it again exits with status 1. If at pexpr() we do not encounter a left parenthesis first, then we can determine that the token is a variable, and thus, add it to the output string.

After the variable has been added to the output string within the variable case of pexpr(), if the next token is a variable or a left parenthesis, it will append a right parenthesis to the output string. The expr() function also keeps track of the starting position of the current part of the expression and at the end of the function it calls a helper function that counts, from the given start position, the number of each type of parenthesis and adds the necessary number of left parenthesis.

The expr_p() function does this whole process again, in the instance that there are more valid expressions. However, the expr_p() function does not necessarily need to add anything, and can essentially be "empty". If called, the program would return to the previous operation in the stack.

After every token has been analyzed, and the line has been parsed, the parse() function outputs the string. If the token does not equal to either EOL or EOF after every valid expression has been analyzed and parsed, then the parse() function will exit with status 1 (string not fully parsed). 

If by then EOF has not been reached, then the main driver would parse the next line with the exact same steps.
