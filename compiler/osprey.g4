grammar osprey;

// ---------- PARSER RULES ----------

program         : statement* EOF ;

statement
    : importStmt
    | letDecl
    | assignStmt
    | fnDecl
    | externDecl
    | typeDecl
    | effectDecl
    | moduleDecl
    | exprStmt
    ;

importStmt      : IMPORT ID (DOT ID)* ;

letDecl         : (LET | MUT) ID (COLON type)? EQ expr ;

assignStmt       : ID EQ expr ;

fnDecl          : docComment? FN ID LPAREN paramList? RPAREN (ARROW type)? effectSet? (EQ expr | LBRACE blockBody RBRACE) ;

externDecl      : docComment? EXTERN FN ID LPAREN externParamList? RPAREN (ARROW type)? ;

externParamList : externParam (COMMA externParam)* ;

externParam     : ID COLON type ;

paramList       : param (COMMA param)* ;

param           : ID (COLON type)? ;

typeDecl        : docComment? TYPE ID (LT typeParamList GT)? EQ (unionType | recordType) typeValidation? ;

typeParamList   : ID (COMMA ID)* ;

unionType       : variant (BAR variant)* ;

recordType      : LBRACE fieldDeclarations RBRACE ;

variant         : ID (LBRACE fieldDeclarations RBRACE)? ;

fieldDeclarations : fieldDeclaration (COMMA fieldDeclaration)* ;
fieldDeclaration  : ID COLON type (WHERE functionCall)? ;

typeValidation  : WHERE ID ;

// Effect declarations
effectDecl      : docComment? EFFECT ID LBRACE opDecl* RBRACE ;
opDecl          : ID COLON type ;

// Effect sets for function types
effectSet       : NOT_OP ID                              // Single effect: !Effect
                | NOT_OP LSQUARE effectList RSQUARE      // Multiple effects: ![Effect1, Effect2]
                ;

effectList      : ID (COMMA ID)* ;

// Handler expressions - implementing spec syntax
handlerExpr     : HANDLE ID handlerArm+ IN expr ;              // handle Logger log msg => ... in expr
handlerArm      : ID handlerParams? LAMBDA expr ;
handlerParams   : ID+ ;

functionCall    : ID LPAREN argList? RPAREN ;

booleanExpr     : comparisonExpr ;

fieldList       : field (COMMA field)* ;
field           : ID COLON type ;

type            : LPAREN typeList? RPAREN ARROW type  // Function types like (Int, String) -> Bool
                | FN LPAREN typeList? RPAREN ARROW type  // Function types like fn(Int, String) -> Bool
                | ID (LT typeList GT)?  // Generic types like Result<String, Error>
                | ID LSQUARE type RSQUARE  // Array types like [String]
                | ID ;

typeList        : type (COMMA type)* ;

exprStmt        : expr ;

expr
    : matchExpr
    ;

matchExpr
    : MATCH binaryExpr LBRACE matchArm+ RBRACE
    | selectExpr
    | binaryExpr
    ;



selectExpr
    : SELECT LBRACE selectArm+ RBRACE
    ;

selectArm
    : pattern LAMBDA expr                         // pattern => expr
    | UNDERSCORE LAMBDA expr                      // _ => expr (default)
    ;

binaryExpr
    : ternaryExpr
    ;

ternaryExpr
    : cond=comparisonExpr LBRACE pat=fieldPattern RBRACE QUESTION thenExpr=ternaryExpr COLON elseExpr=ternaryExpr    // expr { pattern } ? then : else
    | comparisonExpr QUESTION thenExpr=ternaryExpr COLON elseExpr=ternaryExpr                                        // expr ? then : else
    | comparisonExpr QUESTION COLON elseExpr=ternaryExpr                                                             // expr ?: else (Elvis operator)
    | comparisonExpr
    ;



comparisonExpr
    : addExpr ((EQ_OP | NE_OP | LT | GT | LE_OP | GE_OP) addExpr)*
    ;

addExpr
    : mulExpr ((PLUS | MINUS) mulExpr)*
    ;

mulExpr
    : unaryExpr ((STAR | SLASH | MOD_OP) unaryExpr)*
    ;

unaryExpr
    : (PLUS | MINUS | NOT_OP | AWAIT)? pipeExpr
    ;

pipeExpr
    : callExpr (PIPE callExpr)*
    ;

callExpr
    : primary (DOT ID)+ (LPAREN argList? RPAREN)?  // Field access with optional final method call: obj.field or obj.field.method()
    | primary (DOT ID (LPAREN argList? RPAREN))+   // Method chaining: obj.method().chain() (at least one method call)
    | primary (LPAREN argList? RPAREN)?            // Function call with optional parentheses
    ;

argList
    : namedArgList                                 // Named arguments (for multi-param functions)
    | expr (COMMA expr)*                          // Traditional positional arguments
    ;

namedArgList
    : namedArg (COMMA namedArg)+                  // At least 2 named args
    ;

namedArg
    : ID COLON expr                               // paramName: value
    ;

