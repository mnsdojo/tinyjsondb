package main

import (
	"fmt"
	"log"

	"github.com/mnsdojo/tinyjsondb/pkg/tinydb"
)

type User struct {
	Name  string
	Email string
	Age   int
}

func main() {
	// Create a new TinyDB with default 10-minute cache
	db, err := tinydb.NewTinyDB("users.db")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create users
	users := map[string]User{
		"user1": {Name: "Alice", Email: "alice@example.com", Age: 30},
		"user2": {Name: "Bob", Email: "bob@example.com", Age: 25},
	}

	// Insert users
	for key, user := range users {
		if err := db.Create(key, user); err != nil {
			log.Printf("Error creating user %s: %v", key, err)
		}
	}

	// Read a user
	user, err := db.Read("user1")
	if err != nil {
		log.Printf("Error reading user: %v", err)
	} else {
		fmt.Printf("Read user: %+v\n", user)
	}

	// Update with error handling
	updatedUser := User{Name: "Alice Smith", Email: "alice.smith@example.com", Age: 31}
	if err := db.Update("user1", updatedUser); err != nil {
		log.Printf("Error updating user: %v", err)
	} else {
		fmt.Println("User updated successfully")
	}

	// Read all users
	allUsers := db.ReadAll()
	fmt.Println("All users:")
	for key, user := range allUsers {
		fmt.Printf("%s: %+v\n", key, user)
	}

	// Delete a user
	if err := db.Delete("user2"); err != nil {
		log.Printf("Error deleting user: %v", err)
	}
}
