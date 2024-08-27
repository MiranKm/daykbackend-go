package third_party

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var ctx = context.Background()

var firebaseApp *firebase.App
var firestoreClient *firestore.Client

type User struct {
	Name      string
	CreatedAt time.Time
}

func NewUser(name string, createdAt time.Time) *User {
	return &User{
		Name:      name,
		CreatedAt: createdAt,
	}
}

func InitFirebase() {
	var err error
	var projectID = os.Getenv("PROJECT_ID")

	var conf = &firebase.Config{ProjectID: projectID}

	firebaseApp, err = firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	initializeFirestore()
}

func initializeFirestore() {
	var err error
	firestoreClient, err = firebaseApp.Firestore(ctx)

	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
		return
	}

}

func GetRegisteredUsersSize(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	if firestoreClient == nil {
		fmt.Println(errors.New("firebase Not initialized"))
		return
	}

	docs, err := firestoreClient.Collection("users").Documents(ctx).GetAll()

	if err != nil {
		fmt.Printf("can't fetch docs for users %v\n", err)
		return
	}

	fmt.Println(len(docs))
}

func GetNewlyRegisteredUsers(limit int, fromDate string) ([]User, error) {
	startTime, err := time.Parse(time.DateOnly, fromDate)

	if err != nil {
		return nil, err
	}

	var users []User

	if firestoreClient == nil {
		return users, errors.New("firebase Not initialized")
	}

	query := firestoreClient.Collection("users").Where("created_at", ">", startTime.String())

	if limit > 0 && limit < 100 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(10)
	}

	docs, err := query.Documents(ctx).GetAll()

	if err != nil {
		log.Fatalf("can't fetch docs for users %v\n", err)
		return users, err
	}

	fmt.Printf("Docs size: %d at %v\n", len(docs), startTime)

	for _, doc := range docs {
		var createdAt, err = time.Parse("yyyy-MM-d\\TH:mm:ss.yyyy", doc.Data()["created_at"].(string))
		if err != nil {
			createdAt = time.Now()
		}

		var user = NewUser(
			doc.Data()["full_name"].(string),
			createdAt,
		)

		users = append(users, *user)
	}

	return users, nil
}
