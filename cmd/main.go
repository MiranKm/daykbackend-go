package main

import (
	"daykbackend/third_party"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"sync"
)

// run with go run cmd/main.go 0 2024-08-24
func main() {

	fmt.Println("Starting app")
	err := godotenv.Load("./env/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wg := sync.WaitGroup{}

	third_party.InitFirebase()

	go third_party.GetRegisteredUsersSize(&wg)
	args := os.Args[1:]

	fmt.Printf("Limit:%s Date: %s\n", args[0], args[1])

	limitArg := args[0]
	limit, _ := strconv.Atoi(limitArg)

	users, err := third_party.GetNewlyRegisteredUsers(limit, args[1])

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Docs size: %d", len(users))

	for _, user := range users {
		fmt.Printf("User %s registered on %s \n", user.Name, user.CreatedAt.Format("2006-01-02"))
	}

	wg.Wait()
}
