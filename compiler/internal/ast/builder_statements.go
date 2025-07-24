package ast

import (
	"fmt"

	"github.com/christianfindlay/osprey/parser"
)

func (b *Builder) buildImport(ctx parser.IImportStmtContext) *ImportStatement {
	modules := make([]string, 0)
	for _, id := range ctx.AllID() {
		modules = append(modules, id.GetText())
	}

	return &ImportStatement{Module: modules}
}

func (b *Builder) buildLetDecl(ctx parser.ILetDeclContext) *LetDeclaration {
	name := ctx.ID().GetText()
	mutable := ctx.MUT() != nil
	value := b.buildExpression(ctx.Expr())

	// Parse type annotation if present
	var typeAnnotation *TypeExpression
	if ctx.Type_() != nil {
		typeAnnotation = b.buildTypeExpression(ctx.Type_())
	}

	return &LetDeclaration{
		Name:     name,
		Mutable:  mutable,
		Type:     typeAnnotation,
		Value:    value,
		Position: b.getPositionFromContext(ctx),
	}
}

func (b *Builder) buildAssignStmt(ctx parser.IAssignStmtContext) *AssignmentStatement {
	name := ctx.ID().GetText()
	value := b.buildExpression(ctx.Expr())

	return &AssignmentStatement{
		Name:     name,
		Value:    value,
		Position: b.getPositionFromContext(ctx),
	}
}

func (b *Builder) buildFnDecl(ctx parser.IFnDeclContext) *FunctionDeclaration {
	name := ctx.ID().GetText()

	params := make([]Parameter, 0)

	if ctx.ParamList() != nil {
		for _, paramCtx := range ctx.ParamList().AllParam() {
			param := Parameter{
				Name: paramCtx.ID().GetText(),
				Type: nil, // Parse type annotation if present
			}

			// Parse parameter type annotation if present
			if paramCtx.Type_() != nil {
				param.Type = b.buildTypeExpression(paramCtx.Type_())
			}

			params = append(params, param)
		}
	}

	// Handle both expression bodies (= expr) and block bodies ({ ... })
	var body Expression
	if ctx.Expr() != nil {
		// Expression-bodied function: fn name() = expr
		body = b.buildExpression(ctx.Expr())
	} else if ctx.BlockBody() != nil {
		// Block-bodied function: fn name() { statements }
		body = b.buildBlockBody(ctx.BlockBody())
	}

	// Parse return type annotation if present
	var returnType *TypeExpression
	if ctx.Type_() != nil {
		returnType = b.buildTypeExpression(ctx.Type_())
	}

	// Parse effect signatures (!Logger, ![IO, Net])
	var effects []string
	if ctx.EffectSet() != nil {
		effects = b.buildEffectSet(ctx.EffectSet())
	}

	return &FunctionDeclaration{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		Effects:    effects, // CRITICAL: Include parsed effects
		Body:       body,
		Position:   b.getPositionFromContext(ctx),
	}
}

func (b *Builder) buildExternDecl(ctx parser.IExternDeclContext) *ExternDeclaration {
	name := ctx.ID().GetText()

	params := make([]ExternParameter, 0)

	if ctx.ExternParamList() != nil {
		for _, paramCtx := range ctx.ExternParamList().AllExternParam() {
			param := ExternParameter{
				Name: paramCtx.ID().GetText(),
				Type: *b.buildTypeExpression(paramCtx.Type_()),
			}
			params = append(params, param)
		}
	}

	// Parse return type annotation if present
	var returnType *TypeExpression
	if ctx.Type_() != nil {
		returnType = b.buildTypeExpression(ctx.Type_())
	}

	return &ExternDeclaration{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
	}
}

