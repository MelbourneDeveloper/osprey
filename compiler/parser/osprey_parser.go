// Code generated from osprey.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // osprey
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type ospreyParser struct {
	*antlr.BaseParser
}

var OspreyParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func ospreyParserInit() {
	staticData := &OspreyParserStaticData
	staticData.LiteralNames = []string{
		"", "'match'", "'if'", "'else'", "'select'", "'fn'", "'extern'", "'import'",
		"'type'", "'module'", "'let'", "'mut'", "'effect'", "'perform'", "'handle'",
		"'in'", "'do'", "'spawn'", "'yield'", "'await'", "'fiber'", "'channel'",
		"'send'", "'recv'", "'true'", "'false'", "'where'", "'|>'", "'->'",
		"'=>'", "'_'", "'='", "'=='", "'!='", "'<='", "'>='", "'!'", "'%'",
		"':'", "';'", "','", "'.'", "'|'", "'<'", "'>'", "'('", "')'", "'{'",
		"'}'", "'['", "']'", "'+'", "'-'", "'*'", "'/'",
	}
	staticData.SymbolicNames = []string{
		"", "MATCH", "IF", "ELSE", "SELECT", "FN", "EXTERN", "IMPORT", "TYPE",
		"MODULE", "LET", "MUT", "EFFECT", "PERFORM", "HANDLE", "IN", "DO", "SPAWN",
		"YIELD", "AWAIT", "FIBER", "CHANNEL", "SEND", "RECV", "TRUE", "FALSE",
		"WHERE", "PIPE", "ARROW", "LAMBDA", "UNDERSCORE", "EQ", "EQ_OP", "NE_OP",
		"LE_OP", "GE_OP", "NOT_OP", "MOD_OP", "COLON", "SEMI", "COMMA", "DOT",
		"BAR", "LT", "GT", "LPAREN", "RPAREN", "LBRACE", "RBRACE", "LSQUARE",
		"RSQUARE", "PLUS", "MINUS", "STAR", "SLASH", "INT", "INTERPOLATED_STRING",
		"STRING", "ID", "WS", "DOC_COMMENT", "COMMENT",
	}
	staticData.RuleNames = []string{
		"program", "statement", "importStmt", "letDecl", "assignStmt", "fnDecl",
		"externDecl", "externParamList", "externParam", "paramList", "param",
		"typeDecl", "typeParamList", "unionType", "recordType", "variant", "fieldDeclarations",
		"fieldDeclaration", "typeValidation", "effectDecl", "opDecl", "effectSet",
		"effectList", "handlerExpr", "handlerArm", "handlerParams", "functionCall",
		"booleanExpr", "fieldList", "field", "type", "typeList", "exprStmt",
		"expr", "matchExpr", "selectExpr", "selectArm", "binaryExpr", "comparisonExpr",
		"addExpr", "mulExpr", "unaryExpr", "pipeExpr", "callExpr", "argList",
		"namedArgList", "namedArg", "primary", "typeConstructor", "typeArgs",
		"fieldAssignments", "fieldAssignment", "lambdaExpr", "updateExpr", "blockExpr",
		"literal", "listLiteral", "docComment", "moduleDecl", "moduleBody",
		"moduleStatement", "matchArm", "pattern", "fieldPattern", "blockBody",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 61, 761, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57, 7, 57,
		2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7, 62, 2,
		63, 7, 63, 2, 64, 7, 64, 1, 0, 5, 0, 132, 8, 0, 10, 0, 12, 0, 135, 9, 0,
		1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1,
		148, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 154, 8, 2, 10, 2, 12, 2, 157,
		9, 2, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 163, 8, 3, 1, 3, 1, 3, 1, 3, 1, 4,
		1, 4, 1, 4, 1, 4, 1, 5, 3, 5, 173, 8, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5,
		179, 8, 5, 1, 5, 1, 5, 1, 5, 3, 5, 184, 8, 5, 1, 5, 3, 5, 187, 8, 5, 1,
		5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 195, 8, 5, 1, 6, 3, 6, 198, 8, 6,
		1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 205, 8, 6, 1, 6, 1, 6, 1, 6, 3, 6,
		210, 8, 6, 1, 7, 1, 7, 1, 7, 5, 7, 215, 8, 7, 10, 7, 12, 7, 218, 9, 7,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 5, 9, 227, 8, 9, 10, 9, 12, 9,
		230, 9, 9, 1, 10, 1, 10, 1, 10, 3, 10, 235, 8, 10, 1, 11, 3, 11, 238, 8,
		11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 246, 8, 11, 1, 11,
		1, 11, 1, 11, 3, 11, 251, 8, 11, 1, 11, 3, 11, 254, 8, 11, 1, 12, 1, 12,
		1, 12, 5, 12, 259, 8, 12, 10, 12, 12, 12, 262, 9, 12, 1, 13, 1, 13, 1,
		13, 5, 13, 267, 8, 13, 10, 13, 12, 13, 270, 9, 13, 1, 14, 1, 14, 1, 14,
		1, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 281, 8, 15, 1, 16, 1,
		16, 1, 16, 5, 16, 286, 8, 16, 10, 16, 12, 16, 289, 9, 16, 1, 17, 1, 17,
		1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 1, 19, 3, 19, 299, 8, 19, 1, 19, 1,
		19, 1, 19, 1, 19, 5, 19, 305, 8, 19, 10, 19, 12, 19, 308, 9, 19, 1, 19,
		1, 19, 1, 20, 1, 20, 1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1,
		21, 1, 21, 3, 21, 323, 8, 21, 1, 22, 1, 22, 1, 22, 5, 22, 328, 8, 22, 10,
		22, 12, 22, 331, 9, 22, 1, 23, 1, 23, 1, 23, 4, 23, 336, 8, 23, 11, 23,
		12, 23, 337, 1, 23, 1, 23, 1, 23, 1, 24, 1, 24, 3, 24, 345, 8, 24, 1, 24,
		1, 24, 1, 24, 1, 25, 4, 25, 351, 8, 25, 11, 25, 12, 25, 352, 1, 26, 1,
		26, 1, 26, 3, 26, 358, 8, 26, 1, 26, 1, 26, 1, 27, 1, 27, 1, 28, 1, 28,
		1, 28, 5, 28, 367, 8, 28, 10, 28, 12, 28, 370, 9, 28, 1, 29, 1, 29, 1,
		29, 1, 29, 1, 30, 1, 30, 3, 30, 378, 8, 30, 1, 30, 1, 30, 1, 30, 1, 30,
		1, 30, 1, 30, 3, 30, 386, 8, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1,
		30, 1, 30, 1, 30, 3, 30, 396, 8, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30,
		1, 30, 3, 30, 404, 8, 30, 1, 31, 1, 31, 1, 31, 5, 31, 409, 8, 31, 10, 31,
		12, 31, 412, 9, 31, 1, 32, 1, 32, 1, 33, 1, 33, 1, 34, 1, 34, 1, 34, 1,
		34, 4, 34, 422, 8, 34, 11, 34, 12, 34, 423, 1, 34, 1, 34, 1, 34, 1, 34,
		3, 34, 430, 8, 34, 1, 35, 1, 35, 1, 35, 4, 35, 435, 8, 35, 11, 35, 12,
		35, 436, 1, 35, 1, 35, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36,
		3, 36, 448, 8, 36, 1, 37, 1, 37, 1, 38, 1, 38, 1, 38, 5, 38, 455, 8, 38,
		10, 38, 12, 38, 458, 9, 38, 1, 39, 1, 39, 1, 39, 5, 39, 463, 8, 39, 10,
		39, 12, 39, 466, 9, 39, 1, 40, 1, 40, 1, 40, 5, 40, 471, 8, 40, 10, 40,
		12, 40, 474, 9, 40, 1, 41, 3, 41, 477, 8, 41, 1, 41, 1, 41, 1, 42, 1, 42,
		1, 42, 5, 42, 484, 8, 42, 10, 42, 12, 42, 487, 9, 42, 1, 43, 1, 43, 1,
		43, 4, 43, 492, 8, 43, 11, 43, 12, 43, 493, 1, 43, 1, 43, 3, 43, 498, 8,
		43, 1, 43, 3, 43, 501, 8, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 3, 43,
		508, 8, 43, 1, 43, 4, 43, 511, 8, 43, 11, 43, 12, 43, 512, 1, 43, 1, 43,
		1, 43, 3, 43, 518, 8, 43, 1, 43, 3, 43, 521, 8, 43, 3, 43, 523, 8, 43,
		1, 44, 1, 44, 1, 44, 1, 44, 5, 44, 529, 8, 44, 10, 44, 12, 44, 532, 9,
		44, 3, 44, 534, 8, 44, 1, 45, 1, 45, 1, 45, 4, 45, 539, 8, 45, 11, 45,
		12, 45, 540, 1, 46, 1, 46, 1, 46, 1, 46, 1, 47, 1, 47, 1, 47, 1, 47, 3,
		47, 551, 8, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47,
		1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1,
		47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 3, 47, 578, 8, 47, 1, 47,
		1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1,
		47, 1, 47, 1, 47, 1, 47, 1, 47, 3, 47, 596, 8, 47, 1, 48, 1, 48, 3, 48,
		600, 8, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 49, 1, 49, 1, 49, 1, 49, 1,
		50, 1, 50, 1, 50, 5, 50, 613, 8, 50, 10, 50, 12, 50, 616, 9, 50, 1, 51,
		1, 51, 1, 51, 1, 51, 1, 52, 1, 52, 1, 52, 3, 52, 625, 8, 52, 1, 52, 1,
		52, 1, 52, 3, 52, 630, 8, 52, 1, 52, 1, 52, 1, 52, 1, 52, 3, 52, 636, 8,
		52, 1, 52, 1, 52, 1, 52, 3, 52, 641, 8, 52, 1, 53, 1, 53, 1, 53, 1, 53,
		1, 53, 1, 54, 1, 54, 1, 54, 1, 54, 1, 55, 1, 55, 1, 55, 1, 55, 1, 55, 1,
		55, 3, 55, 658, 8, 55, 1, 56, 1, 56, 1, 56, 1, 56, 5, 56, 664, 8, 56, 10,
		56, 12, 56, 667, 9, 56, 3, 56, 669, 8, 56, 1, 56, 1, 56, 1, 57, 4, 57,
		674, 8, 57, 11, 57, 12, 57, 675, 1, 58, 3, 58, 679, 8, 58, 1, 58, 1, 58,
		1, 58, 1, 58, 1, 58, 1, 58, 1, 59, 5, 59, 688, 8, 59, 10, 59, 12, 59, 691,
		9, 59, 1, 60, 1, 60, 1, 60, 3, 60, 696, 8, 60, 1, 61, 1, 61, 1, 61, 1,
		61, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 3, 62, 708, 8, 62, 1, 62,
		1, 62, 1, 62, 1, 62, 1, 62, 5, 62, 715, 8, 62, 10, 62, 12, 62, 718, 9,
		62, 1, 62, 1, 62, 3, 62, 722, 8, 62, 1, 62, 1, 62, 3, 62, 726, 8, 62, 1,
		62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62, 1, 62,
		1, 62, 1, 62, 1, 62, 3, 62, 742, 8, 62, 1, 63, 1, 63, 1, 63, 5, 63, 747,
		8, 63, 10, 63, 12, 63, 750, 9, 63, 1, 64, 5, 64, 753, 8, 64, 10, 64, 12,
		64, 756, 9, 64, 1, 64, 3, 64, 759, 8, 64, 1, 64, 0, 0, 65, 0, 2, 4, 6,
		8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
		44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78,
		80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112,
		114, 116, 118, 120, 122, 124, 126, 128, 0, 5, 1, 0, 10, 11, 2, 0, 32, 35,
		43, 44, 1, 0, 51, 52, 2, 0, 37, 37, 53, 54, 3, 0, 19, 19, 36, 36, 51, 52,
		813, 0, 133, 1, 0, 0, 0, 2, 147, 1, 0, 0, 0, 4, 149, 1, 0, 0, 0, 6, 158,
		1, 0, 0, 0, 8, 167, 1, 0, 0, 0, 10, 172, 1, 0, 0, 0, 12, 197, 1, 0, 0,
		0, 14, 211, 1, 0, 0, 0, 16, 219, 1, 0, 0, 0, 18, 223, 1, 0, 0, 0, 20, 231,
		1, 0, 0, 0, 22, 237, 1, 0, 0, 0, 24, 255, 1, 0, 0, 0, 26, 263, 1, 0, 0,
		0, 28, 271, 1, 0, 0, 0, 30, 275, 1, 0, 0, 0, 32, 282, 1, 0, 0, 0, 34, 290,
		1, 0, 0, 0, 36, 294, 1, 0, 0, 0, 38, 298, 1, 0, 0, 0, 40, 311, 1, 0, 0,
		0, 42, 322, 1, 0, 0, 0, 44, 324, 1, 0, 0, 0, 46, 332, 1, 0, 0, 0, 48, 342,
		1, 0, 0, 0, 50, 350, 1, 0, 0, 0, 52, 354, 1, 0, 0, 0, 54, 361, 1, 0, 0,
		0, 56, 363, 1, 0, 0, 0, 58, 371, 1, 0, 0, 0, 60, 403, 1, 0, 0, 0, 62, 405,
		1, 0, 0, 0, 64, 413, 1, 0, 0, 0, 66, 415, 1, 0, 0, 0, 68, 429, 1, 0, 0,
		0, 70, 431, 1, 0, 0, 0, 72, 447, 1, 0, 0, 0, 74, 449, 1, 0, 0, 0, 76, 451,
		1, 0, 0, 0, 78, 459, 1, 0, 0, 0, 80, 467, 1, 0, 0, 0, 82, 476, 1, 0, 0,
		0, 84, 480, 1, 0, 0, 0, 86, 522, 1, 0, 0, 0, 88, 533, 1, 0, 0, 0, 90, 535,
		1, 0, 0, 0, 92, 542, 1, 0, 0, 0, 94, 595, 1, 0, 0, 0, 96, 597, 1, 0, 0,
		0, 98, 605, 1, 0, 0, 0, 100, 609, 1, 0, 0, 0, 102, 617, 1, 0, 0, 0, 104,
		640, 1, 0, 0, 0, 106, 642, 1, 0, 0, 0, 108, 647, 1, 0, 0, 0, 110, 657,
		1, 0, 0, 0, 112, 659, 1, 0, 0, 0, 114, 673, 1, 0, 0, 0, 116, 678, 1, 0,
		0, 0, 118, 689, 1, 0, 0, 0, 120, 695, 1, 0, 0, 0, 122, 697, 1, 0, 0, 0,
		124, 741, 1, 0, 0, 0, 126, 743, 1, 0, 0, 0, 128, 754, 1, 0, 0, 0, 130,
		132, 3, 2, 1, 0, 131, 130, 1, 0, 0, 0, 132, 135, 1, 0, 0, 0, 133, 131,
		1, 0, 0, 0, 133, 134, 1, 0, 0, 0, 134, 136, 1, 0, 0, 0, 135, 133, 1, 0,
		0, 0, 136, 137, 5, 0, 0, 1, 137, 1, 1, 0, 0, 0, 138, 148, 3, 4, 2, 0, 139,
		148, 3, 6, 3, 0, 140, 148, 3, 8, 4, 0, 141, 148, 3, 10, 5, 0, 142, 148,
		3, 12, 6, 0, 143, 148, 3, 22, 11, 0, 144, 148, 3, 38, 19, 0, 145, 148,
		3, 116, 58, 0, 146, 148, 3, 64, 32, 0, 147, 138, 1, 0, 0, 0, 147, 139,
		1, 0, 0, 0, 147, 140, 1, 0, 0, 0, 147, 141, 1, 0, 0, 0, 147, 142, 1, 0,
		0, 0, 147, 143, 1, 0, 0, 0, 147, 144, 1, 0, 0, 0, 147, 145, 1, 0, 0, 0,
		147, 146, 1, 0, 0, 0, 148, 3, 1, 0, 0, 0, 149, 150, 5, 7, 0, 0, 150, 155,
		5, 58, 0, 0, 151, 152, 5, 41, 0, 0, 152, 154, 5, 58, 0, 0, 153, 151, 1,
		0, 0, 0, 154, 157, 1, 0, 0, 0, 155, 153, 1, 0, 0, 0, 155, 156, 1, 0, 0,
		0, 156, 5, 1, 0, 0, 0, 157, 155, 1, 0, 0, 0, 158, 159, 7, 0, 0, 0, 159,
		162, 5, 58, 0, 0, 160, 161, 5, 38, 0, 0, 161, 163, 3, 60, 30, 0, 162, 160,
		1, 0, 0, 0, 162, 163, 1, 0, 0, 0, 163, 164, 1, 0, 0, 0, 164, 165, 5, 31,
		0, 0, 165, 166, 3, 66, 33, 0, 166, 7, 1, 0, 0, 0, 167, 168, 5, 58, 0, 0,
		168, 169, 5, 31, 0, 0, 169, 170, 3, 66, 33, 0, 170, 9, 1, 0, 0, 0, 171,
		173, 3, 114, 57, 0, 172, 171, 1, 0, 0, 0, 172, 173, 1, 0, 0, 0, 173, 174,
		1, 0, 0, 0, 174, 175, 5, 5, 0, 0, 175, 176, 5, 58, 0, 0, 176, 178, 5, 45,
		0, 0, 177, 179, 3, 18, 9, 0, 178, 177, 1, 0, 0, 0, 178, 179, 1, 0, 0, 0,
		179, 180, 1, 0, 0, 0, 180, 183, 5, 46, 0, 0, 181, 182, 5, 28, 0, 0, 182,
		184, 3, 60, 30, 0, 183, 181, 1, 0, 0, 0, 183, 184, 1, 0, 0, 0, 184, 186,
		1, 0, 0, 0, 185, 187, 3, 42, 21, 0, 186, 185, 1, 0, 0, 0, 186, 187, 1,
		0, 0, 0, 187, 194, 1, 0, 0, 0, 188, 189, 5, 31, 0, 0, 189, 195, 3, 66,
		33, 0, 190, 191, 5, 47, 0, 0, 191, 192, 3, 128, 64, 0, 192, 193, 5, 48,
		0, 0, 193, 195, 1, 0, 0, 0, 194, 188, 1, 0, 0, 0, 194, 190, 1, 0, 0, 0,
		195, 11, 1, 0, 0, 0, 196, 198, 3, 114, 57, 0, 197, 196, 1, 0, 0, 0, 197,
		198, 1, 0, 0, 0, 198, 199, 1, 0, 0, 0, 199, 200, 5, 6, 0, 0, 200, 201,
		5, 5, 0, 0, 201, 202, 5, 58, 0, 0, 202, 204, 5, 45, 0, 0, 203, 205, 3,
		14, 7, 0, 204, 203, 1, 0, 0, 0, 204, 205, 1, 0, 0, 0, 205, 206, 1, 0, 0,
		0, 206, 209, 5, 46, 0, 0, 207, 208, 5, 28, 0, 0, 208, 210, 3, 60, 30, 0,
		209, 207, 1, 0, 0, 0, 209, 210, 1, 0, 0, 0, 210, 13, 1, 0, 0, 0, 211, 216,
		3, 16, 8, 0, 212, 213, 5, 40, 0, 0, 213, 215, 3, 16, 8, 0, 214, 212, 1,
		0, 0, 0, 215, 218, 1, 0, 0, 0, 216, 214, 1, 0, 0, 0, 216, 217, 1, 0, 0,
		0, 217, 15, 1, 0, 0, 0, 218, 216, 1, 0, 0, 0, 219, 220, 5, 58, 0, 0, 220,
		221, 5, 38, 0, 0, 221, 222, 3, 60, 30, 0, 222, 17, 1, 0, 0, 0, 223, 228,
		3, 20, 10, 0, 224, 225, 5, 40, 0, 0, 225, 227, 3, 20, 10, 0, 226, 224,
		1, 0, 0, 0, 227, 230, 1, 0, 0, 0, 228, 226, 1, 0, 0, 0, 228, 229, 1, 0,
		0, 0, 229, 19, 1, 0, 0, 0, 230, 228, 1, 0, 0, 0, 231, 234, 5, 58, 0, 0,
		232, 233, 5, 38, 0, 0, 233, 235, 3, 60, 30, 0, 234, 232, 1, 0, 0, 0, 234,
		235, 1, 0, 0, 0, 235, 21, 1, 0, 0, 0, 236, 238, 3, 114, 57, 0, 237, 236,
		1, 0, 0, 0, 237, 238, 1, 0, 0, 0, 238, 239, 1, 0, 0, 0, 239, 240, 5, 8,
		0, 0, 240, 245, 5, 58, 0, 0, 241, 242, 5, 43, 0, 0, 242, 243, 3, 24, 12,
		0, 243, 244, 5, 44, 0, 0, 244, 246, 1, 0, 0, 0, 245, 241, 1, 0, 0, 0, 245,
		246, 1, 0, 0, 0, 246, 247, 1, 0, 0, 0, 247, 250, 5, 31, 0, 0, 248, 251,
		3, 26, 13, 0, 249, 251, 3, 28, 14, 0, 250, 248, 1, 0, 0, 0, 250, 249, 1,
		0, 0, 0, 251, 253, 1, 0, 0, 0, 252, 254, 3, 36, 18, 0, 253, 252, 1, 0,
		0, 0, 253, 254, 1, 0, 0, 0, 254, 23, 1, 0, 0, 0, 255, 260, 5, 58, 0, 0,
		256, 257, 5, 40, 0, 0, 257, 259, 5, 58, 0, 0, 258, 256, 1, 0, 0, 0, 259,
		262, 1, 0, 0, 0, 260, 258, 1, 0, 0, 0, 260, 261, 1, 0, 0, 0, 261, 25, 1,
		0, 0, 0, 262, 260, 1, 0, 0, 0, 263, 268, 3, 30, 15, 0, 264, 265, 5, 42,
		0, 0, 265, 267, 3, 30, 15, 0, 266, 264, 1, 0, 0, 0, 267, 270, 1, 0, 0,
		0, 268, 266, 1, 0, 0, 0, 268, 269, 1, 0, 0, 0, 269, 27, 1, 0, 0, 0, 270,
		268, 1, 0, 0, 0, 271, 272, 5, 47, 0, 0, 272, 273, 3, 32, 16, 0, 273, 274,
		5, 48, 0, 0, 274, 29, 1, 0, 0, 0, 275, 280, 5, 58, 0, 0, 276, 277, 5, 47,
		0, 0, 277, 278, 3, 32, 16, 0, 278, 279, 5, 48, 0, 0, 279, 281, 1, 0, 0,
		0, 280, 276, 1, 0, 0, 0, 280, 281, 1, 0, 0, 0, 281, 31, 1, 0, 0, 0, 282,
		287, 3, 34, 17, 0, 283, 284, 5, 40, 0, 0, 284, 286, 3, 34, 17, 0, 285,
		283, 1, 0, 0, 0, 286, 289, 1, 0, 0, 0, 287, 285, 1, 0, 0, 0, 287, 288,
		1, 0, 0, 0, 288, 33, 1, 0, 0, 0, 289, 287, 1, 0, 0, 0, 290, 291, 5, 58,
		0, 0, 291, 292, 5, 38, 0, 0, 292, 293, 3, 60, 30, 0, 293, 35, 1, 0, 0,
		0, 294, 295, 5, 26, 0, 0, 295, 296, 5, 58, 0, 0, 296, 37, 1, 0, 0, 0, 297,
		299, 3, 114, 57, 0, 298, 297, 1, 0, 0, 0, 298, 299, 1, 0, 0, 0, 299, 300,
		1, 0, 0, 0, 300, 301, 5, 12, 0, 0, 301, 302, 5, 58, 0, 0, 302, 306, 5,
		47, 0, 0, 303, 305, 3, 40, 20, 0, 304, 303, 1, 0, 0, 0, 305, 308, 1, 0,
		0, 0, 306, 304, 1, 0, 0, 0, 306, 307, 1, 0, 0, 0, 307, 309, 1, 0, 0, 0,
		308, 306, 1, 0, 0, 0, 309, 310, 5, 48, 0, 0, 310, 39, 1, 0, 0, 0, 311,
		312, 5, 58, 0, 0, 312, 313, 5, 38, 0, 0, 313, 314, 3, 60, 30, 0, 314, 41,
		1, 0, 0, 0, 315, 316, 5, 36, 0, 0, 316, 323, 5, 58, 0, 0, 317, 318, 5,
		36, 0, 0, 318, 319, 5, 49, 0, 0, 319, 320, 3, 44, 22, 0, 320, 321, 5, 50,
		0, 0, 321, 323, 1, 0, 0, 0, 322, 315, 1, 0, 0, 0, 322, 317, 1, 0, 0, 0,
		323, 43, 1, 0, 0, 0, 324, 329, 5, 58, 0, 0, 325, 326, 5, 40, 0, 0, 326,
		328, 5, 58, 0, 0, 327, 325, 1, 0, 0, 0, 328, 331, 1, 0, 0, 0, 329, 327,
		1, 0, 0, 0, 329, 330, 1, 0, 0, 0, 330, 45, 1, 0, 0, 0, 331, 329, 1, 0,
		0, 0, 332, 333, 5, 14, 0, 0, 333, 335, 5, 58, 0, 0, 334, 336, 3, 48, 24,
		0, 335, 334, 1, 0, 0, 0, 336, 337, 1, 0, 0, 0, 337, 335, 1, 0, 0, 0, 337,
		338, 1, 0, 0, 0, 338, 339, 1, 0, 0, 0, 339, 340, 5, 15, 0, 0, 340, 341,
		3, 66, 33, 0, 341, 47, 1, 0, 0, 0, 342, 344, 5, 58, 0, 0, 343, 345, 3,
		50, 25, 0, 344, 343, 1, 0, 0, 0, 344, 345, 1, 0, 0, 0, 345, 346, 1, 0,
		0, 0, 346, 347, 5, 29, 0, 0, 347, 348, 3, 66, 33, 0, 348, 49, 1, 0, 0,
		0, 349, 351, 5, 58, 0, 0, 350, 349, 1, 0, 0, 0, 351, 352, 1, 0, 0, 0, 352,
		350, 1, 0, 0, 0, 352, 353, 1, 0, 0, 0, 353, 51, 1, 0, 0, 0, 354, 355, 5,
		58, 0, 0, 355, 357, 5, 45, 0, 0, 356, 358, 3, 88, 44, 0, 357, 356, 1, 0,
		0, 0, 357, 358, 1, 0, 0, 0, 358, 359, 1, 0, 0, 0, 359, 360, 5, 46, 0, 0,
		360, 53, 1, 0, 0, 0, 361, 362, 3, 76, 38, 0, 362, 55, 1, 0, 0, 0, 363,
		368, 3, 58, 29, 0, 364, 365, 5, 40, 0, 0, 365, 367, 3, 58, 29, 0, 366,
		364, 1, 0, 0, 0, 367, 370, 1, 0, 0, 0, 368, 366, 1, 0, 0, 0, 368, 369,
		1, 0, 0, 0, 369, 57, 1, 0, 0, 0, 370, 368, 1, 0, 0, 0, 371, 372, 5, 58,
		0, 0, 372, 373, 5, 38, 0, 0, 373, 374, 3, 60, 30, 0, 374, 59, 1, 0, 0,
		0, 375, 377, 5, 45, 0, 0, 376, 378, 3, 62, 31, 0, 377, 376, 1, 0, 0, 0,
		377, 378, 1, 0, 0, 0, 378, 379, 1, 0, 0, 0, 379, 380, 5, 46, 0, 0, 380,
		381, 5, 28, 0, 0, 381, 404, 3, 60, 30, 0, 382, 383, 5, 5, 0, 0, 383, 385,
		5, 45, 0, 0, 384, 386, 3, 62, 31, 0, 385, 384, 1, 0, 0, 0, 385, 386, 1,
		0, 0, 0, 386, 387, 1, 0, 0, 0, 387, 388, 5, 46, 0, 0, 388, 389, 5, 28,
		0, 0, 389, 404, 3, 60, 30, 0, 390, 395, 5, 58, 0, 0, 391, 392, 5, 43, 0,
		0, 392, 393, 3, 62, 31, 0, 393, 394, 5, 44, 0, 0, 394, 396, 1, 0, 0, 0,
		395, 391, 1, 0, 0, 0, 395, 396, 1, 0, 0, 0, 396, 404, 1, 0, 0, 0, 397,
		398, 5, 58, 0, 0, 398, 399, 5, 49, 0, 0, 399, 400, 3, 60, 30, 0, 400, 401,
		5, 50, 0, 0, 401, 404, 1, 0, 0, 0, 402, 404, 5, 58, 0, 0, 403, 375, 1,
		0, 0, 0, 403, 382, 1, 0, 0, 0, 403, 390, 1, 0, 0, 0, 403, 397, 1, 0, 0,
		0, 403, 402, 1, 0, 0, 0, 404, 61, 1, 0, 0, 0, 405, 410, 3, 60, 30, 0, 406,
		407, 5, 40, 0, 0, 407, 409, 3, 60, 30, 0, 408, 406, 1, 0, 0, 0, 409, 412,
		1, 0, 0, 0, 410, 408, 1, 0, 0, 0, 410, 411, 1, 0, 0, 0, 411, 63, 1, 0,
		0, 0, 412, 410, 1, 0, 0, 0, 413, 414, 3, 66, 33, 0, 414, 65, 1, 0, 0, 0,
		415, 416, 3, 68, 34, 0, 416, 67, 1, 0, 0, 0, 417, 418, 5, 1, 0, 0, 418,
		419, 3, 66, 33, 0, 419, 421, 5, 47, 0, 0, 420, 422, 3, 122, 61, 0, 421,
		420, 1, 0, 0, 0, 422, 423, 1, 0, 0, 0, 423, 421, 1, 0, 0, 0, 423, 424,
		1, 0, 0, 0, 424, 425, 1, 0, 0, 0, 425, 426, 5, 48, 0, 0, 426, 430, 1, 0,
		0, 0, 427, 430, 3, 70, 35, 0, 428, 430, 3, 74, 37, 0, 429, 417, 1, 0, 0,
		0, 429, 427, 1, 0, 0, 0, 429, 428, 1, 0, 0, 0, 430, 69, 1, 0, 0, 0, 431,
		432, 5, 4, 0, 0, 432, 434, 5, 47, 0, 0, 433, 435, 3, 72, 36, 0, 434, 433,
		1, 0, 0, 0, 435, 436, 1, 0, 0, 0, 436, 434, 1, 0, 0, 0, 436, 437, 1, 0,
		0, 0, 437, 438, 1, 0, 0, 0, 438, 439, 5, 48, 0, 0, 439, 71, 1, 0, 0, 0,
		440, 441, 3, 124, 62, 0, 441, 442, 5, 29, 0, 0, 442, 443, 3, 66, 33, 0,
		443, 448, 1, 0, 0, 0, 444, 445, 5, 30, 0, 0, 445, 446, 5, 29, 0, 0, 446,
		448, 3, 66, 33, 0, 447, 440, 1, 0, 0, 0, 447, 444, 1, 0, 0, 0, 448, 73,
		1, 0, 0, 0, 449, 450, 3, 76, 38, 0, 450, 75, 1, 0, 0, 0, 451, 456, 3, 78,
		39, 0, 452, 453, 7, 1, 0, 0, 453, 455, 3, 78, 39, 0, 454, 452, 1, 0, 0,
		0, 455, 458, 1, 0, 0, 0, 456, 454, 1, 0, 0, 0, 456, 457, 1, 0, 0, 0, 457,
		77, 1, 0, 0, 0, 458, 456, 1, 0, 0, 0, 459, 464, 3, 80, 40, 0, 460, 461,
		7, 2, 0, 0, 461, 463, 3, 80, 40, 0, 462, 460, 1, 0, 0, 0, 463, 466, 1,
		0, 0, 0, 464, 462, 1, 0, 0, 0, 464, 465, 1, 0, 0, 0, 465, 79, 1, 0, 0,
		0, 466, 464, 1, 0, 0, 0, 467, 472, 3, 82, 41, 0, 468, 469, 7, 3, 0, 0,
		469, 471, 3, 82, 41, 0, 470, 468, 1, 0, 0, 0, 471, 474, 1, 0, 0, 0, 472,
		470, 1, 0, 0, 0, 472, 473, 1, 0, 0, 0, 473, 81, 1, 0, 0, 0, 474, 472, 1,
		0, 0, 0, 475, 477, 7, 4, 0, 0, 476, 475, 1, 0, 0, 0, 476, 477, 1, 0, 0,
		0, 477, 478, 1, 0, 0, 0, 478, 479, 3, 84, 42, 0, 479, 83, 1, 0, 0, 0, 480,
		485, 3, 86, 43, 0, 481, 482, 5, 27, 0, 0, 482, 484, 3, 86, 43, 0, 483,
		481, 1, 0, 0, 0, 484, 487, 1, 0, 0, 0, 485, 483, 1, 0, 0, 0, 485, 486,
		1, 0, 0, 0, 486, 85, 1, 0, 0, 0, 487, 485, 1, 0, 0, 0, 488, 491, 3, 94,
		47, 0, 489, 490, 5, 41, 0, 0, 490, 492, 5, 58, 0, 0, 491, 489, 1, 0, 0,
		0, 492, 493, 1, 0, 0, 0, 493, 491, 1, 0, 0, 0, 493, 494, 1, 0, 0, 0, 494,
		500, 1, 0, 0, 0, 495, 497, 5, 45, 0, 0, 496, 498, 3, 88, 44, 0, 497, 496,
		1, 0, 0, 0, 497, 498, 1, 0, 0, 0, 498, 499, 1, 0, 0, 0, 499, 501, 5, 46,
		0, 0, 500, 495, 1, 0, 0, 0, 500, 501, 1, 0, 0, 0, 501, 523, 1, 0, 0, 0,
		502, 510, 3, 94, 47, 0, 503, 504, 5, 41, 0, 0, 504, 505, 5, 58, 0, 0, 505,
		507, 5, 45, 0, 0, 506, 508, 3, 88, 44, 0, 507, 506, 1, 0, 0, 0, 507, 508,
		1, 0, 0, 0, 508, 509, 1, 0, 0, 0, 509, 511, 5, 46, 0, 0, 510, 503, 1, 0,
		0, 0, 511, 512, 1, 0, 0, 0, 512, 510, 1, 0, 0, 0, 512, 513, 1, 0, 0, 0,
		513, 523, 1, 0, 0, 0, 514, 520, 3, 94, 47, 0, 515, 517, 5, 45, 0, 0, 516,
		518, 3, 88, 44, 0, 517, 516, 1, 0, 0, 0, 517, 518, 1, 0, 0, 0, 518, 519,
		1, 0, 0, 0, 519, 521, 5, 46, 0, 0, 520, 515, 1, 0, 0, 0, 520, 521, 1, 0,
		0, 0, 521, 523, 1, 0, 0, 0, 522, 488, 1, 0, 0, 0, 522, 502, 1, 0, 0, 0,
		522, 514, 1, 0, 0, 0, 523, 87, 1, 0, 0, 0, 524, 534, 3, 90, 45, 0, 525,
		530, 3, 66, 33, 0, 526, 527, 5, 40, 0, 0, 527, 529, 3, 66, 33, 0, 528,
		526, 1, 0, 0, 0, 529, 532, 1, 0, 0, 0, 530, 528, 1, 0, 0, 0, 530, 531,
		1, 0, 0, 0, 531, 534, 1, 0, 0, 0, 532, 530, 1, 0, 0, 0, 533, 524, 1, 0,
		0, 0, 533, 525, 1, 0, 0, 0, 534, 89, 1, 0, 0, 0, 535, 538, 3, 92, 46, 0,
		536, 537, 5, 40, 0, 0, 537, 539, 3, 92, 46, 0, 538, 536, 1, 0, 0, 0, 539,
		540, 1, 0, 0, 0, 540, 538, 1, 0, 0, 0, 540, 541, 1, 0, 0, 0, 541, 91, 1,
		0, 0, 0, 542, 543, 5, 58, 0, 0, 543, 544, 5, 38, 0, 0, 544, 545, 3, 66,
		33, 0, 545, 93, 1, 0, 0, 0, 546, 547, 5, 17, 0, 0, 547, 596, 3, 66, 33,
		0, 548, 550, 5, 18, 0, 0, 549, 551, 3, 66, 33, 0, 550, 549, 1, 0, 0, 0,
		550, 551, 1, 0, 0, 0, 551, 596, 1, 0, 0, 0, 552, 553, 5, 19, 0, 0, 553,
		554, 5, 45, 0, 0, 554, 555, 3, 66, 33, 0, 555, 556, 5, 46, 0, 0, 556, 596,
		1, 0, 0, 0, 557, 558, 5, 22, 0, 0, 558, 559, 5, 45, 0, 0, 559, 560, 3,
		66, 33, 0, 560, 561, 5, 40, 0, 0, 561, 562, 3, 66, 33, 0, 562, 563, 5,
		46, 0, 0, 563, 596, 1, 0, 0, 0, 564, 565, 5, 23, 0, 0, 565, 566, 5, 45,
		0, 0, 566, 567, 3, 66, 33, 0, 567, 568, 5, 46, 0, 0, 568, 596, 1, 0, 0,
		0, 569, 570, 5, 4, 0, 0, 570, 596, 3, 70, 35, 0, 571, 572, 5, 13, 0, 0,
		572, 573, 5, 58, 0, 0, 573, 574, 5, 41, 0, 0, 574, 575, 5, 58, 0, 0, 575,
		577, 5, 45, 0, 0, 576, 578, 3, 88, 44, 0, 577, 576, 1, 0, 0, 0, 577, 578,
		1, 0, 0, 0, 578, 579, 1, 0, 0, 0, 579, 596, 5, 46, 0, 0, 580, 596, 3, 46,
		23, 0, 581, 596, 3, 96, 48, 0, 582, 596, 3, 106, 53, 0, 583, 596, 3, 108,
		54, 0, 584, 596, 3, 110, 55, 0, 585, 596, 3, 104, 52, 0, 586, 587, 5, 58,
		0, 0, 587, 588, 5, 49, 0, 0, 588, 589, 5, 55, 0, 0, 589, 596, 5, 50, 0,
		0, 590, 596, 5, 58, 0, 0, 591, 592, 5, 45, 0, 0, 592, 593, 3, 66, 33, 0,
		593, 594, 5, 46, 0, 0, 594, 596, 1, 0, 0, 0, 595, 546, 1, 0, 0, 0, 595,
		548, 1, 0, 0, 0, 595, 552, 1, 0, 0, 0, 595, 557, 1, 0, 0, 0, 595, 564,
		1, 0, 0, 0, 595, 569, 1, 0, 0, 0, 595, 571, 1, 0, 0, 0, 595, 580, 1, 0,
		0, 0, 595, 581, 1, 0, 0, 0, 595, 582, 1, 0, 0, 0, 595, 583, 1, 0, 0, 0,
		595, 584, 1, 0, 0, 0, 595, 585, 1, 0, 0, 0, 595, 586, 1, 0, 0, 0, 595,
		590, 1, 0, 0, 0, 595, 591, 1, 0, 0, 0, 596, 95, 1, 0, 0, 0, 597, 599, 5,
		58, 0, 0, 598, 600, 3, 98, 49, 0, 599, 598, 1, 0, 0, 0, 599, 600, 1, 0,
		0, 0, 600, 601, 1, 0, 0, 0, 601, 602, 5, 47, 0, 0, 602, 603, 3, 100, 50,
		0, 603, 604, 5, 48, 0, 0, 604, 97, 1, 0, 0, 0, 605, 606, 5, 43, 0, 0, 606,
		607, 3, 62, 31, 0, 607, 608, 5, 44, 0, 0, 608, 99, 1, 0, 0, 0, 609, 614,
		3, 102, 51, 0, 610, 611, 5, 40, 0, 0, 611, 613, 3, 102, 51, 0, 612, 610,
		1, 0, 0, 0, 613, 616, 1, 0, 0, 0, 614, 612, 1, 0, 0, 0, 614, 615, 1, 0,
		0, 0, 615, 101, 1, 0, 0, 0, 616, 614, 1, 0, 0, 0, 617, 618, 5, 58, 0, 0,
		618, 619, 5, 38, 0, 0, 619, 620, 3, 66, 33, 0, 620, 103, 1, 0, 0, 0, 621,
		622, 5, 5, 0, 0, 622, 624, 5, 45, 0, 0, 623, 625, 3, 18, 9, 0, 624, 623,
		1, 0, 0, 0, 624, 625, 1, 0, 0, 0, 625, 626, 1, 0, 0, 0, 626, 629, 5, 46,
		0, 0, 627, 628, 5, 28, 0, 0, 628, 630, 3, 60, 30, 0, 629, 627, 1, 0, 0,
		0, 629, 630, 1, 0, 0, 0, 630, 631, 1, 0, 0, 0, 631, 632, 5, 29, 0, 0, 632,
		641, 3, 66, 33, 0, 633, 635, 5, 42, 0, 0, 634, 636, 3, 18, 9, 0, 635, 634,
		1, 0, 0, 0, 635, 636, 1, 0, 0, 0, 636, 637, 1, 0, 0, 0, 637, 638, 5, 42,
		0, 0, 638, 639, 5, 29, 0, 0, 639, 641, 3, 66, 33, 0, 640, 621, 1, 0, 0,
		0, 640, 633, 1, 0, 0, 0, 641, 105, 1, 0, 0, 0, 642, 643, 5, 58, 0, 0, 643,
		644, 5, 47, 0, 0, 644, 645, 3, 100, 50, 0, 645, 646, 5, 48, 0, 0, 646,
		107, 1, 0, 0, 0, 647, 648, 5, 47, 0, 0, 648, 649, 3, 128, 64, 0, 649, 650,
		5, 48, 0, 0, 650, 109, 1, 0, 0, 0, 651, 658, 5, 55, 0, 0, 652, 658, 5,
		57, 0, 0, 653, 658, 5, 56, 0, 0, 654, 658, 5, 24, 0, 0, 655, 658, 5, 25,
		0, 0, 656, 658, 3, 112, 56, 0, 657, 651, 1, 0, 0, 0, 657, 652, 1, 0, 0,
		0, 657, 653, 1, 0, 0, 0, 657, 654, 1, 0, 0, 0, 657, 655, 1, 0, 0, 0, 657,
		656, 1, 0, 0, 0, 658, 111, 1, 0, 0, 0, 659, 668, 5, 49, 0, 0, 660, 665,
		3, 66, 33, 0, 661, 662, 5, 40, 0, 0, 662, 664, 3, 66, 33, 0, 663, 661,
		1, 0, 0, 0, 664, 667, 1, 0, 0, 0, 665, 663, 1, 0, 0, 0, 665, 666, 1, 0,
		0, 0, 666, 669, 1, 0, 0, 0, 667, 665, 1, 0, 0, 0, 668, 660, 1, 0, 0, 0,
		668, 669, 1, 0, 0, 0, 669, 670, 1, 0, 0, 0, 670, 671, 5, 50, 0, 0, 671,
		113, 1, 0, 0, 0, 672, 674, 5, 60, 0, 0, 673, 672, 1, 0, 0, 0, 674, 675,
		1, 0, 0, 0, 675, 673, 1, 0, 0, 0, 675, 676, 1, 0, 0, 0, 676, 115, 1, 0,
		0, 0, 677, 679, 3, 114, 57, 0, 678, 677, 1, 0, 0, 0, 678, 679, 1, 0, 0,
		0, 679, 680, 1, 0, 0, 0, 680, 681, 5, 9, 0, 0, 681, 682, 5, 58, 0, 0, 682,
		683, 5, 47, 0, 0, 683, 684, 3, 118, 59, 0, 684, 685, 5, 48, 0, 0, 685,
		117, 1, 0, 0, 0, 686, 688, 3, 120, 60, 0, 687, 686, 1, 0, 0, 0, 688, 691,
		1, 0, 0, 0, 689, 687, 1, 0, 0, 0, 689, 690, 1, 0, 0, 0, 690, 119, 1, 0,
		0, 0, 691, 689, 1, 0, 0, 0, 692, 696, 3, 6, 3, 0, 693, 696, 3, 10, 5, 0,
		694, 696, 3, 22, 11, 0, 695, 692, 1, 0, 0, 0, 695, 693, 1, 0, 0, 0, 695,
		694, 1, 0, 0, 0, 696, 121, 1, 0, 0, 0, 697, 698, 3, 124, 62, 0, 698, 699,
		5, 29, 0, 0, 699, 700, 3, 66, 33, 0, 700, 123, 1, 0, 0, 0, 701, 742, 3,
		82, 41, 0, 702, 707, 5, 58, 0, 0, 703, 704, 5, 47, 0, 0, 704, 705, 3, 126,
		63, 0, 705, 706, 5, 48, 0, 0, 706, 708, 1, 0, 0, 0, 707, 703, 1, 0, 0,
		0, 707, 708, 1, 0, 0, 0, 708, 742, 1, 0, 0, 0, 709, 721, 5, 58, 0, 0, 710,
		711, 5, 45, 0, 0, 711, 716, 3, 124, 62, 0, 712, 713, 5, 40, 0, 0, 713,
		715, 3, 124, 62, 0, 714, 712, 1, 0, 0, 0, 715, 718, 1, 0, 0, 0, 716, 714,
		1, 0, 0, 0, 716, 717, 1, 0, 0, 0, 717, 719, 1, 0, 0, 0, 718, 716, 1, 0,
		0, 0, 719, 720, 5, 46, 0, 0, 720, 722, 1, 0, 0, 0, 721, 710, 1, 0, 0, 0,
		721, 722, 1, 0, 0, 0, 722, 742, 1, 0, 0, 0, 723, 725, 5, 58, 0, 0, 724,
		726, 5, 58, 0, 0, 725, 724, 1, 0, 0, 0, 725, 726, 1, 0, 0, 0, 726, 742,
		1, 0, 0, 0, 727, 728, 5, 58, 0, 0, 728, 729, 5, 38, 0, 0, 729, 742, 3,
		60, 30, 0, 730, 731, 5, 58, 0, 0, 731, 732, 5, 38, 0, 0, 732, 733, 5, 47,
		0, 0, 733, 734, 3, 126, 63, 0, 734, 735, 5, 48, 0, 0, 735, 742, 1, 0, 0,
		0, 736, 737, 5, 47, 0, 0, 737, 738, 3, 126, 63, 0, 738, 739, 5, 48, 0,
		0, 739, 742, 1, 0, 0, 0, 740, 742, 5, 30, 0, 0, 741, 701, 1, 0, 0, 0, 741,
		702, 1, 0, 0, 0, 741, 709, 1, 0, 0, 0, 741, 723, 1, 0, 0, 0, 741, 727,
		1, 0, 0, 0, 741, 730, 1, 0, 0, 0, 741, 736, 1, 0, 0, 0, 741, 740, 1, 0,
		0, 0, 742, 125, 1, 0, 0, 0, 743, 748, 5, 58, 0, 0, 744, 745, 5, 40, 0,
		0, 745, 747, 5, 58, 0, 0, 746, 744, 1, 0, 0, 0, 747, 750, 1, 0, 0, 0, 748,
		746, 1, 0, 0, 0, 748, 749, 1, 0, 0, 0, 749, 127, 1, 0, 0, 0, 750, 748,
		1, 0, 0, 0, 751, 753, 3, 2, 1, 0, 752, 751, 1, 0, 0, 0, 753, 756, 1, 0,
		0, 0, 754, 752, 1, 0, 0, 0, 754, 755, 1, 0, 0, 0, 755, 758, 1, 0, 0, 0,
		756, 754, 1, 0, 0, 0, 757, 759, 3, 66, 33, 0, 758, 757, 1, 0, 0, 0, 758,
		759, 1, 0, 0, 0, 759, 129, 1, 0, 0, 0, 81, 133, 147, 155, 162, 172, 178,
		183, 186, 194, 197, 204, 209, 216, 228, 234, 237, 245, 250, 253, 260, 268,
		280, 287, 298, 306, 322, 329, 337, 344, 352, 357, 368, 377, 385, 395, 403,
		410, 423, 429, 436, 447, 456, 464, 472, 476, 485, 493, 497, 500, 507, 512,
		517, 520, 522, 530, 533, 540, 550, 577, 595, 599, 614, 624, 629, 635, 640,
		657, 665, 668, 675, 678, 689, 695, 707, 716, 721, 725, 741, 748, 754, 758,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// ospreyParserInit initializes any static state used to implement ospreyParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewospreyParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func OspreyParserInit() {
	staticData := &OspreyParserStaticData
	staticData.once.Do(ospreyParserInit)
}

// NewospreyParser produces a new parser instance for the optional input antlr.TokenStream.
func NewospreyParser(input antlr.TokenStream) *ospreyParser {
	OspreyParserInit()
	this := new(ospreyParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &OspreyParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "osprey.g4"

	return this
}

// ospreyParser tokens.
const (
	ospreyParserEOF                 = antlr.TokenEOF
	ospreyParserMATCH               = 1
	ospreyParserIF                  = 2
	ospreyParserELSE                = 3
	ospreyParserSELECT              = 4
	ospreyParserFN                  = 5
	ospreyParserEXTERN              = 6
	ospreyParserIMPORT              = 7
	ospreyParserTYPE                = 8
	ospreyParserMODULE              = 9
	ospreyParserLET                 = 10
	ospreyParserMUT                 = 11
	ospreyParserEFFECT              = 12
	ospreyParserPERFORM             = 13
	ospreyParserHANDLE              = 14
	ospreyParserIN                  = 15
	ospreyParserDO                  = 16
	ospreyParserSPAWN               = 17
	ospreyParserYIELD               = 18
	ospreyParserAWAIT               = 19
	ospreyParserFIBER               = 20
	ospreyParserCHANNEL             = 21
	ospreyParserSEND                = 22
	ospreyParserRECV                = 23
	ospreyParserTRUE                = 24
	ospreyParserFALSE               = 25
	ospreyParserWHERE               = 26
	ospreyParserPIPE                = 27
	ospreyParserARROW               = 28
	ospreyParserLAMBDA              = 29
	ospreyParserUNDERSCORE          = 30
	ospreyParserEQ                  = 31
	ospreyParserEQ_OP               = 32
	ospreyParserNE_OP               = 33
	ospreyParserLE_OP               = 34
	ospreyParserGE_OP               = 35
	ospreyParserNOT_OP              = 36
	ospreyParserMOD_OP              = 37
	ospreyParserCOLON               = 38
	ospreyParserSEMI                = 39
	ospreyParserCOMMA               = 40
	ospreyParserDOT                 = 41
	ospreyParserBAR                 = 42
	ospreyParserLT                  = 43
	ospreyParserGT                  = 44
	ospreyParserLPAREN              = 45
	ospreyParserRPAREN              = 46
	ospreyParserLBRACE              = 47
	ospreyParserRBRACE              = 48
	ospreyParserLSQUARE             = 49
	ospreyParserRSQUARE             = 50
	ospreyParserPLUS                = 51
	ospreyParserMINUS               = 52
	ospreyParserSTAR                = 53
	ospreyParserSLASH               = 54
	ospreyParserINT                 = 55
	ospreyParserINTERPOLATED_STRING = 56
	ospreyParserSTRING              = 57
	ospreyParserID                  = 58
	ospreyParserWS                  = 59
	ospreyParserDOC_COMMENT         = 60
	ospreyParserCOMMENT             = 61
)

// ospreyParser rules.
const (
	ospreyParserRULE_program           = 0
	ospreyParserRULE_statement         = 1
	ospreyParserRULE_importStmt        = 2
	ospreyParserRULE_letDecl           = 3
	ospreyParserRULE_assignStmt        = 4
	ospreyParserRULE_fnDecl            = 5
	ospreyParserRULE_externDecl        = 6
	ospreyParserRULE_externParamList   = 7
	ospreyParserRULE_externParam       = 8
	ospreyParserRULE_paramList         = 9
	ospreyParserRULE_param             = 10
	ospreyParserRULE_typeDecl          = 11
	ospreyParserRULE_typeParamList     = 12
	ospreyParserRULE_unionType         = 13
	ospreyParserRULE_recordType        = 14
	ospreyParserRULE_variant           = 15
	ospreyParserRULE_fieldDeclarations = 16
	ospreyParserRULE_fieldDeclaration  = 17
	ospreyParserRULE_typeValidation    = 18
	ospreyParserRULE_effectDecl        = 19
	ospreyParserRULE_opDecl            = 20
	ospreyParserRULE_effectSet         = 21
	ospreyParserRULE_effectList        = 22
	ospreyParserRULE_handlerExpr       = 23
	ospreyParserRULE_handlerArm        = 24
	ospreyParserRULE_handlerParams     = 25
	ospreyParserRULE_functionCall      = 26
	ospreyParserRULE_booleanExpr       = 27
	ospreyParserRULE_fieldList         = 28
	ospreyParserRULE_field             = 29
	ospreyParserRULE_type              = 30
	ospreyParserRULE_typeList          = 31
	ospreyParserRULE_exprStmt          = 32
	ospreyParserRULE_expr              = 33
	ospreyParserRULE_matchExpr         = 34
	ospreyParserRULE_selectExpr        = 35
	ospreyParserRULE_selectArm         = 36
	ospreyParserRULE_binaryExpr        = 37
	ospreyParserRULE_comparisonExpr    = 38
	ospreyParserRULE_addExpr           = 39
	ospreyParserRULE_mulExpr           = 40
	ospreyParserRULE_unaryExpr         = 41
	ospreyParserRULE_pipeExpr          = 42
	ospreyParserRULE_callExpr          = 43
	ospreyParserRULE_argList           = 44
	ospreyParserRULE_namedArgList      = 45
	ospreyParserRULE_namedArg          = 46
	ospreyParserRULE_primary           = 47
	ospreyParserRULE_typeConstructor   = 48
	ospreyParserRULE_typeArgs          = 49
	ospreyParserRULE_fieldAssignments  = 50
	ospreyParserRULE_fieldAssignment   = 51
	ospreyParserRULE_lambdaExpr        = 52
	ospreyParserRULE_updateExpr        = 53
	ospreyParserRULE_blockExpr         = 54
	ospreyParserRULE_literal           = 55
	ospreyParserRULE_listLiteral       = 56
	ospreyParserRULE_docComment        = 57
	ospreyParserRULE_moduleDecl        = 58
	ospreyParserRULE_moduleBody        = 59
	ospreyParserRULE_moduleStatement   = 60
	ospreyParserRULE_matchArm          = 61
	ospreyParserRULE_pattern           = 62
	ospreyParserRULE_fieldPattern      = 63
	ospreyParserRULE_blockBody         = 64
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	AllStatement() []IStatementContext
	Statement(i int) IStatementContext

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_program
	return p
}

func InitEmptyProgramContext(p *ProgramContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_program
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) EOF() antlr.TerminalNode {
	return s.GetToken(ospreyParserEOF, 0)
}

func (s *ProgramContext) AllStatement() []IStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatementContext); ok {
			len++
		}
	}

	tst := make([]IStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatementContext); ok {
			tst[i] = t.(IStatementContext)
			i++
		}
	}

	return tst
}

