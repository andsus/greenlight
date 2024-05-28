package data

import (
	"time"

	"github.com/andsus/greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`                // Unique integer ID for the movie
	CreatedAt time.Time `json:"-"`                 // Timestamp for when the movie is added to our database
	Title     string    `json:"title"`             // Movie title
	Year      int32     `json:"year,omitempty"`    // movie release year
	Runtime   Runtime   `json:"runtime,omitempty"` // Movie runtime (in minutes)
	Genres    []string  `json:"generes,omitempty"` // Slice of genres for the movie (romance, comedy, etc.)
	Version   int32     `json:"version"`           // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

/*
	func (m Movie) MarshalJSON() ([]byte, error) {
		// Create a variable holding the custom runtime string, just like before.
		var runtime string
		if m.Runtime != 0 {
			runtime = fmt.Sprintf("%d mins", m.Runtime)
		}
		// Define a MovieAlias type which has the underlying type Movie. Due to the way that
		// Go handles type definitions (https://golang.org/ref/spec#Type_definitions) the
		// MovieAlias type will contain all the fields that our Movie struct has but,
		// importantly, none of the methods.
		type MovieAlias Movie
		// Embed the MovieAlias type inside the anonymous struct, along with a Runtime field
		// that has the type string and the necessary struct tags. It's important that we
		// embed the MovieAlias type here, rather than the Movie type directly, to avoid
		// inheriting the MarshalJSON() method of the Movie type (which would result in an
		// infinite loop during encoding).
		aux := struct {
			MovieAlias
			Runtime string `json:"runtime,omitempty"`
		}{
			MovieAlias: MovieAlias(m),
			Runtime:    runtime,
		}
		return json.Marshal(aux)
	}
*/

// Use the Check() method to execute our validation checks. This will add the
// provided key and error message to the errors map if the check does not evaluate
// to true. For example, in the first line here we "check that the title is not
// equal to the empty string". In the second, we "check that the length of the title
// is less than or equal to 500 bytes" and so on.
// BODY='{"title":"","year":1000,"runtime":"-123 mins","genres":["sci-fi","sci-fi"]}'
// curl -i -d "$BODY" localhost:4000/v1/movies
// Positive: BODY='{"title":"Moana","year":2016,"runtime":"107 mins","genres":["animation","adventure"]}'
func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
