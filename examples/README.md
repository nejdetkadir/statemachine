# Examples

## Traffic Light
A simple traffic light state machine example.

```go
package main

import (
    "fmt"
    "github.com/nejdetkadir/statemachine"
)

func main() {
    // Define states
    states := []string{"Red", "Green", "Yellow"}
    initialState := "Red"

    // Create state machine
    sm, err := statemachine.New(states, initialState)
	
    if err != nil {
        fmt.Println("Error creating state machine:", err)
		
        return
    }

    // Define events
    events := []statemachine.Event{
        {
            Name: "SwitchToGreen",
            From: []string{"Red"},
            To:   "Green",
            Before: func() {
                fmt.Println("Switching from Red to Green")
            },
            After: func() {
                fmt.Println("Switched to Green")
            },
        },
        {
            Name: "SwitchToYellow",
            From: []string{"Green"},
            To:   "Yellow",
            Before: func() {
                fmt.Println("Switching from Green to Yellow")
            },
            After: func() {
                fmt.Println("Switched to Yellow")
            },
        },
        {
            Name: "SwitchToRed",
            From: []string{"Yellow"},
            To:   "Red",
            Before: func() {
                fmt.Println("Switching from Yellow to Red")
            },
            After: func() {
                fmt.Println("Switched to Red")
            },
        },
    }

    // Register events
    err = sm.RegisterEvents(events)
    if err != nil {
        fmt.Println("Error registering events:", err)
		
        return
    }

    // Simulate the traffic light system
    fmt.Println("Initial State:", sm.CurrentState())
	
    sm.Fire("SwitchToGreen")
    sm.Fire("SwitchToYellow")
    sm.Fire("SwitchToRed")
	
    fmt.Println("Final State:", sm.CurrentState())
}
```

## Order Processing System
A simple order processing system state machine example.

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nejdetkadir/statemachine"
)

type Order struct {
	ID          int
	Description string
	State       string
}

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=username dbname=mydb sslmode=disable password=mypassword"
	db, err = sql.Open("postgres", connStr)
	
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func getOrder(id int) (*Order, error) {
	order := &Order{}
	row := db.QueryRow("SELECT id, description, state FROM orders WHERE id = $1", id)
	
	if err := row.Scan(&order.ID, &order.Description, &order.State); err != nil {
		return nil, err
	}
	
	return order, nil
}

func updateOrderState(id int, state string) error {
	_, err := db.Exec("UPDATE orders SET state = $1 WHERE id = $2", state, id)
	return err
}

func main() {
	initDB()
	defer db.Close()

	// Fetch order from database
	orderID := 1
	order, err := getOrder(orderID)
	
	if err != nil {
		log.Fatal("Error fetching order:", err)
	}

	fmt.Println("Current State:", order.State)

	// Define states and initial state
	states := []string{"New", "Processing", "Shipped", "Delivered"}
	initialState := order.State

	// Create state machine
	sm, err := statemachine.New(states, initialState)
	
	if err != nil {
		fmt.Println("Error creating state machine:", err)
		
		return
	}

	// Define events
	events := []statemachine.Event{
		{
			Name: "ProcessOrder",
			From: []string{"New"},
			To:   "Processing",
			Before: func() {
				fmt.Println("Processing the order")
			},
			After: func() {
				fmt.Println("Order is now Processing")
			},
		},
		{
			Name: "ShipOrder",
			From: []string{"Processing"},
			To:   "Shipped",
			Before: func() {
				fmt.Println("Shipping the order")
			},
			After: func() {
				fmt.Println("Order is now Shipped")
			},
		},
		{
			Name: "DeliverOrder",
			From: []string{"Shipped"},
			To:   "Delivered",
			Before: func() {
				fmt.Println("Delivering the order")
			},
			After: func() {
				fmt.Println("Order is now Delivered")
			},
		},
	}

	// Register events
	err = sm.RegisterEvents(events)
	
	if err != nil {
		fmt.Println("Error registering events:", err)
		
		return
	}

	// Set AfterAll callback to update the order state in the database
	sm.AfterAll(func(event string, from string, to string) {
		fmt.Printf("Event '%s' transitioned from '%s' to '%s'\n", event, from, to)
		
		err := updateOrderState(orderID, to)
		
		if err != nil {
			fmt.Printf("Error updating order state in database: %v\n", err)
		} else {
			fmt.Printf("Order state updated to '%s' in database\n", to)
		}
	})

	// Simulate the order processing system
	fmt.Println("Initial State:", sm.CurrentState())
	
	sm.Fire("ProcessOrder")
	sm.Fire("ShipOrder")
	sm.Fire("DeliverOrder")
	
	fmt.Println("Final State:", sm.CurrentState())
}
```