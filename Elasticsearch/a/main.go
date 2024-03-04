package main

import (
	"context"
	"ngoc/logic"
)

func main() {

	ctx := context.Background()

	ctx = logic.LoadMoviesFromFile(ctx)
	ctx = logic.ConnectWithElasticsearch(ctx)
	logic.IndexMoviesAsDocuments(ctx)
	logic.QueryMovieByDocumentID(ctx)
	logic.BestKeanuActionMovies(ctx)
	logic.MovieCountPerGenreAgg(ctx)

}
