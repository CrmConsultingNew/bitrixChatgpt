package office

import (
	"fmt"
	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"log"
)

func AddToTableDataLeft(tableDataLeft *[][]string, numberOfDocument, documentOperation, debit, credit string) {
	*tableDataLeft = append(*tableDataLeft, []string{numberOfDocument, documentOperation, debit, credit})
}

var isLicenseSet = false

func initializeLicense() error {
	if !isLicenseSet {
		err := license.SetMeteredKey("2bee95a134aa3dc14e8e4931aaa8f5cbf643d12fa8ac7397bdd5707f77d5700d")
		if err != nil {
			return fmt.Errorf("ошибка установки лицензии: %v", err)
		}
		isLicenseSet = true
		log.Println("Лицензия установлена")
	}
	return nil
}

func StartWord(dateFrom, dateTo, dateForInvoice, dateForReturnPayment, numberOfCompletedWorks, dateOfCompletedWorks, dateForSaldo, sumOfSaldo, secondCompany string, tableDataLeft [][]string) error {
	log.Println("StartWord was called")

	if err := initializeLicense(); err != nil {
		return err
	}

	doc := document.New()
	defer doc.Close()

	// Заголовок документа
	addCenteredBoldLargeParagraph(doc, "Акт сверки")
	doc.AddParagraph()
	addCenteredParagraph(doc, fmt.Sprintf("взаимных расчетов по состоянию с %s по %s \nмежду ИП Юмагужиной Ралией Махмутовной и %s\n\n. Мы, нижеподписавшиеся, ИП Юмагужина Ралия Махмутовна, с одной стороны, %s, с другой стороны, составили настоящий акт сверки в том, что состояние взаимных расчетов по данным учета следующее:", dateFrom, dateTo, secondCompany, secondCompany))
	doc.AddParagraph()

	createReconciliationTable(doc, dateForInvoice, dateForReturnPayment, numberOfCompletedWorks, dateOfCompletedWorks, dateForSaldo, sumOfSaldo, tableDataLeft)
	doc.AddParagraph()

	textUnderTable := fmt.Sprintf("По данным \n ИП Юмагужина Р.М. на %s", dateTo)

	// Текст ниже таблицы
	addLeftAlignedParagraph(doc, textUnderTable)

	totalDebit, totalCredit := calculateTotals(tableDataLeft)

	totalSum := totalDebit - totalCredit

	generateTotalSumText(doc, dateTo, totalSum, secondCompany)
	doc.AddParagraph()

	addSingleRowWithTwoTexts(doc, "От ИП Юмагужиной Р.М.", "От "+secondCompany)
	doc.AddParagraph()

	addSignaturesRow(doc)
	doc.AddParagraph()

	if err := doc.Validate(); err != nil {
		log.Printf("Ошибка при валидации: %v", err)
		return fmt.Errorf("ошибка при валидации: %v", err)
	}

	if err := doc.SaveToFile("/root/bitrixChatgpt/tables.docx"); err != nil {
		log.Printf("Ошибка при сохранении файла: %v", err)
		return fmt.Errorf("ошибка при сохранении файла: %v", err)
	}

	log.Println("Документ успешно сохранен как /tables.docx")
	log.Println("StartWord was ended")
	return nil
}

