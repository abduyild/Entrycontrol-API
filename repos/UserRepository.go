package repos

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	FirstName string `bson:"FirstName" json:"FirstName"`
	LastName  string `bson:"LastName" json:"LastName"`
	Phone     string `bson:"Phone" json:"Phone"`
	Address   string `bson:"Address" json:"Address"`
	Time      string `bson:"Time" json:"Time"`
	Location  string `bson:"Location" json:"Location"`
}

var clientOptions *options.ClientOptions
var dbclient *mongo.Client

func InitDB() error {
	// Define Address of Database
	clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
	// Try to connect to Database, save error if one is thrown
	client, err := mongo.Connect(context.TODO(), clientOptions)
	// If there was an error connecting to the DB (DB not running, wrong URI, ...) return the error
	if err != nil {
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	dbclient = client
	return nil
}

func getDB(dbname string) *mongo.Database {
	return dbclient.Database(dbname)
}

func DoesDBExist(mosqueid string) bool {
	names, err := dbclient.ListDatabaseNames(context.TODO(), bson.D{{}})
	if err != nil {
		return false
	}
	for _, name := range names {
		if name == mosqueid {
			return true
		}
	}
	return false
}

func GetEntriesForDate(mosqueid string, date string) ([]User, error) {
	collection := getDB(mosqueid).Collection(date)
	var user User
	var users []User
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		cur.Decode(&user)
		users = append(users, user)
	}
	return users, nil
}

func PushToDB(mosque string, user User) {
	db := getDB(mosque)
	db.Collection(GetCurrentDate()).InsertOne(context.TODO(), user)
}

func GetCurrentDate() string {
	currentTime := time.Now()
	day := fmt.Sprintf("%02d", currentTime.Day())
	month := fmt.Sprintf("%02d", int(currentTime.Month()))
	return day + "-" + month + "-" + strconv.Itoa(currentTime.Year())
}

func StringToUser(firstName string, lastName string, phone string, address string, time string, location string) User {
	user := User{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Address:   address,
		Time:      time,
		Location:  location,
	}
	return user
}
