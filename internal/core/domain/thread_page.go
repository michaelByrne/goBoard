package domain

type ThreadPage struct {
	PageNum            int
	TotalPages         int
	Threads            []Thread
	DefaultThreadLimit int
	HasNextPage        bool
	HasPrevPage        bool
}
