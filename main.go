package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"time"
)

type meal struct {
	ID      int64
	Title   string
	garnish string
	main    string
}

type Day struct {
	IdMeal sql.NullInt64
}

type week struct {
	ID         int64
	Days       [7]Day
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
	DBconnection()

	//GET_LAST_WEEK
	last_week, err := GetLastWeek()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Information_about_LAST_week: %v\n ", last_week)

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
	next_week, err := CreateNextWeek(last_week, AllMeals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n Next_week: ", next_week)

	//Add_NEXT_WEEK_TO_DB
	Date, err := time.Parse("2006-01-02", last_week.Date_start)
	if err != nil {
		log.Fatal(err)
	}
	Date_Next_week := Date.AddDate(0, 0, 7)
	var LastInsertId int64
	LastInsertId, err = AddWeekToDB(next_week, Date_Next_week)
	fmt.Printf("LastInsertID: %v\n", LastInsertId)

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
		lastweek.Days[0].IdMeal.Int64,
		lastweek.Days[1].IdMeal.Int64,
		lastweek.Days[2].IdMeal.Int64,
		lastweek.Days[3].IdMeal.Int64,
		lastweek.Days[4].IdMeal.Int64,
		lastweek.Days[5].IdMeal.Int64,
		lastweek.Days[6].IdMeal.Int64,
	)
	fmt.Printf("LastWeekMealsID: %v\n", LastWeekMealsID)

	//PART: Delete all index MealsIDLastWeak from AllMeals and create var Meals_ID_Without_Last_Week
	Meals_ID_Without_Last_Week := AllMealsID
	for _, lastweekmealID := range LastWeekMealsID {
		// If we don't have meal this weak
		if lastweekmealID == 0 {
			continue
		}
		ZeroingElementFromSlice(&Meals_ID_Without_Last_Week, lastweekmealID-1)
	}
	fmt.Printf("Meals_ID_Without_Last_Week: %v, len: %v\n", Meals_ID_Without_Last_Week, len(Meals_ID_Without_Last_Week))

	//PART: Shuffle Meals_ID_Without_Last_Week
	rand.Shuffle(len(Meals_ID_Without_Last_Week), func(i, j int) {
		Meals_ID_Without_Last_Week[i], Meals_ID_Without_Last_Week[j] = Meals_ID_Without_Last_Week[j], Meals_ID_Without_Last_Week[i]
	})
	fmt.Printf("Random Meals_ID_Without_Last_Week: %v\n", Meals_ID_Without_Last_Week)

	//PART: Filling the days for next week in var "nextweek"
	var Days = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for DayIndex := 0; DayIndex < 7; DayIndex++ {
		var MealIndex int64 = 0
		fmt.Printf("!!!Day #%v\n", DayIndex)
		fmt.Printf("Meals_ID_Without_Last_Week %v\n", Meals_ID_Without_Last_Week)
		for Meals_ID_Without_Last_Week[MealIndex] == 0 {
			MealIndex += 1
			if MealIndex == 26 {
				return nextweek, errors.New("BREAK: You don't have enough meals in database for rules, MealIndex = 26\n")
			}
		}
		if DayIndex == 1 {
			for allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].main == allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64-1].main ||
				allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].garnish == allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64-1].garnish {
				MealIndex += 1
				if MealIndex == 26 {
					return nextweek, errors.New("BREAK: You don't have enough meals in database for rules, MealIndex = 26\n")
				}
				for Meals_ID_Without_Last_Week[MealIndex] == 0 {
					MealIndex += 1
					if MealIndex == 26 {
						return nextweek, errors.New("BREAK: You don't have enough meals in database for rules, MealIndex = 26\n")
					}
				}
			}
			fmt.Printf("MEALINDEX: %v\n", MealIndex)
			fmt.Printf("Meals_ID_Without_Last_Week[MealIndex]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].main)
			fmt.Printf(" allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64].main: %v\n", allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64-1].main)
			fmt.Printf("Meals_ID_Without_Last_Week[MealIndex]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].garnish)
			fmt.Printf(" allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64].garnish: %v\n", allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64-1].garnish)
		} else if DayIndex > 1 {
			for allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].main == allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64-1].main ||
				allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].garnish == allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64-1].garnish ||
				allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].main == allmeals[nextweek.Days[DayIndex-2].IdMeal.Int64-1].main ||
				allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].garnish == allmeals[nextweek.Days[DayIndex-2].IdMeal.Int64-1].garnish {
				MealIndex += 1
				if MealIndex == 26 {
					return nextweek, errors.New("BREAK: You don't have enough meals in database for rules, MealIndex = 26\n")
				}
				for Meals_ID_Without_Last_Week[MealIndex] == 0 {
					MealIndex += 1
					if MealIndex == 26 {
						return nextweek, errors.New("BREAK: You don't have enough meals in database for rules, MealIndex = 26\n")
					}
				}
			}
			fmt.Printf("MEALINDEX: %v\n", MealIndex)
			fmt.Printf("Meals_ID_Without_Last_Week[MealIndex]].main: %v\n", allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].main)
			fmt.Printf(" allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64].main: %v\n", allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64-1].main)
			fmt.Printf("Meals_ID_Without_Last_Week[MealIndex]].garnish: %v\n", allmeals[Meals_ID_Without_Last_Week[MealIndex]-1].garnish)
			fmt.Printf(" allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64].garnish: %v\n", allmeals[nextweek.Days[DayIndex-1].IdMeal.Int64-1].garnish)
		}
		fmt.Printf("NextWeek %v Name потребовалось %v попыток для нахождения блюда\n", Days[DayIndex], MealIndex)
		nextweek.Days[DayIndex].IdMeal.Int64 = Meals_ID_Without_Last_Week[MealIndex]
		nextweek.Days[DayIndex].IdMeal.Valid = true
		ZeroingElementFromSlice(&Meals_ID_Without_Last_Week, MealIndex)
		fmt.Printf("NextWeek %vID: %v\n", Days[DayIndex], nextweek.Days[DayIndex].IdMeal.Int64)
		fmt.Printf("NextWeek %vName: %v\n", Days[DayIndex], allmeals[nextweek.Days[DayIndex].IdMeal.Int64-1].Title)
	}
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
	row := db.QueryRow("SELECT * FROM Week WHERE id = (SELECT MAX(id) FROM Week)")
	err := row.Scan(&lastweek.ID,
		&lastweek.Days[0].IdMeal,
		&lastweek.Days[1].IdMeal,
		&lastweek.Days[2].IdMeal,
		&lastweek.Days[3].IdMeal,
		&lastweek.Days[4].IdMeal,
		&lastweek.Days[5].IdMeal,
		&lastweek.Days[6].IdMeal,
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

func ZeroingElementFromSlice(slice *[]int64, index int64) {
	(*slice)[index] = 0
}

func AddWeekToDB(week week, Date_last_week time.Time) (int64, error) {
	result, err := db.Exec("INSERT INTO Week (Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Date_start)"+
		" VALUES (?, ?, ?, ?, ?, ?, ?, ?)", week.Days[0].IdMeal, week.Days[1].IdMeal, week.Days[2].IdMeal,
		week.Days[3].IdMeal, week.Days[4].IdMeal, week.Days[5].IdMeal, week.Days[6].IdMeal, Date_last_week)
	if err != nil {
		return 0, fmt.Errorf("AddWeekToDB: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddWeekToDB: %v", err)
	}
	return id, nil
}

func DBconnection() {
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
		log.Fatal(fmt.Sprintf("Error: %v ", err))
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(fmt.Sprintf("Error: %v and DBUSER: %v and DBPASS: %v ", pingErr, MysqlConfig.User, MysqlConfig.Passwd))
	}
	fmt.Println("Connected!")
}
