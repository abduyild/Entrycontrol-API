package repos

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
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

type Mosque struct {
	Name     string `bson:"Name" json:"Name"`
	Location string `bson:"Location" json:"Location"`
}

var clientOptions *options.ClientOptions
var dbclient *mongo.Client

func InitDB() error {
	// Define Address of Database
	clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
	// Define connection context to have a timeout on connections
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Try to connect to Database, save error if one is thrown
	client, err := mongo.Connect(ctx, clientOptions)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	names, err := dbclient.ListDatabaseNames(ctx, bson.D{{}})
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		cur.Decode(&user)
		users = append(users, GetDecryptedUser(user, mosqueid))
	}
	return users, nil
}

func PushToDB(mosque string, user User) {
	db := getDB(mosque)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := db.Collection(GetCurrentDate()).InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
	}
}

func AddMosque(mosqueid string, mosque Mosque) {
	db := getDB(mosqueid)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := db.Collection(GetCurrentDate()).InsertOne(ctx, mosque)
	if err != nil {
		log.Println(err)
	}
}

func GetCurrentDate() string {
	currentTime := time.Now()
	day := fmt.Sprintf("%02d", currentTime.Day())
	month := fmt.Sprintf("%02d", int(currentTime.Month()))
	return day + "-" + month + "-" + strconv.Itoa(currentTime.Year())
}

func getTime() string {
	return time.Now().Format("15:04")
}

func StringToUser(firstName string, lastName string, phone string, address string, location string) User {
	user := User{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Address:   address,
		Time:      getTime(),
		Location:  location,
	}
	return user
}

func getHash(passphrase string) string {
	hasher := md5.New()
	hasher.Write([]byte(passphrase))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetEncryptedUser(user User, passphrase string) User {
	encryptedUser := user
	encryptedUser.FirstName = encrypt(user.FirstName, passphrase)
	encryptedUser.LastName = encrypt(user.LastName, passphrase)
	encryptedUser.Phone = encrypt(user.Phone, passphrase)
	encryptedUser.Address = encrypt(user.Address, passphrase)
	return encryptedUser
}

func GetDecryptedUser(user User, passphrase string) User {
	decryptedUser := user
	decryptedUser.FirstName = decrypt(user.FirstName, passphrase)
	decryptedUser.LastName = decrypt(user.LastName, passphrase)
	decryptedUser.Phone = decrypt(user.Phone, passphrase)
	decryptedUser.Address = decrypt(user.Address, passphrase)
	return decryptedUser
}

func encrypt(plaintext string, passphrase string) string {
	block, _ := aes.NewCipher([]byte(getHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := string(gcm.Seal(nonce, nonce, []byte(plaintext), nil))
	return ciphertext
}

func decrypt(ciphertext string, passphrase string) string {
	key := []byte(getHash(passphrase))
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, _ := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	return string(plaintext)
}
