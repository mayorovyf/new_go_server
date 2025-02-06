package database

// User представляет пользователя с расширенной информацией
type User struct {
	FirstName  string `bson:"first_name" json:"first_name"`
	LastName   string `bson:"last_name" json:"last_name"`
	MiddleName string `bson:"middle_name" json:"middle_name"`
	UniqueID   string `bson:"unique_id" json:"unique_id"`
	Class      string `bson:"class" json:"class"`
	Username   string `bson:"username" json:"username"`
	Password   string `bson:"password" json:"password"`
}
