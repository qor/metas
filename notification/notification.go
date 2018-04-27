package notification

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/qor/admin"
	"github.com/qor/qor/resource"
)

func init() {
	admin.RegisterViewPath("github.com/qor/metas/notification/views")
}

// Notification notification struct
type Notification struct {
	Emails []string
}

// Config notification config
type Config struct {
	Subject       string
	Template      string
	DefaultEmails []string
}

// Scan scan notification value
func (notification *Notification) Scan(value interface{}) error {
	switch data := value.(type) {
	case []byte:
		return json.Unmarshal(data, notification)
	case string:
		return notification.Scan([]byte(data))
	case []string:
		for _, str := range data {
			if err := notification.Scan([]byte(str)); err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported data")
	}
	return nil
}

// Value get value of notification
func (notification Notification) Value() (driver.Value, error) {
	return json.Marshal(notification)
}

// ConfigureQorMeta configure qor meta
func (notification Notification) ConfigureQorMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*admin.Meta); ok {
		meta.Type = "notification"
	}
}
