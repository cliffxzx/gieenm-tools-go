package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// BulkInsertParameter Generate []*struct or []struct to col, row, args
func BulkInsertParameter(unsavedRows interface{}) (*string, *string, []interface{}, error) {
	rows, err := InterfaceToSlice(unsavedRows)
	if err != nil {
		return nil, nil, nil, err
	}

	colStrings := []string{}

	if len(rows) < 1 {
		fmt.Println("The parameter must greater than zero")
		return nil, nil, nil, errors.New("The parameter must greater than zero")
	}

	v := reflect.ValueOf(rows[0])
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		colStrings = append(colStrings, typeOfS.Field(i).Tag.Get("db"))
	}

	rowStrings := []string{}
	rowArgs := []interface{}{}
	i := 0
	for _, post := range rows {
		v := reflect.ValueOf(post)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		valueStrings := []string{}
		for j := 0; j < v.NumField(); j++ {
			i++
			valueStrings = append(valueStrings, fmt.Sprintf("$%d", i))
			rowArgs = append(rowArgs, v.Field(j).Interface())
		}
		rowStrings = append(rowStrings, fmt.Sprintf("(%s)", strings.Join(valueStrings, ",")))
	}

	colString := fmt.Sprintf("(%s)", strings.Join(colStrings, ","))
	rowString := strings.Join(rowStrings, ",")

	return &colString, &rowString, rowArgs, nil
}
