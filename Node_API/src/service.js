const UserRepository = require('./repository');

class UserService {
  constructor() {
    this.repository = new UserRepository();
  }

  list() {
    return this.repository.list();
  }

  async get(id) {
    const user = await this.repository.get(id);
    if (!user) {
      const error = new Error('user not found');
      error.status = 404;
      throw error;
    }
    return user;
  }

  create(user) {
    return this.repository.create(user);
  }

  async update(id, user) {
    const updatedUser = await this.repository.update(id, user);
    if (!updatedUser) {
      const error = new Error('user not found');
      error.status = 404;
      throw error;
    }
    return updatedUser;
  }

  async delete(id) {
    const deleted = await this.repository.delete(id);
    if (!deleted) {
      const error = new Error('user not found');
      error.status = 404;
      throw error;
    }
  }
}

module.exports = UserService;