func (b *Builder) buildTypeDecl(ctx parser.ITypeDeclContext) *TypeDeclaration {
	name := ctx.ID().GetText()

	// Handle generic type parameters
	typeParams := make([]string, 0)

	if ctx.TypeParamList() != nil {
		for _, id := range ctx.TypeParamList().AllID() {
			typeParams = append(typeParams, id.GetText())
		}
	}

	variants := make([]TypeVariant, 0)

	// Handle union types
	if ctx.UnionType() != nil {
		for _, variantCtx := range ctx.UnionType().AllVariant() {
			variant := b.buildVariant(variantCtx)
			variants = append(variants, variant)
		}
	}

	// Handle record types (which are essentially single variants)
	if ctx.RecordType() != nil {
		// Create a single variant with the type name and record fields
		fields := make([]TypeField, 0)
		if ctx.RecordType().FieldDeclarations() != nil {
			for _, fieldCtx := range ctx.RecordType().FieldDeclarations().AllFieldDeclaration() {
				field := TypeField{
					Name: fieldCtx.ID().GetText(),
					Type: fieldCtx.Type_().ID().GetText(),
				}

				// Parse WHERE constraint if present
				if fieldCtx.FunctionCall() != nil {
					constraint := &FunctionCallExpression{
						Function:  fieldCtx.FunctionCall().ID().GetText(),
						Arguments: make([]Expression, 0),
					}
					
					// Parse function call arguments if any
					if fieldCtx.FunctionCall().ArgList() != nil {
						for _, argCtx := range fieldCtx.FunctionCall().ArgList().AllExpr() {
							arg := b.buildExpression(argCtx)
							constraint.Arguments = append(constraint.Arguments, arg)
						}
					}
					
					field.Constraint = constraint
				}

				fields = append(fields, field)
			}
		}

		variant := TypeVariant{
			Name:   name, // Use the type name as the variant name for record types
			Fields: fields,
		}
		variants = append(variants, variant)
	}

	// Handle optional type validation
	var validationFunc *string
	if ctx.TypeValidation() != nil {
		funcName := ctx.TypeValidation().ID().GetText()
		validationFunc = &funcName
	}

	return &TypeDeclaration{
		Name:           name,
		TypeParams:     typeParams,
		Variants:       variants,
		ValidationFunc: validationFunc,
	}
}

func (b *Builder) buildVariant(ctx parser.IVariantContext) TypeVariant {
	name := ctx.ID().GetText()

	fields := make([]TypeField, 0)

	if ctx.FieldDeclarations() != nil {
		for _, fieldCtx := range ctx.FieldDeclarations().AllFieldDeclaration() {
			field := TypeField{
				Name: fieldCtx.ID().GetText(),
				Type: fieldCtx.Type_().ID().GetText(),
			}

			// Parse WHERE constraint if present
			if fieldCtx.FunctionCall() != nil {
				constraint := &FunctionCallExpression{
					Function:  fieldCtx.FunctionCall().ID().GetText(),
					Arguments: make([]Expression, 0),
				}
				
				// Parse function call arguments if any
				if fieldCtx.FunctionCall().ArgList() != nil {
					for _, argCtx := range fieldCtx.FunctionCall().ArgList().AllExpr() {
						arg := b.buildExpression(argCtx)
						constraint.Arguments = append(constraint.Arguments, arg)
					}
				}
				
				field.Constraint = constraint
			}

			fields = append(fields, field)
		}
	}

	return TypeVariant{
		Name:   name,
		Fields: fields,
	}
}

// buildTypeExpression builds a TypeExpression from a parser type context.
func (b *Builder) buildTypeExpression(ctx parser.ITypeContext) *TypeExpression {
	if ctx == nil {
		return nil
	}

	// Check if this is a function type: (Type, Type, ...) -> Type
	if ctx.LPAREN() != nil && ctx.ARROW() != nil {
		// This is a function type
		var paramTypes []TypeExpression

		// Parse parameter types from TypeList
		if ctx.TypeList() != nil {
			for _, typeCtx := range ctx.TypeList().AllType_() {
				paramType := b.buildTypeExpression(typeCtx)
				if paramType != nil {
					paramTypes = append(paramTypes, *paramType)
				}
			}
		}

		// Parse return type (the Type_ after the arrow)
		returnType := b.buildTypeExpression(ctx.Type_())

		return &TypeExpression{
			IsFunction:     true,
			ParameterTypes: paramTypes,
			ReturnType:     returnType,
		}
	}

	// Regular type (ID-based)
	if ctx.ID() == nil {
		// If no ID, this might be a nested function type or error case
		return nil
	}

	typeName := ctx.ID().GetText()
	typeExpr := &TypeExpression{
		Name: typeName,
	}

	// Handle generic types like Result<String, Error>
	if ctx.LT() != nil && ctx.GT() != nil && ctx.TypeList() != nil {
		for _, typeCtx := range ctx.TypeList().AllType_() {
			genericParam := b.buildTypeExpression(typeCtx)
			if genericParam != nil {
				typeExpr.GenericParams = append(typeExpr.GenericParams, *genericParam)
			}
		}
	}

	// Handle array types like [String]
	if ctx.LSQUARE() != nil && ctx.RSQUARE() != nil {
		typeExpr.IsArray = true
		if ctx.Type_() != nil {
			arrayElement := b.buildTypeExpression(ctx.Type_())
			typeExpr.ArrayElement = arrayElement
		}
	}

	return typeExpr
}

