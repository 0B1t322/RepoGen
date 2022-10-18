package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"golang.org/x/tools/go/ast/inspector"
)

//go:embed edges_builder.tmpl
var edgesBuilderTemplate string

//go:embed filter_builder.tmpl
var filterBuilderTemplate string

//go:embed sort_builder.tmpl
var sortBuilderTemplate string

const (
	edgesBuilderFileName  = "edges_builder.go"
	filterBuilderFileName = "filter_builder.go"
	sortBuilderFileName   = "sort_builder.go"
)

type (
	RepoFileAST struct {
		PackageName string
		Edges       []EdgesAST
		Filter      FilterAST
		Sort        SortAST
	}

	EdgesAST struct {
		TypeSpec   *ast.TypeSpec
		StructType *ast.StructType
	}

	FilterAST struct {
		FilterQueryTypeSpec    *ast.TypeSpec
		FilterFieldsTypeSpec   *ast.TypeSpec
		FilterFieldsStructType *ast.StructType
		Imports                []*ast.ImportSpec
	}

	SortAST struct {
		SortFieldsTypeSpec       *ast.TypeSpec
		SortFieldsStructTypeSpec *ast.StructType
	}
)

func (q FilterAST) IsPresent() bool {
	return q.FilterQueryTypeSpec != nil && q.FilterFieldsTypeSpec != nil && q.FilterFieldsStructType != nil
}

func (s SortAST) IsPresent() bool {
	return s.SortFieldsTypeSpec != nil && s.SortFieldsStructTypeSpec != nil
}

type RepoFile struct {
	EdgesData  *EdgesBuilderData
	FilterData *FilterBuilderData
	SortData   *SortBuilderData
}

type EdgesBuilderData struct {
	Package     string
	EdgesFields []EdgesField
}

type EdgesField struct {
	TypeName        string
	Name            string
	Edges           []EdgesField
	VariablesFields []VariablesField
	WithEdgesField  bool
}

type VariablesField struct {
	Name         string
	VariableType string
}

type FilterBuilderData struct {
	Package          string
	Imports          []string
	FilterQueryType  string
	FilterFieldsType string
	FilterFields     []FilterField
}

type FilterField struct {
	Name string
	Type string
	// if require import some package will contain last package word
	RequireImport *string
}

type SortBuilderData struct {
	Package        string
	SortFieldsType string
	SortFields     []SortField
}

type SortField struct {
	Name string
}

func main() {
	file := os.Getenv("GOFILE")
	if file == "" {
		log.Fatal("GOFILE must be set")
	}

	if !isFileNeedRegeneration(file) {
		return
	}

	repoFile := parseRepoFile(file)

	if repoFile.EdgesData != nil {
		processTemplate(
			"edges_builder.tmpl",
			edgesBuilderTemplate,
			edgesBuilderFileName,
			*repoFile.EdgesData,
		)
	}

	if repoFile.FilterData != nil {
		processTemplate(
			"filter_builder.tmpl",
			filterBuilderTemplate,
			filterBuilderFileName,
			*repoFile.FilterData,
		)
	}

	if repoFile.SortData != nil {
		processTemplate(
			"sort_builder.tmpl",
			sortBuilderTemplate,
			sortBuilderFileName,
			*repoFile.SortData,
		)
	}
}

func isFileNeedRegeneration(fileName string) bool {
	file, err := os.Stat(fileName)
	if err != nil {
		return true
	}

	// Check that one of file is exist
	if _, err := os.Open(edgesBuilderFileName); err != nil {
		return true
	}

	if _, err := os.Open(filterBuilderFileName); err != nil {
		return true
	}

	if _, err := os.Open(sortBuilderFileName); err != nil {
		return true
	}

	if stat, err := os.Stat(edgesBuilderFileName); err == nil {
		return file.ModTime().After(stat.ModTime())
	}

	if stat, err := os.Stat(filterBuilderFileName); err == nil {
		return file.ModTime().After(stat.ModTime())
	}

	if stat, err := os.Stat(sortBuilderFileName); err == nil {
		return file.ModTime().After(stat.ModTime())
	}

	return false
}

func parseRepoFile(fileName string) RepoFile {
	file, err := parser.ParseFile(
		token.NewFileSet(),
		fileName,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		log.Fatal(err)
	}

	repoFileAst := &RepoFileAST{
		PackageName: file.Name.Name,
	}

	parseEdges(repoFileAst, file)
	parseFilter(repoFileAst, file)
	parseSort(repoFileAst, file)

	return fileRepoAstToFileRepo(repoFileAst)
}

func parseEdges(src *RepoFileAST, file *ast.File) {
	i := inspector.New([]*ast.File{file})

	iFilter := []ast.Node{
		&ast.GenDecl{},
	}

	i.Nodes(
		iFilter,
		func(n ast.Node, push bool) bool {
			genDecl := n.(*ast.GenDecl)
			if genDecl.Doc == nil {
				return false
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					return false
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return false
				}

				for _, comment := range genDecl.Doc.List {
					switch comment.Text {
					case "//repogen:edges":
						src.Edges = append(src.Edges, EdgesAST{
							TypeSpec:   typeSpec,
							StructType: structType,
						})
					}
				}
			}

			return false
		},
	)
}

