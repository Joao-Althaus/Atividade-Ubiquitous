package edu.users.controller;

import edu.users.model.User;
import edu.users.model.UserInput;
import edu.users.service.UserService;
import jakarta.validation.Valid;
import java.util.List;
import java.util.Map;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.http.converter.HttpMessageNotReadableException;

@RestController
public class UserController {
    private final UserService service;

    public UserController(UserService service) {
        this.service = service;
    }

    @GetMapping("/health")
    public Map<String, String> health() {
        return Map.of("status", "ok");
    }

    @GetMapping("/users")
    public List<User> list() {
        return service.list();
    }

    @GetMapping("/users/{id}")
    public User get(@PathVariable long id) {
        validateId(id);
        return service.get(id);
    }

    @PostMapping("/users")
    public ResponseEntity<User> create(@Valid @RequestBody UserInput input) {
        User user = service.create(input);
        return ResponseEntity.status(HttpStatus.CREATED).body(user);
    }

    @PutMapping("/users/{id}")
    public User update(@PathVariable long id, @Valid @RequestBody UserInput input) {
        validateId(id);
        return service.update(id, input);
    }

    @DeleteMapping("/users/{id}")
    public ResponseEntity<Void> delete(@PathVariable long id) {
        validateId(id);
        service.delete(id);
        return ResponseEntity.noContent().build();
    }

    private void validateId(long id) {
        if (id < 1) {
            throw new IllegalArgumentException("id deve ser um inteiro positivo");
        }
    }

    @ExceptionHandler(UserService.UserNotFoundException.class)
    public ResponseEntity<Map<String, String>> missing(Exception exception) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND)
                .body(Map.of("error", exception.getMessage()));
    }

    @ExceptionHandler(DuplicateKeyException.class)
    public ResponseEntity<Map<String, String>> duplicate() {
        return ResponseEntity.status(HttpStatus.CONFLICT)
                .body(Map.of("error", "email or cpf already exists"));
    }

    @ExceptionHandler({
        IllegalArgumentException.class,
        HttpMessageNotReadableException.class,
        MethodArgumentNotValidException.class
    })
    public ResponseEntity<Map<String, String>> badRequest(Exception exception) {
        return ResponseEntity.badRequest()
                .body(Map.of("error", "invalid request"));
    }
}
