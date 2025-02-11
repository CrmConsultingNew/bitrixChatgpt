package torgi_new

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartTorgi(w http.ResponseWriter, r *http.Request) {
	/*torgi_online.TestTorgi()
	return*/

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Генерация уникального идентификатора лога
	logTraceID := fmt.Sprintf("p%d-%d", time.Now().Unix(), 100+RandInt(0, 899))

	// 1. Получаем список компаний и пользователей из Bitrix24
	arrayAllComp, arUsers := GetCompaniesAndUsers()
	if len(arrayAllComp) == 0 {
		log.Println("Нет компаний для обработки.")
		return
	}

	// 2. Загружаем уже обработанные закупки из файла hrefs.json
	allOrdersHrefs := LoadHrefs(HrefFile)

	// 3. Обрабатываем компании (парсим сайт, ищем закупки)
	allOrdersHrefs, newCompaniesListArr := ProcessCompanies(arrayAllComp, allOrdersHrefs, arUsers, logTraceID)

	// 4. Сохраняем обновленные закупки в файл hrefs.json
	SaveHrefs(HrefFile, allOrdersHrefs)

	// 5. Отправляем результаты в чат Bitrix24
	ProcessNewCompanies(time.Now().Format("02.01.2006"), len(newCompaniesListArr), allOrdersHrefs, newCompaniesListArr)

	fmt.Println("✅ Обработка завершена. Данные сохранены в hrefs.json и отправлены в Bitrix24.")

	//torgi_new.GetCompaniesAndUsers()
	//torgi_new.ProcessCompanies()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ежедневный отчет успешно отправлен"))

}

func StartTorgiDebug() {
	/*torgi_online.TestTorgi()
	return*/
	// Генерация уникального идентификатора лога
	logTraceID := fmt.Sprintf("p%d-%d", time.Now().Unix(), 100+RandInt(0, 899))

	// 1. Получаем список компаний и пользователей из Bitrix24
	arrayAllComp, arUsers := GetCompaniesAndUsers()
	if len(arrayAllComp) == 0 {
		log.Println("Нет компаний для обработки.")
		return
	}

	// 2. Загружаем уже обработанные закупки из файла hrefs.json
	allOrdersHrefs := LoadHrefs(HrefFile)

	// 3. Обрабатываем компании (парсим сайт, ищем закупки)
	allOrdersHrefs, newCompaniesListArr := ProcessCompanies(arrayAllComp, allOrdersHrefs, arUsers, logTraceID)

	// 4. Сохраняем обновленные закупки в файл hrefs.json
	SaveHrefs(HrefFile, allOrdersHrefs)

	// 5. Отправляем результаты в чат Bitrix24
	ProcessNewCompanies(time.Now().Format("02.01.2006"), len(newCompaniesListArr), allOrdersHrefs, newCompaniesListArr)

	fmt.Println("✅ Обработка завершена. Данные сохранены в hrefs.json и отправлены в Bitrix24.")

	//torgi_new.GetCompaniesAndUsers()
	//torgi_new.ProcessCompanies()
}
