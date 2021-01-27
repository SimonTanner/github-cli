package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/prometheus/common/log"
)

var (
	userDataFile = "userData.json"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Users map[string]User

func SaveUser(name string, user User) (string, User, error) {
	dataBytes, err := openOrCreateFile()
	if err != nil {
		return name, user, err
	}

	data := Users{}

	jsonErr := json.Unmarshal(dataBytes, &data)
	if jsonErr != nil {
		return name, user, jsonErr
	}

	fmt.Println("data length", len(data))

	if len(data) == 0 {
		data[name] = user
		// return name, user, nil
	}

	// for u, uD := range data {
	// 	if u == name {
	// 		return name, user, fmt.Errorf("profile %s already exists, user.name=%s, user.email=%s", name, uD.Name, uD.Email)
	// 	}

	// 	data[name] = user
	// }

	fmt.Println("data", data)

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
	fmt.Println("user:", users)

	for n, uD := range users {
		if n == name {
			return uD, nil
		}
	}

	return user, fmt.Errorf("no user found with profile name %s", name)
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
