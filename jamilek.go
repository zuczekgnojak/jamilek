package jamilek

import "fmt"
import "io"
import "errors"
import "strings"
import "strconv"

func Hello() {
	fmt.Println("The best parser of them all")
}

type Parser struct {
	tokenizer Tokenizer
}

func convValue(value string) interface{} {
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	number, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return number
	}

	first := value[0]
	last := value[len(value)-1]
	if first != '"' || last != '"' {
		panic("invalid string")
	}
	if len(value) == 2 {
		return ""
	}

	return value[1 : len(value)-1]
}
func (p Parser) parseObject() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for {
		token, _ := p.tokenizer.NextToken()
		if token.Type == ObjectEnd || token.Type == EOF {
			return result, nil
		}
		if token.Type != Key {
			panic("expected object key")
		}

		key := strings.Replace(token.Value, ":", "", -1)

		token, _ = p.tokenizer.NextToken()

		var value interface{}
		if token.Type == ObjectStart {
			value, _ = p.parseObject()
		}
		if token.Type == ArrayStart {
			value, _ = p.parseArray()
		}

		if token.Type == Value {
			value = convValue(token.Value)
		}
		result[key] = value
	}
	return result, nil
}

func (p Parser) parseArray() ([]interface{}, error) {
	result := make([]interface{}, 0)
	for {
		token, _ := p.tokenizer.NextToken()

		if token.Type == ArrayEnd || token.Type == EOF {
			return result, nil
		}

		var value interface{}
		if token.Type == ObjectStart {
			value, _ = p.parseObject()
		}
		if token.Type == ArrayStart {
			value, _ = p.parseArray()
		}
		if token.Type == Value {
			value = convValue(token.Value)
		}
		result = append(result, value)
	}
	return result, nil
}

func (p Parser) Parse() (map[string]interface{}, error) {

	token, _ := p.tokenizer.NextToken()
	if token.Type != ObjectStart {
		fmt.Println("ERRROR", token)
		return nil, errors.New("expected object at the root of document")
	}

	return p.parseObject()
}

func NewParser(reader io.Reader) Parser {
	tokenizer := NewTokenizer(reader)
	return Parser{tokenizer}
}
