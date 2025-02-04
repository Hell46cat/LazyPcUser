package main

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	// Настройка маршрутов
	r := mux.NewRouter()
	r.HandleFunc("/volume", handleVolume).Methods("GET")
	r.HandleFunc("/whatvolume", handleWhatVolume).Methods("GET")

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:8088")
	http.ListenAndServe(":8088", r)
}

// Обработчик для установки громкости
func handleVolume(w http.ResponseWriter, r *http.Request) {
	volumeStr := r.URL.Query().Get("level")
	volume, err := strconv.Atoi(volumeStr)
	if err != nil || volume < 0 || volume > 100 {
		http.Error(w, "Неверный уровень громкости. Укажите число от 0 до 100.", http.StatusBadRequest)
		return
	}
	// Установка громкости
	cmd := exec.Command("amixer", "set", "Master", strconv.Itoa(volume)+"%")
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
	cmd := exec.Command("amixer", "get", "Master")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	// Используем регулярное выражение для поиска уровня громкости
	re := regexp.MustCompile(`\[(\d+)%\]`)
	match := re.FindStringSubmatch(out.String())
	if len(match) == 2 {
		volume, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}
		return volume, nil
	}

	return 0, nil
}