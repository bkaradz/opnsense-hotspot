package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/templui/templui-quickstart/internal/database"
	"github.com/templui/templui-quickstart/internal/printing"
)

func PrintingHandler(ctx context.Context, queries *database.Queries, w http.ResponseWriter, r *http.Request) error {

	var payload struct {
		Ids     []string `json:"ids"`
		Printer string   `json:"printer"`
	}

	// --- Read Body ---
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return fmt.Errorf("read body: %w", err)
	}
	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return fmt.Errorf("empty body")
	}

	// --- Decode JSON ---
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return fmt.Errorf("json unmarshal: %w", err)
	}

	// --- Process Each Voucher ---
	for _, vID := range payload.Ids {

		voucherInt, err := strconv.ParseInt(vID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid voucher ID: "+vID, http.StatusBadRequest)
			return fmt.Errorf("invalid voucher ID %s: %w", vID, err)
		}

		selectedVoucher, err := queries.GetVoucherByID(ctx, voucherInt)
		if err != nil {
			http.Error(w, "Voucher not found: "+vID, http.StatusNotFound)
			return fmt.Errorf("get voucher %d: %w", voucherInt, err)
		}

		// Print
		if err := printing.PrintVoucher(selectedVoucher, payload.Printer); err != nil {
			http.Error(w, "Failed to print voucher "+vID, http.StatusInternalServerError)
			return fmt.Errorf("print voucher %d: %w", voucherInt, err)
		}

		// Mark printed
		if err := queries.UpdateVouchersPrintedField(ctx, voucherInt); err != nil {
			http.Error(w, "Failed to mark voucher printed: "+vID, http.StatusInternalServerError)
			return fmt.Errorf("update printed field for %d: %w", voucherInt, err)
		}
	}

	// --- Write simple OK response (HTMX / fetch expects this) ---
	w.Write([]byte("OK"))

	return nil
}