primary
    : SPAWN expr                                  // spawn expr
    | YIELD expr?                                 // yield or yield expr
    | AWAIT LPAREN expr RPAREN                    // await(fiber) - function call style
    | SEND LPAREN expr COMMA expr RPAREN          // send(channel, value)
    | RECV LPAREN expr RPAREN                     // recv(channel)
    | SELECT selectExpr                           // select { ... }
    | PERFORM ID DOT ID LPAREN argList? RPAREN    // perform EffectName.operation(args)
    | handlerExpr                                 // handle EffectName ... in expr
    | typeConstructor                             // Type construction (Fiber<T> { ... })
    | updateExpr                                  // Non-destructive update (record { field: newValue })
    | objectLiteral                               // Anonymous object literal { field: value }
    | blockExpr                                   // Block expressions
    | literal                                     // String, number, boolean literals
    | lambdaExpr                                  // Lambda expressions
    | ID LSQUARE INT RSQUARE                      // List access: list[0] -> Result<T, IndexError>
    | ID                                          // Variable reference
    | LPAREN expr RPAREN                          // Parenthesized expression
    ;

// Anonymous object literal: { field: value, field2: value2 }
objectLiteral
    : LBRACE fieldAssignments RBRACE
    ;

// Type construction for Fiber<T> { ... } and Channel<T> { ... }
typeConstructor
    : ID typeArgs? LBRACE fieldAssignments RBRACE
    ;

typeArgs
    : LT typeList GT
    ;

fieldAssignments
    : fieldAssignment (COMMA fieldAssignment)*
    ;

fieldAssignment
    : ID COLON expr
    ;

lambdaExpr
    : FN LPAREN paramList? RPAREN (ARROW type)? LAMBDA expr      // fn(x, y) => x + y
    | BAR paramList? BAR LAMBDA expr               // |x, y| => x + y (short syntax)
    ;

// Non-destructive update: record { field: newValue }
updateExpr
    : ID LBRACE fieldAssignments RBRACE
    ;

// Block expressions for local scope and sequential execution
blockExpr
    : LBRACE blockBody RBRACE
    ;

literal
    : INT
    | STRING
    | INTERPOLATED_STRING
    | TRUE
    | FALSE
    | listLiteral
    ;

listLiteral
    : LSQUARE (expr (COMMA expr)*)? RSQUARE ;

docComment      : DOC_COMMENT+ ;

moduleDecl      : docComment? MODULE ID LBRACE moduleBody RBRACE ;

moduleBody      : moduleStatement* ;

moduleStatement : letDecl | fnDecl | typeDecl ;

matchArm
    : pattern LAMBDA expr ;

pattern
    : unaryExpr                                   // Support negative numbers: -1, +42, etc.
    | ID (LBRACE fieldPattern RBRACE)?          // Pattern destructuring: Ok { value }
    | ID (LPAREN pattern (COMMA pattern)* RPAREN)?  // Constructor patterns
    | ID (ID)?                                   // Variable capture
    | ID COLON type                              // Type annotation pattern: value: Int
    | ID COLON LBRACE fieldPattern RBRACE       // Named structural: person: { name, age }
    | LBRACE fieldPattern RBRACE                // Anonymous structural: { name, age }
    | UNDERSCORE                                 // Wildcard
    ;

fieldPattern    : ID (COMMA ID)* ;

blockBody       : statement* expr? ;

// ---------- LEXER RULES ----------
// Keywords MUST come before ID to ensure proper tokenization

// Control flow keywords
MATCH       : 'match';
IF          : 'if';
ELSE        : 'else';
SELECT      : 'select';

// Function keywords  
FN          : 'fn';
EXTERN      : 'extern';

// Declaration keywords
IMPORT      : 'import';
TYPE        : 'type';
MODULE      : 'module';
LET         : 'let';
MUT         : 'mut';

// Effect system keywords - CRITICAL ORDER FOR PROPER TOKENIZATION
EFFECT      : 'effect';
PERFORM     : 'perform';
HANDLE      : 'handle';
IN          : 'in';
DO          : 'do';

// Concurrency keywords
SPAWN       : 'spawn';
YIELD       : 'yield';
AWAIT       : 'await';
FIBER       : 'fiber';
CHANNEL     : 'channel';
SEND        : 'send';
RECV        : 'recv';

// Boolean keywords
TRUE        : 'true';
FALSE       : 'false';

// Constraint keyword
WHERE       : 'where';

// Operators and symbols
PIPE        : '|>';
ARROW       : '->';
LAMBDA      : '=>';
UNDERSCORE  : '_';

EQ          : '=';
EQ_OP       : '==';
NE_OP       : '!=';
LE_OP       : '<=';
GE_OP       : '>=';
NOT_OP      : '!';
MOD_OP      : '%';
COLON       : ':';
SEMI        : ';';
COMMA       : ',';
DOT         : '.';
BAR         : '|';
LT          : '<';
GT          : '>';
LPAREN      : '(';
RPAREN      : ')';
LBRACE      : '{';
RBRACE      : '}';
LSQUARE     : '[';
RSQUARE     : ']';
QUESTION    : '?';

PLUS        : '+';
MINUS       : '-';
STAR        : '*';
SLASH       : '/';

// Literals and identifiers - MUST come after keywords
INT         : [0-9]+ ;
INTERPOLATED_STRING : '"' (~["\\$] | '\\' . | '$' ~[{])* ('${' ~[}]* '}' (~["\\$] | '\\' . | '$' ~[{])*)+ '"' ;
STRING      : '"' (~["\\] | '\\' .)* '"' ;
ID          : [a-zA-Z_][a-zA-Z0-9_]* ;

// Whitespace and comments - MUST be at the end
WS          : [ \t\r\n]+ -> skip ;
DOC_COMMENT : '///' ~[\r\n]* ;
COMMENT     : '//' ~[\r\n]* -> skip ;
