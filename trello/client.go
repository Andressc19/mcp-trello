package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://api.trello.com/1"

// TrelloClient makes authenticated requests to the Trello REST API.
type TrelloClient struct {
	APIKey  string
	Token   string
	BaseURL string
}

// NewTrelloClient creates a TrelloClient with the given credentials.
func NewTrelloClient(apiKey, token string) *TrelloClient {
	return &TrelloClient{
		APIKey:  apiKey,
		Token:   token,
		BaseURL: baseURL,
	}
}

// request performs an HTTP request with key/token query param auth.
func (c *TrelloClient) request(method, endpoint string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	q.Set("key", c.APIKey)
	q.Set("token", c.Token)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	log.Printf("[trello] %s %s", method, u.Path)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("trello API error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// ============ BOARDS ============

// ListBoards returns all boards for the authenticated user.
func (c *TrelloClient) ListBoards() ([]Board, error) {
	body, err := c.request("GET", "/members/me/boards", map[string]string{
		"fields": "name,id,url,shortUrl",
	})
	if err != nil {
		return nil, err
	}
	var boards []Board
	if err := json.Unmarshal(body, &boards); err != nil {
		return nil, fmt.Errorf("parsing boards: %w", err)
	}
	return boards, nil
}

// GetBoard returns details for a specific board.
func (c *TrelloClient) GetBoard(boardID string) (*Board, error) {
	body, err := c.request("GET", "/boards/"+boardID, map[string]string{
		"fields": "name,id,url,desc,closed",
		"lists":  "open",
		"labels": "all",
	})
	if err != nil {
		return nil, err
	}
	var board Board
	if err := json.Unmarshal(body, &board); err != nil {
		return nil, fmt.Errorf("parsing board: %w", err)
	}
	return &board, nil
}

// GetBoardLabels returns all labels on a board.
func (c *TrelloClient) GetBoardLabels(boardID string) ([]Label, error) {
	body, err := c.request("GET", "/boards/"+boardID+"/labels", map[string]string{
		"fields": "name,color,uses",
	})
	if err != nil {
		return nil, err
	}
	var labels []Label
	if err := json.Unmarshal(body, &labels); err != nil {
		return nil, fmt.Errorf("parsing labels: %w", err)
	}
	return labels, nil
}

// ============ LISTS ============

// GetLists returns all open lists on a board.
func (c *TrelloClient) GetLists(boardID string) ([]List, error) {
	body, err := c.request("GET", "/boards/"+boardID+"/lists", map[string]string{
		"fields": "name,id,closed,pos",
	})
	if err != nil {
		return nil, err
	}
	var lists []List
	if err := json.Unmarshal(body, &lists); err != nil {
		return nil, fmt.Errorf("parsing lists: %w", err)
	}
	return lists, nil
}

// CreateList creates a new list on a board.
func (c *TrelloClient) CreateList(boardID, name string) (*List, error) {
	body, err := c.request("POST", "/lists", map[string]string{
		"idBoard": boardID,
		"name":    name,
	})
	if err != nil {
		return nil, err
	}
	var list List
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parsing list: %w", err)
	}
	return &list, nil
}

// UpdateList updates a list's name and/or closed state.
func (c *TrelloClient) UpdateList(listID string, params map[string]string) (*List, error) {
	body, err := c.request("PUT", "/lists/"+listID, params)
	if err != nil {
		return nil, err
	}
	var list List
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parsing list: %w", err)
	}
	return &list, nil
}

// ArchiveList archives (closes) a list.
func (c *TrelloClient) ArchiveList(listID string) (*List, error) {
	body, err := c.request("PUT", "/lists/"+listID, map[string]string{
		"closed": "true",
	})
	if err != nil {
		return nil, err
	}
	var list List
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("parsing list: %w", err)
	}
	return &list, nil
}

// ============ CARDS ============

// GetCardsByList returns all cards in a list.
func (c *TrelloClient) GetCardsByList(listID string) ([]Card, error) {
	body, err := c.request("GET", "/lists/"+listID+"/cards", map[string]string{
		"fields":          "name,id,desc,closed,due,dueComplete,idList,idBoard",
		"checkItemStates": "true",
		"labels":          "true",
		"members":         "true",
	})
	if err != nil {
		return nil, err
	}
	var cards []Card
	if err := json.Unmarshal(body, &cards); err != nil {
		return nil, fmt.Errorf("parsing cards: %w", err)
	}
	return cards, nil
}

// GetCard returns details for a specific card.
func (c *TrelloClient) GetCard(cardID string) (*Card, error) {
	body, err := c.request("GET", "/cards/"+cardID, map[string]string{
		"fields":     "name,id,desc,closed,due,dueComplete,idList,idBoard,url,shortUrl",
		"checklists": "all",
		"labels":     "all",
		"members":    "all",
		"actions":    "all",
	})
	if err != nil {
		return nil, err
	}
	var card Card
	if err := json.Unmarshal(body, &card); err != nil {
		return nil, fmt.Errorf("parsing card: %w", err)
	}
	return &card, nil
}

