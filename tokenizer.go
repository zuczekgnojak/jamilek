package jamilek

import (
	"bufio"
	"fmt"
	"io"
)

type TokenType int64

const (
	ObjectStart TokenType = iota
	ObjectEnd
	ArrayStart
	ArrayEnd
	Key
	Value
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Tokenizer struct {
	scanner *bufio.Scanner
	peeked  *Token
}

func (t Tokenizer) nextWord() (string, error) {
	scanning := t.scanner.Scan()
	if !scanning {
		return "", t.scanner.Err()
	}

	return t.scanner.Text(), nil
}

func (t Tokenizer) NextToken() (Token, error) {
	word, err := t.nextWord()

	if word == "" {
		return Token{EOF, ""}, err
	}
	if word == "{" {
		return Token{ObjectStart, word}, nil
	}
	if word == "}" {
		return Token{ObjectEnd, word}, nil
	}
	if word == "[" {
		return Token{ArrayStart, word}, nil
	}
	if word == "]" {
		return Token{ArrayEnd, word}, nil
	}
	if word[len(word)-1] == ':' {
		return Token{Key, word}, nil
	}
	return Token{Value, word}, nil
}

func (t *Tokenizer) Next() (*Token, error) {
	if t.peeked != nil {
		token := t.peeked
		t.peeked = nil
		fmt.Println("NEXT RETURNING CACHED", token)
		return token, nil
	}

	token, err := t.NextToken()
	if err != nil {
		return nil, err
	}

	fmt.Println("NEXT RETURNING FRESH", token)

	return &token, nil
}

func (t *Tokenizer) Peek() (*Token, error) {
	if t.peeked != nil {
		fmt.Println("PEEK RETURNING CHACHED", t.peeked)
		return t.peeked, nil
	}

	token, err := t.Next()
	if err != nil {
		return nil, err
	}
	t.peeked = token
	fmt.Println("PEEK RETURNING FRESH", t.peeked)
	return token, err
}

func NewTokenizer(reader io.Reader) Tokenizer {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	return Tokenizer{scanner, nil}
}
