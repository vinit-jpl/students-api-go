package student

import (
	"log/slog"
	"net/http"

	"github.com/vinit-jpl/students-api-go/internal/types"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		slog.Info("Creating a student")
		w.Write([]byte("Welcome To Students API!"))
	}
}
