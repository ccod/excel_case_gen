package main

import (
	"errors"
	"fmt"
	"github.com/go-faker/faker/v4"
	"math/rand"
	"strconv"
)

/*
	A Student will have one ID, sometimes in multiple grades at once, have more than one email and phone number,
	possibly more than one homeroom teacher or Room

	Some of this data can be split among multiple rows
*/

type studentBasic struct {
	ID          string
	FirstName   string `faker:"first_name_male"`
	LastName    string `faker:"first_name_male"`
	Grades      []string
	Emails      []string
	HomeRooms   []nameGroup
	PhoneNumber string `faker:"phone_number"`
}

type nameGroup struct {
	FirstNameMale string `faker:"first_name_male"`
	LastName      string `faker:"last_name"`
}

type emailGroup struct {
	Emails []string `faker:"email slice_len=2"`
}

type studentConfig struct {
	gradeLower         int
	gradeUpper         int
	combineStudentName bool
	includeHomeroom    bool
	homeroomDispatch   map[string][]string
	includeRoomNumber  bool
	roomNumber         map[string][]string
}

func gradeRange(lower int, higher int) (string, error) {
	delta := (higher + 1) - lower
	if delta < 0 {
		return "", errors.New("first arg must be the same or smaller than the second arg")
	}

	num := (rand.Int() % delta) + lower
	if num < 1 {
		switch num {
		case 0:
			return "KG", nil
		case -1:
			return "PK", nil
		default:
			return "", errors.New("Don't have grade equivelent for that")
		}
	}

	if num > 12 {
		return "", errors.New("I don't have anything beyond seniors")
	}

	return strconv.Itoa(num), nil
}

func genName() string {
	name := nameGroup{}
	err := faker.FakeData(&name)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s, %s", name.FirstNameMale, name.LastName)
}

func genEmails() []string {
	names := emailGroup{}
	err := faker.FakeData(&names)
	if err != nil {
		panic(err)
	}

	return names.Emails
}

func genStudentBasic() studentBasic {
	student := studentBasic{}
	err := faker.FakeData(&student)
	if err != nil {
		panic(err)
	}

	return student
}

func genTable(studentNum int, config studentConfig) [][]string {
	rows := make([][]string, studentNum)
	for i := range studentNum {
		row := []string{}

		id := strconv.Itoa(i)
		row = append(row, id)

		// studentName
		row = append(row, genName())

		grade, err := gradeRange(config.gradeLower, config.gradeUpper)
		if err != nil {
			panic(err)
		}
		row = append(row, grade)

		if config.includeHomeroom {
			if config.homeroomDispatch == nil {
				config.homeroomDispatch = make(map[string][]string)
			}

			val, ok := config.homeroomDispatch[grade]
			if !ok {
				val = []string{genName(), genName()}
				config.homeroomDispatch[grade] = val
			}

			idx := rand.Int() % len(val)
			row = append(row, val[idx])
		}
		rows[i] = row
	}

	return rows
}

func main() {
	config := studentConfig{
		includeHomeroom: true,
		gradeLower:      6,
		gradeUpper:      8,
	}

	table := genTable(10, config)
	for row := range table {
		fmt.Println(table[row])
	}

	// for range 10 {
	// 	fmt.Println(genStudentBasic())
	//
	// }
}