// buildModuleDecl builds a ModuleDeclaration from a parser module context.
func (b *Builder) buildModuleDecl(ctx parser.IModuleDeclContext) *ModuleDeclaration {
	name := ctx.ID().GetText()

	statements := make([]Statement, 0)
	if ctx.ModuleBody() != nil {
		for _, stmtCtx := range ctx.ModuleBody().AllModuleStatement() {
			var stmt Statement

			if stmtCtx.LetDecl() != nil {
				stmt = b.buildLetDecl(stmtCtx.LetDecl())
			} else if stmtCtx.FnDecl() != nil {
				stmt = b.buildFnDecl(stmtCtx.FnDecl())
			} else if stmtCtx.TypeDecl() != nil {
				stmt = b.buildTypeDecl(stmtCtx.TypeDecl())
			}

			if stmt != nil {
				statements = append(statements, stmt)
			}
		}
	}

	return &ModuleDeclaration{
		Name:       name,
		Statements: statements,
	}
}

// buildFunctionCall removed - no longer needed for type-level validation

// buildBlockBody builds a BlockExpression from a parser block body context.
func (b *Builder) buildBlockBody(ctx parser.IBlockBodyContext) *BlockExpression {
	statements := make([]Statement, 0)

	// Build all statements in the block
	for _, stmtCtx := range ctx.AllStatement() {
		stmt := b.buildStatement(stmtCtx)
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	// Build optional return expression
	var expr Expression
	if ctx.Expr() != nil {
		expr = b.buildExpression(ctx.Expr())
	}

	return &BlockExpression{
		Statements: statements,
		Expression: expr,
	}
}

// buildEffectDecl builds an EffectDeclaration from parser context - ALGEBRAIC EFFECTS SUPREMO!
func (b *Builder) buildEffectDecl(ctx parser.IEffectDeclContext) *EffectDeclaration {
	if ctx == nil {
		return nil
	}

	name := ctx.ID().GetText()
	operations := make([]EffectOperation, 0)

	// Parse effect operations
	for _, opCtx := range ctx.AllOpDecl() {
		operation := EffectOperation{
			Name: opCtx.ID().GetText(),
			Type: opCtx.Type_().GetText(), // Parse type as string for now
		}

		// Parse function type to extract parameters and return type
		typeExpr := b.buildTypeExpression(opCtx.Type_())
		if typeExpr != nil && typeExpr.IsFunction {
			// Extract parameters from function type
			operation.Parameters = make([]Parameter, len(typeExpr.ParameterTypes))
			for i, paramType := range typeExpr.ParameterTypes {
				operation.Parameters[i] = Parameter{
					Name: fmt.Sprintf("param%d", i), // Generate parameter names
					Type: &paramType,
				}
			}

			// Extract return type
			if typeExpr.ReturnType != nil {
				operation.ReturnType = typeExpr.ReturnType.Name
			}
		}

		operations = append(operations, operation)
	}

	return &EffectDeclaration{
		Name:       name,
		Operations: operations,
	}
}

// buildEffectSet parses effect signatures like !Logger or ![IO, Net]
func (b *Builder) buildEffectSet(ctx parser.IEffectSetContext) []string {
	if ctx == nil {
		return nil
	}

	effects := make([]string, 0)

	// Single effect: !Effect
	if ctx.ID() != nil {
		effects = append(effects, ctx.ID().GetText())
	}

	// Multiple effects: ![Effect1, Effect2]
	if ctx.EffectList() != nil {
		for _, idNode := range ctx.EffectList().AllID() {
			effects = append(effects, idNode.GetText())
		}
	}

	return effects
}
