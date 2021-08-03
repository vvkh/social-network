package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
)

type requestMaker func(string, string) string

var (
	tests = map[string]requestMaker{
		"register": func(firstName string, lastName string) string {
			return makeRegisterForm(makeUser(firstName, lastName))
		},
		"search": func(firstName string, lastName string) string {
			return makeSearchByNameForm(firstName, lastName)
		},
	}
)

func main() {
	var namesFilePath string
	var requestsFilePath string
	var testName string

	flag.StringVar(&requestsFilePath, "output", "requests.txt", "path to output with generated requests")
	flag.StringVar(&namesFilePath, "names", "names.txt", "path to sample user names")
	flag.Func("test-name", "test name: (register, search)", func(value string) error {
		if _, ok := tests[value]; !ok {
			return errors.New("invalid test name, possible values are: register, search")
		}
		testName = value
		return nil
	})
	flag.Parse()

	err := run(namesFilePath, requestsFilePath, testName)
	if err != nil {
		fmt.Printf("error while generating data: %v\n", err)
		os.Exit(1)
	}
}

func run(namesFilePath string, requestsFilePath string, testName string) error {
	namesFile, err := os.Open(namesFilePath)
	if err != nil {
		return err
	}
	defer namesFile.Close()

	requestsFile, err := os.Create(requestsFilePath)
	if err != nil {
		return err
	}
	defer requestsFile.Close()

	makeRequest, ok := tests[testName]
	if !ok {
		return errors.New("invalid test name")
	}

	reader := bufio.NewReader(namesFile)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		fullName := strings.Split(string(line), " ")
		if len(fullName) != 2 {
			return errors.New("got line with number of words not equal to 2")
		}
		firstName, lastName := fullName[0], fullName[1]
		request := makeRequest(firstName, lastName)
		_, err = requestsFile.WriteString(request + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

type User struct {
	UserName  string
	Password  string
	FirstName string
	LastName  string
	Age       string
	Sex       string
	Location  string
	About     string
}

var (
	possibleSex       = []string{"male", "female"}
	possibleLocations = []string{"USA", "Russia", "UK", "Canada", "Germany", "Norway", "Poland"}
	minPossibleAge    = 18
	maxPossibleAge    = 100
)

func makeUser(firstName string, lastName string) User {
	return User{
		UserName:  strings.ToLower(firstName + lastName),
		Password:  strings.ToLower(firstName + lastName),
		Age:       fmt.Sprintf("%d", minPossibleAge+rand.Intn(maxPossibleAge-minPossibleAge)), //nolint:gosec
		FirstName: firstName,
		LastName:  lastName,
		Sex:       possibleSex[rand.Intn(len(possibleSex))],             //nolint:gosec
		Location:  possibleLocations[rand.Intn(len(possibleLocations))], //nolint:gosec
		About:     "",
	}
}

func makeRegisterForm(user User) string {
	form := url.Values{}
	form.Set("username", user.UserName)
	form.Set("password", user.Password)
	form.Set("first_name", user.FirstName)
	form.Set("last_name", user.LastName)
	form.Set("sex", user.Sex)
	form.Set("about", user.About)
	form.Set("age", user.Age)
	form.Set("location", user.Location)
	return form.Encode()
}

func makeSearchByNameForm(firstName string, lastName string) string {
	firstNamePrefix := firstName[:1+rand.Intn(len(firstName)-1)]
	lastNamePrefix := lastName[:1+rand.Intn(len(lastName)-1)]
	form := url.Values{}
	form.Set("first_name", firstNamePrefix)
	form.Set("last_name", lastNamePrefix)
	return form.Encode()
}
