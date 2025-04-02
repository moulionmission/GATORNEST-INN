// Footer.test.jsx
import { render, screen } from '@testing-library/react';
import Footer from './Footer';

test('renders video and contact section', () => {
  render(<Footer />);
  expect(screen.getByText(/KEEP IN TOUCH/i)).toBeInTheDocument();
  expect(screen.getByText(/Big things are coming/i)).toBeInTheDocument();
});

test('renders footer links and social icons', () => {
  render(<Footer />);
  expect(screen.getByText(/OUR AGENCY/i)).toBeInTheDocument();
  expect(screen.getByText(/PARTNERS/i)).toBeInTheDocument();
});