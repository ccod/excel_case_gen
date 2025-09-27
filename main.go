package main

import (
	"errors"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/xuri/excelize/v2"
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

type email struct {
	Email string `faker:"email"`
}

type studentConfig struct {
	gradeLower         int
	gradeUpper         int
	combineStudentName bool
	includeHomeroom    bool
	homeroomDispatch   map[string][]string
	includeRoomNumber  bool
	roomNumberDispatch map[string][]string
}

func pickUniqueRandomNumbers(count, min, max int) ([]int, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than 0")
	}
	if min >= max {
		return nil, fmt.Errorf("min must be less than max")
	}
	if count > (max - min + 1) {
		return nil, fmt.Errorf("cannot pick %d unique numbers from a range of %d numbers", count, (max - min + 1))
	}

	uniqueNumbers := make(map[int]bool)
	var result []int

	for len(result) < count {
		num := rand.Intn(max-min+1) + min

		if _, exists := uniqueNumbers[num]; !exists {
			uniqueNumbers[num] = true
			result = append(result, num)
		}
	}
	return result, nil
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

func genEmail() string {
	name := email{}
	err := faker.FakeData(&name)
	if err != nil {
		panic(err)
	}

	return name.Email
}

func genStudentBasic() studentBasic {
	student := studentBasic{}
	err := faker.FakeData(&student)
	if err != nil {
		panic(err)
	}

	return student
}

func appendEmail(table *[][]string) {
	for idx, row := range *table {
		row = append(row, genEmail())
		(*table)[idx] = row
	}
}

func addDuplicateAsEmail(count, emailIdx int, table *[][]string) {
	indices, err := pickUniqueRandomNumbers(count, 0, len(*table))
	if err != nil {
		panic(err)
	}

	for _, idx := range indices {
		row := make([]string, len((*table)[idx]))
		copy(row, (*table)[idx])

		row[emailIdx] = genEmail()
		*table = append(*table, row)
	}
}

func genTable(offsetID, studentNum int, config *studentConfig) [][]string {
	rows := make([][]string, studentNum)
	for i := range studentNum {
		row := []string{}

		id := strconv.Itoa(2000 + i + offsetID)
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

			idx := rand.Intn(len(val))
			row = append(row, val[idx])
		}

		if config.includeRoomNumber {
			if config.roomNumberDispatch == nil {
				config.roomNumberDispatch = make(map[string][]string)
			}

			val, ok := config.roomNumberDispatch[grade]
			if !ok {
				one := strconv.Itoa(200 + rand.Intn(20))
				two := strconv.Itoa(100 + rand.Intn(20))
				val = []string{one, two}
				config.roomNumberDispatch[grade] = val
			}

			idx := rand.Intn(len(val))
			row = append(row, val[idx])
		}
		rows[i] = row
	}

	return rows
}

func addA(i int) string {
	return "A" + strconv.Itoa(i)
}

func main() {
	config := studentConfig{
		includeHomeroom:   true,
		includeRoomNumber: true,
		gradeLower:        6,
		gradeUpper:        8,
	}

	headers := []string{"ID", "Studen Name", "Grade", "Teacher", "HomeRoom", "Email"}
	table := genTable(0, 50, &config)
	appendEmail(&table)
	addDuplicateAsEmail(15, 5, &table)

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	const sheetName = "Sheet1"

	f.SetSheetRow(sheetName, addA(1), &headers)
	for i, row := range table {
		err := f.SetSheetRow(sheetName, addA(i+2), &row)
		if err != nil {
			fmt.Println(err)
		}
	}

	if err := f.SaveAs("TestBook.xlsx"); err != nil {
		fmt.Println(err)
	}
}
