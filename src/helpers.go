package main

import (
	"strconv"

	"github.com/Mt-Lampert/htmx_on_fiber/src/internal/db"
)

type ProperContact struct {
	ID    string
	First string
	Last  string
	Phone string
	Email string
}

// converts raw contacts from SQLc into workable contacts
func getProperContacts(rContacts []db.Contact) []ProperContact {
	var pContacts []ProperContact

	for _, rc := range rContacts {
		var pc ProperContact
		if rc.ID.Valid {
			pc.ID = strconv.FormatInt(rc.ID.Int64, 10)
		}
		if rc.LastName.Valid {
			pc.Last = rc.LastName.String
		}
		if rc.FirstName.Valid {
			pc.First = rc.FirstName.String
		}
		if rc.Phone.Valid {
			pc.Phone = rc.Phone.String
		}
		if rc.Email.Valid {
			pc.Email = rc.Email.String
		}

		pContacts = append(pContacts, pc)
	}

	return pContacts
}

// vim: foldmethod=indent
