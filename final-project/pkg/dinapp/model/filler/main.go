package filler

import (
	"final-project/pkg/dinapp/model"
)

func PopulateDatabase(models model.Models) error {
	for _, genre := range genres {
		models.Genres.InsertG(&genre)
	}

	for _, movie := range movies {
		models.Movies.Insert(&movie)
	}

	// for _, book := range books {
	// 	models.Books.Insert(&book)
	// }

	return nil
}

var genres = []model.Genres{
	{Title: "Drama"},
	{Title: "Horror"},
	{Title: "Comedy"},
}

var movies = []model.Movies{
	{Title: "The Shawshank Redemption", Description: "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.", YearOfProduction: 1994, GenreId: "1"},                                                                                                        // Drama
	{Title: "The Godfather", Description: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.", YearOfProduction: 1972, GenreId: "1"},                                                                                                                      // Drama
	{Title: "The Dark Knight", Description: "When the menace known as The Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.", YearOfProduction: 2008, GenreId: "3"},                                          // Action
	{Title: "Pulp Fiction", Description: "The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.", YearOfProduction: 1994, GenreId: "2"},                                                                                           // Crime
	{Title: "Forrest Gump", Description: "The presidencies of Kennedy and Johnson, the events of Vietnam, Watergate, and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.", YearOfProduction: 1994, GenreId: "1"}, // Drama
}

// var books = []model.Book{
//     {Title: "The Shawshank Redemption", Description: "The novella written by Stephen King, upon which the movie is based.", Genre: 1, MovieRec: 1}, // Drama
//     {Title: "The Godfather", Description: "Mario Puzo's classic novel about the Corleone crime family.", Genre: 1, MovieRec: 2}, // Drama
//     {Title: "The Dark Knight Returns", Description: "A graphic novel by Frank Miller, one of the inspirations for the movie 'The Dark Knight'.", Genre: 3, MovieRec: 3}, // Action
//     {Title: "Pulp Fiction: The Complete Screenplay", Description: "The complete screenplay of Quentin Tarantino's masterpiece.", Genre: 2, MovieRec: 4}, // Crime
//     {Title: "Forrest Gump", Description: "The novel by Winston Groom, on which the movie is based.", Genre: 1, MovieRec: 5}, // Drama
// }
