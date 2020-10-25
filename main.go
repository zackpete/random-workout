package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/regex"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
)

type (
	Plan struct {
		Definitions []Definition `@@*`
	}

	Definition struct {
		Name string `@Ident`
		Item Item   `@@`
	}

	Item struct {
		Choose   []Weight `"{" @@ { "|" @@ } "}" |`
		Sequence []Item   `"[" @@ { "," @@ } "]" |`
		Alias    string   `"(" @Ident ")"`
	}

	Weight struct {
		Name  string `@Ident ":"`
		Value string `@Ident`
	}
)

func init() {
	buf := make([]byte, 8)
	_, _ = crand.Read(buf)
	rand.Seed(int64(binary.BigEndian.Uint64(buf)))
}

const example = `
workout { run : 2 | swim : 1 }
run { tempo : 2 | LSD : 1 | interval: 5 }
interval [ (lap), (rest) ]
lap { 1/4 mile : 1 | 1/2 mile : 1 }
rest { 2 minutes : 2 | 3 minutes : 2 | 4 minutes : 1 }
`

const grammar = `
	Space = \s+
	Ident = [\w\d /\-\\\.]+
	Punct = [\{\}\[\]\(\)\|\:\,]
`

func main() {
	if file, err := ioutil.ReadFile("plan.txt"); err != nil {
		fmt.Println("could not read plan.txt")
	} else if err := parse(file); err != nil {
		fmt.Println("error: " + err.Error())
	}

	fmt.Scanln()
}

func parse(data []byte) error {
	var plan Plan

	parser := participle.MustBuild(
		&plan,
		participle.Lexer(lexer.Must(regex.New(grammar))),
		participle.Elide("Space"),
	)

	if err := parser.ParseBytes(data, &plan); err != nil {
		panic(err)
	}

	if len(plan.Definitions) == 0 {
		return errors.New("no definitions")
	}

	write(plan.Definitions[0].Name, 0)
	return resolve(&plan, &plan.Definitions[0].Item, 1)
}

func write(value string, indent int) {
	value = strings.TrimSpace(value)
	fmt.Printf("%s%s\n", strings.Repeat("  ", indent), value)
}

func resolve(root *Plan, current *Item, indent int) error {
	if current.Choose != nil {
		if next, err := choose(current.Choose); err != nil {
			return err
		} else {
			write(next, indent)
			for _, d := range root.Definitions {
				if strings.TrimSpace(d.Name) == next {
					resolve(root, &d.Item, indent+1)
					break
				}
			}
		}
	} else if current.Sequence != nil {
		for _, s := range current.Sequence {
			resolve(root, &s, indent)
		}
	} else {
		for _, d := range root.Definitions {
			if strings.TrimSpace(d.Name) == current.Alias {
				write(current.Alias, indent)
				resolve(root, &d.Item, indent+1)
				return nil
			}
		}

		write(current.Alias, indent)
	}

	return nil
}

func choose(weights []Weight) (string, error) {
	m := make(map[string]float64)
	for _, w := range weights {
		w.Value = strings.TrimSpace(w.Value)
		if value, err := strconv.ParseFloat(w.Value, 64); err != nil {
			return "", errors.New("invalid number: " + w.Value)
		} else {
			m[strings.TrimSpace(w.Name)] = value
		}
	}

	return pick(m), nil
}

func pick(m map[string]float64) string {
	type Item struct {
		Name   string
		Weight float64
	}

	var (
		items []Item
		sum   float64
	)

	for k, v := range m {
		sum += v
		items = append(items, Item{k, sum})
	}

	r := rand.Float64() * sum

	for _, i := range items {
		if r <= i.Weight {
			return i.Name
		}
	}

	panic("unreachable")
}
