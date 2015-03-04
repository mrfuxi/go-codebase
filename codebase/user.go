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
func (c *CodeBaseAPI) User(firstName string) (user User) {
    for _, u := range c.UsersInProject() {
        if strings.ToLower(u.FirstName) == strings.ToLower(firstName) {
            return u
        }
    }

    log.Fatalln("Cound not find user:", firstName)
    return
}

func (c *CodeBaseAPI) UsersInProject() []User {
    type userArray struct {
        Users []User `xml:"user"`
    }

    proxy := userArray{}

    if err := c.fetchFromCodebase("assignments", &proxy, userQueryOptions{}); err != nil {
        log.Fatalln("Could not fetch assigned users:", err)
    }

    return proxy.Users
}

func (u *User) userNameAuth() string {
    return fmt.Sprintf("%s/%s", strings.ToLower(u.Company), u.UserName)
}

// Strinify user to "FistName LastName"
func (u User) String() string {
    return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
