// Main.test.jsx
import { render, screen, fireEvent } from '@testing-library/react';
import Main from './Main';

test('renders list of room cards', () => {
  render(<Main />);
  expect(screen.getByText(/Popular Rooms/i)).toBeInTheDocument();
});

test('opens booking modal on BOOK click', () => {
  render(<Main />);
  const bookButtons = screen.getAllByRole('button');
  fireEvent.click(bookButtons[0]);
});
