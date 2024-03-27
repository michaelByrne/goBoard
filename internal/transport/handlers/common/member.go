package common

import (
	"errors"

	"github.com/gorilla/sessions"
)

const DefaultViewable = 3

type Member struct {
	ID       int
	Username string
	Viewable int
}

func GetMember(session sessions.Session) (Member, error) {
	id, ok := session.Values["id"]
	if !ok {
		return Member{}, errors.New("no id in session")
	}

	idInt, ok := id.(int)
	if !ok {
		return Member{}, errors.New("id is not int")
	}

	username, ok := session.Values["name"]
	if !ok {
		return Member{}, errors.New("no name in session")
	}

	usernameStr, ok := username.(string)
	if !ok {
		return Member{}, errors.New("name is not string")
	}

	viewable, ok := session.Values["collapseopen"]
	if !ok {
		viewable = DefaultViewable
	}

	viewableInt, ok := viewable.(int)
	if !ok {
		return Member{}, errors.New("collapseopen is not int")
	}

	return Member{
		ID:       idInt,
		Username: usernameStr,
		Viewable: viewableInt,
	}, nil
}
