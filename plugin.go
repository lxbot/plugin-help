package main

import (
	"encoding/gob"
	"github.com/lxbot/lxlib"
	"log"
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
			m, err := lxlib.NewLXMessage(msg)
			if err != nil {
				log.Println(err)
				return nil
			}
			text := m.Message.Text
			if text == p+"help" {
				go showHelp(m)
			}
			return msg
		},
	}
}

func showHelp(m *lxlib.LXMessage) {
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
	msg, err := m.SetText(helpText).Reply().ToMap()
	if err != nil {
		log.Println(err)
		return
	}

	*ch <- msg
}
