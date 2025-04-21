// Navbar.test.jsx
import { render, screen, fireEvent } from '@testing-library/react';
import Navbar from './Navbar';

test('renders navbar with logo and nav links', () => {
  render(<Navbar />);
  expect(screen.getByText(/GatorNest - INN/i)).toBeInTheDocument();
  expect(screen.getByText(/Home/i)).toBeInTheDocument();
});

test('toggles navbar visibility on icon click', () => {
  render(<Navbar />);
  const toggleButton = screen.getByRole('button', { hidden: true });
  fireEvent.click(toggleButton);
});
