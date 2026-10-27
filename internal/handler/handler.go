package handler

import (
	"net/http"
	"strconv"

	"delo-backend/internal/model"
	"delo-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.ProjectService
}

func NewHandler(s *service.ProjectService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	// CORS для работы с React фронтендом
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Проекты (по лекции 3: GET и POST)
	r.GET("/projects", h.ListProjects)
	r.POST("/projects", h.CreateProject)

	// Задачи внутри проектов
	r.GET("/tasks", h.ListTasks)
	r.POST("/tasks", h.CreateTask)
	r.POST("/tasks/status", h.UpdateTaskStatus)

	return r
}

func (h *Handler) ListProjects(c *gin.Context) {
	projects := h.service.GetProjects()
	c.JSON(http.StatusOK, projects)
}

func (h *Handler) CreateProject(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Owner       string `json:"owner"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные: " + err.Error()})
		return
	}

	project, err := h.service.CreateProject(req.Title, req.Description, req.Owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *Handler) ListTasks(c *gin.Context) {
	projectIDStr := c.Query("project_id")
	var projectID int
	if projectIDStr != "" {
		projectID, _ = strconv.Atoi(projectIDStr)
	}

	tasks := h.service.GetTasks(projectID)
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) CreateTask(c *gin.Context) {
	var req struct {
		ProjectID int    `json:"project_id" binding:"required"`
		Title     string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные: " + err.Error()})
		return
	}

	task, err := h.service.CreateTask(req.ProjectID, req.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *Handler) UpdateTaskStatus(c *gin.Context) {
	var req struct {
		TaskID int          `json:"task_id" binding:"required"`
		Status model.Status `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные: " + err.Error()})
		return
	}

	task, err := h.service.UpdateTaskStatus(req.TaskID, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}
