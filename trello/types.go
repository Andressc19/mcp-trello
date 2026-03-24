package trello

import "encoding/json"

// Board represents a Trello board.
type Board struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Desc     string  `json:"desc,omitempty"`
	URL      string  `json:"url,omitempty"`
	ShortURL string  `json:"shortUrl,omitempty"`
	Closed   bool    `json:"closed,omitempty"`
	Lists    []List  `json:"lists,omitempty"`
	Labels   []Label `json:"labels,omitempty"`
}

// List represents a Trello list.
type List struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Closed  bool   `json:"closed,omitempty"`
	Pos     string `json:"pos,omitempty"`
	IDBoard string `json:"idBoard,omitempty"`
}

// Card represents a Trello card.
type Card struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Desc        string            `json:"desc,omitempty"`
	Closed      bool              `json:"closed,omitempty"`
	Due         string            `json:"due,omitempty"`
	DueComplete bool              `json:"dueComplete,omitempty"`
	IDList      string            `json:"idList,omitempty"`
	IDBoard     string            `json:"idBoard,omitempty"`
	URL         string            `json:"url,omitempty"`
	ShortURL    string            `json:"shortUrl,omitempty"`
	IDLabels    []string          `json:"idLabels,omitempty"`
	Labels      []Label           `json:"labels,omitempty"`
	Members     []Member          `json:"members,omitempty"`
	Checklists  []Checklist       `json:"checklists,omitempty"`
	Actions     []json.RawMessage `json:"actions,omitempty"`
}

// Label represents a Trello label.
type Label struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Uses    int    `json:"uses,omitempty"`
	IDBoard string `json:"idBoard,omitempty"`
}

// Member represents a Trello board member.
type Member struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Username string `json:"username,omitempty"`
}

// Checklist represents a Trello checklist.
type Checklist struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IDBoard    string      `json:"idBoard,omitempty"`
	IDCard     string      `json:"idCard,omitempty"`
	CheckItems []CheckItem `json:"checkItems,omitempty"`
}

// CheckItem represents an item in a checklist.
type CheckItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state,omitempty"`
	IDChecklist string `json:"idChecklist,omitempty"`
}

// CardInput holds params for creating a card.
type CardInput struct {
	ListID      string   `json:"listId"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	DueDate     string   `json:"dueDate,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

// CardUpdate holds params for updating a card.
type CardUpdate struct {
	CardID      string `json:"cardId"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"dueDate,omitempty"`
	DueComplete bool   `json:"dueComplete,omitempty"`
	ListID      string `json:"listId,omitempty"`
}
