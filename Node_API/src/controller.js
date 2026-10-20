const express = require('express');
const UserService = require('./service');

const router = express.Router();
const service = new UserService();
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function createHttpError(message, status) {
  const error = new Error(message);
  error.status = status;
  return error;
}

function validateName(name) {
  const validName =
    typeof name === 'string' &&
    name.trim().length > 0 &&
    name.trim().length <= 120;

  if (!validName) {
    throw createHttpError('name is required (max 120 characters)', 400);
  }
  return name.trim();
}

function validateEmail(email) {
  const validEmail =
    typeof email === 'string' &&
    email.trim().length <= 254 &&
    emailPattern.test(email.trim());

  if (!validEmail) {
    throw createHttpError('email is invalid', 400);
  }
  return email.trim();
}

function validateCpf(cpf) {
  const validCpf = typeof cpf === 'string' && /^\d{11}$/.test(cpf.trim());
  if (!validCpf) {
    throw createHttpError('cpf must contain exactly 11 digits', 400);
  }
  return cpf.trim();
}

function validateBirthday(birthday) {
  const hasValidFormat =
    typeof birthday === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(birthday);
  const parsedBirthday = hasValidFormat
    ? new Date(`${birthday}T00:00:00Z`)
    : null;
  const validDate =
    parsedBirthday !== null &&
    !Number.isNaN(parsedBirthday.getTime()) &&
    parsedBirthday.toISOString().slice(0, 10) === birthday;

  if (!validDate) {
    throw createHttpError('birthday must use YYYY-MM-DD', 400);
  }
  return birthday;
}

function validateUserInput(body) {
  const input = body || {};

  return {
    name: validateName(input.name),
    email: validateEmail(input.email),
    birthday: validateBirthday(input.birthday),
    cpf: validateCpf(input.cpf),
  };
}

function parseUserId(value) {
  const id = Number(value);
  const validId = Number.isSafeInteger(id) && id >= 1;

  if (!validId) {
    throw createHttpError('id must be a positive integer', 400);
  }
  return id;
}

function handleAsync(handler) {
  return function handleRequest(request, response, next) {
    Promise.resolve(handler(request, response)).catch(next);
  };
}

async function listUsers(_request, response) {
  const users = await service.list();
  response.status(200).json(users);
}

async function getUser(request, response) {
  const id = parseUserId(request.params.id);
  const user = await service.get(id);
  response.status(200).json(user);
}

async function createUser(request, response) {
  const input = validateUserInput(request.body);
  const user = await service.create(input);
  response.status(201).json(user);
}

async function updateUser(request, response) {
  const id = parseUserId(request.params.id);
  const input = validateUserInput(request.body);
  const user = await service.update(id, input);
  response.status(200).json(user);
}

async function deleteUser(request, response) {
  const id = parseUserId(request.params.id);
  await service.delete(id);
  response.status(204).end();
}

router.get('/users', handleAsync(listUsers));
router.get('/users/:id', handleAsync(getUser));
router.post('/users', handleAsync(createUser));
router.put('/users/:id', handleAsync(updateUser));
router.delete('/users/:id', handleAsync(deleteUser));

function errorHandler(error, _request, response, _next) {
  if (error.code === '23505') {
    return response
      .status(409)
      .json({ error: 'email or cpf already exists' });
  }

  const status = error.status || 500;
  const message = status === 500 ? 'internal server error' : error.message;
  return response.status(status).json({ error: message });
}

module.exports = {
  router,
  errorHandler,
};
