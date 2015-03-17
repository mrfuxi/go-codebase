package codebase

import (
    "errors"
    "fmt"
    "log"
    "sort"
    "strings"
    "time"
)

type milestoneArray struct {
    Milestones []Milestone `xml:"ticketing-milestone"`
}

type Milestone struct {
    Id          int    `xml:"id"`
    Name        string `xml:"name"`
    StartAt     string `xml:"start-at"`
    Deadline    string `xml:"deadline"`
    Status      string `xml:"status"`
    Description string `xml:"description"`
}

type byMostRecent struct {
    ms []Milestone
}

func (m byMostRecent) Len() int           { return len(m.ms) }
func (m byMostRecent) Swap(i, j int)      { m.ms[i], m.ms[j] = m.ms[j], m.ms[i] }
func (m byMostRecent) Less(i, j int) bool { return m.ms[i].toEadges() < m.ms[j].toEadges() }

func (m *Milestone) IsActive() bool {
    return m.Status == "active"
}

func (m *Milestone) toEadges() float64 {
    layout := "2006-01-02"

    start_date := time.Unix(0, 0)
    end_date := time.Unix(0, 0)
    if m.StartAt != "" {
        start_date, _ = time.Parse(layout, m.StartAt)
    }
    if m.Deadline != "" {
        end_date, _ = time.Parse(layout, m.Deadline)
    }

    to_start := time.Since(start_date).Hours()
    to_end := time.Since(end_date).Hours()

    return to_start + to_end
}

func (m *Milestone) DaysToEnd() (int, error) {
    layout := "2006-01-02"

    if m.Deadline == "" {
        return 0, errors.New("Deadline not available")
    }

    end_date, err := time.Parse(layout, m.Deadline)
    if err != nil {
        return 0, err
    }

    beginningOfToday := time.Now().Truncate(24 * time.Hour)
    hoursToEnd := end_date.Sub(beginningOfToday).Hours()
    return int(hoursToEnd / 24), nil
}

func (m *Milestone) DaysToEndStr() string {
    daysToEnd, err := m.DaysToEnd()
    if err != nil {
        return ""
    }

    if daysToEnd == 0 {
        return "Ends today"
    }

    ending := ""
    if daysToEnd < -1 || 1 < daysToEnd {
        ending = "s"
    }

    if daysToEnd > 0 {
        return fmt.Sprintf("Ends in %v day%v", daysToEnd, ending)
    } else {
        return fmt.Sprintf("Ended %v day%v ago", -daysToEnd, ending)
    }
}

func (c *CodeBaseAPI) Milesones() (milestones []Milestone) {
    queryOpts := baseQueryOptions{}

    proxy := milestoneArray{}
    if err := c.fetchFromCodebase("milestones", &proxy, queryOpts); err != nil {
        log.Fatalln("Could not fetch milestones:", err)
    }

    milestones = proxy.Milestones
    return
}

func (c *CodeBaseAPI) CurrentMilestone() (current Milestone) {
    allMilestones := c.Milesones()

    var activeMilestones []Milestone
    for _, milestone := range allMilestones {
        if !milestone.IsActive() {
            continue
        }

        activeMilestones = append(activeMilestones, milestone)
    }

    sort.Reverse(byMostRecent{activeMilestones})
    current = activeMilestones[0]

    return
}

func (m Milestone) String() string {
    daysToEnd := m.DaysToEndStr()
    if daysToEnd != "" {
        return fmt.Sprintf("%v (%v)", m.Name, strings.ToLower(daysToEnd))
    }

    return m.Name
}
