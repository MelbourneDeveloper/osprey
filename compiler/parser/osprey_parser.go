// Code generated from osprey.g4 by ANTLR 4.13.1. DO NOT EDIT.

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
		"'let'", "'mut'", "'if'", "'else'", "'spawn'", "'yield'", "'await'",
		"'fiber'", "'channel'", "'send'", "'recv'", "'select'", "'true'", "'false'",
		"'where'", "'->'", "'=>'", "'_'", "'='", "'=='", "'!='", "'<='", "'>='",
		"'!'", "'%'", "':'", "';'", "','", "'.'", "'|'", "'<'", "'>'", "'('",
		"')'", "'{'", "'}'", "'['", "']'", "'+'", "'-'", "'*'", "'/'",
	}
	staticData.SymbolicNames = []string{
		"", "PIPE", "MATCH", "FN", "EXTERN", "IMPORT", "TYPE", "MODULE", "LET",
		"MUT", "IF", "ELSE", "SPAWN", "YIELD", "AWAIT", "FIBER", "CHANNEL",
		"SEND", "RECV", "SELECT", "TRUE", "FALSE", "WHERE", "ARROW", "LAMBDA",
		"UNDERSCORE", "EQ", "EQ_OP", "NE_OP", "LE_OP", "GE_OP", "NOT_OP", "MOD_OP",
		"COLON", "SEMI", "COMMA", "DOT", "BAR", "LT", "GT", "LPAREN", "RPAREN",
		"LBRACE", "RBRACE", "LSQUARE", "RSQUARE", "PLUS", "MINUS", "STAR", "SLASH",
		"INT", "INTERPOLATED_STRING", "STRING", "ID", "WS", "DOC_COMMENT", "COMMENT",
	}
	staticData.RuleNames = []string{
		"program", "statement", "importStmt", "letDecl", "fnDecl", "externDecl",
		"externParamList", "externParam", "paramList", "param", "typeDecl",
		"typeParamList", "unionType", "recordType", "variant", "fieldDeclarations",
		"fieldDeclaration", "constraint", "functionCall", "booleanExpr", "fieldList",
		"field", "type", "typeList", "exprStmt", "expr", "matchExpr", "selectExpr",
		"selectArm", "binaryExpr", "comparisonExpr", "addExpr", "mulExpr", "unaryExpr",
		"pipeExpr", "callExpr", "argList", "namedArgList", "namedArg", "primary",
		"typeConstructor", "typeArgs", "fieldAssignments", "fieldAssignment",
		"lambdaExpr", "updateExpr", "blockExpr", "literal", "listLiteral", "docComment",
		"moduleDecl", "moduleBody", "moduleStatement", "matchArm", "pattern",
		"fieldPattern", "blockBody",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 56, 660, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 1, 0, 5, 0,
		116, 8, 0, 10, 0, 12, 0, 119, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 3, 1, 130, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 136, 8,
		2, 10, 2, 12, 2, 139, 9, 2, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 145, 8, 3, 1,
		3, 1, 3, 1, 3, 1, 4, 3, 4, 151, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 157,
		8, 4, 1, 4, 1, 4, 1, 4, 3, 4, 162, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4,
		1, 4, 3, 4, 170, 8, 4, 1, 5, 3, 5, 173, 8, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1,
		5, 3, 5, 180, 8, 5, 1, 5, 1, 5, 1, 5, 3, 5, 185, 8, 5, 1, 6, 1, 6, 1, 6,
		5, 6, 190, 8, 6, 10, 6, 12, 6, 193, 9, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8,
		1, 8, 1, 8, 5, 8, 202, 8, 8, 10, 8, 12, 8, 205, 9, 8, 1, 9, 1, 9, 1, 9,
		3, 9, 210, 8, 9, 1, 10, 3, 10, 213, 8, 10, 1, 10, 1, 10, 1, 10, 1, 10,
		1, 10, 1, 10, 3, 10, 221, 8, 10, 1, 10, 1, 10, 1, 10, 3, 10, 226, 8, 10,
		1, 11, 1, 11, 1, 11, 5, 11, 231, 8, 11, 10, 11, 12, 11, 234, 9, 11, 1,
		12, 1, 12, 1, 12, 5, 12, 239, 8, 12, 10, 12, 12, 12, 242, 9, 12, 1, 13,
		1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 3, 14, 253, 8,
		14, 1, 15, 1, 15, 1, 15, 5, 15, 258, 8, 15, 10, 15, 12, 15, 261, 9, 15,
		1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 267, 8, 16, 1, 17, 1, 17, 1, 17, 1,
		18, 1, 18, 1, 18, 3, 18, 275, 8, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1, 20,
		1, 20, 1, 20, 5, 20, 284, 8, 20, 10, 20, 12, 20, 287, 9, 20, 1, 21, 1,
		21, 1, 21, 1, 21, 1, 22, 1, 22, 3, 22, 295, 8, 22, 1, 22, 1, 22, 1, 22,
		1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 3, 22, 305, 8, 22, 1, 22, 1, 22, 1,
		22, 1, 22, 1, 22, 1, 22, 3, 22, 313, 8, 22, 1, 23, 1, 23, 1, 23, 5, 23,
		318, 8, 23, 10, 23, 12, 23, 321, 9, 23, 1, 24, 1, 24, 1, 25, 1, 25, 1,
		26, 1, 26, 1, 26, 1, 26, 4, 26, 331, 8, 26, 11, 26, 12, 26, 332, 1, 26,
		1, 26, 1, 26, 1, 26, 3, 26, 339, 8, 26, 1, 27, 1, 27, 1, 27, 4, 27, 344,
		8, 27, 11, 27, 12, 27, 345, 1, 27, 1, 27, 1, 28, 1, 28, 1, 28, 1, 28, 1,
		28, 1, 28, 1, 28, 3, 28, 357, 8, 28, 1, 29, 1, 29, 1, 30, 1, 30, 1, 30,
		5, 30, 364, 8, 30, 10, 30, 12, 30, 367, 9, 30, 1, 31, 1, 31, 1, 31, 5,
		31, 372, 8, 31, 10, 31, 12, 31, 375, 9, 31, 1, 32, 1, 32, 1, 32, 5, 32,
		380, 8, 32, 10, 32, 12, 32, 383, 9, 32, 1, 33, 3, 33, 386, 8, 33, 1, 33,
		1, 33, 1, 34, 1, 34, 1, 34, 5, 34, 393, 8, 34, 10, 34, 12, 34, 396, 9,
		34, 1, 35, 1, 35, 1, 35, 4, 35, 401, 8, 35, 11, 35, 12, 35, 402, 1, 35,
		1, 35, 3, 35, 407, 8, 35, 1, 35, 3, 35, 410, 8, 35, 1, 35, 1, 35, 1, 35,
		1, 35, 1, 35, 3, 35, 417, 8, 35, 1, 35, 4, 35, 420, 8, 35, 11, 35, 12,
		35, 421, 1, 35, 1, 35, 1, 35, 3, 35, 427, 8, 35, 1, 35, 3, 35, 430, 8,
		35, 3, 35, 432, 8, 35, 1, 36, 1, 36, 1, 36, 1, 36, 5, 36, 438, 8, 36, 10,
		36, 12, 36, 441, 9, 36, 3, 36, 443, 8, 36, 1, 37, 1, 37, 1, 37, 4, 37,
		448, 8, 37, 11, 37, 12, 37, 449, 1, 38, 1, 38, 1, 38, 1, 38, 1, 39, 1,
		39, 1, 39, 1, 39, 3, 39, 460, 8, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39,
		1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1,
		39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39,
		1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 3, 39, 495, 8, 39, 1,
		40, 1, 40, 3, 40, 499, 8, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 41, 1, 41,
		1, 41, 1, 41, 1, 42, 1, 42, 1, 42, 5, 42, 512, 8, 42, 10, 42, 12, 42, 515,
		9, 42, 1, 43, 1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1, 44, 3, 44, 524, 8,
		44, 1, 44, 1, 44, 1, 44, 3, 44, 529, 8, 44, 1, 44, 1, 44, 1, 44, 1, 44,
		3, 44, 535, 8, 44, 1, 44, 1, 44, 1, 44, 3, 44, 540, 8, 44, 1, 45, 1, 45,
		1, 45, 1, 45, 1, 45, 1, 46, 1, 46, 1, 46, 1, 46, 1, 47, 1, 47, 1, 47, 1,
		47, 1, 47, 1, 47, 3, 47, 557, 8, 47, 1, 48, 1, 48, 1, 48, 1, 48, 5, 48,
		563, 8, 48, 10, 48, 12, 48, 566, 9, 48, 3, 48, 568, 8, 48, 1, 48, 1, 48,
		1, 49, 4, 49, 573, 8, 49, 11, 49, 12, 49, 574, 1, 50, 3, 50, 578, 8, 50,
		1, 50, 1, 50, 1, 50, 1, 50, 1, 50, 1, 50, 1, 51, 5, 51, 587, 8, 51, 10,
		51, 12, 51, 590, 9, 51, 1, 52, 1, 52, 1, 52, 3, 52, 595, 8, 52, 1, 53,
		1, 53, 1, 53, 1, 53, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 3, 54, 607,
		8, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 5, 54, 614, 8, 54, 10, 54, 12,
		54, 617, 9, 54, 1, 54, 1, 54, 3, 54, 621, 8, 54, 1, 54, 1, 54, 3, 54, 625,
		8, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1,
		54, 1, 54, 1, 54, 1, 54, 1, 54, 3, 54, 641, 8, 54, 1, 55, 1, 55, 1, 55,
		5, 55, 646, 8, 55, 10, 55, 12, 55, 649, 9, 55, 1, 56, 5, 56, 652, 8, 56,
		10, 56, 12, 56, 655, 9, 56, 1, 56, 3, 56, 658, 8, 56, 1, 56, 0, 0, 57,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36,
		38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72,
		74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106,
		108, 110, 112, 0, 5, 1, 0, 8, 9, 2, 0, 27, 30, 38, 39, 1, 0, 46, 47, 2,
		0, 32, 32, 48, 49, 3, 0, 14, 14, 31, 31, 46, 47, 705, 0, 117, 1, 0, 0,
		0, 2, 129, 1, 0, 0, 0, 4, 131, 1, 0, 0, 0, 6, 140, 1, 0, 0, 0, 8, 150,
		1, 0, 0, 0, 10, 172, 1, 0, 0, 0, 12, 186, 1, 0, 0, 0, 14, 194, 1, 0, 0,
		0, 16, 198, 1, 0, 0, 0, 18, 206, 1, 0, 0, 0, 20, 212, 1, 0, 0, 0, 22, 227,
		1, 0, 0, 0, 24, 235, 1, 0, 0, 0, 26, 243, 1, 0, 0, 0, 28, 247, 1, 0, 0,
		0, 30, 254, 1, 0, 0, 0, 32, 262, 1, 0, 0, 0, 34, 268, 1, 0, 0, 0, 36, 271,
		1, 0, 0, 0, 38, 278, 1, 0, 0, 0, 40, 280, 1, 0, 0, 0, 42, 288, 1, 0, 0,
		0, 44, 312, 1, 0, 0, 0, 46, 314, 1, 0, 0, 0, 48, 322, 1, 0, 0, 0, 50, 324,
		1, 0, 0, 0, 52, 338, 1, 0, 0, 0, 54, 340, 1, 0, 0, 0, 56, 356, 1, 0, 0,
		0, 58, 358, 1, 0, 0, 0, 60, 360, 1, 0, 0, 0, 62, 368, 1, 0, 0, 0, 64, 376,
		1, 0, 0, 0, 66, 385, 1, 0, 0, 0, 68, 389, 1, 0, 0, 0, 70, 431, 1, 0, 0,
		0, 72, 442, 1, 0, 0, 0, 74, 444, 1, 0, 0, 0, 76, 451, 1, 0, 0, 0, 78, 494,
		1, 0, 0, 0, 80, 496, 1, 0, 0, 0, 82, 504, 1, 0, 0, 0, 84, 508, 1, 0, 0,
		0, 86, 516, 1, 0, 0, 0, 88, 539, 1, 0, 0, 0, 90, 541, 1, 0, 0, 0, 92, 546,
		1, 0, 0, 0, 94, 556, 1, 0, 0, 0, 96, 558, 1, 0, 0, 0, 98, 572, 1, 0, 0,
		0, 100, 577, 1, 0, 0, 0, 102, 588, 1, 0, 0, 0, 104, 594, 1, 0, 0, 0, 106,
		596, 1, 0, 0, 0, 108, 640, 1, 0, 0, 0, 110, 642, 1, 0, 0, 0, 112, 653,
		1, 0, 0, 0, 114, 116, 3, 2, 1, 0, 115, 114, 1, 0, 0, 0, 116, 119, 1, 0,
		0, 0, 117, 115, 1, 0, 0, 0, 117, 118, 1, 0, 0, 0, 118, 120, 1, 0, 0, 0,
		119, 117, 1, 0, 0, 0, 120, 121, 5, 0, 0, 1, 121, 1, 1, 0, 0, 0, 122, 130,
		3, 4, 2, 0, 123, 130, 3, 6, 3, 0, 124, 130, 3, 8, 4, 0, 125, 130, 3, 10,
		5, 0, 126, 130, 3, 20, 10, 0, 127, 130, 3, 100, 50, 0, 128, 130, 3, 48,
		24, 0, 129, 122, 1, 0, 0, 0, 129, 123, 1, 0, 0, 0, 129, 124, 1, 0, 0, 0,
		129, 125, 1, 0, 0, 0, 129, 126, 1, 0, 0, 0, 129, 127, 1, 0, 0, 0, 129,
		128, 1, 0, 0, 0, 130, 3, 1, 0, 0, 0, 131, 132, 5, 5, 0, 0, 132, 137, 5,
		53, 0, 0, 133, 134, 5, 36, 0, 0, 134, 136, 5, 53, 0, 0, 135, 133, 1, 0,
		0, 0, 136, 139, 1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 137, 138, 1, 0, 0, 0,
		138, 5, 1, 0, 0, 0, 139, 137, 1, 0, 0, 0, 140, 141, 7, 0, 0, 0, 141, 144,
		5, 53, 0, 0, 142, 143, 5, 33, 0, 0, 143, 145, 3, 44, 22, 0, 144, 142, 1,
		0, 0, 0, 144, 145, 1, 0, 0, 0, 145, 146, 1, 0, 0, 0, 146, 147, 5, 26, 0,
		0, 147, 148, 3, 50, 25, 0, 148, 7, 1, 0, 0, 0, 149, 151, 3, 98, 49, 0,
		150, 149, 1, 0, 0, 0, 150, 151, 1, 0, 0, 0, 151, 152, 1, 0, 0, 0, 152,
		153, 5, 3, 0, 0, 153, 154, 5, 53, 0, 0, 154, 156, 5, 40, 0, 0, 155, 157,
		3, 16, 8, 0, 156, 155, 1, 0, 0, 0, 156, 157, 1, 0, 0, 0, 157, 158, 1, 0,
		0, 0, 158, 161, 5, 41, 0, 0, 159, 160, 5, 23, 0, 0, 160, 162, 3, 44, 22,
		0, 161, 159, 1, 0, 0, 0, 161, 162, 1, 0, 0, 0, 162, 169, 1, 0, 0, 0, 163,
		164, 5, 26, 0, 0, 164, 170, 3, 50, 25, 0, 165, 166, 5, 42, 0, 0, 166, 167,
		3, 112, 56, 0, 167, 168, 5, 43, 0, 0, 168, 170, 1, 0, 0, 0, 169, 163, 1,
		0, 0, 0, 169, 165, 1, 0, 0, 0, 170, 9, 1, 0, 0, 0, 171, 173, 3, 98, 49,
		0, 172, 171, 1, 0, 0, 0, 172, 173, 1, 0, 0, 0, 173, 174, 1, 0, 0, 0, 174,
		175, 5, 4, 0, 0, 175, 176, 5, 3, 0, 0, 176, 177, 5, 53, 0, 0, 177, 179,
		5, 40, 0, 0, 178, 180, 3, 12, 6, 0, 179, 178, 1, 0, 0, 0, 179, 180, 1,
		0, 0, 0, 180, 181, 1, 0, 0, 0, 181, 184, 5, 41, 0, 0, 182, 183, 5, 23,
		0, 0, 183, 185, 3, 44, 22, 0, 184, 182, 1, 0, 0, 0, 184, 185, 1, 0, 0,
		0, 185, 11, 1, 0, 0, 0, 186, 191, 3, 14, 7, 0, 187, 188, 5, 35, 0, 0, 188,
		190, 3, 14, 7, 0, 189, 187, 1, 0, 0, 0, 190, 193, 1, 0, 0, 0, 191, 189,
		1, 0, 0, 0, 191, 192, 1, 0, 0, 0, 192, 13, 1, 0, 0, 0, 193, 191, 1, 0,
		0, 0, 194, 195, 5, 53, 0, 0, 195, 196, 5, 33, 0, 0, 196, 197, 3, 44, 22,
		0, 197, 15, 1, 0, 0, 0, 198, 203, 3, 18, 9, 0, 199, 200, 5, 35, 0, 0, 200,
		202, 3, 18, 9, 0, 201, 199, 1, 0, 0, 0, 202, 205, 1, 0, 0, 0, 203, 201,
		1, 0, 0, 0, 203, 204, 1, 0, 0, 0, 204, 17, 1, 0, 0, 0, 205, 203, 1, 0,
		0, 0, 206, 209, 5, 53, 0, 0, 207, 208, 5, 33, 0, 0, 208, 210, 3, 44, 22,
		0, 209, 207, 1, 0, 0, 0, 209, 210, 1, 0, 0, 0, 210, 19, 1, 0, 0, 0, 211,
		213, 3, 98, 49, 0, 212, 211, 1, 0, 0, 0, 212, 213, 1, 0, 0, 0, 213, 214,
		1, 0, 0, 0, 214, 215, 5, 6, 0, 0, 215, 220, 5, 53, 0, 0, 216, 217, 5, 38,
		0, 0, 217, 218, 3, 22, 11, 0, 218, 219, 5, 39, 0, 0, 219, 221, 1, 0, 0,
		0, 220, 216, 1, 0, 0, 0, 220, 221, 1, 0, 0, 0, 221, 222, 1, 0, 0, 0, 222,
		225, 5, 26, 0, 0, 223, 226, 3, 24, 12, 0, 224, 226, 3, 26, 13, 0, 225,
		223, 1, 0, 0, 0, 225, 224, 1, 0, 0, 0, 226, 21, 1, 0, 0, 0, 227, 232, 5,
		53, 0, 0, 228, 229, 5, 35, 0, 0, 229, 231, 5, 53, 0, 0, 230, 228, 1, 0,
		0, 0, 231, 234, 1, 0, 0, 0, 232, 230, 1, 0, 0, 0, 232, 233, 1, 0, 0, 0,
		233, 23, 1, 0, 0, 0, 234, 232, 1, 0, 0, 0, 235, 240, 3, 28, 14, 0, 236,
		237, 5, 37, 0, 0, 237, 239, 3, 28, 14, 0, 238, 236, 1, 0, 0, 0, 239, 242,
		1, 0, 0, 0, 240, 238, 1, 0, 0, 0, 240, 241, 1, 0, 0, 0, 241, 25, 1, 0,
		0, 0, 242, 240, 1, 0, 0, 0, 243, 244, 5, 42, 0, 0, 244, 245, 3, 30, 15,
		0, 245, 246, 5, 43, 0, 0, 246, 27, 1, 0, 0, 0, 247, 252, 5, 53, 0, 0, 248,
		249, 5, 42, 0, 0, 249, 250, 3, 30, 15, 0, 250, 251, 5, 43, 0, 0, 251, 253,
		1, 0, 0, 0, 252, 248, 1, 0, 0, 0, 252, 253, 1, 0, 0, 0, 253, 29, 1, 0,
		0, 0, 254, 259, 3, 32, 16, 0, 255, 256, 5, 35, 0, 0, 256, 258, 3, 32, 16,
		0, 257, 255, 1, 0, 0, 0, 258, 261, 1, 0, 0, 0, 259, 257, 1, 0, 0, 0, 259,
		260, 1, 0, 0, 0, 260, 31, 1, 0, 0, 0, 261, 259, 1, 0, 0, 0, 262, 263, 5,
		53, 0, 0, 263, 264, 5, 33, 0, 0, 264, 266, 3, 44, 22, 0, 265, 267, 3, 34,
		17, 0, 266, 265, 1, 0, 0, 0, 266, 267, 1, 0, 0, 0, 267, 33, 1, 0, 0, 0,
		268, 269, 5, 22, 0, 0, 269, 270, 3, 36, 18, 0, 270, 35, 1, 0, 0, 0, 271,
		272, 5, 53, 0, 0, 272, 274, 5, 40, 0, 0, 273, 275, 3, 72, 36, 0, 274, 273,
		1, 0, 0, 0, 274, 275, 1, 0, 0, 0, 275, 276, 1, 0, 0, 0, 276, 277, 5, 41,
		0, 0, 277, 37, 1, 0, 0, 0, 278, 279, 3, 60, 30, 0, 279, 39, 1, 0, 0, 0,
		280, 285, 3, 42, 21, 0, 281, 282, 5, 35, 0, 0, 282, 284, 3, 42, 21, 0,
		283, 281, 1, 0, 0, 0, 284, 287, 1, 0, 0, 0, 285, 283, 1, 0, 0, 0, 285,
		286, 1, 0, 0, 0, 286, 41, 1, 0, 0, 0, 287, 285, 1, 0, 0, 0, 288, 289, 5,
		53, 0, 0, 289, 290, 5, 33, 0, 0, 290, 291, 3, 44, 22, 0, 291, 43, 1, 0,
		0, 0, 292, 294, 5, 40, 0, 0, 293, 295, 3, 46, 23, 0, 294, 293, 1, 0, 0,
		0, 294, 295, 1, 0, 0, 0, 295, 296, 1, 0, 0, 0, 296, 297, 5, 41, 0, 0, 297,
		298, 5, 23, 0, 0, 298, 313, 3, 44, 22, 0, 299, 304, 5, 53, 0, 0, 300, 301,
		5, 38, 0, 0, 301, 302, 3, 46, 23, 0, 302, 303, 5, 39, 0, 0, 303, 305, 1,
		0, 0, 0, 304, 300, 1, 0, 0, 0, 304, 305, 1, 0, 0, 0, 305, 313, 1, 0, 0,
		0, 306, 307, 5, 53, 0, 0, 307, 308, 5, 44, 0, 0, 308, 309, 3, 44, 22, 0,
		309, 310, 5, 45, 0, 0, 310, 313, 1, 0, 0, 0, 311, 313, 5, 53, 0, 0, 312,
		292, 1, 0, 0, 0, 312, 299, 1, 0, 0, 0, 312, 306, 1, 0, 0, 0, 312, 311,
		1, 0, 0, 0, 313, 45, 1, 0, 0, 0, 314, 319, 3, 44, 22, 0, 315, 316, 5, 35,
		0, 0, 316, 318, 3, 44, 22, 0, 317, 315, 1, 0, 0, 0, 318, 321, 1, 0, 0,
		0, 319, 317, 1, 0, 0, 0, 319, 320, 1, 0, 0, 0, 320, 47, 1, 0, 0, 0, 321,
		319, 1, 0, 0, 0, 322, 323, 3, 50, 25, 0, 323, 49, 1, 0, 0, 0, 324, 325,
		3, 52, 26, 0, 325, 51, 1, 0, 0, 0, 326, 327, 5, 2, 0, 0, 327, 328, 3, 50,
		25, 0, 328, 330, 5, 42, 0, 0, 329, 331, 3, 106, 53, 0, 330, 329, 1, 0,
		0, 0, 331, 332, 1, 0, 0, 0, 332, 330, 1, 0, 0, 0, 332, 333, 1, 0, 0, 0,
		333, 334, 1, 0, 0, 0, 334, 335, 5, 43, 0, 0, 335, 339, 1, 0, 0, 0, 336,
		339, 3, 54, 27, 0, 337, 339, 3, 58, 29, 0, 338, 326, 1, 0, 0, 0, 338, 336,
		1, 0, 0, 0, 338, 337, 1, 0, 0, 0, 339, 53, 1, 0, 0, 0, 340, 341, 5, 19,
		0, 0, 341, 343, 5, 42, 0, 0, 342, 344, 3, 56, 28, 0, 343, 342, 1, 0, 0,
		0, 344, 345, 1, 0, 0, 0, 345, 343, 1, 0, 0, 0, 345, 346, 1, 0, 0, 0, 346,
		347, 1, 0, 0, 0, 347, 348, 5, 43, 0, 0, 348, 55, 1, 0, 0, 0, 349, 350,
		3, 108, 54, 0, 350, 351, 5, 24, 0, 0, 351, 352, 3, 50, 25, 0, 352, 357,
		1, 0, 0, 0, 353, 354, 5, 25, 0, 0, 354, 355, 5, 24, 0, 0, 355, 357, 3,
		50, 25, 0, 356, 349, 1, 0, 0, 0, 356, 353, 1, 0, 0, 0, 357, 57, 1, 0, 0,
		0, 358, 359, 3, 60, 30, 0, 359, 59, 1, 0, 0, 0, 360, 365, 3, 62, 31, 0,
		361, 362, 7, 1, 0, 0, 362, 364, 3, 62, 31, 0, 363, 361, 1, 0, 0, 0, 364,
		367, 1, 0, 0, 0, 365, 363, 1, 0, 0, 0, 365, 366, 1, 0, 0, 0, 366, 61, 1,
		0, 0, 0, 367, 365, 1, 0, 0, 0, 368, 373, 3, 64, 32, 0, 369, 370, 7, 2,
		0, 0, 370, 372, 3, 64, 32, 0, 371, 369, 1, 0, 0, 0, 372, 375, 1, 0, 0,
		0, 373, 371, 1, 0, 0, 0, 373, 374, 1, 0, 0, 0, 374, 63, 1, 0, 0, 0, 375,
		373, 1, 0, 0, 0, 376, 381, 3, 66, 33, 0, 377, 378, 7, 3, 0, 0, 378, 380,
		3, 66, 33, 0, 379, 377, 1, 0, 0, 0, 380, 383, 1, 0, 0, 0, 381, 379, 1,
		0, 0, 0, 381, 382, 1, 0, 0, 0, 382, 65, 1, 0, 0, 0, 383, 381, 1, 0, 0,
		0, 384, 386, 7, 4, 0, 0, 385, 384, 1, 0, 0, 0, 385, 386, 1, 0, 0, 0, 386,
		387, 1, 0, 0, 0, 387, 388, 3, 68, 34, 0, 388, 67, 1, 0, 0, 0, 389, 394,
		3, 70, 35, 0, 390, 391, 5, 1, 0, 0, 391, 393, 3, 70, 35, 0, 392, 390, 1,
		0, 0, 0, 393, 396, 1, 0, 0, 0, 394, 392, 1, 0, 0, 0, 394, 395, 1, 0, 0,
		0, 395, 69, 1, 0, 0, 0, 396, 394, 1, 0, 0, 0, 397, 400, 3, 78, 39, 0, 398,
		399, 5, 36, 0, 0, 399, 401, 5, 53, 0, 0, 400, 398, 1, 0, 0, 0, 401, 402,
		1, 0, 0, 0, 402, 400, 1, 0, 0, 0, 402, 403, 1, 0, 0, 0, 403, 409, 1, 0,
		0, 0, 404, 406, 5, 40, 0, 0, 405, 407, 3, 72, 36, 0, 406, 405, 1, 0, 0,
		0, 406, 407, 1, 0, 0, 0, 407, 408, 1, 0, 0, 0, 408, 410, 5, 41, 0, 0, 409,
		404, 1, 0, 0, 0, 409, 410, 1, 0, 0, 0, 410, 432, 1, 0, 0, 0, 411, 419,
		3, 78, 39, 0, 412, 413, 5, 36, 0, 0, 413, 414, 5, 53, 0, 0, 414, 416, 5,
		40, 0, 0, 415, 417, 3, 72, 36, 0, 416, 415, 1, 0, 0, 0, 416, 417, 1, 0,
		0, 0, 417, 418, 1, 0, 0, 0, 418, 420, 5, 41, 0, 0, 419, 412, 1, 0, 0, 0,
		420, 421, 1, 0, 0, 0, 421, 419, 1, 0, 0, 0, 421, 422, 1, 0, 0, 0, 422,
		432, 1, 0, 0, 0, 423, 429, 3, 78, 39, 0, 424, 426, 5, 40, 0, 0, 425, 427,
		3, 72, 36, 0, 426, 425, 1, 0, 0, 0, 426, 427, 1, 0, 0, 0, 427, 428, 1,
		0, 0, 0, 428, 430, 5, 41, 0, 0, 429, 424, 1, 0, 0, 0, 429, 430, 1, 0, 0,
		0, 430, 432, 1, 0, 0, 0, 431, 397, 1, 0, 0, 0, 431, 411, 1, 0, 0, 0, 431,
		423, 1, 0, 0, 0, 432, 71, 1, 0, 0, 0, 433, 443, 3, 74, 37, 0, 434, 439,
		3, 50, 25, 0, 435, 436, 5, 35, 0, 0, 436, 438, 3, 50, 25, 0, 437, 435,
		1, 0, 0, 0, 438, 441, 1, 0, 0, 0, 439, 437, 1, 0, 0, 0, 439, 440, 1, 0,
		0, 0, 440, 443, 1, 0, 0, 0, 441, 439, 1, 0, 0, 0, 442, 433, 1, 0, 0, 0,
		442, 434, 1, 0, 0, 0, 443, 73, 1, 0, 0, 0, 444, 447, 3, 76, 38, 0, 445,
		446, 5, 35, 0, 0, 446, 448, 3, 76, 38, 0, 447, 445, 1, 0, 0, 0, 448, 449,
		1, 0, 0, 0, 449, 447, 1, 0, 0, 0, 449, 450, 1, 0, 0, 0, 450, 75, 1, 0,
		0, 0, 451, 452, 5, 53, 0, 0, 452, 453, 5, 33, 0, 0, 453, 454, 3, 50, 25,
		0, 454, 77, 1, 0, 0, 0, 455, 456, 5, 12, 0, 0, 456, 495, 3, 50, 25, 0,
		457, 459, 5, 13, 0, 0, 458, 460, 3, 50, 25, 0, 459, 458, 1, 0, 0, 0, 459,
		460, 1, 0, 0, 0, 460, 495, 1, 0, 0, 0, 461, 462, 5, 14, 0, 0, 462, 463,
		5, 40, 0, 0, 463, 464, 3, 50, 25, 0, 464, 465, 5, 41, 0, 0, 465, 495, 1,
		0, 0, 0, 466, 467, 5, 17, 0, 0, 467, 468, 5, 40, 0, 0, 468, 469, 3, 50,
		25, 0, 469, 470, 5, 35, 0, 0, 470, 471, 3, 50, 25, 0, 471, 472, 5, 41,
		0, 0, 472, 495, 1, 0, 0, 0, 473, 474, 5, 18, 0, 0, 474, 475, 5, 40, 0,
		0, 475, 476, 3, 50, 25, 0, 476, 477, 5, 41, 0, 0, 477, 495, 1, 0, 0, 0,
		478, 479, 5, 19, 0, 0, 479, 495, 3, 54, 27, 0, 480, 495, 3, 80, 40, 0,
		481, 495, 3, 90, 45, 0, 482, 495, 3, 92, 46, 0, 483, 495, 3, 94, 47, 0,
		484, 495, 3, 88, 44, 0, 485, 486, 5, 53, 0, 0, 486, 487, 5, 44, 0, 0, 487,
		488, 5, 50, 0, 0, 488, 495, 5, 45, 0, 0, 489, 495, 5, 53, 0, 0, 490, 491,
		5, 40, 0, 0, 491, 492, 3, 50, 25, 0, 492, 493, 5, 41, 0, 0, 493, 495, 1,
		0, 0, 0, 494, 455, 1, 0, 0, 0, 494, 457, 1, 0, 0, 0, 494, 461, 1, 0, 0,
		0, 494, 466, 1, 0, 0, 0, 494, 473, 1, 0, 0, 0, 494, 478, 1, 0, 0, 0, 494,
		480, 1, 0, 0, 0, 494, 481, 1, 0, 0, 0, 494, 482, 1, 0, 0, 0, 494, 483,
		1, 0, 0, 0, 494, 484, 1, 0, 0, 0, 494, 485, 1, 0, 0, 0, 494, 489, 1, 0,
		0, 0, 494, 490, 1, 0, 0, 0, 495, 79, 1, 0, 0, 0, 496, 498, 5, 53, 0, 0,
		497, 499, 3, 82, 41, 0, 498, 497, 1, 0, 0, 0, 498, 499, 1, 0, 0, 0, 499,
		500, 1, 0, 0, 0, 500, 501, 5, 42, 0, 0, 501, 502, 3, 84, 42, 0, 502, 503,
		5, 43, 0, 0, 503, 81, 1, 0, 0, 0, 504, 505, 5, 38, 0, 0, 505, 506, 3, 46,
		23, 0, 506, 507, 5, 39, 0, 0, 507, 83, 1, 0, 0, 0, 508, 513, 3, 86, 43,
		0, 509, 510, 5, 35, 0, 0, 510, 512, 3, 86, 43, 0, 511, 509, 1, 0, 0, 0,
		512, 515, 1, 0, 0, 0, 513, 511, 1, 0, 0, 0, 513, 514, 1, 0, 0, 0, 514,
		85, 1, 0, 0, 0, 515, 513, 1, 0, 0, 0, 516, 517, 5, 53, 0, 0, 517, 518,
		5, 33, 0, 0, 518, 519, 3, 50, 25, 0, 519, 87, 1, 0, 0, 0, 520, 521, 5,
		3, 0, 0, 521, 523, 5, 40, 0, 0, 522, 524, 3, 16, 8, 0, 523, 522, 1, 0,
		0, 0, 523, 524, 1, 0, 0, 0, 524, 525, 1, 0, 0, 0, 525, 528, 5, 41, 0, 0,
		526, 527, 5, 23, 0, 0, 527, 529, 3, 44, 22, 0, 528, 526, 1, 0, 0, 0, 528,
		529, 1, 0, 0, 0, 529, 530, 1, 0, 0, 0, 530, 531, 5, 24, 0, 0, 531, 540,
		3, 50, 25, 0, 532, 534, 5, 37, 0, 0, 533, 535, 3, 16, 8, 0, 534, 533, 1,
		0, 0, 0, 534, 535, 1, 0, 0, 0, 535, 536, 1, 0, 0, 0, 536, 537, 5, 37, 0,
		0, 537, 538, 5, 24, 0, 0, 538, 540, 3, 50, 25, 0, 539, 520, 1, 0, 0, 0,
		539, 532, 1, 0, 0, 0, 540, 89, 1, 0, 0, 0, 541, 542, 5, 53, 0, 0, 542,
		543, 5, 42, 0, 0, 543, 544, 3, 84, 42, 0, 544, 545, 5, 43, 0, 0, 545, 91,
		1, 0, 0, 0, 546, 547, 5, 42, 0, 0, 547, 548, 3, 112, 56, 0, 548, 549, 5,
		43, 0, 0, 549, 93, 1, 0, 0, 0, 550, 557, 5, 50, 0, 0, 551, 557, 5, 52,
		0, 0, 552, 557, 5, 51, 0, 0, 553, 557, 5, 20, 0, 0, 554, 557, 5, 21, 0,
		0, 555, 557, 3, 96, 48, 0, 556, 550, 1, 0, 0, 0, 556, 551, 1, 0, 0, 0,
		556, 552, 1, 0, 0, 0, 556, 553, 1, 0, 0, 0, 556, 554, 1, 0, 0, 0, 556,
		555, 1, 0, 0, 0, 557, 95, 1, 0, 0, 0, 558, 567, 5, 44, 0, 0, 559, 564,
		3, 50, 25, 0, 560, 561, 5, 35, 0, 0, 561, 563, 3, 50, 25, 0, 562, 560,
		1, 0, 0, 0, 563, 566, 1, 0, 0, 0, 564, 562, 1, 0, 0, 0, 564, 565, 1, 0,
		0, 0, 565, 568, 1, 0, 0, 0, 566, 564, 1, 0, 0, 0, 567, 559, 1, 0, 0, 0,
		567, 568, 1, 0, 0, 0, 568, 569, 1, 0, 0, 0, 569, 570, 5, 45, 0, 0, 570,
		97, 1, 0, 0, 0, 571, 573, 5, 55, 0, 0, 572, 571, 1, 0, 0, 0, 573, 574,
		1, 0, 0, 0, 574, 572, 1, 0, 0, 0, 574, 575, 1, 0, 0, 0, 575, 99, 1, 0,
		0, 0, 576, 578, 3, 98, 49, 0, 577, 576, 1, 0, 0, 0, 577, 578, 1, 0, 0,
		0, 578, 579, 1, 0, 0, 0, 579, 580, 5, 7, 0, 0, 580, 581, 5, 53, 0, 0, 581,
		582, 5, 42, 0, 0, 582, 583, 3, 102, 51, 0, 583, 584, 5, 43, 0, 0, 584,
		101, 1, 0, 0, 0, 585, 587, 3, 104, 52, 0, 586, 585, 1, 0, 0, 0, 587, 590,
		1, 0, 0, 0, 588, 586, 1, 0, 0, 0, 588, 589, 1, 0, 0, 0, 589, 103, 1, 0,
		0, 0, 590, 588, 1, 0, 0, 0, 591, 595, 3, 6, 3, 0, 592, 595, 3, 8, 4, 0,
		593, 595, 3, 20, 10, 0, 594, 591, 1, 0, 0, 0, 594, 592, 1, 0, 0, 0, 594,
		593, 1, 0, 0, 0, 595, 105, 1, 0, 0, 0, 596, 597, 3, 108, 54, 0, 597, 598,
		5, 24, 0, 0, 598, 599, 3, 50, 25, 0, 599, 107, 1, 0, 0, 0, 600, 641, 3,
		66, 33, 0, 601, 606, 5, 53, 0, 0, 602, 603, 5, 42, 0, 0, 603, 604, 3, 110,
		55, 0, 604, 605, 5, 43, 0, 0, 605, 607, 1, 0, 0, 0, 606, 602, 1, 0, 0,
		0, 606, 607, 1, 0, 0, 0, 607, 641, 1, 0, 0, 0, 608, 620, 5, 53, 0, 0, 609,
		610, 5, 40, 0, 0, 610, 615, 3, 108, 54, 0, 611, 612, 5, 35, 0, 0, 612,
		614, 3, 108, 54, 0, 613, 611, 1, 0, 0, 0, 614, 617, 1, 0, 0, 0, 615, 613,
		1, 0, 0, 0, 615, 616, 1, 0, 0, 0, 616, 618, 1, 0, 0, 0, 617, 615, 1, 0,
		0, 0, 618, 619, 5, 41, 0, 0, 619, 621, 1, 0, 0, 0, 620, 609, 1, 0, 0, 0,
		620, 621, 1, 0, 0, 0, 621, 641, 1, 0, 0, 0, 622, 624, 5, 53, 0, 0, 623,
		625, 5, 53, 0, 0, 624, 623, 1, 0, 0, 0, 624, 625, 1, 0, 0, 0, 625, 641,
		1, 0, 0, 0, 626, 627, 5, 53, 0, 0, 627, 628, 5, 33, 0, 0, 628, 641, 3,
		44, 22, 0, 629, 630, 5, 53, 0, 0, 630, 631, 5, 33, 0, 0, 631, 632, 5, 42,
		0, 0, 632, 633, 3, 110, 55, 0, 633, 634, 5, 43, 0, 0, 634, 641, 1, 0, 0,
		0, 635, 636, 5, 42, 0, 0, 636, 637, 3, 110, 55, 0, 637, 638, 5, 43, 0,
		0, 638, 641, 1, 0, 0, 0, 639, 641, 5, 25, 0, 0, 640, 600, 1, 0, 0, 0, 640,
		601, 1, 0, 0, 0, 640, 608, 1, 0, 0, 0, 640, 622, 1, 0, 0, 0, 640, 626,
		1, 0, 0, 0, 640, 629, 1, 0, 0, 0, 640, 635, 1, 0, 0, 0, 640, 639, 1, 0,
		0, 0, 641, 109, 1, 0, 0, 0, 642, 647, 5, 53, 0, 0, 643, 644, 5, 35, 0,
		0, 644, 646, 5, 53, 0, 0, 645, 643, 1, 0, 0, 0, 646, 649, 1, 0, 0, 0, 647,
		645, 1, 0, 0, 0, 647, 648, 1, 0, 0, 0, 648, 111, 1, 0, 0, 0, 649, 647,
		1, 0, 0, 0, 650, 652, 3, 2, 1, 0, 651, 650, 1, 0, 0, 0, 652, 655, 1, 0,
		0, 0, 653, 651, 1, 0, 0, 0, 653, 654, 1, 0, 0, 0, 654, 657, 1, 0, 0, 0,
		655, 653, 1, 0, 0, 0, 656, 658, 3, 50, 25, 0, 657, 656, 1, 0, 0, 0, 657,
		658, 1, 0, 0, 0, 658, 113, 1, 0, 0, 0, 71, 117, 129, 137, 144, 150, 156,
		161, 169, 172, 179, 184, 191, 203, 209, 212, 220, 225, 232, 240, 252, 259,
		266, 274, 285, 294, 304, 312, 319, 332, 338, 345, 356, 365, 373, 381, 385,
		394, 402, 406, 409, 416, 421, 426, 429, 431, 439, 442, 449, 459, 494, 498,
		513, 523, 528, 534, 539, 556, 564, 567, 574, 577, 588, 594, 606, 615, 620,
		624, 640, 647, 653, 657,
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
	ospreyParserSPAWN               = 12
	ospreyParserYIELD               = 13
	ospreyParserAWAIT               = 14
	ospreyParserFIBER               = 15
	ospreyParserCHANNEL             = 16
	ospreyParserSEND                = 17
	ospreyParserRECV                = 18
	ospreyParserSELECT              = 19
	ospreyParserTRUE                = 20
	ospreyParserFALSE               = 21
	ospreyParserWHERE               = 22
	ospreyParserARROW               = 23
	ospreyParserLAMBDA              = 24
	ospreyParserUNDERSCORE          = 25
	ospreyParserEQ                  = 26
	ospreyParserEQ_OP               = 27
	ospreyParserNE_OP               = 28
	ospreyParserLE_OP               = 29
	ospreyParserGE_OP               = 30
	ospreyParserNOT_OP              = 31
	ospreyParserMOD_OP              = 32
	ospreyParserCOLON               = 33
	ospreyParserSEMI                = 34
	ospreyParserCOMMA               = 35
	ospreyParserDOT                 = 36
	ospreyParserBAR                 = 37
	ospreyParserLT                  = 38
	ospreyParserGT                  = 39
	ospreyParserLPAREN              = 40
	ospreyParserRPAREN              = 41
	ospreyParserLBRACE              = 42
	ospreyParserRBRACE              = 43
	ospreyParserLSQUARE             = 44
	ospreyParserRSQUARE             = 45
	ospreyParserPLUS                = 46
	ospreyParserMINUS               = 47
	ospreyParserSTAR                = 48
	ospreyParserSLASH               = 49
	ospreyParserINT                 = 50
	ospreyParserINTERPOLATED_STRING = 51
	ospreyParserSTRING              = 52
	ospreyParserID                  = 53
	ospreyParserWS                  = 54
	ospreyParserDOC_COMMENT         = 55
	ospreyParserCOMMENT             = 56
)

