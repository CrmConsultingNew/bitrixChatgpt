package reports

import (
	"fmt"
	"os"
	"path/filepath"
)

// _log записывает данные в файл лога
func _log(title string, data interface{}) {
	if DEBUG {
		dir := LOG_DIR   // Название директории для лога
		file := LOG_FILE // Имя файла лога
		if dir != "" {
			file = filepath.Join(dir, file) // Полный путь к файлу лога
		}

		// Текст лога
		content := fmt.Sprintf(
			"================ %s ================\n%v\n",
			title,
			data,
		)

		// Создание директории лога, если она отсутствует
		if dir != "" {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err = os.Mkdir(dir, 0777)
				if err != nil {
					fmt.Printf("Ошибка создания директории: %v\n", err)
					return
				}
				err = os.Chmod(dir, 0777) // Установка прав на директорию
				if err != nil {
					fmt.Printf("Ошибка установки прав на директорию: %v\n", err)
					return
				}
			}
		}

		// Запись лога
		f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Ошибка открытия файла: %v\n", err)
			return
		}
		defer f.Close()

		if _, err := f.WriteString(content); err != nil {
			fmt.Printf("Ошибка записи в файл: %v\n", err)
		}
	}
}
