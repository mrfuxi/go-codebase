package codebase

import (
    "fmt"
    "log"
    "strings"
)

type User struct {
    FirstName string `xml:"first-name"`
    LastName  string `xml:"last-name"`
    UserName  string `xml:"username"`
    UserID    int    `xml:"id"`
    Email     string `xml:"email-address"`
    Company   string `xml:"company"`
}

type userQueryOptions struct {
    *baseQueryOptions
}

// Fetch User object for currently authenticated user
func (c *CodeBaseAPI) AuthUser() (user User) {
    for _, u := range c.UsersInProject() {
        if u.userNameAuth() == c.userNameAuth {
            return u
        }
    }

    log.Fatalln("Cound not find your user:", c.userNameAuth)
    return
}

// Fetch User based on first name
func (c *CodeBaseAPI) User(firstName, lastname, username string) (user User) {
    for _, u := range c.UsersInProject() {
        if username != "" && strings.ToLower(u.UserName) == strings.ToLower(username) {
            return u
        }

        if lastname != "" && strings.ToLower(u.LastName) == strings.ToLower(lastname) {
            return u
        }

        if firstName != "" && strings.ToLower(u.FirstName) == strings.ToLower(firstName) {
            return u
        }
    }

    names := []string{firstName, lastname, username}

    log.Fatalln("Cound not find user:", strings.Join(names, " "))
    return
}

func (c *CodeBaseAPI) UsersInProject() []User {
    if cachedUsers, ok := c.users[c.Project]; ok {
        return cachedUsers
    }

    type userArray struct {
        Users []User `xml:"user"`
    }
    proxy := userArray{}

    if err := c.fetchFromCodebase("assignments", &proxy, userQueryOptions{}); err != nil {
        log.Fatalln("Could not fetch assigned users:", err)
    }

    c.users[c.Project] = proxy.Users

    return proxy.Users
}

func (u *User) userNameAuth() string {
    return fmt.Sprintf("%s/%s", strings.ToLower(u.Company), u.UserName)
}

// Strinify user to "FistName LastName"
func (u User) String() string {
    return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
