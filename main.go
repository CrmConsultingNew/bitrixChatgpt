package main

import (
	"bitrix_app/backend/bitrix/endpoints"
	"bitrix_app/backend/mail"
	"bitrix_app/backend/routes"
	"bitrix_app/medi"
	"bitrix_app/metrologiya"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {

	if err := godotenv.Load(filepath.Join(".env")); err != nil {
		log.Print("No .env file found")
	} else {
		fmt.Println("Loaded .env file")
	}

	//OpenAI.ChatGptTest()

	_, err := metrologiya.FetchUserIDsByDepartment()
	if err != nil {
		log.Printf("Error fetching user IDs: %v", err)
	}

	// SCHEDULERS
	metrologiya.StartScheduler()
	medi.StartMediScheduler()

	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(filepath.Join(".env")); err != nil {
		log.Print("No .env file found")
	} else {
		fmt.Println("Loaded .env file")
	}

	// Установка домена для Bitrix24
	endpoints.BitrixDomain = os.Getenv("BITRIX_DOMAIN")
	endpoints.NewBitrixDomain = os.Getenv("NEW_BITRIX_DOMAIN")

	// Инициализация маршрутов
	routes.Router()

	smtpConfig := mail.SMTPConfig{
		Host:     "smtp.beget.com",
		Port:     "465",
		Username: "crm@crmconsulting-api.ru",
		Password: "AA1379657aa!",
		From:     "crm@crmconsulting-api.ru",
	}

	// Проверка SMTP соединения
	if err := mail.TestSMTPConnection(smtpConfig); err != nil {
		log.Fatalf("SMTP connection failed: %v", err)
	} else {
		log.Println("SMTP connection successful")
	}

	// Запуск сервера
	server := &http.Server{
		Addr:              ":9090",
		ReadHeaderTimeout: 3 * time.Second,
	}

	fmt.Printf("server started on addr: %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Server started with error")
		panic(err)
	}
}

/*var tableDataLeft [][]string

// Добавление данных в таблицу Left с заполненными дебетом и кредитом
office.AddToTableDataLeft(&tableDataLeft, "1", "Оплата по счету от 15.01.2024", "", "1000")
office.AddToTableDataLeft(&tableDataLeft, "2", "Возврат средств от 20.02.2024", "", "")
office.AddToTableDataLeft(&tableDataLeft, "3", "Акт выполненных работ №123 от 10.03.2024", "", "1000")

// Вызов функции StartWord с новыми аргументами, включая второй компанию
err := office.StartWord(
	"01.01.2024",       // dateFrom
	"31.12.2024",       // dateTo
	"15.01.2024",       // dateForInvoice
	"20.02.2024",       // dateForReturnPayment
	"123",              // numberOfCompletedWorks
	"10.03.2024",       // dateOfCompletedWorks
	"31.12.2024",       // dateForSaldo
	"10000",            // sumOfSaldo
	"ООО «Компания 2»", // secondCompany
	tableDataLeft,
)
if err != nil {
	log.Fatalf("Error creating Word document: %v", err)
	return
}*/

//agency.xlsx

type DealsList struct {
	ID string `json:"ID"`
}

func GetDealsList() {

}

func UpdateDeals() {

}

//crm.company.list "filter:NAME" - return IDs
//crm.deal.list - "filter: >UF_CRM_1725716514219:0 ""return IDs
//crm.deal.update "UF_CRM_1734368465" : "COMPANY_ID"

// ProcessAgencyFile обрабатывает файл agency.xlsx и возвращает map[int]string с данными
func ProcessAgencyFile(filePath string) (map[int]string, error) {
	// Открываем файл
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Переходим на первый лист (или укажите имя листа явно)
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheet found in file")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error reading rows from sheet '%s': %v", sheetName, err)
	}

	// Ищем индексы заголовков
	var idColumnIndex, agencyColumnIndex int = -1, -1
	for colIndex, cell := range rows[0] {
		switch strings.TrimSpace(strings.ToLower(cell)) {
		case "id":
			idColumnIndex = colIndex
		case "агентство недвижимости":
			agencyColumnIndex = colIndex
		}
	}

	if idColumnIndex == -1 || agencyColumnIndex == -1 {
		return nil, fmt.Errorf("required columns 'id' or 'агентство недвижимости' not found")
	}

	// Создаём map для хранения данных
	result := make(map[int]string)

	// Обрабатываем строки, начиная со второй (первая — заголовок)
	for _, row := range rows[1:] {
		// Проверяем, чтобы индексы не выходили за границы строки
		if len(row) <= idColumnIndex || len(row) <= agencyColumnIndex {
			continue
		}

		idValue := strings.TrimSpace(row[idColumnIndex])
		agencyValue := strings.TrimSpace(row[agencyColumnIndex])

		if idValue == "" || agencyValue == "" {
			continue
		}

		// Преобразуем ID в int
		id, err := strconv.Atoi(idValue)
		if err != nil {
			log.Printf("Skipping invalid ID '%s': %v", idValue, err)
			continue
		}

		// Удаляем возможный префикс "АН " из названия агентства
		agency := strings.TrimPrefix(agencyValue, "АН ")
		result[id] = agency
	}

	return result, nil
}
