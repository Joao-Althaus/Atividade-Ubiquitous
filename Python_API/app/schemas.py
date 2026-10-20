from datetime import date, datetime
from typing import Annotated

from pydantic import BaseModel, ConfigDict, EmailStr, Field, StringConstraints

CPF = Annotated[
    str,
    StringConstraints(pattern=r"^[0-9]{11}$", min_length=11, max_length=11),
]


class UserInput(BaseModel):
    model_config = ConfigDict(str_strip_whitespace=True)

    name: str = Field(min_length=1, max_length=120)
    email: EmailStr
    birthday: date
    cpf: CPF


class User(UserInput):
    id: int
    inserted_at: datetime
    updated_at: datetime


class Error(BaseModel):
    error: str
