package codebase

import (
    "errors"
    "fmt"
    "log"
    "sort"
    "strings"
    "time"
)

const (
    timeLayout = "2006-01-02"
)

type milestoneArray struct {
    Milestones []Milestone `xml:"ticketing-milestone"`
}

type Milestone struct {
    Id            int    `xml:"id"`
    Name          string `xml:"name"`
    StartAt       string `xml:"start-at"`
    Deadline      string `xml:"deadline"`
    Status        string `xml:"status"`
    Description   string `xml:"description"`
    EstimatedTime int64  `xml:"estimated-time"`
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

func (m *Milestone) StartAtTime() (t time.Time, err error) {
    if m.StartAt != "" {
        t, err = time.Parse(timeLayout, m.StartAt)
    } else {
        err = errors.New("Start time not available")
    }
    return
}

func (m *Milestone) DeadlineTime() (t time.Time, err error) {
    if m.Deadline != "" {
        t, err = time.Parse(timeLayout, m.Deadline)
    } else {
        err = errors.New("Deadline not available")
    }
    return
}

func (m *Milestone) toEadges() float64 {
    start_date := time.Unix(0, 0)
    end_date := time.Unix(0, 0)
    if m.StartAt != "" {
        start_date, _ = time.Parse(timeLayout, m.StartAt)
    }
    if m.Deadline != "" {
        end_date, _ = time.Parse(timeLayout, m.Deadline)
    }

    to_start := time.Since(start_date).Hours()
    to_end := time.Since(end_date).Hours()

    return to_start + to_end
}

func (m *Milestone) DaysToEnd() (int64, error) {
    end_date, err := m.DeadlineTime()
    if err != nil {
        return 0, err
    }

    days := WorkingDays(end_date).SinceUnix() - WorkingDays(time.Now()).SinceUnix()
    return days, nil
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

func (m *Milestone) EstimatedTimeDuration() time.Duration {
    if m.EstimatedTime <= 0 {
        return 0
    }

    return time.Duration(m.EstimatedTime) * time.Minute
}