func parseFilter(src *RepoFileAST, file *ast.File) {
	i := inspector.New([]*ast.File{file})

	iFilter := []ast.Node{
		&ast.GenDecl{},
	}

	i.Nodes(
		iFilter,
		func(n ast.Node, push bool) bool {
			genDecl := n.(*ast.GenDecl)

			if genDecl.Tok == token.IMPORT {
				for _, spec := range genDecl.Specs {
					importSpec, ok := spec.(*ast.ImportSpec)
					if !ok {
						continue
					}

					src.Filter.Imports = append(src.Filter.Imports, importSpec)
				}
			}

			if genDecl.Doc == nil {
				return false
			}

			for _, spec := range genDecl.Specs {
				for _, comment := range genDecl.Doc.List {
					switch comment.Text {
					case "//repogen:filter":
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							return false
						}

						switch t := typeSpec.Type.(type) {
						case *ast.StructType:
							src.Filter.FilterFieldsTypeSpec = typeSpec
							src.Filter.FilterFieldsStructType = t
						case *ast.IndexExpr:
							src.Filter.FilterQueryTypeSpec = typeSpec
						}

					}
				}
			}

			return false
		},
	)
}

func parseSort(src *RepoFileAST, file *ast.File) {
	i := inspector.New([]*ast.File{file})

	iFilter := []ast.Node{
		&ast.GenDecl{},
	}

	i.Nodes(
		iFilter,
		func(n ast.Node, push bool) bool {
			genDecl := n.(*ast.GenDecl)
			if genDecl.Doc == nil {
				return false
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					return false
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return false
				}

				for _, comment := range genDecl.Doc.List {
					switch comment.Text {
					case "//repogen:sort":
						src.Sort.SortFieldsTypeSpec = typeSpec
						src.Sort.SortFieldsStructTypeSpec = structType
					}
				}
			}

			return false
		},
	)
}

func fileRepoAstToFileRepo(repoFile *RepoFileAST) RepoFile {
	f := RepoFile{
		EdgesData:  fileRepoAstToEdge(repoFile),
		FilterData: fileRepoAstToQuery(repoFile),
		SortData:   fileRepoAstToSort(repoFile),
	}

	return f
}

func fileRepoAstToQuery(repoFile *RepoFileAST) *FilterBuilderData {
	if query := repoFile.Filter; query.IsPresent() {
		queryData := &FilterBuilderData{
			Package: repoFile.PackageName,
		}

		queryData.FilterQueryType = query.FilterQueryTypeSpec.Name.Name
		queryData.FilterFieldsType = query.FilterFieldsTypeSpec.Name.Name

		queryData.FilterFields = initFilterFields(query.FilterFieldsStructType.Fields.List)

		// find required imports
		for _, field := range queryData.FilterFields {
			if requiredImport := field.RequireImport; requiredImport != nil {
				queryData.Imports = append(
					queryData.Imports,
					func() string {
						imp, ok := lo.Find(
							repoFile.Filter.Imports,
							func(importSpec *ast.ImportSpec) bool {
								name, _ := lo.Last(strings.Split(importSpec.Path.Value, "/"))
								return strings.ReplaceAll(name, `"`, "") == *requiredImport
							},
						)
						if !ok {
							return ""
						}

						return imp.Path.Value
					}(),
				)
			}
		}

		queryData.Imports = lo.Uniq(queryData.Imports)

		return queryData
	}

	return nil
}

func fileRepoAstToEdge(repoFile *RepoFileAST) *EdgesBuilderData {
	if len(repoFile.Edges) > 0 {
		edgesData := &EdgesBuilderData{
			Package: repoFile.PackageName,
		}

		// Init root edge
		rootEdge := initRootEdge(repoFile.Edges)

		edgesData.EdgesFields = append(edgesData.EdgesFields, rootEdge)
		edgesData.EdgesFields = append(edgesData.EdgesFields, initEdges(repoFile.Edges[1:])...)

		return edgesData
	}

	return nil
}

func fileRepoAstToSort(repoFile *RepoFileAST) *SortBuilderData {
	if sort := repoFile.Sort; sort.IsPresent() {
		sortData := &SortBuilderData{
			Package: repoFile.PackageName,
		}

		sortData.SortFieldsType = sort.SortFieldsTypeSpec.Name.Name
		sortData.SortFields = initSortFields(sort.SortFieldsStructTypeSpec.Fields.List)

		return sortData
	}

	return nil
}

func initRootEdge(edgesAST []EdgesAST) EdgesField {
	rootEdge := EdgesField{}
	{
		edge := edgesAST[0]

		rootEdge.TypeName = edge.TypeSpec.Name.Name
		for _, field := range edge.StructType.Fields.List {
			rootEdge.Edges = append(
				rootEdge.Edges,
				EdgesField{
					TypeName:        field.Type.(*ast.IndexExpr).Index.(*ast.Ident).Name,
					Name:            field.Names[0].Name,
					Edges:           []EdgesField{},
					VariablesFields: []VariablesField{},
				},
			)
		}
	}

	return rootEdge
}

