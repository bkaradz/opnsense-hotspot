package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/a-h/templ"
	"github.com/templui/templui-quickstart/internal/database"
	"github.com/templui/templui-quickstart/internal/pagination"
	"github.com/templui/templui-quickstart/internal/printing"
	"github.com/templui/templui-quickstart/ui/layouts"
	"github.com/templui/templui-quickstart/ui/pages"
)

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func int32Ptr(i int32, isSet bool) *int32 {
	if !isSet {
		return nil
	}
	return &i
}

func GetVouchersData(ctx context.Context, queries *database.Queries, q url.Values) layouts.VouchersResponse {

	validityStr := q.Get("validity")
	limitStr := q.Get("limit")
	pageStr := q.Get("page")
	state := q.Get("state")
	search := q.Get("search")
	group := q.Get("group_name")
	printed := q.Get("printed")
	selectedPrinter := q.Get("printer")

	page := int64(1)
	limit := int64(10)
	var validity int32
	hasValidity := false

	if validityStr != "" {
		if v, err := strconv.Atoi(validityStr); err == nil {
			validity = int32(v)
			hasValidity = true
		}
	}

	if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
		limit = int64(v)
	}
	if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
		page = int64(v)
	}

	offset := (page - 1) * limit

	groupNameUnq, err := queries.GroupNameListUniqueVouchers(ctx, int32Ptr(validity, hasValidity))
	if err != nil {
		log.Fatal("Group Name unique Error: ", err)
	}

	if len(groupNameUnq) > 0 && !slices.Contains(groupNameUnq, group) {
		group = groupNameUnq[0]
	}

	vouchers, err := queries.ListVouchers(ctx, database.ListVouchersParams{
		State:     strPtr(state),
		Search:    strPtr(search),
		GroupName: strPtr(group),
		Printed:   strPtr(printed),
		Validity:  int32Ptr(validity, hasValidity),
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return layouts.VouchersResponse{}
	}

	total, err := queries.CountVouchers(ctx, database.CountVouchersParams{
		State:     strPtr(state),
		Search:    strPtr(search),
		GroupName: strPtr(group),
		Printed:   strPtr(printed),
		Validity:  int32Ptr(validity, hasValidity),
	})
	if err != nil {
		return layouts.VouchersResponse{}
	}

	return apiCommon(ctx, queries, vouchers, int(limit), int(total), page, validity, hasValidity, state, group, printed, int(limit), groupNameUnq, selectedPrinter)

}

func apiCommon(ctx context.Context, queries *database.Queries, vouchers []database.ListVouchersRow, limit int, total int, page int64, validity int32, hasValidity bool, currentState string, currentGroupName string, currentPrinted string, currentLimit int, groupNameUnq []string, selectedPrinter string) layouts.VouchersResponse {
	statesUnq, err := queries.StateListUniqueVouchers(ctx)

	if err != nil {
		log.Fatal("State unique Error: ", err)
	}

	validityUnq, err := queries.ValidityListUniqueVouchers(ctx)

	if err != nil {
		log.Fatal("Validity unique Error: ", err)
	}

	printedUnq, err := queries.PrintedListUniqueVouchers(ctx)

	if err != nil {
		log.Fatal("Printerd unique Error: ", err)
	}

	printers := printing.GetPrinterName()

	// Default to the first available printer if none is selected
	if selectedPrinter == "" && len(printers) > 0 {
		selectedPrinter = printers[0]
	}

	filterValues := layouts.FilterValues{
		StatesUnique:    statesUnq,
		ValidityUnique:  validityUnq,
		GroupNameUnique: groupNameUnq,
		PrintedUnique:   printedUnq,
		Printers:        printers,
	}

	pagination := pagination.CalculateState(int(page), int(limit), int(total))

	apiVouchers := layouts.VouchersResponse{
		Vouchers:         make([]layouts.Voucher, 0, len(vouchers)),
		FilterValues:     filterValues,
		CurrentState:     currentState,
		CurrentValidity:  validity,
		CurrentGroupName: currentGroupName,
		CurrentPrinted:   currentPrinted,
		CurrentLimit:     currentLimit,
		Pagination:       pagination,
		SelectedPrinter:  selectedPrinter,
	}

	for _, r := range vouchers {
		apiVouchers.Vouchers = append(apiVouchers.Vouchers, layouts.Voucher{
			ID:        int(r.ID),
			Username:  r.Username,
			Validity:  int(r.Validity),
			State:     r.State,
			GroupName: r.GroupName,
			Printed:   r.Printed,
			ProviderName: layouts.ProviderName{
				String: r.ProviderName.String,
				Valid:  r.ProviderName.Valid,
			},
		})
	}

	return apiVouchers

}

func GetVouchersHandler(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := GetVouchersData(r.Context(), queries, r.URL.Query())

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func UpdateVouchersHandler(ctx context.Context, queries *database.Queries, q url.Values) (templ.Component, error) {

	resp := GetVouchersData(ctx, queries, q)

	// fmt.Println("Vouchers: ", resp.Vouchers)
	// fmt.Println("FilterValues: ", resp.FilterValues)
	// fmt.Println("CurrentState: ", resp.CurrentState)
	// fmt.Println("CurrentValidity: ", resp.CurrentValidity)
	// fmt.Println("CurrentGroupName: ", resp.CurrentGroupName)
	// fmt.Println("CurrentPrinted: ", resp.CurrentPrinted)
	// fmt.Println("CurrentLimit: ", resp.CurrentLimit)
	// fmt.Println("Pagination: ", resp.Pagination)

	return templ.Join(
		pages.VoucherTable(resp),
		pages.VoucherGroupName(resp),
	), nil
}
