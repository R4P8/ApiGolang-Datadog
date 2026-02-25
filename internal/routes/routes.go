package routes

import (
	"net/http"

	"ApiGolangv2/internal/handler"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func NewRouter(taskHandler *handler.TaskHandler) http.Handler {
	mux := httptrace.NewServeMux(
		httptrace.WithServiceName("task-manager"),
	)

	mux.HandleFunc("GET    /tasks", taskHandler.GetAllTasks)
	mux.HandleFunc("POST   /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET    /tasks/{id}", taskHandler.GetTaskByID)
	mux.HandleFunc("PUT    /tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	return mux
}
