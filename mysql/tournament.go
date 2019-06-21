package mysql

// Tournament represents a tournament in a social tournaments service.
type Tournament struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Deposit int64   `json:"deposit"`
	Prize   int64   `json:"prize"`
	Winner  int64   `json:"winner"`
	Users   []int64 `json:"users"`
}
