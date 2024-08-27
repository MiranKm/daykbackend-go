package main

import (
	"daykbackend/third_party"
	"fmt"
	"os"
	"strconv"
)

func main() {

	fmt.Println("Starting app")
	third_party.InitFirebase()

	//size, err := third_party.GetRegisteredUsersSize()
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

}
