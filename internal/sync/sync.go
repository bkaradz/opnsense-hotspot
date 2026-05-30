package sync

import (
	"database/sql"
	"fmt"

	"github.com/templui/templui-quickstart/internal/api"
)

type Voucher struct {
	Username   string  `json:"username"`
	Validity   int     `json:"validity"`
	Expirytime int64   `json:"expirytime"`
	Starttime  float64 `json:"starttime"`
	Endtime    int64   `json:"endtime"`
	State      string  `json:"state"`
}

// SyncAll syncs providers, voucher groups, and vouchers into SQLite
// without ever increasing autoincrement sequences unnecessarily.
func SyncAll(db *sql.DB, baseURL string) error {
	var providers []string
	if err := api.FetchVouchersJSON(baseURL+"/list_providers", &providers); err != nil {
		return fmt.Errorf("fetch providers: %w", err)
	}

	for _, provider := range providers {
		var providerID int
		err := db.QueryRow(`SELECT id FROM providers WHERE name=?`, provider).Scan(&providerID)
		if err == sql.ErrNoRows {
			res, err := db.Exec(`INSERT INTO providers (name) VALUES (?)`, provider)
			if err != nil {
				return fmt.Errorf("insert provider: %w", err)
			}
			lastID, _ := res.LastInsertId()
			providerID = int(lastID)
		} else if err != nil {
			return fmt.Errorf("select provider: %w", err)
		}

		// --- Fetch voucher groups for this provider ---
		var groups []string
		groupURL := fmt.Sprintf("%s/list_voucher_groups/%s", baseURL, provider)
		if err := api.FetchVouchersJSON(groupURL, &groups); err != nil {
			return fmt.Errorf("fetch groups for %s: %w", provider, err)
		}

		for _, group := range groups {
			var groupID int
			err = db.QueryRow(`SELECT id FROM voucher_groups WHERE provider_id=? AND name=?`, providerID, group).Scan(&groupID)
			if err == sql.ErrNoRows {
				res, err := db.Exec(`INSERT INTO voucher_groups (provider_id, name) VALUES (?, ?)`, providerID, group)
				if err != nil {
					return fmt.Errorf("insert group: %w", err)
				}
				lastID, _ := res.LastInsertId()
				groupID = int(lastID)
			} else if err != nil {
				return fmt.Errorf("select group: %w", err)
			}

			// --- Fetch vouchers for this group ---
			var vouchers []Voucher
			voucherURL := fmt.Sprintf("%s/list_vouchers/%s/%s", baseURL, provider, group)
			if err := api.FetchVouchersJSON(voucherURL, &vouchers); err != nil {
				return fmt.Errorf("fetch vouchers for %s/%s: %w", provider, group, err)
			}

			// --- Sequence-safe batch insert/update ---
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("begin tx: %w", err)
			}

			selectStmt, err := tx.Prepare(`SELECT id FROM vouchers WHERE username=? AND group_id=?`)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("prepare select: %w", err)
			}
			insertStmt, err := tx.Prepare(`
				INSERT INTO vouchers (group_id, username, validity, expirytime, starttime, endtime, state)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("prepare insert: %w", err)
			}
			updateStmt, err := tx.Prepare(`
				UPDATE vouchers
				SET validity=?, expirytime=?, starttime=?, endtime=?, state=?
				WHERE username=? AND group_id=?
			`)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("prepare update: %w", err)
			}

			defer selectStmt.Close()
			defer insertStmt.Close()
			defer updateStmt.Close()

			for _, v := range vouchers {
				var existingID int
				err := selectStmt.QueryRow(v.Username, groupID).Scan(&existingID)
				if err == sql.ErrNoRows {
					// Only insert if new (increments sequence once)
					_, err = insertStmt.Exec(groupID, v.Username, v.Validity, v.Expirytime, v.Starttime, v.Endtime, v.State)
				} else if err == nil {
					// Update existing (no sequence bump)
					_, err = updateStmt.Exec(v.Validity, v.Expirytime, v.Starttime, v.Endtime, v.State, v.Username, groupID)
				}
				if err != nil {
					tx.Rollback()
					return fmt.Errorf("sync voucher (%s): %w", v.Username, err)
				}
			}

			if err := tx.Commit(); err != nil {
				return fmt.Errorf("commit tx: %w", err)
			}
		}
	}

	return nil
}
