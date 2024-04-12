package main

import (
	"fmt"
	"functions/functions"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// определение структуры файла
func main() {
	rootflag := "root"
	sortflag := "sort"
	root, sort, err := functions.GetFilePathFromCommand(rootflag, sortflag)
	if err != nil {
		fmt.Println(err)
	}
	arrayfiles, err := functions.GetInfo(root)

	if err != nil {
		fmt.Println(err)
		return
	}

	fileStructArray := []functions.File{}
	var wg sync.WaitGroup
	for _, item := range arrayfiles {
		wg.Add(1)
		go func(item os.FileInfo) {
			defer wg.Done()
			switch mode := item.Mode(); {
			case mode.IsDir():
				var size int64 = 0
				err := filepath.Walk(item.Name(), func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					size += info.Size()
					return nil
				})
				fileStructArray = append(fileStructArray, functions.Newfile(item.Name(), "Каталог", size))
				if err != nil {
					fmt.Println(err)
				}
			case mode.IsRegular():
				fileStructArray = append(fileStructArray, functions.Newfile(item.Name(), "Файл", item.Size()))
			}

		}(item)
	}
	wg.Wait()
	functions.SelectSort(fileStructArray, root, sort)
	//http.Handle("/", http.FileServer(http.Dir("./static")))
	//table := functions.GetFilesFromDirectory(root)
	// Create a new ServeMux
	mux := http.NewServeMux()
	//Serve static files from the "static" directory
	//fileServer := http.FileServer(http.Dir("../ui"))
	mux.Handle("/ui/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("../ui/static"))))
	//staticDir := "./static"
	// Servir archivos estáticos desde el directorio especificado
	//fs := http.FileServer(http.Dir(staticDir))

	// Manejador para servir archivos estáticos
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	path := functions.Root{Name: root}
	subDir := path
	table, err := subDir.GetDir()
	if err != nil {
		fmt.Println(err)
	}

	functions.FletchHandler(table)
	functions.TemplateHTML(table)
	functions.ListenAndServer(":8080")

}
