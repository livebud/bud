package session

import "time"

type Session struct {
	UserID   *int
	ID       string    `session:"id"`
	IP       string    `session:"ip_address"`
	UA       string    `session:"user_agent"`
	LastSeen time.Time `session:"last_seen"`
	Expires  time.Time `session:"expires"`
}

func (s *Session) Visitor() bool {
	return s.UserID == nil
}

// func LoadUser(db *sql.Database, session *Session) (*User, error) {
//   if session.ID == nil {
//     return nil, nil
//   }
//   // hypothetical database call
//   return db.FindUserByID(*session.ID)
// }

// type User struct {
//   *table.User
// }
