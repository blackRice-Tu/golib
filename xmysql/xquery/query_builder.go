package xquery

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type QueryMap map[string]any

const (
	commandSep = "__"

	commandExact       = "exact"
	commandIn          = "in"
	commandContains    = "contains"
	commandIContains   = "icontains"
	commandGt          = "gt"
	commandGte         = "gte"
	commandLt          = "lt"
	commandLte         = "lte"
	commandStartswith  = "startswith"
	commandIStartswith = "istartswith"
	commandEndswith    = "endswith"
	commandIEndswith   = "iendswith"

	orderSep = ","
)

// like django
var sqlCommandOperateMap = map[string]string{
	commandExact:       "= ?",
	commandIn:          "IN (?)",
	commandContains:    "LIKE BINARY ?",
	commandIContains:   "LIKE ?",
	commandGt:          "> ?",
	commandGte:         ">= ?",
	commandLt:          "< ?",
	commandLte:         "<= ?",
	commandStartswith:  "LIKE BINARY ?",
	commandIStartswith: "LIKE ?",
	commandEndswith:    "LIKE BINARY ?",
	commandIEndswith:   "LIKE ?",
}

var stringCommandMap = map[string]struct{}{
	commandContains:    {},
	commandIContains:   {},
	commandStartswith:  {},
	commandIStartswith: {},
	commandEndswith:    {},
	commandIEndswith:   {},
}

func isStringCommand(command string) bool {
	_, ok := stringCommandMap[command]
	return ok
}

func SqlQueryBuilderV1(ctx context.Context, queryMap QueryMap, fieldMap map[string]string) (query string, args []any, e error) {
	queryList := make([]string, 0)
	for k, v := range queryMap {
		field, command := getFieldAndCommand(k)
		if field == "" {
			continue
		}
		if newField, ok := fieldMap[field]; ok {
			field = newField
		}
		operate, value, err := getOperateAndValue(command, v)
		if err != nil {
			e = errors.WithMessagef(err, "field '%s'", k)
			return
		}
		queryList = append(queryList, field+" "+operate)
		args = append(args, value)
	}
	if len(queryList) == 0 {
		queryList = append(queryList, "1=1")
		args = make([]any, 0)
	}
	query = strings.Join(queryList, " AND ")
	return
}

func getFieldAndCommand(key string) (field string, command string) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}

	command = commandExact
	items := strings.Split(key, commandSep)
	size := len(items)
	switch size {
	case 1:
		field = items[0]
	case 2:
		field = items[0]
		command = items[1]
	default:
		field = strings.Join(items[0:size-1], commandSep)
		command = items[size-1]
	}
	return
}

func getOperateAndValue(command string, rawValue any) (operate string, value any, e error) {
	operate, ok := sqlCommandOperateMap[command]
	if !ok {
		e = errors.Errorf("command '%s' is invalided", command)
		return
	}

	kind := reflect.ValueOf(rawValue).Kind()
	if isStringCommand(command) && kind != reflect.String {
		e = errors.Errorf("value %v can only be string", rawValue)
		return
	}
	value = rawValue
	switch command {
	case commandContains, commandIContains:
		value = "%%" + value.(string) + "%%"
	case commandStartswith, commandIStartswith:
		value = value.(string) + "%%"
	case commandEndswith, commandIEndswith:
		value = "%%" + value.(string)
	}
	return
}

func SqlOrderBuilderV1(ctx context.Context, s string, fieldMap map[string]string) (order string) {
	sort := strings.TrimSpace(s)
	if sort == "" {
		return
	}
	orderList := make([]string, 0)
	items := strings.Split(s, orderSep)
	for i, _ := range items {
		item := strings.TrimSpace(items[i])
		field := ""
		desc := ""
		if strings.HasPrefix(item, "-") {
			desc = " DESC"
			field = item[1:]
		} else {
			field = item
		}
		if newField, ok := fieldMap[field]; ok {
			field = newField
		}

		if item != "" {
			orderList = append(orderList, fmt.Sprintf("%s %s", field, desc))
		}
	}
	if len(orderList) == 0 {
		return
	}
	order = strings.Join(orderList, ",")
	return
}

func GetOrNewQueryMap(queryMap QueryMap) QueryMap {
	if queryMap == nil {
		queryMap = make(QueryMap)
	}
	return queryMap
}
