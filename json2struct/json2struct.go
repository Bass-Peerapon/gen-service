package json2struct

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

// Given a JSON string representation of an object and a name structName,
// attemp to generate a struct definition
func Generate(input io.Reader, structName, pkgName string) ([]byte, error) {
	var iresult interface{}
	var result map[string]interface{}
	if err := json.NewDecoder(input).Decode(&iresult); err != nil {
		return nil, err
	}

	switch iresult := iresult.(type) {
	case map[string]interface{}:
		result = iresult
	case []map[string]interface{}:
		if len(iresult) > 0 {
			result = iresult[0]
		} else {
			return nil, fmt.Errorf("empty array")
		}
	case []interface{}:
		src := fmt.Sprintf("package %s\n\ntype %s %s\n",
			pkgName,
			structName,
			"[]interface{}")
		return []byte(src), nil

	default:
		return nil, fmt.Errorf("unexpected type: %T", iresult)
	}

	src := fmt.Sprintf("package %s\ntype %s %s}",
		pkgName,
		structName,
		generateTypes(structName, result, 0))
	formatted, err := format.Source([]byte(src))
	if err != nil {
		err = fmt.Errorf("error formatting: %s, was formatting\n%s", err, src)
	}
	return formatted, err
}

// Generate go struct entries for a map[string]interface{} structure
func generateTypes(structName string, obj map[string]interface{}, depth int) string {
	structure := fmt.Sprintf("struct { TableName         struct{}               `json:\"-\" db:\"%s\" pk:\"Id\"`", strcase.ToSnake(structName))

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := obj[key]
		valueType := typeForValue(value)
		typ := valueType
		db := key

		//If a nested value, recurse
		switch value := value.(type) {
		case []map[string]interface{}:
			typ = "-"
			db = "-"
			valueType = "[]" + generateTypes(fmtFieldName(strcase.ToSnake(key)), value[0], depth+1) + "}"
			goto FILENAME
		case map[string]interface{}:
			typ = "-"
			db = "-"
			valueType = generateTypes(fmtFieldName(strcase.ToSnake(key)), value, depth+1) + "}"
			goto FILENAME
		}
		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice:
			typ = "-"
			db = "-"

		case reflect.Array:
			typ = "-"
			db = "-"
		}

	FILENAME:
		fieldName := fmtFieldName(stringifyFirstChar(key))
		structure += fmt.Sprintf("\n%s %s `json:\"%s\" db:\"%s\" type:\"%s\"`",
			fieldName,
			valueType,
			key,
			db,
			typ,
		)
	}
	return structure
}

// fmtFieldName formats a string as a struct key
//
// Example:
//
//	fmtFieldName("foo_id")
//
// Output: FooID
func fmtFieldName(s string) string {
	name := lintFieldName(s)
	runes := []rune(name)
	for i, c := range runes {
		ok := unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = '_'
		}
	}
	return string(runes)
}

func lintFieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		runes := []rune(name)
		if u := strings.ToUpper(name); commonInitialisms[u] {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return string(runes)
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))

		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

// generate an appropriate struct type entry
func typeForValue(value interface{}) string {
	//Check if this is an array
	if objects, ok := value.([]interface{}); ok {
		types := make(map[reflect.Type]bool, 0)
		for _, o := range objects {
			types[reflect.TypeOf(o)] = true
		}
		if len(types) == 1 {
			return "[]" + typeForValue(objects[0])
		}
		return "[]interface{}"
	} else if object, ok := value.(map[string]interface{}); ok {
		return generateTypes("", object, 0) + "}"
	} else if reflect.TypeOf(value) == nil {
		return "interface{}"
	}
	v := reflect.TypeOf(value).Name()
	if v == "float64" {
		v = disambiguateFloatInt(value)
	}
	return v
}

// All numbers will initially be read as float64
// If the number appears to be an integer value, use int instead
func disambiguateFloatInt(value interface{}) string {
	const epsilon = .0001
	vfloat := value.(float64)
	if math.Abs(vfloat-math.Floor(vfloat+epsilon)) < epsilon {
		var tmp int = 1
		return reflect.TypeOf(tmp).Name()
	}
	return reflect.TypeOf(value).Name()
}

// convert first character ints to strings
func stringifyFirstChar(str string) string {
	first := str[:1]

	i, err := strconv.ParseInt(first, 10, 8)

	if err != nil {
		return str
	}

	return intToWordMap[i] + "_" + str[1:]
}
