package model

type Rsvp struct{
	Rsvp_id int `json:"r_id"`
	Event_id int `json:"e_id"`
	Attendee_id int `json:"a_id"`
}