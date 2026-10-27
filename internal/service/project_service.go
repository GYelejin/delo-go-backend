package service

import (
	"errors"
	"sync"
	"time"

	"delo-backend/internal/model"
)

type ProjectService struct {
	mu       sync.RWMutex
	projects []model.Project
	tasks    []model.Task
	nextPID  int
	nextTID  int
}

func NewProjectService() *ProjectService {
	s := &ProjectService{
		nextPID: 3,
		nextTID: 4,
	}

	// Моковые данные по стандартам курса ДЕЛО
	s.projects = []model.Project{
		{ID: 1, Title: "Дело", Description: "Образовательная платформа для проектной деятельности", Owner: "Тигран", CreatedAt: time.Now()},
		{ID: 2, Title: "Витрина проектов", Description: "Биржа студенческих стартапов и команд", Owner: "Артем", CreatedAt: time.Now()},
	}

	s.tasks = []model.Task{
		{ID: 1, ProjectID: 1, Title: "Разработать каркас HTTP-сервера", Status: model.StatusDone, CreatedAt: time.Now()},
		{ID: 2, ProjectID: 1, Title: "Собрать Dockerfile", Status: model.StatusInProgress, CreatedAt: time.Now()},
		{ID: 3, ProjectID: 2, Title: "Интеграция с React-фронтендом", Status: model.StatusOpen, CreatedAt: time.Now()},
	}

	return s
}

func (s *ProjectService) GetProjects() []model.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]model.Project, len(s.projects))
	copy(res, s.projects)
	return res
}

func (s *ProjectService) CreateProject(title, description, owner string) (model.Project, error) {
	if title == "" {
		return model.Project{}, errors.New("название проекта не может быть пустым")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	p := model.Project{
		ID:          s.nextPID,
		Title:       title,
		Description: description,
		Owner:       owner,
		CreatedAt:   time.Now(),
	}
	s.nextPID++
	s.projects = append(s.projects, p)
	return p, nil
}

func (s *ProjectService) GetTasks(projectID int) []model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []model.Task
	for _, t := range s.tasks {
		if projectID == 0 || t.ProjectID == projectID {
			res = append(res, t)
		}
	}
	return res
}

func (s *ProjectService) CreateTask(projectID int, title string) (model.Task, error) {
	if title == "" {
		return model.Task{}, errors.New("заголовок задачи не может быть пустым")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := model.Task{
		ID:        s.nextTID,
		ProjectID: projectID,
		Title:     title,
		Status:    model.StatusOpen,
		CreatedAt: time.Now(),
	}
	s.nextTID++
	s.tasks = append(s.tasks, task)
	return task, nil
}

func (s *ProjectService) UpdateTaskStatus(taskID int, status model.Status) (model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == taskID {
			s.tasks[i].SetStatus(status)
			return s.tasks[i], nil
		}
	}
	return model.Task{}, errors.New("задача не найдена")
}
