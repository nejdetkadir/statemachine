package statemachine

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("should return error if initial state is not in states", func(t *testing.T) {
		_, err := New([]string{"A", "B"}, "C")

		assert.Error(t, err)
		assert.Equal(t, "initial state must be one of the states", err.Error())
	})

	t.Run("should return a new state machine", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)
		assert.NotNil(t, sm)
		assert.Equal(t, states, sm.Context().states)
		assert.Equal(t, initialState, sm.CurrentState())
	})
}

func TestStateMachine_CurrentState(t *testing.T) {
	t.Run("should return the current state", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)
		assert.Equal(t, initialState, sm.CurrentState())
	})

	t.Run("should return the current state after firing an event", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "B"})

		assert.NoError(t, err)

		err = sm.Fire("event")

		assert.NoError(t, err)
		assert.Equal(t, "B", sm.CurrentState())
	})
}

func TestStateMachine_Context(t *testing.T) {
	t.Run("should return the context", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)
		assert.Equal(t, sm.Context(), sm.(*Context))
		assert.Equal(t, states, sm.Context().states)
		assert.Equal(t, initialState, sm.Context().currentState)
		assert.Empty(t, sm.Context().events)
		assert.Nil(t, sm.Context().beforeAll)
		assert.Nil(t, sm.Context().afterAll)
	})
}

func TestStateMachine_RegisterEvent(t *testing.T) {
	t.Run("should register an event", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "B"})

		assert.NoError(t, err)
		assert.Len(t, sm.Context().events, 1)
		assert.Equal(t, "event", sm.Context().events[0].name)
		assert.Equal(t, []string{"A"}, sm.Context().events[0].from)
		assert.Equal(t, "B", sm.Context().events[0].to)
		assert.Nil(t, sm.Context().events[0].before)
		assert.Nil(t, sm.Context().events[0].after)
	})

	t.Run("should return error if to state is not in states", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "C"})

		assert.Error(t, err)
		assert.Equal(t, "to state must be one of: [A B]", err.Error())
	})

	t.Run("should return error if from states are not in states", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A", "C"}, to: "B"})

		assert.Error(t, err)
		assert.Equal(t, "from states must be one of: [A B]", err.Error())
	})

	t.Run("should return error if from and to states are the same", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "A"})

		assert.Error(t, err)
		assert.Equal(t, "from and to states cannot be the same", err.Error())
	})

	t.Run("should return error if event name is not unique", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "B"})

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "B"})

		assert.Error(t, err)
		assert.Equal(t, "event name must be unique", err.Error())
	})
}

func TestStateMachine_RegisterEvents(t *testing.T) {
	t.Run("should register multiple events", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvents([]Event{
			{
				name: "event1",
				from: []string{"A"},
				to:   "B",
			},
			{
				name: "event2",
				from: []string{"B"},
				to:   "A",
			},
		})

		assert.NoError(t, err)
		assert.Len(t, sm.Context().events, 2)
		assert.Equal(t, "event1", sm.Context().events[0].name)
		assert.Equal(t, "event2", sm.Context().events[1].name)
	})

	t.Run("should return error if one of the events is invalid", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvents([]Event{
			{
				name: "event1",
				from: []string{"A"},
				to:   "B",
			},
			{
				name: "event2",
				from: []string{"C"},
				to:   "A",
			},
		})

		assert.Error(t, err)
		assert.Equal(t, "from states must be one of: [A B]", err.Error())
		assert.Empty(t, sm.Context().events)
	})

	t.Run("should return error if one of the events is invalid and clear all events", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvents([]Event{
			{
				name: "event1",
				from: []string{"A"},
				to:   "B",
			},
			{
				name: "event2",
				from: []string{"C"},
				to:   "A",
			},
		})

		assert.Error(t, err)
		assert.Equal(t, "from states must be one of: [A B]", err.Error())
		assert.Empty(t, sm.Context().events)
	})
}

