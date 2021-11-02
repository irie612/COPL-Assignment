YEAR: 2021-2022

COURSE: CONCEPTS OF PROGRAMMING LANGUAGES (ASSIGNMENT 1)

CLASS: 2

GROUP: 11

AUTHOR(S): Emanuele Greco (s3375951), Irie Railton (s3292037) & Kah Ming wong (s2641976)
OPERATING SYSTEM(S): Ubuntu 20.04

The program works correctly on most parts. Aspects that may be called deficiencies are:
- The program is not able to correctly parse an expression with dots as part of the lambda abstractions

## Reduction strategy
A)
It checks if a reduction is possible on the root of the tree through the use 
of the betaReduction() function.
If not, we recursively call the function first on the left child node and then, 
if no reduction has been found there, on the right node.
If no reduction has been found, we exit

B)
It checks for reductions starting from the top of the tree.
if two reductions are possible on the same level of the tree, the left one is preferred
