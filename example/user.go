package main

import (
	"encoding/xml"
	"errors"
	"net/http"
	"sync"
)

type User struct {
	XMLName  xml.Name `json:"-" xml:"user"`
	Name     string   `json:"name" xml:"name"`
	Password string   `json:"password,omitempty" xml:"password,omitempty"`
}

func (u *User) SafeExport() *User {
	return &User{
		Name: u.Name,
	}
}

type userResource struct {
	sync.RWMutex
	users map[string]*User
}

func newUserResource() *userResource {
	return &userResource{users: make(map[string]*User)}
}

func (r *userResource) New() interface{} {
	return &User{}
}

func (r *userResource) Id(value interface{}) string {
	u := value.(*User)
	return u.Name
}

func (r *userResource) List(req *http.Request) (int, interface{}) {
	r.RLock()
	defer r.RUnlock()

	us := make([]*User, len(r.users), len(r.users))
	idx := 0
	for _, u := range r.users {
		us[idx] = u.SafeExport()
		idx++
	}
	return http.StatusOK, us
}

func (r *userResource) Get(req *http.Request, id string) (int, interface{}) {
	r.RLock()
	defer r.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return http.StatusNotFound, nil
	}
	return http.StatusOK, u.SafeExport()
}

func (r *userResource) Create(req *http.Request, value interface{}) (int, interface{}) {
	u := value.(*User)
	r.Lock()
	defer r.Unlock()

	_, exists := r.users[u.Name]
	if exists {
		return http.StatusBadRequest, errors.New("Already exists")
	}
	r.users[u.Name] = u
	return http.StatusCreated, u.SafeExport()
}

func (r *userResource) Update(req *http.Request, id string, value interface{}) (int, interface{}) {
	u := value.(*User)
	r.Lock()
	defer r.Unlock()

	_, exists := r.users[id]
	if !exists {
		return http.StatusNotFound, nil
	}
	r.users[u.Name] = u
	return http.StatusOK, u.SafeExport()
}

func (r *userResource) Delete(req *http.Request, id string) (int, interface{}) {
	r.Lock()
	defer r.Unlock()

	delete(r.users, id)
	return http.StatusNoContent, nil
}
