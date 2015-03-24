package codebase

import (
    "fmt"
    "log"
    "strconv"
    "time"
)

type ticketArray struct {
    Tickets []Ticket `xml:"ticket"`
}

type Ticket struct {
    Id            int    `xml:"ticket-id"`
    Summary       string `xml:"summary"`
    TicketType    string `xml:"ticket-type"`
    Assignee      string `xml:"assignee"`
    Status        string `xml:"status>name"`
    Priority      string `xml:"priority>name"`
    EstimationRaw string `xml:"estimated-time"`
    Reporter      string `xml:"reporter"`
    Category      string `xml:"category>name"`

    Milestone      Milestone `xml:"milestone"`
    StartOn        string    `xml:"start-on"`
    Deadline       string    `xml:"deadline"`
    UpdatedAt      time.Time `xml:"updated-at"`
    CreatedAt      time.Time `xml:"created-at"`
    TotalTimeSpent int       `xml:"total-time-spent"`
    BlockedBy      []int     `xml:"blocked-by>blocked-by"`
    Blocking       []int     `xml:"blocking>blocking"`
}

type TicketByAssignee []Ticket

func (t TicketByAssignee) Len() int           { return len(t) }
func (t TicketByAssignee) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TicketByAssignee) Less(i, j int) bool { return t[i].Assignee < t[j].Assignee }

type ticketQueryOptions struct {
    baseQueryOptions
    Query string `url:"query,omitempty"`
}

func (t *Ticket) IsAssigned() bool {
    return t.Assignee != ""
}

func (t *Ticket) Estimation() int64 {
    if t.EstimationRaw == "" {
        return 0
    }

    val, err := strconv.ParseInt(t.EstimationRaw, 10, 0)
    if err != nil {
        log.Fatalln(err)
    }
    return val
}

func (c *CodeBaseAPI) TicketsForMilestone(milestone Milestone) (tickets []Ticket) {
    query := fmt.Sprintf("resolution:open milestone:%s", strconv.Quote(milestone.Name))

    queryOpts := ticketQueryOptions{
        Query: query,
    }
    queryOpts.Page = 1

    for {
        proxy := ticketArray{}
        if err := c.fetchFromCodebase("tickets", &proxy, queryOpts); err != nil {
            log.Fatalln("Could not fetch tickets:", err)
        }

        if len(proxy.Tickets) == 0 {
            break
        }

        tickets = append(tickets, proxy.Tickets...)
        queryOpts.Page++
    }

    return
}

func (c *CodeBaseAPI) TicketById(id int) (ticket Ticket, err error) {
    endpoint := fmt.Sprintf("tickets/%v", id)
    err = c.fetchFromCodebase(endpoint, &ticket, nil)
    return
}
