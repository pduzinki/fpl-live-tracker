package domain

// Club represents a football club that plays in the Premier League (e.g. Liverpool, Manchester City, Chelsea)
type Club struct {
	ID        int    `bson:"ID"`
	Name      string `bson:"Name"`
	Shortname string `bson:"Shortname"`
}

type ClubRepository interface {
	Add(club Club) error
	AddMany(clubs []Club) error
	GetByID(id int) (Club, error)
}
