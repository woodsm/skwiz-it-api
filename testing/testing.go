package testing

import "../notification"

func TestEmail() {
	notification.SendEmail("ben@krashidbuilt.com", 80)
}