func (s *ProgramContext) Statement(i int) IStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *ProgramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (p *ospreyParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, ospreyParserRULE_program)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(133)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1700852197976080370) != 0 {
		{
			p.SetState(130)
			p.Statement()
		}

		p.SetState(135)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(136)
		p.Match(ospreyParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ImportStmt() IImportStmtContext
	LetDecl() ILetDeclContext
	AssignStmt() IAssignStmtContext
	FnDecl() IFnDeclContext
	ExternDecl() IExternDeclContext
	TypeDecl() ITypeDeclContext
	EffectDecl() IEffectDeclContext
	ModuleDecl() IModuleDeclContext
	ExprStmt() IExprStmtContext

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_statement
	return p
}

func InitEmptyStatementContext(p *StatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_statement
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) ImportStmt() IImportStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImportStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImportStmtContext)
}

func (s *StatementContext) LetDecl() ILetDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILetDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILetDeclContext)
}

func (s *StatementContext) AssignStmt() IAssignStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssignStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssignStmtContext)
}

func (s *StatementContext) FnDecl() IFnDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFnDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFnDeclContext)
}

func (s *StatementContext) ExternDecl() IExternDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExternDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExternDeclContext)
}

func (s *StatementContext) TypeDecl() ITypeDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeDeclContext)
}

