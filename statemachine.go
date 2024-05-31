package statemachine

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"slices"
)

type (
	Context struct {
		states       []string
		initialState string
		currentState string
		events       []Event
		beforeAll    func(event string, from string, to string)
		afterAll     func(event string, from string, to string)
	}
	Event struct {
		name     string
		to       string
		from     []string
		before   func()
		after    func()
		validate func(from string, to string) error
	}
	StateMachine interface {
		CurrentState() string
		Fire(event string) error
		RegisterEvent(event Event) error
		RegisterEvents(events []Event) error
		BeforeAll(before func(event string, from string, to string))
		AfterAll(after func(event string, from string, to string))
		RenderGraph()
		Context() *Context
		SetCurrentState(state string) error
	}
)

func New(states []string, initialState string) (StateMachine, error) {
	if slices.Contains(states, initialState) == false {
		return nil, errors.New("initial state must be one of the states")
	}

	return &Context{
		states:       states,
		initialState: initialState,
		currentState: initialState,
	}, nil
}

func (c *Context) CurrentState() string {
	return c.currentState
}

func (c *Context) Context() *Context {
	return c
}

func (c *Context) RegisterEvent(event Event) error {
	if slices.Contains(c.states, event.to) == false {
		return errors.New(fmt.Sprintf("to state must be one of: %v", c.states))
	}

	if slices.ContainsFunc(event.from, func(s string) bool {
		return slices.Contains(c.states, s) == false
	}) {
		return errors.New(fmt.Sprintf("from states must be one of: %v", c.states))
	}

	if slices.Contains(event.from, event.to) {
		return errors.New("from and to states cannot be the same")
	}

	if slices.ContainsFunc(c.events, func(e Event) bool {
		return e.name == event.name
	}) {
		return errors.New("event name must be unique")
	}

	c.events = append(c.events, event)

	return nil
}

func (c *Context) RegisterEvents(events []Event) error {
	var err error

	for _, e := range events {
		err = c.RegisterEvent(e)
	}

	if err != nil {
		c.events = []Event{}

		return err
	}

	return nil
}

func (c *Context) BeforeAll(before func(event string, from string, to string)) {
	c.beforeAll = before
}

func (c *Context) AfterAll(after func(event string, from string, to string)) {
	c.afterAll = after
}

func (c *Context) RenderGraph() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Event", "From", "To"})

	for _, e := range c.events {
		t.AppendRow([]interface{}{e.name, e.from, e.to})
	}

	t.Render()
}

func (c *Context) Fire(event string) error {
	var currentEvent *Event

	for _, e := range c.events {
		if e.name == event {
			currentEvent = &e
			break
		}
	}

	if currentEvent == nil {
		return errors.New(fmt.Sprintf("%s event is not registered", event))
	}

	if slices.Contains(currentEvent.from, c.currentState) == false {
		return errors.New(fmt.Sprintf("cannot fire the %s event from the %s state", currentEvent.name, c.currentState))
	}

	if c.beforeAll != nil {
		c.beforeAll(currentEvent.name, c.currentState, currentEvent.to)
	}

	if currentEvent.before != nil {
		currentEvent.before()
	}

	if currentEvent.validate != nil {
		err := currentEvent.validate(c.currentState, currentEvent.to)

		if err != nil {
			return err
		}
	}

	c.currentState = currentEvent.to

	if currentEvent.after != nil {
		currentEvent.after()
	}

	if c.afterAll != nil {
		c.afterAll(currentEvent.name, c.currentState, currentEvent.to)
	}

	return nil
}

func (c *Context) SetCurrentState(state string) error {
	if !slices.Contains(c.states, state) {
		return errors.New(fmt.Sprintf("state must be one of: %v", c.states))
	}

	c.currentState = state

	return nil
}
