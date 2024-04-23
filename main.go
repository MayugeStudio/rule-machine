package main

import (
	"fmt"
	"os"
	"strings"
)

// ----------- Structs -----------

type Rule struct {
	name  string
	read  string
	write string
	step  string
	next  string
}

type Program struct {
	tokens []string
}

func (p *Program) Pop() string {
	s := p.tokens[0]
	p.tokens = p.tokens[1:]
	return s
}

func (p *Program) Count() int {
	return len(p.tokens)
}

type Tape struct {
	seq []string
}

func (t *Tape) Read(idx int) string {
	return t.seq[idx]
}

func (t *Tape) Write(idx int, s string) {
	t.seq[idx] = s
}

type Machine struct {
	rules   []*Rule
	current string
	tape    Tape
	head    int
	halt    bool
}

func (m *Machine) Next() {
	for _, rule := range m.rules {
		if rule.name == m.current && rule.read == m.tape.Read(m.head) {
			m.tape.Write(m.head, rule.write)
			switch rule.step {
			case "<-":
				m.head--
			case "->":
				m.head++
			case "<>": // Nothing have to do
			}
			m.current = rule.next
			m.halt = false
			break
		}
	}
}

func (m *Machine) Print() {
	fmt.Printf("%s\n", m.current)
	for _, t := range m.tape.seq {
		fmt.Printf("%s", t)
	}
	fmt.Printf("\n")
	fmt.Printf("%s^\n", strings.Repeat(" ", m.head))
	fmt.Printf("--------------------------------\n")
}

// ----------- Functions -----------

func ParseStep(name string) (string, error) {
	switch name {
	case "->", "<-", "<>":
		return name, nil
	default:
		return name, fmt.Errorf("expected <- or -> but got %s", name)
	}
}

func ParseRule(program *Program) (*Rule, error) {
	name := program.Pop()
	in := program.Pop()
	out := program.Pop()
	step, err := ParseStep(program.Pop())
	if err != nil {
		return nil, err
	}
	next := program.Pop()
	return &Rule{name, in, out, step, next}, nil
}

func ParseProgram(program Program) ([]*Rule, error) {
	var rules []*Rule
	for program.Count() > 0 {
		token := program.Pop()

		switch token {
		case "case":
			{
				rule, err := ParseRule(&program)
				if err != nil {
					return nil, err
				}
				rules = append(rules, rule)
			}
		default:
			{
				return nil, fmt.Errorf("unknown statement %s", token)
			}
		}
	}
	if len(rules) == 0 {
		return nil, fmt.Errorf("at least one rule has to be presant")
	}

	return rules, nil
}

// ----------- Main Function -----------

func main() {
	argv := os.Args[1:]
	if len(argv) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: expected 1 arguments but got %d\n", len(argv))
		os.Exit(1)
	}

	rule_file_name := argv[0]

	tokens, err := ReadFile(rule_file_name)
	program := Program{tokens}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not read file: %s\n", err)
		os.Exit(1)
	}
	rules, err := ParseProgram(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	machine := Machine{
		rules:   rules,
		current: "Entry",
		tape: Tape{
			seq: []string{"$", "0", "0", "0", "0", "1"},
		},
		head: 0,
		halt: false,
	}

	for !machine.halt {
		machine.Print()
		machine.halt = true
		machine.Next()
	}
}