func (s *StatementContext) EffectDecl() IEffectDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEffectDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEffectDeclContext)
}

func (s *StatementContext) ModuleDecl() IModuleDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IModuleDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IModuleDeclContext)
}

func (s *StatementContext) ExprStmt() IExprStmtContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprStmtContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprStmtContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (p *ospreyParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, ospreyParserRULE_statement)
	p.SetState(147)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(138)
			p.ImportStmt()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(139)
			p.LetDecl()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(140)
			p.AssignStmt()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(141)
			p.FnDecl()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(142)
			p.ExternDecl()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(143)
			p.TypeDecl()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(144)
			p.EffectDecl()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(145)
			p.ModuleDecl()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(146)
			p.ExprStmt()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IImportStmtContext is an interface to support dynamic dispatch.
type IImportStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IMPORT() antlr.TerminalNode
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsImportStmtContext differentiates from other interfaces.
	IsImportStmtContext()
}

type ImportStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportStmtContext() *ImportStmtContext {
	var p = new(ImportStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_importStmt
	return p
}

func InitEmptyImportStmtContext(p *ImportStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_importStmt
}

func (*ImportStmtContext) IsImportStmtContext() {}

func NewImportStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportStmtContext {
	var p = new(ImportStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_importStmt

	return p
}

func (s *ImportStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportStmtContext) IMPORT() antlr.TerminalNode {
	return s.GetToken(ospreyParserIMPORT, 0)
}

func (s *ImportStmtContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *ImportStmtContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
}

func (s *ImportStmtContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserDOT)
}

func (s *ImportStmtContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserDOT, i)
}

func (s *ImportStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterImportStmt(s)
	}
}

func (s *ImportStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitImportStmt(s)
	}
}

