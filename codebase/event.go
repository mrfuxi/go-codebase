package codebase

import (
    "fmt"
    "log"
    "time"
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

type EventChanges struct {
    Status   []string `xml:"status-id>status-id"`
    Assignee []string `xml:"assignee-id>assignee-id"`
    Subject  []string `xml:"subject>subject"`
    Priority []string `xml:"priority-id>priority-id"`
}

type eventQueryOptions struct {
    baseQueryOptions
    Raw   bool      `url:"raw,omitempty"`
    Since time.Time `url:"since,omitempty"`
}

type ChangeMapping struct {
    Status map[string]string
}

func (c *CodeBaseAPI) Activities(since time.Time, user User) (events []Event) {
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

            if event.HasChanges(ChangeMapping{}) == false {
                continue
            }

            events = append(events, event)
        }
        queryOpts.Page++
    }

    return
}

func (e *Event) HasChanges(mapping ChangeMapping) bool {
    return e.Raw.Changes.Changes(mapping) != ""
}

func (e *Event) TicketUrl(company string) string {
    return fmt.Sprintf("https://%s.codebasehq.com/projects/%s/tickets/%v", company, e.Raw.ProjectPermalink, e.Raw.TicketID)
}

func (e *Event) Day() time.Weekday {
    date, _ := time.Parse("2006-01-02 15:04:05 MST", e.Date)
    return date.Weekday()
}

func (e *EventChanges) Changes(mapping ChangeMapping) string {
    changes := ""

    workType := mapping.Status

    if len(e.Status) == 1 {
        changes += e.Status[0]
    } else if len(e.Status) == 2 {
        change := fmt.Sprintf("%s -> %s", e.Status[0], e.Status[1])
        if description, ok := workType[change]; ok {
            change = description
        }

        if change != "" {
            changes += change
        }
    }

    return changes
}