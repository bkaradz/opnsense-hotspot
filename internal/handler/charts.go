package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/templui/templui-quickstart/internal/database"
	"github.com/templui/templui-quickstart/ui/components/chart"
)

func GetDailyChartData(ctx context.Context, queries *database.Queries) chart.Data {
	now := time.Now()

	// Calculate end time (today at 23:59:59)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// Calculate start time (6 days ago at 00:00:00)
	startTime := endTime.AddDate(0, 0, -6)
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())

	expiredVouchers, err := queries.GetExpiredVouchersPerDay(ctx, database.GetExpiredVouchersPerDayParams{
		Endtime:   sql.NullInt64{Int64: startTime.Unix(), Valid: true},
		Endtime_2: sql.NullInt64{Int64: endTime.Unix(), Valid: true},
	})
	if err != nil {
		log.Println("Error fetching daily chart data:", err)
		return chart.Data{}
	}

	// Generate labels for the last 7 days
	labels := make([]string, 7)
	for i := 0; i < 7; i++ {
		day := startTime.AddDate(0, 0, i)
		labels[i] = day.Format("Mon")
	}

	// Process data for chart
	// Map: Validity -> DayIndex (0-6) -> Count
	dataMap := make(map[int64]map[int]float64)
	validities := []int64{}

	for _, v := range expiredVouchers {
		// v.DayOfWeek is string "0" (Sun) to "6" (Sat) from sqlite strftime '%w'
		dayOfWeek, _ := strconv.Atoi(v.DayOfWeek.(string))

		// Find which index (0-6) this day corresponds to in our labels
		// We iterate through our 7 days window to match the day of week
		dayIndex := -1
		for i := 0; i < 7; i++ {
			// Check if the day of week matches
			checkDay := startTime.AddDate(0, 0, i)
			if int(checkDay.Weekday()) == dayOfWeek {
				dayIndex = i
				break
			}
		}

		if dayIndex != -1 {
			if _, ok := dataMap[v.Validity]; !ok {
				dataMap[v.Validity] = make(map[int]float64)
				validities = append(validities, v.Validity)
			}
			dataMap[v.Validity][dayIndex] += float64(v.Count)
		}
	}

	datasets := []chart.Dataset{}
	for _, validity := range validities {
		data := make([]float64, 7)
		for i := 0; i < 7; i++ {
			data[i] = dataMap[validity][i]
		}
		datasets = append(datasets, chart.Dataset{
			Label: fmt.Sprintf("%d", validity),
			Data:  data,
		})
	}

	return chart.Data{
		Labels:   labels,
		Datasets: datasets,
	}
}

func GetWeeklyChartData(ctx context.Context, queries *database.Queries) chart.Data {
	now := time.Now()

	// Calculate end time (today at 23:59:59)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// Calculate start time (7 weeks ago at 00:00:00, total 8 weeks)
	startTime := endTime.AddDate(0, 0, -7*7)
	// Adjust to start of the week (Monday)
	for startTime.Weekday() != time.Monday {
		startTime = startTime.AddDate(0, 0, -1)
	}
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())

	expiredVouchers, err := queries.GetExpiredVouchersPerWeek(ctx, database.GetExpiredVouchersPerWeekParams{
		Endtime:   sql.NullInt64{Int64: startTime.Unix(), Valid: true},
		Endtime_2: sql.NullInt64{Int64: endTime.Unix(), Valid: true},
	})
	if err != nil {
		log.Println("Error fetching weekly chart data:", err)
		return chart.Data{}
	}

	// Generate labels for the last 8 weeks
	labels := make([]string, 8)
	weekMap := make(map[string]int) // Map "YYYY-WW" to index 0-7
	for i := 0; i < 8; i++ {
		weekStart := startTime.AddDate(0, 0, i*7)
		_, week := weekStart.ISOWeek()
		year := weekStart.Year()
		weekStr := fmt.Sprintf("%d-%02d", year, week)
		labels[i] = fmt.Sprintf("W%d", week)
		weekMap[weekStr] = i
	}

	// Process data for chart
	dataMap := make(map[int64]map[int]float64)
	validities := []int64{}

	for _, v := range expiredVouchers {
		weekStr := v.Week.(string)
		if idx, ok := weekMap[weekStr]; ok {
			if _, exists := dataMap[v.Validity]; !exists {
				dataMap[v.Validity] = make(map[int]float64)
				validities = append(validities, v.Validity)
			}
			dataMap[v.Validity][idx] += float64(v.Count)
		}
	}

	datasets := []chart.Dataset{}
	for _, validity := range validities {
		data := make([]float64, 8)
		for i := 0; i < 8; i++ {
			data[i] = dataMap[validity][i]
		}
		datasets = append(datasets, chart.Dataset{
			Label: fmt.Sprintf("%d", validity),
			Data:  data,
		})
	}

	return chart.Data{
		Labels:   labels,
		Datasets: datasets,
	}
}

func GetMonthlyChartData(ctx context.Context, queries *database.Queries) chart.Data {
	now := time.Now()

	// Calculate end time (today at 23:59:59)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// Calculate start time (11 months ago to cover a year, or just current year)
	// Let's show last 12 months
	startTime := endTime.AddDate(0, -11, 0)
	startTime = time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, startTime.Location())

	expiredVouchers, err := queries.GetExpiredVouchersPerMonth(ctx, database.GetExpiredVouchersPerMonthParams{
		Endtime:   sql.NullInt64{Int64: startTime.Unix(), Valid: true},
		Endtime_2: sql.NullInt64{Int64: endTime.Unix(), Valid: true},
	})
	if err != nil {
		log.Println("Error fetching monthly chart data:", err)
		return chart.Data{}
	}

	// Generate labels for the last 12 months
	labels := make([]string, 12)
	monthMap := make(map[string]int) // Map "YYYY-MM" to index 0-11
	for i := 0; i < 12; i++ {
		month := startTime.AddDate(0, i, 0)
		monthStr := month.Format("2006-01")
		labels[i] = month.Format("Jan")
		monthMap[monthStr] = i
	}

	// Process data for chart
	dataMap := make(map[int64]map[int]float64)
	validities := []int64{}

	for _, v := range expiredVouchers {
		monthStr := v.Month.(string)
		if idx, ok := monthMap[monthStr]; ok {
			if _, exists := dataMap[v.Validity]; !exists {
				dataMap[v.Validity] = make(map[int]float64)
				validities = append(validities, v.Validity)
			}
			dataMap[v.Validity][idx] += float64(v.Count)
		}
	}

	datasets := []chart.Dataset{}
	for _, validity := range validities {
		data := make([]float64, 12)
		for i := 0; i < 12; i++ {
			data[i] = dataMap[validity][i]
		}
		datasets = append(datasets, chart.Dataset{
			Label: fmt.Sprintf("%d", validity),
			Data:  data,
		})
	}

	return chart.Data{
		Labels:   labels,
		Datasets: datasets,
	}
}
