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
		"ANY_CHAR",
	}
	staticData.RuleNames = []string{
		"program", "statement", "importStmt", "letDecl", "fnDecl", "pluginContent",
		"pluginToken", "pluginReturnType", "externDecl", "externParamList",
		"externParam", "paramList", "param", "typeDecl", "typeParamList", "unionType",
		"recordType", "variant", "fieldDeclarations", "fieldDeclaration", "constraint",
		"functionCall", "booleanExpr", "fieldList", "field", "type", "typeList",
		"exprStmt", "expr", "matchExpr", "selectExpr", "selectArm", "binaryExpr",
		"comparisonExpr", "addExpr", "mulExpr", "unaryExpr", "pipeExpr", "callExpr",
		"argList", "namedArgList", "namedArg", "primary", "typeConstructor",
		"typeArgs", "fieldAssignments", "fieldAssignment", "lambdaExpr", "updateExpr",
		"blockExpr", "literal", "docComment", "moduleDecl", "moduleBody", "moduleStatement",
		"matchArm", "pattern", "fieldPattern", "blockBody",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 59, 677, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
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
		1, 4, 1, 4, 1, 4, 3, 4, 161, 8, 4, 1, 4, 1, 4, 1, 4, 3, 4, 166, 8, 4, 1,
		4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 174, 8, 4, 1, 4, 3, 4, 177, 8, 4,
		1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 184, 8, 4, 1, 4, 1, 4, 3, 4, 188, 8,
		4, 1, 4, 1, 4, 3, 4, 192, 8, 4, 1, 5, 4, 5, 195, 8, 5, 11, 5, 12, 5, 196,
		1, 5, 1, 5, 1, 6, 1, 6, 3, 6, 203, 8, 6, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7,
		209, 8, 7, 1, 8, 3, 8, 212, 8, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 219,
		8, 8, 1, 8, 1, 8, 1, 8, 3, 8, 224, 8, 8, 1, 9, 1, 9, 1, 9, 5, 9, 229, 8,
		9, 10, 9, 12, 9, 232, 9, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1,
		11, 5, 11, 241, 8, 11, 10, 11, 12, 11, 244, 9, 11, 1, 12, 1, 12, 1, 12,
		3, 12, 249, 8, 12, 1, 13, 3, 13, 252, 8, 13, 1, 13, 1, 13, 1, 13, 1, 13,
		1, 13, 1, 13, 3, 13, 260, 8, 13, 1, 13, 1, 13, 1, 13, 3, 13, 265, 8, 13,
		1, 14, 1, 14, 1, 14, 5, 14, 270, 8, 14, 10, 14, 12, 14, 273, 9, 14, 1,
		15, 1, 15, 1, 15, 5, 15, 278, 8, 15, 10, 15, 12, 15, 281, 9, 15, 1, 16,
		1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 3, 17, 292, 8,
		17, 1, 18, 1, 18, 1, 18, 5, 18, 297, 8, 18, 10, 18, 12, 18, 300, 9, 18,
		1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 306, 8, 19, 1, 20, 1, 20, 1, 20, 1,
		21, 1, 21, 1, 21, 3, 21, 314, 8, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 23,
		1, 23, 1, 23, 5, 23, 323, 8, 23, 10, 23, 12, 23, 326, 9, 23, 1, 24, 1,
		24, 1, 24, 1, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 337, 8, 25,
		1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 346, 8, 25, 1,
		25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 353, 8, 25, 1, 26, 1, 26, 1, 26,
		5, 26, 358, 8, 26, 10, 26, 12, 26, 361, 9, 26, 1, 27, 1, 27, 1, 28, 1,
		28, 1, 29, 1, 29, 1, 29, 1, 29, 4, 29, 371, 8, 29, 11, 29, 12, 29, 372,
		1, 29, 1, 29, 1, 29, 1, 29, 3, 29, 379, 8, 29, 1, 30, 1, 30, 1, 30, 4,
		30, 384, 8, 30, 11, 30, 12, 30, 385, 1, 30, 1, 30, 1, 31, 1, 31, 1, 31,
		1, 31, 1, 31, 1, 31, 1, 31, 3, 31, 397, 8, 31, 1, 32, 1, 32, 1, 33, 1,
		33, 1, 33, 5, 33, 404, 8, 33, 10, 33, 12, 33, 407, 9, 33, 1, 34, 1, 34,
		1, 34, 5, 34, 412, 8, 34, 10, 34, 12, 34, 415, 9, 34, 1, 35, 1, 35, 1,
		35, 5, 35, 420, 8, 35, 10, 35, 12, 35, 423, 9, 35, 1, 36, 3, 36, 426, 8,
		36, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 5, 37, 433, 8, 37, 10, 37, 12, 37,
		436, 9, 37, 1, 38, 1, 38, 1, 38, 4, 38, 441, 8, 38, 11, 38, 12, 38, 442,
		1, 38, 1, 38, 3, 38, 447, 8, 38, 1, 38, 3, 38, 450, 8, 38, 1, 38, 1, 38,
		1, 38, 1, 38, 1, 38, 3, 38, 457, 8, 38, 1, 38, 4, 38, 460, 8, 38, 11, 38,
		12, 38, 461, 1, 38, 1, 38, 1, 38, 3, 38, 467, 8, 38, 1, 38, 3, 38, 470,
		8, 38, 3, 38, 472, 8, 38, 1, 39, 1, 39, 1, 39, 1, 39, 5, 39, 478, 8, 39,
		10, 39, 12, 39, 481, 9, 39, 3, 39, 483, 8, 39, 1, 40, 1, 40, 1, 40, 4,
		40, 488, 8, 40, 11, 40, 12, 40, 489, 1, 41, 1, 41, 1, 41, 1, 41, 1, 42,
		1, 42, 1, 42, 1, 42, 3, 42, 500, 8, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1,
		42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42,
		1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1,
		42, 1, 42, 1, 42, 1, 42, 3, 42, 531, 8, 42, 1, 43, 1, 43, 3, 43, 535, 8,
		43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1, 44, 1, 44, 1, 45, 1, 45,
		1, 45, 5, 45, 548, 8, 45, 10, 45, 12, 45, 551, 9, 45, 1, 46, 1, 46, 1,
		46, 1, 46, 1, 47, 1, 47, 1, 47, 3, 47, 560, 8, 47, 1, 47, 1, 47, 1, 47,
		3, 47, 565, 8, 47, 1, 47, 1, 47, 1, 47, 1, 47, 3, 47, 571, 8, 47, 1, 47,
		1, 47, 1, 47, 3, 47, 576, 8, 47, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1,
		49, 1, 49, 1, 49, 1, 49, 1, 50, 1, 50, 1, 51, 4, 51, 590, 8, 51, 11, 51,
		12, 51, 591, 1, 52, 3, 52, 595, 8, 52, 1, 52, 1, 52, 1, 52, 1, 52, 1, 52,
		1, 52, 1, 53, 5, 53, 604, 8, 53, 10, 53, 12, 53, 607, 9, 53, 1, 54, 1,
		54, 1, 54, 3, 54, 612, 8, 54, 1, 55, 1, 55, 1, 55, 1, 55, 1, 56, 1, 56,
		1, 56, 1, 56, 1, 56, 1, 56, 3, 56, 624, 8, 56, 1, 56, 1, 56, 1, 56, 1,
		56, 1, 56, 5, 56, 631, 8, 56, 10, 56, 12, 56, 634, 9, 56, 1, 56, 1, 56,
		3, 56, 638, 8, 56, 1, 56, 1, 56, 3, 56, 642, 8, 56, 1, 56, 1, 56, 1, 56,
		1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1,
		56, 3, 56, 658, 8, 56, 1, 57, 1, 57, 1, 57, 5, 57, 663, 8, 57, 10, 57,
		12, 57, 666, 9, 57, 1, 58, 5, 58, 669, 8, 58, 10, 58, 12, 58, 672, 9, 58,
		1, 58, 3, 58, 675, 8, 58, 1, 58, 0, 0, 59, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52,
		54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88,
		90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 112, 114, 116, 0, 7,
		1, 0, 8, 9, 1, 0, 36, 36, 2, 0, 29, 32, 40, 41, 1, 0, 48, 49, 2, 0, 34,
		34, 50, 51, 3, 0, 15, 15, 33, 33, 48, 49, 2, 0, 21, 22, 52, 54, 718, 0,
		121, 1, 0, 0, 0, 2, 133, 1, 0, 0, 0, 4, 135, 1, 0, 0, 0, 6, 144, 1, 0,
		0, 0, 8, 191, 1, 0, 0, 0, 10, 194, 1, 0, 0, 0, 12, 202, 1, 0, 0, 0, 14,
		208, 1, 0, 0, 0, 16, 211, 1, 0, 0, 0, 18, 225, 1, 0, 0, 0, 20, 233, 1,
		0, 0, 0, 22, 237, 1, 0, 0, 0, 24, 245, 1, 0, 0, 0, 26, 251, 1, 0, 0, 0,
		28, 266, 1, 0, 0, 0, 30, 274, 1, 0, 0, 0, 32, 282, 1, 0, 0, 0, 34, 286,
		1, 0, 0, 0, 36, 293, 1, 0, 0, 0, 38, 301, 1, 0, 0, 0, 40, 307, 1, 0, 0,
		0, 42, 310, 1, 0, 0, 0, 44, 317, 1, 0, 0, 0, 46, 319, 1, 0, 0, 0, 48, 327,
		1, 0, 0, 0, 50, 352, 1, 0, 0, 0, 52, 354, 1, 0, 0, 0, 54, 362, 1, 0, 0,
		0, 56, 364, 1, 0, 0, 0, 58, 378, 1, 0, 0, 0, 60, 380, 1, 0, 0, 0, 62, 396,
		1, 0, 0, 0, 64, 398, 1, 0, 0, 0, 66, 400, 1, 0, 0, 0, 68, 408, 1, 0, 0,
		0, 70, 416, 1, 0, 0, 0, 72, 425, 1, 0, 0, 0, 74, 429, 1, 0, 0, 0, 76, 471,
		1, 0, 0, 0, 78, 482, 1, 0, 0, 0, 80, 484, 1, 0, 0, 0, 82, 491, 1, 0, 0,
		0, 84, 530, 1, 0, 0, 0, 86, 532, 1, 0, 0, 0, 88, 540, 1, 0, 0, 0, 90, 544,
		1, 0, 0, 0, 92, 552, 1, 0, 0, 0, 94, 575, 1, 0, 0, 0, 96, 577, 1, 0, 0,
		0, 98, 582, 1, 0, 0, 0, 100, 586, 1, 0, 0, 0, 102, 589, 1, 0, 0, 0, 104,
		594, 1, 0, 0, 0, 106, 605, 1, 0, 0, 0, 108, 611, 1, 0, 0, 0, 110, 613,
		1, 0, 0, 0, 112, 657, 1, 0, 0, 0, 114, 659, 1, 0, 0, 0, 116, 670, 1, 0,
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
		5, 3, 0, 0, 157, 158, 5, 55, 0, 0, 158, 160, 5, 42, 0, 0, 159, 161, 3,
		22, 11, 0, 160, 159, 1, 0, 0, 0, 160, 161, 1, 0, 0, 0, 161, 162, 1, 0,
		0, 0, 162, 165, 5, 43, 0, 0, 163, 164, 5, 25, 0, 0, 164, 166, 3, 50, 25,
		0, 165, 163, 1, 0, 0, 0, 165, 166, 1, 0, 0, 0, 166, 173, 1, 0, 0, 0, 167,
		168, 5, 28, 0, 0, 168, 174, 3, 56, 28, 0, 169, 170, 5, 44, 0, 0, 170, 171,
		3, 116, 58, 0, 171, 172, 5, 45, 0, 0, 172, 174, 1, 0, 0, 0, 173, 167, 1,
		0, 0, 0, 173, 169, 1, 0, 0, 0, 174, 192, 1, 0, 0, 0, 175, 177, 3, 102,
		51, 0, 176, 175, 1, 0, 0, 0, 176, 177, 1, 0, 0, 0, 177, 178, 1, 0, 0, 0,
		178, 179, 5, 3, 0, 0, 179, 180, 5, 55, 0, 0, 180, 181, 5, 55, 0, 0, 181,
		183, 5, 42, 0, 0, 182, 184, 3, 22, 11, 0, 183, 182, 1, 0, 0, 0, 183, 184,
		1, 0, 0, 0, 184, 185, 1, 0, 0, 0, 185, 187, 5, 43, 0, 0, 186, 188, 3, 14,
		7, 0, 187, 186, 1, 0, 0, 0, 187, 188, 1, 0, 0, 0, 188, 189, 1, 0, 0, 0,
		189, 190, 5, 28, 0, 0, 190, 192, 3, 10, 5, 0, 191, 154, 1, 0, 0, 0, 191,
		176, 1, 0, 0, 0, 192, 9, 1, 0, 0, 0, 193, 195, 3, 12, 6, 0, 194, 193, 1,
		0, 0, 0, 195, 196, 1, 0, 0, 0, 196, 194, 1, 0, 0, 0, 196, 197, 1, 0, 0,
		0, 197, 198, 1, 0, 0, 0, 198, 199, 5, 36, 0, 0, 199, 11, 1, 0, 0, 0, 200,
		203, 8, 1, 0, 0, 201, 203, 5, 59, 0, 0, 202, 200, 1, 0, 0, 0, 202, 201,
		1, 0, 0, 0, 203, 13, 1, 0, 0, 0, 204, 205, 5, 25, 0, 0, 205, 209, 3, 50,
		25, 0, 206, 207, 5, 24, 0, 0, 207, 209, 3, 50, 25, 0, 208, 204, 1, 0, 0,
		0, 208, 206, 1, 0, 0, 0, 209, 15, 1, 0, 0, 0, 210, 212, 3, 102, 51, 0,
		211, 210, 1, 0, 0, 0, 211, 212, 1, 0, 0, 0, 212, 213, 1, 0, 0, 0, 213,
		214, 5, 4, 0, 0, 214, 215, 5, 3, 0, 0, 215, 216, 5, 55, 0, 0, 216, 218,
		5, 42, 0, 0, 217, 219, 3, 18, 9, 0, 218, 217, 1, 0, 0, 0, 218, 219, 1,
		0, 0, 0, 219, 220, 1, 0, 0, 0, 220, 223, 5, 43, 0, 0, 221, 222, 5, 25,
		0, 0, 222, 224, 3, 50, 25, 0, 223, 221, 1, 0, 0, 0, 223, 224, 1, 0, 0,
		0, 224, 17, 1, 0, 0, 0, 225, 230, 3, 20, 10, 0, 226, 227, 5, 37, 0, 0,
		227, 229, 3, 20, 10, 0, 228, 226, 1, 0, 0, 0, 229, 232, 1, 0, 0, 0, 230,
		228, 1, 0, 0, 0, 230, 231, 1, 0, 0, 0, 231, 19, 1, 0, 0, 0, 232, 230, 1,
		0, 0, 0, 233, 234, 5, 55, 0, 0, 234, 235, 5, 35, 0, 0, 235, 236, 3, 50,
		25, 0, 236, 21, 1, 0, 0, 0, 237, 242, 3, 24, 12, 0, 238, 239, 5, 37, 0,
		0, 239, 241, 3, 24, 12, 0, 240, 238, 1, 0, 0, 0, 241, 244, 1, 0, 0, 0,
		242, 240, 1, 0, 0, 0, 242, 243, 1, 0, 0, 0, 243, 23, 1, 0, 0, 0, 244, 242,
		1, 0, 0, 0, 245, 248, 5, 55, 0, 0, 246, 247, 5, 35, 0, 0, 247, 249, 3,
		50, 25, 0, 248, 246, 1, 0, 0, 0, 248, 249, 1, 0, 0, 0, 249, 25, 1, 0, 0,
		0, 250, 252, 3, 102, 51, 0, 251, 250, 1, 0, 0, 0, 251, 252, 1, 0, 0, 0,
		252, 253, 1, 0, 0, 0, 253, 254, 5, 6, 0, 0, 254, 259, 5, 55, 0, 0, 255,
		256, 5, 40, 0, 0, 256, 257, 3, 28, 14, 0, 257, 258, 5, 41, 0, 0, 258, 260,
		1, 0, 0, 0, 259, 255, 1, 0, 0, 0, 259, 260, 1, 0, 0, 0, 260, 261, 1, 0,
		0, 0, 261, 264, 5, 28, 0, 0, 262, 265, 3, 30, 15, 0, 263, 265, 3, 32, 16,
		0, 264, 262, 1, 0, 0, 0, 264, 263, 1, 0, 0, 0, 265, 27, 1, 0, 0, 0, 266,
		271, 5, 55, 0, 0, 267, 268, 5, 37, 0, 0, 268, 270, 5, 55, 0, 0, 269, 267,
		1, 0, 0, 0, 270, 273, 1, 0, 0, 0, 271, 269, 1, 0, 0, 0, 271, 272, 1, 0,
		0, 0, 272, 29, 1, 0, 0, 0, 273, 271, 1, 0, 0, 0, 274, 279, 3, 34, 17, 0,
		275, 276, 5, 39, 0, 0, 276, 278, 3, 34, 17, 0, 277, 275, 1, 0, 0, 0, 278,
		281, 1, 0, 0, 0, 279, 277, 1, 0, 0, 0, 279, 280, 1, 0, 0, 0, 280, 31, 1,
		0, 0, 0, 281, 279, 1, 0, 0, 0, 282, 283, 5, 44, 0, 0, 283, 284, 3, 36,
		18, 0, 284, 285, 5, 45, 0, 0, 285, 33, 1, 0, 0, 0, 286, 291, 5, 55, 0,
		0, 287, 288, 5, 44, 0, 0, 288, 289, 3, 36, 18, 0, 289, 290, 5, 45, 0, 0,
		290, 292, 1, 0, 0, 0, 291, 287, 1, 0, 0, 0, 291, 292, 1, 0, 0, 0, 292,
		35, 1, 0, 0, 0, 293, 298, 3, 38, 19, 0, 294, 295, 5, 37, 0, 0, 295, 297,
		3, 38, 19, 0, 296, 294, 1, 0, 0, 0, 297, 300, 1, 0, 0, 0, 298, 296, 1,
		0, 0, 0, 298, 299, 1, 0, 0, 0, 299, 37, 1, 0, 0, 0, 300, 298, 1, 0, 0,
		0, 301, 302, 5, 55, 0, 0, 302, 303, 5, 35, 0, 0, 303, 305, 3, 50, 25, 0,
		304, 306, 3, 40, 20, 0, 305, 304, 1, 0, 0, 0, 305, 306, 1, 0, 0, 0, 306,
		39, 1, 0, 0, 0, 307, 308, 5, 23, 0, 0, 308, 309, 3, 42, 21, 0, 309, 41,
		1, 0, 0, 0, 310, 311, 5, 55, 0, 0, 311, 313, 5, 42, 0, 0, 312, 314, 3,
		78, 39, 0, 313, 312, 1, 0, 0, 0, 313, 314, 1, 0, 0, 0, 314, 315, 1, 0,
		0, 0, 315, 316, 5, 43, 0, 0, 316, 43, 1, 0, 0, 0, 317, 318, 3, 66, 33,
		0, 318, 45, 1, 0, 0, 0, 319, 324, 3, 48, 24, 0, 320, 321, 5, 37, 0, 0,
		321, 323, 3, 48, 24, 0, 322, 320, 1, 0, 0, 0, 323, 326, 1, 0, 0, 0, 324,
		322, 1, 0, 0, 0, 324, 325, 1, 0, 0, 0, 325, 47, 1, 0, 0, 0, 326, 324, 1,
		0, 0, 0, 327, 328, 5, 55, 0, 0, 328, 329, 5, 35, 0, 0, 329, 330, 3, 50,
		25, 0, 330, 49, 1, 0, 0, 0, 331, 336, 5, 55, 0, 0, 332, 333, 5, 40, 0,
		0, 333, 334, 3, 52, 26, 0, 334, 335, 5, 41, 0, 0, 335, 337, 1, 0, 0, 0,
		336, 332, 1, 0, 0, 0, 336, 337, 1, 0, 0, 0, 337, 338, 1, 0, 0, 0, 338,
		339, 5, 46, 0, 0, 339, 353, 5, 47, 0, 0, 340, 345, 5, 55, 0, 0, 341, 342,
		5, 40, 0, 0, 342, 343, 3, 52, 26, 0, 343, 344, 5, 41, 0, 0, 344, 346, 1,
		0, 0, 0, 345, 341, 1, 0, 0, 0, 345, 346, 1, 0, 0, 0, 346, 353, 1, 0, 0,
		0, 347, 348, 5, 55, 0, 0, 348, 349, 5, 46, 0, 0, 349, 350, 3, 50, 25, 0,
		350, 351, 5, 47, 0, 0, 351, 353, 1, 0, 0, 0, 352, 331, 1, 0, 0, 0, 352,
		340, 1, 0, 0, 0, 352, 347, 1, 0, 0, 0, 353, 51, 1, 0, 0, 0, 354, 359, 3,
		50, 25, 0, 355, 356, 5, 37, 0, 0, 356, 358, 3, 50, 25, 0, 357, 355, 1,
		0, 0, 0, 358, 361, 1, 0, 0, 0, 359, 357, 1, 0, 0, 0, 359, 360, 1, 0, 0,
		0, 360, 53, 1, 0, 0, 0, 361, 359, 1, 0, 0, 0, 362, 363, 3, 56, 28, 0, 363,
		55, 1, 0, 0, 0, 364, 365, 3, 58, 29, 0, 365, 57, 1, 0, 0, 0, 366, 367,
		5, 2, 0, 0, 367, 368, 3, 56, 28, 0, 368, 370, 5, 44, 0, 0, 369, 371, 3,
		110, 55, 0, 370, 369, 1, 0, 0, 0, 371, 372, 1, 0, 0, 0, 372, 370, 1, 0,
		0, 0, 372, 373, 1, 0, 0, 0, 373, 374, 1, 0, 0, 0, 374, 375, 5, 45, 0, 0,
		375, 379, 1, 0, 0, 0, 376, 379, 3, 60, 30, 0, 377, 379, 3, 64, 32, 0, 378,
		366, 1, 0, 0, 0, 378, 376, 1, 0, 0, 0, 378, 377, 1, 0, 0, 0, 379, 59, 1,
		0, 0, 0, 380, 381, 5, 20, 0, 0, 381, 383, 5, 44, 0, 0, 382, 384, 3, 62,
		31, 0, 383, 382, 1, 0, 0, 0, 384, 385, 1, 0, 0, 0, 385, 383, 1, 0, 0, 0,
		385, 386, 1, 0, 0, 0, 386, 387, 1, 0, 0, 0, 387, 388, 5, 45, 0, 0, 388,
		61, 1, 0, 0, 0, 389, 390, 3, 112, 56, 0, 390, 391, 5, 26, 0, 0, 391, 392,
		3, 56, 28, 0, 392, 397, 1, 0, 0, 0, 393, 394, 5, 27, 0, 0, 394, 395, 5,
		26, 0, 0, 395, 397, 3, 56, 28, 0, 396, 389, 1, 0, 0, 0, 396, 393, 1, 0,
		0, 0, 397, 63, 1, 0, 0, 0, 398, 399, 3, 66, 33, 0, 399, 65, 1, 0, 0, 0,
		400, 405, 3, 68, 34, 0, 401, 402, 7, 2, 0, 0, 402, 404, 3, 68, 34, 0, 403,
		401, 1, 0, 0, 0, 404, 407, 1, 0, 0, 0, 405, 403, 1, 0, 0, 0, 405, 406,
		1, 0, 0, 0, 406, 67, 1, 0, 0, 0, 407, 405, 1, 0, 0, 0, 408, 413, 3, 70,
		35, 0, 409, 410, 7, 3, 0, 0, 410, 412, 3, 70, 35, 0, 411, 409, 1, 0, 0,
		0, 412, 415, 1, 0, 0, 0, 413, 411, 1, 0, 0, 0, 413, 414, 1, 0, 0, 0, 414,
		69, 1, 0, 0, 0, 415, 413, 1, 0, 0, 0, 416, 421, 3, 72, 36, 0, 417, 418,
		7, 4, 0, 0, 418, 420, 3, 72, 36, 0, 419, 417, 1, 0, 0, 0, 420, 423, 1,
		0, 0, 0, 421, 419, 1, 0, 0, 0, 421, 422, 1, 0, 0, 0, 422, 71, 1, 0, 0,
		0, 423, 421, 1, 0, 0, 0, 424, 426, 7, 5, 0, 0, 425, 424, 1, 0, 0, 0, 425,
		426, 1, 0, 0, 0, 426, 427, 1, 0, 0, 0, 427, 428, 3, 74, 37, 0, 428, 73,
		1, 0, 0, 0, 429, 434, 3, 76, 38, 0, 430, 431, 5, 1, 0, 0, 431, 433, 3,
		76, 38, 0, 432, 430, 1, 0, 0, 0, 433, 436, 1, 0, 0, 0, 434, 432, 1, 0,
		0, 0, 434, 435, 1, 0, 0, 0, 435, 75, 1, 0, 0, 0, 436, 434, 1, 0, 0, 0,
		437, 440, 3, 84, 42, 0, 438, 439, 5, 38, 0, 0, 439, 441, 5, 55, 0, 0, 440,
		438, 1, 0, 0, 0, 441, 442, 1, 0, 0, 0, 442, 440, 1, 0, 0, 0, 442, 443,
		1, 0, 0, 0, 443, 449, 1, 0, 0, 0, 444, 446, 5, 42, 0, 0, 445, 447, 3, 78,
		39, 0, 446, 445, 1, 0, 0, 0, 446, 447, 1, 0, 0, 0, 447, 448, 1, 0, 0, 0,
		448, 450, 5, 43, 0, 0, 449, 444, 1, 0, 0, 0, 449, 450, 1, 0, 0, 0, 450,
		472, 1, 0, 0, 0, 451, 459, 3, 84, 42, 0, 452, 453, 5, 38, 0, 0, 453, 454,
		5, 55, 0, 0, 454, 456, 5, 42, 0, 0, 455, 457, 3, 78, 39, 0, 456, 455, 1,
		0, 0, 0, 456, 457, 1, 0, 0, 0, 457, 458, 1, 0, 0, 0, 458, 460, 5, 43, 0,
		0, 459, 452, 1, 0, 0, 0, 460, 461, 1, 0, 0, 0, 461, 459, 1, 0, 0, 0, 461,
		462, 1, 0, 0, 0, 462, 472, 1, 0, 0, 0, 463, 469, 3, 84, 42, 0, 464, 466,
		5, 42, 0, 0, 465, 467, 3, 78, 39, 0, 466, 465, 1, 0, 0, 0, 466, 467, 1,
		0, 0, 0, 467, 468, 1, 0, 0, 0, 468, 470, 5, 43, 0, 0, 469, 464, 1, 0, 0,
		0, 469, 470, 1, 0, 0, 0, 470, 472, 1, 0, 0, 0, 471, 437, 1, 0, 0, 0, 471,
		451, 1, 0, 0, 0, 471, 463, 1, 0, 0, 0, 472, 77, 1, 0, 0, 0, 473, 483, 3,
		80, 40, 0, 474, 479, 3, 56, 28, 0, 475, 476, 5, 37, 0, 0, 476, 478, 3,
		56, 28, 0, 477, 475, 1, 0, 0, 0, 478, 481, 1, 0, 0, 0, 479, 477, 1, 0,
		0, 0, 479, 480, 1, 0, 0, 0, 480, 483, 1, 0, 0, 0, 481, 479, 1, 0, 0, 0,
		482, 473, 1, 0, 0, 0, 482, 474, 1, 0, 0, 0, 483, 79, 1, 0, 0, 0, 484, 487,
		3, 82, 41, 0, 485, 486, 5, 37, 0, 0, 486, 488, 3, 82, 41, 0, 487, 485,
		1, 0, 0, 0, 488, 489, 1, 0, 0, 0, 489, 487, 1, 0, 0, 0, 489, 490, 1, 0,
		0, 0, 490, 81, 1, 0, 0, 0, 491, 492, 5, 55, 0, 0, 492, 493, 5, 35, 0, 0,
		493, 494, 3, 56, 28, 0, 494, 83, 1, 0, 0, 0, 495, 496, 5, 13, 0, 0, 496,
		531, 3, 56, 28, 0, 497, 499, 5, 14, 0, 0, 498, 500, 3, 56, 28, 0, 499,
		498, 1, 0, 0, 0, 499, 500, 1, 0, 0, 0, 500, 531, 1, 0, 0, 0, 501, 502,
		5, 15, 0, 0, 502, 503, 5, 42, 0, 0, 503, 504, 3, 56, 28, 0, 504, 505, 5,
		43, 0, 0, 505, 531, 1, 0, 0, 0, 506, 507, 5, 18, 0, 0, 507, 508, 5, 42,
		0, 0, 508, 509, 3, 56, 28, 0, 509, 510, 5, 37, 0, 0, 510, 511, 3, 56, 28,
		0, 511, 512, 5, 43, 0, 0, 512, 531, 1, 0, 0, 0, 513, 514, 5, 19, 0, 0,
		514, 515, 5, 42, 0, 0, 515, 516, 3, 56, 28, 0, 516, 517, 5, 43, 0, 0, 517,
		531, 1, 0, 0, 0, 518, 519, 5, 20, 0, 0, 519, 531, 3, 60, 30, 0, 520, 531,
		3, 86, 43, 0, 521, 531, 3, 96, 48, 0, 522, 531, 3, 98, 49, 0, 523, 531,
		3, 100, 50, 0, 524, 531, 3, 94, 47, 0, 525, 531, 5, 55, 0, 0, 526, 527,
		5, 42, 0, 0, 527, 528, 3, 56, 28, 0, 528, 529, 5, 43, 0, 0, 529, 531, 1,
		0, 0, 0, 530, 495, 1, 0, 0, 0, 530, 497, 1, 0, 0, 0, 530, 501, 1, 0, 0,
		0, 530, 506, 1, 0, 0, 0, 530, 513, 1, 0, 0, 0, 530, 518, 1, 0, 0, 0, 530,
		520, 1, 0, 0, 0, 530, 521, 1, 0, 0, 0, 530, 522, 1, 0, 0, 0, 530, 523,
		1, 0, 0, 0, 530, 524, 1, 0, 0, 0, 530, 525, 1, 0, 0, 0, 530, 526, 1, 0,
		0, 0, 531, 85, 1, 0, 0, 0, 532, 534, 5, 55, 0, 0, 533, 535, 3, 88, 44,
		0, 534, 533, 1, 0, 0, 0, 534, 535, 1, 0, 0, 0, 535, 536, 1, 0, 0, 0, 536,
		537, 5, 44, 0, 0, 537, 538, 3, 90, 45, 0, 538, 539, 5, 45, 0, 0, 539, 87,
		1, 0, 0, 0, 540, 541, 5, 40, 0, 0, 541, 542, 3, 52, 26, 0, 542, 543, 5,
		41, 0, 0, 543, 89, 1, 0, 0, 0, 544, 549, 3, 92, 46, 0, 545, 546, 5, 37,
		0, 0, 546, 548, 3, 92, 46, 0, 547, 545, 1, 0, 0, 0, 548, 551, 1, 0, 0,
		0, 549, 547, 1, 0, 0, 0, 549, 550, 1, 0, 0, 0, 550, 91, 1, 0, 0, 0, 551,
		549, 1, 0, 0, 0, 552, 553, 5, 55, 0, 0, 553, 554, 5, 35, 0, 0, 554, 555,
		3, 56, 28, 0, 555, 93, 1, 0, 0, 0, 556, 557, 5, 3, 0, 0, 557, 559, 5, 42,
		0, 0, 558, 560, 3, 22, 11, 0, 559, 558, 1, 0, 0, 0, 559, 560, 1, 0, 0,
		0, 560, 561, 1, 0, 0, 0, 561, 564, 5, 43, 0, 0, 562, 563, 5, 25, 0, 0,
		563, 565, 3, 50, 25, 0, 564, 562, 1, 0, 0, 0, 564, 565, 1, 0, 0, 0, 565,
		566, 1, 0, 0, 0, 566, 567, 5, 26, 0, 0, 567, 576, 3, 56, 28, 0, 568, 570,
		5, 39, 0, 0, 569, 571, 3, 22, 11, 0, 570, 569, 1, 0, 0, 0, 570, 571, 1,
		0, 0, 0, 571, 572, 1, 0, 0, 0, 572, 573, 5, 39, 0, 0, 573, 574, 5, 26,
		0, 0, 574, 576, 3, 56, 28, 0, 575, 556, 1, 0, 0, 0, 575, 568, 1, 0, 0,
		0, 576, 95, 1, 0, 0, 0, 577, 578, 5, 55, 0, 0, 578, 579, 5, 44, 0, 0, 579,
		580, 3, 90, 45, 0, 580, 581, 5, 45, 0, 0, 581, 97, 1, 0, 0, 0, 582, 583,
		5, 44, 0, 0, 583, 584, 3, 116, 58, 0, 584, 585, 5, 45, 0, 0, 585, 99, 1,
		0, 0, 0, 586, 587, 7, 6, 0, 0, 587, 101, 1, 0, 0, 0, 588, 590, 5, 57, 0,
		0, 589, 588, 1, 0, 0, 0, 590, 591, 1, 0, 0, 0, 591, 589, 1, 0, 0, 0, 591,
		592, 1, 0, 0, 0, 592, 103, 1, 0, 0, 0, 593, 595, 3, 102, 51, 0, 594, 593,
		1, 0, 0, 0, 594, 595, 1, 0, 0, 0, 595, 596, 1, 0, 0, 0, 596, 597, 5, 7,
		0, 0, 597, 598, 5, 55, 0, 0, 598, 599, 5, 44, 0, 0, 599, 600, 3, 106, 53,
		0, 600, 601, 5, 45, 0, 0, 601, 105, 1, 0, 0, 0, 602, 604, 3, 108, 54, 0,
		603, 602, 1, 0, 0, 0, 604, 607, 1, 0, 0, 0, 605, 603, 1, 0, 0, 0, 605,
		606, 1, 0, 0, 0, 606, 107, 1, 0, 0, 0, 607, 605, 1, 0, 0, 0, 608, 612,
		3, 6, 3, 0, 609, 612, 3, 8, 4, 0, 610, 612, 3, 26, 13, 0, 611, 608, 1,
		0, 0, 0, 611, 609, 1, 0, 0, 0, 611, 610, 1, 0, 0, 0, 612, 109, 1, 0, 0,
		0, 613, 614, 3, 112, 56, 0, 614, 615, 5, 26, 0, 0, 615, 616, 3, 56, 28,
		0, 616, 111, 1, 0, 0, 0, 617, 658, 3, 72, 36, 0, 618, 623, 5, 55, 0, 0,
		619, 620, 5, 44, 0, 0, 620, 621, 3, 114, 57, 0, 621, 622, 5, 45, 0, 0,
		622, 624, 1, 0, 0, 0, 623, 619, 1, 0, 0, 0, 623, 624, 1, 0, 0, 0, 624,
		658, 1, 0, 0, 0, 625, 637, 5, 55, 0, 0, 626, 627, 5, 42, 0, 0, 627, 632,
		3, 112, 56, 0, 628, 629, 5, 37, 0, 0, 629, 631, 3, 112, 56, 0, 630, 628,
		1, 0, 0, 0, 631, 634, 1, 0, 0, 0, 632, 630, 1, 0, 0, 0, 632, 633, 1, 0,
		0, 0, 633, 635, 1, 0, 0, 0, 634, 632, 1, 0, 0, 0, 635, 636, 5, 43, 0, 0,
		636, 638, 1, 0, 0, 0, 637, 626, 1, 0, 0, 0, 637, 638, 1, 0, 0, 0, 638,
		658, 1, 0, 0, 0, 639, 641, 5, 55, 0, 0, 640, 642, 5, 55, 0, 0, 641, 640,
		1, 0, 0, 0, 641, 642, 1, 0, 0, 0, 642, 658, 1, 0, 0, 0, 643, 644, 5, 55,
		0, 0, 644, 645, 5, 35, 0, 0, 645, 658, 3, 50, 25, 0, 646, 647, 5, 55, 0,
		0, 647, 648, 5, 35, 0, 0, 648, 649, 5, 44, 0, 0, 649, 650, 3, 114, 57,
		0, 650, 651, 5, 45, 0, 0, 651, 658, 1, 0, 0, 0, 652, 653, 5, 44, 0, 0,
		653, 654, 3, 114, 57, 0, 654, 655, 5, 45, 0, 0, 655, 658, 1, 0, 0, 0, 656,
		658, 5, 27, 0, 0, 657, 617, 1, 0, 0, 0, 657, 618, 1, 0, 0, 0, 657, 625,
		1, 0, 0, 0, 657, 639, 1, 0, 0, 0, 657, 643, 1, 0, 0, 0, 657, 646, 1, 0,
		0, 0, 657, 652, 1, 0, 0, 0, 657, 656, 1, 0, 0, 0, 658, 113, 1, 0, 0, 0,
		659, 664, 5, 55, 0, 0, 660, 661, 5, 37, 0, 0, 661, 663, 5, 55, 0, 0, 662,
		660, 1, 0, 0, 0, 663, 666, 1, 0, 0, 0, 664, 662, 1, 0, 0, 0, 664, 665,
		1, 0, 0, 0, 665, 115, 1, 0, 0, 0, 666, 664, 1, 0, 0, 0, 667, 669, 3, 2,
		1, 0, 668, 667, 1, 0, 0, 0, 669, 672, 1, 0, 0, 0, 670, 668, 1, 0, 0, 0,
		670, 671, 1, 0, 0, 0, 671, 674, 1, 0, 0, 0, 672, 670, 1, 0, 0, 0, 673,
		675, 3, 56, 28, 0, 674, 673, 1, 0, 0, 0, 674, 675, 1, 0, 0, 0, 675, 117,
		1, 0, 0, 0, 75, 121, 133, 141, 148, 154, 160, 165, 173, 176, 183, 187,
		191, 196, 202, 208, 211, 218, 223, 230, 242, 248, 251, 259, 264, 271, 279,
		291, 298, 305, 313, 324, 336, 345, 352, 359, 372, 378, 385, 396, 405, 413,
		421, 425, 434, 442, 446, 449, 456, 461, 466, 469, 471, 479, 482, 489, 499,
		530, 534, 549, 559, 564, 570, 575, 591, 594, 605, 611, 623, 632, 637, 641,
		657, 664, 670, 674,
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
	ospreyParserANY_CHAR            = 59
)

