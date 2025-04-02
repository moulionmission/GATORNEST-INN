// Home.test.jsx
import { render, screen } from '@testing-library/react';
import Home from './Home';

test('renders video background and main content', () => {
  render(<Home />);
  expect(screen.getByText(/Our Packages/i)).toBeInTheDocument();
  expect(screen.getByText(/Search your Hotels/i)).toBeInTheDocument();
});

test('renders all input sections like name and date', () => {
  render(<Home />);
  expect(screen.getByLabelText(/Name/i)).toBeInTheDocument();
  expect(screen.getByLabelText(/Select your date/i)).toBeInTheDocument();
});