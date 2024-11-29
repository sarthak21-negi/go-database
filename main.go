package main 

import (
	"fmt"
	"os"
	"encoding/json"
	"sync"
	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type (
	Logger interface{
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}
	Driver struct{
		mutex sync.Mutex
		mutexes map[string]*sync.Mutex
		dir string
		log Logger
	}
)

type Options struct{
	Logger
}

func New(){

}

type Address struct{
	City string
	State string
	Country string
	Pincode json.Number
}

type User struct{
	Name string 
	Contact  json.Number
	Age string
	Company string 
	Address Address
}

func main(){
	dir := "./"

	db, err := New(dir, nil)
	 if err != nil{
		fmt.Println("Error", err)
	 }
	 employees := []User{
		{"Himanshu", "9286457312", "22", "Jane-Street",Address{"New York", "New York State", "USA", "10001"}},
		{"Chinmoy", "8786444382", "22", "Apple",Address{"San Fransico", "California", "USA", "10251"}},
		{"Mohit", "9582457312", "22", "Walmart",Address{"Banglore", "Karnataka", "India", "410001"}},
		{"Sarthak", "9785207312", "21", "Microsoft",Address{"Redmond", "Washington", "USA", "10501"}},
		{"Uday", "8896542374", "22", "Morgan Stanley",Address{"Banglore", "Karnataka","India","410013"}},
		{"Akash", "9271257452", "24", "Qualcomm",Address{"New York", "New York State", "USA", "10001"}},
	 }
	 for _, value := range employees{
		db.Write("users", value.Name, User{
			Name: value.Name,
			Contact: value.Contact,
			Age: value.Age,
			Company: value.Company,
			Address: value.Address,
		})
	 }

	 records, err := db.ReadAll("users")
	 if err != nil {
		fmt.Println("Error", err)
	 }
	 fmt.Println(records)

	 allusers := []User{}

	 for _, f := range records{
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil{
			fmt.Println("Error", err)
		}
		allusers = append(allusers, employeeFound)
	 }
	 fmt.Println(allusers)

	 if err := db.Delete("user", "john"); err != nil{
		fmt.Println("Error", err)
	 }

	 if err := db.Delete("user", ""); err != nil{
		fmt.Println("Error", err)
	 }
}