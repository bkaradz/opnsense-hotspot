package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/templui/templui-quickstart/internal/sync"
)

func SyncHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := os.Getenv("BASE_URL")
		if err := sync.SyncAll(db, baseURL); err != nil {
			http.Error(w, fmt.Sprintf("Sync failed: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", `{"show-toast": {"title": "Sync Completed", "description": "Database synchronization finished successfully.", "variant": "success"}}`)
		w.WriteHeader(http.StatusOK)
	}
}
