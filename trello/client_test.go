package trello

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func setupClientTest(t *testing.T) (*TrelloClient, *[]capturedRequest) {
	t.Helper()
	var reqs []capturedRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := capturedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.Query(),
		}
		reqs = append(reqs, req)

		handler := getRouteHandler(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(handler.status)
		w.Write(handler.body)
	}))
	t.Cleanup(server.Close)

	client := &TrelloClient{
		APIKey:  "test-key",
		Token:   "test-token",
		BaseURL: server.URL,
	}
	return client, &reqs
}

type capturedRequest struct {
	Method string
	Path   string
	Query  url.Values
}

type routeHandler struct {
	status int
	body   []byte
}

func getRouteHandler(method, path string) routeHandler {
	boards := []Board{
		{ID: "board-1", Name: "Test Board", URL: "https://trello.com/b/abc"},
	}
	board := Board{ID: "board-1", Name: "Test Board", Desc: "A board", URL: "https://trello.com/b/abc"}
	labels := []Label{
		{ID: "label-1", Name: "Bug", Color: "red", Uses: 5},
	}
	lists := []List{
		{ID: "list-1", Name: "To Do", IDBoard: "board-1"},
	}
	list := List{ID: "list-1", Name: "To Do", IDBoard: "board-1"}
	cards := []Card{
		{ID: "card-1", Name: "Task", Desc: "Do something", IDList: "list-1"},
	}
	card := Card{ID: "card-1", Name: "Task", Desc: "Do something", IDList: "list-1", IDBoard: "board-1", URL: "https://trello.com/c/xyz"}
	label := Label{ID: "label-1", Name: "Bug", Color: "red"}
	checklist := Checklist{ID: "cl-1", Name: "Checklist", IDCard: "card-1"}
	checkItem := CheckItem{ID: "ci-1", Name: "Item 1", State: "incomplete"}

	switch {
	// Boards
	case method == "GET" && path == "/members/me/boards":
		return routeHandler{200, mustJSON(boards)}
	case method == "GET" && path == "/boards/board-1":
		return routeHandler{200, mustJSON(board)}
	case method == "GET" && path == "/boards/board-1/labels":
		return routeHandler{200, mustJSON(labels)}
	// Lists
	case method == "GET" && path == "/boards/board-1/lists":
		return routeHandler{200, mustJSON(lists)}
	case method == "POST" && path == "/lists":
		return routeHandler{200, mustJSON(list)}
	case method == "PUT" && path == "/lists/list-1":
		return routeHandler{200, mustJSON(list)}
	// Cards
	case method == "GET" && path == "/lists/list-1/cards":
		return routeHandler{200, mustJSON(cards)}
	case method == "GET" && path == "/cards/card-1":
		return routeHandler{200, mustJSON(card)}
	case method == "POST" && path == "/cards":
		return routeHandler{200, mustJSON(card)}
	case method == "PUT" && path == "/cards/card-1":
		return routeHandler{200, mustJSON(card)}
	case method == "DELETE" && path == "/cards/card-1":
		return routeHandler{200, []byte("{}")}
	// Labels
	case method == "POST" && path == "/boards/board-1/labels":
		return routeHandler{200, mustJSON(label)}
	case method == "POST" && path == "/cards/card-1/idLabels":
		return routeHandler{200, []byte("{}")}
	case method == "DELETE" && path == "/cards/card-1/idLabels/label-1":
		return routeHandler{200, []byte("{}")}
	// Checklists
	case method == "POST" && path == "/checklists":
		return routeHandler{200, mustJSON(checklist)}
	case method == "POST" && path == "/checklists/cl-1/checkItems":
		return routeHandler{200, mustJSON(checkItem)}
	case method == "PUT" && path == "/cards/card-1/checklist/cl-1/checkItem/ci-1":
		return routeHandler{200, []byte("{}")}
	// Error routes
	case path == "/error/401":
		return routeHandler{401, []byte("Unauthorized")}
	case path == "/error/500":
		return routeHandler{500, []byte("Internal Server Error")}
	default:
		return routeHandler{404, []byte("not found")}
	}
}

