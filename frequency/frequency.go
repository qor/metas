package frequency

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
)

func init() {
	admin.RegisterViewPath("github.com/qor/metas/frequency/views")
}

// Frequencier frequencier interface
type Frequencier interface {
	GetFrequency() *Frequency
}

// Frequency frequency struct
type Frequency struct {
	ScheduledStartAt *time.Time
	ScheduledEndAt   *time.Time
	IntervalUnit     string // hour, day, week, month
	Interval         *int
}

// Scan scan frequency value
func (frequency *Frequency) Scan(value interface{}) error {
	switch data := value.(type) {
	case []byte:
		return json.Unmarshal(data, frequency)
	case string:
		return frequency.Scan([]byte(data))
	case []string:
		for _, str := range data {
			if err := frequency.Scan([]byte(str)); err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported data")
	}
	return nil
}

// Value get value of frequency
func (frequency Frequency) Value() (driver.Value, error) {
	return json.Marshal(frequency)
}

// GetFrequency get frequency
func (frequency Frequency) GetFrequency() *Frequency {
	return &frequency
}

// ConfigureQorMeta configure qor meta
func (frequency Frequency) ConfigureQorMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*admin.Meta); ok {
		meta.Type = "frequency"
		meta.SetSetter(func(res interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			var (
				intervalUnit   = utils.ToString(metaValue.MetaValues.Get("IntervalUnit").Value)
				activeMonths   = utils.ToString(metaValue.MetaValues.Get("ActiveMonths").Value)
				monthDay, err1 = strconv.Atoi(utils.ToString(metaValue.MetaValues.Get("MonthDay").Value))
				weekDay, err2  = strconv.Atoi(utils.ToString(metaValue.MetaValues.Get("WeekDay").Value))
			)

			if err1 != nil || err2 != nil {
				context.AddError(err1, err2)
				return
			}

			today := time.Now()
			if frequencier, ok := (res).(Frequencier); ok {
				frequency := frequencier.GetFrequency()
				frequency.IntervalUnit = intervalUnit
				one := 1
				switch intervalUnit {
				case "once":
				case "daily":
					frequency.Interval = &one
				case "weekly":
					frequency.Interval = &one

					if int(today.Weekday()) > weekDay {
						since := today.AddDate(0, 0, weekDay-int(today.Weekday()))
						frequency.ScheduledStartAt = &since
					} else if int(today.Weekday()) > weekDay {
						since := today.AddDate(0, 0, 7-int(today.Weekday()))
						frequency.ScheduledStartAt = &since
					}
				case "monthly":
					frequency.Interval = &one
					if time.Now().Day() <= monthDay {
						since := today.AddDate(0, 0, monthDay-time.Now().Day())
						frequency.ScheduledStartAt = &since
					} else {
						since := now.New(today).EndOfMonth().Add(time.Second).AddDate(0, 0, monthDay)
						frequency.ScheduledStartAt = &since
					}
				}

				switch activeMonths {
				case "1":
					end := today.AddDate(0, 1, 0)
					frequency.ScheduledEndAt = &end
				case "3":
					end := today.AddDate(0, 3, 0)
					frequency.ScheduledEndAt = &end
				case "6":
					end := today.AddDate(0, 6, 0)
					frequency.ScheduledEndAt = &end
				case "12":
					end := today.AddDate(0, 12, 0)
					frequency.ScheduledEndAt = &end
				}

				reflect.Indirect(reflect.ValueOf(res)).FieldByName(meta.FieldName).Set(reflect.ValueOf(*frequency))
			}
		})
	}
}