func (p *ospreyParser) ImportStmt() (localctx IImportStmtContext) {
	localctx = NewImportStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, ospreyParserRULE_importStmt)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(149)
		p.Match(ospreyParserIMPORT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(150)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(155)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserDOT {
		{
			p.SetState(151)
			p.Match(ospreyParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(152)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(157)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILetDeclContext is an interface to support dynamic dispatch.
type ILetDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	EQ() antlr.TerminalNode
	Expr() IExprContext
	LET() antlr.TerminalNode
	MUT() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext

	// IsLetDeclContext differentiates from other interfaces.
	IsLetDeclContext()
}

type LetDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLetDeclContext() *LetDeclContext {
	var p = new(LetDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_letDecl
	return p
}

func InitEmptyLetDeclContext(p *LetDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_letDecl
}

func (*LetDeclContext) IsLetDeclContext() {}

func NewLetDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LetDeclContext {
	var p = new(LetDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_letDecl

	return p
}

func (s *LetDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *LetDeclContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *LetDeclContext) EQ() antlr.TerminalNode {
	return s.GetToken(ospreyParserEQ, 0)
}

func (s *LetDeclContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *LetDeclContext) LET() antlr.TerminalNode {
	return s.GetToken(ospreyParserLET, 0)
}

func (s *LetDeclContext) MUT() antlr.TerminalNode {
	return s.GetToken(ospreyParserMUT, 0)
}

func (s *LetDeclContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *LetDeclContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *LetDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LetDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LetDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterLetDecl(s)
	}
}

func (s *LetDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitLetDecl(s)
	}
}

func (p *ospreyParser) LetDecl() (localctx ILetDeclContext) {
	localctx = NewLetDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, ospreyParserRULE_letDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(158)
		_la = p.GetTokenStream().LA(1)

		if !(_la == ospreyParserLET || _la == ospreyParserMUT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(159)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(162)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserCOLON {
		{
			p.SetState(160)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(161)
			p.Type_()
		}

	}
	{
		p.SetState(164)
		p.Match(ospreyParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(165)
		p.Expr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAssignStmtContext is an interface to support dynamic dispatch.
type IAssignStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	EQ() antlr.TerminalNode
	Expr() IExprContext

	// IsAssignStmtContext differentiates from other interfaces.
	IsAssignStmtContext()
}

type AssignStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAssignStmtContext() *AssignStmtContext {
	var p = new(AssignStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_assignStmt
	return p
}

func InitEmptyAssignStmtContext(p *AssignStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_assignStmt
}

func (*AssignStmtContext) IsAssignStmtContext() {}

func NewAssignStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignStmtContext {
	var p = new(AssignStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_assignStmt

	return p
}

func (s *AssignStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignStmtContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *AssignStmtContext) EQ() antlr.TerminalNode {
	return s.GetToken(ospreyParserEQ, 0)
}

func (s *AssignStmtContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *AssignStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterAssignStmt(s)
	}
}

func (s *AssignStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitAssignStmt(s)
	}
}

func (p *ospreyParser) AssignStmt() (localctx IAssignStmtContext) {
	localctx = NewAssignStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, ospreyParserRULE_assignStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(167)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(168)
		p.Match(ospreyParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(169)
		p.Expr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFnDeclContext is an interface to support dynamic dispatch.
type IFnDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FN() antlr.TerminalNode
	ID() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	EQ() antlr.TerminalNode
	Expr() IExprContext
	LBRACE() antlr.TerminalNode
	BlockBody() IBlockBodyContext
	RBRACE() antlr.TerminalNode
	DocComment() IDocCommentContext
	ParamList() IParamListContext
	ARROW() antlr.TerminalNode
	Type_() ITypeContext
	EffectSet() IEffectSetContext

	// IsFnDeclContext differentiates from other interfaces.
	IsFnDeclContext()
}

type FnDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFnDeclContext() *FnDeclContext {
	var p = new(FnDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fnDecl
	return p
}

func InitEmptyFnDeclContext(p *FnDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fnDecl
}

func (*FnDeclContext) IsFnDeclContext() {}

func NewFnDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FnDeclContext {
	var p = new(FnDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_fnDecl

	return p
}

func (s *FnDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *FnDeclContext) FN() antlr.TerminalNode {
	return s.GetToken(ospreyParserFN, 0)
}

func (s *FnDeclContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *FnDeclContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, 0)
}

func (s *FnDeclContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, 0)
}

func (s *FnDeclContext) EQ() antlr.TerminalNode {
	return s.GetToken(ospreyParserEQ, 0)
}

func (s *FnDeclContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *FnDeclContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *FnDeclContext) BlockBody() IBlockBodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockBodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockBodyContext)
}

func (s *FnDeclContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *FnDeclContext) DocComment() IDocCommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDocCommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDocCommentContext)
}

func (s *FnDeclContext) ParamList() IParamListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamListContext)
}

func (s *FnDeclContext) ARROW() antlr.TerminalNode {
	return s.GetToken(ospreyParserARROW, 0)
}

func (s *FnDeclContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *FnDeclContext) EffectSet() IEffectSetContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEffectSetContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEffectSetContext)
}

func (s *FnDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FnDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FnDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterFnDecl(s)
	}
}

func (s *FnDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitFnDecl(s)
	}
}

func (p *ospreyParser) FnDecl() (localctx IFnDeclContext) {
	localctx = NewFnDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, ospreyParserRULE_fnDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(172)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(171)
			p.DocComment()
		}

	}
	{
		p.SetState(174)
		p.Match(ospreyParserFN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(175)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(176)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(178)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserID {
		{
			p.SetState(177)
			p.ParamList()
		}

	}
	{
		p.SetState(180)
		p.Match(ospreyParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(183)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserARROW {
		{
			p.SetState(181)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(182)
			p.Type_()
		}

	}
	p.SetState(186)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserNOT_OP {
		{
			p.SetState(185)
			p.EffectSet()
		}

	}
	p.SetState(194)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserEQ:
		{
			p.SetState(188)
			p.Match(ospreyParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(189)
			p.Expr()
		}

	case ospreyParserLBRACE:
		{
			p.SetState(190)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(191)
			p.BlockBody()
		}
		{
			p.SetState(192)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExternDeclContext is an interface to support dynamic dispatch.
type IExternDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXTERN() antlr.TerminalNode
	FN() antlr.TerminalNode
	ID() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	DocComment() IDocCommentContext
	ExternParamList() IExternParamListContext
	ARROW() antlr.TerminalNode
	Type_() ITypeContext

	// IsExternDeclContext differentiates from other interfaces.
	IsExternDeclContext()
}

type ExternDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExternDeclContext() *ExternDeclContext {
	var p = new(ExternDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_externDecl
	return p
}

func InitEmptyExternDeclContext(p *ExternDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_externDecl
}

func (*ExternDeclContext) IsExternDeclContext() {}

func NewExternDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternDeclContext {
	var p = new(ExternDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_externDecl

	return p
}

func (s *ExternDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *ExternDeclContext) EXTERN() antlr.TerminalNode {
	return s.GetToken(ospreyParserEXTERN, 0)
}

func (s *ExternDeclContext) FN() antlr.TerminalNode {
	return s.GetToken(ospreyParserFN, 0)
}

func (s *ExternDeclContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *ExternDeclContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, 0)
}

func (s *ExternDeclContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, 0)
}

func (s *ExternDeclContext) DocComment() IDocCommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDocCommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDocCommentContext)
}

func (s *ExternDeclContext) ExternParamList() IExternParamListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExternParamListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExternParamListContext)
}

func (s *ExternDeclContext) ARROW() antlr.TerminalNode {
	return s.GetToken(ospreyParserARROW, 0)
}

func (s *ExternDeclContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *ExternDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExternDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExternDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterExternDecl(s)
	}
}

func (s *ExternDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitExternDecl(s)
	}
}

func (p *ospreyParser) ExternDecl() (localctx IExternDeclContext) {
	localctx = NewExternDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, ospreyParserRULE_externDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(197)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(196)
			p.DocComment()
		}

	}
	{
		p.SetState(199)
		p.Match(ospreyParserEXTERN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(200)
		p.Match(ospreyParserFN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(201)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(202)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(204)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserID {
		{
			p.SetState(203)
			p.ExternParamList()
		}

	}
	{
		p.SetState(206)
		p.Match(ospreyParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(209)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserARROW {
		{
			p.SetState(207)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(208)
			p.Type_()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExternParamListContext is an interface to support dynamic dispatch.
type IExternParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExternParam() []IExternParamContext
	ExternParam(i int) IExternParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsExternParamListContext differentiates from other interfaces.
	IsExternParamListContext()
}

type ExternParamListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExternParamListContext() *ExternParamListContext {
	var p = new(ExternParamListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_externParamList
	return p
}

func InitEmptyExternParamListContext(p *ExternParamListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_externParamList
}

func (*ExternParamListContext) IsExternParamListContext() {}

func NewExternParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternParamListContext {
	var p = new(ExternParamListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_externParamList

	return p
}

func (s *ExternParamListContext) GetParser() antlr.Parser { return s.parser }

func (s *ExternParamListContext) AllExternParam() []IExternParamContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExternParamContext); ok {
			len++
		}
	}

	tst := make([]IExternParamContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExternParamContext); ok {
			tst[i] = t.(IExternParamContext)
			i++
		}
	}

	return tst
}

func (s *ExternParamListContext) ExternParam(i int) IExternParamContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExternParamContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExternParamContext)
}

func (s *ExternParamListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *ExternParamListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *ExternParamListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExternParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExternParamListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterExternParamList(s)
	}
}

func (s *ExternParamListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitExternParamList(s)
	}
}

func (p *ospreyParser) ExternParamList() (localctx IExternParamListContext) {
	localctx = NewExternParamListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, ospreyParserRULE_externParamList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(211)
		p.ExternParam()
	}
	p.SetState(216)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(212)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(213)
			p.ExternParam()
		}

		p.SetState(218)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExternParamContext is an interface to support dynamic dispatch.
type IExternParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext

	// IsExternParamContext differentiates from other interfaces.
	IsExternParamContext()
}

type ExternParamContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExternParamContext() *ExternParamContext {
	var p = new(ExternParamContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_externParam
	return p
}

func InitEmptyExternParamContext(p *ExternParamContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_externParam
}

func (*ExternParamContext) IsExternParamContext() {}

func NewExternParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternParamContext {
	var p = new(ExternParamContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_externParam

	return p
}

func (s *ExternParamContext) GetParser() antlr.Parser { return s.parser }

func (s *ExternParamContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *ExternParamContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *ExternParamContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *ExternParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExternParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExternParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterExternParam(s)
	}
}

func (s *ExternParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitExternParam(s)
	}
}

func (p *ospreyParser) ExternParam() (localctx IExternParamContext) {
	localctx = NewExternParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, ospreyParserRULE_externParam)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(219)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(220)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(221)
		p.Type_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IParamListContext is an interface to support dynamic dispatch.
type IParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllParam() []IParamContext
	Param(i int) IParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsParamListContext differentiates from other interfaces.
	IsParamListContext()
}

type ParamListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamListContext() *ParamListContext {
	var p = new(ParamListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_paramList
	return p
}

func InitEmptyParamListContext(p *ParamListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_paramList
}

func (*ParamListContext) IsParamListContext() {}

func NewParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamListContext {
	var p = new(ParamListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_paramList

	return p
}

func (s *ParamListContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamListContext) AllParam() []IParamContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IParamContext); ok {
			len++
		}
	}

	tst := make([]IParamContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IParamContext); ok {
			tst[i] = t.(IParamContext)
			i++
		}
	}

	return tst
}

func (s *ParamListContext) Param(i int) IParamContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ParamListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *ParamListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *ParamListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterParamList(s)
	}
}

func (s *ParamListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitParamList(s)
	}
}

func (p *ospreyParser) ParamList() (localctx IParamListContext) {
	localctx = NewParamListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, ospreyParserRULE_paramList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(223)
		p.Param()
	}
	p.SetState(228)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(224)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(225)
			p.Param()
		}

		p.SetState(230)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IParamContext is an interface to support dynamic dispatch.
type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}

type ParamContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamContext() *ParamContext {
	var p = new(ParamContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_param
	return p
}

func InitEmptyParamContext(p *ParamContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_param
}

func (*ParamContext) IsParamContext() {}

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext {
	var p = new(ParamContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_param

	return p
}

func (s *ParamContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *ParamContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *ParamContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *ParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterParam(s)
	}
}

func (s *ParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitParam(s)
	}
}

func (p *ospreyParser) Param() (localctx IParamContext) {
	localctx = NewParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, ospreyParserRULE_param)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(231)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(234)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserCOLON {
		{
			p.SetState(232)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(233)
			p.Type_()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeDeclContext is an interface to support dynamic dispatch.
type ITypeDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE() antlr.TerminalNode
	ID() antlr.TerminalNode
	EQ() antlr.TerminalNode
	UnionType() IUnionTypeContext
	RecordType() IRecordTypeContext
	DocComment() IDocCommentContext
	LT() antlr.TerminalNode
	TypeParamList() ITypeParamListContext
	GT() antlr.TerminalNode
	TypeValidation() ITypeValidationContext

	// IsTypeDeclContext differentiates from other interfaces.
	IsTypeDeclContext()
}

type TypeDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeDeclContext() *TypeDeclContext {
	var p = new(TypeDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeDecl
	return p
}

func InitEmptyTypeDeclContext(p *TypeDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeDecl
}

func (*TypeDeclContext) IsTypeDeclContext() {}

func NewTypeDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeDeclContext {
	var p = new(TypeDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_typeDecl

	return p
}

func (s *TypeDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeDeclContext) TYPE() antlr.TerminalNode {
	return s.GetToken(ospreyParserTYPE, 0)
}

func (s *TypeDeclContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *TypeDeclContext) EQ() antlr.TerminalNode {
	return s.GetToken(ospreyParserEQ, 0)
}

func (s *TypeDeclContext) UnionType() IUnionTypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionTypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionTypeContext)
}

func (s *TypeDeclContext) RecordType() IRecordTypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRecordTypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRecordTypeContext)
}

func (s *TypeDeclContext) DocComment() IDocCommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDocCommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDocCommentContext)
}

func (s *TypeDeclContext) LT() antlr.TerminalNode {
	return s.GetToken(ospreyParserLT, 0)
}

func (s *TypeDeclContext) TypeParamList() ITypeParamListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeParamListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeParamListContext)
}

func (s *TypeDeclContext) GT() antlr.TerminalNode {
	return s.GetToken(ospreyParserGT, 0)
}

func (s *TypeDeclContext) TypeValidation() ITypeValidationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeValidationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeValidationContext)
}

func (s *TypeDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterTypeDecl(s)
	}
}

func (s *TypeDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitTypeDecl(s)
	}
}

func (p *ospreyParser) TypeDecl() (localctx ITypeDeclContext) {
	localctx = NewTypeDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, ospreyParserRULE_typeDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(237)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(236)
			p.DocComment()
		}

	}
	{
		p.SetState(239)
		p.Match(ospreyParserTYPE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(240)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(245)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserLT {
		{
			p.SetState(241)
			p.Match(ospreyParserLT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(242)
			p.TypeParamList()
		}
		{
			p.SetState(243)
			p.Match(ospreyParserGT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(247)
		p.Match(ospreyParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(250)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserID:
		{
			p.SetState(248)
			p.UnionType()
		}

	case ospreyParserLBRACE:
		{
			p.SetState(249)
			p.RecordType()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserWHERE {
		{
			p.SetState(252)
			p.TypeValidation()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeParamListContext is an interface to support dynamic dispatch.
type ITypeParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTypeParamListContext differentiates from other interfaces.
	IsTypeParamListContext()
}

type TypeParamListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeParamListContext() *TypeParamListContext {
	var p = new(TypeParamListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeParamList
	return p
}

func InitEmptyTypeParamListContext(p *TypeParamListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeParamList
}

func (*TypeParamListContext) IsTypeParamListContext() {}

func NewTypeParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamListContext {
	var p = new(TypeParamListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_typeParamList

	return p
}

func (s *TypeParamListContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeParamListContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *TypeParamListContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
}

func (s *TypeParamListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *TypeParamListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *TypeParamListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeParamListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterTypeParamList(s)
	}
}

func (s *TypeParamListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitTypeParamList(s)
	}
}

func (p *ospreyParser) TypeParamList() (localctx ITypeParamListContext) {
	localctx = NewTypeParamListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, ospreyParserRULE_typeParamList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(255)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(260)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(256)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(257)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(262)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnionTypeContext is an interface to support dynamic dispatch.
type IUnionTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllVariant() []IVariantContext
	Variant(i int) IVariantContext
	AllBAR() []antlr.TerminalNode
	BAR(i int) antlr.TerminalNode

	// IsUnionTypeContext differentiates from other interfaces.
	IsUnionTypeContext()
}

type UnionTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnionTypeContext() *UnionTypeContext {
	var p = new(UnionTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_unionType
	return p
}

func InitEmptyUnionTypeContext(p *UnionTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_unionType
}

func (*UnionTypeContext) IsUnionTypeContext() {}

func NewUnionTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionTypeContext {
	var p = new(UnionTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_unionType

	return p
}

func (s *UnionTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionTypeContext) AllVariant() []IVariantContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IVariantContext); ok {
			len++
		}
	}

	tst := make([]IVariantContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IVariantContext); ok {
			tst[i] = t.(IVariantContext)
			i++
		}
	}

	return tst
}

func (s *UnionTypeContext) Variant(i int) IVariantContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariantContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariantContext)
}

func (s *UnionTypeContext) AllBAR() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserBAR)
}

func (s *UnionTypeContext) BAR(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserBAR, i)
}

func (s *UnionTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterUnionType(s)
	}
}

func (s *UnionTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitUnionType(s)
	}
}

func (p *ospreyParser) UnionType() (localctx IUnionTypeContext) {
	localctx = NewUnionTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, ospreyParserRULE_unionType)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(263)
		p.Variant()
	}
	p.SetState(268)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(264)
				p.Match(ospreyParserBAR)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(265)
				p.Variant()
			}

		}
		p.SetState(270)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRecordTypeContext is an interface to support dynamic dispatch.
type IRecordTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	FieldDeclarations() IFieldDeclarationsContext
	RBRACE() antlr.TerminalNode

	// IsRecordTypeContext differentiates from other interfaces.
	IsRecordTypeContext()
}

type RecordTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRecordTypeContext() *RecordTypeContext {
	var p = new(RecordTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_recordType
	return p
}

func InitEmptyRecordTypeContext(p *RecordTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_recordType
}

func (*RecordTypeContext) IsRecordTypeContext() {}

func NewRecordTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RecordTypeContext {
	var p = new(RecordTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_recordType

	return p
}

func (s *RecordTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *RecordTypeContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *RecordTypeContext) FieldDeclarations() IFieldDeclarationsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldDeclarationsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldDeclarationsContext)
}

func (s *RecordTypeContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *RecordTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecordTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RecordTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterRecordType(s)
	}
}

func (s *RecordTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitRecordType(s)
	}
}

func (p *ospreyParser) RecordType() (localctx IRecordTypeContext) {
	localctx = NewRecordTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, ospreyParserRULE_recordType)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(271)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(272)
		p.FieldDeclarations()
	}
	{
		p.SetState(273)
		p.Match(ospreyParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVariantContext is an interface to support dynamic dispatch.
type IVariantContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	FieldDeclarations() IFieldDeclarationsContext
	RBRACE() antlr.TerminalNode

	// IsVariantContext differentiates from other interfaces.
	IsVariantContext()
}

type VariantContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariantContext() *VariantContext {
	var p = new(VariantContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_variant
	return p
}

func InitEmptyVariantContext(p *VariantContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_variant
}

func (*VariantContext) IsVariantContext() {}

func NewVariantContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariantContext {
	var p = new(VariantContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_variant

	return p
}

func (s *VariantContext) GetParser() antlr.Parser { return s.parser }

func (s *VariantContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *VariantContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *VariantContext) FieldDeclarations() IFieldDeclarationsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldDeclarationsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldDeclarationsContext)
}

func (s *VariantContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *VariantContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariantContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariantContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterVariant(s)
	}
}

func (s *VariantContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitVariant(s)
	}
}

func (p *ospreyParser) Variant() (localctx IVariantContext) {
	localctx = NewVariantContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, ospreyParserRULE_variant)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(275)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(280)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(276)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(277)
			p.FieldDeclarations()
		}
		{
			p.SetState(278)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldDeclarationsContext is an interface to support dynamic dispatch.
type IFieldDeclarationsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllFieldDeclaration() []IFieldDeclarationContext
	FieldDeclaration(i int) IFieldDeclarationContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFieldDeclarationsContext differentiates from other interfaces.
	IsFieldDeclarationsContext()
}

type FieldDeclarationsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldDeclarationsContext() *FieldDeclarationsContext {
	var p = new(FieldDeclarationsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldDeclarations
	return p
}

func InitEmptyFieldDeclarationsContext(p *FieldDeclarationsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldDeclarations
}

func (*FieldDeclarationsContext) IsFieldDeclarationsContext() {}

func NewFieldDeclarationsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldDeclarationsContext {
	var p = new(FieldDeclarationsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_fieldDeclarations

	return p
}

func (s *FieldDeclarationsContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldDeclarationsContext) AllFieldDeclaration() []IFieldDeclarationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldDeclarationContext); ok {
			len++
		}
	}

	tst := make([]IFieldDeclarationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldDeclarationContext); ok {
			tst[i] = t.(IFieldDeclarationContext)
			i++
		}
	}

	return tst
}

