package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	_ "embed"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/templui/templui-quickstart/assets"
	"github.com/templui/templui-quickstart/internal/database"
	"github.com/templui/templui-quickstart/internal/handler"
	"github.com/templui/templui-quickstart/internal/sync"
	"github.com/templui/templui-quickstart/ui/pages"
	_ "modernc.org/sqlite"
)

//go:embed database/schema.sql
var schema string

func openDB() *sql.DB {
	conn, err := sql.Open("sqlite", "file:mydb.sqlite3?_foreign_keys=on")
	conn.SetMaxOpenConns(1)
	conn.Exec("PRAGMA journal_mode=WAL;")
	conn.Exec("PRAGMA journal_mode=MEMORY;")
	conn.Exec("PRAGMA _synchronous=OFF;")

	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func main() {

	InitDotEnv()

	dbConn := openDB()
	defer dbConn.Close()

	ctx := context.Background()

	// Run migrations (apply schema.sql)
	if _, err := dbConn.Exec(schema); err != nil {
		log.Fatal(err)
	}

	// Populate the database
	var baseURL = os.Getenv("BASE_URL")

	if err := sync.SyncAll(dbConn, baseURL); err != nil {
		log.Fatal("Sync failed:", err)
	}

	queries := database.New(dbConn)

	params := url.Values{}

	params.Set("search", os.Getenv("SEARCH"))
	params.Set("state", os.Getenv("STATE"))
	params.Set("validity", os.Getenv("VALIDITY"))
	params.Set("group_name", os.Getenv("GROUP_NAME"))
	params.Set("printed", os.Getenv("PRINTED"))
	params.Set("limit", os.Getenv("LIMIT"))

	apiVouchers := handler.GetVouchersData(ctx, queries, params)

	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	SetupAssetsRoutes(mux)

	// Chart Data Logic
	dailyChart := handler.GetDailyChartData(ctx, queries)
	weeklyChart := handler.GetWeeklyChartData(ctx, queries)
	monthlyChart := handler.GetMonthlyChartData(ctx, queries)

	mux.Handle("/", templ.Handler(pages.Landing(dailyChart, weeklyChart, monthlyChart)))
	mux.Handle("/vouchers", templ.Handler(pages.Vouchers(apiVouchers)))
	mux.HandleFunc("/api/vouchers", handler.GetVouchersHandler(queries))
	mux.HandleFunc("/api/sync", handler.SyncHandler(dbConn))
	mux.HandleFunc("/api/print", func(w http.ResponseWriter, r *http.Request) {
		handler.PrintingHandler(r.Context(), queries, w, r)
	})
	mux.HandleFunc("/api/update", func(w http.ResponseWriter, r *http.Request) {
		comp, err := handler.UpdateVouchersHandler(r.Context(), queries, r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comp.Render(r.Context(), w)
	})
	fmt.Println("Server is running on http://localhost" + port)
	log.Println(http.ListenAndServe(port, mux))
}

func InitDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	var isDevelopment = os.Getenv("GO_ENV") != "production"

	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isDevelopment {
			log.Println("Serving assets from disk (Development Mode)")
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			log.Println("Serving embedded assets (Production Mode)")
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("/assets/", http.StripPrefix("/assets/", assetHandler))
}