func TestStateMachine_BeforeAll(t *testing.T) {
	t.Run("should run before all function before firing an event", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		beforeAll := false

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		sm.BeforeAll(func(event string, from string, to string) {
			beforeAll = true
		})

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "B"})

		assert.NoError(t, err)

		err = sm.Fire("event")

		assert.NoError(t, err)
		assert.True(t, beforeAll)
	})
}

func TestStateMachine_AfterAll(t *testing.T) {
	t.Run("should run after all function after firing an event", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		afterAll := false

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		sm.AfterAll(func(event string, from string, to string) {
			afterAll = true
		})

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "B"})

		assert.NoError(t, err)

		err = sm.Fire("event")

		assert.NoError(t, err)
		assert.True(t, afterAll)
	})
}

func TestStateMachine_RenderGraph(t *testing.T) {
	t.Run("should render the graph", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "event", from: []string{"A"}, to: "B"})

		assert.NoError(t, err)

		sm.RenderGraph()
	})

	t.Run("should render the graph with multiple events", func(t *testing.T) {
		states := []string{"A", "B", "C"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvents([]Event{
			{
				name: "event1",
				from: []string{"A"},
				to:   "B",
			},
			{
				name: "event2",
				from: []string{"B"},
				to:   "C",
			},
		})

		assert.NoError(t, err)

		sm.RenderGraph()
	})
}

func TestStateMachine_Fire(t *testing.T) {
	t.Run("should return error if event name is not registered", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.Fire("test1")

		assert.Error(t, err)
		assert.Equal(t, "test1 event is not registered", err.Error())
	})

	t.Run("should return error if event is not allowed in the current state", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "test1", from: []string{"B"}, to: "A"})

		assert.NoError(t, err)

		err = sm.Fire("test1")

		assert.Error(t, err)
		assert.Equal(t, "cannot fire the test1 event from the A state", err.Error())
	})

	t.Run("should return error if validate function returns an error", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{
			name: "test1",
			from: []string{"A"},
			to:   "B",
			validate: func(from string, to string) error {
				return errors.New("test1 event is not allowed")
			},
		})

		assert.NoError(t, err)

		err = sm.Fire("test1")

		assert.Error(t, err)
		assert.Equal(t, "test1 event is not allowed", err.Error())
	})

	t.Run("should run before function before firing an event", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		before := false

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{
			name: "test1",
			from: []string{"A"},
			to:   "B",
			before: func() {
				before = true
			},
		})

		assert.NoError(t, err)

		err = sm.Fire("test1")

		assert.NoError(t, err)
		assert.True(t, before)
	})

	t.Run("should run after function after firing an event", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		after := false

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{
			name: "test1",
			from: []string{"A"},
			to:   "B",
			after: func() {
				after = true
			},
		})

		assert.NoError(t, err)

		err = sm.Fire("test1")

		assert.NoError(t, err)
		assert.True(t, after)
	})

	t.Run("should run before and after functions before and after firing an event", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"
		before := false
		after := false

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{
			name: "test1",
			from: []string{"A"},
			to:   "B",
			before: func() {
				before = true
			},
			after: func() {
				after = true
			},
		})

		assert.NoError(t, err)

		err = sm.Fire("test1")

		assert.NoError(t, err)
		assert.True(t, before)
		assert.True(t, after)
	})

	t.Run("should change the current state after firing an event", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.RegisterEvent(Event{name: "test1", from: []string{"A"}, to: "B"})

		assert.NoError(t, err)

		err = sm.Fire("test1")

		assert.NoError(t, err)
		assert.Equal(t, "B", sm.CurrentState())
	})
}

func TestStateMachine_SetCurrentState(t *testing.T) {
	t.Run("should return error if state is not in states", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.SetCurrentState("C")

		assert.Error(t, err)
		assert.Equal(t, "state must be one of: [A B]", err.Error())
	})

	t.Run("should set the current state", func(t *testing.T) {
		states := []string{"A", "B"}
		initialState := "A"

		sm, err := New(states, initialState)

		assert.NoError(t, err)

		err = sm.SetCurrentState("B")

		assert.NoError(t, err)
		assert.Equal(t, "B", sm.CurrentState())
	})
}
