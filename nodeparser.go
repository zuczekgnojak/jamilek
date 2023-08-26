package jamilek

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type NodeParser struct {
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

func (p NodeParser) parseObject() (*Node, error) {
	object := make(map[string]*Node)

	token, err := p.tokenizer.Next()
	if token.Type != ObjectStart {
		fmt.Println(token, err)
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
		fmt.Println("key: ", key)

		node, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		object[key] = node
	}

	token, err = p.tokenizer.Next()
	fmt.Println(token, err)
	if token.Type != ObjectEnd {
		return nil, errors.New("expecting object end")
	}
	if err != nil {
		return nil, err
	}

	node := &Node{Object, object}
	return node, nil
}

func (p NodeParser) parseArray() (*Node, error) {
	token, err := p.tokenizer.Next()
	if err != nil {
		return nil, err
	}

	if token.Type != ArrayStart {
		return nil, errors.New("not an array")
	}

	value := make([]*Node, 0)
	i := 0
	for {
		token, err = p.tokenizer.Peek()
		fmt.Println("PEEKED", token)
		if token.Type == Value {
			fmt.Println("parsing in array", token)
			node, err := p.parseValue()
			fmt.Println("NODE", node, err)
			if err != nil {
				return nil, err
			}
			value = append(value, node)
			// p.tokenizer.Next()
		} else {
			break
		}
		i += 1
		if i > 5 {
			panic("too much looping")
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

func (p NodeParser) parseValue() (*Node, error) {

	token, err := p.tokenizer.Peek()
	fmt.Println("parsing a value", token)
	if token.Type == ObjectStart {
		return p.parseObject()
	}
	if token.Type == ArrayStart {
		return p.parseArray()
	}
	if token.Type != Value {
		return nil, errors.New("expecting value")
	}

	fmt.Println("parseValue NEXT")
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

func (p NodeParser) Parse() (*Node, error) {
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

func NewNodeParser(reader io.Reader) NodeParser {
	tokenizer := NewTokenizer(reader)
	return NodeParser{&tokenizer}
}
