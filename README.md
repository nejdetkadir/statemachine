![Build and test](https://github.com/nejdetkadir/statemachine/actions/workflows/main.yml/badge.svg?branch=main)
![Go Version](https://img.shields.io/badge/go_version-_1.22.2-007d9c.svg)


![cover](docs/cover.png)

# StateMachine

StateMachine is a lightweight, easy-to-use state machine library for Go. It allows you to define states and events to manage complex state transitions in a clean and organized way.

## Features
- Define states and an initial state
- Register events with transitions between states
- Execute before and after hooks for state transitions
- Validate state transitions
- Render state transition graph for visualization

## Installation
To install StateMachine, use go get:

```bash
$ go get github.com/nejdetkadir/statemachine
```

## Usage

### Creating a StateMachine
Creating a new state machine is simple. Define your states and initial state:

```go
package main

import (
    "fmt"
    "github.com/nejdetkadir/statemachine"
)

func main() {
    states := []string{"A", "B"}
    initialState := "A"
    
    sm, err := statemachine.New(states, initialState)
	
    if err != nil {
        fmt.Println("Error creating state machine:", err)
		
        return
    }

    fmt.Println("Initial State:", sm.CurrentState())
}
```

### Registering Events
Register events that define transitions between states:

```go
func main() {
    states := []string{"A", "B"}
    initialState := "A"
    
    sm, err := statemachine.New(states, initialState)
	
    if err != nil {
        fmt.Println("Error creating state machine:", err)
		
        return
    }

    event := statemachine.Event{
        Name: "event1",
        From: []string{"A"},
        To:   "B",
    }

    err = sm.RegisterEvent(event)
	
    if err != nil {
        fmt.Println("Error registering event:", err)
		
        return
    }

    fmt.Println("Event registered successfully")
}
```

### Firing Events
Fire events to transition between states:

```go
func main() {
    states := []string{"A", "B"}
    initialState := "A"
    
    sm, err := statemachine.New(states, initialState)
	
    if err != nil {
        fmt.Println("Error creating state machine:", err)
		
        return
    }

    event := statemachine.Event{
        Name: "event1",
        From: []string{"A"},
        To:   "B",
    }

    err = sm.RegisterEvent(event)
	
    if err != nil {
        fmt.Println("Error registering event:", err)
		
        return
    }

    err = sm.Fire("event1")
	
    if err != nil {
        fmt.Println("Error firing event:", err)
		
        return
    }

    fmt.Println("Current State:", sm.CurrentState())
}
```

### Hooks and Validation
You can define before, after, and validate hooks for each event:

```go
event := statemachine.Event{
    Name: "event1",
    From: []string{"A"},
    To:   "B",
    Before: func() {
        fmt.Println("Before transition")
    },
    After: func() {
        fmt.Println("After transition")
    },
    Validate: func(from string, to string) error {
        if from == "A" && to == "B" {
            return nil
        }
		
        return fmt.Errorf("invalid transition from %s to %s", from, to)
    },
}
```

### Rendering Graph
Render the state transition graph to visualize events:

```go
func main() {
    states := []string{"A", "B"}
    initialState := "A"
    
    sm, err := statemachine.New(states, initialState)
	
    if err != nil {
        fmt.Println("Error creating state machine:", err)
		
        return
    }

    event := statemachine.Event{
        Name: "event1",
        From: []string{"A"},
        To:   "B",
    }

    err = sm.RegisterEvent(event)
	
    if err != nil {
        fmt.Println("Error registering event:", err)
		
        return
    }

    sm.RenderGraph()
}
```

## Testing
To run the tests, use go test:

```bash
$ go test ./...
```

## Contributing
Bug reports and pull requests are welcome on GitHub at https://github.com/nejdetkadir/statemachine. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [code of conduct](https://github.com/nejdetkadir/statemachine/blob/main/CODE_OF_CONDUCT.md).

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
