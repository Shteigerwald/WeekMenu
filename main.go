package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"
)

type meal struct {
	ID      int64
	Title   string
	garnish string
	main    string
}

type week struct {
	ID         int64
	Sunday     sql.NullInt64
	Monday     sql.NullInt64
	Tuesday    sql.NullInt64
	Wednesday  sql.NullInt64
	Thursday   sql.NullInt64
	Friday     sql.NullInt64
	Saturday   sql.NullInt64
	Date_start string
}

// SQL PART
var db *sql.DB

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	// random
	rand.Seed(time.Now().UnixNano())
}

func main() {
	//Capture connection properties.
	MysqlConfig := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "menuapp",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", MysqlConfig.FormatDSN())
	if err != nil {
		log.Fatal("Error: %v ", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(fmt.Sprintf("Error: %v and DBUSER: %u and DBPASS: %p ", pingErr, MysqlConfig.User, MysqlConfig.Passwd))
	}
	fmt.Println("Connected!")

	//GET_LAST_WEEK
	last_week, err := GetLastWeek()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Information_about_week: %v ", last_week)

	//GET_ALL_MEALS_FROM_DATABASE
	AllMeals, err := GetAllMeals()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All meals id and title: \n")
	for _, meal := range AllMeals {
		fmt.Printf(" %v - %v, \n", meal.ID, meal.Title)
	}

	//GET_NEXT_WEEK
	AllMealID, err := CreateNextWeek(last_week, AllMeals)
	fmt.Print("\n ЗАГЛУШКА", AllMealID)

}
func CreateNextWeek(lastweek week, allmeals []meal) (week, error) {

	var nextweek week

	//PART: GetAllMealID from allmeals
	var AllMealsID []int64
	for _, Meal := range allmeals {
		MealID := Meal.ID
		AllMealsID = append(AllMealsID, MealID)
	}
	fmt.Printf("AllIdMeals: %v, len: %v\n", AllMealsID, len(AllMealsID))

	//PART: GetAllMealID from lastweek
	var LastWeekMealsID []int64
	LastWeekMealsID = append(
		LastWeekMealsID,
		lastweek.Sunday.Int64,
		lastweek.Monday.Int64,
		lastweek.Tuesday.Int64,
		lastweek.Wednesday.Int64,
		lastweek.Thursday.Int64,
		lastweek.Friday.Int64,
		lastweek.Saturday.Int64)
	fmt.Printf("LastWeekMealsID: %v\n", LastWeekMealsID)

	//PART: Delete all index MealsIDLastWeak from AllMeals and create var Meals_ID_Without_Last_Week
	Meals_ID_Without_Last_Week := AllMealsID
	for _, lastweekmealID := range LastWeekMealsID {
		if lastweekmealID == 0 {
			continue
		}
		Meals_ID_Without_Last_Week = RemoveElementFromSlice(Meals_ID_Without_Last_Week, lastweekmealID)
	}
	fmt.Printf("Meals_ID_Without_Last_Week: %v, len: %v\n", Meals_ID_Without_Last_Week, len(Meals_ID_Without_Last_Week))

	//PART: Shuffle Meals_ID_Without_Last_Week
	rand.Shuffle(len(Meals_ID_Without_Last_Week), func(i, j int) {
		Meals_ID_Without_Last_Week[i], Meals_ID_Without_Last_Week[j] = Meals_ID_Without_Last_Week[j], Meals_ID_Without_Last_Week[i]
	})
	fmt.Printf("Random Meals_ID_Without_Last_Week: %v\n", Meals_ID_Without_Last_Week)

	//PART: Filling the days for next week in var "nextweek"
	//Monday
	var indexmeal = int64(1)
	nextweek.Monday.Int64 = Meals_ID_Without_Last_Week[indexmeal]
	nextweek.Monday.Valid = true
	Meals_ID_Without_Last_Week = RemoveElementFromSlice(Meals_ID_Without_Last_Week, indexmeal)
	fmt.Printf("NextWeek MondayID: %v\n", nextweek.Monday.Int64)
	fmt.Printf("NextWeek MondayName: %v\n", allmeals[nextweek.Monday.Int64-1].Title)

	//Tuesday
	indexmeal = 0
	var Days = [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	for _, Day := range Days {
		reflectNextweek := reflect.ValueOf(nextweek).FieldByName(Day)
		fmt.Printf("reflectNextweek: %v", reflectNextweek)
		for allmeals[Meals_ID_Without_Last_Week[indexmeal]].main == allmeals[nextweek.Monday.Int64].main || allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish == allmeals[nextweek.Monday.Int64].garnish {
			fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].main)
			fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].main: %v\n", allmeals[nextweek.Monday.Int64].main)
			fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish)
			fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].garnish: %v\n", allmeals[nextweek.Monday.Int64].garnish)
			fmt.Printf("Tuesday indexmeal: %v\n", indexmeal)
			indexmeal += 1
		}
	}
	//
	//fmt.Println(val)
	//for allmeals[Meals_ID_Without_Last_Week[indexmeal]].main == allmeals[nextweek.Monday.Int64].main || allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish == allmeals[nextweek.Monday.Int64].garnish {
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].main)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].main: %v\n", allmeals[nextweek.Monday.Int64].main)
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].garnish: %v\n", allmeals[nextweek.Monday.Int64].garnish)
	//	fmt.Printf("Tuesday indexmeal: %v\n", indexmeal)
	//	indexmeal += 1
	//}
	//nextweek.Tuesday.Int64 = Meals_ID_Without_Last_Week[indexmeal]
	//nextweek.Tuesday.Valid = true
	//Meals_ID_Without_Last_Week = RemoveElementFromSlice(Meals_ID_Without_Last_Week, indexmeal)
	//fmt.Printf("NextWeek TuesdayID: %v\n", nextweek.Tuesday.Int64)
	//fmt.Printf("NextWeek TuesdayName: %v\n", allmeals[nextweek.Tuesday.Int64-indexmeal].Title)
	//
	////Wednesday
	//indexmeal = 0
	//for allmeals[Meals_ID_Without_Last_Week[indexmeal]].main == allmeals[nextweek.Tuesday.Int64].main || allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish == allmeals[nextweek.Tuesday.Int64].garnish {
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].main)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].main: %v\n", allmeals[nextweek.Tuesday.Int64].main)
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].garnish: %v\n", allmeals[nextweek.Tuesday.Int64].garnish)
	//	fmt.Printf("Tuesday indexmeal: %v\n", indexmeal)
	//	indexmeal += 1
	//}
	//nextweek.Wednesday.Int64 = Meals_ID_Without_Last_Week[indexmeal]
	//nextweek.Wednesday.Valid = true
	//Meals_ID_Without_Last_Week = RemoveElementFromSlice(Meals_ID_Without_Last_Week, indexmeal)
	//fmt.Printf("NextWeek WednesdayID: %v\n", nextweek.Wednesday.Int64)
	//fmt.Printf("NextWeek WednesdayName: %v\n", allmeals[nextweek.Wednesday.Int64-indexmeal].Title)
	//
	////Thursday
	//indexmeal = 0
	//for allmeals[Meals_ID_Without_Last_Week[indexmeal]].main == allmeals[nextweek.Wednesday.Int64].main || allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish == allmeals[nextweek.Wednesday.Int64].garnish {
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].main)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].main: %v\n", allmeals[nextweek.Wednesday.Int64].main)
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].garnish: %v\n", allmeals[nextweek.Wednesday.Int64].garnish)
	//	fmt.Printf("Tuesday indexmeal: %v\n", indexmeal)
	//	indexmeal += 1
	//}
	//nextweek.Thursday.Int64 = Meals_ID_Without_Last_Week[1]
	//nextweek.Thursday.Valid = true
	//Meals_ID_Without_Last_Week = RemoveElementFromSlice(Meals_ID_Without_Last_Week, 1)
	//fmt.Printf("NextWeek ThursdayID: %v\n", nextweek.Thursday.Int64)
	//fmt.Printf("NextWeek ThursdayName: %v\n", allmeals[nextweek.Thursday.Int64-1].Title)
	//
	////Friday
	//indexmeal = 0
	//for allmeals[Meals_ID_Without_Last_Week[indexmeal]].main == allmeals[nextweek.Thursday.Int64].main || allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish == allmeals[nextweek.Thursday.Int64].garnish {
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].main)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].main: %v\n", allmeals[nextweek.Thursday.Int64].main)
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].garnish: %v\n", allmeals[nextweek.Thursday.Int64].garnish)
	//	fmt.Printf("Tuesday indexmeal: %v\n", indexmeal)
	//	indexmeal += 1
	//}
	//nextweek.Friday.Int64 = Meals_ID_Without_Last_Week[1]
	//nextweek.Friday.Valid = true
	//Meals_ID_Without_Last_Week = RemoveElementFromSlice(Meals_ID_Without_Last_Week, 1)
	//fmt.Printf("NextWeek FridayID: %v\n", nextweek.Friday.Int64)
	//fmt.Printf("NextWeek FridayName: %v\n", allmeals[nextweek.Friday.Int64-1].Title)
	//
	////Saturday
	//indexmeal = 0
	//for allmeals[Meals_ID_Without_Last_Week[indexmeal]].main == allmeals[nextweek.Friday.Int64].main || allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish == allmeals[nextweek.Friday.Int64].garnish {
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].main)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].main: %v\n", allmeals[nextweek.Friday.Int64].main)
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].garnish: %v\n", allmeals[nextweek.Friday.Int64].garnish)
	//	fmt.Printf("Tuesday indexmeal: %v\n", indexmeal)
	//	indexmeal += 1
	//}
	//nextweek.Saturday.Int64 = Meals_ID_Without_Last_Week[1]
	//nextweek.Saturday.Valid = true
	//Meals_ID_Without_Last_Week = RemoveElementFromSlice(Meals_ID_Without_Last_Week, 1)
	//fmt.Printf("NextWeek SaturdayID: %v\n", nextweek.Saturday.Int64)
	//fmt.Printf("NextWeek SaturdayName: %v\n", allmeals[nextweek.Saturday.Int64-1].Title)
	//
	////Sunday
	//indexmeal = 0
	//for allmeals[Meals_ID_Without_Last_Week[indexmeal]].main == allmeals[nextweek.Saturday.Int64].main || allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish == allmeals[nextweek.Saturday.Int64].garnish {
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].main)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].main: %v\n", allmeals[nextweek.Saturday.Int64].main)
	//	fmt.Printf("Tuesday allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[indexmeal]].garnish)
	//	fmt.Printf("Tuesday allmeals[nextweek.Monday.Int64].garnish: %v\n", allmeals[nextweek.Saturday.Int64].garnish)
	//	fmt.Printf("Tuesday indexmeal: %v\n", indexmeal)
	//	indexmeal += 1
	//}
	//nextweek.Sunday.Int64 = Meals_ID_Without_Last_Week[1]
	//nextweek.Sunday.Valid = true
	//Meals_ID_Without_Last_Week = RemoveElementFromSlice(Meals_ID_Without_Last_Week, 1)
	//fmt.Printf("NextWeek SundayID: %v\n", nextweek.Sunday.Int64)
	//fmt.Printf("NextWeek SundayName: %v\n", allmeals[nextweek.Sunday.Int64-1].Title)

	return nextweek, nil
}

