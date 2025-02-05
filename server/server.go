package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/mux"
)

func main() {
	// Создаем графическое приложение
	a := app.New()
	w := a.NewWindow("Сервер управления громкостью")
	w.SetContent(createUI(w))
	w.Resize(fyne.NewSize(300, 200)) // Устанавливаем размер окна
	w.ShowAndRun()
}

// Функция для создания интерфейса
func createUI(window fyne.Window) fyne.CanvasObject {
	// Текст "Сервер работает"
	statusLabel := widget.NewLabel("Сервер работает")

	// Кнопка для выключения сервера
	stopButton := widget.NewButton("Выключить сервер", func() {
		log.Println("Выключение сервера...")
		window.Close() // Закрываем окно приложения
	})

	// Запуск HTTP сервера в отдельной горутине
	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/volume", handleVolume).Methods("GET") // Обработчик для установки громкости
		r.HandleFunc("/whatvolume", handleWhatVolume).Methods("GET") // Обработчик для получения громкости

		log.Println("Сервер запущен на http://localhost:8088")
		if err := http.ListenAndServe(":8088", r); err != nil && err != http.ErrServerClosed {
			log.Fatal("Ошибка сервера:", err)
		}
	}()

	// Возвращаем интерфейс с кнопкой и текстом
	return container.NewVBox(
		statusLabel,
		stopButton,
	)
}

// Обработчик для установки громкости
func handleVolume(w http.ResponseWriter, r *http.Request) {
	// Получаем уровень громкости из URL параметра "level"
	volumeStr := r.URL.Query().Get("level")
	volume, err := strconv.Atoi(volumeStr)
	if err != nil || volume < 0 || volume > 100 {
		http.Error(w, "Неверный уровень громкости. Укажите число от 0 до 100.", http.StatusBadRequest)
		return
	}

	// Формируем команду в зависимости от операционной системы
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("amixer", "set", "Master", strconv.Itoa(volume)+"%")
	case "darwin":
		cmd = exec.Command("osascript", "-e", "set volume output volume "+strconv.Itoa(volume))
	default:
		http.Error(w, "Операционная система не поддерживается.", http.StatusBadRequest)
		return
	}

	// Выполняем команду для установки громкости
	if err := cmd.Run(); err != nil {
		http.Error(w, "Не удалось установить громкость.", http.StatusBadRequest)
		return
	}

	// Отправляем ответ
	w.Write([]byte("Громкость установлена на " + volumeStr + "%"))
}

// Обработчик для получения текущей громкости
func handleWhatVolume(w http.ResponseWriter, r *http.Request) {
	// Получаем текущий уровень громкости
	volume, err := findVolume()
	if err != nil {
		http.Error(w, "Не удалось получить уровень громкости.", http.StatusBadRequest)
		return
	}

	// Отправляем текущую громкость
	w.Write([]byte("Текущая громкость: " + strconv.Itoa(volume) + "%"))
}

// Функция для получения текущей громкости
func findVolume() (int, error) {
	// Формируем команду для получения громкости в зависимости от операционной системы
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("amixer", "get", "Master")
	case "darwin":
		cmd = exec.Command("osascript", "-e", "output volume of (get volume settings)")
	default:
		return 0, fmt.Errorf("операционная система не поддерживается")
	}

	// Выполняем команду и считываем результат
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0, err
	}

	// Парсим результат с помощью регулярных выражений
	var re *regexp.Regexp
	switch runtime.GOOS {
	case "linux":
		re = regexp.MustCompile(`\[(\d+)%\]`) // Регулярное выражение для Linux
	case "darwin":
		re = regexp.MustCompile(`(\d+)`) // Регулярное выражение для macOS
	}

	// Ищем совпадения
	match := re.FindStringSubmatch(out.String())
	if len(match) == 2 {
		volume, err := strconv.Atoi(match[1]) // Преобразуем строку в число
		if err != nil {
			return 0, err
		}
		return volume, nil
	}

	return 0, fmt.Errorf("не удалось найти уровень громкости")
}