func mustJSON(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}

// --- Auth tests ---

func TestRequest_AuthParams(t *testing.T) {
	client, reqs := setupClientTest(t)

	_, err := client.ListBoards()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Query.Get("key") != "test-key" {
		t.Errorf("expected key=test-key, got %s", r.Query.Get("key"))
	}
	if r.Query.Get("token") != "test-token" {
		t.Errorf("expected token=test-token, got %s", r.Query.Get("token"))
	}
}

// --- Boards ---

func TestListBoards(t *testing.T) {
	client, reqs := setupClientTest(t)

	boards, err := client.ListBoards()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(boards) != 1 {
		t.Fatalf("expected 1 board, got %d", len(boards))
	}
	if boards[0].Name != "Test Board" {
		t.Errorf("expected name=Test Board, got %s", boards[0].Name)
	}

	r := (*reqs)[0]
	if r.Method != "GET" {
		t.Errorf("expected GET, got %s", r.Method)
	}
	if r.Path != "/members/me/boards" {
		t.Errorf("expected path=/members/me/boards, got %s", r.Path)
	}
}

func TestGetBoard(t *testing.T) {
	client, reqs := setupClientTest(t)

	board, err := client.GetBoard("board-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if board.ID != "board-1" {
		t.Errorf("expected id=board-1, got %s", board.ID)
	}

	r := (*reqs)[0]
	if r.Path != "/boards/board-1" {
		t.Errorf("expected path=/boards/board-1, got %s", r.Path)
	}
	if r.Query.Get("fields") != "name,id,url,desc,closed" {
		t.Errorf("unexpected fields param: %s", r.Query.Get("fields"))
	}
}

func TestGetBoardLabels(t *testing.T) {
	client, reqs := setupClientTest(t)

	labels, err := client.GetBoardLabels("board-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(labels))
	}
	if labels[0].Color != "red" {
		t.Errorf("expected color=red, got %s", labels[0].Color)
	}

	r := (*reqs)[0]
	if r.Path != "/boards/board-1/labels" {
		t.Errorf("expected path=/boards/board-1/labels, got %s", r.Path)
	}
}

// --- Lists ---

func TestGetLists(t *testing.T) {
	client, reqs := setupClientTest(t)

	lists, err := client.GetLists("board-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lists) != 1 {
		t.Fatalf("expected 1 list, got %d", len(lists))
	}

	r := (*reqs)[0]
	if r.Method != "GET" {
		t.Errorf("expected GET, got %s", r.Method)
	}
	if r.Path != "/boards/board-1/lists" {
		t.Errorf("expected path=/boards/board-1/lists, got %s", r.Path)
	}
}

func TestCreateList(t *testing.T) {
	client, reqs := setupClientTest(t)

	list, err := client.CreateList("board-1", "New List")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if list.Name != "To Do" {
		t.Errorf("expected name=To Do from mock, got %s", list.Name)
	}

	r := (*reqs)[0]
	if r.Method != "POST" {
		t.Errorf("expected POST, got %s", r.Method)
	}
	if r.Path != "/lists" {
		t.Errorf("expected path=/lists, got %s", r.Path)
	}
	if r.Query.Get("idBoard") != "board-1" {
		t.Errorf("expected idBoard=board-1, got %s", r.Query.Get("idBoard"))
	}
	if r.Query.Get("name") != "New List" {
		t.Errorf("expected name=New List, got %s", r.Query.Get("name"))
	}
}