func createReconciliationTable(doc *document.Document, dateForInvoice, dateForReturnPayment, numberOfCompletedWorks, dateOfCompletedWorks, dateForSaldo, sumOfSaldo string, tableDataLeft [][]string) {
	table := doc.AddTable()

	table.Properties().SetWidthPercent(100)

	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

	headerRow := table.AddRow()
	headerCellLeft := headerRow.AddCell()
	headerCellLeft.Properties().SetColumnSpan(4)
	addCenteredText(headerCellLeft, "По данным")

	headerCellRight := headerRow.AddCell()
	headerCellRight.Properties().SetColumnSpan(4)
	addCenteredText(headerCellRight, "По данным")

	subHeaderRow := table.AddRow()
	for i := 0; i < 2; i++ {
		addCenteredText(subHeaderRow.AddCell(), "№")
		addCenteredText(subHeaderRow.AddCell(), "Операция")
		addCenteredText(subHeaderRow.AddCell(), "Дебет")
		addCenteredText(subHeaderRow.AddCell(), "Кредит")
	}

	// Заполнение данных для левой части таблицы
	for _, rowLeft := range tableDataLeft {
		row := table.AddRow()
		addCenteredText(row.AddCell(), rowLeft[0]) // Номер документа
		addCenteredText(row.AddCell(), rowLeft[1]) // Название операции
		addCenteredText(row.AddCell(), rowLeft[2]) // Дебет
		addCenteredText(row.AddCell(), rowLeft[3]) // Кредит
		for i := 0; i < 4; i++ {
			addCenteredText(row.AddCell(), "") // Пустые ячейки для правой части таблицы
		}
	}

	// Добавляем строку "Оборот за период"
	totalDebit, totalCredit := calculateTotals(tableDataLeft)
	totalRow := table.AddRow()
	addCenteredText(totalRow.AddCell(), "")
	addCenteredText(totalRow.AddCell(), "Оборот за период")
	addCenteredText(totalRow.AddCell(), fmt.Sprintf("%.0f", totalDebit))
	addCenteredText(totalRow.AddCell(), fmt.Sprintf("%.0f", totalCredit))
	for i := 0; i < 4; i++ {
		addCenteredText(totalRow.AddCell(), "")
	}

	// Вычисление итогового сальдо на дату
	finalSaldo := totalDebit - totalCredit

	// Добавляем строку "Сальдо"
	finalSaldoRow := table.AddRow()
	finalSaldoCell := finalSaldoRow.AddCell()
	finalSaldoCell.Properties().SetColumnSpan(8)
	addLeftAlignedBoldText(finalSaldoCell, fmt.Sprintf("Сальдо на %s %.0f руб.", dateForSaldo, finalSaldo))
}

func generateTotalSumText(doc *document.Document, dateTo string, totalSum float64, secondCompany string) {
	para := doc.AddParagraph()
	para.SetAlignment(wml.ST_JcLeft)
	run := para.AddRun()

	if totalSum > 0 {
		run.AddText(fmt.Sprintf("На %s задолженность в пользу компании ИП Юмагужина Р.М. составляет %.0f руб.", dateTo, totalSum))
	} else if totalSum == 0 {
		run.AddText("Задолженность отсутствует")
	} else {
		run.AddText(fmt.Sprintf("Задолженность в пользу %s составляет %.0f руб.", secondCompany, -totalSum))
	}
}

func calculateTotals(tableDataLeft [][]string) (totalDebit, totalCredit float64) {
	for _, row := range tableDataLeft {
		var debit, credit float64
		fmt.Sscanf(row[2], "%f", &debit)
		fmt.Sscanf(row[3], "%f", &credit)
		totalDebit += debit
		totalCredit += credit
	}
	return
}

// Остальные функции (addCenteredBoldLargeParagraph, addLeftAlignedBoldText, addSingleRowWithTwoTexts и другие) остаются без изменений...

// Функция для добавления жирного текста большего размера по центру страницы
func addCenteredBoldLargeParagraph(doc *document.Document, text string) {
	para := doc.AddParagraph()
	para.SetAlignment(wml.ST_JcCenter)
	run := para.AddRun()
	run.AddText(text)
	run.Properties().SetBold(true)
	run.Properties().SetSize(20) // Установленный размер текста 20
}

// Функция для добавления текста по центру страницы
func addCenteredParagraph(doc *document.Document, text string) {
	para := doc.AddParagraph()
	para.SetAlignment(wml.ST_JcCenter)
	run := para.AddRun()
	run.AddText(text)
}

// Функция для добавления жирного текста по левому краю страницы
func addLeftAlignedBoldParagraph(doc *document.Document, text string) {
	para := doc.AddParagraph()
	para.SetAlignment(wml.ST_JcLeft)
	run := para.AddRun()
	run.AddText(text)
	run.Properties().SetBold(true)
}

// Функция для добавления текста по левому краю страницы
func addLeftAlignedParagraph(doc *document.Document, text string) {
	para := doc.AddParagraph()
	para.SetAlignment(wml.ST_JcLeft)
	run := para.AddRun()
	run.AddText(text)
}

// Функция для добавления текста по центру в ячейку
func addCenteredText(cell document.Cell, text string) {
	para := cell.AddParagraph()
	para.SetAlignment(wml.ST_JcCenter)
	run := para.AddRun()
	run.AddText(text)
}

// Функция для добавления жирного текста по левому краю в ячейку
func addLeftAlignedBoldText(cell document.Cell, text string) {
	para := cell.AddParagraph()
	para.SetAlignment(wml.ST_JcLeft)
	run := para.AddRun()
	run.AddText(text)
	run.Properties().SetBold(true)
}

