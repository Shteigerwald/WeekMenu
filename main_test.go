package main

import (
	"fmt"
	"log"
	"testing"
)

func TestCreateNextWeek(t *testing.T) {
	DBconnection()
	AllMeals, err := GetAllMeals()
	if err != nil {
		fmt.Printf("GetAllMeals Error")
	}
	last_week, err := GetLastWeek()
	if err != nil {
		log.Fatal(err)
	}
	for NewIdMeal := 0; NewIdMeal < 26; NewIdMeal++ {
		AllMeals[NewIdMeal].ID = 0
	}
	for NewIdMeal := 0; NewIdMeal < 7; NewIdMeal++ {
		last_week.Days[NewIdMeal].IdMeal.Int64 = 0
	}
	next_week, err := CreateNextWeek(last_week, AllMeals)
	fmt.Printf("next_week: %v\n", next_week)
	fmt.Printf("Error: %v\n", err)
	if err == nil {
		t.Errorf("Doesn't have error with zero last_week, AllMeals")
	}

}
func TestGetAllMeals_Right_ID(t *testing.T) {
	DBconnection()
	rightOrder := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26}
	AllMeals, err := GetAllMeals()
	if err != nil {
		log.Fatal(err)
	}
	for _, right_idmeal := range rightOrder {
		if right_idmeal != AllMeals[right_idmeal-1].ID {
			t.Errorf("GetAllMeals has wrong order")
		}
	}
}
