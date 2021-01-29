package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	osUs "os/user"
)

var (
	filePath       = getHomePath()
	userDataFile   = ".userData.json"
	fullPath       = fmt.Sprintf("%s/%s", filePath, userDataFile)
	ErrNoUserFound = errors.New("no user profile found")
	ErrUserExists  = errors.New("profile already exists")
)

func getHomePath() string {
	osUser, err := osUs.Current()
	if err != nil {
		log.Fatal(err)
	}

	return osUser.HomeDir
}

type User struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}

type Users map[string]User

func SaveUser(name string, user User, overwrite bool) (string, User, error) {
	dataBytes, err := openOrCreateFile()
	if err != nil {
		return name, user, err
	}

	users := Users{}

	jsonErr := json.Unmarshal(dataBytes, &users)
	if jsonErr != nil {
		return name, user, jsonErr
	}

	if overwrite {
		users[name] = user
	} else {
		for u := range users {
			if u == name {
				return name, user, ErrUserExists
			}
		}
		users[name] = user
	}

	// fmt.Println("data", data)

	dataWBytes, mErr := json.MarshalIndent(users, "", "    ")
	if mErr != nil {
		return name, user, mErr
	}

	err = ioutil.WriteFile(fullPath, dataWBytes, 0644)
	if err != nil {
		return name, user, mErr
	}

	return name, user, nil
}

func GetUser(name string) (User, error) {
	var (
		user  User
		users Users
	)

	db, err := openOrCreateFile()
	if err != nil {
		return user, err
	}

	jsonErr := json.Unmarshal(db, &users)
	if jsonErr != nil {
		return user, jsonErr
	}

	for n, uD := range users {
		if n == name {
			return uD, nil
		}
	}

	return user, ErrNoUserFound
}

func openOrCreateFile() ([]byte, error) {
	var dataBytes []byte
	_, err := os.Stat(fullPath)

	if err != nil {
		file, fileErr := os.Create(fullPath)
		if fileErr != nil {
			// log.Error("error creating file: %w", fileErr)
			return dataBytes, fileErr
		}
		defer file.Close()

		users := Users{}
		usersBytes, jsonErr := json.Marshal(users)
		if jsonErr != nil {
			return nil, jsonErr
		}
		wErr := ioutil.WriteFile(fullPath, usersBytes, 0644)
		if wErr != nil {
			return nil, wErr
		}
	}

	dataBytes, err = ioutil.ReadFile(fullPath)
	if err != nil {
		return dataBytes, err
	}
	return dataBytes, nil
}
