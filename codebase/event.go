package codebase

import (
    "fmt"
    "log"
    "time"
)

const (
    TYPE_NEW_TICKET = "ticketing_ticket"
)

type Event struct {
    Id        int    `xml:"id"`
    Title     string `xml:"title"`
    Date      string `xml:"timestamp"`
    Type      string `xml:"type"`
    HtmlTitle string `xml:"html-title"`
    HtmlText  string `xml:"html-text"`
    UserID    int    `xml:"user-id"`
    UserEmail string `xml:"actor-email"`
    ProjectID int    `xml:"project-id"`
    Deleted   bool   `xml:"deleted"`
    Avatar    string `xml:"avatar-url"`

    Raw struct {
        TicketID         int          `xml:"number"`
        Subject          string       `xml:"subject"`
        Content          string       `xml:"content"`
        Changes          EventChanges `xml:"changes"`
        ProjectPermalink string       `xml:"project-permalink"`
    }   `xml:"raw-properties"`
}

type eventQueryOptions struct {
    baseQueryOptions
    Raw   bool      `url:"raw,omitempty"`
    Since time.Time `url:"since,omitempty"`
}

type Descriptor interface {
    MapChange(field, before, after string) (description string)
}

func (c *CodeBaseAPI) Activities(since time.Time, user User, descriptor Descriptor) (events []Event) {
    type eventArray struct {
        Events []Event `xml:"event"`
    }

    queryOpts := eventQueryOptions{
        Raw:   true,
        Since: since,
    }
    queryOpts.Page = 1

    for {
        proxy := eventArray{}
        if err := c.fetchFromCodebase("activity", &proxy, queryOpts); err != nil {
            log.Fatalln("Could not fetch activities:", err)
        }

        if len(proxy.Events) == 0 {
            break
        }

        for _, event := range proxy.Events {
            if user.UserID != 0 && event.UserID != user.UserID {
                continue
            }

            if event.HasChanges(descriptor) == false {
                continue
            }

            events = append(events, event)
        }
        queryOpts.Page++
    }

    return
}

func (e *Event) HasChanges(descriptor Descriptor) bool {
    return e.Changes(descriptor) != ""
}

func (e *Event) TicketUrl(company string) string {
    return fmt.Sprintf("https://%s.codebasehq.com/projects/%s/tickets/%v", company, e.Raw.ProjectPermalink, e.Raw.TicketID)
}

func (e *Event) Day() time.Weekday {
    date, _ := time.Parse("2006-01-02 15:04:05 MST", e.Date)
    return date.Weekday()
}

func (e *Event) Changes(descriptor Descriptor) string {
    chagnesToMap := e.Raw.Changes.mappedChanges()

    if e.Type == TYPE_NEW_TICKET {
        change := changeToMap{changeType: CHANGE_NEW_TICKET}
        chagnesToMap = append(chagnesToMap, change)
    }

    return describeChanges(chagnesToMap, descriptor)
}