func TestUpdateList(t *testing.T) {
	client, reqs := setupClientTest(t)

	_, err := client.UpdateList("list-1", map[string]string{"name": "Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "PUT" {
		t.Errorf("expected PUT, got %s", r.Method)
	}
	if r.Path != "/lists/list-1" {
		t.Errorf("expected path=/lists/list-1, got %s", r.Path)
	}
}

func TestArchiveList(t *testing.T) {
	client, reqs := setupClientTest(t)

	_, err := client.ArchiveList("list-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "PUT" {
		t.Errorf("expected PUT, got %s", r.Method)
	}
	if r.Query.Get("closed") != "true" {
		t.Errorf("expected closed=true, got %s", r.Query.Get("closed"))
	}
}

// --- Cards ---

func TestGetCardsByList(t *testing.T) {
	client, reqs := setupClientTest(t)

	cards, err := client.GetCardsByList("list-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(cards))
	}

	r := (*reqs)[0]
	if r.Path != "/lists/list-1/cards" {
		t.Errorf("expected path=/lists/list-1/cards, got %s", r.Path)
	}
}

func TestGetCard(t *testing.T) {
	client, reqs := setupClientTest(t)

	card, err := client.GetCard("card-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.Desc != "Do something" {
		t.Errorf("expected desc=Do something, got %s", card.Desc)
	}

	r := (*reqs)[0]
	if r.Path != "/cards/card-1" {
		t.Errorf("expected path=/cards/card-1, got %s", r.Path)
	}
}

func TestAddCard(t *testing.T) {
	client, reqs := setupClientTest(t)

	input := CardInput{
		ListID:      "list-1",
		Name:        "New Card",
		Description: "A card",
		Labels:      []string{"label-1", "label-2"},
	}
	_, err := client.AddCard(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "POST" {
		t.Errorf("expected POST, got %s", r.Method)
	}
	if r.Path != "/cards" {
		t.Errorf("expected path=/cards, got %s", r.Path)
	}
	if r.Query.Get("idList") != "list-1" {
		t.Errorf("expected idList=list-1, got %s", r.Query.Get("idList"))
	}
	if r.Query.Get("name") != "New Card" {
		t.Errorf("expected name=New Card, got %s", r.Query.Get("name"))
	}
	if r.Query.Get("desc") != "A card" {
		t.Errorf("expected desc=A card, got %s", r.Query.Get("desc"))
	}
	if r.Query.Get("idLabels") != "label-1,label-2" {
		t.Errorf("expected idLabels=label-1,label-2, got %s", r.Query.Get("idLabels"))
	}
}

func TestUpdateCard(t *testing.T) {
	client, reqs := setupClientTest(t)

	input := CardUpdate{
		CardID: "card-1",
		Name:   "Updated",
	}
	_, err := client.UpdateCard(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "PUT" {
		t.Errorf("expected PUT, got %s", r.Method)
	}
	if r.Path != "/cards/card-1" {
		t.Errorf("expected path=/cards/card-1, got %s", r.Path)
	}
	if r.Query.Get("name") != "Updated" {
		t.Errorf("expected name=Updated, got %s", r.Query.Get("name"))
	}
}

func TestMoveCard(t *testing.T) {
	client, reqs := setupClientTest(t)

	_, err := client.MoveCard("card-1", "list-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "PUT" {
		t.Errorf("expected PUT, got %s", r.Method)
	}
	if r.Query.Get("idList") != "list-2" {
		t.Errorf("expected idList=list-2, got %s", r.Query.Get("idList"))
	}
}

func TestArchiveCard(t *testing.T) {
	client, reqs := setupClientTest(t)

	_, err := client.ArchiveCard("card-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Query.Get("closed") != "true" {
		t.Errorf("expected closed=true, got %s", r.Query.Get("closed"))
	}
}

func TestDeleteCard(t *testing.T) {
	client, reqs := setupClientTest(t)

	err := client.DeleteCard("card-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", r.Method)
	}
	if r.Path != "/cards/card-1" {
		t.Errorf("expected path=/cards/card-1, got %s", r.Path)
	}
}

// --- Labels ---

func TestCreateLabel(t *testing.T) {
	client, reqs := setupClientTest(t)

	label, err := client.CreateLabel("board-1", "Bug", "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.Name != "Bug" {
		t.Errorf("expected name=Bug, got %s", label.Name)
	}

	r := (*reqs)[0]
	if r.Method != "POST" {
		t.Errorf("expected POST, got %s", r.Method)
	}
	if r.Path != "/boards/board-1/labels" {
		t.Errorf("expected path=/boards/board-1/labels, got %s", r.Path)
	}
	if r.Query.Get("name") != "Bug" {
		t.Errorf("expected name=Bug, got %s", r.Query.Get("name"))
	}
	if r.Query.Get("color") != "red" {
		t.Errorf("expected color=red, got %s", r.Query.Get("color"))
	}
}

func TestAddLabelToCard(t *testing.T) {
	client, reqs := setupClientTest(t)

	err := client.AddLabelToCard("card-1", "label-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "POST" {
		t.Errorf("expected POST, got %s", r.Method)
	}
	if r.Path != "/cards/card-1/idLabels" {
		t.Errorf("expected path=/cards/card-1/idLabels, got %s", r.Path)
	}
	if r.Query.Get("value") != "label-1" {
		t.Errorf("expected value=label-1, got %s", r.Query.Get("value"))
	}
}

func TestRemoveLabelFromCard(t *testing.T) {
	client, reqs := setupClientTest(t)

	err := client.RemoveLabelFromCard("card-1", "label-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", r.Method)
	}
	if r.Path != "/cards/card-1/idLabels/label-1" {
		t.Errorf("expected path=/cards/card-1/idLabels/label-1, got %s", r.Path)
	}
}

// --- Checklists ---

func TestCreateChecklist(t *testing.T) {
	client, reqs := setupClientTest(t)

	cl, err := client.CreateChecklist("card-1", "My Checklist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cl.Name != "Checklist" {
		t.Errorf("expected name=Checklist from mock, got %s", cl.Name)
	}

	r := (*reqs)[0]
	if r.Method != "POST" {
		t.Errorf("expected POST, got %s", r.Method)
	}
	if r.Path != "/checklists" {
		t.Errorf("expected path=/checklists, got %s", r.Path)
	}
	if r.Query.Get("idCard") != "card-1" {
		t.Errorf("expected idCard=card-1, got %s", r.Query.Get("idCard"))
	}
}

func TestAddChecklistItem(t *testing.T) {
	client, reqs := setupClientTest(t)

	item, err := client.AddChecklistItem("cl-1", "Do the thing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Name != "Item 1" {
		t.Errorf("expected name=Item 1 from mock, got %s", item.Name)
	}

	r := (*reqs)[0]
	if r.Method != "POST" {
		t.Errorf("expected POST, got %s", r.Method)
	}
	if r.Path != "/checklists/cl-1/checkItems" {
		t.Errorf("expected path=/checklists/cl-1/checkItems, got %s", r.Path)
	}
	if r.Query.Get("name") != "Do the thing" {
		t.Errorf("expected name=Do the thing, got %s", r.Query.Get("name"))
	}
}

func TestCompleteCheckItem(t *testing.T) {
	client, reqs := setupClientTest(t)

	err := client.CompleteCheckItem("card-1", "cl-1", "ci-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Method != "PUT" {
		t.Errorf("expected PUT, got %s", r.Method)
	}
	if r.Path != "/cards/card-1/checklist/cl-1/checkItem/ci-1" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.Query.Get("state") != "complete" {
		t.Errorf("expected state=complete, got %s", r.Query.Get("state"))
	}
}

func TestUncompleteCheckItem(t *testing.T) {
	client, reqs := setupClientTest(t)

	err := client.UncompleteCheckItem("card-1", "cl-1", "ci-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := (*reqs)[0]
	if r.Query.Get("state") != "incomplete" {
		t.Errorf("expected state=incomplete, got %s", r.Query.Get("state"))
	}
}

// --- Error handling ---

func TestRequest_401Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
	}))
	t.Cleanup(server.Close)

	client := &TrelloClient{
		APIKey:  "bad-key",
		Token:   "bad-token",
		BaseURL: server.URL,
	}

	_, err := client.ListBoards()
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestRequest_500Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
	}))
	t.Cleanup(server.Close)

	client := &TrelloClient{
		APIKey:  "test-key",
		Token:   "test-token",
		BaseURL: server.URL,
	}

	_, err := client.ListBoards()
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