func initEdges(edgesAST []EdgesAST) []EdgesField {
	var edges []EdgesField
	{
		for _, edge := range edgesAST {
			edgeField := EdgesField{}
			// If contains Edges field mean that we need init variables and edges in another field
			if lo.ContainsBy(
				edge.StructType.Fields.List,
				func(field *ast.Field) bool {
					return field.Names[0].Name == "Edges"
				},
			) {
				edgeField.TypeName = edge.TypeSpec.Name.Name
				edgeField.Edges = append(
					edgeField.Edges,
					func() (slice []EdgesField) {
						for _, field := range edge.StructType.Fields.List {
							if field.Names[0].Name == "Edges" {
								edgeType := field.Type.(*ast.Ident).Name
								edgeAst, _ := lo.Find(
									edgesAST,
									func(edge EdgesAST) bool {
										return edge.TypeSpec.Name.Name == edgeType
									},
								)

								edges := initEdges([]EdgesAST{edgeAst})

								lo.ForEach(
									edges,
									func(e EdgesField, _ int) {
										lo.ForEach(
											e.Edges,
											func(_ EdgesField, i int) {
												e.Edges[i].WithEdgesField = true
											},
										)
										slice = append(
											slice,
											e.Edges...,
										)
									},
								)
							}
						}

						return slice
					}()...,
				)

				edgeField.VariablesFields = append(
					edgeField.VariablesFields,
					func() (slice []VariablesField) {
						for _, field := range edge.StructType.Fields.List {
							if field.Names[0].Name != "Edges" {
								slice = append(slice, initVariableField(field))
							}
						}
						return slice
					}()...,
				)
			} else {
				edgeField = initRootEdge([]EdgesAST{edge})
			}

			edges = append(edges, edgeField)
		}
	}

	return edges
}

func initSubEdges(edgesAST []EdgesAST, rootEdge *EdgesField) {
	for i, edge := range rootEdge.Edges {
		// Find edge in AST
		spec, _ := lo.Find(
			edgesAST,
			func(e EdgesAST) bool {
				return e.TypeSpec.Name.Name == edge.TypeName
			},
		)

		for _, field := range spec.StructType.Fields.List {
			// Find field field that we need init as edge
			if field.Names[0].Name == "Edges" {
				typeName := field.Type.(*ast.Ident).Name

				newEdgeAst, _ := lo.Find(
					edgesAST,
					func(e EdgesAST) bool {
						return e.TypeSpec.Name.Name == typeName
					},
				)

				newEdge := initRootEdge([]EdgesAST{newEdgeAst})
				initSubEdges(edgesAST, &newEdge)
				edge.Edges = append(edge.Edges, newEdge)

			} else {
				edge.VariablesFields = append(edge.VariablesFields, initVariableField(field))
			}
		}
		rootEdge.Edges[i] = edge
	}
}

func initVariableField(field *ast.Field) VariablesField {
	variableField := VariablesField{}
	{
		variableField.Name = field.Names[0].Name
		variableField.VariableType = field.Type.(*ast.IndexExpr).Index.(*ast.Ident).Name
	}

	return variableField
}

func initFilterFields(fields []*ast.Field) []FilterField {
	var filterFields []FilterField
	{
		for _, field := range fields {
			filterField := FilterField{}
			{
				filterField.Name = field.Names[0].Name

				if filterField.Name == "AttachedTo" {

				}

				wrappedType := field.Type.(*ast.IndexExpr).Index.(*ast.IndexExpr)

				switch t := wrappedType.Index.(type) {
				case *ast.Ident:
					filterField.Type = t.Name
				case *ast.ArrayType:
					filterField.Type = "[]" + t.Elt.(*ast.Ident).Name
				case *ast.SelectorExpr:
					packageLastWord := t.X.(*ast.Ident).Name
					filterField.Type = packageLastWord + "." + t.Sel.Name
					filterField.RequireImport = &packageLastWord
				}
			}

			filterFields = append(filterFields, filterField)
		}
	}

	return filterFields
}

func initSortFields(fields []*ast.Field) []SortField {
	var sortFields []SortField
	{
		for _, field := range fields {
			sortField := SortField{}
			{
				sortField.Name = field.Names[0].Name
			}

			sortFields = append(sortFields, sortField)
		}
	}

	return sortFields
}

func processTemplate(tmplName string, tmplText string, outFile string, data any) {
	funcMap := template.FuncMap{
		"ToLowerCamel": strcase.ToLowerCamel,
	}

	tmpl := template.Must(
		template.New(tmplName).Funcs(funcMap).Parse(
			tmplText,
		),
	)

	tmpl = tmpl.Funcs(funcMap)

	var processed bytes.Buffer

	err := tmpl.ExecuteTemplate(&processed, tmplName, data)
	if err != nil {
		panic(err)
	}

	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		log.Fatalf("Could not format processed template: %v\n", err)
	}

	outputPath := outFile
	f, _ := os.Create(outputPath)
	w := bufio.NewWriter(f)
	w.WriteString(string(formatted))
	w.Flush()
}
