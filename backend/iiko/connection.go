package iiko

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

//f95fa107-04ab-1353-24bd-e82186a35d21
//getDepartmentsIDS : : : https://harizma-co.iiko.it:443/resto/api/corporation/departments?key=ae8d5fbe-9235-6f8d-14e9-263f4ad28bc3
//getOpportunityReport : : : https://harizma-co.iiko.it:443/resto/api/reports/sales?key=ae8d5fbe-9235-6f8d-14e9-263f4ad28bc3&dateFrom=17.04.2024&dateTo=18.04.2024&dishDetails=true&allRevenue=false&department=969e5b39-5306-0777-0189-306abd4d0011
//getOLAP_report : : : https://harizma-co.iiko.it:443/resto/api/reports/olap?key=f95fa107-04ab-1353-24bd-e82186a35d21&report=SALES&from=18.04.2024&to=19.04.2024&groupRow=WaiterName&groupRow=CloseTime&groupRow=Department&groupRow=averageByGuest&agr=fullSum&agr=OrderNum&groupRow=HourClose&groupRow=OperationType&groupRow=OrderType&groupRow=SessionNum&groupRow=TableNum

var IikoApiKey string

func GetSales() {
	//https://harizma-co.iiko.it:443/resto/api/reports/olap?key=90568e67-cce3-183f-c9bf-7740a920bcd4&report=SALES&from=01.04.2024&to=01.04.2024&groupRow=Department&groupRow=OpenTime&agr=fullSum&departmentId=79963b6f-2e77-4906-ba1d-ee761d558340

	dateStart := "01.04.2024"
	dateFinish := "01.04.2024"
	uri := fmt.Sprintf("https://harizma-co.iiko.it:443/resto/api/reports/olap?key=%s&report=SALES&from=%s&to=%s&groupRow=Department&groupRow=OpenTime&agr=fullSum", GlobalToken, dateStart, dateFinish)

	method := "GET"
	body := []byte(``)
	req, err := http.NewRequest(method, uri, bytes.NewReader(body))
	if err != nil {
		log.Println("Error creating request GetUsersInfo:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request GetUsersInfo:", err)
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body GetUsersInfo:", err)
		return
	}
	log.Println("Response body BS GetUsersInfo:", string(bs))

	var report Report
	err = xml.Unmarshal(bs, &report)
	if err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	// Output results
	for _, record := range report.Records {
		fmt.Printf("OpenTime: %s, Department: %s, FullSum: %s\n", record.OpenTime, record.Department, record.FullSum)
	}

}

var GlobalToken string

func IikoGetToken() {
	//appKey := os.Getenv("TUYA_APP_KEY")
	//todo page_no must be count, no more than 10
	//uri := fmt.Sprintf("/v1.0/apps/%v/users?page_no=%v&page_size=100&access_token=%v&sign=&t=", appKey, page, AccessToken)

	password := "12345"
	h := sha1.New()
	h.Write([]byte(password))
	sha1PasswordHash := hex.EncodeToString(h.Sum(nil))

	uri := fmt.Sprintf("https://harizma-co.iiko.it:443/resto/api/auth?login=dsannikovharizmaclub&pass=%s", sha1PasswordHash)

	method := "GET"
	body := []byte(``)
	req, err := http.NewRequest(method, uri, bytes.NewReader(body))
	if err != nil {
		log.Println("Error creating request GetUsersInfo:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request GetUsersInfo:", err)
		return
	}
	defer resp.Body.Close()

	//var token string

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body GetUsersInfo:", err)
		return
	}
	log.Println("Response body BS GetUsersInfo:", string(bs))
	//err = json.Unmarshal(bs, &token)
	IikoApiKey = string(bs)
	//fmt.Println("GlobalToken: ", GlobalToken)
}

type Report struct {
	XMLName xml.Name `xml:"report"`
	Records []Record `xml:"r"`
}

type Record struct {
	OpenTime   string `xml:"OpenTime"`
	Department string `xml:"Department"`
	FullSum    string `xml:"fullSum"`
}
