package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ApiGolangv2/internal/entity"
	"ApiGolangv2/internal/repository"

	"github.com/redis/go-redis/v9"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	cacheAllTasks   = "tasks:all"
	cacheTaskByID   = "tasks:id:%d"
	cacheTaskStatus = "tasks:status:%s"
	cacheTTL        = 60 * time.Second
)

type TaskService interface {
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
	GetTaskByID(ctx context.Context, id int) (*entity.Task, error)
	GetTaskByStatus(ctx context.Context, status string) ([]entity.Task, error)
	CreateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

type taskService struct {
	repo  repository.TaskRepository
	redis *redis.Client
}

func NewTaskService(repo repository.TaskRepository, redis *redis.Client) TaskService {
	return &taskService{repo: repo, redis: redis}
}

func (s *taskService) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "service.GetAllTasks",
		tracer.ServiceName("task-manager"),
		tracer.ResourceName("GetAllTasks"),
	)
	defer span.Finish()

	// cek cache
	cacheSpan, cacheCtx := tracer.StartSpanFromContext(ctx, "redis.Get",
		tracer.ServiceName("task-manager-redis"),
		tracer.ResourceName(cacheAllTasks),
	)
	cached, err := s.redis.Get(cacheCtx, cacheAllTasks).Result()
	cacheSpan.Finish()

	if err == nil {
		var tasks []entity.Task
		if err := json.Unmarshal([]byte(cached), &tasks); err == nil {
			span.SetTag("cache", "hit")
			return tasks, nil
		}
	}
	span.SetTag("cache", "miss")

	// ambil dari DB
	tasks, err := s.repo.FindAll(ctx)
	if err != nil {
		span.Finish(tracer.WithError(err))
		return nil, err
	}

	// simpan ke cache
	setSpan, setCtx := tracer.StartSpanFromContext(ctx, "redis.SetEx",
		tracer.ServiceName("task-manager-redis"),
		tracer.ResourceName(cacheAllTasks),
	)
	data, _ := json.Marshal(tasks)
	s.redis.SetEx(setCtx, cacheAllTasks, data, cacheTTL)
	setSpan.Finish()

	return tasks, nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id int) (*entity.Task, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "service.GetTaskByID",
		tracer.ServiceName("task-manager"),
		tracer.ResourceName("GetTaskByID"),
		tracer.Tag("task.id", id),
	)
	defer span.Finish()

	key := fmt.Sprintf(cacheTaskByID, id)

	cacheSpan, cacheCtx := tracer.StartSpanFromContext(ctx, "redis.Get",
		tracer.ServiceName("task-manager-redis"),
		tracer.ResourceName(key),
	)
	cached, err := s.redis.Get(cacheCtx, key).Result()
	cacheSpan.Finish()

	if err == nil {
		var task entity.Task
		if err := json.Unmarshal([]byte(cached), &task); err == nil {
			span.SetTag("cache", "hit")
			return &task, nil
		}
	}
	span.SetTag("cache", "miss")

	task, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.Finish(tracer.WithError(err))
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	setSpan, setCtx := tracer.StartSpanFromContext(ctx, "redis.SetEx",
		tracer.ServiceName("task-manager-redis"),
		tracer.ResourceName(key),
	)
	data, _ := json.Marshal(task)
	s.redis.SetEx(setCtx, key, data, cacheTTL)
	setSpan.Finish()

	return task, nil
}

func (s *taskService) GetTaskByStatus(ctx context.Context, status string) ([]entity.Task, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "service.GetTaskByStatus",
		tracer.ServiceName("task-manager"),
		tracer.ResourceName("GetTaskByStatus"),
		tracer.Tag("task.status", status),
	)
	defer span.Finish()

	key := fmt.Sprintf(cacheTaskStatus, status)

	cacheSpan, cacheCtx := tracer.StartSpanFromContext(ctx, "redis.Get",
		tracer.ServiceName("task-manager-redis"),
		tracer.ResourceName(key),
	)
	cached, err := s.redis.Get(cacheCtx, key).Result()
	cacheSpan.Finish()

	if err == nil {
		var tasks []entity.Task
		if err := json.Unmarshal([]byte(cached), &tasks); err == nil {
			span.SetTag("cache", "hit")
			return tasks, nil
		}
	}
	span.SetTag("cache", "miss")

	tasks, err := s.repo.FindByStatus(ctx, status)
	if err != nil {
		span.Finish(tracer.WithError(err))
		return nil, err
	}

	setSpan, setCtx := tracer.StartSpanFromContext(ctx, "redis.SetEx",
		tracer.ServiceName("task-manager-redis"),
		tracer.ResourceName(key),
	)
	data, _ := json.Marshal(tasks)
	s.redis.SetEx(setCtx, key, data, cacheTTL)
	setSpan.Finish()

	return tasks, nil
}

func (s *taskService) CreateTask(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "service.CreateTask",
		tracer.ServiceName("task-manager"),
		tracer.ResourceName("CreateTask"),
		tracer.Tag("task.title", task.Title),
		tracer.Tag("task.priority", task.Priority),
	)
	defer span.Finish()

	created, err := s.repo.Create(ctx, task)
	if err != nil {
		span.Finish(tracer.WithError(err))
		return nil, err
	}

	s.invalidateCache(ctx, 0)

	span.SetTag("task.id", created.IDTask)
	return created, nil
}

func (s *taskService) UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "service.UpdateTask",
		tracer.ServiceName("task-manager"),
		tracer.ResourceName("UpdateTask"),
		tracer.Tag("task.id", task.IDTask),
		tracer.Tag("task.status", task.Status),
	)
	defer span.Finish()

	updated, err := s.repo.Update(ctx, task)
	if err != nil {
		span.Finish(tracer.WithError(err))
		return nil, err
	}
	if updated == nil {
		return nil, nil
	}

	s.invalidateCache(ctx, task.IDTask)

	return updated, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id int) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "service.DeleteTask",
		tracer.ServiceName("task-manager"),
		tracer.ResourceName("DeleteTask"),
		tracer.Tag("task.id", id),
	)
	defer span.Finish()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		span.Finish(tracer.WithError(err))
		return err
	}

	s.invalidateCache(ctx, id)

	return nil
}

func (s *taskService) invalidateCache(ctx context.Context, id int) {
	span, ctx := tracer.StartSpanFromContext(ctx, "redis.invalidateCache",
		tracer.ServiceName("task-manager-redis"),
		tracer.ResourceName("invalidateCache"),
	)
	defer span.Finish()

	s.redis.Del(ctx, cacheAllTasks)
	if id != 0 {
		s.redis.Del(ctx, fmt.Sprintf(cacheTaskByID, id))
	}
	for _, status := range []string{"pending", "in_progress", "done", "cancelled"} {
		s.redis.Del(ctx, fmt.Sprintf(cacheTaskStatus, status))
	}
}
