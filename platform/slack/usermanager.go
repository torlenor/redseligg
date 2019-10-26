package slack

import "fmt"

type userManager struct {
	knownUsers     map[string]user   // key is UserID
	knownUserNames map[string]string // mapping of UserName to UserID
	knownUserIDs   map[string]string // mapping of UserID to UserUserNameName
}

func newUserManager() userManager {
	return userManager{
		knownUsers:     make(map[string]user),
		knownUserNames: make(map[string]string),
		knownUserIDs:   make(map[string]string),
	}
}

func (um *userManager) addKnownUser(u user) {
	um.knownUsers[u.ID] = u
	um.knownUserNames[u.Name] = u.ID
	um.knownUserIDs[u.ID] = u.Name
}

func (um userManager) getUserByID(id string) (user, error) {
	if u, ok := um.knownUsers[id]; ok {
		return u, nil
	}
	return user{}, fmt.Errorf("User with ID %s not known", id)
}

func (um userManager) getUserByName(name string) (user, error) {
	if id, ok := um.knownUserNames[name]; ok {
		return um.knownUsers[id], nil
	}
	return user{}, fmt.Errorf("User with Name %s not known", name)
}

func (um userManager) getUserNameByID(id string) (string, error) {
	if name, ok := um.knownUserIDs[id]; ok {
		return name, nil
	}
	return "", fmt.Errorf("User with ID %s not known", id)
}

func (um userManager) getUserIDByName(name string) (string, error) {
	if id, ok := um.knownUserNames[name]; ok {
		return id, nil
	}
	return "", fmt.Errorf("User with Name %s not known", name)
}

func (um userManager) isUserIDKnown(id string) bool {
	_, ok := um.knownUserIDs[id]
	return ok
}

func (um userManager) isUserNameKnown(name string) bool {
	_, ok := um.knownUserNames[name]
	return ok
}

func (um userManager) Len() int {
	return len(um.knownUsers)
}
