package daterange

import "time"

// RangeType range type
type RangeType = string

var (
	// Today today range type
	Today RangeType = "today"
	// Yesterday yesterday
	Yesterday RangeType = "yesterday"
	// LastWeek last week
	LastWeek RangeType = "last_week"
	// LastMonth last month
	LastMonth RangeType = "last_month"
	// Last7Days last 7 days
	Last7Days RangeType = "last_7_days"
	// Last30Days last 30 days
	Last30Days RangeType = "last_30_days"
)

// DateRange date range type
type DateRange struct {
	Type  RangeType
	From  *time.Time
	Until *time.Time
}
