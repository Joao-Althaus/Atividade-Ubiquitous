package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	cpfPattern   = regexp.MustCompile(`^[0-9]{11}$`)
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (controller *Controller) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", controller.handleHealth)
	mux.HandleFunc("/users", controller.handleUsers)
	mux.HandleFunc("/users/", controller.handleUser)
}

func (controller *Controller) handleHealth(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.Header().Set("Allow", http.MethodGet)
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(response, http.StatusOK, map[string]string{"status": "ok"})
}

func (controller *Controller) handleUsers(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		users, err := controller.service.List(request.Context())
		if err != nil {
			writeServiceError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, users)
	case http.MethodPost:
		input, ok := decodeInput(response, request)
		if !ok {
			return
		}
		user, err := controller.service.Create(request.Context(), input)
		if err != nil {
			writeServiceError(response, err)
			return
		}
		writeJSON(response, http.StatusCreated, user)
	default:
		response.Header().Set("Allow", "GET, POST")
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (controller *Controller) handleUser(response http.ResponseWriter, request *http.Request) {
	id, ok := parseID(response, request.URL.Path)
	if !ok {
		return
	}

	switch request.Method {
	case http.MethodGet:
		user, err := controller.service.Get(request.Context(), id)
		if err != nil {
			writeServiceError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, user)
	case http.MethodPut:
		input, ok := decodeInput(response, request)
		if !ok {
			return
		}
		user, err := controller.service.Update(request.Context(), id, input)
		if err != nil {
			writeServiceError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, user)
	case http.MethodDelete:
		if err := controller.service.Delete(request.Context(), id); err != nil {
			writeServiceError(response, err)
			return
		}
		response.WriteHeader(http.StatusNoContent)
	default:
		response.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(response, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func parseID(response http.ResponseWriter, path string) (int64, bool) {
	idText := strings.TrimPrefix(path, "/users/")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id < 1 || strings.Contains(idText, "/") {
		writeError(response, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}

func decodeInput(response http.ResponseWriter, request *http.Request) (Input, bool) {
	var input Input
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid JSON body")
		return Input{}, false
	}

	if err := validateInput(input); err != nil {
		writeError(response, http.StatusBadRequest, err.Error())
		return Input{}, false
	}
	return input, true
}

func validateInput(input Input) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.CPF = strings.TrimSpace(input.CPF)

	if input.Name == "" || len(input.Name) > 120 {
		return ErrName
	}
	if len(input.Email) > 254 || !emailPattern.MatchString(input.Email) {
		return ErrEmail
	}
	if !cpfPattern.MatchString(input.CPF) {
		return ErrCPF
	}
	if _, err := parseBirthday(input.Birthday); err != nil {
		return ErrBirthday
	}
	return nil
}

func writeServiceError(response http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeError(response, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrName), errors.Is(err, ErrEmail), errors.Is(err, ErrCPF), errors.Is(err, ErrBirthday):
		writeError(response, http.StatusBadRequest, err.Error())
	case IsConflict(err):
		writeError(response, http.StatusConflict, "email or cpf already exists")
	case IsInvalidDate(err):
		writeError(response, http.StatusBadRequest, ErrBirthday.Error())
	default:
		log.Printf("request error: %v", err)
		writeError(response, http.StatusInternalServerError, "internal server error")
	}
}

func writeError(response http.ResponseWriter, status int, message string) {
	writeJSON(response, status, map[string]string{"error": message})
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	if err := json.NewEncoder(response).Encode(value); err != nil {
		log.Printf("could not encode response: %v", err)
	}
}