func (s *FieldDeclarationsContext) FieldDeclaration(i int) IFieldDeclarationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldDeclarationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldDeclarationContext)
}

func (s *FieldDeclarationsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *FieldDeclarationsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *FieldDeclarationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldDeclarationsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldDeclarationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterFieldDeclarations(s)
	}
}

func (s *FieldDeclarationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitFieldDeclarations(s)
	}
}

func (p *ospreyParser) FieldDeclarations() (localctx IFieldDeclarationsContext) {
	localctx = NewFieldDeclarationsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, ospreyParserRULE_fieldDeclarations)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(282)
		p.FieldDeclaration()
	}
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(283)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(284)
			p.FieldDeclaration()
		}

		p.SetState(289)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldDeclarationContext is an interface to support dynamic dispatch.
type IFieldDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext

	// IsFieldDeclarationContext differentiates from other interfaces.
	IsFieldDeclarationContext()
}

type FieldDeclarationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldDeclarationContext() *FieldDeclarationContext {
	var p = new(FieldDeclarationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldDeclaration
	return p
}

func InitEmptyFieldDeclarationContext(p *FieldDeclarationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldDeclaration
}

func (*FieldDeclarationContext) IsFieldDeclarationContext() {}

func NewFieldDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldDeclarationContext {
	var p = new(FieldDeclarationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_fieldDeclaration

	return p
}

func (s *FieldDeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldDeclarationContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *FieldDeclarationContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *FieldDeclarationContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *FieldDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldDeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterFieldDeclaration(s)
	}
}

func (s *FieldDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitFieldDeclaration(s)
	}
}

func (p *ospreyParser) FieldDeclaration() (localctx IFieldDeclarationContext) {
	localctx = NewFieldDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, ospreyParserRULE_fieldDeclaration)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(290)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(291)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(292)
		p.Type_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeValidationContext is an interface to support dynamic dispatch.
type ITypeValidationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHERE() antlr.TerminalNode
	ID() antlr.TerminalNode

	// IsTypeValidationContext differentiates from other interfaces.
	IsTypeValidationContext()
}

type TypeValidationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeValidationContext() *TypeValidationContext {
	var p = new(TypeValidationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeValidation
	return p
}

func InitEmptyTypeValidationContext(p *TypeValidationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeValidation
}

func (*TypeValidationContext) IsTypeValidationContext() {}

func NewTypeValidationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeValidationContext {
	var p = new(TypeValidationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_typeValidation

	return p
}

func (s *TypeValidationContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeValidationContext) WHERE() antlr.TerminalNode {
	return s.GetToken(ospreyParserWHERE, 0)
}

func (s *TypeValidationContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *TypeValidationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeValidationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeValidationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterTypeValidation(s)
	}
}

func (s *TypeValidationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitTypeValidation(s)
	}
}

func (p *ospreyParser) TypeValidation() (localctx ITypeValidationContext) {
	localctx = NewTypeValidationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, ospreyParserRULE_typeValidation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(294)
		p.Match(ospreyParserWHERE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(295)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEffectDeclContext is an interface to support dynamic dispatch.
type IEffectDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EFFECT() antlr.TerminalNode
	ID() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	DocComment() IDocCommentContext
	AllOpDecl() []IOpDeclContext
	OpDecl(i int) IOpDeclContext

	// IsEffectDeclContext differentiates from other interfaces.
	IsEffectDeclContext()
}

type EffectDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEffectDeclContext() *EffectDeclContext {
	var p = new(EffectDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_effectDecl
	return p
}

func InitEmptyEffectDeclContext(p *EffectDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_effectDecl
}

func (*EffectDeclContext) IsEffectDeclContext() {}

func NewEffectDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EffectDeclContext {
	var p = new(EffectDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_effectDecl

	return p
}

func (s *EffectDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *EffectDeclContext) EFFECT() antlr.TerminalNode {
	return s.GetToken(ospreyParserEFFECT, 0)
}

func (s *EffectDeclContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *EffectDeclContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *EffectDeclContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *EffectDeclContext) DocComment() IDocCommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDocCommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDocCommentContext)
}

func (s *EffectDeclContext) AllOpDecl() []IOpDeclContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOpDeclContext); ok {
			len++
		}
	}

	tst := make([]IOpDeclContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOpDeclContext); ok {
			tst[i] = t.(IOpDeclContext)
			i++
		}
	}

	return tst
}

func (s *EffectDeclContext) OpDecl(i int) IOpDeclContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOpDeclContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOpDeclContext)
}

func (s *EffectDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EffectDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EffectDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterEffectDecl(s)
	}
}

func (s *EffectDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitEffectDecl(s)
	}
}

func (p *ospreyParser) EffectDecl() (localctx IEffectDeclContext) {
	localctx = NewEffectDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, ospreyParserRULE_effectDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(298)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(297)
			p.DocComment()
		}

	}
	{
		p.SetState(300)
		p.Match(ospreyParserEFFECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(301)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(302)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(306)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserID {
		{
			p.SetState(303)
			p.OpDecl()
		}

		p.SetState(308)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(309)
		p.Match(ospreyParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOpDeclContext is an interface to support dynamic dispatch.
type IOpDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext

	// IsOpDeclContext differentiates from other interfaces.
	IsOpDeclContext()
}

type OpDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOpDeclContext() *OpDeclContext {
	var p = new(OpDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_opDecl
	return p
}

func InitEmptyOpDeclContext(p *OpDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_opDecl
}

func (*OpDeclContext) IsOpDeclContext() {}

func NewOpDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OpDeclContext {
	var p = new(OpDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_opDecl

	return p
}

func (s *OpDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *OpDeclContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *OpDeclContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *OpDeclContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *OpDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OpDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OpDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterOpDecl(s)
	}
}

func (s *OpDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitOpDecl(s)
	}
}

func (p *ospreyParser) OpDecl() (localctx IOpDeclContext) {
	localctx = NewOpDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, ospreyParserRULE_opDecl)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(311)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(312)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(313)
		p.Type_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEffectSetContext is an interface to support dynamic dispatch.
type IEffectSetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NOT_OP() antlr.TerminalNode
	ID() antlr.TerminalNode
	LSQUARE() antlr.TerminalNode
	EffectList() IEffectListContext
	RSQUARE() antlr.TerminalNode

	// IsEffectSetContext differentiates from other interfaces.
	IsEffectSetContext()
}

type EffectSetContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEffectSetContext() *EffectSetContext {
	var p = new(EffectSetContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_effectSet
	return p
}

func InitEmptyEffectSetContext(p *EffectSetContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_effectSet
}

func (*EffectSetContext) IsEffectSetContext() {}

func NewEffectSetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EffectSetContext {
	var p = new(EffectSetContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_effectSet

	return p
}

func (s *EffectSetContext) GetParser() antlr.Parser { return s.parser }

func (s *EffectSetContext) NOT_OP() antlr.TerminalNode {
	return s.GetToken(ospreyParserNOT_OP, 0)
}

func (s *EffectSetContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *EffectSetContext) LSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLSQUARE, 0)
}

func (s *EffectSetContext) EffectList() IEffectListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEffectListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEffectListContext)
}

func (s *EffectSetContext) RSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRSQUARE, 0)
}

func (s *EffectSetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EffectSetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EffectSetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterEffectSet(s)
	}
}

func (s *EffectSetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitEffectSet(s)
	}
}

func (p *ospreyParser) EffectSet() (localctx IEffectSetContext) {
	localctx = NewEffectSetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, ospreyParserRULE_effectSet)
	p.SetState(322)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(315)
			p.Match(ospreyParserNOT_OP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(316)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(317)
			p.Match(ospreyParserNOT_OP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(318)
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(319)
			p.EffectList()
		}
		{
			p.SetState(320)
			p.Match(ospreyParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEffectListContext is an interface to support dynamic dispatch.
type IEffectListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsEffectListContext differentiates from other interfaces.
	IsEffectListContext()
}

type EffectListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEffectListContext() *EffectListContext {
	var p = new(EffectListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_effectList
	return p
}

func InitEmptyEffectListContext(p *EffectListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_effectList
}

func (*EffectListContext) IsEffectListContext() {}

func NewEffectListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EffectListContext {
	var p = new(EffectListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_effectList

	return p
}

func (s *EffectListContext) GetParser() antlr.Parser { return s.parser }

func (s *EffectListContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *EffectListContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
}

func (s *EffectListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *EffectListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *EffectListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EffectListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EffectListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterEffectList(s)
	}
}

func (s *EffectListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitEffectList(s)
	}
}

func (p *ospreyParser) EffectList() (localctx IEffectListContext) {
	localctx = NewEffectListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, ospreyParserRULE_effectList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(324)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(329)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(325)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(326)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(331)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHandlerExprContext is an interface to support dynamic dispatch.
type IHandlerExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HANDLE() antlr.TerminalNode
	ID() antlr.TerminalNode
	IN() antlr.TerminalNode
	Expr() IExprContext
	AllHandlerArm() []IHandlerArmContext
	HandlerArm(i int) IHandlerArmContext

	// IsHandlerExprContext differentiates from other interfaces.
	IsHandlerExprContext()
}

type HandlerExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHandlerExprContext() *HandlerExprContext {
	var p = new(HandlerExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_handlerExpr
	return p
}

func InitEmptyHandlerExprContext(p *HandlerExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_handlerExpr
}

func (*HandlerExprContext) IsHandlerExprContext() {}

func NewHandlerExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HandlerExprContext {
	var p = new(HandlerExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_handlerExpr

	return p
}

func (s *HandlerExprContext) GetParser() antlr.Parser { return s.parser }

func (s *HandlerExprContext) HANDLE() antlr.TerminalNode {
	return s.GetToken(ospreyParserHANDLE, 0)
}

func (s *HandlerExprContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *HandlerExprContext) IN() antlr.TerminalNode {
	return s.GetToken(ospreyParserIN, 0)
}

func (s *HandlerExprContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *HandlerExprContext) AllHandlerArm() []IHandlerArmContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IHandlerArmContext); ok {
			len++
		}
	}

	tst := make([]IHandlerArmContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IHandlerArmContext); ok {
			tst[i] = t.(IHandlerArmContext)
			i++
		}
	}

	return tst
}

func (s *HandlerExprContext) HandlerArm(i int) IHandlerArmContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHandlerArmContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHandlerArmContext)
}

func (s *HandlerExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HandlerExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HandlerExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterHandlerExpr(s)
	}
}

func (s *HandlerExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitHandlerExpr(s)
	}
}

func (p *ospreyParser) HandlerExpr() (localctx IHandlerExprContext) {
	localctx = NewHandlerExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, ospreyParserRULE_handlerExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(332)
		p.Match(ospreyParserHANDLE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(333)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(335)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserID {
		{
			p.SetState(334)
			p.HandlerArm()
		}

		p.SetState(337)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(339)
		p.Match(ospreyParserIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(340)
		p.Expr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHandlerArmContext is an interface to support dynamic dispatch.
type IHandlerArmContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	LAMBDA() antlr.TerminalNode
	Expr() IExprContext
	HandlerParams() IHandlerParamsContext

	// IsHandlerArmContext differentiates from other interfaces.
	IsHandlerArmContext()
}

type HandlerArmContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHandlerArmContext() *HandlerArmContext {
	var p = new(HandlerArmContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_handlerArm
	return p
}

func InitEmptyHandlerArmContext(p *HandlerArmContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_handlerArm
}

func (*HandlerArmContext) IsHandlerArmContext() {}

func NewHandlerArmContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HandlerArmContext {
	var p = new(HandlerArmContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_handlerArm

	return p
}

func (s *HandlerArmContext) GetParser() antlr.Parser { return s.parser }

func (s *HandlerArmContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *HandlerArmContext) LAMBDA() antlr.TerminalNode {
	return s.GetToken(ospreyParserLAMBDA, 0)
}

func (s *HandlerArmContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *HandlerArmContext) HandlerParams() IHandlerParamsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHandlerParamsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHandlerParamsContext)
}

func (s *HandlerArmContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HandlerArmContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HandlerArmContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterHandlerArm(s)
	}
}

func (s *HandlerArmContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitHandlerArm(s)
	}
}

func (p *ospreyParser) HandlerArm() (localctx IHandlerArmContext) {
	localctx = NewHandlerArmContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, ospreyParserRULE_handlerArm)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(342)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(344)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserID {
		{
			p.SetState(343)
			p.HandlerParams()
		}

	}
	{
		p.SetState(346)
		p.Match(ospreyParserLAMBDA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(347)
		p.Expr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHandlerParamsContext is an interface to support dynamic dispatch.
type IHandlerParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode

	// IsHandlerParamsContext differentiates from other interfaces.
	IsHandlerParamsContext()
}

type HandlerParamsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHandlerParamsContext() *HandlerParamsContext {
	var p = new(HandlerParamsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_handlerParams
	return p
}

func InitEmptyHandlerParamsContext(p *HandlerParamsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_handlerParams
}

func (*HandlerParamsContext) IsHandlerParamsContext() {}

func NewHandlerParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HandlerParamsContext {
	var p = new(HandlerParamsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_handlerParams

	return p
}

func (s *HandlerParamsContext) GetParser() antlr.Parser { return s.parser }

func (s *HandlerParamsContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *HandlerParamsContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
}

func (s *HandlerParamsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HandlerParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HandlerParamsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterHandlerParams(s)
	}
}

func (s *HandlerParamsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitHandlerParams(s)
	}
}

func (p *ospreyParser) HandlerParams() (localctx IHandlerParamsContext) {
	localctx = NewHandlerParamsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, ospreyParserRULE_handlerParams)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserID {
		{
			p.SetState(349)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(352)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ArgList() IArgListContext

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_functionCall
	return p
}

func InitEmptyFunctionCallContext(p *FunctionCallContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_functionCall
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *FunctionCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, 0)
}

func (s *FunctionCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, 0)
}

func (s *FunctionCallContext) ArgList() IArgListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgListContext)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (p *ospreyParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, ospreyParserRULE_functionCall)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(354)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(355)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(357)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930693369225266) != 0 {
		{
			p.SetState(356)
			p.ArgList()
		}

	}
	{
		p.SetState(359)
		p.Match(ospreyParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBooleanExprContext is an interface to support dynamic dispatch.
type IBooleanExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ComparisonExpr() IComparisonExprContext

	// IsBooleanExprContext differentiates from other interfaces.
	IsBooleanExprContext()
}

type BooleanExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanExprContext() *BooleanExprContext {
	var p = new(BooleanExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_booleanExpr
	return p
}

func InitEmptyBooleanExprContext(p *BooleanExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_booleanExpr
}

func (*BooleanExprContext) IsBooleanExprContext() {}

func NewBooleanExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanExprContext {
	var p = new(BooleanExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_booleanExpr

	return p
}

func (s *BooleanExprContext) GetParser() antlr.Parser { return s.parser }

func (s *BooleanExprContext) ComparisonExpr() IComparisonExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonExprContext)
}

func (s *BooleanExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BooleanExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterBooleanExpr(s)
	}
}

func (s *BooleanExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitBooleanExpr(s)
	}
}

func (p *ospreyParser) BooleanExpr() (localctx IBooleanExprContext) {
	localctx = NewBooleanExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, ospreyParserRULE_booleanExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(361)
		p.ComparisonExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldListContext is an interface to support dynamic dispatch.
type IFieldListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllField() []IFieldContext
	Field(i int) IFieldContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFieldListContext differentiates from other interfaces.
	IsFieldListContext()
}

type FieldListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldListContext() *FieldListContext {
	var p = new(FieldListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldList
	return p
}

func InitEmptyFieldListContext(p *FieldListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldList
}

func (*FieldListContext) IsFieldListContext() {}

func NewFieldListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldListContext {
	var p = new(FieldListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_fieldList

	return p
}

func (s *FieldListContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldListContext) AllField() []IFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldContext); ok {
			len++
		}
	}

	tst := make([]IFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldContext); ok {
			tst[i] = t.(IFieldContext)
			i++
		}
	}

	return tst
}

func (s *FieldListContext) Field(i int) IFieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *FieldListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *FieldListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *FieldListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterFieldList(s)
	}
}

func (s *FieldListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitFieldList(s)
	}
}

func (p *ospreyParser) FieldList() (localctx IFieldListContext) {
	localctx = NewFieldListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, ospreyParserRULE_fieldList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(363)
		p.Field()
	}
	p.SetState(368)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(364)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(365)
			p.Field()
		}

		p.SetState(370)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldContext is an interface to support dynamic dispatch.
type IFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext

	// IsFieldContext differentiates from other interfaces.
	IsFieldContext()
}

type FieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldContext() *FieldContext {
	var p = new(FieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_field
	return p
}

func InitEmptyFieldContext(p *FieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_field
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	var p = new(FieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *FieldContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *FieldContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *FieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterField(s)
	}
}

func (s *FieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitField(s)
	}
}

func (p *ospreyParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, ospreyParserRULE_field)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(371)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(372)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(373)
		p.Type_()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeContext is an interface to support dynamic dispatch.
type ITypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ARROW() antlr.TerminalNode
	Type_() ITypeContext
	TypeList() ITypeListContext
	FN() antlr.TerminalNode
	ID() antlr.TerminalNode
	LT() antlr.TerminalNode
	GT() antlr.TerminalNode
	LSQUARE() antlr.TerminalNode
	RSQUARE() antlr.TerminalNode

	// IsTypeContext differentiates from other interfaces.
	IsTypeContext()
}

type TypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeContext() *TypeContext {
	var p = new(TypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_type
	return p
}

func InitEmptyTypeContext(p *TypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_type
}

func (*TypeContext) IsTypeContext() {}

func NewTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeContext {
	var p = new(TypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_type

	return p
}

func (s *TypeContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, 0)
}

func (s *TypeContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, 0)
}

func (s *TypeContext) ARROW() antlr.TerminalNode {
	return s.GetToken(ospreyParserARROW, 0)
}

func (s *TypeContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *TypeContext) TypeList() ITypeListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeListContext)
}

func (s *TypeContext) FN() antlr.TerminalNode {
	return s.GetToken(ospreyParserFN, 0)
}

func (s *TypeContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *TypeContext) LT() antlr.TerminalNode {
	return s.GetToken(ospreyParserLT, 0)
}

func (s *TypeContext) GT() antlr.TerminalNode {
	return s.GetToken(ospreyParserGT, 0)
}

func (s *TypeContext) LSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLSQUARE, 0)
}

func (s *TypeContext) RSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRSQUARE, 0)
}

func (s *TypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterType(s)
	}
}

func (s *TypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitType(s)
	}
}

