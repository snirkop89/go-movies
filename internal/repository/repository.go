package repository

import (
	"database/sql"

	"github.com/snirkop89/go-movies/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB

	// Movies models
	AllMovies(genre ...int) ([]*models.Movie, error)
	OneMovie(id int) (*models.Movie, error)
	EditMovie(id int) (*models.Movie, []*models.Genre, error)
	InsertMovie(movie models.Movie) (int, error)
	UpdateMovie(movie models.Movie) error
	UpdateMovieGenres(id int, genresIDs []int) error
	DeleteMovie(id int) error

	AllGenres() ([]*models.Genre, error)

	// User models
	UserByEmail(email string) (*models.User, error)
	UserByID(id int) (*models.User, error)
}
