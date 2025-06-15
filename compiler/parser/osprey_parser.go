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
		"", "'|>'", "'match'", "'fn'", "'extern'", "'import'", "'type'", "'module'",
		"'let'", "'mut'", "'if'", "'else'", "'loop'", "'spawn'", "'yield'",
		"'await'", "'fiber'", "'channel'", "'send'", "'recv'", "'select'", "'true'",
		"'false'", "'where'", "'as'", "'->'", "'=>'", "'_'", "'='", "'=='",
		"'!='", "'<='", "'>='", "'!'", "'%'", "':'", "';'", "','", "'.'", "'|'",
		"'<'", "'>'", "'('", "')'", "'{'", "'}'", "'['", "']'", "'+'", "'-'",
		"'*'", "'/'",
	}
	staticData.SymbolicNames = []string{
		"", "PIPE", "MATCH", "FN", "EXTERN", "IMPORT", "TYPE", "MODULE", "LET",
		"MUT", "IF", "ELSE", "LOOP", "SPAWN", "YIELD", "AWAIT", "FIBER", "CHANNEL",
		"SEND", "RECV", "SELECT", "TRUE", "FALSE", "WHERE", "AS", "ARROW", "LAMBDA",
		"UNDERSCORE", "EQ", "EQ_OP", "NE_OP", "LE_OP", "GE_OP", "NOT_OP", "MOD_OP",
		"COLON", "SEMI", "COMMA", "DOT", "BAR", "LT", "GT", "LPAREN", "RPAREN",
		"LBRACE", "RBRACE", "LSQUARE", "RSQUARE", "PLUS", "MINUS", "STAR", "SLASH",
		"INT", "INTERPOLATED_STRING", "STRING", "ID", "WS", "DOC_COMMENT", "COMMENT",
		"PLUGIN_CONTENT",
	}
	staticData.RuleNames = []string{
		"program", "statement", "importStmt", "letDecl", "fnDecl", "pluginContent",
		"pluginName", "pluginReturnType", "externDecl", "externParamList", "externParam",
		"paramList", "param", "typeDecl", "typeParamList", "unionType", "recordType",
		"variant", "fieldDeclarations", "fieldDeclaration", "constraint", "functionCall",
		"booleanExpr", "fieldList", "field", "type", "typeList", "exprStmt",
		"expr", "matchExpr", "selectExpr", "selectArm", "binaryExpr", "comparisonExpr",
		"addExpr", "mulExpr", "unaryExpr", "pipeExpr", "callExpr", "argList",
		"namedArgList", "namedArg", "primary", "typeConstructor", "typeArgs",
		"fieldAssignments", "fieldAssignment", "lambdaExpr", "updateExpr", "blockExpr",
		"literal", "docComment", "moduleDecl", "moduleBody", "moduleStatement",
		"matchArm", "pattern", "fieldPattern", "blockBody",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 59, 671, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
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
		2, 58, 7, 58, 1, 0, 5, 0, 120, 8, 0, 10, 0, 12, 0, 123, 9, 0, 1, 0, 1,
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 134, 8, 1, 1, 2, 1,
		2, 1, 2, 1, 2, 5, 2, 140, 8, 2, 10, 2, 12, 2, 143, 9, 2, 1, 3, 1, 3, 1,
		3, 1, 3, 3, 3, 149, 8, 3, 1, 3, 1, 3, 1, 3, 1, 4, 3, 4, 155, 8, 4, 1, 4,
		1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 162, 8, 4, 1, 4, 1, 4, 3, 4, 166, 8, 4, 1,
		4, 1, 4, 1, 4, 1, 4, 3, 4, 172, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 178,
		8, 4, 1, 4, 1, 4, 1, 4, 3, 4, 183, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4,
		1, 4, 3, 4, 191, 8, 4, 3, 4, 193, 8, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1,
		7, 1, 7, 1, 7, 3, 7, 203, 8, 7, 1, 8, 3, 8, 206, 8, 8, 1, 8, 1, 8, 1, 8,
		1, 8, 1, 8, 3, 8, 213, 8, 8, 1, 8, 1, 8, 1, 8, 3, 8, 218, 8, 8, 1, 9, 1,
		9, 1, 9, 5, 9, 223, 8, 9, 10, 9, 12, 9, 226, 9, 9, 1, 10, 1, 10, 1, 10,
		1, 10, 1, 11, 1, 11, 1, 11, 5, 11, 235, 8, 11, 10, 11, 12, 11, 238, 9,
		11, 1, 12, 1, 12, 1, 12, 3, 12, 243, 8, 12, 1, 13, 3, 13, 246, 8, 13, 1,
		13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13, 254, 8, 13, 1, 13, 1, 13,
		1, 13, 3, 13, 259, 8, 13, 1, 14, 1, 14, 1, 14, 5, 14, 264, 8, 14, 10, 14,
		12, 14, 267, 9, 14, 1, 15, 1, 15, 1, 15, 5, 15, 272, 8, 15, 10, 15, 12,
		15, 275, 9, 15, 1, 16, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 1, 17,
		1, 17, 3, 17, 286, 8, 17, 1, 18, 1, 18, 1, 18, 5, 18, 291, 8, 18, 10, 18,
		12, 18, 294, 9, 18, 1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 300, 8, 19, 1, 20,
		1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 3, 21, 308, 8, 21, 1, 21, 1, 21, 1,
		22, 1, 22, 1, 23, 1, 23, 1, 23, 5, 23, 317, 8, 23, 10, 23, 12, 23, 320,
		9, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3,
		25, 331, 8, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25,
		340, 8, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 347, 8, 25, 1, 26,
		1, 26, 1, 26, 5, 26, 352, 8, 26, 10, 26, 12, 26, 355, 9, 26, 1, 27, 1,
		27, 1, 28, 1, 28, 1, 29, 1, 29, 1, 29, 1, 29, 4, 29, 365, 8, 29, 11, 29,
		12, 29, 366, 1, 29, 1, 29, 1, 29, 1, 29, 3, 29, 373, 8, 29, 1, 30, 1, 30,
		1, 30, 4, 30, 378, 8, 30, 11, 30, 12, 30, 379, 1, 30, 1, 30, 1, 31, 1,
		31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 3, 31, 391, 8, 31, 1, 32, 1, 32,
		1, 33, 1, 33, 1, 33, 5, 33, 398, 8, 33, 10, 33, 12, 33, 401, 9, 33, 1,
		34, 1, 34, 1, 34, 5, 34, 406, 8, 34, 10, 34, 12, 34, 409, 9, 34, 1, 35,
		1, 35, 1, 35, 5, 35, 414, 8, 35, 10, 35, 12, 35, 417, 9, 35, 1, 36, 3,
		36, 420, 8, 36, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 5, 37, 427, 8, 37, 10,
		37, 12, 37, 430, 9, 37, 1, 38, 1, 38, 1, 38, 4, 38, 435, 8, 38, 11, 38,
		12, 38, 436, 1, 38, 1, 38, 3, 38, 441, 8, 38, 1, 38, 3, 38, 444, 8, 38,
		1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 3, 38, 451, 8, 38, 1, 38, 4, 38, 454,
		8, 38, 11, 38, 12, 38, 455, 1, 38, 1, 38, 1, 38, 3, 38, 461, 8, 38, 1,
		38, 3, 38, 464, 8, 38, 3, 38, 466, 8, 38, 1, 39, 1, 39, 1, 39, 1, 39, 5,
		39, 472, 8, 39, 10, 39, 12, 39, 475, 9, 39, 3, 39, 477, 8, 39, 1, 40, 1,
		40, 1, 40, 4, 40, 482, 8, 40, 11, 40, 12, 40, 483, 1, 41, 1, 41, 1, 41,
		1, 41, 1, 42, 1, 42, 1, 42, 1, 42, 3, 42, 494, 8, 42, 1, 42, 1, 42, 1,
		42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42,
		1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1,
		42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 3, 42, 525, 8, 42, 1, 43, 1, 43,
		3, 43, 529, 8, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1, 44, 1,
		44, 1, 45, 1, 45, 1, 45, 5, 45, 542, 8, 45, 10, 45, 12, 45, 545, 9, 45,
		1, 46, 1, 46, 1, 46, 1, 46, 1, 47, 1, 47, 1, 47, 3, 47, 554, 8, 47, 1,
		47, 1, 47, 1, 47, 3, 47, 559, 8, 47, 1, 47, 1, 47, 1, 47, 1, 47, 3, 47,
		565, 8, 47, 1, 47, 1, 47, 1, 47, 3, 47, 570, 8, 47, 1, 48, 1, 48, 1, 48,
		1, 48, 1, 48, 1, 49, 1, 49, 1, 49, 1, 49, 1, 50, 1, 50, 1, 51, 4, 51, 584,
		8, 51, 11, 51, 12, 51, 585, 1, 52, 3, 52, 589, 8, 52, 1, 52, 1, 52, 1,
		52, 1, 52, 1, 52, 1, 52, 1, 53, 5, 53, 598, 8, 53, 10, 53, 12, 53, 601,
		9, 53, 1, 54, 1, 54, 1, 54, 3, 54, 606, 8, 54, 1, 55, 1, 55, 1, 55, 1,
		55, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 3, 56, 618, 8, 56, 1, 56,
		1, 56, 1, 56, 1, 56, 1, 56, 5, 56, 625, 8, 56, 10, 56, 12, 56, 628, 9,
		56, 1, 56, 1, 56, 3, 56, 632, 8, 56, 1, 56, 1, 56, 3, 56, 636, 8, 56, 1,
		56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56,
		1, 56, 1, 56, 1, 56, 3, 56, 652, 8, 56, 1, 57, 1, 57, 1, 57, 5, 57, 657,
		8, 57, 10, 57, 12, 57, 660, 9, 57, 1, 58, 5, 58, 663, 8, 58, 10, 58, 12,
		58, 666, 9, 58, 1, 58, 3, 58, 669, 8, 58, 1, 58, 0, 0, 59, 0, 2, 4, 6,
		8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
		44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78,
		80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112,
		114, 116, 0, 6, 1, 0, 8, 9, 2, 0, 29, 32, 40, 41, 1, 0, 48, 49, 2, 0, 34,
		34, 50, 51, 3, 0, 15, 15, 33, 33, 48, 49, 2, 0, 21, 22, 52, 54, 710, 0,
		121, 1, 0, 0, 0, 2, 133, 1, 0, 0, 0, 4, 135, 1, 0, 0, 0, 6, 144, 1, 0,
		0, 0, 8, 192, 1, 0, 0, 0, 10, 194, 1, 0, 0, 0, 12, 196, 1, 0, 0, 0, 14,
		202, 1, 0, 0, 0, 16, 205, 1, 0, 0, 0, 18, 219, 1, 0, 0, 0, 20, 227, 1,
		0, 0, 0, 22, 231, 1, 0, 0, 0, 24, 239, 1, 0, 0, 0, 26, 245, 1, 0, 0, 0,
		28, 260, 1, 0, 0, 0, 30, 268, 1, 0, 0, 0, 32, 276, 1, 0, 0, 0, 34, 280,
		1, 0, 0, 0, 36, 287, 1, 0, 0, 0, 38, 295, 1, 0, 0, 0, 40, 301, 1, 0, 0,
		0, 42, 304, 1, 0, 0, 0, 44, 311, 1, 0, 0, 0, 46, 313, 1, 0, 0, 0, 48, 321,
		1, 0, 0, 0, 50, 346, 1, 0, 0, 0, 52, 348, 1, 0, 0, 0, 54, 356, 1, 0, 0,
		0, 56, 358, 1, 0, 0, 0, 58, 372, 1, 0, 0, 0, 60, 374, 1, 0, 0, 0, 62, 390,
		1, 0, 0, 0, 64, 392, 1, 0, 0, 0, 66, 394, 1, 0, 0, 0, 68, 402, 1, 0, 0,
		0, 70, 410, 1, 0, 0, 0, 72, 419, 1, 0, 0, 0, 74, 423, 1, 0, 0, 0, 76, 465,
		1, 0, 0, 0, 78, 476, 1, 0, 0, 0, 80, 478, 1, 0, 0, 0, 82, 485, 1, 0, 0,
		0, 84, 524, 1, 0, 0, 0, 86, 526, 1, 0, 0, 0, 88, 534, 1, 0, 0, 0, 90, 538,
		1, 0, 0, 0, 92, 546, 1, 0, 0, 0, 94, 569, 1, 0, 0, 0, 96, 571, 1, 0, 0,
		0, 98, 576, 1, 0, 0, 0, 100, 580, 1, 0, 0, 0, 102, 583, 1, 0, 0, 0, 104,
		588, 1, 0, 0, 0, 106, 599, 1, 0, 0, 0, 108, 605, 1, 0, 0, 0, 110, 607,
		1, 0, 0, 0, 112, 651, 1, 0, 0, 0, 114, 653, 1, 0, 0, 0, 116, 664, 1, 0,
		0, 0, 118, 120, 3, 2, 1, 0, 119, 118, 1, 0, 0, 0, 120, 123, 1, 0, 0, 0,
		121, 119, 1, 0, 0, 0, 121, 122, 1, 0, 0, 0, 122, 124, 1, 0, 0, 0, 123,
		121, 1, 0, 0, 0, 124, 125, 5, 0, 0, 1, 125, 1, 1, 0, 0, 0, 126, 134, 3,
		4, 2, 0, 127, 134, 3, 6, 3, 0, 128, 134, 3, 8, 4, 0, 129, 134, 3, 16, 8,
		0, 130, 134, 3, 26, 13, 0, 131, 134, 3, 104, 52, 0, 132, 134, 3, 54, 27,
		0, 133, 126, 1, 0, 0, 0, 133, 127, 1, 0, 0, 0, 133, 128, 1, 0, 0, 0, 133,
		129, 1, 0, 0, 0, 133, 130, 1, 0, 0, 0, 133, 131, 1, 0, 0, 0, 133, 132,
		1, 0, 0, 0, 134, 3, 1, 0, 0, 0, 135, 136, 5, 5, 0, 0, 136, 141, 5, 55,
		0, 0, 137, 138, 5, 38, 0, 0, 138, 140, 5, 55, 0, 0, 139, 137, 1, 0, 0,
		0, 140, 143, 1, 0, 0, 0, 141, 139, 1, 0, 0, 0, 141, 142, 1, 0, 0, 0, 142,
		5, 1, 0, 0, 0, 143, 141, 1, 0, 0, 0, 144, 145, 7, 0, 0, 0, 145, 148, 5,
		55, 0, 0, 146, 147, 5, 35, 0, 0, 147, 149, 3, 50, 25, 0, 148, 146, 1, 0,
		0, 0, 148, 149, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150, 151, 5, 28, 0, 0,
		151, 152, 3, 56, 28, 0, 152, 7, 1, 0, 0, 0, 153, 155, 3, 102, 51, 0, 154,
		153, 1, 0, 0, 0, 154, 155, 1, 0, 0, 0, 155, 156, 1, 0, 0, 0, 156, 157,
		5, 3, 0, 0, 157, 158, 3, 12, 6, 0, 158, 159, 5, 55, 0, 0, 159, 161, 5,
		42, 0, 0, 160, 162, 3, 22, 11, 0, 161, 160, 1, 0, 0, 0, 161, 162, 1, 0,
		0, 0, 162, 163, 1, 0, 0, 0, 163, 165, 5, 43, 0, 0, 164, 166, 3, 14, 7,
		0, 165, 164, 1, 0, 0, 0, 165, 166, 1, 0, 0, 0, 166, 167, 1, 0, 0, 0, 167,
		168, 5, 28, 0, 0, 168, 169, 3, 10, 5, 0, 169, 193, 1, 0, 0, 0, 170, 172,
		3, 102, 51, 0, 171, 170, 1, 0, 0, 0, 171, 172, 1, 0, 0, 0, 172, 173, 1,
		0, 0, 0, 173, 174, 5, 3, 0, 0, 174, 175, 5, 55, 0, 0, 175, 177, 5, 42,
		0, 0, 176, 178, 3, 22, 11, 0, 177, 176, 1, 0, 0, 0, 177, 178, 1, 0, 0,
		0, 178, 179, 1, 0, 0, 0, 179, 182, 5, 43, 0, 0, 180, 181, 5, 25, 0, 0,
		181, 183, 3, 50, 25, 0, 182, 180, 1, 0, 0, 0, 182, 183, 1, 0, 0, 0, 183,
		190, 1, 0, 0, 0, 184, 185, 5, 28, 0, 0, 185, 191, 3, 56, 28, 0, 186, 187,
		5, 44, 0, 0, 187, 188, 3, 116, 58, 0, 188, 189, 5, 45, 0, 0, 189, 191,
		1, 0, 0, 0, 190, 184, 1, 0, 0, 0, 190, 186, 1, 0, 0, 0, 191, 193, 1, 0,
		0, 0, 192, 154, 1, 0, 0, 0, 192, 171, 1, 0, 0, 0, 193, 9, 1, 0, 0, 0, 194,
		195, 5, 59, 0, 0, 195, 11, 1, 0, 0, 0, 196, 197, 5, 55, 0, 0, 197, 13,
		1, 0, 0, 0, 198, 199, 5, 25, 0, 0, 199, 203, 3, 50, 25, 0, 200, 201, 5,
		24, 0, 0, 201, 203, 3, 50, 25, 0, 202, 198, 1, 0, 0, 0, 202, 200, 1, 0,
		0, 0, 203, 15, 1, 0, 0, 0, 204, 206, 3, 102, 51, 0, 205, 204, 1, 0, 0,
		0, 205, 206, 1, 0, 0, 0, 206, 207, 1, 0, 0, 0, 207, 208, 5, 4, 0, 0, 208,
		209, 5, 3, 0, 0, 209, 210, 5, 55, 0, 0, 210, 212, 5, 42, 0, 0, 211, 213,
		3, 18, 9, 0, 212, 211, 1, 0, 0, 0, 212, 213, 1, 0, 0, 0, 213, 214, 1, 0,
		0, 0, 214, 217, 5, 43, 0, 0, 215, 216, 5, 25, 0, 0, 216, 218, 3, 50, 25,
		0, 217, 215, 1, 0, 0, 0, 217, 218, 1, 0, 0, 0, 218, 17, 1, 0, 0, 0, 219,
		224, 3, 20, 10, 0, 220, 221, 5, 37, 0, 0, 221, 223, 3, 20, 10, 0, 222,
		220, 1, 0, 0, 0, 223, 226, 1, 0, 0, 0, 224, 222, 1, 0, 0, 0, 224, 225,
		1, 0, 0, 0, 225, 19, 1, 0, 0, 0, 226, 224, 1, 0, 0, 0, 227, 228, 5, 55,
		0, 0, 228, 229, 5, 35, 0, 0, 229, 230, 3, 50, 25, 0, 230, 21, 1, 0, 0,
		0, 231, 236, 3, 24, 12, 0, 232, 233, 5, 37, 0, 0, 233, 235, 3, 24, 12,
		0, 234, 232, 1, 0, 0, 0, 235, 238, 1, 0, 0, 0, 236, 234, 1, 0, 0, 0, 236,
		237, 1, 0, 0, 0, 237, 23, 1, 0, 0, 0, 238, 236, 1, 0, 0, 0, 239, 242, 5,
		55, 0, 0, 240, 241, 5, 35, 0, 0, 241, 243, 3, 50, 25, 0, 242, 240, 1, 0,
		0, 0, 242, 243, 1, 0, 0, 0, 243, 25, 1, 0, 0, 0, 244, 246, 3, 102, 51,
		0, 245, 244, 1, 0, 0, 0, 245, 246, 1, 0, 0, 0, 246, 247, 1, 0, 0, 0, 247,
		248, 5, 6, 0, 0, 248, 253, 5, 55, 0, 0, 249, 250, 5, 40, 0, 0, 250, 251,
		3, 28, 14, 0, 251, 252, 5, 41, 0, 0, 252, 254, 1, 0, 0, 0, 253, 249, 1,
		0, 0, 0, 253, 254, 1, 0, 0, 0, 254, 255, 1, 0, 0, 0, 255, 258, 5, 28, 0,
		0, 256, 259, 3, 30, 15, 0, 257, 259, 3, 32, 16, 0, 258, 256, 1, 0, 0, 0,
		258, 257, 1, 0, 0, 0, 259, 27, 1, 0, 0, 0, 260, 265, 5, 55, 0, 0, 261,
		262, 5, 37, 0, 0, 262, 264, 5, 55, 0, 0, 263, 261, 1, 0, 0, 0, 264, 267,
		1, 0, 0, 0, 265, 263, 1, 0, 0, 0, 265, 266, 1, 0, 0, 0, 266, 29, 1, 0,
		0, 0, 267, 265, 1, 0, 0, 0, 268, 273, 3, 34, 17, 0, 269, 270, 5, 39, 0,
		0, 270, 272, 3, 34, 17, 0, 271, 269, 1, 0, 0, 0, 272, 275, 1, 0, 0, 0,
		273, 271, 1, 0, 0, 0, 273, 274, 1, 0, 0, 0, 274, 31, 1, 0, 0, 0, 275, 273,
		1, 0, 0, 0, 276, 277, 5, 44, 0, 0, 277, 278, 3, 36, 18, 0, 278, 279, 5,
		45, 0, 0, 279, 33, 1, 0, 0, 0, 280, 285, 5, 55, 0, 0, 281, 282, 5, 44,
		0, 0, 282, 283, 3, 36, 18, 0, 283, 284, 5, 45, 0, 0, 284, 286, 1, 0, 0,
		0, 285, 281, 1, 0, 0, 0, 285, 286, 1, 0, 0, 0, 286, 35, 1, 0, 0, 0, 287,
		292, 3, 38, 19, 0, 288, 289, 5, 37, 0, 0, 289, 291, 3, 38, 19, 0, 290,
		288, 1, 0, 0, 0, 291, 294, 1, 0, 0, 0, 292, 290, 1, 0, 0, 0, 292, 293,
		1, 0, 0, 0, 293, 37, 1, 0, 0, 0, 294, 292, 1, 0, 0, 0, 295, 296, 5, 55,
		0, 0, 296, 297, 5, 35, 0, 0, 297, 299, 3, 50, 25, 0, 298, 300, 3, 40, 20,
		0, 299, 298, 1, 0, 0, 0, 299, 300, 1, 0, 0, 0, 300, 39, 1, 0, 0, 0, 301,
		302, 5, 23, 0, 0, 302, 303, 3, 42, 21, 0, 303, 41, 1, 0, 0, 0, 304, 305,
		5, 55, 0, 0, 305, 307, 5, 42, 0, 0, 306, 308, 3, 78, 39, 0, 307, 306, 1,
		0, 0, 0, 307, 308, 1, 0, 0, 0, 308, 309, 1, 0, 0, 0, 309, 310, 5, 43, 0,
		0, 310, 43, 1, 0, 0, 0, 311, 312, 3, 66, 33, 0, 312, 45, 1, 0, 0, 0, 313,
		318, 3, 48, 24, 0, 314, 315, 5, 37, 0, 0, 315, 317, 3, 48, 24, 0, 316,
		314, 1, 0, 0, 0, 317, 320, 1, 0, 0, 0, 318, 316, 1, 0, 0, 0, 318, 319,
		1, 0, 0, 0, 319, 47, 1, 0, 0, 0, 320, 318, 1, 0, 0, 0, 321, 322, 5, 55,
		0, 0, 322, 323, 5, 35, 0, 0, 323, 324, 3, 50, 25, 0, 324, 49, 1, 0, 0,
		0, 325, 330, 5, 55, 0, 0, 326, 327, 5, 40, 0, 0, 327, 328, 3, 52, 26, 0,
		328, 329, 5, 41, 0, 0, 329, 331, 1, 0, 0, 0, 330, 326, 1, 0, 0, 0, 330,
		331, 1, 0, 0, 0, 331, 332, 1, 0, 0, 0, 332, 333, 5, 46, 0, 0, 333, 347,
		5, 47, 0, 0, 334, 339, 5, 55, 0, 0, 335, 336, 5, 40, 0, 0, 336, 337, 3,
		52, 26, 0, 337, 338, 5, 41, 0, 0, 338, 340, 1, 0, 0, 0, 339, 335, 1, 0,
		0, 0, 339, 340, 1, 0, 0, 0, 340, 347, 1, 0, 0, 0, 341, 342, 5, 55, 0, 0,
		342, 343, 5, 46, 0, 0, 343, 344, 3, 50, 25, 0, 344, 345, 5, 47, 0, 0, 345,
		347, 1, 0, 0, 0, 346, 325, 1, 0, 0, 0, 346, 334, 1, 0, 0, 0, 346, 341,
		1, 0, 0, 0, 347, 51, 1, 0, 0, 0, 348, 353, 3, 50, 25, 0, 349, 350, 5, 37,
		0, 0, 350, 352, 3, 50, 25, 0, 351, 349, 1, 0, 0, 0, 352, 355, 1, 0, 0,
		0, 353, 351, 1, 0, 0, 0, 353, 354, 1, 0, 0, 0, 354, 53, 1, 0, 0, 0, 355,
		353, 1, 0, 0, 0, 356, 357, 3, 56, 28, 0, 357, 55, 1, 0, 0, 0, 358, 359,
		3, 58, 29, 0, 359, 57, 1, 0, 0, 0, 360, 361, 5, 2, 0, 0, 361, 362, 3, 56,
		28, 0, 362, 364, 5, 44, 0, 0, 363, 365, 3, 110, 55, 0, 364, 363, 1, 0,
		0, 0, 365, 366, 1, 0, 0, 0, 366, 364, 1, 0, 0, 0, 366, 367, 1, 0, 0, 0,
		367, 368, 1, 0, 0, 0, 368, 369, 5, 45, 0, 0, 369, 373, 1, 0, 0, 0, 370,
		373, 3, 60, 30, 0, 371, 373, 3, 64, 32, 0, 372, 360, 1, 0, 0, 0, 372, 370,
		1, 0, 0, 0, 372, 371, 1, 0, 0, 0, 373, 59, 1, 0, 0, 0, 374, 375, 5, 20,
		0, 0, 375, 377, 5, 44, 0, 0, 376, 378, 3, 62, 31, 0, 377, 376, 1, 0, 0,
		0, 378, 379, 1, 0, 0, 0, 379, 377, 1, 0, 0, 0, 379, 380, 1, 0, 0, 0, 380,
		381, 1, 0, 0, 0, 381, 382, 5, 45, 0, 0, 382, 61, 1, 0, 0, 0, 383, 384,
		3, 112, 56, 0, 384, 385, 5, 26, 0, 0, 385, 386, 3, 56, 28, 0, 386, 391,
		1, 0, 0, 0, 387, 388, 5, 27, 0, 0, 388, 389, 5, 26, 0, 0, 389, 391, 3,
		56, 28, 0, 390, 383, 1, 0, 0, 0, 390, 387, 1, 0, 0, 0, 391, 63, 1, 0, 0,
		0, 392, 393, 3, 66, 33, 0, 393, 65, 1, 0, 0, 0, 394, 399, 3, 68, 34, 0,
		395, 396, 7, 1, 0, 0, 396, 398, 3, 68, 34, 0, 397, 395, 1, 0, 0, 0, 398,
		401, 1, 0, 0, 0, 399, 397, 1, 0, 0, 0, 399, 400, 1, 0, 0, 0, 400, 67, 1,
		0, 0, 0, 401, 399, 1, 0, 0, 0, 402, 407, 3, 70, 35, 0, 403, 404, 7, 2,
		0, 0, 404, 406, 3, 70, 35, 0, 405, 403, 1, 0, 0, 0, 406, 409, 1, 0, 0,
		0, 407, 405, 1, 0, 0, 0, 407, 408, 1, 0, 0, 0, 408, 69, 1, 0, 0, 0, 409,
		407, 1, 0, 0, 0, 410, 415, 3, 72, 36, 0, 411, 412, 7, 3, 0, 0, 412, 414,
		3, 72, 36, 0, 413, 411, 1, 0, 0, 0, 414, 417, 1, 0, 0, 0, 415, 413, 1,
		0, 0, 0, 415, 416, 1, 0, 0, 0, 416, 71, 1, 0, 0, 0, 417, 415, 1, 0, 0,
		0, 418, 420, 7, 4, 0, 0, 419, 418, 1, 0, 0, 0, 419, 420, 1, 0, 0, 0, 420,
		421, 1, 0, 0, 0, 421, 422, 3, 74, 37, 0, 422, 73, 1, 0, 0, 0, 423, 428,
		3, 76, 38, 0, 424, 425, 5, 1, 0, 0, 425, 427, 3, 76, 38, 0, 426, 424, 1,
		0, 0, 0, 427, 430, 1, 0, 0, 0, 428, 426, 1, 0, 0, 0, 428, 429, 1, 0, 0,
		0, 429, 75, 1, 0, 0, 0, 430, 428, 1, 0, 0, 0, 431, 434, 3, 84, 42, 0, 432,
		433, 5, 38, 0, 0, 433, 435, 5, 55, 0, 0, 434, 432, 1, 0, 0, 0, 435, 436,
		1, 0, 0, 0, 436, 434, 1, 0, 0, 0, 436, 437, 1, 0, 0, 0, 437, 443, 1, 0,
		0, 0, 438, 440, 5, 42, 0, 0, 439, 441, 3, 78, 39, 0, 440, 439, 1, 0, 0,
		0, 440, 441, 1, 0, 0, 0, 441, 442, 1, 0, 0, 0, 442, 444, 5, 43, 0, 0, 443,
		438, 1, 0, 0, 0, 443, 444, 1, 0, 0, 0, 444, 466, 1, 0, 0, 0, 445, 453,
		3, 84, 42, 0, 446, 447, 5, 38, 0, 0, 447, 448, 5, 55, 0, 0, 448, 450, 5,
		42, 0, 0, 449, 451, 3, 78, 39, 0, 450, 449, 1, 0, 0, 0, 450, 451, 1, 0,
		0, 0, 451, 452, 1, 0, 0, 0, 452, 454, 5, 43, 0, 0, 453, 446, 1, 0, 0, 0,
		454, 455, 1, 0, 0, 0, 455, 453, 1, 0, 0, 0, 455, 456, 1, 0, 0, 0, 456,
		466, 1, 0, 0, 0, 457, 463, 3, 84, 42, 0, 458, 460, 5, 42, 0, 0, 459, 461,
		3, 78, 39, 0, 460, 459, 1, 0, 0, 0, 460, 461, 1, 0, 0, 0, 461, 462, 1,
		0, 0, 0, 462, 464, 5, 43, 0, 0, 463, 458, 1, 0, 0, 0, 463, 464, 1, 0, 0,
		0, 464, 466, 1, 0, 0, 0, 465, 431, 1, 0, 0, 0, 465, 445, 1, 0, 0, 0, 465,
		457, 1, 0, 0, 0, 466, 77, 1, 0, 0, 0, 467, 477, 3, 80, 40, 0, 468, 473,
		3, 56, 28, 0, 469, 470, 5, 37, 0, 0, 470, 472, 3, 56, 28, 0, 471, 469,
		1, 0, 0, 0, 472, 475, 1, 0, 0, 0, 473, 471, 1, 0, 0, 0, 473, 474, 1, 0,
		0, 0, 474, 477, 1, 0, 0, 0, 475, 473, 1, 0, 0, 0, 476, 467, 1, 0, 0, 0,
		476, 468, 1, 0, 0, 0, 477, 79, 1, 0, 0, 0, 478, 481, 3, 82, 41, 0, 479,
		480, 5, 37, 0, 0, 480, 482, 3, 82, 41, 0, 481, 479, 1, 0, 0, 0, 482, 483,
		1, 0, 0, 0, 483, 481, 1, 0, 0, 0, 483, 484, 1, 0, 0, 0, 484, 81, 1, 0,
		0, 0, 485, 486, 5, 55, 0, 0, 486, 487, 5, 35, 0, 0, 487, 488, 3, 56, 28,
		0, 488, 83, 1, 0, 0, 0, 489, 490, 5, 13, 0, 0, 490, 525, 3, 56, 28, 0,
		491, 493, 5, 14, 0, 0, 492, 494, 3, 56, 28, 0, 493, 492, 1, 0, 0, 0, 493,
		494, 1, 0, 0, 0, 494, 525, 1, 0, 0, 0, 495, 496, 5, 15, 0, 0, 496, 497,
		5, 42, 0, 0, 497, 498, 3, 56, 28, 0, 498, 499, 5, 43, 0, 0, 499, 525, 1,
		0, 0, 0, 500, 501, 5, 18, 0, 0, 501, 502, 5, 42, 0, 0, 502, 503, 3, 56,
		28, 0, 503, 504, 5, 37, 0, 0, 504, 505, 3, 56, 28, 0, 505, 506, 5, 43,
		0, 0, 506, 525, 1, 0, 0, 0, 507, 508, 5, 19, 0, 0, 508, 509, 5, 42, 0,
		0, 509, 510, 3, 56, 28, 0, 510, 511, 5, 43, 0, 0, 511, 525, 1, 0, 0, 0,
		512, 513, 5, 20, 0, 0, 513, 525, 3, 60, 30, 0, 514, 525, 3, 86, 43, 0,
		515, 525, 3, 96, 48, 0, 516, 525, 3, 98, 49, 0, 517, 525, 3, 100, 50, 0,
		518, 525, 3, 94, 47, 0, 519, 525, 5, 55, 0, 0, 520, 521, 5, 42, 0, 0, 521,
		522, 3, 56, 28, 0, 522, 523, 5, 43, 0, 0, 523, 525, 1, 0, 0, 0, 524, 489,
		1, 0, 0, 0, 524, 491, 1, 0, 0, 0, 524, 495, 1, 0, 0, 0, 524, 500, 1, 0,
		0, 0, 524, 507, 1, 0, 0, 0, 524, 512, 1, 0, 0, 0, 524, 514, 1, 0, 0, 0,
		524, 515, 1, 0, 0, 0, 524, 516, 1, 0, 0, 0, 524, 517, 1, 0, 0, 0, 524,
		518, 1, 0, 0, 0, 524, 519, 1, 0, 0, 0, 524, 520, 1, 0, 0, 0, 525, 85, 1,
		0, 0, 0, 526, 528, 5, 55, 0, 0, 527, 529, 3, 88, 44, 0, 528, 527, 1, 0,
		0, 0, 528, 529, 1, 0, 0, 0, 529, 530, 1, 0, 0, 0, 530, 531, 5, 44, 0, 0,
		531, 532, 3, 90, 45, 0, 532, 533, 5, 45, 0, 0, 533, 87, 1, 0, 0, 0, 534,
		535, 5, 40, 0, 0, 535, 536, 3, 52, 26, 0, 536, 537, 5, 41, 0, 0, 537, 89,
		1, 0, 0, 0, 538, 543, 3, 92, 46, 0, 539, 540, 5, 37, 0, 0, 540, 542, 3,
		92, 46, 0, 541, 539, 1, 0, 0, 0, 542, 545, 1, 0, 0, 0, 543, 541, 1, 0,
		0, 0, 543, 544, 1, 0, 0, 0, 544, 91, 1, 0, 0, 0, 545, 543, 1, 0, 0, 0,
		546, 547, 5, 55, 0, 0, 547, 548, 5, 35, 0, 0, 548, 549, 3, 56, 28, 0, 549,
		93, 1, 0, 0, 0, 550, 551, 5, 3, 0, 0, 551, 553, 5, 42, 0, 0, 552, 554,
		3, 22, 11, 0, 553, 552, 1, 0, 0, 0, 553, 554, 1, 0, 0, 0, 554, 555, 1,
		0, 0, 0, 555, 558, 5, 43, 0, 0, 556, 557, 5, 25, 0, 0, 557, 559, 3, 50,
		25, 0, 558, 556, 1, 0, 0, 0, 558, 559, 1, 0, 0, 0, 559, 560, 1, 0, 0, 0,
		560, 561, 5, 26, 0, 0, 561, 570, 3, 56, 28, 0, 562, 564, 5, 39, 0, 0, 563,
		565, 3, 22, 11, 0, 564, 563, 1, 0, 0, 0, 564, 565, 1, 0, 0, 0, 565, 566,
		1, 0, 0, 0, 566, 567, 5, 39, 0, 0, 567, 568, 5, 26, 0, 0, 568, 570, 3,
		56, 28, 0, 569, 550, 1, 0, 0, 0, 569, 562, 1, 0, 0, 0, 570, 95, 1, 0, 0,
		0, 571, 572, 5, 55, 0, 0, 572, 573, 5, 44, 0, 0, 573, 574, 3, 90, 45, 0,
		574, 575, 5, 45, 0, 0, 575, 97, 1, 0, 0, 0, 576, 577, 5, 44, 0, 0, 577,
		578, 3, 116, 58, 0, 578, 579, 5, 45, 0, 0, 579, 99, 1, 0, 0, 0, 580, 581,
		7, 5, 0, 0, 581, 101, 1, 0, 0, 0, 582, 584, 5, 57, 0, 0, 583, 582, 1, 0,
		0, 0, 584, 585, 1, 0, 0, 0, 585, 583, 1, 0, 0, 0, 585, 586, 1, 0, 0, 0,
		586, 103, 1, 0, 0, 0, 587, 589, 3, 102, 51, 0, 588, 587, 1, 0, 0, 0, 588,
		589, 1, 0, 0, 0, 589, 590, 1, 0, 0, 0, 590, 591, 5, 7, 0, 0, 591, 592,
		5, 55, 0, 0, 592, 593, 5, 44, 0, 0, 593, 594, 3, 106, 53, 0, 594, 595,
		5, 45, 0, 0, 595, 105, 1, 0, 0, 0, 596, 598, 3, 108, 54, 0, 597, 596, 1,
		0, 0, 0, 598, 601, 1, 0, 0, 0, 599, 597, 1, 0, 0, 0, 599, 600, 1, 0, 0,
		0, 600, 107, 1, 0, 0, 0, 601, 599, 1, 0, 0, 0, 602, 606, 3, 6, 3, 0, 603,
		606, 3, 8, 4, 0, 604, 606, 3, 26, 13, 0, 605, 602, 1, 0, 0, 0, 605, 603,
		1, 0, 0, 0, 605, 604, 1, 0, 0, 0, 606, 109, 1, 0, 0, 0, 607, 608, 3, 112,
		56, 0, 608, 609, 5, 26, 0, 0, 609, 610, 3, 56, 28, 0, 610, 111, 1, 0, 0,
		0, 611, 652, 3, 72, 36, 0, 612, 617, 5, 55, 0, 0, 613, 614, 5, 44, 0, 0,
		614, 615, 3, 114, 57, 0, 615, 616, 5, 45, 0, 0, 616, 618, 1, 0, 0, 0, 617,
		613, 1, 0, 0, 0, 617, 618, 1, 0, 0, 0, 618, 652, 1, 0, 0, 0, 619, 631,
		5, 55, 0, 0, 620, 621, 5, 42, 0, 0, 621, 626, 3, 112, 56, 0, 622, 623,
		5, 37, 0, 0, 623, 625, 3, 112, 56, 0, 624, 622, 1, 0, 0, 0, 625, 628, 1,
		0, 0, 0, 626, 624, 1, 0, 0, 0, 626, 627, 1, 0, 0, 0, 627, 629, 1, 0, 0,
		0, 628, 626, 1, 0, 0, 0, 629, 630, 5, 43, 0, 0, 630, 632, 1, 0, 0, 0, 631,
		620, 1, 0, 0, 0, 631, 632, 1, 0, 0, 0, 632, 652, 1, 0, 0, 0, 633, 635,
		5, 55, 0, 0, 634, 636, 5, 55, 0, 0, 635, 634, 1, 0, 0, 0, 635, 636, 1,
		0, 0, 0, 636, 652, 1, 0, 0, 0, 637, 638, 5, 55, 0, 0, 638, 639, 5, 35,
		0, 0, 639, 652, 3, 50, 25, 0, 640, 641, 5, 55, 0, 0, 641, 642, 5, 35, 0,
		0, 642, 643, 5, 44, 0, 0, 643, 644, 3, 114, 57, 0, 644, 645, 5, 45, 0,
		0, 645, 652, 1, 0, 0, 0, 646, 647, 5, 44, 0, 0, 647, 648, 3, 114, 57, 0,
		648, 649, 5, 45, 0, 0, 649, 652, 1, 0, 0, 0, 650, 652, 5, 27, 0, 0, 651,
		611, 1, 0, 0, 0, 651, 612, 1, 0, 0, 0, 651, 619, 1, 0, 0, 0, 651, 633,
		1, 0, 0, 0, 651, 637, 1, 0, 0, 0, 651, 640, 1, 0, 0, 0, 651, 646, 1, 0,
		0, 0, 651, 650, 1, 0, 0, 0, 652, 113, 1, 0, 0, 0, 653, 658, 5, 55, 0, 0,
		654, 655, 5, 37, 0, 0, 655, 657, 5, 55, 0, 0, 656, 654, 1, 0, 0, 0, 657,
		660, 1, 0, 0, 0, 658, 656, 1, 0, 0, 0, 658, 659, 1, 0, 0, 0, 659, 115,
		1, 0, 0, 0, 660, 658, 1, 0, 0, 0, 661, 663, 3, 2, 1, 0, 662, 661, 1, 0,
		0, 0, 663, 666, 1, 0, 0, 0, 664, 662, 1, 0, 0, 0, 664, 665, 1, 0, 0, 0,
		665, 668, 1, 0, 0, 0, 666, 664, 1, 0, 0, 0, 667, 669, 3, 56, 28, 0, 668,
		667, 1, 0, 0, 0, 668, 669, 1, 0, 0, 0, 669, 117, 1, 0, 0, 0, 73, 121, 133,
		141, 148, 154, 161, 165, 171, 177, 182, 190, 192, 202, 205, 212, 217, 224,
		236, 242, 245, 253, 258, 265, 273, 285, 292, 299, 307, 318, 330, 339, 346,
		353, 366, 372, 379, 390, 399, 407, 415, 419, 428, 436, 440, 443, 450, 455,
		460, 463, 465, 473, 476, 483, 493, 524, 528, 543, 553, 558, 564, 569, 585,
		588, 599, 605, 617, 626, 631, 635, 651, 658, 664, 668,
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
	ospreyParserPIPE                = 1
	ospreyParserMATCH               = 2
	ospreyParserFN                  = 3
	ospreyParserEXTERN              = 4
	ospreyParserIMPORT              = 5
	ospreyParserTYPE                = 6
	ospreyParserMODULE              = 7
	ospreyParserLET                 = 8
	ospreyParserMUT                 = 9
	ospreyParserIF                  = 10
	ospreyParserELSE                = 11
	ospreyParserLOOP                = 12
	ospreyParserSPAWN               = 13
	ospreyParserYIELD               = 14
	ospreyParserAWAIT               = 15
	ospreyParserFIBER               = 16
	ospreyParserCHANNEL             = 17
	ospreyParserSEND                = 18
	ospreyParserRECV                = 19
	ospreyParserSELECT              = 20
	ospreyParserTRUE                = 21
	ospreyParserFALSE               = 22
	ospreyParserWHERE               = 23
	ospreyParserAS                  = 24
	ospreyParserARROW               = 25
	ospreyParserLAMBDA              = 26
	ospreyParserUNDERSCORE          = 27
	ospreyParserEQ                  = 28
	ospreyParserEQ_OP               = 29
	ospreyParserNE_OP               = 30
	ospreyParserLE_OP               = 31
	ospreyParserGE_OP               = 32
	ospreyParserNOT_OP              = 33
	ospreyParserMOD_OP              = 34
	ospreyParserCOLON               = 35
	ospreyParserSEMI                = 36
	ospreyParserCOMMA               = 37
	ospreyParserDOT                 = 38
	ospreyParserBAR                 = 39
	ospreyParserLT                  = 40
	ospreyParserGT                  = 41
	ospreyParserLPAREN              = 42
	ospreyParserRPAREN              = 43
	ospreyParserLBRACE              = 44
	ospreyParserRBRACE              = 45
	ospreyParserLSQUARE             = 46
	ospreyParserRSQUARE             = 47
	ospreyParserPLUS                = 48
	ospreyParserMINUS               = 49
	ospreyParserSTAR                = 50
	ospreyParserSLASH               = 51
	ospreyParserINT                 = 52
	ospreyParserINTERPOLATED_STRING = 53
	ospreyParserSTRING              = 54
	ospreyParserID                  = 55
	ospreyParserWS                  = 56
	ospreyParserDOC_COMMENT         = 57
	ospreyParserCOMMENT             = 58
	ospreyParserPLUGIN_CONTENT      = 59
)

