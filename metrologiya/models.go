package metrologiya

// Константы для ASSIGNED_BY_ID (ID сотрудников)
const (
	BakhtegareevDenis  = "17"    // Бахтегареев Денис +
	IslamgareevLinar   = "37"    // Исламгареев Линар +
	TelyakovRadik      = "103"   // Теляков Радик +
	ShamsutdinovRim    = "63"    // Шамсутдинов Рим +
	ShaykhutdinovIlgiz = "99"    // Шайхутдинов Ильгизар +
	SayfutdiyarovTimur = "129"   // Сайфутдияров Тимур +
	BuranbaevMansur    = "127"   // Буранбаев Мансур +
	KlimenkoIlya       = "135"   // Клименко Илья +
	ShaniyazovArtur    = "8515"  // Шаниязов Артур +
	ZaripovSiren       = "9347"  // Зарипов Сирен
	YermakovArtem      = "13595" // Ермаков Артем +
	AzamatAkhmadiev    = "34373" // Азамат +
)

// UF_CRM_1650442073068 - дата выезда

// Структура для полей UF_CRM_ (кастомные поля)
type UFCRMFields struct {
	GasMeterVerification        string `json:"UF_CRM_1717764898"`    // Поверка газ. счетчика
	GasBoilerDiagnostics        string `json:"UF_CRM_1720091191105"` // Диагностика газ. котла
	GasStoveDiagnostics         string `json:"UF_CRM_1717765382"`    // Диагностика газовой плиты
	WaterMeterVerification      string `json:"UF_CRM_1717765494"`    // Поверка водяного счетчика
	IndustrialBoilerDiagnostics string `json:"UF_CRM_1729173539594"` // Диагностика промышленных котлов
	SalesHose                   string `json:"UF_CRM_1717765416"`    // Продажа шланга
	SalesGasMeter               string `json:"UF_CRM_1717765660"`    // Продажа газового счетчика
	SalesWaterMeter             string `json:"UF_CRM_1718280444"`    // Продажа водяного счетчика
	QRCode                      string `json:"UF_CRM_1725598657143"` // QR-код
	NonCashTransfer             string `json:"UF_CRM_1718021528233"` // Безналичный перевод
	CompanySum                  string `json:"UF_CRM_1716890120764"` // Сумма компании
	TechnicianPercentage        string `json:"UF_CRM_1716890078180"` // Процент техника поверителя
	GSM                         string `json:"UF_CRM_1716890031608"` // ГСМ
	ApartmentDaily              string `json:"UF_CRM_1717765798"`    // Квартира за сутки
	MiscellaneousExpenses       string `json:"UF_CRM_1717765861"`    // Прочие расходы
}
