package users

type UserPaged struct {
	Users    []User `json:"users"`
	ScrollId string `json:"scrollId"`
}
