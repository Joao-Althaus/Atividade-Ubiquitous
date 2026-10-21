package edu.users.model;

import java.time.LocalDate;
import java.time.OffsetDateTime;

public record User(
        long id,
        String name,
        String email,
        LocalDate birthday,
        String cpf,
        OffsetDateTime inserted_at,
        OffsetDateTime updated_at) {
}
