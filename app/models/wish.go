package models

type Wish struct {
	ID        int
	UID       int
	Name      string
	Content   string
	IsClaim   bool
	ClaimName string
}
