/**
 * @Time : 11/12/2019 11:56 AM
 * @Author : solacowa@gmail.com
 * @File : response
 * @Software: GoLand
 */

package record

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/nsini/cardbill/src/repository/types"
	"net/http"
	"strconv"
	"time"
)

type ExportResponse struct {
	Err  error                   `json:"err"`
	Data []*types.ExpensesRecord `json:"data"`
}

func (r ExportResponse) Failed() error {
	return r.Err
}

type Failure interface {
	Failed() error
}

func encodeExportResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	resp := response.(ExportResponse).Data

	f := excelize.NewFile()

	var sheetName = "消费记录"

	index := f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", "时间")
	f.SetCellValue(sheetName, "B1", "银行")
	f.SetCellValue(sheetName, "C1", "信用卡名称")
	f.SetCellValue(sheetName, "D1", "信用卡后四位")
	f.SetCellValue(sheetName, "E1", "商户类型")
	f.SetCellValue(sheetName, "F1", "商户类型码")
	f.SetCellValue(sheetName, "G1", "商户名称")
	f.SetCellValue(sheetName, "H1", "费率")
	f.SetCellValue(sheetName, "I1", "金额")
	f.SetCellValue(sheetName, "J1", "到账")

	var k = 1
	for _, v := range resp {
		k++
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", k), v.CreatedAt.String())
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", k), v.CreditCard.Bank.BankName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", k), v.CreditCard.CardName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", k), v.CreditCard.TailNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", k), v.Business.BusinessName+" - "+strconv.Itoa(int(v.Business.Code)))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", k), strconv.Itoa(int(v.Business.Code)))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", k), v.BusinessName)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", k), v.Rate)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", k), v.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", k), v.Arrival)
	}

	f.SetActiveSheet(index)

	var fileName = fmt.Sprintf("card-bill-%s.xlsx", time.Now().Format("2006-01-02"))

	w.Header().Add("Accept-Ranges", "bytes")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Add("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Expires", "0")

	_ = f.Write(w)

	return

}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	_ = json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	return http.StatusOK
}

type errorWrapper struct {
	Error string `json:"error"`
}
