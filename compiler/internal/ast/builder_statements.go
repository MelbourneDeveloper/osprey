package ast

import (
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

	return &LetDeclaration{
		Name:    name,
		Mutable: mutable,
		Type:    nil, // For now, type annotation support comes later
		Value:   value,
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

	return &FunctionDeclaration{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
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

				// Handle optional constraint
				if fieldCtx.Constraint() != nil {
					constraint := b.buildFunctionCall(fieldCtx.Constraint().FunctionCall())
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

	return &TypeDeclaration{
		Name:       name,
		TypeParams: typeParams,
		Variants:   variants,
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

			// Handle optional constraint
			if fieldCtx.Constraint() != nil {
				constraint := b.buildFunctionCall(fieldCtx.Constraint().FunctionCall())
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

	return &TypeExpression{
		Name: ctx.ID().GetText(),
		// TODO: Add support for generic types and arrays when needed
	}
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

// buildFunctionCall builds a FunctionCallExpression from a parser function call context.
func (b *Builder) buildFunctionCall(ctx parser.IFunctionCallContext) *FunctionCallExpression {
	if ctx == nil {
		return nil
	}

	functionName := ctx.ID().GetText()

	var arguments []Expression
	if ctx.ArgList() != nil {
		args, _ := b.buildArguments(ctx.ArgList()) // Ignore named args for now
		arguments = args
	}

	return &FunctionCallExpression{
		Function:  functionName,
		Arguments: arguments,
	}
}

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
