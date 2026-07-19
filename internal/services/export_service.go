package services

import (
	"bytes"
	"fmt"
	"time"

	"finance-app-backend/internal/repositories"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ExportService struct {
	db       *gorm.DB
	userRepo repositories.UserRepositoryInterface
	txRepo   repositories.TransactionRepositoryInterface
}

func NewExportService(db *gorm.DB, userRepo repositories.UserRepositoryInterface, txRepo repositories.TransactionRepositoryInterface) *ExportService {
	return &ExportService{
		db:       db,
		userRepo: userRepo,
		txRepo:   txRepo,
	}
}

func (s *ExportService) getMonthNameIndo(month int) string {
	names := map[int]string{
		1: "Januari", 2: "Februari", 3: "Maret", 4: "April",
		5: "Mei", 6: "Juni", 7: "Juli", 8: "Agustus",
		9: "September", 10: "Oktober", 11: "November", 12: "Desember",
	}
	return names[month]
}

func (s *ExportService) GenerateExcel(userID uint, month, year int) ([]byte, string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, "", err
	}

	txs, err := s.txRepo.GetForUser(userID, month, year)
	if err != nil {
		return nil, "", err
	}

	monthName := fmt.Sprintf("%s %d", s.getMonthNameIndo(month), year)
	f := excelize.NewFile()
	sheet := "Transaksi"
	f.SetSheetName("Sheet1", sheet)

	// Title
	f.MergeCell(sheet, "A1", "F1")
	f.SetCellValue(sheet, "A1", fmt.Sprintf("Laporan Keuangan — %s", user.Name))

	f.MergeCell(sheet, "A2", "F2")
	f.SetCellValue(sheet, "A2", monthName)

	// Headers
	headers := []string{"Tanggal", "Tipe", "Kategori", "Dompet", "Nominal", "Catatan"}
	cols := []string{"A", "B", "C", "D", "E", "F"}
	headerRow := 4

	for i, h := range headers {
		f.SetCellValue(sheet, fmt.Sprintf("%s%d", cols[i], headerRow), h)
	}

	// Data rows
	row := headerRow + 1
	var totalIncome, totalExpense float64

	for _, t := range txs {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), t.Date.Format("02/01/2006"))

		typeName := "Pengeluaran"
		if t.Type == "income" {
			typeName = "Pemasukan"
			totalIncome += t.Amount
		} else {
			totalExpense += t.Amount
		}
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), typeName)

		catName := "-"
		if t.Category != nil {
			catName = t.Category.Name
		}
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), catName)

		wName := "-"
		if t.Wallet != nil {
			wName = t.Wallet.Name
		}
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), wName)

		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), t.Amount)

		note := "-"
		if t.Note != nil && *t.Note != "" {
			note = *t.Note
		}
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), note)

		row++
	}

	// Summary section
	row++
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Total Pemasukan")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), totalIncome)

	row++
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Total Pengeluaran")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), totalExpense)

	row++
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Selisih (Net)")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), totalIncome-totalExpense)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("transaksi_%d_%d.xlsx", year, month)
	return buf.Bytes(), filename, nil
}

func (s *ExportService) GeneratePDF(userID uint, month, year int) ([]byte, string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, "", err
	}

	txs, err := s.txRepo.GetForUser(userID, month, year)
	if err != nil {
		return nil, "", err
	}

	var totalIncome, totalExpense float64
	for _, t := range txs {
		if t.Type == "income" {
			totalIncome += t.Amount
		} else {
			totalExpense += t.Amount
		}
	}

	monthName := fmt.Sprintf("%s %d", s.getMonthNameIndo(month), year)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Title
	pdf.CellFormat(190, 10, "Laporan Keuangan", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "I", 11)
	pdf.CellFormat(190, 7, fmt.Sprintf("%s — %s", user.Name, monthName), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Summary Box
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(95, 7, "Total Pemasukan:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Rp %.0f", totalIncome), "1", 1, "R", false, 0, "")

	pdf.CellFormat(95, 7, "Total Pengeluaran:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Rp %.0f", totalExpense), "1", 1, "R", false, 0, "")

	pdf.CellFormat(95, 7, "Selisih (Net):", "1", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Rp %.0f", totalIncome-totalExpense), "1", 1, "R", false, 0, "")
	pdf.Ln(10)

	// Table Headers
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(25, 8, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Tipe", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 8, "Kategori", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 8, "Dompet", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 8, "Nominal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 8, "Catatan", "1", 1, "C", false, 0, "")

	// Table Data
	pdf.SetFont("Arial", "", 9)
	for _, t := range txs {
		catName := "-"
		if t.Category != nil {
			catName = t.Category.Name
		}
		wName := "-"
		if t.Wallet != nil {
			wName = t.Wallet.Name
		}
		tType := "Pengeluaran"
		if t.Type == "income" {
			tType = "Pemasukan"
		}
		note := "-"
		if t.Note != nil && *t.Note != "" {
			note = *t.Note
		}

		pdf.CellFormat(25, 7, t.Date.Format("02/01/2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 7, tType, "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 7, catName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(35, 7, wName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(35, 7, fmt.Sprintf("Rp %.0f", t.Amount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 7, note, "1", 1, "L", false, 0, "")
	}

	// Footer timestamp
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(190, 5, fmt.Sprintf("Dibuat otomatis oleh Pencatat Keuangan — %s", time.Now().Format("02/01/2006 15:04")), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("transaksi_%d_%d.pdf", year, month)
	return buf.Bytes(), filename, nil
}