func (p *ospreyParser) Type_() (localctx ITypeContext) {
	localctx = NewTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, ospreyParserRULE_type)
	var _la int

	p.SetState(403)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 35, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(375)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(377)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&288265560523800608) != 0 {
			{
				p.SetState(376)
				p.TypeList()
			}

		}
		{
			p.SetState(379)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(380)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(381)
			p.Type_()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(382)
			p.Match(ospreyParserFN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(383)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(385)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&288265560523800608) != 0 {
			{
				p.SetState(384)
				p.TypeList()
			}

		}
		{
			p.SetState(387)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(388)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(389)
			p.Type_()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(390)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(395)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLT {
			{
				p.SetState(391)
				p.Match(ospreyParserLT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(392)
				p.TypeList()
			}
			{
				p.SetState(393)
				p.Match(ospreyParserGT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(397)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(398)
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(399)
			p.Type_()
		}
		{
			p.SetState(400)
			p.Match(ospreyParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(402)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeListContext is an interface to support dynamic dispatch.
type ITypeListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllType_() []ITypeContext
	Type_(i int) ITypeContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTypeListContext differentiates from other interfaces.
	IsTypeListContext()
}

type TypeListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeListContext() *TypeListContext {
	var p = new(TypeListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeList
	return p
}

func InitEmptyTypeListContext(p *TypeListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeList
}

func (*TypeListContext) IsTypeListContext() {}

func NewTypeListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeListContext {
	var p = new(TypeListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_typeList

	return p
}

func (s *TypeListContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeListContext) AllType_() []ITypeContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeContext); ok {
			len++
		}
	}

	tst := make([]ITypeContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeContext); ok {
			tst[i] = t.(ITypeContext)
			i++
		}
	}

	return tst
}

func (s *TypeListContext) Type_(i int) ITypeContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *TypeListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *TypeListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *TypeListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterTypeList(s)
	}
}

func (s *TypeListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitTypeList(s)
	}
}

func (p *ospreyParser) TypeList() (localctx ITypeListContext) {
	localctx = NewTypeListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, ospreyParserRULE_typeList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(405)
		p.Type_()
	}
	p.SetState(410)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(406)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(407)
			p.Type_()
		}

		p.SetState(412)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExprStmtContext is an interface to support dynamic dispatch.
type IExprStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expr() IExprContext

	// IsExprStmtContext differentiates from other interfaces.
	IsExprStmtContext()
}

type ExprStmtContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprStmtContext() *ExprStmtContext {
	var p = new(ExprStmtContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_exprStmt
	return p
}

func InitEmptyExprStmtContext(p *ExprStmtContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_exprStmt
}

func (*ExprStmtContext) IsExprStmtContext() {}

func NewExprStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprStmtContext {
	var p = new(ExprStmtContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_exprStmt

	return p
}

func (s *ExprStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprStmtContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterExprStmt(s)
	}
}

func (s *ExprStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitExprStmt(s)
	}
}

func (p *ospreyParser) ExprStmt() (localctx IExprStmtContext) {
	localctx = NewExprStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, ospreyParserRULE_exprStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(413)
		p.Expr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MatchExpr() IMatchExprContext

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_expr
	return p
}

func InitEmptyExprContext(p *ExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_expr
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) MatchExpr() IMatchExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMatchExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMatchExprContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitExpr(s)
	}
}

func (p *ospreyParser) Expr() (localctx IExprContext) {
	localctx = NewExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, ospreyParserRULE_expr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(415)
		p.MatchExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMatchExprContext is an interface to support dynamic dispatch.
type IMatchExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MATCH() antlr.TerminalNode
	Expr() IExprContext
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllMatchArm() []IMatchArmContext
	MatchArm(i int) IMatchArmContext
	SelectExpr() ISelectExprContext
	BinaryExpr() IBinaryExprContext

	// IsMatchExprContext differentiates from other interfaces.
	IsMatchExprContext()
}

type MatchExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMatchExprContext() *MatchExprContext {
	var p = new(MatchExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_matchExpr
	return p
}

func InitEmptyMatchExprContext(p *MatchExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_matchExpr
}

func (*MatchExprContext) IsMatchExprContext() {}

func NewMatchExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MatchExprContext {
	var p = new(MatchExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_matchExpr

	return p
}

func (s *MatchExprContext) GetParser() antlr.Parser { return s.parser }

func (s *MatchExprContext) MATCH() antlr.TerminalNode {
	return s.GetToken(ospreyParserMATCH, 0)
}

func (s *MatchExprContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MatchExprContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *MatchExprContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *MatchExprContext) AllMatchArm() []IMatchArmContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMatchArmContext); ok {
			len++
		}
	}

	tst := make([]IMatchArmContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMatchArmContext); ok {
			tst[i] = t.(IMatchArmContext)
			i++
		}
	}

	return tst
}

func (s *MatchExprContext) MatchArm(i int) IMatchArmContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMatchArmContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMatchArmContext)
}

func (s *MatchExprContext) SelectExpr() ISelectExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectExprContext)
}

func (s *MatchExprContext) BinaryExpr() IBinaryExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBinaryExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBinaryExprContext)
}

func (s *MatchExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MatchExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MatchExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterMatchExpr(s)
	}
}

func (s *MatchExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitMatchExpr(s)
	}
}

func (p *ospreyParser) MatchExpr() (localctx IMatchExprContext) {
	localctx = NewMatchExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, ospreyParserRULE_matchExpr)
	var _la int

	p.SetState(429)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(417)
			p.Match(ospreyParserMATCH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(418)
			p.Expr()
		}
		{
			p.SetState(419)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(421)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930694442967088) != 0) {
			{
				p.SetState(420)
				p.MatchArm()
			}

			p.SetState(423)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(425)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(427)
			p.SelectExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(428)
			p.BinaryExpr()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectExprContext is an interface to support dynamic dispatch.
type ISelectExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELECT() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllSelectArm() []ISelectArmContext
	SelectArm(i int) ISelectArmContext

	// IsSelectExprContext differentiates from other interfaces.
	IsSelectExprContext()
}

type SelectExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectExprContext() *SelectExprContext {
	var p = new(SelectExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_selectExpr
	return p
}

func InitEmptySelectExprContext(p *SelectExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_selectExpr
}

func (*SelectExprContext) IsSelectExprContext() {}

func NewSelectExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectExprContext {
	var p = new(SelectExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_selectExpr

	return p
}

func (s *SelectExprContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectExprContext) SELECT() antlr.TerminalNode {
	return s.GetToken(ospreyParserSELECT, 0)
}

func (s *SelectExprContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *SelectExprContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *SelectExprContext) AllSelectArm() []ISelectArmContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISelectArmContext); ok {
			len++
		}
	}

	tst := make([]ISelectArmContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISelectArmContext); ok {
			tst[i] = t.(ISelectArmContext)
			i++
		}
	}

	return tst
}

func (s *SelectExprContext) SelectArm(i int) ISelectArmContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectArmContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectArmContext)
}

func (s *SelectExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterSelectExpr(s)
	}
}

func (s *SelectExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitSelectExpr(s)
	}
}

func (p *ospreyParser) SelectExpr() (localctx ISelectExprContext) {
	localctx = NewSelectExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, ospreyParserRULE_selectExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(431)
		p.Match(ospreyParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(432)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(434)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930694442967088) != 0) {
		{
			p.SetState(433)
			p.SelectArm()
		}

		p.SetState(436)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(438)
		p.Match(ospreyParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectArmContext is an interface to support dynamic dispatch.
type ISelectArmContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Pattern() IPatternContext
	LAMBDA() antlr.TerminalNode
	Expr() IExprContext
	UNDERSCORE() antlr.TerminalNode

	// IsSelectArmContext differentiates from other interfaces.
	IsSelectArmContext()
}

type SelectArmContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectArmContext() *SelectArmContext {
	var p = new(SelectArmContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_selectArm
	return p
}

func InitEmptySelectArmContext(p *SelectArmContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_selectArm
}

func (*SelectArmContext) IsSelectArmContext() {}

func NewSelectArmContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectArmContext {
	var p = new(SelectArmContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_selectArm

	return p
}

func (s *SelectArmContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectArmContext) Pattern() IPatternContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPatternContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPatternContext)
}

func (s *SelectArmContext) LAMBDA() antlr.TerminalNode {
	return s.GetToken(ospreyParserLAMBDA, 0)
}

func (s *SelectArmContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SelectArmContext) UNDERSCORE() antlr.TerminalNode {
	return s.GetToken(ospreyParserUNDERSCORE, 0)
}

func (s *SelectArmContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectArmContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectArmContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterSelectArm(s)
	}
}

func (s *SelectArmContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitSelectArm(s)
	}
}

func (p *ospreyParser) SelectArm() (localctx ISelectArmContext) {
	localctx = NewSelectArmContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, ospreyParserRULE_selectArm)
	p.SetState(447)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(440)
			p.Pattern()
		}
		{
			p.SetState(441)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(442)
			p.Expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(444)
			p.Match(ospreyParserUNDERSCORE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(445)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(446)
			p.Expr()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBinaryExprContext is an interface to support dynamic dispatch.
type IBinaryExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ComparisonExpr() IComparisonExprContext

	// IsBinaryExprContext differentiates from other interfaces.
	IsBinaryExprContext()
}

type BinaryExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinaryExprContext() *BinaryExprContext {
	var p = new(BinaryExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_binaryExpr
	return p
}

func InitEmptyBinaryExprContext(p *BinaryExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_binaryExpr
}

func (*BinaryExprContext) IsBinaryExprContext() {}

func NewBinaryExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryExprContext {
	var p = new(BinaryExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_binaryExpr

	return p
}

func (s *BinaryExprContext) GetParser() antlr.Parser { return s.parser }

func (s *BinaryExprContext) ComparisonExpr() IComparisonExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonExprContext)
}

func (s *BinaryExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BinaryExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterBinaryExpr(s)
	}
}

func (s *BinaryExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitBinaryExpr(s)
	}
}

func (p *ospreyParser) BinaryExpr() (localctx IBinaryExprContext) {
	localctx = NewBinaryExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, ospreyParserRULE_binaryExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(449)
		p.ComparisonExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComparisonExprContext is an interface to support dynamic dispatch.
type IComparisonExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAddExpr() []IAddExprContext
	AddExpr(i int) IAddExprContext
	AllEQ_OP() []antlr.TerminalNode
	EQ_OP(i int) antlr.TerminalNode
	AllNE_OP() []antlr.TerminalNode
	NE_OP(i int) antlr.TerminalNode
	AllLT() []antlr.TerminalNode
	LT(i int) antlr.TerminalNode
	AllGT() []antlr.TerminalNode
	GT(i int) antlr.TerminalNode
	AllLE_OP() []antlr.TerminalNode
	LE_OP(i int) antlr.TerminalNode
	AllGE_OP() []antlr.TerminalNode
	GE_OP(i int) antlr.TerminalNode

	// IsComparisonExprContext differentiates from other interfaces.
	IsComparisonExprContext()
}

type ComparisonExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonExprContext() *ComparisonExprContext {
	var p = new(ComparisonExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_comparisonExpr
	return p
}

func InitEmptyComparisonExprContext(p *ComparisonExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_comparisonExpr
}

func (*ComparisonExprContext) IsComparisonExprContext() {}

func NewComparisonExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonExprContext {
	var p = new(ComparisonExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_comparisonExpr

	return p
}

func (s *ComparisonExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonExprContext) AllAddExpr() []IAddExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAddExprContext); ok {
			len++
		}
	}

	tst := make([]IAddExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAddExprContext); ok {
			tst[i] = t.(IAddExprContext)
			i++
		}
	}

	return tst
}

func (s *ComparisonExprContext) AddExpr(i int) IAddExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAddExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAddExprContext)
}

func (s *ComparisonExprContext) AllEQ_OP() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserEQ_OP)
}

func (s *ComparisonExprContext) EQ_OP(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserEQ_OP, i)
}

func (s *ComparisonExprContext) AllNE_OP() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserNE_OP)
}

func (s *ComparisonExprContext) NE_OP(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserNE_OP, i)
}

func (s *ComparisonExprContext) AllLT() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserLT)
}

func (s *ComparisonExprContext) LT(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserLT, i)
}

func (s *ComparisonExprContext) AllGT() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserGT)
}

func (s *ComparisonExprContext) GT(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserGT, i)
}

func (s *ComparisonExprContext) AllLE_OP() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserLE_OP)
}

func (s *ComparisonExprContext) LE_OP(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserLE_OP, i)
}

func (s *ComparisonExprContext) AllGE_OP() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserGE_OP)
}

func (s *ComparisonExprContext) GE_OP(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserGE_OP, i)
}

func (s *ComparisonExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterComparisonExpr(s)
	}
}

func (s *ComparisonExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitComparisonExpr(s)
	}
}

func (p *ospreyParser) ComparisonExpr() (localctx IComparisonExprContext) {
	localctx = NewComparisonExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, ospreyParserRULE_comparisonExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(451)
		p.AddExpr()
	}
	p.SetState(456)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(452)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&26452703576064) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(453)
				p.AddExpr()
			}

		}
		p.SetState(458)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAddExprContext is an interface to support dynamic dispatch.
type IAddExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllMulExpr() []IMulExprContext
	MulExpr(i int) IMulExprContext
	AllPLUS() []antlr.TerminalNode
	PLUS(i int) antlr.TerminalNode
	AllMINUS() []antlr.TerminalNode
	MINUS(i int) antlr.TerminalNode

	// IsAddExprContext differentiates from other interfaces.
	IsAddExprContext()
}

type AddExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAddExprContext() *AddExprContext {
	var p = new(AddExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_addExpr
	return p
}

func InitEmptyAddExprContext(p *AddExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_addExpr
}

func (*AddExprContext) IsAddExprContext() {}

func NewAddExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AddExprContext {
	var p = new(AddExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_addExpr

	return p
}

func (s *AddExprContext) GetParser() antlr.Parser { return s.parser }

func (s *AddExprContext) AllMulExpr() []IMulExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMulExprContext); ok {
			len++
		}
	}

	tst := make([]IMulExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMulExprContext); ok {
			tst[i] = t.(IMulExprContext)
			i++
		}
	}

	return tst
}

func (s *AddExprContext) MulExpr(i int) IMulExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMulExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMulExprContext)
}

func (s *AddExprContext) AllPLUS() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserPLUS)
}

func (s *AddExprContext) PLUS(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserPLUS, i)
}

func (s *AddExprContext) AllMINUS() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserMINUS)
}

func (s *AddExprContext) MINUS(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserMINUS, i)
}

func (s *AddExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AddExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AddExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterAddExpr(s)
	}
}

func (s *AddExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitAddExpr(s)
	}
}

func (p *ospreyParser) AddExpr() (localctx IAddExprContext) {
	localctx = NewAddExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, ospreyParserRULE_addExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(459)
		p.MulExpr()
	}
	p.SetState(464)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(460)
				_la = p.GetTokenStream().LA(1)

				if !(_la == ospreyParserPLUS || _la == ospreyParserMINUS) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(461)
				p.MulExpr()
			}

		}
		p.SetState(466)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMulExprContext is an interface to support dynamic dispatch.
type IMulExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllUnaryExpr() []IUnaryExprContext
	UnaryExpr(i int) IUnaryExprContext
	AllSTAR() []antlr.TerminalNode
	STAR(i int) antlr.TerminalNode
	AllSLASH() []antlr.TerminalNode
	SLASH(i int) antlr.TerminalNode
	AllMOD_OP() []antlr.TerminalNode
	MOD_OP(i int) antlr.TerminalNode

	// IsMulExprContext differentiates from other interfaces.
	IsMulExprContext()
}

type MulExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMulExprContext() *MulExprContext {
	var p = new(MulExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_mulExpr
	return p
}

func InitEmptyMulExprContext(p *MulExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_mulExpr
}

func (*MulExprContext) IsMulExprContext() {}

func NewMulExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MulExprContext {
	var p = new(MulExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_mulExpr

	return p
}

func (s *MulExprContext) GetParser() antlr.Parser { return s.parser }

func (s *MulExprContext) AllUnaryExpr() []IUnaryExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IUnaryExprContext); ok {
			len++
		}
	}

	tst := make([]IUnaryExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IUnaryExprContext); ok {
			tst[i] = t.(IUnaryExprContext)
			i++
		}
	}

	return tst
}

func (s *MulExprContext) UnaryExpr(i int) IUnaryExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnaryExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnaryExprContext)
}

func (s *MulExprContext) AllSTAR() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserSTAR)
}

func (s *MulExprContext) STAR(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserSTAR, i)
}

func (s *MulExprContext) AllSLASH() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserSLASH)
}

func (s *MulExprContext) SLASH(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserSLASH, i)
}

func (s *MulExprContext) AllMOD_OP() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserMOD_OP)
}

func (s *MulExprContext) MOD_OP(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserMOD_OP, i)
}

func (s *MulExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MulExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MulExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterMulExpr(s)
	}
}

func (s *MulExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitMulExpr(s)
	}
}

func (p *ospreyParser) MulExpr() (localctx IMulExprContext) {
	localctx = NewMulExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, ospreyParserRULE_mulExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(467)
		p.UnaryExpr()
	}
	p.SetState(472)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(468)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&27021735203176448) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(469)
				p.UnaryExpr()
			}

		}
		p.SetState(474)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnaryExprContext is an interface to support dynamic dispatch.
type IUnaryExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PipeExpr() IPipeExprContext
	PLUS() antlr.TerminalNode
	MINUS() antlr.TerminalNode
	NOT_OP() antlr.TerminalNode
	AWAIT() antlr.TerminalNode

	// IsUnaryExprContext differentiates from other interfaces.
	IsUnaryExprContext()
}

type UnaryExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnaryExprContext() *UnaryExprContext {
	var p = new(UnaryExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_unaryExpr
	return p
}

func InitEmptyUnaryExprContext(p *UnaryExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_unaryExpr
}

func (*UnaryExprContext) IsUnaryExprContext() {}

func NewUnaryExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryExprContext {
	var p = new(UnaryExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_unaryExpr

	return p
}

func (s *UnaryExprContext) GetParser() antlr.Parser { return s.parser }

func (s *UnaryExprContext) PipeExpr() IPipeExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPipeExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPipeExprContext)
}

func (s *UnaryExprContext) PLUS() antlr.TerminalNode {
	return s.GetToken(ospreyParserPLUS, 0)
}

func (s *UnaryExprContext) MINUS() antlr.TerminalNode {
	return s.GetToken(ospreyParserMINUS, 0)
}

func (s *UnaryExprContext) NOT_OP() antlr.TerminalNode {
	return s.GetToken(ospreyParserNOT_OP, 0)
}

func (s *UnaryExprContext) AWAIT() antlr.TerminalNode {
	return s.GetToken(ospreyParserAWAIT, 0)
}

func (s *UnaryExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnaryExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterUnaryExpr(s)
	}
}

func (s *UnaryExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitUnaryExpr(s)
	}
}

func (p *ospreyParser) UnaryExpr() (localctx IUnaryExprContext) {
	localctx = NewUnaryExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, ospreyParserRULE_unaryExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(476)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(475)
			_la = p.GetTokenStream().LA(1)

			if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&6755468161056768) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(478)
		p.PipeExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPipeExprContext is an interface to support dynamic dispatch.
type IPipeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCallExpr() []ICallExprContext
	CallExpr(i int) ICallExprContext
	AllPIPE() []antlr.TerminalNode
	PIPE(i int) antlr.TerminalNode

	// IsPipeExprContext differentiates from other interfaces.
	IsPipeExprContext()
}

type PipeExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPipeExprContext() *PipeExprContext {
	var p = new(PipeExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pipeExpr
	return p
}

func InitEmptyPipeExprContext(p *PipeExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pipeExpr
}

func (*PipeExprContext) IsPipeExprContext() {}

func NewPipeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PipeExprContext {
	var p = new(PipeExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_pipeExpr

	return p
}

func (s *PipeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *PipeExprContext) AllCallExpr() []ICallExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICallExprContext); ok {
			len++
		}
	}

	tst := make([]ICallExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICallExprContext); ok {
			tst[i] = t.(ICallExprContext)
			i++
		}
	}

	return tst
}

func (s *PipeExprContext) CallExpr(i int) ICallExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICallExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICallExprContext)
}

func (s *PipeExprContext) AllPIPE() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserPIPE)
}

func (s *PipeExprContext) PIPE(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserPIPE, i)
}

func (s *PipeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PipeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PipeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterPipeExpr(s)
	}
}

func (s *PipeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitPipeExpr(s)
	}
}

