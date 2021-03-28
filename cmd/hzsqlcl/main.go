package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	hz "github.com/hazelcast/hazelcast-go-client/v4/hazelcast"
	"log"
	"os"
	"strings"
)

type sqlCl struct {
	client *hz.Client
}

func NewSqlCl(addr string) (*sqlCl, error) {
	cb := hz.NewClientConfigBuilder()
	cb.Network().SetAddrs(addr)
	if config, err := cb.Config(); err != nil {
		return nil, err
	} else if client, err := hz.StartNewClientWithConfig(config); err != nil {
		return nil, err
	} else {
		return &sqlCl{client}, nil
	}
}

func (cl *sqlCl) completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	if strings.Index(d.Text, "select") >= 0 {
		return []prompt.Suggest{
			{Text: "*"},
			{Text: "from...", Description: "Choose the map to query"},
		}
	}
	s := []prompt.Suggest{
		{Text: "select"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (cl *sqlCl) executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}
	if result, err := cl.client.ExecuteSQL(s); err != nil {
		log.Println(err.Error())
	} else {
		fmt.Println("Result:", result)
	}
}

func main() {
	cl, err := NewSqlCl("localhost:5701")
	if err != nil {
		log.Fatal(err)
	}
	p := prompt.New(
		cl.executor,
		cl.completer,
		prompt.OptionTitle("HzSQLCl: Interactive Hazelcast SQL client"),
		prompt.OptionPrefix("SQL> "),
		//prompt.OptionInputTextColor(prompt.Yellow),
		//prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}
