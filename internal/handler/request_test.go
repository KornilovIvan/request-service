package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"request-service/internal/model"
	"request-service/internal/testdb"
)

func TestListActiveExcludesPending(t *testing.T) {
	srv := newTestServer(t)
	customerID := createCustomer(t, srv, "Альфа")
	activeID := createRequest(t, srv, customerID, "Активная", "active")
	createRequest(t, srv, customerID, "Ожидает", "pending")

	status, body := call(t, srv, http.MethodGet, "/api/v1/requests?status=active&created_from=2022-01-01&limit=20", nil)
	if status != http.StatusOK {
		t.Fatalf("status = %d, body = %s", status, body)
	}

	var requests []model.Request
	decode(t, body, &requests)
	found := false
	for _, request := range requests {
		if request.Status == "pending" {
			t.Fatalf("pending request in list: %+v", request)
		}
		if request.ID == activeID {
			found = true
			if request.CustomerName != "Альфа" {
				t.Fatalf("customer_name = %q", request.CustomerName)
			}
		}
	}
	if !found {
		t.Fatalf("active request %d not in %s", activeID, body)
	}
}

func TestListEmptyStatusKeepsActive(t *testing.T) {
	srv := newTestServer(t)
	customerID := createCustomer(t, srv, "Альфа")
	activeID := createRequest(t, srv, customerID, "Активная", "active")

	status, body := call(t, srv, http.MethodGet, "/api/v1/requests?status=", nil)
	if status != http.StatusOK {
		t.Fatalf("status = %d, body = %s", status, body)
	}

	var requests []model.Request
	decode(t, body, &requests)
	for _, request := range requests {
		if request.ID == activeID {
			return
		}
	}
	t.Fatalf("active request %d not in %s", activeID, body)
}

func TestListInvalidStatus(t *testing.T) {
	srv := newTestServer(t)
	status, body := call(t, srv, http.MethodGet, "/api/v1/requests?status=done", nil)
	assertError(t, status, body, http.StatusBadRequest, "invalid status")
}

func TestListInvalidLimit(t *testing.T) {
	srv := newTestServer(t)
	status, body := call(t, srv, http.MethodGet, "/api/v1/requests?limit=abc", nil)
	assertError(t, status, body, http.StatusBadRequest, "invalid limit")
}

func TestListInvalidDate(t *testing.T) {
	srv := newTestServer(t)
	status, body := call(t, srv, http.MethodGet, "/api/v1/requests?created_from=2999-01-01", nil)
	assertError(t, status, body, http.StatusBadRequest, "invalid date")
}

func TestCreateCustomerName(t *testing.T) {
	srv := newTestServer(t)

	status, body := call(t, srv, http.MethodPost, "/api/v1/customers", map[string]string{"name": "   "})
	assertError(t, status, body, http.StatusBadRequest, "name is required")

	status, body = call(t, srv, http.MethodPost, "/api/v1/customers", map[string]string{"name": "Альфа"})
	if status != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", status, body)
	}
}

func TestCreateRequestCustomer(t *testing.T) {
	srv := newTestServer(t)

	status, body := call(t, srv, http.MethodPost, "/api/v1/requests", map[string]any{
		"customer_id": 0,
		"title":       "x",
		"status":      "active",
	})
	assertError(t, status, body, http.StatusBadRequest, "customer not found")

	customerID := createCustomer(t, srv, "Альфа")
	status, body = call(t, srv, http.MethodPost, "/api/v1/requests", map[string]any{
		"customer_id": customerID,
		"title":       "Подключить отчёт",
		"status":      "active",
	})
	if status != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", status, body)
	}
}

func TestPatchClosedSetsClosedAt(t *testing.T) {
	srv := newTestServer(t)
	customerID := createCustomer(t, srv, "Альфа")
	requestID := createRequest(t, srv, customerID, "Закрыть", "active")

	status, body := call(t, srv, http.MethodPatch, "/api/v1/requests/"+strconv.Itoa(requestID)+"/status", map[string]string{
		"status": "closed",
	})
	if status != http.StatusOK {
		t.Fatalf("status = %d, body = %s", status, body)
	}

	var request model.Request
	decode(t, body, &request)
	if request.ClosedAt == nil {
		t.Fatalf("closed_at is null: %s", body)
	}
}

