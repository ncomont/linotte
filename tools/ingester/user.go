package ingester

import (
	"fmt"

	"git.ve.home/nicolasc/linotte/libs/helpers"
	"git.ve.home/nicolasc/linotte/models"
	"git.ve.home/nicolasc/linotte/tools/console"
)

type UserIngester struct {
	userAccessor *models.UserAccessor
}

var userIngester *UserIngester = nil

func InitializeUserIngester(
	userAccessor *models.UserAccessor) *UserIngester {

	if userIngester == nil {
		userIngester = &UserIngester{userAccessor}
	}

	return userIngester
}

func (ingester *UserIngester) Ingest() error {
	fmt.Println("Adding users ... ")
	console, err := console.Initialize()
	helpers.HandleError(err)

	if user, err := console.Users.Create("clementf", "50LciYbq01F5"); err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	} else {
		fmt.Printf("User created: %v\n", user)
	}

	if user, err := console.Users.Create("vincentl", "iRF2g1878112"); err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	} else {
		fmt.Printf("User created: %v\n", user)
	}

	if user, err := console.Users.Create("helenel", "A6Z62tTd34h8"); err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	} else {
		fmt.Printf("User created: %v\n", user)
	}

	if user, err := console.Users.Create("nicolasc", "s87kX2bdMH72"); err != nil {
		fmt.Printf("Error creating user: %v\n", err)
	} else {
		fmt.Printf("User created: %v\n", user)
	}
	fmt.Println("Users added.")

	return nil
}
