package functions

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// определение структуры файла
type File struct {
	Typefile string `json:"typefile"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
}

// Метод для создания структурой файла
func Newfile(typefile string, name string, size int64) File {
	return File{
		Typefile: typefile,
		Name:     name,
		Size:     size,
	}
}

// Интерфейс(ReadPath) с методом для считания root
type ReadPath interface {
	GetsubDir(root string) ([]File, error)
}

// Папка структуры
type Root struct {
	Name string
}

// Реализация метода(GetsubDir) интерфейса
func (root *Root) GetSubDir() ([]File, error) {
	var files []File
	var size int64 = 0
	filepath.Walk(root.Name, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		switch mode := info.Mode(); {
		case mode.IsDir():
			size += info.Size()
			files = append(files, Newfile("Каталог", path, size))
		case mode.IsRegular():
			files = append(files, Newfile("Файл", path, info.Size()))
		}
		return nil
	})
	return files[1:len(files)], nil
}

// определение функции для ввода информации классы Files в консоль
func (ob *File) print() {
	fmt.Println("Type:", ob.Typefile, "Name:", ob.Name, "FileSize/byte", ob.Size)
}

// определение функции для получения строк через консоль
func GetFilePathFromCommand(root string, sort string) (string, string, error) {
	if root == "None" || sort == "None" {
		log.Fatal("->Введите правильную командную строку:(--root=/pathfile  --sort=Desc) or --root=/pathfile")
	} else if root == "None" && sort != "" {
		log.Fatal("->Введите правильную командную строку:(--root=/pathfile  --sort=Desc) or --root=/pathfile")
	}

	var sourcepath *string
	var sortflag *string
	sourcepath = flag.String(root, "None", "")
	sortflag = flag.String(sort, "None", "")
	flag.Parse()

	return *sourcepath, *sortflag, nil

}

// функция для проверкаи попки
func RootExist(root string) bool {
	if _, err := os.Stat(root); err != nil {
		log.Fatal("Root не существует...!")
	}
	return true
}

// метод для получения значения size класса
func (ob *File) getSize() int64 {
	return ob.Size
}

// метод для получения значения name класса
func (ob *File) getName() string {
	return ob.Name
}

// метод для получения значения Extension класса
func (ob *File) getExtension() string {
	return ob.Typefile
}

// функция для получения значения  size
func Getsize(filename string) (int64, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}
	return stat.Size(), nil
}

// функция для Обработки сортировки по Убывающий
func SortAsc(arr []File) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Size < arr[j].Size
	})
}

// функция для Обработки сортировки по возврастающий
func SortDesc(arr []File) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Size > arr[j].Size
	})
}

// Чтение файлы из католога(Root)
func GetInfo(dirname string) ([]fs.FileInfo, error) {
	dir, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	arrayInfo, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}
	return arrayInfo, nil
}

// выборка сортировки
func SelectSort(files []File, root string, sortMode string) {
	switch {
	case root != "None" && sortMode == "None":
		SortAsc(files)
		for i := range files {
			files[i].print()
		}
	case sortMode == "Desc" && root != "None":
		SortDesc(files)
		for i := range files {
			files[i].print()
		}
	default:
		log.Fatal("->Введите правильную командную строку:(--root=/pathfile  --sort=Desc) or --root=/pathfile")
	}

}
