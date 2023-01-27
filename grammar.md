program        → declaration* EOF ;

declaration    → varDecl | statement ;
varDecl        → "var" IDENTIFIER ("=" expression) ;
statement      → exprStmt | printStmt | block | ifStmt | forStmt ;
forStmt        → "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ; 
ifStmt         → "if" "(" expression ")" statement ("else" statement)?
block          → "{" declaration* "}"
exprStmt       → expression ";" ;
printStmt      → "print" expression ";" ;
expression     → assignment ;
assignment     → IDENTIFIER "=" assignment | logical_or ;
logical_or     → logical_and ( "or" logical_and )* ;
logical_and    → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER;