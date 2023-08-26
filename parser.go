package jamilek

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

type Parser struct {
	tokenizer *Tokenizer
}

func convBool(value string) (bool, error) {
	if value == "true" {
		return true, nil
	}
	if value == "false" {
		return false, nil
	}

	return false, errors.New("not a bool")
}

func convNumber(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func convString(value string) (string, error) {
	first := value[0]
	last := value[len(value)-1]
	if first != '"' || last != '"' {
		return "", errors.New("invalid string")
	}
	return value[1 : len(value)-1], nil
}

func (p Parser) parseObject() (*Node, error) {
	object := make(map[string]*Node)

	token, err := p.tokenizer.Next()
	if token.Type != ObjectStart {
		return nil, errors.New("expecting object start")
	}
	if err != nil {
		return nil, err
	}

	for {
		token, err = p.tokenizer.Peek()
		if err != nil {
			return nil, err
		}
		if token.Type != Key {
			break
		}

		token, err = p.tokenizer.Next()
		if err != nil {
			return nil, err
		}
		key := strings.Replace(token.Value, ":", "", -1)

		node, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		object[key] = node
	}

	token, err = p.tokenizer.Next()
	if token.Type != ObjectEnd {
		return nil, errors.New("expecting object end")
	}
	if err != nil {
		return nil, err
	}

	node := &Node{Object, object}
	return node, nil
}

func (p Parser) parseArray() (*Node, error) {
	token, err := p.tokenizer.Next()
	if err != nil {
		return nil, err
	}

	if token.Type != ArrayStart {
		return nil, errors.New("not an array")
	}

	value := make([]*Node, 0)
	for {
		token, err = p.tokenizer.Peek()
		if token.Type == Value {
			node, err := p.parseValue()
			if err != nil {
				return nil, err
			}
			value = append(value, node)
		} else {
			break
		}
	}

	token, err = p.tokenizer.Next()
	if err != nil {
		return nil, err
	}
	if token.Type != ArrayEnd {
		return nil, errors.New("invalid array end")
	}
	return &Node{Array, value}, nil
}

func (p Parser) parseValue() (*Node, error) {
	token, err := p.tokenizer.Peek()

	if token.Type == ObjectStart {
		return p.parseObject()
	}
	if token.Type == ArrayStart {
		return p.parseArray()
	}
	if token.Type != Value {
		return nil, errors.New("expecting value")
	}

	p.tokenizer.Next()

	if err != nil {
		return nil, err
	}
	var value interface{}
	value, err = convBool(token.Value)
	if err == nil {
		return &Node{Bool, value}, nil
	}

	value, err = convNumber(token.Value)
	if err == nil {
		return &Node{Number, value}, nil
	}

	value, err = convString(token.Value)
	if err != nil {
		return nil, err
	}
	return &Node{String, value}, nil
}

func (p Parser) Parse() (*Node, error) {
	node, err := p.parseObject()
	if err != nil {
		return nil, err
	}
	token, err := p.tokenizer.Next()
	if err != nil {
		return nil, err
	}
	if token.Type != EOF {
		return nil, errors.New("expecting file end")
	}
	return node, nil
}

func NewParser(reader io.Reader) Parser {
	tokenizer := NewTokenizer(reader)
	return Parser{&tokenizer}
}
