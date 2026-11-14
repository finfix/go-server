package service

import (
	"context"
	"time"

	"server/internal/modules/statistic/model"
	"server/internal/utils/errors"
)

// Ключи рядов без разбивки
const (
	keyIncome   = "income"   // Доходы
	keyExpense  = "expense"  // Расходы
	keyDelta    = "delta"    // Доходы минус расходы
	keyBalance  = "balance"  // В наличии одной линией
	keyUntagged = "untagged" // Транзакции без тегов при разбивке по тегам
)

// GetStatistic считает ряды графика по бакетам периода. Правила как в приложении: расход — сумма
// пополнения расходных счетов, доход — сумма списания доходных счетов; наличность — счета regular и
// debt, остаток собирается от первой транзакции. Фильтр типов транзакций к графику не применяется
func (s *StatisticService) GetStatistic(ctx context.Context, req model.GetStatisticReq) (res model.Statistic, err error) {
	ctx, span := tracer.Start(ctx, "GetStatistic")
	defer span.End()

	// Проверяем параметры
	if err = req.ChartType.Validate(ctx); err != nil {
		return res, err
	}
	if err = req.Period.Validate(ctx); err != nil {
		return res, err
	}
	if req.GroupBy == "" {
		req.GroupBy = model.GroupByAccount
	}
	if err = req.GroupBy.Validate(ctx); err != nil {
		return res, err
	}

	// Бакеты периода
	res.DateTo = startOf(time.Now(), model.PeriodDay)
	if req.Filters.DateTo != nil {
		res.DateTo = startOf(*req.Filters.DateTo, model.PeriodDay)
	}
	res.DateFrom = previousBuckets(startOf(res.DateTo, req.Period), req.Period, defaultBucketCount[req.Period]-1)
	if req.Filters.DateFrom != nil {
		res.DateFrom = startOf(*req.Filters.DateFrom, req.Period)
	}
	if res.DateFrom.After(res.DateTo) {
		return res, errors.BadRequest.New("date_from is after date_to").WithCustomHumanText("Дата начала позже даты конца")
	}
	for bucket := res.DateFrom; !bucket.After(res.DateTo); bucket = nextBucket(bucket, req.Period) {
		res.Buckets = append(res.Buckets, bucket)
		if len(res.Buckets) > maxBuckets {
			return res, errors.BadRequest.New("too many buckets").
				WithCustomHumanText("Слишком много точек: увеличьте период или сократите интервал дат")
		}
	}
	endOfPeriod := res.DateTo.AddDate(0, 0, 1)
	req.Filters.DateFrom, req.Filters.DateTo = &res.DateFrom, &res.DateTo

	// Данные выбранных групп
	data, err := s.loadDataset(ctx, req.Necessary, req.Filters.AccountGroupIDs, req.Currency)
	if err != nil {
		return res, err
	}
	res.Currency = data.currency

	// Считаем ряды по виду графика
	calculator := seriesCalculator{data: data, buckets: res.Buckets, endOfPeriod: endOfPeriod, aggregateIntoParents: req.AggregateIntoParents, series: map[string]*model.Series{}}
	filter := data.newTransactionFilter(req.Filters, true)
	switch req.ChartType {
	case model.ChartIncome, model.ChartExpense:
		if err = calculator.addFlows(filter, req.ChartType, req.GroupBy); err != nil {
			return res, err
		}
		if req.GroupBy == model.GroupByTag {
			if err = s.nameTagSeries(ctx, req, data, calculator.series); err != nil {
				return res, err
			}
		}
		if req.ChartType == model.ChartExpense && req.GroupBy == model.GroupByAccount {
			if err = s.addBudgets(ctx, req, data, calculator, res.DateTo); err != nil {
				return res, err
			}
		}
	case model.ChartIncomeExpense, model.ChartDelta:
		if err = calculator.addFlows(filter, req.ChartType, ""); err != nil {
			return res, err
		}
	case model.ChartBalanceDetailed, model.ChartBalanceTotal:
		if err = calculator.addBalances(req.Filters, req.ChartType == model.ChartBalanceTotal); err != nil {
			return res, err
		}
	}

	res.Series = calculator.result(req.ChartType)
	return res, nil
}
