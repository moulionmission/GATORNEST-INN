import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import Home from './Home.jsx';

describe('Home Component', () => {
  test('renders the main titles and package text', () => {
    render(<Home />);
    expect(screen.getByText(/Our Packages/i)).toBeInTheDocument();
    expect(screen.getByText(/Search your Hotels/i)).toBeInTheDocument();
  });

  test('renders the location input with placeholder text', () => {
    render(<Home />);
    const locationInput = screen.getByPlaceholderText(/Enter location here.../i);
    expect(locationInput).toBeInTheDocument();
  });

  test('renders the date input element', () => {
    render(<Home />);
    const dateLabel = screen.getByText(/Select your date:/i);
    expect(dateLabel).toBeInTheDocument();
    const { container } = render(<Home />);
    const dateInput = container.querySelector('input[type="date"]');
    expect(dateInput).toBeInTheDocument();
  });

  test('renders the price input range and displays total price', () => {
    const { container } = render(<Home />);
    expect(screen.getByText(/\$1000/i)).toBeInTheDocument();
    const rangeInput = container.querySelector('input[type="range"]');
    expect(rangeInput).toBeInTheDocument();
    expect(rangeInput).toHaveAttribute('max', '1000');
    expect(rangeInput).toHaveAttribute('min', '100');
  });

  test('renders the MOREFILTERS button', () => {
    render(<Home />);
    expect(screen.getByText(/MOREFILTERS/i)).toBeInTheDocument();
  });

  test('renders the video element with correct attributes', () => {
    const { container } = render(<Home />);
    const videoElement = container.querySelector('video');
    expect(videoElement).toBeInTheDocument();
    expect(videoElement.muted).toBe(true);
    expect(videoElement.autoplay).toBe(true);
    expect(videoElement.loop).toBe(true);
    });
});
