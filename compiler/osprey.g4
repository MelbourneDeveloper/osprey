grammar osprey;

// ---------- PARSER RULES ----------

program         : statement* EOF ;

statement
    : importStmt
    | letDecl
    | fnDecl
    | externDecl
    | typeDecl
    | moduleDecl
    | exprStmt
    ;

importStmt      : IMPORT ID (DOT ID)* ;

letDecl         : (LET | MUT) ID (COLON type)? EQ expr ;

fnDecl          : docComment? FN pluginName ID LPAREN paramList? RPAREN pluginReturnType? EQ expr        // Plugin function  
                | docComment? FN ID LPAREN paramList? RPAREN (ARROW type)? (EQ expr | LBRACE blockBody RBRACE)    // Regular function
                ;

pluginName      : ID ;

pluginReturnType : ARROW type          // -> Type (existing type)
                 | AS type             // as Type (generate type)
                 ;

externDecl      : docComment? EXTERN FN ID LPAREN externParamList? RPAREN (ARROW type)? ;

externParamList : externParam (COMMA externParam)* ;

externParam     : ID COLON type ;

paramList       : param (COMMA param)* ;

param           : ID (COLON type)? ;

typeDecl        : docComment? TYPE ID (LT typeParamList GT)? EQ (unionType | recordType) ;

typeParamList   : ID (COMMA ID)* ;

unionType       : variant (BAR variant)* ;

recordType      : LBRACE fieldDeclarations RBRACE ;

variant         : ID (LBRACE fieldDeclarations RBRACE)? ;

fieldDeclarations : fieldDeclaration (COMMA fieldDeclaration)* ;
fieldDeclaration  : ID COLON type constraint? ;

constraint      : WHERE functionCall ;

functionCall    : ID LPAREN argList? RPAREN ;

booleanExpr     : comparisonExpr ;

fieldList       : field (COMMA field)* ;
field           : ID COLON type ;

type            : ID (LT typeList GT)? LSQUARE RSQUARE  // Array types like Type[]
                | ID (LT typeList GT)?                     // Generic types like Result<String, Error>
                | ID LSQUARE type RSQUARE                  // Legacy array types like [String]
                ;

typeList        : type (COMMA type)* ;

exprStmt        : expr ;

expr
    : matchExpr
    ;

matchExpr
    : MATCH expr LBRACE matchArm+ RBRACE
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
    : comparisonExpr
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
    | typeConstructor                             // Type construction (Fiber<T> { ... })
    | updateExpr                                  // Non-destructive update (record { field: newValue })
    | blockExpr                                   // Block expressions
    | literal                                     // String, number, boolean literals
    | lambdaExpr                                  // Lambda expressions
    | ID                                          // Variable reference
    | LPAREN expr RPAREN                          // Parenthesized expression
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
    | FALSE ;

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

PIPE        : '|>';
MATCH       : 'match';
FN          : 'fn';
EXTERN      : 'extern';
IMPORT      : 'import';
TYPE        : 'type';
MODULE      : 'module';
LET         : 'let';
MUT         : 'mut';
IF          : 'if';
ELSE        : 'else';
LOOP        : 'loop';
SPAWN       : 'spawn';
YIELD       : 'yield';
AWAIT       : 'await';
FIBER       : 'fiber';
CHANNEL     : 'channel';
SEND        : 'send';
RECV        : 'recv';
SELECT      : 'select';
TRUE        : 'true';
FALSE       : 'false';
WHERE       : 'where';
AS          : 'as';

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

PLUS        : '+';
MINUS       : '-';
STAR        : '*';
SLASH       : '/';

INT         : [0-9]+ ;
INTERPOLATED_STRING : '"' (~["\\$] | '\\' . | '$' ~[{])* ('${' ~[}]* '}' (~["\\$] | '\\' . | '$' ~[{])*)+ '"' ;
STRING      : '"' (~["\\] | '\\' .)* '"' ;
ID          : [a-zA-Z_][a-zA-Z0-9_]* ;

WS          : [ \t\r\n]+ -> skip ;
DOC_COMMENT : '///' ~[\r\n]* ;
COMMENT     : '//' ~[\r\n]* -> skip ;
// PLUGIN_CONTENT must be last to avoid conflicts - matches everything until newline
PLUGIN_CONTENT : ~[\r\n]+ ;

