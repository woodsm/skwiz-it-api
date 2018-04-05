package model

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    int64  `json:"id"`
}

//make sure that this a valid user
// TODO: should actually do some validation of email, id etc
func (u *User) IsValid() bool {
	if u == nil {
		return false
	}

	if len(u.Email) < 5 {
		return false
	}

	if len(u.Name) < 1 {
		return false
	}

	return true
}