// albumsByArtist queries for albums that have the specified artist name.
func GetAllMeals() ([]meal, error) {
	// An albums slice to hold data from returned rows.
	var meals []meal

	rows, err := db.Query("SELECT * FROM Meals")
	if err != nil {
		return nil, fmt.Errorf("Meals error: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var meal meal
		if err := rows.Scan(&meal.ID, &meal.Title, &meal.garnish, &meal.main); err != nil {
			return nil, fmt.Errorf("GetAllMeals error: %v", err)
		}
		meals = append(meals, meal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllMeals error: %v", err)
	}
	return meals, nil
}

func GetLastWeek() (week, error) {
	var lastweek week
	row := db.QueryRow("SELECT * FROM Week WHERE id = 1")
	err := row.Scan(&lastweek.ID,
		&lastweek.Sunday,
		&lastweek.Monday,
		&lastweek.Tuesday,
		&lastweek.Wednesday,
		&lastweek.Thursday,
		&lastweek.Friday,
		&lastweek.Saturday,
		&lastweek.Date_start,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return lastweek, fmt.Errorf("GetLastWeek: no such Week")
		}
		return lastweek, fmt.Errorf("GetLastWeek error:  %v", err)
	}
	return lastweek, nil
}

func RemoveElementFromSlice(slice []int64, index int64) []int64 {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
