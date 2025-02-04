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
	w.Resize(fyne.NewSize(300, 200)) // Используем fyne.Size
	w.ShowAndRun()
}

// Функция для создания интерфейса
func createUI(window fyne.Window) fyne.CanvasObject {
	// Текст "Сервер работает"
	statusLabel := widget.NewLabel("Сервер работает")

	// Кнопка для выключения сервера
	stopButton := widget.NewButton("Выключить сервер", func() {
		log.Println("Выключение сервера...")
		window.Close() // Закрываем окно
	})

	// Запуск сервера в отдельной горутине
	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/volume", handleVolume).Methods("GET")
		r.HandleFunc("/whatvolume", handleWhatVolume).Methods("GET")

		log.Println("Сервер запущен на http://localhost:8088")
		if err := http.ListenAndServe(":8088", r); err != nil && err != http.ErrServerClosed {
			log.Fatal("Ошибка сервера:", err)
		}
	}()

	// Возвращаем интерфейс
	return container.NewVBox(
		statusLabel,
		stopButton,
	)
}

// Обработчик для установки громкости
func handleVolume(w http.ResponseWriter, r *http.Request) {
	volumeStr := r.URL.Query().Get("level")
	volume, err := strconv.Atoi(volumeStr)
	if err != nil || volume < 0 || volume > 100 {
		http.Error(w, "Неверный уровень громкости. Укажите число от 0 до 100.", http.StatusBadRequest)
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("amixer", "set", "Master", strconv.Itoa(volume)+"%")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("osascript", "-e", "set volume output volume "+strconv.Itoa(volume))
	} else {
		http.Error(w, "Операционная система не поддерживается.", http.StatusBadRequest)
		return
	}

	err = cmd.Run()
	if err != nil {
		http.Error(w, "Не удалось установить громкость.", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Громкость установлена на " + volumeStr + "%"))
}

// Обработчик для получения текущей громкости
func handleWhatVolume(w http.ResponseWriter, r *http.Request) {
	volume, err := findVolume()
	if err != nil {
		http.Error(w, "Не удалось получить уровень громкости.", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Текущая громкость: " + strconv.Itoa(volume) + "%"))
}

// Функция для получения текущей громкости
func findVolume() (int, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("amixer", "get", "Master")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("osascript", "-e", "output volume of (get volume settings)")
	} else {
		return 0, fmt.Errorf("операционная система не поддерживается")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	var re *regexp.Regexp
	if runtime.GOOS == "linux" {
		re = regexp.MustCompile(`\[(\d+)%\]`)
	} else if runtime.GOOS == "darwin" {
		re = regexp.MustCompile(`(\d+)`)
	} else {
		return 0, fmt.Errorf("операционная система не поддерживается")
	}

	match := re.FindStringSubmatch(out.String())
	if len(match) == 2 {
		volume, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}
		return volume, nil
	}

	return 0, fmt.Errorf("не удалось найти уровень громкости")
}
