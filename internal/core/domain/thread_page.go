package domain

type ThreadPage struct {
	PageNum    int
	TotalPages int
	Threads    []Thread
}
