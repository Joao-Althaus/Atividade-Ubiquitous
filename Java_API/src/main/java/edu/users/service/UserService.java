package edu.users.service;

import edu.users.model.User;
import edu.users.model.UserInput;
import edu.users.repository.UserRepository;
import java.util.List;
import org.springframework.stereotype.Service;

@Service
public class UserService {
    private final UserRepository repository;

    public UserService(UserRepository repository) {
        this.repository = repository;
    }

    public List<User> list() {
        return repository.list();
    }

    public User get(long id) {
        return repository.get(id).orElseThrow(UserNotFoundException::new);
    }

    public User create(UserInput input) {
        return repository.create(input);
    }

    public User update(long id, UserInput input) {
        return repository.update(id, input).orElseThrow(UserNotFoundException::new);
    }

    public void delete(long id) {
        boolean deleted = repository.delete(id);
        if (!deleted) {
            throw new UserNotFoundException();
        }
    }

    public static class UserNotFoundException extends RuntimeException {
        public UserNotFoundException() {
            super("user not found");
        }
    }
}