func (p *ospreyParser) PipeExpr() (localctx IPipeExprContext) {
	localctx = NewPipeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, ospreyParserRULE_pipeExpr)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(480)
		p.CallExpr()
	}
	p.SetState(485)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 45, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(481)
				p.Match(ospreyParserPIPE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(482)
				p.CallExpr()
			}

		}
		p.SetState(487)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 45, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICallExprContext is an interface to support dynamic dispatch.
type ICallExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Primary() IPrimaryContext
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllLPAREN() []antlr.TerminalNode
	LPAREN(i int) antlr.TerminalNode
	AllRPAREN() []antlr.TerminalNode
	RPAREN(i int) antlr.TerminalNode
	AllArgList() []IArgListContext
	ArgList(i int) IArgListContext

	// IsCallExprContext differentiates from other interfaces.
	IsCallExprContext()
}

type CallExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCallExprContext() *CallExprContext {
	var p = new(CallExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_callExpr
	return p
}

func InitEmptyCallExprContext(p *CallExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_callExpr
}

func (*CallExprContext) IsCallExprContext() {}

func NewCallExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CallExprContext {
	var p = new(CallExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_callExpr

	return p
}

func (s *CallExprContext) GetParser() antlr.Parser { return s.parser }

func (s *CallExprContext) Primary() IPrimaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *CallExprContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserDOT)
}

func (s *CallExprContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserDOT, i)
}

func (s *CallExprContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *CallExprContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
}

func (s *CallExprContext) AllLPAREN() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserLPAREN)
}

func (s *CallExprContext) LPAREN(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, i)
}

func (s *CallExprContext) AllRPAREN() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserRPAREN)
}

func (s *CallExprContext) RPAREN(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, i)
}

func (s *CallExprContext) AllArgList() []IArgListContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IArgListContext); ok {
			len++
		}
	}

	tst := make([]IArgListContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IArgListContext); ok {
			tst[i] = t.(IArgListContext)
			i++
		}
	}

	return tst
}

func (s *CallExprContext) ArgList(i int) IArgListContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgListContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgListContext)
}

func (s *CallExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CallExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CallExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterCallExpr(s)
	}
}

func (s *CallExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitCallExpr(s)
	}
}

func (p *ospreyParser) CallExpr() (localctx ICallExprContext) {
	localctx = NewCallExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, ospreyParserRULE_callExpr)
	var _la int

	var _alt int

	p.SetState(522)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 53, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(488)
			p.Primary()
		}
		p.SetState(491)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(489)
					p.Match(ospreyParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(490)
					p.Match(ospreyParserID)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(493)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		p.SetState(500)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(495)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(497)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930693369225266) != 0 {
				{
					p.SetState(496)
					p.ArgList()
				}

			}
			{
				p.SetState(499)
				p.Match(ospreyParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(502)
			p.Primary()
		}
		p.SetState(510)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(503)
					p.Match(ospreyParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(504)
					p.Match(ospreyParserID)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				{
					p.SetState(505)
					p.Match(ospreyParserLPAREN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(507)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930693369225266) != 0 {
					{
						p.SetState(506)
						p.ArgList()
					}

				}
				{
					p.SetState(509)
					p.Match(ospreyParserRPAREN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(512)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(514)
			p.Primary()
		}
		p.SetState(520)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 52, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(515)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(517)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930693369225266) != 0 {
				{
					p.SetState(516)
					p.ArgList()
				}

			}
			{
				p.SetState(519)
				p.Match(ospreyParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArgListContext is an interface to support dynamic dispatch.
type IArgListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NamedArgList() INamedArgListContext
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsArgListContext differentiates from other interfaces.
	IsArgListContext()
}

type ArgListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgListContext() *ArgListContext {
	var p = new(ArgListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_argList
	return p
}

func InitEmptyArgListContext(p *ArgListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_argList
}

func (*ArgListContext) IsArgListContext() {}

func NewArgListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgListContext {
	var p = new(ArgListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_argList

	return p
}

func (s *ArgListContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgListContext) NamedArgList() INamedArgListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamedArgListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamedArgListContext)
}

func (s *ArgListContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ArgListContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ArgListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *ArgListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *ArgListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterArgList(s)
	}
}

func (s *ArgListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitArgList(s)
	}
}

func (p *ospreyParser) ArgList() (localctx IArgListContext) {
	localctx = NewArgListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, ospreyParserRULE_argList)
	var _la int

	p.SetState(533)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 55, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(524)
			p.NamedArgList()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(525)
			p.Expr()
		}
		p.SetState(530)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == ospreyParserCOMMA {
			{
				p.SetState(526)
				p.Match(ospreyParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(527)
				p.Expr()
			}

			p.SetState(532)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INamedArgListContext is an interface to support dynamic dispatch.
type INamedArgListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNamedArg() []INamedArgContext
	NamedArg(i int) INamedArgContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsNamedArgListContext differentiates from other interfaces.
	IsNamedArgListContext()
}

type NamedArgListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamedArgListContext() *NamedArgListContext {
	var p = new(NamedArgListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_namedArgList
	return p
}

func InitEmptyNamedArgListContext(p *NamedArgListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_namedArgList
}

func (*NamedArgListContext) IsNamedArgListContext() {}

func NewNamedArgListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamedArgListContext {
	var p = new(NamedArgListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_namedArgList

	return p
}

func (s *NamedArgListContext) GetParser() antlr.Parser { return s.parser }

func (s *NamedArgListContext) AllNamedArg() []INamedArgContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INamedArgContext); ok {
			len++
		}
	}

	tst := make([]INamedArgContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INamedArgContext); ok {
			tst[i] = t.(INamedArgContext)
			i++
		}
	}

	return tst
}

func (s *NamedArgListContext) NamedArg(i int) INamedArgContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamedArgContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamedArgContext)
}

func (s *NamedArgListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *NamedArgListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *NamedArgListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamedArgListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamedArgListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterNamedArgList(s)
	}
}

func (s *NamedArgListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitNamedArgList(s)
	}
}

func (p *ospreyParser) NamedArgList() (localctx INamedArgListContext) {
	localctx = NewNamedArgListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, ospreyParserRULE_namedArgList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(535)
		p.NamedArg()
	}
	p.SetState(538)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserCOMMA {
		{
			p.SetState(536)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(537)
			p.NamedArg()
		}

		p.SetState(540)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INamedArgContext is an interface to support dynamic dispatch.
type INamedArgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Expr() IExprContext

	// IsNamedArgContext differentiates from other interfaces.
	IsNamedArgContext()
}

type NamedArgContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamedArgContext() *NamedArgContext {
	var p = new(NamedArgContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_namedArg
	return p
}

func InitEmptyNamedArgContext(p *NamedArgContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_namedArg
}

func (*NamedArgContext) IsNamedArgContext() {}

func NewNamedArgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamedArgContext {
	var p = new(NamedArgContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_namedArg

	return p
}

func (s *NamedArgContext) GetParser() antlr.Parser { return s.parser }

func (s *NamedArgContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *NamedArgContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *NamedArgContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *NamedArgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamedArgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamedArgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterNamedArg(s)
	}
}

func (s *NamedArgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitNamedArg(s)
	}
}

func (p *ospreyParser) NamedArg() (localctx INamedArgContext) {
	localctx = NewNamedArgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, ospreyParserRULE_namedArg)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(542)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(543)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(544)
		p.Expr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrimaryContext is an interface to support dynamic dispatch.
type IPrimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SPAWN() antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	YIELD() antlr.TerminalNode
	AWAIT() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	SEND() antlr.TerminalNode
	COMMA() antlr.TerminalNode
	RECV() antlr.TerminalNode
	SELECT() antlr.TerminalNode
	SelectExpr() ISelectExprContext
	PERFORM() antlr.TerminalNode
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	DOT() antlr.TerminalNode
	ArgList() IArgListContext
	HandlerExpr() IHandlerExprContext
	TypeConstructor() ITypeConstructorContext
	UpdateExpr() IUpdateExprContext
	BlockExpr() IBlockExprContext
	Literal() ILiteralContext
	LambdaExpr() ILambdaExprContext
	LSQUARE() antlr.TerminalNode
	INT() antlr.TerminalNode
	RSQUARE() antlr.TerminalNode

	// IsPrimaryContext differentiates from other interfaces.
	IsPrimaryContext()
}

type PrimaryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryContext() *PrimaryContext {
	var p = new(PrimaryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_primary
	return p
}

func InitEmptyPrimaryContext(p *PrimaryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_primary
}

func (*PrimaryContext) IsPrimaryContext() {}

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext {
	var p = new(PrimaryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_primary

	return p
}

func (s *PrimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryContext) SPAWN() antlr.TerminalNode {
	return s.GetToken(ospreyParserSPAWN, 0)
}

func (s *PrimaryContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *PrimaryContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *PrimaryContext) YIELD() antlr.TerminalNode {
	return s.GetToken(ospreyParserYIELD, 0)
}

func (s *PrimaryContext) AWAIT() antlr.TerminalNode {
	return s.GetToken(ospreyParserAWAIT, 0)
}

func (s *PrimaryContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, 0)
}

func (s *PrimaryContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, 0)
}

func (s *PrimaryContext) SEND() antlr.TerminalNode {
	return s.GetToken(ospreyParserSEND, 0)
}

func (s *PrimaryContext) COMMA() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, 0)
}

func (s *PrimaryContext) RECV() antlr.TerminalNode {
	return s.GetToken(ospreyParserRECV, 0)
}

func (s *PrimaryContext) SELECT() antlr.TerminalNode {
	return s.GetToken(ospreyParserSELECT, 0)
}

func (s *PrimaryContext) SelectExpr() ISelectExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectExprContext)
}

func (s *PrimaryContext) PERFORM() antlr.TerminalNode {
	return s.GetToken(ospreyParserPERFORM, 0)
}

func (s *PrimaryContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *PrimaryContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
}

func (s *PrimaryContext) DOT() antlr.TerminalNode {
	return s.GetToken(ospreyParserDOT, 0)
}

func (s *PrimaryContext) ArgList() IArgListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgListContext)
}

func (s *PrimaryContext) HandlerExpr() IHandlerExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHandlerExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHandlerExprContext)
}

func (s *PrimaryContext) TypeConstructor() ITypeConstructorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeConstructorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeConstructorContext)
}

func (s *PrimaryContext) UpdateExpr() IUpdateExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUpdateExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUpdateExprContext)
}

func (s *PrimaryContext) BlockExpr() IBlockExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockExprContext)
}

func (s *PrimaryContext) Literal() ILiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *PrimaryContext) LambdaExpr() ILambdaExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILambdaExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILambdaExprContext)
}

func (s *PrimaryContext) LSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLSQUARE, 0)
}

func (s *PrimaryContext) INT() antlr.TerminalNode {
	return s.GetToken(ospreyParserINT, 0)
}

func (s *PrimaryContext) RSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRSQUARE, 0)
}

func (s *PrimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterPrimary(s)
	}
}

func (s *PrimaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitPrimary(s)
	}
}

func (p *ospreyParser) Primary() (localctx IPrimaryContext) {
	localctx = NewPrimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, ospreyParserRULE_primary)
	var _la int

	p.SetState(595)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 59, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(546)
			p.Match(ospreyParserSPAWN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(547)
			p.Expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(548)
			p.Match(ospreyParserYIELD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(550)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 57, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(549)
				p.Expr()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(552)
			p.Match(ospreyParserAWAIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(553)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(554)
			p.Expr()
		}
		{
			p.SetState(555)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(557)
			p.Match(ospreyParserSEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(558)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(559)
			p.Expr()
		}
		{
			p.SetState(560)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(561)
			p.Expr()
		}
		{
			p.SetState(562)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(564)
			p.Match(ospreyParserRECV)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(565)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(566)
			p.Expr()
		}
		{
			p.SetState(567)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(569)
			p.Match(ospreyParserSELECT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(570)
			p.SelectExpr()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(571)
			p.Match(ospreyParserPERFORM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(572)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(573)
			p.Match(ospreyParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(574)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(575)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(577)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930693369225266) != 0 {
			{
				p.SetState(576)
				p.ArgList()
			}

		}
		{
			p.SetState(579)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(580)
			p.HandlerExpr()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(581)
			p.TypeConstructor()
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(582)
			p.UpdateExpr()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(583)
			p.BlockExpr()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(584)
			p.Literal()
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(585)
			p.LambdaExpr()
		}

	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(586)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(587)
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(588)
			p.Match(ospreyParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(589)
			p.Match(ospreyParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 15:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(590)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 16:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(591)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(592)
			p.Expr()
		}
		{
			p.SetState(593)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeConstructorContext is an interface to support dynamic dispatch.
type ITypeConstructorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	FieldAssignments() IFieldAssignmentsContext
	RBRACE() antlr.TerminalNode
	TypeArgs() ITypeArgsContext

	// IsTypeConstructorContext differentiates from other interfaces.
	IsTypeConstructorContext()
}

type TypeConstructorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeConstructorContext() *TypeConstructorContext {
	var p = new(TypeConstructorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeConstructor
	return p
}

func InitEmptyTypeConstructorContext(p *TypeConstructorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeConstructor
}

func (*TypeConstructorContext) IsTypeConstructorContext() {}

func NewTypeConstructorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeConstructorContext {
	var p = new(TypeConstructorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_typeConstructor

	return p
}

func (s *TypeConstructorContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeConstructorContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *TypeConstructorContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *TypeConstructorContext) FieldAssignments() IFieldAssignmentsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldAssignmentsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldAssignmentsContext)
}

func (s *TypeConstructorContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *TypeConstructorContext) TypeArgs() ITypeArgsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeArgsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeArgsContext)
}

func (s *TypeConstructorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeConstructorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeConstructorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterTypeConstructor(s)
	}
}

func (s *TypeConstructorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitTypeConstructor(s)
	}
}

func (p *ospreyParser) TypeConstructor() (localctx ITypeConstructorContext) {
	localctx = NewTypeConstructorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, ospreyParserRULE_typeConstructor)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(597)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(599)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserLT {
		{
			p.SetState(598)
			p.TypeArgs()
		}

	}
	{
		p.SetState(601)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(602)
		p.FieldAssignments()
	}
	{
		p.SetState(603)
		p.Match(ospreyParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeArgsContext is an interface to support dynamic dispatch.
type ITypeArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LT() antlr.TerminalNode
	TypeList() ITypeListContext
	GT() antlr.TerminalNode

	// IsTypeArgsContext differentiates from other interfaces.
	IsTypeArgsContext()
}

type TypeArgsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeArgsContext() *TypeArgsContext {
	var p = new(TypeArgsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeArgs
	return p
}

func InitEmptyTypeArgsContext(p *TypeArgsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_typeArgs
}

func (*TypeArgsContext) IsTypeArgsContext() {}

func NewTypeArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeArgsContext {
	var p = new(TypeArgsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_typeArgs

	return p
}

func (s *TypeArgsContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeArgsContext) LT() antlr.TerminalNode {
	return s.GetToken(ospreyParserLT, 0)
}

func (s *TypeArgsContext) TypeList() ITypeListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeListContext)
}

func (s *TypeArgsContext) GT() antlr.TerminalNode {
	return s.GetToken(ospreyParserGT, 0)
}

func (s *TypeArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterTypeArgs(s)
	}
}

func (s *TypeArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitTypeArgs(s)
	}
}

func (p *ospreyParser) TypeArgs() (localctx ITypeArgsContext) {
	localctx = NewTypeArgsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, ospreyParserRULE_typeArgs)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(605)
		p.Match(ospreyParserLT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(606)
		p.TypeList()
	}
	{
		p.SetState(607)
		p.Match(ospreyParserGT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldAssignmentsContext is an interface to support dynamic dispatch.
type IFieldAssignmentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllFieldAssignment() []IFieldAssignmentContext
	FieldAssignment(i int) IFieldAssignmentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFieldAssignmentsContext differentiates from other interfaces.
	IsFieldAssignmentsContext()
}

type FieldAssignmentsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldAssignmentsContext() *FieldAssignmentsContext {
	var p = new(FieldAssignmentsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldAssignments
	return p
}

func InitEmptyFieldAssignmentsContext(p *FieldAssignmentsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldAssignments
}

func (*FieldAssignmentsContext) IsFieldAssignmentsContext() {}

func NewFieldAssignmentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldAssignmentsContext {
	var p = new(FieldAssignmentsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_fieldAssignments

	return p
}

func (s *FieldAssignmentsContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldAssignmentsContext) AllFieldAssignment() []IFieldAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldAssignmentContext); ok {
			len++
		}
	}

	tst := make([]IFieldAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldAssignmentContext); ok {
			tst[i] = t.(IFieldAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *FieldAssignmentsContext) FieldAssignment(i int) IFieldAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldAssignmentContext)
}

func (s *FieldAssignmentsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *FieldAssignmentsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *FieldAssignmentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldAssignmentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldAssignmentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterFieldAssignments(s)
	}
}

func (s *FieldAssignmentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitFieldAssignments(s)
	}
}

func (p *ospreyParser) FieldAssignments() (localctx IFieldAssignmentsContext) {
	localctx = NewFieldAssignmentsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, ospreyParserRULE_fieldAssignments)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(609)
		p.FieldAssignment()
	}
	p.SetState(614)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(610)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(611)
			p.FieldAssignment()
		}

		p.SetState(616)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldAssignmentContext is an interface to support dynamic dispatch.
type IFieldAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Expr() IExprContext

	// IsFieldAssignmentContext differentiates from other interfaces.
	IsFieldAssignmentContext()
}

type FieldAssignmentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldAssignmentContext() *FieldAssignmentContext {
	var p = new(FieldAssignmentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldAssignment
	return p
}

func InitEmptyFieldAssignmentContext(p *FieldAssignmentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldAssignment
}

func (*FieldAssignmentContext) IsFieldAssignmentContext() {}

func NewFieldAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldAssignmentContext {
	var p = new(FieldAssignmentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_fieldAssignment

	return p
}

func (s *FieldAssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldAssignmentContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *FieldAssignmentContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *FieldAssignmentContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *FieldAssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldAssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterFieldAssignment(s)
	}
}

func (s *FieldAssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitFieldAssignment(s)
	}
}

func (p *ospreyParser) FieldAssignment() (localctx IFieldAssignmentContext) {
	localctx = NewFieldAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, ospreyParserRULE_fieldAssignment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(617)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(618)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(619)
		p.Expr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILambdaExprContext is an interface to support dynamic dispatch.
type ILambdaExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FN() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	LAMBDA() antlr.TerminalNode
	Expr() IExprContext
	ParamList() IParamListContext
	ARROW() antlr.TerminalNode
	Type_() ITypeContext
	AllBAR() []antlr.TerminalNode
	BAR(i int) antlr.TerminalNode

	// IsLambdaExprContext differentiates from other interfaces.
	IsLambdaExprContext()
}

type LambdaExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLambdaExprContext() *LambdaExprContext {
	var p = new(LambdaExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_lambdaExpr
	return p
}

func InitEmptyLambdaExprContext(p *LambdaExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_lambdaExpr
}

func (*LambdaExprContext) IsLambdaExprContext() {}

func NewLambdaExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LambdaExprContext {
	var p = new(LambdaExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_lambdaExpr

	return p
}

func (s *LambdaExprContext) GetParser() antlr.Parser { return s.parser }

func (s *LambdaExprContext) FN() antlr.TerminalNode {
	return s.GetToken(ospreyParserFN, 0)
}

func (s *LambdaExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, 0)
}

func (s *LambdaExprContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, 0)
}

func (s *LambdaExprContext) LAMBDA() antlr.TerminalNode {
	return s.GetToken(ospreyParserLAMBDA, 0)
}

func (s *LambdaExprContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *LambdaExprContext) ParamList() IParamListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamListContext)
}

func (s *LambdaExprContext) ARROW() antlr.TerminalNode {
	return s.GetToken(ospreyParserARROW, 0)
}

func (s *LambdaExprContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *LambdaExprContext) AllBAR() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserBAR)
}

