package main

import (
	"bytes"
	"encoding/gob"
	"os"
	"plugin"
	"strings"
)

type M = map[string]interface{}

var (
	store *plugin.Plugin
	scripts []*plugin.Plugin
	ch *chan M
)

func Boot(st *plugin.Plugin, sc []*plugin.Plugin, c *chan M) {
	store = st
	scripts = sc
	ch = c

	gob.Register(M{})
	gob.Register([]interface{}{})
}

func BeforeScripts() []func(M) M {
	p := os.Getenv("LXBOT_COMMAND_PREFIX")

	return []func(M) M{
		func(msg M) M {
			text := msg["message"].(M)["text"].(string)
			if text == p+"help" {
				m, _ := deepCopy(msg)
				go showHelp(m)
			}
			return msg
		},
	}
}

func deepCopy(msg M) (M, error) {
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	d := gob.NewDecoder(&b)
	if err := e.Encode(msg); err != nil {
		return nil, err
	}
	r := map[string]interface{}{}
	if err := d.Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}

func showHelp(msg M) {
	helpTexts := make([]string, 0)
	helpTexts = append(helpTexts, "\n")
	for _, s := range scripts {
		if fn, err := s.Lookup("Help"); err == nil{
			help := fn.(func() string)()
			if !strings.HasSuffix(help, "\n") {
				help += "\n"
			}
			helpTexts = append(helpTexts, help)
		}
	}
	helpText := strings.Join(helpTexts, "\n")
	msg["mode"] = "reply"
	msg["message"].(M)["text"] = helpText
	*ch <- msg
}
