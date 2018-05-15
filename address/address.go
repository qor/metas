package address

import "github.com/qor/admin"

func init() {
	admin.RegisterViewPath("github.com/qor/metas/address/views")
}

// Address address picker
type Address struct {
	Address string
}
