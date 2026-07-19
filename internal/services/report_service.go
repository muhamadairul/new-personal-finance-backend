package services

import (
	"sort"
	"time"

	"finance-app-backend/internal/models"
	"finance-app-backend/internal/repositories"
	"finance-app-backend/internal/resources"

	"gorm.io/gorm"
)

type ReportService struct {
	db         *gorm.DB
	walletRepo repositories.WalletRepositoryInterface
	txRepo     repositories.TransactionRepositoryInterface
}

func NewReportService(db *gorm.DB, walletRepo repositories.WalletRepositoryInterface, txRepo repositories.TransactionRepositoryInterface) *ReportService {
	return &ReportService{
		db:         db,
		walletRepo: walletRepo,
		txRepo:     txRepo,
	}
}

func (s *ReportService) GetDashboardData(userID uint, appURL string) (map[string]interface{}, error) {
	now := time.Now()

	// 1. Wallets & Total Balance
	wallets, err := s.walletRepo.GetForUser(userID)
	if err != nil {
		return nil, err
	}

	totalBalance := 0.0
	for _, w := range wallets {
		totalBalance += w.Balance
	}

	// 2. Monthly Income
	var monthlyIncome float64
	s.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND type = 'income' AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", userID, int(now.Month()), now.Year()).
		Row().Scan(&monthlyIncome)

	// 3. Monthly Expense
	var monthlyExpense float64
	s.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND type = 'expense' AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", userID, int(now.Month()), now.Year()).
		Row().Scan(&monthlyExpense)

	// 4. Weekly Expenses (last 7 days)
	weeklyExpenses := make([]float64, 7)
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		endOfDay := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 999999999, day.Location())

		var dayExpense float64
		s.db.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("user_id = ? AND type = 'expense' AND date >= ? AND date <= ?", userID, startOfDay, endOfDay).
			Row().Scan(&dayExpense)

		weeklyExpenses[6-i] = dayExpense
	}

	// 5. Recent Transactions (last 5)
	var recentTxs []models.Transaction
	s.db.Where("user_id = ?", userID).
		Preload("Category").
		Preload("Wallet").
		Order("date DESC").
		Order("created_at DESC").
		Limit(5).
		Find(&recentTxs)

	return map[string]interface{}{
		"total_balance":       totalBalance,
		"monthly_income":      monthlyIncome,
		"monthly_expense":     monthlyExpense,
		"weekly_expenses":     weeklyExpenses,
		"recent_transactions": resources.ToTransactionCollection(recentTxs),
		"wallets":             resources.ToWalletCollection(wallets),
	}, nil
}

func (s *ReportService) GetMonthlyReport(userID uint, month, year int) (map[string]interface{}, error) {
	targetDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	monthlyIncome := make([]float64, 6)
	monthlyExpense := make([]float64, 6)
	monthLabels := make([]string, 6)

	monthNamesIndo := map[time.Month]string{
		time.January:   "Jan",
		time.February:  "Feb",
		time.March:     "Mar",
		time.April:     "Apr",
		time.May:       "Mei",
		time.June:      "Jun",
		time.July:      "Jul",
		time.August:    "Agu",
		time.September: "Sep",
		time.October:   "Okt",
		time.November:  "Nov",
		time.December:  "Des",
	}

	for i := 5; i >= 0; i-- {
		d := targetDate.AddDate(0, -i, 0)
		idx := 5 - i
		monthLabels[idx] = monthNamesIndo[d.Month()]

		var inc float64
		s.db.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("user_id = ? AND type = 'income' AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", userID, int(d.Month()), d.Year()).
			Row().Scan(&inc)

		var exp float64
		s.db.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("user_id = ? AND type = 'expense' AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", userID, int(d.Month()), d.Year()).
			Row().Scan(&exp)

		monthlyIncome[idx] = inc
		monthlyExpense[idx] = exp
	}

	totalIncome := 0.0
	for _, v := range monthlyIncome {
		totalIncome += v
	}

	totalExpense := 0.0
	for _, v := range monthlyExpense {
		totalExpense += v
	}

	return map[string]interface{}{
		"total_income":    totalIncome,
		"total_expense":   totalExpense,
		"net":             totalIncome - totalExpense,
		"monthly_income":  monthlyIncome,
		"monthly_expense": monthlyExpense,
		"month_labels":    monthLabels,
	}, nil
}

type CategoryBreakdownItem struct {
	Name   string
	Amount float64
	Color  int64
}

func (s *ReportService) GetCategoryReport(userID uint, month, year int) (map[string]interface{}, error) {
	var txs []models.Transaction
	err := s.db.Where("user_id = ? AND type = 'expense' AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", userID, month, year).
		Preload("Category").
		Find(&txs).Error
	if err != nil {
		return nil, err
	}

	sums := make(map[string]float64)
	colors := make(map[string]int64)

	for _, t := range txs {
		if t.Category != nil {
			sums[t.Category.Name] += t.Amount
			colors[t.Category.Name] = t.Category.Color
		}
	}

	// Sort items descending by amount
	var items []CategoryBreakdownItem
	for name, amt := range sums {
		items = append(items, CategoryBreakdownItem{
			Name:   name,
			Amount: amt,
			Color:  colors[name],
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Amount > items[j].Amount
	})

	orderedBreakdown := make(map[string]float64)
	orderedColors := make(map[string]int64)

	for _, item := range items {
		orderedBreakdown[item.Name] = item.Amount
		orderedColors[item.Name] = item.Color
	}

	return map[string]interface{}{
		"category_breakdown": orderedBreakdown,
		"category_colors":    orderedColors,
	}, nil
}
