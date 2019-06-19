package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

type ApiError struct {
	Message string `json:"message"`
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func WriteError(w http.ResponseWriter, message string, err error) {
	log.Printf("%s (%v)", message, err)
	json.NewEncoder(w).Encode(&ApiError{message})
}

func FilterGormError(err error) error {
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}
