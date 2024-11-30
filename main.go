package main 

import (
	"fmt"
	"os"
	"encoding/json"
	"sync"
	"io/ioutil"
	"path/filepath"
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

func New(dir string, option *Options)(*Driver, error){
	dir = filepath.Clean(dir)

	optns := Options{}

	if option != nil{
		optns = *option
	}

	if optns.Logger == nil{
		optns.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir: dir,
		mutexes: make(map[string]*sync.Mutex),
		log: optns.Logger,
	}

	if _, err := os.Stat(dir); err == nil{
		optns.Logger.Debug("Using '%s' (database already exist)/n",dir)
		return &driver, nil
	}

	optns.Logger.Debug("Creating the database at '%s'..../n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d* Driver) Write(collection, resource string, v interface{}) error{

	if collection == ""{
		return fmt.Errorf("Missing collection - no place to save records!")
	}

	if resource == ""{
		return fmt.Errorf("Missing resource - unable to save records (no name)!")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil{
		return err
	}
	b = append(b, byte('\n'))

	if err := os.WriteFile(tmpPath, b , 0644); err != nil{
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

func (d* Driver) Read(collection, resource string, v interface{}) error{

	if collection == ""{
		return fmt.Errorf("Missing collection - no place to save records!")
	}

	if resource == ""{
		return fmt.Errorf("Missing resource - unable to save records (no name)!")
	}

	records := filepath.Join(d.dir, collection, resource)

	if _, err := stat(records); err != nil{
		return err
	}

	b, err := os.ReadFile(records +".json")

	if err != nil {
		return err
	}
	return json.Unmarshal(b, &v)
}

func (d* Driver) ReadAll()(){

}

func (d* Driver) Delete() error{

}

func (d* Driver) getOrCreateMutex(collection string)* sync.Mutex{
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok{
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

func stat(path string)(fi os.FileInfo, err error){
	if fi, err =  os.Stat(path); os.IsNotExist(err){
		fi, err = os.Stat(path + ".json") 
	}
	return fi, err
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

	 if err := db.Delete("user", "Akash"); err != nil{
		fmt.Println("Error", err)
	 }

	 if err := db.Delete("user", ""); err != nil{
		fmt.Println("Error", err)
	 }
}