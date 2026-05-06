package fileio

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// ReadFile Парсит строку с маршрутом и возвращает список строк.
func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// SearchBatFiles Ищет все .bat файлы в текущей директории и возвращает их список, исключая указанный файл.
// Если указанный файл найден, он удаляется из списка и возвращается обновленный список .bat файлов.
// Если указанный файл не найден, возвращается полный список .bat файлов.
func SearchBatFiles(finalFile string) []string {
	files, err := filepath.Glob("*.bat")
	if err != nil {
		fmt.Println(err)
	}
	if slices.Contains(files, finalFile) {
		index := slices.Index(files, finalFile)
		files = append(files[:index], files[index+1:]...)
	}
	return files
}

func WriteLinesToFile(finalFile string, lines []string) error {
	file, err := os.Create(finalFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return w.Flush()
}