// Функция для добавления текста по левому краю в ячейку
func addLeftAlignedText(cell document.Cell, text string) {
	para := cell.AddParagraph()
	para.SetAlignment(wml.ST_JcLeft)
	run := para.AddRun()
	run.AddText(text)
}

// Функция для добавления строки с двумя текстами в одной строке
func addSingleRowWithTwoTexts(doc *document.Document, leftText, rightText string) {
	table := doc.AddTable()
	table.Properties().SetWidthPercent(100)

	row := table.AddRow()

	// Первая ячейка с левым текстом
	leftCell := row.AddCell()
	addLeftAlignedText(leftCell, leftText)

	// Вторая ячейка с правым текстом
	rightCell := row.AddCell()
	addRightAlignedText(rightCell, rightText)
}

// Функция для добавления текста по правому краю в ячейку
func addRightAlignedText(cell document.Cell, text string) {
	para := cell.AddParagraph()
	para.SetAlignment(wml.ST_JcRight)
	run := para.AddRun()
	run.AddText(text)
}

// Функция для добавления строки с подписями, подчеркиванием и изображением
func addSignaturesRow(doc *document.Document) {
	table := doc.AddTable()
	table.Properties().SetWidthPercent(100)

	// Строка с подписями и подчеркиванием
	row := table.AddRow()

	// Левое подчеркивание с изображением
	leftCell := row.AddCell()
	leftCellPara := leftCell.AddParagraph()
	leftCellPara.SetAlignment(wml.ST_JcLeft)

	// Добавление изображения подписи в левую ячейку
	imgFile := "/root/bitrixChatgpt/sign.jpg"
	img, err := common.ImageFromFile(imgFile)
	if err != nil {
		log.Printf("Ошибка при загрузке изображения: %v", err)
		return
	}

	imgRef, err := doc.AddImage(img)
	if err != nil {
		log.Printf("Ошибка при добавлении изображения в документ: %v", err)
		return
	}

	anchored, err := leftCellPara.AddRun().AddDrawingAnchored(imgRef)
	if err != nil {
		log.Printf("Ошибка при добавлении изображения: %v", err)
		return
	}

	anchored.SetSize(2*measurement.Inch, 1*measurement.Inch) // Установите размер изображения по необходимости
	anchored.SetHAlignment(wml.WdST_AlignHCenter)
	anchored.SetYOffset(0.1 * measurement.Inch)
	anchored.SetTextWrapSquare(wml.WdST_WrapTextBothSides)

	// Добавление текста после изображения
	leftRun := leftCellPara.AddRun()
	leftRun.AddText("_______________ подпись")

	// Правое подчеркивание
	rightCell := row.AddCell()
	addSignatureUnderline(rightCell, wml.ST_JcRight)

	// Строка с текстом "подпись"
	row = table.AddRow()

	// Левый текст "подпись" (пусто, так как уже добавлен вместе с изображением)
	leftCell = row.AddCell()
	addSignatureText(leftCell, "", wml.ST_JcLeft)

	// Правый текст "подпись"
	rightCell = row.AddCell()
	addSignatureText(rightCell, "подпись", wml.ST_JcRight)

	// Следующая строка с подчеркиванием
	row = table.AddRow()

	// Левое подчеркивание
	leftCell = row.AddCell()
	addSignatureUnderline(leftCell, wml.ST_JcLeft)

	// Правое подчеркивание
	rightCell = row.AddCell()
	addSignatureUnderline(rightCell, wml.ST_JcRight)

	// Строка с текстом "расшифровка подписи"
	row = table.AddRow()

	// Левый текст "расшифровка подписи"
	leftCell = row.AddCell()
	addSignatureText(leftCell, "расшифровка подписи", wml.ST_JcLeft)

	// Правый текст "расшифровка подписи"
	rightCell = row.AddCell()
	addSignatureText(rightCell, "расшифровка подписи", wml.ST_JcRight)
}

// Функция для добавления подчеркивания
func addSignatureUnderline(cell document.Cell, alignment wml.ST_Jc) {
	para := cell.AddParagraph()
	para.SetAlignment(alignment)
	run := para.AddRun()
	run.AddText("_______________")
}

// Функция для добавления текста "подпись" или "расшифровка подписи"
func addSignatureText(cell document.Cell, text string, alignment wml.ST_Jc) {
	para := cell.AddParagraph()
	para.SetAlignment(alignment)
	run := para.AddRun()
	run.AddText(text)
	run.Properties().SetSize(10)
}
