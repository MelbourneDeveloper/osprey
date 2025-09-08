// Code generated from osprey.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // osprey
import "github.com/antlr4-go/antlr/v4"

// BaseospreyListener is a complete listener for a parse tree produced by ospreyParser.
type BaseospreyListener struct{}

var _ ospreyListener = &BaseospreyListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseospreyListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseospreyListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseospreyListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseospreyListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BaseospreyListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BaseospreyListener) ExitProgram(ctx *ProgramContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseospreyListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseospreyListener) ExitStatement(ctx *StatementContext) {}

// EnterImportStmt is called when production importStmt is entered.
func (s *BaseospreyListener) EnterImportStmt(ctx *ImportStmtContext) {}

// ExitImportStmt is called when production importStmt is exited.
func (s *BaseospreyListener) ExitImportStmt(ctx *ImportStmtContext) {}

// EnterLetDecl is called when production letDecl is entered.
func (s *BaseospreyListener) EnterLetDecl(ctx *LetDeclContext) {}

// ExitLetDecl is called when production letDecl is exited.
func (s *BaseospreyListener) ExitLetDecl(ctx *LetDeclContext) {}

// EnterAssignStmt is called when production assignStmt is entered.
func (s *BaseospreyListener) EnterAssignStmt(ctx *AssignStmtContext) {}

// ExitAssignStmt is called when production assignStmt is exited.
func (s *BaseospreyListener) ExitAssignStmt(ctx *AssignStmtContext) {}

// EnterFnDecl is called when production fnDecl is entered.
func (s *BaseospreyListener) EnterFnDecl(ctx *FnDeclContext) {}

// ExitFnDecl is called when production fnDecl is exited.
func (s *BaseospreyListener) ExitFnDecl(ctx *FnDeclContext) {}

// EnterExternDecl is called when production externDecl is entered.
func (s *BaseospreyListener) EnterExternDecl(ctx *ExternDeclContext) {}

// ExitExternDecl is called when production externDecl is exited.
func (s *BaseospreyListener) ExitExternDecl(ctx *ExternDeclContext) {}

// EnterExternParamList is called when production externParamList is entered.
func (s *BaseospreyListener) EnterExternParamList(ctx *ExternParamListContext) {}

// ExitExternParamList is called when production externParamList is exited.
func (s *BaseospreyListener) ExitExternParamList(ctx *ExternParamListContext) {}

// EnterExternParam is called when production externParam is entered.
func (s *BaseospreyListener) EnterExternParam(ctx *ExternParamContext) {}

// ExitExternParam is called when production externParam is exited.
func (s *BaseospreyListener) ExitExternParam(ctx *ExternParamContext) {}

// EnterParamList is called when production paramList is entered.
func (s *BaseospreyListener) EnterParamList(ctx *ParamListContext) {}

// ExitParamList is called when production paramList is exited.
func (s *BaseospreyListener) ExitParamList(ctx *ParamListContext) {}

// EnterParam is called when production param is entered.
func (s *BaseospreyListener) EnterParam(ctx *ParamContext) {}

// ExitParam is called when production param is exited.
func (s *BaseospreyListener) ExitParam(ctx *ParamContext) {}

// EnterTypeDecl is called when production typeDecl is entered.
func (s *BaseospreyListener) EnterTypeDecl(ctx *TypeDeclContext) {}

// ExitTypeDecl is called when production typeDecl is exited.
func (s *BaseospreyListener) ExitTypeDecl(ctx *TypeDeclContext) {}

// EnterTypeParamList is called when production typeParamList is entered.
func (s *BaseospreyListener) EnterTypeParamList(ctx *TypeParamListContext) {}

// ExitTypeParamList is called when production typeParamList is exited.
func (s *BaseospreyListener) ExitTypeParamList(ctx *TypeParamListContext) {}

// EnterUnionType is called when production unionType is entered.
func (s *BaseospreyListener) EnterUnionType(ctx *UnionTypeContext) {}

// ExitUnionType is called when production unionType is exited.
func (s *BaseospreyListener) ExitUnionType(ctx *UnionTypeContext) {}

// EnterRecordType is called when production recordType is entered.
func (s *BaseospreyListener) EnterRecordType(ctx *RecordTypeContext) {}

