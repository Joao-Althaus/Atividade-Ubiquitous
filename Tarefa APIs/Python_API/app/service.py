from .repository import UserRepository


class UserService:
    def __init__(self):
        self.repository = UserRepository()

    def list(self):
        return self.repository.list()

    def get(self, user_id: int):
        user = self.repository.get(user_id)
        if user is None:
            raise LookupError("user not found")
        return user

    def create(self, data):
        return self.repository.create(data)

    def update(self, user_id: int, data):
        user = self.repository.update(user_id, data)
        if user is None:
            raise LookupError("user not found")
        return user

    def delete(self, user_id: int):
        deleted = self.repository.delete(user_id)
        if not deleted:
            raise LookupError("user not found")
