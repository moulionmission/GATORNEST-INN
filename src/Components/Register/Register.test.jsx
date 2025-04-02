import { render, screen, fireEvent } from '@testing-library/react';
import Register from './Register';

test('renders register form inputs and button', () => {
  render(<Register onToggle={() => {}} />);
  expect(screen.getByPlaceholderText(/Email/i)).toBeInTheDocument();
  expect(screen.getByPlaceholderText(/Password/i)).toBeInTheDocument();
  expect(screen.getByPlaceholderText(/Confirm Password/i)).toBeInTheDocument();
});

test('shows alert on password mismatch', () => {
  window.alert = jest.fn();
  render(<Register onToggle={() => {}} />);

  const emailInput = screen.getByPlaceholderText(/Email/i);
  const passwordInput = screen.getByPlaceholderText('Password');
  const confirmInput = screen.getByPlaceholderText('Confirm Password');

  fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
  fireEvent.change(passwordInput, { target: { value: 'abc123' } });
  fireEvent.change(confirmInput, { target: { value: 'xyz456' } });

  fireEvent.click(screen.getByRole('button', { name: /Register/i }));
  expect(window.alert).toHaveBeenCalledWith('Passwords do not match');
});
