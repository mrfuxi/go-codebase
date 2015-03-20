package codebase

import (
    "fmt"
    "strings"
)

const (
    CHANGE_STATUS     = "status"
    CHANGE_MILESTONE  = "milestone"
    CHANGE_CATEGORY   = "category"
    CHANGE_PRIORITY   = "priority"
    CHANGE_NEW_TICKET = "new ticket"
)

type EventChanges struct {
    Status    []string `json:"status_id" xml:"status-id>status-id"`
    Assignee  []string `json:"assignee_id" xml:"assignee-id>assignee-id"`
    Subject   []string `json:"subject" xml:"subject>subject"`
    Priority  []string `json:"priority_id" xml:"priority-id>priority-id"`
    Milestone []string `json:"milestone_id" xml:"milestone-id>milestone-id"`
    Category  []string `json:"category_id" xml:"category-id>category-id"`
}

type changeToMap struct {
    changeType string
    from       string
    to         string
}

func (e *EventChanges) mappedChanges() []changeToMap {
    chagnesToMap := make([]changeToMap, 0)

    if len(e.Status) == 2 {
        change := changeToMap{CHANGE_STATUS, e.Status[0], e.Status[1]}
        chagnesToMap = append(chagnesToMap, change)
    }

    if len(e.Milestone) == 2 {
        change := changeToMap{CHANGE_MILESTONE, e.Milestone[0], e.Milestone[1]}
        chagnesToMap = append(chagnesToMap, change)
    }

    if len(e.Category) == 2 {
        change := changeToMap{CHANGE_CATEGORY, e.Category[0], e.Category[1]}
        chagnesToMap = append(chagnesToMap, change)
    }

    if len(e.Priority) == 2 {
        change := changeToMap{CHANGE_PRIORITY, e.Priority[0], e.Priority[1]}
        chagnesToMap = append(chagnesToMap, change)
    }

    return chagnesToMap
}

func describeChanges(changesToMap []changeToMap, descriptor Descriptor) string {
    changes := make([]string, 0)

    for _, change := range changesToMap {
        changeDescription := fmt.Sprintf("%s -> %s", change.from, change.to)

        if descriptor != nil {
            changeDescription = descriptor.MapChange(change.changeType, change.from, change.to)
        }

        if changeDescription != "" {
            changes = append(changes, changeDescription)
        }
    }

    return strings.Join(changes, ", ")
}
