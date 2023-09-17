package jamilek

import (
	"errors"
	"strconv"
)

type NodeType int

const (
	Object NodeType = iota
	Array
	String
	Float
	Integer
	Bool
)

type Node struct {
	nodeType NodeType
	value    interface{}
}

func (n Node) Type() NodeType {
	return n.nodeType
}

func (n Node) Get(path ...string) (*Node, error) {
	if len(path) == 0 {
		return &n, nil
	}
	key := path[0]

	if n.nodeType == Object {
		object := n.value.(map[string]*Node)
		node := object[key]
		if node == nil {
			return nil, errors.New("invalid key")
		}
		return node.Get(path[1:]...)
	}
	if n.nodeType == Array {
		array := n.value.([]*Node)
		index, err := strconv.Atoi(key)
		if err != nil {
			return nil, errors.New("non numeric array index")
		}
		if index >= len(array) {
			return nil, errors.New("out of range index")
		}
		node := array[index]
		return node.Get(path[1:]...)

	}
	return nil, errors.New("invalid path")
}

func (n Node) GetArray(path ...string) ([]*Node, error) {
	if len(path) != 0 {
		node, err := n.Get(path...)
		if err != nil {
			return nil, err
		}
		return node.GetArray()
	}
	if n.nodeType != Array {
		return nil, errors.New("invalid array")
	}
	return n.value.([]*Node), nil
}

func (n Node) GetString(path ...string) (string, error) {
	if len(path) != 0 {
		node, err := n.Get(path...)
		if err != nil {
			return "", err
		}
		return node.GetString()
	}
	if n.nodeType != String {
		return "", errors.New("invalid string")
	}
	return n.value.(string), nil
}

func (n Node) GetBool(path ...string) (bool, error) {
	if len(path) != 0 {
		node, err := n.Get(path...)
		if err != nil {
			return false, err
		}
		return node.GetBool()
	}
	if n.nodeType != Bool {
		return false, errors.New("invalid bool")
	}
	return n.value.(bool), nil
}

func (n Node) GetFloat(path ...string) (float64, error) {
	if len(path) != 0 {
		node, err := n.Get(path...)
		if err != nil {
			return 0, err
		}
		return node.GetFloat()
	}
	if n.nodeType != Float {
		return 0, errors.New("invalid float")
	}
	return n.value.(float64), nil
}

func (n Node) GetInteger(path ...string) (int64, error) {
	if len(path) != 0 {
		node, err := n.Get(path...)
		if err != nil {
			return 0, err
		}
		return node.GetInteger()
	}
	if n.nodeType != Integer {
		return 0, errors.New("invalid int")
	}
	return n.value.(int64), nil
}


func (n Node) String() string {
	if n.nodeType == String {
		return "\"" + n.value.(string) + "\""
	}
	if n.nodeType == Object {
		result := "{"
		for key, node := range n.value.(map[string]*Node) {
			entry := " " + key + ": " + node.String()
			result += entry
		}
		result += " }"
		return result
	}
	if n.nodeType == Array {
		result := "["
		for _, node := range n.value.([]*Node) {
			entry := " " + node.String()
			result += entry
		}
		result += " ]"
		return result
	}
	if n.nodeType == Integer {
		return strconv.FormatInt(n.value.(int64), 10)
	}
	if n.nodeType == Float {
		return strconv.FormatFloat(n.value.(float64), 'f', -1, 64)
	}
	if n.nodeType == Bool {
		return strconv.FormatBool(n.value.(bool))
	}
	return "ERROR"
}