// AddCard creates a new card on a list.
func (c *TrelloClient) AddCard(input CardInput) (*Card, error) {
	params := map[string]string{
		"idList": input.ListID,
		"name":   input.Name,
	}
	if input.Description != "" {
		params["desc"] = input.Description
	}
	if input.DueDate != "" {
		params["due"] = input.DueDate
	}
	if len(input.Labels) > 0 {
		params["idLabels"] = strings.Join(input.Labels, ",")
	}

	body, err := c.request("POST", "/cards", params)
	if err != nil {
		return nil, err
	}
	var card Card
	if err := json.Unmarshal(body, &card); err != nil {
		return nil, fmt.Errorf("parsing card: %w", err)
	}
	return &card, nil
}

// UpdateCard updates a card's fields.
func (c *TrelloClient) UpdateCard(input CardUpdate) (*Card, error) {
	params := map[string]string{}
	if input.Name != "" {
		params["name"] = input.Name
	}
	if input.Description != "" {
		params["desc"] = input.Description
	}
	if input.DueDate != "" {
		params["due"] = input.DueDate
	}
	if input.DueComplete {
		params["dueComplete"] = "true"
	}
	if input.ListID != "" {
		params["idList"] = input.ListID
	}

	body, err := c.request("PUT", "/cards/"+input.CardID, params)
	if err != nil {
		return nil, err
	}
	var card Card
	if err := json.Unmarshal(body, &card); err != nil {
		return nil, fmt.Errorf("parsing card: %w", err)
	}
	return &card, nil
}

// MoveCard moves a card to a different list.
func (c *TrelloClient) MoveCard(cardID, listID string) (*Card, error) {
	body, err := c.request("PUT", "/cards/"+cardID, map[string]string{
		"idList": listID,
	})
	if err != nil {
		return nil, err
	}
	var card Card
	if err := json.Unmarshal(body, &card); err != nil {
		return nil, fmt.Errorf("parsing card: %w", err)
	}
	return &card, nil
}

// ArchiveCard archives (closes) a card.
func (c *TrelloClient) ArchiveCard(cardID string) (*Card, error) {
	body, err := c.request("PUT", "/cards/"+cardID, map[string]string{
		"closed": "true",
	})
	if err != nil {
		return nil, err
	}
	var card Card
	if err := json.Unmarshal(body, &card); err != nil {
		return nil, fmt.Errorf("parsing card: %w", err)
	}
	return &card, nil
}

// DeleteCard permanently deletes a card.
func (c *TrelloClient) DeleteCard(cardID string) error {
	_, err := c.request("DELETE", "/cards/"+cardID, nil)
	return err
}

// ============ LABELS ============

// CreateLabel creates a new label on a board.
func (c *TrelloClient) CreateLabel(boardID, name, color string) (*Label, error) {
	body, err := c.request("POST", "/boards/"+boardID+"/labels", map[string]string{
		"name":  name,
		"color": color,
	})
	if err != nil {
		return nil, err
	}
	var label Label
	if err := json.Unmarshal(body, &label); err != nil {
		return nil, fmt.Errorf("parsing label: %w", err)
	}
	return &label, nil
}

// AddLabelToCard adds a label to a card.
func (c *TrelloClient) AddLabelToCard(cardID, labelID string) error {
	_, err := c.request("POST", "/cards/"+cardID+"/idLabels", map[string]string{
		"value": labelID,
	})
	return err
}

// RemoveLabelFromCard removes a label from a card.
func (c *TrelloClient) RemoveLabelFromCard(cardID, labelID string) error {
	_, err := c.request("DELETE", "/cards/"+cardID+"/idLabels/"+labelID, nil)
	return err
}

// ============ CHECKLISTS ============

// CreateChecklist creates a new checklist on a card.
func (c *TrelloClient) CreateChecklist(cardID, name string) (*Checklist, error) {
	if name == "" {
		name = "Checklist"
	}
	body, err := c.request("POST", "/checklists", map[string]string{
		"idCard": cardID,
		"name":   name,
	})
	if err != nil {
		return nil, err
	}
	var cl Checklist
	if err := json.Unmarshal(body, &cl); err != nil {
		return nil, fmt.Errorf("parsing checklist: %w", err)
	}
	return &cl, nil
}

// AddChecklistItem adds an item to a checklist.
func (c *TrelloClient) AddChecklistItem(checklistID, text string) (*CheckItem, error) {
	body, err := c.request("POST", "/checklists/"+checklistID+"/checkItems", map[string]string{
		"name": text,
	})
	if err != nil {
		return nil, err
	}
	var item CheckItem
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, fmt.Errorf("parsing checkitem: %w", err)
	}
	return &item, nil
}

// CompleteCheckItem marks a checklist item as complete.
func (c *TrelloClient) CompleteCheckItem(cardID, checklistID, itemID string) error {
	_, err := c.request("PUT",
		"/cards/"+cardID+"/checklist/"+checklistID+"/checkItem/"+itemID,
		map[string]string{"state": "complete"},
	)
	return err
}

// UncompleteCheckItem marks a checklist item as incomplete.
func (c *TrelloClient) UncompleteCheckItem(cardID, checklistID, itemID string) error {
	_, err := c.request("PUT",
		"/cards/"+cardID+"/checklist/"+checklistID+"/checkItem/"+itemID,
		map[string]string{"state": "incomplete"},
	)
	return err
}
