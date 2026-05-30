package printing

import (
	"fmt"
	"log"
	"os"

	// "os"
	"strconv"
	"time"

	"github.com/DevLumuz/go-escpos"
	"github.com/skip2/go-qrcode"
	"github.com/templui/templui-quickstart/internal/database"
)

func PrintVoucher(voucher database.GetVoucherByIDRow, printerName string) error {

	// fmt.Println("Printer:", printerName)

	// 2. Connect to selected printer
	printer, err := escpos.NewWindowsPrinter(printerName)
	if err != nil {
		log.Fatal(err)
	}
	defer printer.Close()

	// 3. Print Initialize
	printer.Initialize()

	// 4. Prepare printed fields
	validityTime := FormatTime(voucher.Validity)
	number := strconv.FormatInt(voucher.ID, 10)

	// 5. Build ticket items
	items := []Item{
		{"Wi-Fi:", os.Getenv("ORGANISATION_NAME")},
		{"Username:", voucher.Username},
		{"Validity:", validityTime},
		{"Ticket No:", number},
		{"VGroup:", voucher.GroupName},
	}

	// 6. Generate QR
	qrURL := os.Getenv("CAPTIVE_PORTAL_URL") + "/?username=" + voucher.Username
	qr, err := qrcode.New(qrURL, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("qr code generation failed: %w", err)
	}
	img := qr.Image(150)
	image := Image(img)

	// 7. Print header
	printer.Justify(escpos.CenterJustify)
	printer.SetCharacterSize(1, 0)
	printer.SetBold(true)
	printer.Println("HOTSPOT VOUCHER")
	printer.SetBold(false)
	// 8. Print QR image
	printer.Write([]byte(image))

	// 9. Print item list
	printer.SetCharacterSize(0, 0)
	printer.Justify(escpos.LeftJustify)

	for _, item := range items {
		printer.Println(itemToString(item))
	}

	// 10. Footer
	printer.Justify(escpos.CenterJustify)
	printer.Println(FormatDate(time.Now()))

	// 11. Cut paper
	printer.FeedLines(4)
	printer.Cut()

	return nil
}

func GetPrinterName() []string {
	// 1. List available printers
	printers, err := escpos.GetInstalledPrinters()

	// for i, printer := range printers {
	// 	fmt.Printf("Printer %d: %s\n", i, printer)
	// }
	if err != nil {
		log.Fatal(err)
	}

	return printers
}
