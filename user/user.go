package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/prometheus/common/log"
)

var (
	userDataFile   = "userData.json"
	ErrNoUserFound = errors.New("no user profile found")
	ErrUserExists  = errors.New("profile already exists")
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Users map[string]User

func SaveUser(name string, user User, overwrite bool) (string, User, error) {
	dataBytes, err := openOrCreateFile()
	if err != nil {
		return name, user, err
	}

	data := Users{}

	jsonErr := json.Unmarshal(dataBytes, &data)
	if jsonErr != nil {
		return name, user, jsonErr
	}

	if overwrite {
		data[name] = user
	} else {
		for u := range data {
			if u == name {
				return name, user, ErrUserExists
			}
		}
		data[name] = user
	}

	// fmt.Println("data", data)

	dataWBytes, mErr := json.MarshalIndent(data, "", "    ")
	if mErr != nil {
		return name, user, mErr
	}

	err = ioutil.WriteFile(userDataFile, dataWBytes, 0644)
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
	fmt.Println("user:", fmt.Sprintf("%+v", user))

	for n, uD := range users {
		if n == name {
			return uD, nil
		}
	}

	return user, ErrNoUserFound
}

func openOrCreateFile() ([]byte, error) {
	var dataBytes []byte
	_, err := os.Stat(userDataFile)
	// fmt.Println(fInf)

	if err != nil {
		fmt.Println(err)
		var fileErr error
		_, fileErr = os.Create(userDataFile)
		if fileErr != nil {
			log.Errorf("error creating file: %w", fileErr)
			return dataBytes, fileErr
		}

		users := Users{}
		usersBytes, jsonErr := json.Marshal(users)
		if jsonErr != nil {
			fmt.Println("J BALLS")
			return nil, jsonErr
		}
		wErr := ioutil.WriteFile(userDataFile, usersBytes, 0644)
		if wErr != nil {
			fmt.Println("W BALLS")
			return nil, wErr
		}
	}

	dataBytes, err = ioutil.ReadFile(userDataFile)
	if err != nil {
		fmt.Println("BALLS")
		return dataBytes, err
	}
	return dataBytes, nil
}
