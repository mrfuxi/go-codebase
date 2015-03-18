package codebase

import (
    "fmt"
    "log"
    "time"
)

type ticketNotesArray struct {
    Notes []TicketNote `xml:"ticket-note"`
}

type TicketNote struct {
    // Attachements []Attachement
    Content    string    `xml:"content"`
    Id         int       `xml:"id"`
    UpdatedAt  time.Time `xml:"updated-at"`
    CreatedAt  time.Time `xml:"created-at"`
    UserID     int       `xml:"user-id"`
    UpdatesRaw string    `xml:"updates"`
}

func (t *TicketNote) ChangesState() bool {
    return t.UpdatesRaw != "{}"
}

func (c *CodeBaseAPI) NotesForTicket(id int) []TicketNote {
    endpoint := fmt.Sprintf("tickets/%v/notes", id)

    proxy := ticketNotesArray{}
    if err := c.fetchFromCodebase(endpoint, &proxy, nil); err != nil {
        log.Fatalln("Could not fetch ticket notes:", err)
    }

    return proxy.Notes
}