func TestGetMissingRequest(t *testing.T) {
	srv := newTestServer(t)
	status, body := call(t, srv, http.MethodGet, "/api/v1/requests/99999", nil)
	assertError(t, status, body, http.StatusNotFound, "request not found")
}

func TestListEmpty(t *testing.T) {
	srv := newTestServer(t)
	status, body := call(t, srv, http.MethodGet, "/api/v1/requests", nil)
	if status != http.StatusOK {
		t.Fatalf("status = %d, body = %s", status, body)
	}
	if string(body) != "[]" {
		t.Fatalf("body = %s, want []", body)
	}
}

func TestCreateDocument(t *testing.T) {
	srv := newTestServer(t)
	customerID := createCustomer(t, srv, "Альфа")
	requestID := createRequest(t, srv, customerID, "Документы", "active")

	status, body := call(t, srv, http.MethodPost, "/api/v1/requests/"+strconv.Itoa(requestID)+"/documents", map[string]string{
		"name": "Договор.pdf",
	})
	if status != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", status, body)
	}

	status, body = call(t, srv, http.MethodPost, "/api/v1/requests/"+strconv.Itoa(requestID)+"/documents", map[string]string{
		"name": "   ",
	})
	assertError(t, status, body, http.StatusBadRequest, "document name is required")
}

func TestListStatuses(t *testing.T) {
	srv := newTestServer(t)
	status, body := call(t, srv, http.MethodGet, "/api/v1/statuses", nil)
	if status != http.StatusOK {
		t.Fatalf("status = %d, body = %s", status, body)
	}

	var statuses []model.Status
	decode(t, body, &statuses)
	if len(statuses) != 3 {
		t.Fatalf("len = %d, body = %s", len(statuses), body)
	}
	got := []string{statuses[0].Code, statuses[1].Code, statuses[2].Code}
	want := []string{"active", "closed", "pending"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("codes = %v, want %v", got, want)
		}
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	db := testdb.Open(t)
	srv := httptest.NewServer(NewRouter(db))
	t.Cleanup(srv.Close)
	return srv
}

func createCustomer(t *testing.T, srv *httptest.Server, name string) int {
	t.Helper()
	status, body := call(t, srv, http.MethodPost, "/api/v1/customers", map[string]string{"name": name})
	if status != http.StatusCreated {
		t.Fatalf("create customer status = %d, body = %s", status, body)
	}
	var customer model.Customer
	decode(t, body, &customer)
	return customer.ID
}

func createRequest(t *testing.T, srv *httptest.Server, customerID int, title, statusCode string) int {
	t.Helper()
	status, body := call(t, srv, http.MethodPost, "/api/v1/requests", map[string]any{
		"customer_id": customerID,
		"title":       title,
		"status":      statusCode,
	})
	if status != http.StatusCreated {
		t.Fatalf("create request status = %d, body = %s", status, body)
	}
	var request model.Request
	decode(t, body, &request)
	return request.ID
}

func call(t *testing.T, srv *httptest.Server, method, path string, body any) (int, []byte) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, srv.URL+path, reader)
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, data
}

func assertError(t *testing.T, status int, body []byte, wantStatus int, wantMessage string) {
	t.Helper()
	if status != wantStatus {
		t.Fatalf("status = %d, body = %s", status, body)
	}
	var payload struct {
		Error string `json:"error"`
	}
	decode(t, body, &payload)
	if payload.Error != wantMessage {
		t.Fatalf("error = %q, want %q", payload.Error, wantMessage)
	}
}

func decode(t *testing.T, body []byte, dest any) {
	t.Helper()
	if err := json.Unmarshal(body, dest); err != nil {
		t.Fatalf("decode %s: %v", body, err)
	}
}