// ospreyParser rules.
const (
	ospreyParserRULE_program           = 0
	ospreyParserRULE_statement         = 1
	ospreyParserRULE_importStmt        = 2
	ospreyParserRULE_letDecl           = 3
	ospreyParserRULE_fnDecl            = 4
	ospreyParserRULE_pluginContent     = 5
	ospreyParserRULE_pluginName        = 6
	ospreyParserRULE_pluginReturnType  = 7
	ospreyParserRULE_externDecl        = 8
	ospreyParserRULE_externParamList   = 9
	ospreyParserRULE_externParam       = 10
	ospreyParserRULE_paramList         = 11
	ospreyParserRULE_param             = 12
	ospreyParserRULE_typeDecl          = 13
	ospreyParserRULE_typeParamList     = 14
	ospreyParserRULE_unionType         = 15
	ospreyParserRULE_recordType        = 16
	ospreyParserRULE_variant           = 17
	ospreyParserRULE_fieldDeclarations = 18
	ospreyParserRULE_fieldDeclaration  = 19
	ospreyParserRULE_constraint        = 20
	ospreyParserRULE_functionCall      = 21
	ospreyParserRULE_booleanExpr       = 22
	ospreyParserRULE_fieldList         = 23
	ospreyParserRULE_field             = 24
	ospreyParserRULE_type              = 25
	ospreyParserRULE_typeList          = 26
	ospreyParserRULE_exprStmt          = 27
	ospreyParserRULE_expr              = 28
	ospreyParserRULE_matchExpr         = 29
	ospreyParserRULE_selectExpr        = 30
	ospreyParserRULE_selectArm         = 31
	ospreyParserRULE_binaryExpr        = 32
	ospreyParserRULE_comparisonExpr    = 33
	ospreyParserRULE_addExpr           = 34
	ospreyParserRULE_mulExpr           = 35
	ospreyParserRULE_unaryExpr         = 36
	ospreyParserRULE_pipeExpr          = 37
	ospreyParserRULE_callExpr          = 38
	ospreyParserRULE_argList           = 39
	ospreyParserRULE_namedArgList      = 40
	ospreyParserRULE_namedArg          = 41
	ospreyParserRULE_primary           = 42
	ospreyParserRULE_typeConstructor   = 43
	ospreyParserRULE_typeArgs          = 44
	ospreyParserRULE_fieldAssignments  = 45
	ospreyParserRULE_fieldAssignment   = 46
	ospreyParserRULE_lambdaExpr        = 47
	ospreyParserRULE_updateExpr        = 48
	ospreyParserRULE_blockExpr         = 49
	ospreyParserRULE_literal           = 50
	ospreyParserRULE_docComment        = 51
	ospreyParserRULE_moduleDecl        = 52
	ospreyParserRULE_moduleBody        = 53
	ospreyParserRULE_moduleStatement   = 54
	ospreyParserRULE_matchArm          = 55
	ospreyParserRULE_pattern           = 56
	ospreyParserRULE_fieldPattern      = 57
	ospreyParserRULE_blockBody         = 58
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
	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&212536156003034108) != 0 {
		{
			p.SetState(118)
			p.Statement()
		}

		p.SetState(123)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(124)
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
	FnDecl() IFnDeclContext
	ExternDecl() IExternDeclContext
	TypeDecl() ITypeDeclContext
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
	p.SetState(133)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(126)
			p.ImportStmt()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(127)
			p.LetDecl()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(128)
			p.FnDecl()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(129)
			p.ExternDecl()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(130)
			p.TypeDecl()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(131)
			p.ModuleDecl()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(132)
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
		p.SetState(135)
		p.Match(ospreyParserIMPORT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(136)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(141)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserDOT {
		{
			p.SetState(137)
			p.Match(ospreyParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(138)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(143)
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
		p.SetState(144)
		_la = p.GetTokenStream().LA(1)

		if !(_la == ospreyParserLET || _la == ospreyParserMUT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(145)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(148)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserCOLON {
		{
			p.SetState(146)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(147)
			p.Type_()
		}

	}
	{
		p.SetState(150)
		p.Match(ospreyParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(151)
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
	PluginName() IPluginNameContext
	ID() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	EQ() antlr.TerminalNode
	PluginContent() IPluginContentContext
	DocComment() IDocCommentContext
	ParamList() IParamListContext
	PluginReturnType() IPluginReturnTypeContext
	Expr() IExprContext
	LBRACE() antlr.TerminalNode
	BlockBody() IBlockBodyContext
	RBRACE() antlr.TerminalNode
	ARROW() antlr.TerminalNode
	Type_() ITypeContext

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

func (s *FnDeclContext) PluginName() IPluginNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPluginNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPluginNameContext)
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

func (s *FnDeclContext) PluginContent() IPluginContentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPluginContentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPluginContentContext)
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

func (s *FnDeclContext) PluginReturnType() IPluginReturnTypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPluginReturnTypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPluginReturnTypeContext)
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
	p.EnterRule(localctx, 8, ospreyParserRULE_fnDecl)
	var _la int

	p.SetState(192)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(154)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserDOC_COMMENT {
			{
				p.SetState(153)
				p.DocComment()
			}

		}
		{
			p.SetState(156)
			p.Match(ospreyParserFN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(157)
			p.PluginName()
		}
		{
			p.SetState(158)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(159)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(161)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(160)
				p.ParamList()
			}

		}
		{
			p.SetState(163)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(165)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserAS || _la == ospreyParserARROW {
			{
				p.SetState(164)
				p.PluginReturnType()
			}

		}
		{
			p.SetState(167)
			p.Match(ospreyParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(168)
			p.PluginContent()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(171)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserDOC_COMMENT {
			{
				p.SetState(170)
				p.DocComment()
			}

		}
		{
			p.SetState(173)
			p.Match(ospreyParserFN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(174)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(175)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(177)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(176)
				p.ParamList()
			}

		}
		{
			p.SetState(179)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(182)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserARROW {
			{
				p.SetState(180)
				p.Match(ospreyParserARROW)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(181)
				p.Type_()
			}

		}
		p.SetState(190)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case ospreyParserEQ:
			{
				p.SetState(184)
				p.Match(ospreyParserEQ)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(185)
				p.Expr()
			}

		case ospreyParserLBRACE:
			{
				p.SetState(186)
				p.Match(ospreyParserLBRACE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(187)
				p.BlockBody()
			}
			{
				p.SetState(188)
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

// IPluginContentContext is an interface to support dynamic dispatch.
type IPluginContentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PLUGIN_CONTENT() antlr.TerminalNode

	// IsPluginContentContext differentiates from other interfaces.
	IsPluginContentContext()
}

type PluginContentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPluginContentContext() *PluginContentContext {
	var p = new(PluginContentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pluginContent
	return p
}

func InitEmptyPluginContentContext(p *PluginContentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pluginContent
}

func (*PluginContentContext) IsPluginContentContext() {}

func NewPluginContentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PluginContentContext {
	var p = new(PluginContentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_pluginContent

	return p
}

func (s *PluginContentContext) GetParser() antlr.Parser { return s.parser }

func (s *PluginContentContext) PLUGIN_CONTENT() antlr.TerminalNode {
	return s.GetToken(ospreyParserPLUGIN_CONTENT, 0)
}

func (s *PluginContentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PluginContentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PluginContentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterPluginContent(s)
	}
}

func (s *PluginContentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitPluginContent(s)
	}
}

func (p *ospreyParser) PluginContent() (localctx IPluginContentContext) {
	localctx = NewPluginContentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, ospreyParserRULE_pluginContent)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(194)
		p.Match(ospreyParserPLUGIN_CONTENT)
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

// IPluginNameContext is an interface to support dynamic dispatch.
type IPluginNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsPluginNameContext differentiates from other interfaces.
	IsPluginNameContext()
}

type PluginNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPluginNameContext() *PluginNameContext {
	var p = new(PluginNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pluginName
	return p
}

func InitEmptyPluginNameContext(p *PluginNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pluginName
}

func (*PluginNameContext) IsPluginNameContext() {}

func NewPluginNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PluginNameContext {
	var p = new(PluginNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_pluginName

	return p
}

func (s *PluginNameContext) GetParser() antlr.Parser { return s.parser }

func (s *PluginNameContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *PluginNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PluginNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PluginNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterPluginName(s)
	}
}

func (s *PluginNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitPluginName(s)
	}
}

func (p *ospreyParser) PluginName() (localctx IPluginNameContext) {
	localctx = NewPluginNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, ospreyParserRULE_pluginName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(196)
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

// IPluginReturnTypeContext is an interface to support dynamic dispatch.
type IPluginReturnTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ARROW() antlr.TerminalNode
	Type_() ITypeContext
	AS() antlr.TerminalNode

	// IsPluginReturnTypeContext differentiates from other interfaces.
	IsPluginReturnTypeContext()
}

type PluginReturnTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPluginReturnTypeContext() *PluginReturnTypeContext {
	var p = new(PluginReturnTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pluginReturnType
	return p
}

func InitEmptyPluginReturnTypeContext(p *PluginReturnTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pluginReturnType
}

func (*PluginReturnTypeContext) IsPluginReturnTypeContext() {}

func NewPluginReturnTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PluginReturnTypeContext {
	var p = new(PluginReturnTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_pluginReturnType

	return p
}

func (s *PluginReturnTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *PluginReturnTypeContext) ARROW() antlr.TerminalNode {
	return s.GetToken(ospreyParserARROW, 0)
}

func (s *PluginReturnTypeContext) Type_() ITypeContext {
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

func (s *PluginReturnTypeContext) AS() antlr.TerminalNode {
	return s.GetToken(ospreyParserAS, 0)
}

func (s *PluginReturnTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PluginReturnTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PluginReturnTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterPluginReturnType(s)
	}
}

func (s *PluginReturnTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitPluginReturnType(s)
	}
}

func (p *ospreyParser) PluginReturnType() (localctx IPluginReturnTypeContext) {
	localctx = NewPluginReturnTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, ospreyParserRULE_pluginReturnType)
	p.SetState(202)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserARROW:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(198)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(199)
			p.Type_()
		}

	case ospreyParserAS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(200)
			p.Match(ospreyParserAS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(201)
			p.Type_()
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
	p.EnterRule(localctx, 16, ospreyParserRULE_externDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(205)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(204)
			p.DocComment()
		}

	}
	{
		p.SetState(207)
		p.Match(ospreyParserEXTERN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(208)
		p.Match(ospreyParserFN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(209)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(210)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserID {
		{
			p.SetState(211)
			p.ExternParamList()
		}

	}
	{
		p.SetState(214)
		p.Match(ospreyParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(217)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserARROW {
		{
			p.SetState(215)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(216)
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
	p.EnterRule(localctx, 18, ospreyParserRULE_externParamList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(219)
		p.ExternParam()
	}
	p.SetState(224)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(220)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(221)
			p.ExternParam()
		}

		p.SetState(226)
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
	p.EnterRule(localctx, 20, ospreyParserRULE_externParam)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(227)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(228)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(229)
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
	p.EnterRule(localctx, 22, ospreyParserRULE_paramList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(231)
		p.Param()
	}
	p.SetState(236)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(232)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(233)
			p.Param()
		}

		p.SetState(238)
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
	p.EnterRule(localctx, 24, ospreyParserRULE_param)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(239)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(242)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserCOLON {
		{
			p.SetState(240)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(241)
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
	p.EnterRule(localctx, 26, ospreyParserRULE_typeDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(245)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(244)
			p.DocComment()
		}

	}
	{
		p.SetState(247)
		p.Match(ospreyParserTYPE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(248)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserLT {
		{
			p.SetState(249)
			p.Match(ospreyParserLT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(250)
			p.TypeParamList()
		}
		{
			p.SetState(251)
			p.Match(ospreyParserGT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(255)
		p.Match(ospreyParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(258)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserID:
		{
			p.SetState(256)
			p.UnionType()
		}

	case ospreyParserLBRACE:
		{
			p.SetState(257)
			p.RecordType()
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
	p.EnterRule(localctx, 28, ospreyParserRULE_typeParamList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(260)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(265)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(261)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(262)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(267)
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
	p.EnterRule(localctx, 30, ospreyParserRULE_unionType)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(268)
		p.Variant()
	}
	p.SetState(273)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(269)
				p.Match(ospreyParserBAR)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(270)
				p.Variant()
			}

		}
		p.SetState(275)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 32, ospreyParserRULE_recordType)
	p.EnterOuterAlt(localctx, 1)
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
	p.EnterRule(localctx, 34, ospreyParserRULE_variant)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(280)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(285)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(281)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(282)
			p.FieldDeclarations()
		}
		{
			p.SetState(283)
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
	p.EnterRule(localctx, 36, ospreyParserRULE_fieldDeclarations)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(287)
		p.FieldDeclaration()
	}
	p.SetState(292)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(288)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(289)
			p.FieldDeclaration()
		}

		p.SetState(294)
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
	Constraint() IConstraintContext

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

func (s *FieldDeclarationContext) Constraint() IConstraintContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConstraintContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConstraintContext)
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
	p.EnterRule(localctx, 38, ospreyParserRULE_fieldDeclaration)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(295)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(296)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(297)
		p.Type_()
	}
	p.SetState(299)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserWHERE {
		{
			p.SetState(298)
			p.Constraint()
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

// IConstraintContext is an interface to support dynamic dispatch.
type IConstraintContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHERE() antlr.TerminalNode
	FunctionCall() IFunctionCallContext

	// IsConstraintContext differentiates from other interfaces.
	IsConstraintContext()
}

type ConstraintContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConstraintContext() *ConstraintContext {
	var p = new(ConstraintContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_constraint
	return p
}

func InitEmptyConstraintContext(p *ConstraintContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_constraint
}

func (*ConstraintContext) IsConstraintContext() {}

func NewConstraintContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstraintContext {
	var p = new(ConstraintContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_constraint

	return p
}

func (s *ConstraintContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstraintContext) WHERE() antlr.TerminalNode {
	return s.GetToken(ospreyParserWHERE, 0)
}

func (s *ConstraintContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *ConstraintContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstraintContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConstraintContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterConstraint(s)
	}
}

func (s *ConstraintContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitConstraint(s)
	}
}

func (p *ospreyParser) Constraint() (localctx IConstraintContext) {
	localctx = NewConstraintContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, ospreyParserRULE_constraint)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(301)
		p.Match(ospreyParserWHERE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(302)
		p.FunctionCall()
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
	p.EnterRule(localctx, 42, ospreyParserRULE_functionCall)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(304)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(305)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(307)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
		{
			p.SetState(306)
			p.ArgList()
		}

	}
	{
		p.SetState(309)
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
	p.EnterRule(localctx, 44, ospreyParserRULE_booleanExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(311)
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
	p.EnterRule(localctx, 46, ospreyParserRULE_fieldList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(313)
		p.Field()
	}
	p.SetState(318)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(314)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(315)
			p.Field()
		}

		p.SetState(320)
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
	p.EnterRule(localctx, 48, ospreyParserRULE_field)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(321)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(322)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(323)
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
	ID() antlr.TerminalNode
	LSQUARE() antlr.TerminalNode
	RSQUARE() antlr.TerminalNode
	LT() antlr.TerminalNode
	TypeList() ITypeListContext
	GT() antlr.TerminalNode
	Type_() ITypeContext

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

func (s *TypeContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
}

func (s *TypeContext) LSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLSQUARE, 0)
}

func (s *TypeContext) RSQUARE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRSQUARE, 0)
}

func (s *TypeContext) LT() antlr.TerminalNode {
	return s.GetToken(ospreyParserLT, 0)
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

func (s *TypeContext) GT() antlr.TerminalNode {
	return s.GetToken(ospreyParserGT, 0)
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
	p.EnterRule(localctx, 50, ospreyParserRULE_type)
	var _la int

	p.SetState(346)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 31, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(325)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(330)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLT {
			{
				p.SetState(326)
				p.Match(ospreyParserLT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(327)
				p.TypeList()
			}
			{
				p.SetState(328)
				p.Match(ospreyParserGT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(332)
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(333)
			p.Match(ospreyParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(334)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(339)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLT {
			{
				p.SetState(335)
				p.Match(ospreyParserLT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(336)
				p.TypeList()
			}
			{
				p.SetState(337)
				p.Match(ospreyParserGT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(341)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(342)
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(343)
			p.Type_()
		}
		{
			p.SetState(344)
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
	p.EnterRule(localctx, 52, ospreyParserRULE_typeList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(348)
		p.Type_()
	}
	p.SetState(353)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(349)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(350)
			p.Type_()
		}

		p.SetState(355)
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
	p.EnterRule(localctx, 54, ospreyParserRULE_exprStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(356)
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
	p.EnterRule(localctx, 56, ospreyParserRULE_expr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(358)
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
	p.EnterRule(localctx, 58, ospreyParserRULE_matchExpr)
	var _la int

	p.SetState(372)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(360)
			p.Match(ospreyParserMATCH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(361)
			p.Expr()
		}
		{
			p.SetState(362)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(364)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420968061394952) != 0) {
			{
				p.SetState(363)
				p.MatchArm()
			}

			p.SetState(366)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(368)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(370)
			p.SelectExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(371)
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
	p.EnterRule(localctx, 60, ospreyParserRULE_selectExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(374)
		p.Match(ospreyParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(375)
		p.Match(ospreyParserLBRACE)
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

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420968061394952) != 0) {
		{
			p.SetState(376)
			p.SelectArm()
		}

		p.SetState(379)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(381)
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
	p.EnterRule(localctx, 62, ospreyParserRULE_selectArm)
	p.SetState(390)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(383)
			p.Pattern()
		}
		{
			p.SetState(384)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(385)
			p.Expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(387)
			p.Match(ospreyParserUNDERSCORE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(388)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(389)
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
	p.EnterRule(localctx, 64, ospreyParserRULE_binaryExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(392)
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
	p.EnterRule(localctx, 66, ospreyParserRULE_comparisonExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(394)
		p.AddExpr()
	}
	p.SetState(399)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(395)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3306587947008) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(396)
				p.AddExpr()
			}

		}
		p.SetState(401)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 68, ospreyParserRULE_addExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(402)
		p.MulExpr()
	}
	p.SetState(407)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(403)
				_la = p.GetTokenStream().LA(1)

				if !(_la == ospreyParserPLUS || _la == ospreyParserMINUS) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(404)
				p.MulExpr()
			}

		}
		p.SetState(409)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 70, ospreyParserRULE_mulExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(410)
		p.UnaryExpr()
	}
	p.SetState(415)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 39, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(411)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3377716900397056) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(412)
				p.UnaryExpr()
			}

		}
		p.SetState(417)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 39, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 72, ospreyParserRULE_unaryExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(419)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(418)
			_la = p.GetTokenStream().LA(1)

			if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&844433520099328) != 0) {
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
		p.SetState(421)
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
	p.EnterRule(localctx, 74, ospreyParserRULE_pipeExpr)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(423)
		p.CallExpr()
	}
	p.SetState(428)
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
				p.SetState(424)
				p.Match(ospreyParserPIPE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(425)
				p.CallExpr()
			}

		}
		p.SetState(430)
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
	p.EnterRule(localctx, 76, ospreyParserRULE_callExpr)
	var _la int

	var _alt int

	p.SetState(465)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(431)
			p.Primary()
		}
		p.SetState(434)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(432)
					p.Match(ospreyParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(433)
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

			p.SetState(436)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		p.SetState(443)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(438)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(440)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
				{
					p.SetState(439)
					p.ArgList()
				}

			}
			{
				p.SetState(442)
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
			p.SetState(445)
			p.Primary()
		}
		p.SetState(453)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(446)
					p.Match(ospreyParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(447)
					p.Match(ospreyParserID)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				{
					p.SetState(448)
					p.Match(ospreyParserLPAREN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(450)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
					{
						p.SetState(449)
						p.ArgList()
					}

				}
				{
					p.SetState(452)
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

			p.SetState(455)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(457)
			p.Primary()
		}
		p.SetState(463)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(458)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(460)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
				{
					p.SetState(459)
					p.ArgList()
				}

			}
			{
				p.SetState(462)
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
	p.EnterRule(localctx, 78, ospreyParserRULE_argList)
	var _la int

	p.SetState(476)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(467)
			p.NamedArgList()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(468)
			p.Expr()
		}
		p.SetState(473)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == ospreyParserCOMMA {
			{
				p.SetState(469)
				p.Match(ospreyParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(470)
				p.Expr()
			}

			p.SetState(475)
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
	p.EnterRule(localctx, 80, ospreyParserRULE_namedArgList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(478)
		p.NamedArg()
	}
	p.SetState(481)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserCOMMA {
		{
			p.SetState(479)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(480)
			p.NamedArg()
		}

		p.SetState(483)
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
	p.EnterRule(localctx, 82, ospreyParserRULE_namedArg)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(485)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(486)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(487)
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
	TypeConstructor() ITypeConstructorContext
	UpdateExpr() IUpdateExprContext
	BlockExpr() IBlockExprContext
	Literal() ILiteralContext
	LambdaExpr() ILambdaExprContext
	ID() antlr.TerminalNode

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

func (s *PrimaryContext) ID() antlr.TerminalNode {
	return s.GetToken(ospreyParserID, 0)
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
	p.EnterRule(localctx, 84, ospreyParserRULE_primary)
	p.SetState(524)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 54, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(489)
			p.Match(ospreyParserSPAWN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(490)
			p.Expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(491)
			p.Match(ospreyParserYIELD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(493)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 53, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(492)
				p.Expr()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(495)
			p.Match(ospreyParserAWAIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(496)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(497)
			p.Expr()
		}
		{
			p.SetState(498)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(500)
			p.Match(ospreyParserSEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(501)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(502)
			p.Expr()
		}
		{
			p.SetState(503)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(504)
			p.Expr()
		}
		{
			p.SetState(505)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(507)
			p.Match(ospreyParserRECV)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(508)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(509)
			p.Expr()
		}
		{
			p.SetState(510)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(512)
			p.Match(ospreyParserSELECT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(513)
			p.SelectExpr()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(514)
			p.TypeConstructor()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(515)
			p.UpdateExpr()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(516)
			p.BlockExpr()
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(517)
			p.Literal()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(518)
			p.LambdaExpr()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(519)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(520)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(521)
			p.Expr()
		}
		{
			p.SetState(522)
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
	p.EnterRule(localctx, 86, ospreyParserRULE_typeConstructor)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(526)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(528)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserLT {
		{
			p.SetState(527)
			p.TypeArgs()
		}

	}
	{
		p.SetState(530)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(531)
		p.FieldAssignments()
	}
	{
		p.SetState(532)
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
	p.EnterRule(localctx, 88, ospreyParserRULE_typeArgs)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(534)
		p.Match(ospreyParserLT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(535)
		p.TypeList()
	}
	{
		p.SetState(536)
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
	p.EnterRule(localctx, 90, ospreyParserRULE_fieldAssignments)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(538)
		p.FieldAssignment()
	}
	p.SetState(543)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(539)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(540)
			p.FieldAssignment()
		}

		p.SetState(545)
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
	p.EnterRule(localctx, 92, ospreyParserRULE_fieldAssignment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(546)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(547)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(548)
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
	p.EnterRule(localctx, 94, ospreyParserRULE_lambdaExpr)
	var _la int

	p.SetState(569)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserFN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(550)
			p.Match(ospreyParserFN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(551)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(553)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(552)
				p.ParamList()
			}

		}
		{
			p.SetState(555)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(558)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserARROW {
			{
				p.SetState(556)
				p.Match(ospreyParserARROW)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(557)
				p.Type_()
			}

		}
		{
			p.SetState(560)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(561)
			p.Expr()
		}

	case ospreyParserBAR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(562)
			p.Match(ospreyParserBAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(564)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(563)
				p.ParamList()
			}

		}
		{
			p.SetState(566)
			p.Match(ospreyParserBAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(567)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(568)
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
	p.EnterRule(localctx, 96, ospreyParserRULE_updateExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(571)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(572)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(573)
		p.FieldAssignments()
	}
	{
		p.SetState(574)
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
	p.EnterRule(localctx, 98, ospreyParserRULE_blockExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(576)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(577)
		p.BlockBody()
	}
	{
		p.SetState(578)
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
	p.EnterRule(localctx, 100, ospreyParserRULE_literal)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(580)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&31525197397884928) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
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
	p.EnterRule(localctx, 102, ospreyParserRULE_docComment)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(583)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(582)
			p.Match(ospreyParserDOC_COMMENT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(585)
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
	p.EnterRule(localctx, 104, ospreyParserRULE_moduleDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(588)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(587)
			p.DocComment()
		}

	}
	{
		p.SetState(590)
		p.Match(ospreyParserMODULE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(591)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(592)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(593)
		p.ModuleBody()
	}
	{
		p.SetState(594)
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
	p.EnterRule(localctx, 106, ospreyParserRULE_moduleBody)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(599)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&144115188075856712) != 0 {
		{
			p.SetState(596)
			p.ModuleStatement()
		}

		p.SetState(601)
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
	p.EnterRule(localctx, 108, ospreyParserRULE_moduleStatement)
	p.SetState(605)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 64, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(602)
			p.LetDecl()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(603)
			p.FnDecl()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(604)
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
	p.EnterRule(localctx, 110, ospreyParserRULE_matchArm)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(607)
		p.Pattern()
	}
	{
		p.SetState(608)
		p.Match(ospreyParserLAMBDA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(609)
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
	p.EnterRule(localctx, 112, ospreyParserRULE_pattern)
	var _la int

	p.SetState(651)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 69, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(611)
			p.UnaryExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(612)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(617)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLBRACE {
			{
				p.SetState(613)
				p.Match(ospreyParserLBRACE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(614)
				p.FieldPattern()
			}
			{
				p.SetState(615)
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
			p.SetState(619)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(631)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLPAREN {
			{
				p.SetState(620)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(621)
				p.Pattern()
			}
			p.SetState(626)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == ospreyParserCOMMA {
				{
					p.SetState(622)
					p.Match(ospreyParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(623)
					p.Pattern()
				}

				p.SetState(628)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(629)
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
			p.SetState(633)
			p.Match(ospreyParserID)
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
			p.SetState(637)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(638)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(639)
			p.Type_()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(640)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(641)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(642)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(643)
			p.FieldPattern()
		}
		{
			p.SetState(644)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(646)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(647)
			p.FieldPattern()
		}
		{
			p.SetState(648)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(650)
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
	p.EnterRule(localctx, 114, ospreyParserRULE_fieldPattern)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(653)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(658)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(654)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(655)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(660)
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
	p.EnterRule(localctx, 116, ospreyParserRULE_blockBody)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(664)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 71, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(661)
				p.Statement()
			}

		}
		p.SetState(666)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 71, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(668)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
		{
			p.SetState(667)
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
