package edu.users.repository;

import edu.users.model.User;
import edu.users.model.UserInput;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.time.ZoneOffset;
import java.util.List;
import java.util.Optional;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Repository;

@Repository
public class UserRepository {
    private static final String COLUMNS =
            "id, name, email, birthday, cpf, inserted_at, updated_at";

    private final JdbcTemplate jdbc;

    private final RowMapper<User> mapper = (resultSet, rowNumber) -> map(resultSet);

    public UserRepository(JdbcTemplate jdbc) {
        this.jdbc = jdbc;
    }

    private User map(ResultSet resultSet) throws SQLException {
        return new User(
                resultSet.getLong("id"),
                resultSet.getString("name"),
                resultSet.getString("email"),
                resultSet.getDate("birthday").toLocalDate(),
                resultSet.getString("cpf").trim(),
                resultSet.getTimestamp("inserted_at")
                        .toInstant()
                        .atOffset(ZoneOffset.UTC),
                resultSet.getTimestamp("updated_at")
                        .toInstant()
                        .atOffset(ZoneOffset.UTC));
    }

    public List<User> list() {
        String query = "SELECT " + COLUMNS + " FROM users ORDER BY id";
        return jdbc.query(query, mapper);
    }

    public Optional<User> get(long id) {
        String query = "SELECT " + COLUMNS + " FROM users WHERE id = ?";
        return jdbc.query(query, mapper, id).stream().findFirst();
    }

    public User create(UserInput user) {
        String query = "INSERT INTO users (name, email, birthday, cpf) "
                + "VALUES (?, ?, ?, ?) RETURNING " + COLUMNS;
        return jdbc.queryForObject(
                query,
                mapper,
                user.name().trim(),
                user.email().trim(),
                user.birthday(),
                user.cpf());
    }

    public Optional<User> update(long id, UserInput user) {
        String query = "UPDATE users "
                + "SET name = ?, email = ?, birthday = ?, cpf = ?, "
                + "updated_at = CURRENT_TIMESTAMP WHERE id = ?";
        int updatedRows = jdbc.update(
                query,
                user.name().trim(),
                user.email().trim(),
                user.birthday(),
                user.cpf(),
                id);

        if (updatedRows == 0) {
            return Optional.empty();
        }
        return get(id);
    }

    public boolean delete(long id) {
        return jdbc.update("DELETE FROM users WHERE id = ?", id) > 0;
    }
}