// ospreyParser rules.
const (
	ospreyParserRULE_program           = 0
	ospreyParserRULE_statement         = 1
	ospreyParserRULE_importStmt        = 2
	ospreyParserRULE_letDecl           = 3
	ospreyParserRULE_fnDecl            = 4
	ospreyParserRULE_pluginContent     = 5
	ospreyParserRULE_pluginToken       = 6
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
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
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
	PluginContent() IPluginContentContext
	PluginReturnType() IPluginReturnTypeContext

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

func (s *FnDeclContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(ospreyParserID)
}

func (s *FnDeclContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(ospreyParserID, i)
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

	p.SetState(191)
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
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(158)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(160)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(159)
				p.ParamList()
			}

		}
		{
			p.SetState(162)
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

		if _la == ospreyParserARROW {
			{
				p.SetState(163)
				p.Match(ospreyParserARROW)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(164)
				p.Type_()
			}

		}
		p.SetState(173)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case ospreyParserEQ:
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
				p.Expr()
			}

		case ospreyParserLBRACE:
			{
				p.SetState(169)
				p.Match(ospreyParserLBRACE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(170)
				p.BlockBody()
			}
			{
				p.SetState(171)
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

	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(176)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserDOC_COMMENT {
			{
				p.SetState(175)
				p.DocComment()
			}

		}
		{
			p.SetState(178)
			p.Match(ospreyParserFN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(179)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(180)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(181)
			p.Match(ospreyParserLPAREN)
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

		if _la == ospreyParserID {
			{
				p.SetState(182)
				p.ParamList()
			}

		}
		{
			p.SetState(185)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(187)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserAS || _la == ospreyParserARROW {
			{
				p.SetState(186)
				p.PluginReturnType()
			}

		}
		{
			p.SetState(189)
			p.Match(ospreyParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(190)
			p.PluginContent()
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
	SEMI() antlr.TerminalNode
	AllPluginToken() []IPluginTokenContext
	PluginToken(i int) IPluginTokenContext

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

func (s *PluginContentContext) SEMI() antlr.TerminalNode {
	return s.GetToken(ospreyParserSEMI, 0)
}

func (s *PluginContentContext) AllPluginToken() []IPluginTokenContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPluginTokenContext); ok {
			len++
		}
	}

	tst := make([]IPluginTokenContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPluginTokenContext); ok {
			tst[i] = t.(IPluginTokenContext)
			i++
		}
	}

	return tst
}

func (s *PluginContentContext) PluginToken(i int) IPluginTokenContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPluginTokenContext); ok {
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

	return t.(IPluginTokenContext)
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
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(194)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1152921435887370238) != 0) {
		{
			p.SetState(193)
			p.PluginToken()
		}

		p.SetState(196)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(198)
		p.Match(ospreyParserSEMI)
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

// IPluginTokenContext is an interface to support dynamic dispatch.
type IPluginTokenContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SEMI() antlr.TerminalNode
	ANY_CHAR() antlr.TerminalNode

	// IsPluginTokenContext differentiates from other interfaces.
	IsPluginTokenContext()
}

type PluginTokenContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPluginTokenContext() *PluginTokenContext {
	var p = new(PluginTokenContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pluginToken
	return p
}

func InitEmptyPluginTokenContext(p *PluginTokenContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_pluginToken
}

func (*PluginTokenContext) IsPluginTokenContext() {}

func NewPluginTokenContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PluginTokenContext {
	var p = new(PluginTokenContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_pluginToken

	return p
}

func (s *PluginTokenContext) GetParser() antlr.Parser { return s.parser }

func (s *PluginTokenContext) SEMI() antlr.TerminalNode {
	return s.GetToken(ospreyParserSEMI, 0)
}

func (s *PluginTokenContext) ANY_CHAR() antlr.TerminalNode {
	return s.GetToken(ospreyParserANY_CHAR, 0)
}

func (s *PluginTokenContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PluginTokenContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PluginTokenContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterPluginToken(s)
	}
}

func (s *PluginTokenContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitPluginToken(s)
	}
}

func (p *ospreyParser) PluginToken() (localctx IPluginTokenContext) {
	localctx = NewPluginTokenContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, ospreyParserRULE_pluginToken)
	var _la int

	p.SetState(202)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(200)
			_la = p.GetTokenStream().LA(1)

			if _la <= 0 || _la == ospreyParserSEMI {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(201)
			p.Match(ospreyParserANY_CHAR)
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
	p.SetState(208)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserARROW:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(204)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(205)
			p.Type_()
		}

	case ospreyParserAS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(206)
			p.Match(ospreyParserAS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(207)
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
	p.SetState(211)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(210)
			p.DocComment()
		}

	}
	{
		p.SetState(213)
		p.Match(ospreyParserEXTERN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(214)
		p.Match(ospreyParserFN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(215)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(216)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(218)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserID {
		{
			p.SetState(217)
			p.ExternParamList()
		}

	}
	{
		p.SetState(220)
		p.Match(ospreyParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(223)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserARROW {
		{
			p.SetState(221)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(222)
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
		p.SetState(225)
		p.ExternParam()
	}
	p.SetState(230)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(226)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(227)
			p.ExternParam()
		}

		p.SetState(232)
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
		p.SetState(233)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(234)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(235)
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
		p.SetState(237)
		p.Param()
	}
	p.SetState(242)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(238)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(239)
			p.Param()
		}

		p.SetState(244)
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
		p.SetState(245)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(248)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserCOLON {
		{
			p.SetState(246)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(247)
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
	p.SetState(251)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(250)
			p.DocComment()
		}

	}
	{
		p.SetState(253)
		p.Match(ospreyParserTYPE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(254)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(259)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserLT {
		{
			p.SetState(255)
			p.Match(ospreyParserLT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(256)
			p.TypeParamList()
		}
		{
			p.SetState(257)
			p.Match(ospreyParserGT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(261)
		p.Match(ospreyParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(264)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserID:
		{
			p.SetState(262)
			p.UnionType()
		}

	case ospreyParserLBRACE:
		{
			p.SetState(263)
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
		p.SetState(266)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(271)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(267)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(268)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(273)
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
		p.SetState(274)
		p.Variant()
	}
	p.SetState(279)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(275)
				p.Match(ospreyParserBAR)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(276)
				p.Variant()
			}

		}
		p.SetState(281)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext())
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
		p.SetState(282)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(283)
		p.FieldDeclarations()
	}
	{
		p.SetState(284)
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
		p.SetState(286)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(291)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(287)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(288)
			p.FieldDeclarations()
		}
		{
			p.SetState(289)
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
		p.SetState(293)
		p.FieldDeclaration()
	}
	p.SetState(298)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(294)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(295)
			p.FieldDeclaration()
		}

		p.SetState(300)
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
		p.SetState(301)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(302)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(303)
		p.Type_()
	}
	p.SetState(305)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserWHERE {
		{
			p.SetState(304)
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
		p.SetState(307)
		p.Match(ospreyParserWHERE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(308)
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
		p.SetState(310)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(311)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(313)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
		{
			p.SetState(312)
			p.ArgList()
		}

	}
	{
		p.SetState(315)
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
		p.SetState(317)
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
		p.SetState(319)
		p.Field()
	}
	p.SetState(324)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(320)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(321)
			p.Field()
		}

		p.SetState(326)
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
		p.SetState(327)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(328)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(329)
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

	p.SetState(352)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(331)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(336)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLT {
			{
				p.SetState(332)
				p.Match(ospreyParserLT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(333)
				p.TypeList()
			}
			{
				p.SetState(334)
				p.Match(ospreyParserGT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(338)
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(339)
			p.Match(ospreyParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(340)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(345)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLT {
			{
				p.SetState(341)
				p.Match(ospreyParserLT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(342)
				p.TypeList()
			}
			{
				p.SetState(343)
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
			p.SetState(347)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(348)
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(349)
			p.Type_()
		}
		{
			p.SetState(350)
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
		p.SetState(354)
		p.Type_()
	}
	p.SetState(359)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(355)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(356)
			p.Type_()
		}

		p.SetState(361)
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
		p.SetState(362)
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
		p.SetState(364)
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

	p.SetState(378)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(366)
			p.Match(ospreyParserMATCH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(367)
			p.Expr()
		}
		{
			p.SetState(368)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(370)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420968061394952) != 0) {
			{
				p.SetState(369)
				p.MatchArm()
			}

			p.SetState(372)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(374)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(376)
			p.SelectExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(377)
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
		p.SetState(380)
		p.Match(ospreyParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(381)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(383)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420968061394952) != 0) {
		{
			p.SetState(382)
			p.SelectArm()
		}

		p.SetState(385)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(387)
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
	p.SetState(396)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(389)
			p.Pattern()
		}
		{
			p.SetState(390)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(391)
			p.Expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(393)
			p.Match(ospreyParserUNDERSCORE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(394)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(395)
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
		p.SetState(398)
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
		p.SetState(400)
		p.AddExpr()
	}
	p.SetState(405)
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
				p.SetState(401)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3306587947008) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(402)
				p.AddExpr()
			}

		}
		p.SetState(407)
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
		p.SetState(408)
		p.MulExpr()
	}
	p.SetState(413)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(409)
				_la = p.GetTokenStream().LA(1)

				if !(_la == ospreyParserPLUS || _la == ospreyParserMINUS) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(410)
				p.MulExpr()
			}

		}
		p.SetState(415)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext())
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
		p.SetState(416)
		p.UnaryExpr()
	}
	p.SetState(421)
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
				p.SetState(417)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3377716900397056) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(418)
				p.UnaryExpr()
			}

		}
		p.SetState(423)
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
	p.SetState(425)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(424)
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
		p.SetState(427)
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
		p.SetState(429)
		p.CallExpr()
	}
	p.SetState(434)
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
				p.SetState(430)
				p.Match(ospreyParserPIPE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(431)
				p.CallExpr()
			}

		}
		p.SetState(436)
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

	p.SetState(471)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(437)
			p.Primary()
		}
		p.SetState(440)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(438)
					p.Match(ospreyParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(439)
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

			p.SetState(442)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		p.SetState(449)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(444)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(446)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
				{
					p.SetState(445)
					p.ArgList()
				}

			}
			{
				p.SetState(448)
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
			p.SetState(451)
			p.Primary()
		}
		p.SetState(459)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(452)
					p.Match(ospreyParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(453)
					p.Match(ospreyParserID)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				{
					p.SetState(454)
					p.Match(ospreyParserLPAREN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(456)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
					{
						p.SetState(455)
						p.ArgList()
					}

				}
				{
					p.SetState(458)
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

			p.SetState(461)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(463)
			p.Primary()
		}
		p.SetState(469)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(464)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(466)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
				{
					p.SetState(465)
					p.ArgList()
				}

			}
			{
				p.SetState(468)
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

	p.SetState(482)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 53, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(473)
			p.NamedArgList()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(474)
			p.Expr()
		}
		p.SetState(479)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == ospreyParserCOMMA {
			{
				p.SetState(475)
				p.Match(ospreyParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(476)
				p.Expr()
			}

			p.SetState(481)
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
		p.SetState(484)
		p.NamedArg()
	}
	p.SetState(487)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserCOMMA {
		{
			p.SetState(485)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(486)
			p.NamedArg()
		}

		p.SetState(489)
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
		p.SetState(491)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(492)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(493)
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
	p.SetState(530)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 56, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(495)
			p.Match(ospreyParserSPAWN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(496)
			p.Expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(497)
			p.Match(ospreyParserYIELD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(499)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 55, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(498)
				p.Expr()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(501)
			p.Match(ospreyParserAWAIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(502)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(503)
			p.Expr()
		}
		{
			p.SetState(504)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(506)
			p.Match(ospreyParserSEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(507)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(508)
			p.Expr()
		}
		{
			p.SetState(509)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(510)
			p.Expr()
		}
		{
			p.SetState(511)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(513)
			p.Match(ospreyParserRECV)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(514)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(515)
			p.Expr()
		}
		{
			p.SetState(516)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(518)
			p.Match(ospreyParserSELECT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(519)
			p.SelectExpr()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(520)
			p.TypeConstructor()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(521)
			p.UpdateExpr()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(522)
			p.BlockExpr()
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(523)
			p.Literal()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(524)
			p.LambdaExpr()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(525)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(526)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(527)
			p.Expr()
		}
		{
			p.SetState(528)
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
		p.SetState(532)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(534)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserLT {
		{
			p.SetState(533)
			p.TypeArgs()
		}

	}
	{
		p.SetState(536)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(537)
		p.FieldAssignments()
	}
	{
		p.SetState(538)
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
		p.SetState(540)
		p.Match(ospreyParserLT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(541)
		p.TypeList()
	}
	{
		p.SetState(542)
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
		p.SetState(544)
		p.FieldAssignment()
	}
	p.SetState(549)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(545)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(546)
			p.FieldAssignment()
		}

		p.SetState(551)
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
		p.SetState(552)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(553)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(554)
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

	p.SetState(575)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserFN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(556)
			p.Match(ospreyParserFN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(557)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(559)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(558)
				p.ParamList()
			}

		}
		{
			p.SetState(561)
			p.Match(ospreyParserRPAREN)
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

		if _la == ospreyParserARROW {
			{
				p.SetState(562)
				p.Match(ospreyParserARROW)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(563)
				p.Type_()
			}

		}
		{
			p.SetState(566)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(567)
			p.Expr()
		}

	case ospreyParserBAR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(568)
			p.Match(ospreyParserBAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(570)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(569)
				p.ParamList()
			}

		}
		{
			p.SetState(572)
			p.Match(ospreyParserBAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(573)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(574)
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
		p.SetState(577)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(578)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(579)
		p.FieldAssignments()
	}
	{
		p.SetState(580)
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
		p.SetState(582)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(583)
		p.BlockBody()
	}
	{
		p.SetState(584)
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
		p.SetState(586)
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
	p.SetState(589)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(588)
			p.Match(ospreyParserDOC_COMMENT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(591)
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
	p.SetState(594)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(593)
			p.DocComment()
		}

	}
	{
		p.SetState(596)
		p.Match(ospreyParserMODULE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(597)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(598)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(599)
		p.ModuleBody()
	}
	{
		p.SetState(600)
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
	p.SetState(605)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&144115188075856712) != 0 {
		{
			p.SetState(602)
			p.ModuleStatement()
		}

		p.SetState(607)
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
	p.SetState(611)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 66, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(608)
			p.LetDecl()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(609)
			p.FnDecl()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(610)
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
		p.SetState(613)
		p.Pattern()
	}
	{
		p.SetState(614)
		p.Match(ospreyParserLAMBDA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(615)
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

	p.SetState(657)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 71, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(617)
			p.UnaryExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(618)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(623)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLBRACE {
			{
				p.SetState(619)
				p.Match(ospreyParserLBRACE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(620)
				p.FieldPattern()
			}
			{
				p.SetState(621)
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
			p.SetState(625)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(637)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLPAREN {
			{
				p.SetState(626)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(627)
				p.Pattern()
			}
			p.SetState(632)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == ospreyParserCOMMA {
				{
					p.SetState(628)
					p.Match(ospreyParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(629)
					p.Pattern()
				}

				p.SetState(634)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(635)
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
			p.SetState(639)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(641)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(640)
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
			p.SetState(643)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(644)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(645)
			p.Type_()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(646)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(647)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(648)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(649)
			p.FieldPattern()
		}
		{
			p.SetState(650)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(652)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(653)
			p.FieldPattern()
		}
		{
			p.SetState(654)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(656)
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
		p.SetState(659)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(664)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(660)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(661)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(666)
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
	p.SetState(670)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 73, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(667)
				p.Statement()
			}

		}
		p.SetState(672)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 73, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(674)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&68420967927177228) != 0 {
		{
			p.SetState(673)
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
