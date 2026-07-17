package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	var paths []string

	for _, arg := range args {
		// Игнорируем типичные unix-флаги для rm
		if strings.HasPrefix(arg, "-") {
			lowerArg := strings.ToLower(arg)
			if lowerArg == "-r" || lowerArg == "-rf" || lowerArg == "-fr" || lowerArg == "-f" || lowerArg == "--recursive" || lowerArg == "--force" {
				continue
			}
			// Пропускаем все неизвестные флаги для простоты
			continue
		}
		paths = append(paths, arg)
	}

	if len(paths) == 0 {
		os.Exit(0)
	}

	for _, p := range paths {
		// Поддержка масок (globbing), например ./folder/*
		matches, err := filepath.Glob(p)
		if err != nil || len(matches) == 0 {
			// Если маска не совпала или возникла ошибка при парсинге маски, 
			// всё равно попытаемся удалить как точный путь 
			// (полезно, если файл содержит спецсимволы, но не является маской).
			_ = os.RemoveAll(p)
			continue
		}

		for _, match := range matches {
			// Используем os.RemoveAll, который удаляет путь со всем содержимым без вопросов (как -rf)
			err = os.RemoveAll(match)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка удаления %s: %v\n", match, err)
			}
		}
	}
}
