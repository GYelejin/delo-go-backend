package main

import (
	"log"

	"delo-backend/internal/handler"
	"delo-backend/internal/service"
)

func main() {
	svc := service.NewProjectService()
	h := handler.NewHandler(svc)
	routes := h.InitRoutes()

	log.Println("Сервер ДЕЛО Backend успешно запущен на порту :8080")
	if err := routes.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