// ospreyParser rules.
const (
	ospreyParserRULE_program           = 0
	ospreyParserRULE_statement         = 1
	ospreyParserRULE_importStmt        = 2
	ospreyParserRULE_letDecl           = 3
	ospreyParserRULE_fnDecl            = 4
	ospreyParserRULE_externDecl        = 5
	ospreyParserRULE_externParamList   = 6
	ospreyParserRULE_externParam       = 7
	ospreyParserRULE_paramList         = 8
	ospreyParserRULE_param             = 9
	ospreyParserRULE_typeDecl          = 10
	ospreyParserRULE_typeParamList     = 11
	ospreyParserRULE_unionType         = 12
	ospreyParserRULE_recordType        = 13
	ospreyParserRULE_variant           = 14
	ospreyParserRULE_fieldDeclarations = 15
	ospreyParserRULE_fieldDeclaration  = 16
	ospreyParserRULE_constraint        = 17
	ospreyParserRULE_functionCall      = 18
	ospreyParserRULE_booleanExpr       = 19
	ospreyParserRULE_fieldList         = 20
	ospreyParserRULE_field             = 21
	ospreyParserRULE_type              = 22
	ospreyParserRULE_typeList          = 23
	ospreyParserRULE_exprStmt          = 24
	ospreyParserRULE_expr              = 25
	ospreyParserRULE_matchExpr         = 26
	ospreyParserRULE_selectExpr        = 27
	ospreyParserRULE_selectArm         = 28
	ospreyParserRULE_binaryExpr        = 29
	ospreyParserRULE_comparisonExpr    = 30
	ospreyParserRULE_addExpr           = 31
	ospreyParserRULE_mulExpr           = 32
	ospreyParserRULE_unaryExpr         = 33
	ospreyParserRULE_pipeExpr          = 34
	ospreyParserRULE_callExpr          = 35
	ospreyParserRULE_argList           = 36
	ospreyParserRULE_namedArgList      = 37
	ospreyParserRULE_namedArg          = 38
	ospreyParserRULE_primary           = 39
	ospreyParserRULE_typeConstructor   = 40
	ospreyParserRULE_typeArgs          = 41
	ospreyParserRULE_fieldAssignments  = 42
	ospreyParserRULE_fieldAssignment   = 43
	ospreyParserRULE_lambdaExpr        = 44
	ospreyParserRULE_updateExpr        = 45
	ospreyParserRULE_blockExpr         = 46
	ospreyParserRULE_literal           = 47
	ospreyParserRULE_listLiteral       = 48
	ospreyParserRULE_docComment        = 49
	ospreyParserRULE_moduleDecl        = 50
	ospreyParserRULE_moduleBody        = 51
	ospreyParserRULE_moduleStatement   = 52
	ospreyParserRULE_matchArm          = 53
	ospreyParserRULE_pattern           = 54
	ospreyParserRULE_fieldPattern      = 55
	ospreyParserRULE_blockBody         = 56
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
	p.SetState(117)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&53151631188849660) != 0 {
		{
			p.SetState(114)
			p.Statement()
		}

		p.SetState(119)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(120)
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
	p.SetState(129)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(122)
			p.ImportStmt()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(123)
			p.LetDecl()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(124)
			p.FnDecl()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(125)
			p.ExternDecl()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(126)
			p.TypeDecl()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(127)
			p.ModuleDecl()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(128)
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
		p.SetState(131)
		p.Match(ospreyParserIMPORT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(132)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserDOT {
		{
			p.SetState(133)
			p.Match(ospreyParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(134)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(139)
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
		p.SetState(140)
		_la = p.GetTokenStream().LA(1)

		if !(_la == ospreyParserLET || _la == ospreyParserMUT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(141)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(144)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserCOLON {
		{
			p.SetState(142)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(143)
			p.Type_()
		}

	}
	{
		p.SetState(146)
		p.Match(ospreyParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(147)
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

	p.EnterOuterAlt(localctx, 1)
	p.SetState(150)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(149)
			p.DocComment()
		}

	}
	{
		p.SetState(152)
		p.Match(ospreyParserFN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(153)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(154)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserID {
		{
			p.SetState(155)
			p.ParamList()
		}

	}
	{
		p.SetState(158)
		p.Match(ospreyParserRPAREN)
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

	if _la == ospreyParserARROW {
		{
			p.SetState(159)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(160)
			p.Type_()
		}

	}
	p.SetState(169)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserEQ:
		{
			p.SetState(163)
			p.Match(ospreyParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(164)
			p.Expr()
		}

	case ospreyParserLBRACE:
		{
			p.SetState(165)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(166)
			p.BlockBody()
		}
		{
			p.SetState(167)
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
	p.EnterRule(localctx, 10, ospreyParserRULE_externDecl)
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
		p.Match(ospreyParserEXTERN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(175)
		p.Match(ospreyParserFN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(176)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(177)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserID {
		{
			p.SetState(178)
			p.ExternParamList()
		}

	}
	{
		p.SetState(181)
		p.Match(ospreyParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserARROW {
		{
			p.SetState(182)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(183)
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
	p.EnterRule(localctx, 12, ospreyParserRULE_externParamList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(186)
		p.ExternParam()
	}
	p.SetState(191)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(187)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(188)
			p.ExternParam()
		}

		p.SetState(193)
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
	p.EnterRule(localctx, 14, ospreyParserRULE_externParam)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(194)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(195)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(196)
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
	p.EnterRule(localctx, 16, ospreyParserRULE_paramList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(198)
		p.Param()
	}
	p.SetState(203)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(199)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(200)
			p.Param()
		}

		p.SetState(205)
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
	p.EnterRule(localctx, 18, ospreyParserRULE_param)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(206)
		p.Match(ospreyParserID)
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

	if _la == ospreyParserCOLON {
		{
			p.SetState(207)
			p.Match(ospreyParserCOLON)
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
	p.EnterRule(localctx, 20, ospreyParserRULE_typeDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(211)
			p.DocComment()
		}

	}
	{
		p.SetState(214)
		p.Match(ospreyParserTYPE)
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
	p.SetState(220)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserLT {
		{
			p.SetState(216)
			p.Match(ospreyParserLT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(217)
			p.TypeParamList()
		}
		{
			p.SetState(218)
			p.Match(ospreyParserGT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(222)
		p.Match(ospreyParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(225)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserID:
		{
			p.SetState(223)
			p.UnionType()
		}

	case ospreyParserLBRACE:
		{
			p.SetState(224)
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
	p.EnterRule(localctx, 22, ospreyParserRULE_typeParamList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(227)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(232)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(228)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(229)
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
	p.EnterRule(localctx, 24, ospreyParserRULE_unionType)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(235)
		p.Variant()
	}
	p.SetState(240)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(236)
				p.Match(ospreyParserBAR)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(237)
				p.Variant()
			}

		}
		p.SetState(242)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 26, ospreyParserRULE_recordType)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(243)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(244)
		p.FieldDeclarations()
	}
	{
		p.SetState(245)
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
	p.EnterRule(localctx, 28, ospreyParserRULE_variant)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(247)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(252)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(248)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(249)
			p.FieldDeclarations()
		}
		{
			p.SetState(250)
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
	p.EnterRule(localctx, 30, ospreyParserRULE_fieldDeclarations)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(254)
		p.FieldDeclaration()
	}
	p.SetState(259)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(255)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(256)
			p.FieldDeclaration()
		}

		p.SetState(261)
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
	p.EnterRule(localctx, 32, ospreyParserRULE_fieldDeclaration)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(262)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(263)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(264)
		p.Type_()
	}
	p.SetState(266)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserWHERE {
		{
			p.SetState(265)
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
	p.EnterRule(localctx, 34, ospreyParserRULE_constraint)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(268)
		p.Match(ospreyParserWHERE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(269)
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
	p.EnterRule(localctx, 36, ospreyParserRULE_functionCall)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(271)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(272)
		p.Match(ospreyParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(274)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17122834169884684) != 0 {
		{
			p.SetState(273)
			p.ArgList()
		}

	}
	{
		p.SetState(276)
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
	p.EnterRule(localctx, 38, ospreyParserRULE_booleanExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(278)
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
	p.EnterRule(localctx, 40, ospreyParserRULE_fieldList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(280)
		p.Field()
	}
	p.SetState(285)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(281)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(282)
			p.Field()
		}

		p.SetState(287)
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
	p.EnterRule(localctx, 42, ospreyParserRULE_field)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(288)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(289)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(290)
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
	p.EnterRule(localctx, 44, ospreyParserRULE_type)
	var _la int

	p.SetState(312)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(292)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(294)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLPAREN || _la == ospreyParserID {
			{
				p.SetState(293)
				p.TypeList()
			}

		}
		{
			p.SetState(296)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(297)
			p.Match(ospreyParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(298)
			p.Type_()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(299)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(304)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLT {
			{
				p.SetState(300)
				p.Match(ospreyParserLT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(301)
				p.TypeList()
			}
			{
				p.SetState(302)
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
			p.SetState(306)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(307)
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(308)
			p.Type_()
		}
		{
			p.SetState(309)
			p.Match(ospreyParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(311)
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
	p.EnterRule(localctx, 46, ospreyParserRULE_typeList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(314)
		p.Type_()
	}
	p.SetState(319)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(315)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(316)
			p.Type_()
		}

		p.SetState(321)
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
	p.EnterRule(localctx, 48, ospreyParserRULE_exprStmt)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(322)
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
	LoopExpr() ILoopExprContext

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

func (s *ExprContext) LoopExpr() ILoopExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILoopExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILoopExprContext)
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
	p.EnterRule(localctx, 50, ospreyParserRULE_expr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(324)
		p.MatchExpr()
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserMATCH, ospreyParserFN, ospreyParserSPAWN, ospreyParserYIELD, ospreyParserAWAIT, ospreyParserSEND, ospreyParserRECV, ospreyParserSELECT, ospreyParserTRUE, ospreyParserFALSE, ospreyParserNOT_OP, ospreyParserBAR, ospreyParserLPAREN, ospreyParserLBRACE, ospreyParserLSQUARE, ospreyParserPLUS, ospreyParserMINUS, ospreyParserINT, ospreyParserINTERPOLATED_STRING, ospreyParserSTRING, ospreyParserID:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(319)
			p.MatchExpr()
		}

	case ospreyParserLOOP:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(320)
			p.LoopExpr()
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
	p.EnterRule(localctx, 52, ospreyParserRULE_matchExpr)
	var _la int

	p.SetState(338)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(326)
			p.Match(ospreyParserMATCH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(327)
			p.Expr()
		}
		{
			p.SetState(328)
			p.Match(ospreyParserLBRACE)
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

		for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17122834203439112) != 0) {
			{
				p.SetState(329)
				p.MatchArm()
			}

			p.SetState(332)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(334)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(336)
			p.SelectExpr()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(337)
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

// ILoopExprContext is an interface to support dynamic dispatch.
type ILoopExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LOOP() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	BlockBody() IBlockBodyContext
	RBRACE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	Expr() IExprContext
	RPAREN() antlr.TerminalNode

	// IsLoopExprContext differentiates from other interfaces.
	IsLoopExprContext()
}

type LoopExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLoopExprContext() *LoopExprContext {
	var p = new(LoopExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_loopExpr
	return p
}

func InitEmptyLoopExprContext(p *LoopExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = ospreyParserRULE_loopExpr
}

func (*LoopExprContext) IsLoopExprContext() {}

func NewLoopExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LoopExprContext {
	var p = new(LoopExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = ospreyParserRULE_loopExpr

	return p
}

func (s *LoopExprContext) GetParser() antlr.Parser { return s.parser }

func (s *LoopExprContext) LOOP() antlr.TerminalNode {
	return s.GetToken(ospreyParserLOOP, 0)
}

func (s *LoopExprContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserLBRACE, 0)
}

func (s *LoopExprContext) BlockBody() IBlockBodyContext {
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

func (s *LoopExprContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(ospreyParserRBRACE, 0)
}

func (s *LoopExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserLPAREN, 0)
}

func (s *LoopExprContext) Expr() IExprContext {
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

func (s *LoopExprContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(ospreyParserRPAREN, 0)
}

func (s *LoopExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LoopExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LoopExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.EnterLoopExpr(s)
	}
}

func (s *LoopExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(ospreyListener); ok {
		listenerT.ExitLoopExpr(s)
	}
}

func (p *ospreyParser) LoopExpr() (localctx ILoopExprContext) {
	localctx = NewLoopExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, ospreyParserRULE_loopExpr)
	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(337)
			p.Match(ospreyParserLOOP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(338)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(339)
			p.BlockBody()
		}
		{
			p.SetState(340)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(342)
			p.Match(ospreyParserLOOP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(343)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(344)
			p.Expr()
		}
		{
			p.SetState(345)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(346)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(347)
			p.BlockBody()
		}
		{
			p.SetState(348)
			p.Match(ospreyParserRBRACE)
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
	p.EnterRule(localctx, 56, ospreyParserRULE_selectExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(340)
		p.Match(ospreyParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(341)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(343)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17122834203439112) != 0) {
		{
			p.SetState(342)
			p.SelectArm()
		}

		p.SetState(345)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(347)
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
	p.EnterRule(localctx, 56, ospreyParserRULE_selectArm)
	p.SetState(356)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 31, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(349)
			p.Pattern()
		}
		{
			p.SetState(350)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(351)
			p.Expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(353)
			p.Match(ospreyParserUNDERSCORE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(354)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(355)
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
	p.EnterRule(localctx, 60, ospreyParserRULE_binaryExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(358)
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
	p.EnterRule(localctx, 62, ospreyParserRULE_comparisonExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(360)
		p.AddExpr()
	}
	p.SetState(365)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(361)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&826646986752) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(362)
				p.AddExpr()
			}

		}
		p.SetState(367)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 64, ospreyParserRULE_addExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(368)
		p.MulExpr()
	}
	p.SetState(373)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(369)
				_la = p.GetTokenStream().LA(1)

				if !(_la == ospreyParserPLUS || _la == ospreyParserMINUS) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(370)
				p.MulExpr()
			}

		}
		p.SetState(375)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 66, ospreyParserRULE_mulExpr)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(376)
		p.UnaryExpr()
	}
	p.SetState(381)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(377)
				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&844429225099264) != 0) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(378)
				p.UnaryExpr()
			}

		}
		p.SetState(383)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 68, ospreyParserRULE_unaryExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(385)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 35, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(384)
			_la = p.GetTokenStream().LA(1)

			if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&211108380033024) != 0) {
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
		p.SetState(387)
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
	p.EnterRule(localctx, 70, ospreyParserRULE_pipeExpr)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(389)
		p.CallExpr()
	}
	p.SetState(394)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(390)
				p.Match(ospreyParserPIPE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(391)
				p.CallExpr()
			}

		}
		p.SetState(396)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 72, ospreyParserRULE_callExpr)
	var _la int

	var _alt int

	p.SetState(431)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(397)
			p.Primary()
		}
		p.SetState(400)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(398)
					p.Match(ospreyParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(399)
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

			p.SetState(402)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		p.SetState(409)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 39, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(404)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(406)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17122834169884684) != 0 {
				{
					p.SetState(405)
					p.ArgList()
				}

			}
			{
				p.SetState(408)
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
			p.SetState(411)
			p.Primary()
		}
		p.SetState(419)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(412)
					p.Match(ospreyParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(413)
					p.Match(ospreyParserID)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				{
					p.SetState(414)
					p.Match(ospreyParserLPAREN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(416)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17122834169884684) != 0 {
					{
						p.SetState(415)
						p.ArgList()
					}

				}
				{
					p.SetState(418)
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

			p.SetState(421)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(423)
			p.Primary()
		}
		p.SetState(429)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(424)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(426)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17122834169884684) != 0 {
				{
					p.SetState(425)
					p.ArgList()
				}

			}
			{
				p.SetState(428)
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
	p.EnterRule(localctx, 74, ospreyParserRULE_argList)
	var _la int

	p.SetState(442)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(433)
			p.NamedArgList()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(434)
			p.Expr()
		}
		p.SetState(439)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == ospreyParserCOMMA {
			{
				p.SetState(435)
				p.Match(ospreyParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(436)
				p.Expr()
			}

			p.SetState(441)
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
	p.EnterRule(localctx, 76, ospreyParserRULE_namedArgList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(444)
		p.NamedArg()
	}
	p.SetState(447)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserCOMMA {
		{
			p.SetState(445)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(446)
			p.NamedArg()
		}

		p.SetState(449)
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
	p.EnterRule(localctx, 78, ospreyParserRULE_namedArg)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(451)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(452)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(453)
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
	p.EnterRule(localctx, 78, ospreyParserRULE_primary)
	p.SetState(494)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 49, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(455)
			p.Match(ospreyParserSPAWN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(456)
			p.Expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(457)
			p.Match(ospreyParserYIELD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(459)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(458)
				p.Expr()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(461)
			p.Match(ospreyParserAWAIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(462)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(463)
			p.Expr()
		}
		{
			p.SetState(464)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(466)
			p.Match(ospreyParserSEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(467)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(468)
			p.Expr()
		}
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
		{
			p.SetState(471)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(473)
			p.Match(ospreyParserRECV)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(474)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(475)
			p.Expr()
		}
		{
			p.SetState(476)
			p.Match(ospreyParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(478)
			p.Match(ospreyParserSELECT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(479)
			p.SelectExpr()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(480)
			p.TypeConstructor()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(481)
			p.UpdateExpr()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(482)
			p.BlockExpr()
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(483)
			p.Literal()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(484)
			p.LambdaExpr()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
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
			p.Match(ospreyParserLSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(487)
			p.Match(ospreyParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(488)
			p.Match(ospreyParserRSQUARE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(489)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(490)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(491)
			p.Expr()
		}
		{
			p.SetState(492)
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
	p.EnterRule(localctx, 82, ospreyParserRULE_typeConstructor)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(496)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(498)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserLT {
		{
			p.SetState(497)
			p.TypeArgs()
		}

	}
	{
		p.SetState(500)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(501)
		p.FieldAssignments()
	}
	{
		p.SetState(502)
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
	p.EnterRule(localctx, 84, ospreyParserRULE_typeArgs)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(504)
		p.Match(ospreyParserLT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(505)
		p.TypeList()
	}
	{
		p.SetState(506)
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
	p.EnterRule(localctx, 86, ospreyParserRULE_fieldAssignments)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(508)
		p.FieldAssignment()
	}
	p.SetState(513)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
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
			p.FieldAssignment()
		}

		p.SetState(515)
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
	p.EnterRule(localctx, 88, ospreyParserRULE_fieldAssignment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(516)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(517)
		p.Match(ospreyParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(518)
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
	p.EnterRule(localctx, 90, ospreyParserRULE_lambdaExpr)
	var _la int

	p.SetState(539)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserFN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(520)
			p.Match(ospreyParserFN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(521)
			p.Match(ospreyParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(523)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserID {
			{
				p.SetState(522)
				p.ParamList()
			}

		}
		{
			p.SetState(525)
			p.Match(ospreyParserRPAREN)
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

		if _la == ospreyParserARROW {
			{
				p.SetState(526)
				p.Match(ospreyParserARROW)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(527)
				p.Type_()
			}

		}
		{
			p.SetState(530)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(531)
			p.Expr()
		}

	case ospreyParserBAR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(532)
			p.Match(ospreyParserBAR)
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

		if _la == ospreyParserID {
			{
				p.SetState(533)
				p.ParamList()
			}

		}
		{
			p.SetState(536)
			p.Match(ospreyParserBAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(537)
			p.Match(ospreyParserLAMBDA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(538)
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
	p.EnterRule(localctx, 92, ospreyParserRULE_updateExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(541)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(542)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(543)
		p.FieldAssignments()
	}
	{
		p.SetState(544)
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
	p.EnterRule(localctx, 94, ospreyParserRULE_blockExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(546)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(547)
		p.BlockBody()
	}
	{
		p.SetState(548)
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
	p.EnterRule(localctx, 94, ospreyParserRULE_literal)
	p.SetState(556)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case ospreyParserINT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(550)
			p.Match(ospreyParserINT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserSTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(551)
			p.Match(ospreyParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserINTERPOLATED_STRING:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(552)
			p.Match(ospreyParserINTERPOLATED_STRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserTRUE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(553)
			p.Match(ospreyParserTRUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserFALSE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(554)
			p.Match(ospreyParserFALSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case ospreyParserLSQUARE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(555)
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
	p.EnterRule(localctx, 96, ospreyParserRULE_listLiteral)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(558)
		p.Match(ospreyParserLSQUARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(567)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17122834169884684) != 0 {
		{
			p.SetState(559)
			p.Expr()
		}
		p.SetState(564)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == ospreyParserCOMMA {
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

			p.SetState(566)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(569)
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
	p.EnterRule(localctx, 98, ospreyParserRULE_docComment)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(572)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(571)
			p.Match(ospreyParserDOC_COMMENT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(574)
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
	p.EnterRule(localctx, 100, ospreyParserRULE_moduleDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(577)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == ospreyParserDOC_COMMENT {
		{
			p.SetState(576)
			p.DocComment()
		}

	}
	{
		p.SetState(579)
		p.Match(ospreyParserMODULE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(580)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(581)
		p.Match(ospreyParserLBRACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(582)
		p.ModuleBody()
	}
	{
		p.SetState(583)
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
	p.EnterRule(localctx, 102, ospreyParserRULE_moduleBody)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(588)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&36028797018964808) != 0 {
		{
			p.SetState(585)
			p.ModuleStatement()
		}

		p.SetState(590)
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
	p.EnterRule(localctx, 104, ospreyParserRULE_moduleStatement)
	p.SetState(594)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 62, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(591)
			p.LetDecl()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(592)
			p.FnDecl()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(593)
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
	p.EnterRule(localctx, 106, ospreyParserRULE_matchArm)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(596)
		p.Pattern()
	}
	{
		p.SetState(597)
		p.Match(ospreyParserLAMBDA)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(598)
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
	p.EnterRule(localctx, 108, ospreyParserRULE_pattern)
	var _la int

	p.SetState(640)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 67, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(600)
			p.UnaryExpr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(601)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(606)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLBRACE {
			{
				p.SetState(602)
				p.Match(ospreyParserLBRACE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(603)
				p.FieldPattern()
			}
			{
				p.SetState(604)
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
			p.SetState(608)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(620)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == ospreyParserLPAREN {
			{
				p.SetState(609)
				p.Match(ospreyParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(610)
				p.Pattern()
			}
			p.SetState(615)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == ospreyParserCOMMA {
				{
					p.SetState(611)
					p.Match(ospreyParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(612)
					p.Pattern()
				}

				p.SetState(617)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(618)
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
			p.SetState(622)
			p.Match(ospreyParserID)
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
			p.SetState(626)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(627)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(628)
			p.Type_()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(629)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(630)
			p.Match(ospreyParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(631)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(632)
			p.FieldPattern()
		}
		{
			p.SetState(633)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(635)
			p.Match(ospreyParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(636)
			p.FieldPattern()
		}
		{
			p.SetState(637)
			p.Match(ospreyParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(639)
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
	p.EnterRule(localctx, 110, ospreyParserRULE_fieldPattern)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(642)
		p.Match(ospreyParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(647)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == ospreyParserCOMMA {
		{
			p.SetState(643)
			p.Match(ospreyParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(644)
			p.Match(ospreyParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(649)
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
	p.EnterRule(localctx, 112, ospreyParserRULE_blockBody)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(653)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 69, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(650)
				p.Statement()
			}

		}
		p.SetState(655)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 69, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(657)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&17122834169884684) != 0 {
		{
			p.SetState(656)
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
