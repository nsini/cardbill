/**
 * @Time: 2019/12/15 17:14
 * @Author: solacowa@gmail.com
 * @File: date
 * @Software: GoLand
 */

package date

import "time"

func ParseCardBillAndHolderDay(billingDay, holderDay int) (billingTime, holderTime time.Time) {
	now := time.Now()

	if billingDay < now.Day() {
		billingTime = time.Date(now.Year(), now.Month(), billingDay, 0, 0, 0, 0, time.Local)
	} else {
		billingTime = time.Date(now.Year(), now.Month()-1, billingDay, 0, 0, 0, 0, time.Local)
	}

	if billingDay >= 15 {
		holderTime = time.Date(now.Year(), now.Month()+1, holderDay, 0, 0, 0, 0, time.Local)
	} else {
		holderTime = time.Date(now.Year(), now.Month(), holderDay, 0, 0, 0, 0, time.Local)
	}

	if holderTime.Unix() <= time.Now().Unix() {
		holderTime = time.Date(now.Year(), now.Month()+1, holderDay, 0, 0, 0, 0, time.Local)
	}

	return
}
