package main

import (
	"encoding/json"
	"fmt"
	"functions/functions"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func templateHTML(tablefiles []functions.File) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles("ui/static/index.html")
		if err != nil {
			http.Error(w, "Ошибка анализа шаблона:"+err.Error(), http.StatusInternalServerError)
			return
		}
		if err = templ.Execute(w, tablefiles); err != nil {
			http.Error(w, "Ошибка выполнения шаблона:"+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// функция которая принимает в качестве аргументов средство записи HTTP-ответа и HTTP-запрос.
func fletchHandler(table []functions.File) {
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(table); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// Настройка Сервера
func listenAndServer(addr string) {
	log.Println("Сервер работает на порту 8080...")
	log.Fatal(http.ListenAndServe(addr, nil))
}

// определение структуры файла
func main() {
	rootflag := "root"
	sortflag := "sort"
	root, _, err := functions.GetFilePathFromCommand(rootflag, sortflag)
	if err != nil {
		fmt.Println(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/ui/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("../ui/static"))))

	path := functions.Root{Name: root}
	Dirpath := path
	table, err := Dirpath.GetSubDir()
	if err != nil {
		fmt.Println(err)
	}

	fletchHandler(table)
	templateHTML(table)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		fmt.Println("Закрытие сервера...")
		os.Exit(0)
	}()
	listenAndServer(":8080")
}
