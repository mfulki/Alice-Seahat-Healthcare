package utils

import (
	"context"
	"fmt"
	"strings"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/sirupsen/logrus"
)

type PaginateQuery struct {
	SelectColumns string
	AdvanceQuery  string
	OrderQuery    string
}

func BuildSortQuery(alias map[string]string, sort, sortDefault string) string {
	if sort == "" {
		return sortDefault
	}

	sortSplit := strings.Split(sort, ",")
	if len(sortSplit) != 2 {
		return sortDefault
	}

	orderColumn := sortSplit[0]
	orderBy := sortSplit[1]

	if !SliceIsContain([]string{"asc", "desc"}, orderBy) {
		return sortDefault
	}

	val, ok := alias[orderColumn]
	if !ok {
		return sortDefault
	}

	return fmt.Sprintf("%s %s", val, orderBy)
}

func BuildFilterQuery(alias map[string]string, clc *entity.Collection, filterDefault string) string {
	builder := new(strings.Builder)
	builder.WriteString(filterDefault)
	argsLen := len(clc.Args)

	for key, val := range clc.Filter {
		if val == "" {
			continue
		}

		column, ok := alias[key]
		if !ok {
			continue
		}

		argsLen++
		clc.Args = append(clc.Args, val)
		builder.WriteString(fmt.Sprintf(" AND %s = $%d", column, argsLen))
	}

	return builder.String()
}

func BuildSearchQuery(columns []string, clc *entity.Collection) string {
	search := clc.Search
	argsLen := len(clc.Args)

	if search == "" {
		return search
	}

	builder := new(strings.Builder)
	colLength := len(columns)

	for index, col := range columns {
		if index == 0 {
			builder.WriteString("AND (")
		} else {
			builder.WriteString(" OR ")
		}

		argsLen++
		clc.Args = append(clc.Args, "%"+search+"%")
		builder.WriteString(fmt.Sprintf("%s ILIKE $%d", col, argsLen))

		if index == colLength-1 {
			builder.WriteString(")")
		}
	}

	return builder.String()
}

func BuildQuery(db transaction.DBTransaction, pq PaginateQuery, collection *entity.Collection) string {
	return buildQuery(db, pq, collection, nil)
}

func BuildQueryWithGroupBy(db transaction.DBTransaction, pq PaginateQuery, collection *entity.Collection) string {
	return buildQuery(db, pq, collection, func(s string) string {
		return fmt.Sprintf("SELECT COUNT(*) FROM (%s)", s)
	})
}

func buildQuery(db transaction.DBTransaction, pq PaginateQuery, collection *entity.Collection, countFunc func(string) string) string {
	commonQuery := "SELECT " + pq.SelectColumns + " FROM " + pq.AdvanceQuery
	query := commonQuery + " ORDER BY " + pq.OrderQuery

	if collection.Page == 0 || collection.Limit == 0 {
		return query
	}

	ctx := context.Background()
	countQuery := "SELECT COUNT(*) FROM " + pq.AdvanceQuery
	if countFunc != nil {
		countQuery = countFunc(countQuery)
	}

	if err := db.QueryRowContext(ctx, countQuery, collection.Args...).Scan(&collection.TotalRecords); err != nil {
		logrus.Error(err)
		return ""
	}

	currentOffset := (collection.Page - 1) * collection.Limit
	pgQuery := fmt.Sprintf(" LIMIT %d OFFSET %d", collection.Limit, currentOffset)
	return query + pgQuery
}
