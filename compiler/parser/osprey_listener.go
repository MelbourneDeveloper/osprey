// Code generated from osprey.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // osprey
import "github.com/antlr4-go/antlr/v4"

// ospreyListener is a complete listener for a parse tree produced by ospreyParser.
type ospreyListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterImportStmt is called when entering the importStmt production.
	EnterImportStmt(c *ImportStmtContext)

	// EnterLetDecl is called when entering the letDecl production.
	EnterLetDecl(c *LetDeclContext)

	// EnterAssignStmt is called when entering the assignStmt production.
	EnterAssignStmt(c *AssignStmtContext)

	// EnterFnDecl is called when entering the fnDecl production.
	EnterFnDecl(c *FnDeclContext)

	// EnterExternDecl is called when entering the externDecl production.
	EnterExternDecl(c *ExternDeclContext)

	// EnterExternParamList is called when entering the externParamList production.
	EnterExternParamList(c *ExternParamListContext)

	// EnterExternParam is called when entering the externParam production.
	EnterExternParam(c *ExternParamContext)

	// EnterParamList is called when entering the paramList production.
	EnterParamList(c *ParamListContext)

	// EnterParam is called when entering the param production.
	EnterParam(c *ParamContext)

	// EnterTypeDecl is called when entering the typeDecl production.
	EnterTypeDecl(c *TypeDeclContext)

	// EnterTypeParamList is called when entering the typeParamList production.
	EnterTypeParamList(c *TypeParamListContext)

	// EnterUnionType is called when entering the unionType production.
	EnterUnionType(c *UnionTypeContext)

	// EnterRecordType is called when entering the recordType production.
	EnterRecordType(c *RecordTypeContext)

	// EnterVariant is called when entering the variant production.
	EnterVariant(c *VariantContext)

	// EnterFieldDeclarations is called when entering the fieldDeclarations production.
	EnterFieldDeclarations(c *FieldDeclarationsContext)

	// EnterFieldDeclaration is called when entering the fieldDeclaration production.
	EnterFieldDeclaration(c *FieldDeclarationContext)

	// EnterTypeValidation is called when entering the typeValidation production.
	EnterTypeValidation(c *TypeValidationContext)

	// EnterEffectDecl is called when entering the effectDecl production.
	EnterEffectDecl(c *EffectDeclContext)

	// EnterOpDecl is called when entering the opDecl production.
	EnterOpDecl(c *OpDeclContext)

	// EnterEffectSet is called when entering the effectSet production.
	EnterEffectSet(c *EffectSetContext)

	// EnterEffectList is called when entering the effectList production.
	EnterEffectList(c *EffectListContext)

	// EnterHandlerExpr is called when entering the handlerExpr production.
	EnterHandlerExpr(c *HandlerExprContext)

	// EnterHandlerArm is called when entering the handlerArm production.
	EnterHandlerArm(c *HandlerArmContext)

	// EnterHandlerParams is called when entering the handlerParams production.
	EnterHandlerParams(c *HandlerParamsContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterBooleanExpr is called when entering the booleanExpr production.
	EnterBooleanExpr(c *BooleanExprContext)

	// EnterFieldList is called when entering the fieldList production.
	EnterFieldList(c *FieldListContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterType is called when entering the type production.
	EnterType(c *TypeContext)

	// EnterTypeList is called when entering the typeList production.
	EnterTypeList(c *TypeListContext)

	// EnterExprStmt is called when entering the exprStmt production.
	EnterExprStmt(c *ExprStmtContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// EnterMatchExpr is called when entering the matchExpr production.
	EnterMatchExpr(c *MatchExprContext)

	// EnterSelectExpr is called when entering the selectExpr production.
	EnterSelectExpr(c *SelectExprContext)

	// EnterSelectArm is called when entering the selectArm production.
	EnterSelectArm(c *SelectArmContext)

	// EnterBinaryExpr is called when entering the binaryExpr production.
	EnterBinaryExpr(c *BinaryExprContext)

	// EnterTernaryExpr is called when entering the ternaryExpr production.
	EnterTernaryExpr(c *TernaryExprContext)

	// EnterComparisonExpr is called when entering the comparisonExpr production.
	EnterComparisonExpr(c *ComparisonExprContext)

	// EnterAddExpr is called when entering the addExpr production.
	EnterAddExpr(c *AddExprContext)

	// EnterMulExpr is called when entering the mulExpr production.
	EnterMulExpr(c *MulExprContext)

	// EnterUnaryExpr is called when entering the unaryExpr production.
	EnterUnaryExpr(c *UnaryExprContext)

	// EnterPipeExpr is called when entering the pipeExpr production.
	EnterPipeExpr(c *PipeExprContext)

	// EnterCallExpr is called when entering the callExpr production.
	EnterCallExpr(c *CallExprContext)

	// EnterArgList is called when entering the argList production.
	EnterArgList(c *ArgListContext)

	// EnterNamedArgList is called when entering the namedArgList production.
	EnterNamedArgList(c *NamedArgListContext)

	// EnterNamedArg is called when entering the namedArg production.
	EnterNamedArg(c *NamedArgContext)

	// EnterPrimary is called when entering the primary production.
	EnterPrimary(c *PrimaryContext)

	// EnterObjectLiteral is called when entering the objectLiteral production.
	EnterObjectLiteral(c *ObjectLiteralContext)

	// EnterTypeConstructor is called when entering the typeConstructor production.
	EnterTypeConstructor(c *TypeConstructorContext)

	// EnterTypeArgs is called when entering the typeArgs production.
	EnterTypeArgs(c *TypeArgsContext)

	// EnterFieldAssignments is called when entering the fieldAssignments production.
	EnterFieldAssignments(c *FieldAssignmentsContext)

	// EnterFieldAssignment is called when entering the fieldAssignment production.
	EnterFieldAssignment(c *FieldAssignmentContext)

	// EnterLambdaExpr is called when entering the lambdaExpr production.
	EnterLambdaExpr(c *LambdaExprContext)

	// EnterUpdateExpr is called when entering the updateExpr production.
	EnterUpdateExpr(c *UpdateExprContext)

	// EnterBlockExpr is called when entering the blockExpr production.
	EnterBlockExpr(c *BlockExprContext)

	// EnterLiteral is called when entering the literal production.
	EnterLiteral(c *LiteralContext)

	// EnterListLiteral is called when entering the listLiteral production.
	EnterListLiteral(c *ListLiteralContext)

	// EnterDocComment is called when entering the docComment production.
	EnterDocComment(c *DocCommentContext)

	// EnterModuleDecl is called when entering the moduleDecl production.
	EnterModuleDecl(c *ModuleDeclContext)

	// EnterModuleBody is called when entering the moduleBody production.
	EnterModuleBody(c *ModuleBodyContext)

	// EnterModuleStatement is called when entering the moduleStatement production.
	EnterModuleStatement(c *ModuleStatementContext)

	// EnterMatchArm is called when entering the matchArm production.
	EnterMatchArm(c *MatchArmContext)

	// EnterPattern is called when entering the pattern production.
	EnterPattern(c *PatternContext)

	// EnterFieldPattern is called when entering the fieldPattern production.
	EnterFieldPattern(c *FieldPatternContext)

	// EnterBlockBody is called when entering the blockBody production.
	EnterBlockBody(c *BlockBodyContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitImportStmt is called when exiting the importStmt production.
	ExitImportStmt(c *ImportStmtContext)

	// ExitLetDecl is called when exiting the letDecl production.
	ExitLetDecl(c *LetDeclContext)

	// ExitAssignStmt is called when exiting the assignStmt production.
	ExitAssignStmt(c *AssignStmtContext)

	// ExitFnDecl is called when exiting the fnDecl production.
	ExitFnDecl(c *FnDeclContext)

	// ExitExternDecl is called when exiting the externDecl production.
	ExitExternDecl(c *ExternDeclContext)

	// ExitExternParamList is called when exiting the externParamList production.
	ExitExternParamList(c *ExternParamListContext)

	// ExitExternParam is called when exiting the externParam production.
	ExitExternParam(c *ExternParamContext)

	// ExitParamList is called when exiting the paramList production.
	ExitParamList(c *ParamListContext)

	// ExitParam is called when exiting the param production.
	ExitParam(c *ParamContext)

	// ExitTypeDecl is called when exiting the typeDecl production.
	ExitTypeDecl(c *TypeDeclContext)

	// ExitTypeParamList is called when exiting the typeParamList production.
	ExitTypeParamList(c *TypeParamListContext)

	// ExitUnionType is called when exiting the unionType production.
	ExitUnionType(c *UnionTypeContext)

	// ExitRecordType is called when exiting the recordType production.
	ExitRecordType(c *RecordTypeContext)

	// ExitVariant is called when exiting the variant production.
	ExitVariant(c *VariantContext)

	// ExitFieldDeclarations is called when exiting the fieldDeclarations production.
	ExitFieldDeclarations(c *FieldDeclarationsContext)

	// ExitFieldDeclaration is called when exiting the fieldDeclaration production.
	ExitFieldDeclaration(c *FieldDeclarationContext)

	// ExitTypeValidation is called when exiting the typeValidation production.
	ExitTypeValidation(c *TypeValidationContext)

	// ExitEffectDecl is called when exiting the effectDecl production.
	ExitEffectDecl(c *EffectDeclContext)

	// ExitOpDecl is called when exiting the opDecl production.
	ExitOpDecl(c *OpDeclContext)

	// ExitEffectSet is called when exiting the effectSet production.
	ExitEffectSet(c *EffectSetContext)

	// ExitEffectList is called when exiting the effectList production.
	ExitEffectList(c *EffectListContext)

	// ExitHandlerExpr is called when exiting the handlerExpr production.
	ExitHandlerExpr(c *HandlerExprContext)

	// ExitHandlerArm is called when exiting the handlerArm production.
	ExitHandlerArm(c *HandlerArmContext)

	// ExitHandlerParams is called when exiting the handlerParams production.
	ExitHandlerParams(c *HandlerParamsContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitBooleanExpr is called when exiting the booleanExpr production.
	ExitBooleanExpr(c *BooleanExprContext)

	// ExitFieldList is called when exiting the fieldList production.
	ExitFieldList(c *FieldListContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitType is called when exiting the type production.
	ExitType(c *TypeContext)

	// ExitTypeList is called when exiting the typeList production.
	ExitTypeList(c *TypeListContext)

	// ExitExprStmt is called when exiting the exprStmt production.
	ExitExprStmt(c *ExprStmtContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

	// ExitMatchExpr is called when exiting the matchExpr production.
	ExitMatchExpr(c *MatchExprContext)

	// ExitSelectExpr is called when exiting the selectExpr production.
	ExitSelectExpr(c *SelectExprContext)

	// ExitSelectArm is called when exiting the selectArm production.
	ExitSelectArm(c *SelectArmContext)

	// ExitBinaryExpr is called when exiting the binaryExpr production.
	ExitBinaryExpr(c *BinaryExprContext)

	// ExitTernaryExpr is called when exiting the ternaryExpr production.
	ExitTernaryExpr(c *TernaryExprContext)

	// ExitComparisonExpr is called when exiting the comparisonExpr production.
	ExitComparisonExpr(c *ComparisonExprContext)

	// ExitAddExpr is called when exiting the addExpr production.
	ExitAddExpr(c *AddExprContext)

	// ExitMulExpr is called when exiting the mulExpr production.
	ExitMulExpr(c *MulExprContext)

	// ExitUnaryExpr is called when exiting the unaryExpr production.
	ExitUnaryExpr(c *UnaryExprContext)

	// ExitPipeExpr is called when exiting the pipeExpr production.
	ExitPipeExpr(c *PipeExprContext)

	// ExitCallExpr is called when exiting the callExpr production.
	ExitCallExpr(c *CallExprContext)

	// ExitArgList is called when exiting the argList production.
	ExitArgList(c *ArgListContext)

	// ExitNamedArgList is called when exiting the namedArgList production.
	ExitNamedArgList(c *NamedArgListContext)

	// ExitNamedArg is called when exiting the namedArg production.
	ExitNamedArg(c *NamedArgContext)

	// ExitPrimary is called when exiting the primary production.
	ExitPrimary(c *PrimaryContext)

	// ExitObjectLiteral is called when exiting the objectLiteral production.
	ExitObjectLiteral(c *ObjectLiteralContext)

	// ExitTypeConstructor is called when exiting the typeConstructor production.
	ExitTypeConstructor(c *TypeConstructorContext)

	// ExitTypeArgs is called when exiting the typeArgs production.
	ExitTypeArgs(c *TypeArgsContext)

	// ExitFieldAssignments is called when exiting the fieldAssignments production.
	ExitFieldAssignments(c *FieldAssignmentsContext)

	// ExitFieldAssignment is called when exiting the fieldAssignment production.
	ExitFieldAssignment(c *FieldAssignmentContext)

	// ExitLambdaExpr is called when exiting the lambdaExpr production.
	ExitLambdaExpr(c *LambdaExprContext)

	// ExitUpdateExpr is called when exiting the updateExpr production.
	ExitUpdateExpr(c *UpdateExprContext)

	// ExitBlockExpr is called when exiting the blockExpr production.
	ExitBlockExpr(c *BlockExprContext)

	// ExitLiteral is called when exiting the literal production.
	ExitLiteral(c *LiteralContext)

	// ExitListLiteral is called when exiting the listLiteral production.
	ExitListLiteral(c *ListLiteralContext)

	// ExitDocComment is called when exiting the docComment production.
	ExitDocComment(c *DocCommentContext)

	// ExitModuleDecl is called when exiting the moduleDecl production.
	ExitModuleDecl(c *ModuleDeclContext)

	// ExitModuleBody is called when exiting the moduleBody production.
	ExitModuleBody(c *ModuleBodyContext)

	// ExitModuleStatement is called when exiting the moduleStatement production.
	ExitModuleStatement(c *ModuleStatementContext)

	// ExitMatchArm is called when exiting the matchArm production.
	ExitMatchArm(c *MatchArmContext)

	// ExitPattern is called when exiting the pattern production.
	ExitPattern(c *PatternContext)

	// ExitFieldPattern is called when exiting the fieldPattern production.
	ExitFieldPattern(c *FieldPatternContext)

	// ExitBlockBody is called when exiting the blockBody production.
	ExitBlockBody(c *BlockBodyContext)
}
