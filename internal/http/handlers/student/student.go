package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/vinit-jpl/students-api-go/internal/storage"
	"github.com/vinit-jpl/students-api-go/internal/types"
	"github.com/vinit-jpl/students-api-go/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	slog.Info("Creating a student")
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) { // handeling error for empty body
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil { // catching all errrors apart from empty body error
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors) // type assertion to get the validation errors
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		} // this will return an error if the struct is not valid

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}
