package main

import (
	"fmt"
	"time"

	"github.com/jpcshka/ip-address-with-script/script/fileio"
	"github.com/jpcshka/ip-address-with-script/script/routes"
)

func main() {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		fmt.Printf("\nВремя выполнения: %s\n", duration)
		fmt.Println("Нажмите Enter для выхода...")
		fmt.Scanln()
	}()

	finalFile := "all.bat"
	batFiles := fileio.SearchBatFiles(finalFile)
	if len(batFiles) == 0 {
		fmt.Println("Не найдено .bat файлов для обработки.")
		return
	}

	var routers []routes.Route
	for _, batFile := range batFiles {
		fileLines, err := fileio.ReadFile(batFile)
		if err != nil {
			fmt.Println(err)
		}
		for _, line := range fileLines {
			prefix, err := routes.ParseRouteFromLine(line, batFile)
			if err != nil {
				fmt.Println(err)
			}
			routers = append(routers, prefix)
		}
	}
	routeList := routes.NewRouteList(routers)
	if routeList == nil {
		fmt.Println("Не найдено маршрутов для обработки.")
		return
	}

	routeList.Sort()
	fmt.Println("Маршруты успешно отсортированы и отфильтрованы.")
	uniqRouteList, err := routeList.UniqueRoutes()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Всего маршрутов:", uniqRouteList.Count)
	fmt.Println("Удалено маршрутов:", uniqRouteList.Deleted)

	var outputLines []string
	for _, route := range uniqRouteList.Routes {
		line := routes.LineToBat(route)
		outputLines = append(outputLines, line)
	}
	err = fileio.WriteLinesToFile("all.bat", outputLines)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Маршруты успешно записаны в файл %s", finalFile)
}