// ExitRecordType is called when production recordType is exited.
func (s *BaseospreyListener) ExitRecordType(ctx *RecordTypeContext) {}

// EnterVariant is called when production variant is entered.
func (s *BaseospreyListener) EnterVariant(ctx *VariantContext) {}

// ExitVariant is called when production variant is exited.
func (s *BaseospreyListener) ExitVariant(ctx *VariantContext) {}

// EnterFieldDeclarations is called when production fieldDeclarations is entered.
func (s *BaseospreyListener) EnterFieldDeclarations(ctx *FieldDeclarationsContext) {}

// ExitFieldDeclarations is called when production fieldDeclarations is exited.
func (s *BaseospreyListener) ExitFieldDeclarations(ctx *FieldDeclarationsContext) {}

// EnterFieldDeclaration is called when production fieldDeclaration is entered.
func (s *BaseospreyListener) EnterFieldDeclaration(ctx *FieldDeclarationContext) {}

// ExitFieldDeclaration is called when production fieldDeclaration is exited.
func (s *BaseospreyListener) ExitFieldDeclaration(ctx *FieldDeclarationContext) {}

// EnterTypeValidation is called when production typeValidation is entered.
func (s *BaseospreyListener) EnterTypeValidation(ctx *TypeValidationContext) {}

// ExitTypeValidation is called when production typeValidation is exited.
func (s *BaseospreyListener) ExitTypeValidation(ctx *TypeValidationContext) {}

// EnterEffectDecl is called when production effectDecl is entered.
func (s *BaseospreyListener) EnterEffectDecl(ctx *EffectDeclContext) {}

// ExitEffectDecl is called when production effectDecl is exited.
func (s *BaseospreyListener) ExitEffectDecl(ctx *EffectDeclContext) {}

// EnterOpDecl is called when production opDecl is entered.
func (s *BaseospreyListener) EnterOpDecl(ctx *OpDeclContext) {}

// ExitOpDecl is called when production opDecl is exited.
func (s *BaseospreyListener) ExitOpDecl(ctx *OpDeclContext) {}

// EnterEffectSet is called when production effectSet is entered.
func (s *BaseospreyListener) EnterEffectSet(ctx *EffectSetContext) {}

// ExitEffectSet is called when production effectSet is exited.
func (s *BaseospreyListener) ExitEffectSet(ctx *EffectSetContext) {}

// EnterEffectList is called when production effectList is entered.
func (s *BaseospreyListener) EnterEffectList(ctx *EffectListContext) {}

// ExitEffectList is called when production effectList is exited.
func (s *BaseospreyListener) ExitEffectList(ctx *EffectListContext) {}

// EnterHandlerExpr is called when production handlerExpr is entered.
func (s *BaseospreyListener) EnterHandlerExpr(ctx *HandlerExprContext) {}

// ExitHandlerExpr is called when production handlerExpr is exited.
func (s *BaseospreyListener) ExitHandlerExpr(ctx *HandlerExprContext) {}

// EnterHandlerArm is called when production handlerArm is entered.
func (s *BaseospreyListener) EnterHandlerArm(ctx *HandlerArmContext) {}

// ExitHandlerArm is called when production handlerArm is exited.
func (s *BaseospreyListener) ExitHandlerArm(ctx *HandlerArmContext) {}

// EnterHandlerParams is called when production handlerParams is entered.
func (s *BaseospreyListener) EnterHandlerParams(ctx *HandlerParamsContext) {}

