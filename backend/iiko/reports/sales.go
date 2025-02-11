package reports

import (
	"bitrix_app/backend/iiko"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func IikoReportSales() {
	from := `01.03.2024`
	to := `01.03.2024`

	//departmentId=79963b6f-2e77-4906-ba1d-ee761d558340& -- Бауманская
	url := fmt.Sprintf(`https://harizma-co.iiko.it:443/resto/api/reports/olap?key=%v&report=SALES&from=%v&to=%v&groupRow=Department&groupRow=OpenTime&agr=fullSum&groupCol=OpenDate.Typed&groupCol=OpenTime&groupCol=CloseTime&groupCol=CookingPlace&groupCol=DishName&groupCol=OrderNum`, iiko.IikoApiKey, from, to)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching report:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var report Report
	err = xml.Unmarshal(body, &report)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return
	}
	var count int
	for _, rep := range report.Records {
		openTime, err := parseDateTime(rep.OpenTime)
		if err != nil {
			fmt.Println("Error parsing OpenTime:", err)
			continue
		}
		closeTime, err := parseDateTime(rep.CloseTime)
		if err != nil {
			fmt.Println("Error parsing CloseTime:", err)
			continue
		}
		openDateTime, err := parseDateTime(rep.OpenDateTime)
		if err != nil {
			fmt.Println("Error parsing CloseTime:", err)
			continue
		}

		// Check if the hour part of OpenTime is greater than 14
		if (openTime.Hour() > 14 || openTime.Hour() < 6) && closeTime.Hour() < 6 && rep.Department == "Бауманская Harizma" && strings.Contains(rep.DishName, "Обслуживание") == true && openDateTime.Day() == 01 {
			fmt.Println("Department:", rep.Department)
			fmt.Println("CookingPlace:", rep.CookingPlace)
			fmt.Println("DishName:", rep.DishName)
			fmt.Println("OrderNum:", rep.OrderNum)
			fmt.Println("OpenTime:", rep.OpenTime)
			fmt.Println("CloseTime:", rep.CloseTime)
			fmt.Println("OpenDateTime:", rep.OpenDateTime)
			fmt.Println("FullSum:", rep.FullSum)
			count++
		}

	}
	fmt.Println("Total count:", count)
	//fmt.Printf("Report details: %+v\n", report)
}

func parseDateTime(dateTimeStr string) (time.Time, error) {
	layout := "Mon Jan 02 15:04:05 MST 2006"
	return time.Parse(layout, dateTimeStr)
}
