package graphql

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/noonyuu/comparison/backend/internal/db"
)

var Schema = `
    type Movie {
        id: ID!
        email: String
        name: String
        lastName: String
        firstName: String
        avatarUrl: String
        nickName: String
        providerId: String
    }

    type Query {
        movies: [Movie!]!
        movie(id: ID!): Movie
    }

    type Mutation {
        CreateMovie(email: String, name: String, last_name: String, first_name: String, avatar_url: String, nick_name: String, provider_id: String): Movie
    }
`

type Resolver struct {
	DB *db.Database
}

func (r *Resolver) Movies() ([]*MovieResolver, error) {
	movies, err := r.DB.GetMovies()
	if err != nil {
		return nil, err
	}

	var resolvers []*MovieResolver
	for _, m := range movies {
		resolvers = append(resolvers, &MovieResolver{m})
	}
	return resolvers, nil
}

func (r *Resolver) Movie(args struct{ ID graphql.ID }) (*MovieResolver, error) {
	movie, err := r.DB.GetMovie(string(args.ID))
	if err != nil {
		return nil, err
	}
	return &MovieResolver{movie}, nil
}

func (r *Resolver) CreateMovie(args struct {
	Email      *string
	Name       *string
	LastName   *string
	FirstName  *string
	AvatarURL  *string
	NickName   *string
	ProviderID *string
}) (*MovieResolver, error) {
	movie := &db.Movie{
		Email:      *args.Email,
		Name:       *args.Name,
		LastName:   *args.LastName,
		FirstName:  *args.FirstName,
		AvatarURL:  *args.AvatarURL,
		NickName:   *args.NickName,
		ProviderID: *args.ProviderID,
	}
	err := r.DB.CreateMovie(movie)
	if err != nil {
		return nil, err
	}
	return &MovieResolver{movie}, nil
}

type MovieResolver struct {
	m *db.Movie
}

func (r *MovieResolver) ID() graphql.ID      { return graphql.ID(r.m.ID) }
func (r *MovieResolver) Email() *string      { return &r.m.Email }
func (r *MovieResolver) Name() *string       { return &r.m.Name }
func (r *MovieResolver) LastName() *string   { return &r.m.LastName }
func (r *MovieResolver) FirstName() *string  { return &r.m.FirstName }
func (r *MovieResolver) AvatarURL() *string  { return &r.m.AvatarURL }
func (r *MovieResolver) NickName() *string   { return &r.m.NickName }
func (r *MovieResolver) ProviderID() *string { return &r.m.ProviderID }

func NewHandler(database *db.Database) *relay.Handler {
	schema := graphql.MustParseSchema(Schema, &Resolver{DB: database})
	return &relay.Handler{Schema: schema}
}
