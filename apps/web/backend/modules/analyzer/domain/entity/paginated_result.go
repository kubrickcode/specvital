package entity

type PaginatedRepositoryCards struct {
	Data       []RepositoryCard
	HasNext    bool
	NextCursor *string
}
