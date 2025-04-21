import { render, screen, fireEvent } from '@testing-library/react';
import Login from './Login';
import { MemoryRouter } from 'react-router-dom';

test('renders login form inputs and button', () => {
  render(
    <MemoryRouter>
      <Login onToggle={() => {}} />
    </MemoryRouter>
  );

  const emailInput = screen.getByPlaceholderText(/Email/i);
  const passwordInput = screen.getByPlaceholderText(/Password/i);
  const loginButton = screen.getByRole('button', { name: /Login/i });

  expect(emailInput).toBeInTheDocument();
  expect(passwordInput).toBeInTheDocument();
  expect(loginButton).toBeInTheDocument();
});

test('submits login form with valid input', () => {
  render(
    <MemoryRouter>
      <Login onToggle={() => {}} />
    </MemoryRouter>
  );

  const emailInput = screen.getByPlaceholderText(/Email/i);
  const passwordInput = screen.getByPlaceholderText(/Password/i);
  const loginButton = screen.getByRole('button', { name: /Login/i });

  fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
  fireEvent.change(passwordInput, { target: { value: 'password123' } });
  fireEvent.click(loginButton);
});