// ExitHandlerParams is called when production handlerParams is exited.
func (s *BaseospreyListener) ExitHandlerParams(ctx *HandlerParamsContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseospreyListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseospreyListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterBooleanExpr is called when production booleanExpr is entered.
func (s *BaseospreyListener) EnterBooleanExpr(ctx *BooleanExprContext) {}

// ExitBooleanExpr is called when production booleanExpr is exited.
func (s *BaseospreyListener) ExitBooleanExpr(ctx *BooleanExprContext) {}

// EnterFieldList is called when production fieldList is entered.
func (s *BaseospreyListener) EnterFieldList(ctx *FieldListContext) {}

// ExitFieldList is called when production fieldList is exited.
func (s *BaseospreyListener) ExitFieldList(ctx *FieldListContext) {}

// EnterField is called when production field is entered.
func (s *BaseospreyListener) EnterField(ctx *FieldContext) {}

// ExitField is called when production field is exited.
func (s *BaseospreyListener) ExitField(ctx *FieldContext) {}

// EnterType is called when production type is entered.
func (s *BaseospreyListener) EnterType(ctx *TypeContext) {}

// ExitType is called when production type is exited.
func (s *BaseospreyListener) ExitType(ctx *TypeContext) {}

// EnterTypeList is called when production typeList is entered.
func (s *BaseospreyListener) EnterTypeList(ctx *TypeListContext) {}

// ExitTypeList is called when production typeList is exited.
func (s *BaseospreyListener) ExitTypeList(ctx *TypeListContext) {}

// EnterExprStmt is called when production exprStmt is entered.
func (s *BaseospreyListener) EnterExprStmt(ctx *ExprStmtContext) {}

// ExitExprStmt is called when production exprStmt is exited.
func (s *BaseospreyListener) ExitExprStmt(ctx *ExprStmtContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseospreyListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseospreyListener) ExitExpr(ctx *ExprContext) {}

// EnterMatchExpr is called when production matchExpr is entered.
func (s *BaseospreyListener) EnterMatchExpr(ctx *MatchExprContext) {}

// ExitMatchExpr is called when production matchExpr is exited.
func (s *BaseospreyListener) ExitMatchExpr(ctx *MatchExprContext) {}

// EnterSelectExpr is called when production selectExpr is entered.
func (s *BaseospreyListener) EnterSelectExpr(ctx *SelectExprContext) {}

// ExitSelectExpr is called when production selectExpr is exited.
func (s *BaseospreyListener) ExitSelectExpr(ctx *SelectExprContext) {}

// EnterSelectArm is called when production selectArm is entered.
func (s *BaseospreyListener) EnterSelectArm(ctx *SelectArmContext) {}

// ExitSelectArm is called when production selectArm is exited.
func (s *BaseospreyListener) ExitSelectArm(ctx *SelectArmContext) {}

// EnterBinaryExpr is called when production binaryExpr is entered.
func (s *BaseospreyListener) EnterBinaryExpr(ctx *BinaryExprContext) {}

// ExitBinaryExpr is called when production binaryExpr is exited.
func (s *BaseospreyListener) ExitBinaryExpr(ctx *BinaryExprContext) {}

// EnterTernaryExpr is called when production ternaryExpr is entered.
func (s *BaseospreyListener) EnterTernaryExpr(ctx *TernaryExprContext) {}

// ExitTernaryExpr is called when production ternaryExpr is exited.
func (s *BaseospreyListener) ExitTernaryExpr(ctx *TernaryExprContext) {}

// EnterComparisonExpr is called when production comparisonExpr is entered.
func (s *BaseospreyListener) EnterComparisonExpr(ctx *ComparisonExprContext) {}

// ExitComparisonExpr is called when production comparisonExpr is exited.
func (s *BaseospreyListener) ExitComparisonExpr(ctx *ComparisonExprContext) {}

// EnterLogicalOrExpr is called when production logicalOrExpr is entered.
func (s *BaseospreyListener) EnterLogicalOrExpr(ctx *LogicalOrExprContext) {}

// ExitLogicalOrExpr is called when production logicalOrExpr is exited.
func (s *BaseospreyListener) ExitLogicalOrExpr(ctx *LogicalOrExprContext) {}

// EnterLogicalAndExpr is called when production logicalAndExpr is entered.
func (s *BaseospreyListener) EnterLogicalAndExpr(ctx *LogicalAndExprContext) {}

// ExitLogicalAndExpr is called when production logicalAndExpr is exited.
func (s *BaseospreyListener) ExitLogicalAndExpr(ctx *LogicalAndExprContext) {}

// EnterAddExpr is called when production addExpr is entered.
func (s *BaseospreyListener) EnterAddExpr(ctx *AddExprContext) {}

// ExitAddExpr is called when production addExpr is exited.
func (s *BaseospreyListener) ExitAddExpr(ctx *AddExprContext) {}

// EnterMulExpr is called when production mulExpr is entered.
func (s *BaseospreyListener) EnterMulExpr(ctx *MulExprContext) {}

// ExitMulExpr is called when production mulExpr is exited.
func (s *BaseospreyListener) ExitMulExpr(ctx *MulExprContext) {}

// EnterUnaryExpr is called when production unaryExpr is entered.
func (s *BaseospreyListener) EnterUnaryExpr(ctx *UnaryExprContext) {}

// ExitUnaryExpr is called when production unaryExpr is exited.
func (s *BaseospreyListener) ExitUnaryExpr(ctx *UnaryExprContext) {}

// EnterPipeExpr is called when production pipeExpr is entered.
func (s *BaseospreyListener) EnterPipeExpr(ctx *PipeExprContext) {}

// ExitPipeExpr is called when production pipeExpr is exited.
func (s *BaseospreyListener) ExitPipeExpr(ctx *PipeExprContext) {}

// EnterCallExpr is called when production callExpr is entered.
func (s *BaseospreyListener) EnterCallExpr(ctx *CallExprContext) {}

// ExitCallExpr is called when production callExpr is exited.
func (s *BaseospreyListener) ExitCallExpr(ctx *CallExprContext) {}

// EnterArgList is called when production argList is entered.
func (s *BaseospreyListener) EnterArgList(ctx *ArgListContext) {}

// ExitArgList is called when production argList is exited.
func (s *BaseospreyListener) ExitArgList(ctx *ArgListContext) {}

// EnterNamedArgList is called when production namedArgList is entered.
func (s *BaseospreyListener) EnterNamedArgList(ctx *NamedArgListContext) {}

// ExitNamedArgList is called when production namedArgList is exited.
func (s *BaseospreyListener) ExitNamedArgList(ctx *NamedArgListContext) {}

// EnterNamedArg is called when production namedArg is entered.
func (s *BaseospreyListener) EnterNamedArg(ctx *NamedArgContext) {}

// ExitNamedArg is called when production namedArg is exited.
func (s *BaseospreyListener) ExitNamedArg(ctx *NamedArgContext) {}

// EnterPrimary is called when production primary is entered.
func (s *BaseospreyListener) EnterPrimary(ctx *PrimaryContext) {}

// ExitPrimary is called when production primary is exited.
func (s *BaseospreyListener) ExitPrimary(ctx *PrimaryContext) {}

// EnterObjectLiteral is called when production objectLiteral is entered.
func (s *BaseospreyListener) EnterObjectLiteral(ctx *ObjectLiteralContext) {}

// ExitObjectLiteral is called when production objectLiteral is exited.
func (s *BaseospreyListener) ExitObjectLiteral(ctx *ObjectLiteralContext) {}

// EnterTypeConstructor is called when production typeConstructor is entered.
func (s *BaseospreyListener) EnterTypeConstructor(ctx *TypeConstructorContext) {}

// ExitTypeConstructor is called when production typeConstructor is exited.
func (s *BaseospreyListener) ExitTypeConstructor(ctx *TypeConstructorContext) {}

// EnterTypeArgs is called when production typeArgs is entered.
func (s *BaseospreyListener) EnterTypeArgs(ctx *TypeArgsContext) {}

// ExitTypeArgs is called when production typeArgs is exited.
func (s *BaseospreyListener) ExitTypeArgs(ctx *TypeArgsContext) {}

// EnterFieldAssignments is called when production fieldAssignments is entered.
func (s *BaseospreyListener) EnterFieldAssignments(ctx *FieldAssignmentsContext) {}

// ExitFieldAssignments is called when production fieldAssignments is exited.
func (s *BaseospreyListener) ExitFieldAssignments(ctx *FieldAssignmentsContext) {}

// EnterFieldAssignment is called when production fieldAssignment is entered.
func (s *BaseospreyListener) EnterFieldAssignment(ctx *FieldAssignmentContext) {}

// ExitFieldAssignment is called when production fieldAssignment is exited.
func (s *BaseospreyListener) ExitFieldAssignment(ctx *FieldAssignmentContext) {}

// EnterLambdaExpr is called when production lambdaExpr is entered.
func (s *BaseospreyListener) EnterLambdaExpr(ctx *LambdaExprContext) {}

// ExitLambdaExpr is called when production lambdaExpr is exited.
func (s *BaseospreyListener) ExitLambdaExpr(ctx *LambdaExprContext) {}

// EnterUpdateExpr is called when production updateExpr is entered.
func (s *BaseospreyListener) EnterUpdateExpr(ctx *UpdateExprContext) {}

// ExitUpdateExpr is called when production updateExpr is exited.
func (s *BaseospreyListener) ExitUpdateExpr(ctx *UpdateExprContext) {}

// EnterBlockExpr is called when production blockExpr is entered.
func (s *BaseospreyListener) EnterBlockExpr(ctx *BlockExprContext) {}

// ExitBlockExpr is called when production blockExpr is exited.
func (s *BaseospreyListener) ExitBlockExpr(ctx *BlockExprContext) {}

// EnterLiteral is called when production literal is entered.
func (s *BaseospreyListener) EnterLiteral(ctx *LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *BaseospreyListener) ExitLiteral(ctx *LiteralContext) {}

// EnterListLiteral is called when production listLiteral is entered.
func (s *BaseospreyListener) EnterListLiteral(ctx *ListLiteralContext) {}

// ExitListLiteral is called when production listLiteral is exited.
func (s *BaseospreyListener) ExitListLiteral(ctx *ListLiteralContext) {}

// EnterMapLiteral is called when production mapLiteral is entered.
func (s *BaseospreyListener) EnterMapLiteral(ctx *MapLiteralContext) {}

// ExitMapLiteral is called when production mapLiteral is exited.
func (s *BaseospreyListener) ExitMapLiteral(ctx *MapLiteralContext) {}

// EnterMapEntry is called when production mapEntry is entered.
func (s *BaseospreyListener) EnterMapEntry(ctx *MapEntryContext) {}

// ExitMapEntry is called when production mapEntry is exited.
func (s *BaseospreyListener) ExitMapEntry(ctx *MapEntryContext) {}

// EnterDocComment is called when production docComment is entered.
func (s *BaseospreyListener) EnterDocComment(ctx *DocCommentContext) {}

// ExitDocComment is called when production docComment is exited.
func (s *BaseospreyListener) ExitDocComment(ctx *DocCommentContext) {}

// EnterModuleDecl is called when production moduleDecl is entered.
func (s *BaseospreyListener) EnterModuleDecl(ctx *ModuleDeclContext) {}

// ExitModuleDecl is called when production moduleDecl is exited.
func (s *BaseospreyListener) ExitModuleDecl(ctx *ModuleDeclContext) {}

// EnterModuleBody is called when production moduleBody is entered.
func (s *BaseospreyListener) EnterModuleBody(ctx *ModuleBodyContext) {}

// ExitModuleBody is called when production moduleBody is exited.
func (s *BaseospreyListener) ExitModuleBody(ctx *ModuleBodyContext) {}

// EnterModuleStatement is called when production moduleStatement is entered.
func (s *BaseospreyListener) EnterModuleStatement(ctx *ModuleStatementContext) {}

// ExitModuleStatement is called when production moduleStatement is exited.
func (s *BaseospreyListener) ExitModuleStatement(ctx *ModuleStatementContext) {}

// EnterMatchArm is called when production matchArm is entered.
func (s *BaseospreyListener) EnterMatchArm(ctx *MatchArmContext) {}

// ExitMatchArm is called when production matchArm is exited.
func (s *BaseospreyListener) ExitMatchArm(ctx *MatchArmContext) {}

// EnterPattern is called when production pattern is entered.
func (s *BaseospreyListener) EnterPattern(ctx *PatternContext) {}

// ExitPattern is called when production pattern is exited.
func (s *BaseospreyListener) ExitPattern(ctx *PatternContext) {}

// EnterFieldPattern is called when production fieldPattern is entered.
func (s *BaseospreyListener) EnterFieldPattern(ctx *FieldPatternContext) {}

// ExitFieldPattern is called when production fieldPattern is exited.
func (s *BaseospreyListener) ExitFieldPattern(ctx *FieldPatternContext) {}

// EnterBlockBody is called when production blockBody is entered.
func (s *BaseospreyListener) EnterBlockBody(ctx *BlockBodyContext) {}

// ExitBlockBody is called when production blockBody is exited.
func (s *BaseospreyListener) ExitBlockBody(ctx *BlockBodyContext) {}
