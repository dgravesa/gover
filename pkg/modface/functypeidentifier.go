package modface

import (
	"fmt"
	"go/ast"
	"strings"
)

// FuncTypeIdentifier identifies a function type.
type FuncTypeIdentifier struct {
	ParamTypes  []TypeIdentifier // TODO: need to handle ...
	ResultTypes []TypeIdentifier
}

func joinTypeIDsToString(typeIDs []TypeIdentifier) string {
	typeStrings := []string{}
	for _, typeID := range typeIDs {
		typeStrings = append(typeStrings, typeID.String())
	}
	return strings.Join(typeStrings, ", ")
}

func (fi FuncTypeIdentifier) String() string {
	paramTypesStr := joinTypeIDsToString(fi.ParamTypes)
	resultTypesStr := joinTypeIDsToString(fi.ResultTypes)

	switch len(fi.ResultTypes) {
	case 0:
		return fmt.Sprintf("func(%s)", paramTypesStr)
	case 1:
		return fmt.Sprintf("func(%s) %s", paramTypesStr, resultTypesStr)
	default:
		return fmt.Sprintf("func(%s) (%s)", paramTypesStr, resultTypesStr)
	}
}

func (fi FuncTypeIdentifier) typeID() string {
	var sb strings.Builder
	sb.WriteString("func(")
	for _, paramType := range fi.ParamTypes {
		sb.WriteString(paramType.typeID())
		sb.WriteString(",")
	}
	sb.WriteString(")")
	for _, resultType := range fi.ResultTypes {
		sb.WriteString(resultType.typeID())
		sb.WriteString(",")
	}
	return sb.String()
}

func parseFuncTypeToTypeIdentifier(funcType *ast.FuncType, cache *importCache) (TypeIdentifier, error) {
	listTypeIdentifiersFromFieldList := func(fieldList *ast.FieldList) ([]TypeIdentifier, error) {
		if fieldList == nil {
			return nil, nil
		}
		typeIDs := []TypeIdentifier{}
		for _, field := range fieldList.List {
			// identify type
			typeID, err := parseExprToTypeIdentifier(field.Type, cache)
			if err != nil {
				return nil, err
			}
			// determine number of consecutive parameters with type
			numParamsWithType := maxInt(1, len(field.Names))
			for i := 0; i < numParamsWithType; i++ {
				typeIDs = append(typeIDs, typeID)
			}
		}
		return typeIDs, nil
	}

	paramTypes, err := listTypeIdentifiersFromFieldList(funcType.Params)
	if err != nil {
		return nil, err
	}

	resultTypes, err := listTypeIdentifiersFromFieldList(funcType.Results)
	if err != nil {
		return nil, err
	}

	return FuncTypeIdentifier{
		ParamTypes:  paramTypes,
		ResultTypes: resultTypes,
	}, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
