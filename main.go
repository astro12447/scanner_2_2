package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// определение структуры файла
type file struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

var s []file

// определение функции для ввода информации классы Files в консоль
func (ob *file) print() {
	fmt.Println("Name:", ob.Name, "Type:", ob.Extension, "FileSize/byte", ob.Size)
}

// определение функции для получения строк через консоль
func getFilePathFromCommand(root string, sort string) (string, string, error) {
	if root == "None" || sort == "None" {
		fmt.Println("->Введите правильную командную строку:(--root=/pathfile  --sort=Desc) or --root=/pathfile")
	}
	var sourcepath *string
	var sortflag *string
	sourcepath = flag.String(root, "None", "")
	sortflag = flag.String(sort, "None", "")
	flag.Parse()
	return *sourcepath, *sortflag, nil
}

// функция для проверкаи попки
func rootExist(root string) (bool, error) {
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		fmt.Println("Root не существует...!")
	}
	return true, nil
}

// метод для получения значения size класса
func (ob *file) getSize() int64 {
	return ob.Size
}

// метод для получения значения name класса
func (ob *file) getName() string {
	return ob.Name
}

// метод для получения значения Extension класса
func (ob *file) getExtension() string {
	return ob.Extension
}

// метод для получение информации о файлах
func getFilesRecurvise(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return true, nil
}

// метод для получение информации католога файлы
func getFileLocation(root string, filename string) (string, error) {
	if root == "" {
		return "", errors.New("Root  пуст!")
	}
	return root + "/" + filename, nil
}

// Получение все файл из котолога
func getAllFromDir(path string) ([]file, error) {
	err := filepath.Walk(path, func(p string, inf os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		size, err := getsize(p)
		if err != nil {
			fmt.Println(err)
		}
		Ext, err := getFileExtension2(p)
		if err != nil {
			fmt.Println(err)
		}
		element := file{Name: p, Extension: Ext, Size: size}
		s = append(s, element)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return s, nil
}

// функция для получения значения  size
func getsize(filename string) (int64, error) {
	f, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}
	return f.Size(), nil
}

// функция для получения значения  Extension
func getFileExtension(root string, filename string) (string, error) {
	f, err := getFileLocation(root, filename)
	if err != nil {
		fmt.Println(err)
	}
	st, err := os.Stat(f)
	if err != nil {
		fmt.Println(err)
	}
	if st.IsDir() {
		return "Каталог", nil
	}
	return "файл", nil
}

// функция для получения значения  Extension2
func getFileExtension2(filename string) (string, error) {
	f, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}
	if f.IsDir() {
		return "Каталог", nil
	}
	return "файл", nil
}

// функция для получения значения  файлы из католога
func getFilesFromDirectory(pathName string) ([]file, error) {
	fi, err := os.Open(pathName)
	if err != nil {
		log.Fatal(err, fi.Name())
	}
	defer fi.Close()
	files, err := os.ReadDir(pathName)
	if err != nil {
		fmt.Print("Невозможно прочитать каталога!", err)
	}
	for _, item := range files {
		p, err := getFileLocation(pathName, item.Name())
		f, err := os.Stat(p)
		if err != nil {
			panic(err)
		}
		Ext, err := getFileExtension(pathName, item.Name())
		name := pathName + "/" + f.Name()
		element := file{Name: name, Extension: Ext, Size: f.Size()}
		s = append(s, element)
		fmt.Println(Ext, name, f.Size())
	}
	return s, nil
}

// функция для Обработки сортировки по Убывающий
func sortAsc(arr []file) {
	if len(arr) < 0 {
		fmt.Println("Массив пуст!")
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Size < arr[j].Size
	})
}

// функция для Обработки сортировки по возврастающий
func sortDesc(arr []file) {
	if len(arr) < 0 {
		fmt.Println("Массив пуст!")
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Size > arr[j].Size
	})
}

// Чтение файлы из католога(Root)
func readDir(root string) ([]fs.FileInfo, error) {
	arrayfiles, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}
	return arrayfiles, nil
}

// выборка сортировки
func selectSort(struc []file, root string, sortMode string) error {
	if len(struc) < 0 {
		log.Panic("Нет элементов в массиве!")
	}
	if root != "None" && sortMode == "None" {
		sortAsc(struc)
		for i := 0; i < len(struc); i++ {
			struc[i].print()
		}
	} else if sortMode == "Desc" && root != "None" {
		sortDesc(struc)
		for i := 0; i < len(struc); i++ {
			struc[i].print()
		}
	}
	return nil
}
func main() {
	rootflag := "root"
	sortflag := "sort"
	root, sort, err := getFilePathFromCommand(rootflag, sortflag)
	if err != nil {
		fmt.Println(err)
	}
	arrayfiles, err := readDir(root)
	if err != nil {
		fmt.Println(err)
		return
	}
	filesArr := []file{}
	var wg sync.WaitGroup
	for _, item := range arrayfiles {
		wg.Add(1)
		go func(item os.FileInfo) {
			defer wg.Done()
			if err != nil {
				fmt.Println(err)
			}
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
				filesArr = append(filesArr, file{
					Name:      item.Name(),
					Extension: "Каталог",
					Size:      size,
				})
				if err != nil {
					log.Println(err)
				}
			case mode.IsRegular():
				filesArr = append(filesArr, file{
					Name:      item.Name(),
					Extension: "Файл",
					Size:      item.Size(),
				})

			}

		}(item)
	}
	wg.Wait()
	selectSort(filesArr, root, sort)
}
