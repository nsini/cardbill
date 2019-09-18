/**
 * @Time: 2019-08-18 11:31
 * @Author: solacowa@gmail.com
 * @File: math
 * @Software: GoLand
 */

package transform

import "math"

func Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}
