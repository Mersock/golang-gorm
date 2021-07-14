package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func main() {
	dsn := "root:root@tcp(mysql-db:3306)/knight?parseTime=true"
	dial := mysql.Open(dsn)
	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}

	//Create Table
	// db.AutoMigrate(Gender{}, Test{}, Customer{})
	// db.Migrator().CreateTable(Test{})
	// db.Migrator().CreateTable(Person{})
	// db.Migrator().CreateTable(OrderDetail{})
	// db.Migrator().CreateTable(Category{})
	// db.Migrator().CreateTable(Model{})
	// db.Migrator().CreateTable(Customer{})

	// CreateGender("Male")
	// CreateGender("xxx")
	// GetGenders()
	// GetGender(10)
	// GetGenderByName("Males")
	// UpdateGender2(4, "yyy")
	// DeleteGender(3)
	// CreateTest(0, "test 1")
	// CreateTest(0, "test 2")
	// CreateTest(0, "test 3")
	// DeleteTest(3)
	// GetTests()
	// CreateCustomer("nid", 2)
	// GetCustomer()
}

type SqlLogger struct {
	logger.Interface
}

//config column
type Test struct {
	gorm.Model
	Code uint   `gorm:"primaryKey;comment:Test Comment"`
	Name string `gorm:"column:myname;type:varchar(100);unique;default:Hello;not null"`
}

//change table name
func (t Test) TableName() string {
	return "mytest"
}

//struct for create table
type Customer struct {
	ID       uint
	Name     string
	Gender   Gender
	GenderID uint
}

//struct for create table
type Person struct {
	ID   uint
	Name string
}

//struct for create table
type OrderDetail struct {
	ID   uint
	Name string
}

//struct for create table
type Category struct {
	ID   uint
	Name string
}

type Gender struct {
	ID   uint
	Name string `gorm:"unique;size(10)"`
}

func CreateCustomer(name string, genderID uint) {
	customer := Customer{
		Name:     name,
		GenderID: genderID,
	}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}

func GetCustomer() {
	customers := []Customer{}
	//preload all fk
	// tx := db.Preload(clause.Associations).Find(&customers)
	//preload specific
	tx := db.Preload("Gender").Find(&customers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	for _, v := range customers {
		fmt.Printf("%v|%v|%v\n", v.ID, v.Name, v.Gender.Name)
	}
}

func CreateTest(code uint, name string) {
	test := Test{Code: code, Name: name}
	db.Create(&test)
}

func GetTests() {
	tests := []Test{}
	db.Find(&tests)
	for _, v := range tests {
		fmt.Printf("%v|%v\n", v.ID, v.Name)
	}
}

func DeleteTest(id uint) {
	//soft delete
	db.Delete(&Test{}, id)
	//delete from table
	// db.Unscoped().Delete(&Test{}, id)
}

func CreateGender(name string) {
	gender := Gender{
		Name: name,
	}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGenders() {
	genders := []Gender{}
	tx := db.Order("id").Find(&genders)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}

func GetGender(id uint) {
	gender := []Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGenderByName(name string) {
	gender := []Gender{}
	tx := db.Where("name", name).First(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
	}
	gender.Name = name
	tx = db.Save(gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
	}
	GetGender(id)
}

func UpdateGender2(id uint, name string) {
	gender := Gender{Name: name}
	// tx := db.Model(&Gender{}).Where("id = ?", id).Updates(gender)
	tx := db.Model(&Gender{}).Where("id = @myid", sql.Named("myid", id)).Updates(gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("Deleted")
	GetGender(id)
}

//config Logger
func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v \n==============================\n", sql)
}