func (s *LambdaExprContext) BAR(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserBAR, i)
}

func (s *LambdaExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LambdaExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LambdaExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterLambdaExpr(s)
	}
}

func (s *LambdaExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitLambdaExpr(s)
	}
}

func (p *ospreyParser) LambdaExpr() (localctx ILambdaExprContext) {
	localctx = NewLambdaExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, ospreyParserRULE_lambdaExpr)
	var _la int

	p.SetState(640)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserFN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(621)
			p.Match(ospreyParserFN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(622)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(624)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(623)
				p.ParamList()
			}

		}
		{
			p.SetState(626)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(629)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserARROW {
			{
				p.SetState(627)
				p.Match(ospreyParserARROW)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(628)
				p.Type_()
			}

		}
		{
			p.SetState(631)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(632)
			p.Expr()
		}

	case ospreyParserBAR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(633)
			p.Match(ospreyParserBAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(635)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(634)
				p.ParamList()
			}

		}
		{
			p.SetState(637)
			p.Match(ospreyParserBAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(638)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(639)
			p.Expr()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUpdateExprContext is an interface to support dynamic dispatch.
type IUpdateExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	FieldAssignments() IFieldAssignmentsContext
	RBRACE() antlr.TerminalNode

	// IsUpdateExprContext differentiates from other interfaces.
	IsUpdateExprContext()
}

type UpdateExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUpdateExprContext() *UpdateExprContext {
	var p = new(UpdateExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_updateExpr
	return p
}

func InitEmptyUpdateExprContext(p *UpdateExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_updateExpr
}

func (*UpdateExprContext) IsUpdateExprContext() {}

func NewUpdateExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UpdateExprContext {
	var p = new(UpdateExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_updateExpr

	return p
}

func (s *UpdateExprContext) GetParser() antlr.Parser { return s.parser }

func (s *UpdateExprContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *UpdateExprContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *UpdateExprContext) FieldAssignments() IFieldAssignmentsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldAssignmentsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldAssignmentsContext)
}

func (s *UpdateExprContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *UpdateExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UpdateExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UpdateExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterUpdateExpr(s)
	}
}

func (s *UpdateExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitUpdateExpr(s)
	}
}

func (p *ospreyParser) UpdateExpr() (localctx IUpdateExprContext) {
	localctx = NewUpdateExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, ospreyParserRULE_updateExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(642)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(643)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(644)
		p.FieldAssignments()
	}
	{
		p.SetState(645)
		p.Match(ospreyParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBlockExprContext is an interface to support dynamic dispatch.
type IBlockExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	BlockBody() IBlockBodyContext
	RBRACE() antlr.TerminalNode

	// IsBlockExprContext differentiates from other interfaces.
	IsBlockExprContext()
}

type BlockExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBlockExprContext() *BlockExprContext {
	var p = new(BlockExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_blockExpr
	return p
}

func InitEmptyBlockExprContext(p *BlockExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_blockExpr
}

func (*BlockExprContext) IsBlockExprContext() {}

func NewBlockExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlockExprContext {
	var p = new(BlockExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_blockExpr

	return p
}

func (s *BlockExprContext) GetParser() antlr.Parser { return s.parser }

func (s *BlockExprContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *BlockExprContext) BlockBody() IBlockBodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockBodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockBodyContext)
}

func (s *BlockExprContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *BlockExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BlockExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BlockExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterBlockExpr(s)
	}
}

func (s *BlockExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitBlockExpr(s)
	}
}

func (p *ospreyParser) BlockExpr() (localctx IBlockExprContext) {
	localctx = NewBlockExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, ospreyParserRULE_blockExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(647)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(648)
		p.BlockBody()
	}
	{
		p.SetState(649)
		p.Match(ospreyParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT() antlr.TerminalNode
	STRING() antlr.TerminalNode
	INTERPOLATED_STRING() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	ListLiteral() IListLiteralContext

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_literal
	return p
}

func InitEmptyLiteralContext(p *LiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_literal
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) INT() antlr.TerminalNode {
	return s.GetToken(ospreyParserINT, 0)
}

func (s *LiteralContext) STRING() antlr.TerminalNode {
	return s.GetToken(ospreyParserSTRING, 0)
}

func (s *LiteralContext) INTERPOLATED_STRING() antlr.TerminalNode {
	return s.GetToken(ospreyParserINTERPOLATED_STRING, 0)
}

func (s *LiteralContext) TRUE() antlr.TerminalNode {
	return s.GetToken(ospreyParserTRUE, 0)
}

func (s *LiteralContext) FALSE() antlr.TerminalNode {
	return s.GetToken(ospreyParserFALSE, 0)
}

func (s *LiteralContext) ListLiteral() IListLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListLiteralContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterLiteral(s)
	}
}

func (s *LiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitLiteral(s)
	}
}

func (p *ospreyParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, ospreyParserRULE_literal)
	p.SetState(657)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserINT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(651)
			p.Match(ospreyParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserSTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(652)
			p.Match(ospreyParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserINTERPOLATED_STRING:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(653)
			p.Match(ospreyParserINTERPOLATED_STRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserTRUE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(654)
			p.Match(ospreyParserTRUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserFALSE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(655)
			p.Match(ospreyParserFALSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserLSQUARE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(656)
			p.ListLiteral()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IListLiteralContext is an interface to support dynamic dispatch.
type IListLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LSQUARE() antlr.TerminalNode
	RSQUARE() antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsListLiteralContext differentiates from other interfaces.
	IsListLiteralContext()
}

type ListLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyListLiteralContext() *ListLiteralContext {
	var p = new(ListLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_listLiteral
	return p
}

func InitEmptyListLiteralContext(p *ListLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_listLiteral
}

func (*ListLiteralContext) IsListLiteralContext() {}

func NewListLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListLiteralContext {
	var p = new(ListLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_listLiteral

	return p
}

func (s *ListLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *ListLiteralContext) LSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLSQUARE, 0)
}

func (s *ListLiteralContext) RSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRSQUARE, 0)
}

func (s *ListLiteralContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ListLiteralContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ListLiteralContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *ListLiteralContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *ListLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterListLiteral(s)
	}
}

func (s *ListLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitListLiteral(s)
	}
}

func (p *ospreyParser) ListLiteral() (localctx IListLiteralContext) {
	localctx = NewListLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, ospreyParserRULE_listLiteral)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(659)
		p.Match(ospreyParserLSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(668)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930693369225266) != 0 {
		{
			p.SetState(660)
			p.Expr()
		}
		p.SetState(665)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == ospreyParserCOMMA {
			{
				p.SetState(661)
				p.Match(ospreyParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(662)
				p.Expr()
			}

			p.SetState(667)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(670)
		p.Match(ospreyParserRSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDocCommentContext is an interface to support dynamic dispatch.
type IDocCommentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDOC_COMMENT() []antlr.TerminalNode
	DOC_COMMENT(i int) antlr.TerminalNode

	// IsDocCommentContext differentiates from other interfaces.
	IsDocCommentContext()
}

type DocCommentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDocCommentContext() *DocCommentContext {
	var p = new(DocCommentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_docComment
	return p
}

func InitEmptyDocCommentContext(p *DocCommentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_docComment
}

func (*DocCommentContext) IsDocCommentContext() {}

func NewDocCommentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DocCommentContext {
	var p = new(DocCommentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_docComment

	return p
}

func (s *DocCommentContext) GetParser() antlr.Parser { return s.parser }

func (s *DocCommentContext) AllDOC_COMMENT() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserDOC_COMMENT)
}

func (s *DocCommentContext) DOC_COMMENT(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserDOC_COMMENT, i)
}

func (s *DocCommentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DocCommentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DocCommentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterDocComment(s)
	}
}

func (s *DocCommentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitDocComment(s)
	}
}

func (p *ospreyParser) DocComment() (localctx IDocCommentContext) {
	localctx = NewDocCommentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, ospreyParserRULE_docComment)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(673)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(672)
			p.Match(ospreyParserDOC_COMMENT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(675)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IModuleDeclContext is an interface to support dynamic dispatch.
type IModuleDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	MODULE() antlr.TerminalNode
	ID() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	ModuleBody() IModuleBodyContext
	RBRACE() antlr.TerminalNode
	DocComment() IDocCommentContext

	// IsModuleDeclContext differentiates from other interfaces.
	IsModuleDeclContext()
}

type ModuleDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModuleDeclContext() *ModuleDeclContext {
	var p = new(ModuleDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_moduleDecl
	return p
}

func InitEmptyModuleDeclContext(p *ModuleDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_moduleDecl
}

func (*ModuleDeclContext) IsModuleDeclContext() {}

func NewModuleDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModuleDeclContext {
	var p = new(ModuleDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_moduleDecl

	return p
}

func (s *ModuleDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *ModuleDeclContext) MODULE() antlr.TerminalNode {
	return s.GetToken(ospreyParserMODULE, 0)
}

func (s *ModuleDeclContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *ModuleDeclContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *ModuleDeclContext) ModuleBody() IModuleBodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IModuleBodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IModuleBodyContext)
}

func (s *ModuleDeclContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *ModuleDeclContext) DocComment() IDocCommentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDocCommentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDocCommentContext)
}

func (s *ModuleDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModuleDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModuleDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterModuleDecl(s)
	}
}

func (s *ModuleDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitModuleDecl(s)
	}
}

func (p *ospreyParser) ModuleDecl() (localctx IModuleDeclContext) {
	localctx = NewModuleDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, ospreyParserRULE_moduleDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(678)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(677)
			p.DocComment()
		}

	}
	{
		p.SetState(680)
		p.Match(ospreyParserMODULE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(681)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(682)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(683)
		p.ModuleBody()
	}
	{
		p.SetState(684)
		p.Match(ospreyParserRBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IModuleBodyContext is an interface to support dynamic dispatch.
type IModuleBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllModuleStatement() []IModuleStatementContext
	ModuleStatement(i int) IModuleStatementContext

	// IsModuleBodyContext differentiates from other interfaces.
	IsModuleBodyContext()
}

type ModuleBodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModuleBodyContext() *ModuleBodyContext {
	var p = new(ModuleBodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_moduleBody
	return p
}

func InitEmptyModuleBodyContext(p *ModuleBodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_moduleBody
}

func (*ModuleBodyContext) IsModuleBodyContext() {}

func NewModuleBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModuleBodyContext {
	var p = new(ModuleBodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_moduleBody

	return p
}

func (s *ModuleBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *ModuleBodyContext) AllModuleStatement() []IModuleStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IModuleStatementContext); ok {
			len++
		}
	}

	tst := make([]IModuleStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IModuleStatementContext); ok {
			tst[i] = t.(IModuleStatementContext)
			i++
		}
	}

	return tst
}

func (s *ModuleBodyContext) ModuleStatement(i int) IModuleStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IModuleStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IModuleStatementContext)
}

func (s *ModuleBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModuleBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModuleBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterModuleBody(s)
	}
}

func (s *ModuleBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitModuleBody(s)
	}
}

func (p *ospreyParser) ModuleBody() (localctx IModuleBodyContext) {
	localctx = NewModuleBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, ospreyParserRULE_moduleBody)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(689)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1152921504606850336) != 0 {
		{
			p.SetState(686)
			p.ModuleStatement()
		}

		p.SetState(691)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IModuleStatementContext is an interface to support dynamic dispatch.
type IModuleStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LetDecl() ILetDeclContext
	FnDecl() IFnDeclContext
	TypeDecl() ITypeDeclContext

	// IsModuleStatementContext differentiates from other interfaces.
	IsModuleStatementContext()
}

type ModuleStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyModuleStatementContext() *ModuleStatementContext {
	var p = new(ModuleStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_moduleStatement
	return p
}

func InitEmptyModuleStatementContext(p *ModuleStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_moduleStatement
}

func (*ModuleStatementContext) IsModuleStatementContext() {}

func NewModuleStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ModuleStatementContext {
	var p = new(ModuleStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_moduleStatement

	return p
}

func (s *ModuleStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ModuleStatementContext) LetDecl() ILetDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILetDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILetDeclContext)
}

func (s *ModuleStatementContext) FnDecl() IFnDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFnDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFnDeclContext)
}

func (s *ModuleStatementContext) TypeDecl() ITypeDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeDeclContext)
}

func (s *ModuleStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ModuleStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ModuleStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterModuleStatement(s)
	}
}

func (s *ModuleStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitModuleStatement(s)
	}
}

func (p *ospreyParser) ModuleStatement() (localctx IModuleStatementContext) {
	localctx = NewModuleStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, ospreyParserRULE_moduleStatement)
	p.SetState(695)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 72, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(692)
			p.LetDecl()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(693)
			p.FnDecl()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(694)
			p.TypeDecl()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMatchArmContext is an interface to support dynamic dispatch.
type IMatchArmContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Pattern() IPatternContext
	LAMBDA() antlr.TerminalNode
	Expr() IExprContext

	// IsMatchArmContext differentiates from other interfaces.
	IsMatchArmContext()
}

type MatchArmContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMatchArmContext() *MatchArmContext {
	var p = new(MatchArmContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_matchArm
	return p
}

func InitEmptyMatchArmContext(p *MatchArmContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_matchArm
}

func (*MatchArmContext) IsMatchArmContext() {}

func NewMatchArmContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MatchArmContext {
	var p = new(MatchArmContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_matchArm

	return p
}

func (s *MatchArmContext) GetParser() antlr.Parser { return s.parser }

func (s *MatchArmContext) Pattern() IPatternContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPatternContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPatternContext)
}

func (s *MatchArmContext) LAMBDA() antlr.TerminalNode {
	return s.GetToken(ospreyParserLAMBDA, 0)
}

func (s *MatchArmContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MatchArmContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MatchArmContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MatchArmContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterMatchArm(s)
	}
}

func (s *MatchArmContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitMatchArm(s)
	}
}

func (p *ospreyParser) MatchArm() (localctx IMatchArmContext) {
	localctx = NewMatchArmContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, ospreyParserRULE_matchArm)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(697)
		p.Pattern()
	}
	{
		p.SetState(698)
		p.Match(ospreyParserLAMBDA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(699)
		p.Expr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPatternContext is an interface to support dynamic dispatch.
type IPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UnaryExpr() IUnaryExprContext
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	FieldPattern() IFieldPatternContext
	RBRACE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllPattern() []IPatternContext
	Pattern(i int) IPatternContext
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	COLON() antlr.TerminalNode
	Type_() ITypeContext
	UNDERSCORE() antlr.TerminalNode

	// IsPatternContext differentiates from other interfaces.
	IsPatternContext()
}

type PatternContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPatternContext() *PatternContext {
	var p = new(PatternContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pattern
	return p
}

func InitEmptyPatternContext(p *PatternContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pattern
}

func (*PatternContext) IsPatternContext() {}

func NewPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PatternContext {
	var p = new(PatternContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_pattern

	return p
}

func (s *PatternContext) GetParser() antlr.Parser { return s.parser }

func (s *PatternContext) UnaryExpr() IUnaryExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnaryExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnaryExprContext)
}

func (s *PatternContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *PatternContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
}

func (s *PatternContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *PatternContext) FieldPattern() IFieldPatternContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldPatternContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldPatternContext)
}

func (s *PatternContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *PatternContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, 0)
}

func (s *PatternContext) AllPattern() []IPatternContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPatternContext); ok {
			len++
		}
	}

	tst := make([]IPatternContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPatternContext); ok {
			tst[i] = t.(IPatternContext)
			i++
		}
	}

	return tst
}

func (s *PatternContext) Pattern(i int) IPatternContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPatternContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPatternContext)
}

func (s *PatternContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, 0)
}

func (s *PatternContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *PatternContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *PatternContext) COLON() antlr.TerminalNode {
	return s.GetToken(ospreyParserCOLON, 0)
}

func (s *PatternContext) Type_() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *PatternContext) UNDERSCORE() antlr.TerminalNode {
	return s.GetToken(ospreyParserUNDERSCORE, 0)
}

func (s *PatternContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PatternContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterPattern(s)
	}
}

func (s *PatternContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitPattern(s)
	}
}

func (p *ospreyParser) Pattern() (localctx IPatternContext) {
	localctx = NewPatternContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, ospreyParserRULE_pattern)
	var _la int

	p.SetState(741)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 77, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(701)
			p.UnaryExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(702)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(707)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLBRACE {
			{
				p.SetState(703)
				p.Match(ospreyParserLBRACE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(704)
				p.FieldPattern()
			}
			{
				p.SetState(705)
				p.Match(ospreyParserRBRACE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(709)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(721)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLPAREN {
			{
				p.SetState(710)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(711)
				p.Pattern()
			}
			p.SetState(716)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == ospreyParserCOMMA {
				{
					p.SetState(712)
					p.Match(ospreyParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(713)
					p.Pattern()
				}

				p.SetState(718)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(719)
				p.Match(ospreyParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(723)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(725)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(724)
				p.Match(ospreyParserID)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(727)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(728)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(729)
			p.Type_()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(730)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(731)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(732)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(733)
			p.FieldPattern()
		}
		{
			p.SetState(734)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(736)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(737)
			p.FieldPattern()
		}
		{
			p.SetState(738)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(740)
			p.Match(ospreyParserUNDERSCORE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldPatternContext is an interface to support dynamic dispatch.
type IFieldPatternContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFieldPatternContext differentiates from other interfaces.
	IsFieldPatternContext()
}

type FieldPatternContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldPatternContext() *FieldPatternContext {
	var p = new(FieldPatternContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldPattern
	return p
}

func InitEmptyFieldPatternContext(p *FieldPatternContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_fieldPattern
}

func (*FieldPatternContext) IsFieldPatternContext() {}

func NewFieldPatternContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldPatternContext {
	var p = new(FieldPatternContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_fieldPattern

	return p
}

func (s *FieldPatternContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldPatternContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *FieldPatternContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
}

func (s *FieldPatternContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserCOMMA)
}

func (s *FieldPatternContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserCOMMA, i)
}

func (s *FieldPatternContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldPatternContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldPatternContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterFieldPattern(s)
	}
}

func (s *FieldPatternContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitFieldPattern(s)
	}
}

func (p *ospreyParser) FieldPattern() (localctx IFieldPatternContext) {
	localctx = NewFieldPatternContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, ospreyParserRULE_fieldPattern)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(743)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(748)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(744)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(745)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(750)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBlockBodyContext is an interface to support dynamic dispatch.
type IBlockBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllStatement() []IStatementContext
	Statement(i int) IStatementContext
	Expr() IExprContext

	// IsBlockBodyContext differentiates from other interfaces.
	IsBlockBodyContext()
}

type BlockBodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBlockBodyContext() *BlockBodyContext {
	var p = new(BlockBodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_blockBody
	return p
}

func InitEmptyBlockBodyContext(p *BlockBodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_blockBody
}

func (*BlockBodyContext) IsBlockBodyContext() {}

func NewBlockBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlockBodyContext {
	var p = new(BlockBodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_blockBody

	return p
}

func (s *BlockBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *BlockBodyContext) AllStatement() []IStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatementContext); ok {
			len++
		}
	}

	tst := make([]IStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatementContext); ok {
			tst[i] = t.(IStatementContext)
			i++
		}
	}

	return tst
}

func (s *BlockBodyContext) Statement(i int) IStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *BlockBodyContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *BlockBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BlockBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BlockBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterBlockBody(s)
	}
}

func (s *BlockBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitBlockBody(s)
	}
}

func (p *ospreyParser) BlockBody() (localctx IBlockBodyContext) {
	localctx = NewBlockBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, ospreyParserRULE_blockBody)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(754)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 79, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(751)
				p.Statement()
			}

		}
		p.SetState(756)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 79, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(758)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&547930693369225266) != 0 {
		{
			p.SetState(757)
			p.Expr()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
