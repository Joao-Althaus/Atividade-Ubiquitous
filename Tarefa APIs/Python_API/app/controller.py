from fastapi import APIRouter, HTTPException, Response
from psycopg.errors import UniqueViolation

from .schemas import User, UserInput
from .service import UserService

router = APIRouter()
service = UserService()


def execute(action):
    try:
        return action()
    except LookupError as exc:
        raise HTTPException(status_code=404, detail=str(exc)) from exc
    except UniqueViolation as exc:
        raise HTTPException(
            status_code=409,
            detail="email or cpf already exists",
        ) from exc


@router.get("/users", response_model=list[User])
def list_users():
    return execute(service.list)


@router.get("/users/{user_id}", response_model=User)
def get_user(user_id: int):
    return execute(lambda: service.get(user_id))


@router.post("/users", response_model=User, status_code=201)
def create_user(data: UserInput):
    return execute(lambda: service.create(data.model_dump()))


@router.put("/users/{user_id}", response_model=User)
def update_user(user_id: int, data: UserInput):
    return execute(lambda: service.update(user_id, data.model_dump()))


@router.delete("/users/{user_id}", status_code=204)
def delete_user(user_id: int):
    execute(lambda: service.delete(user_id))
    return Response(status_code=204)
