import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import Main from './Main.jsx';

describe('Main Component', () => {
  test('renders the section title "Popular Rooms"', () => {
    render(<Main />);
    const titleElement = screen.getByText(/Popular Rooms/i);
    expect(titleElement).toBeInTheDocument();
  });

  test('renders the correct number of room cards', () => {
    render(<Main />);
    const roomImages = screen.getAllByRole('img');
    expect(roomImages.length).toBe(9);
  });

  test('renders each room card with correct details', () => {
    render(<Main />);
    const roomTitle = screen.getByText(/Deluxe Suite/i);
    expect(roomTitle).toBeInTheDocument();

    const detailsButtons = screen.getAllByRole('button', { name: /DETAILS/i });
    expect(detailsButtons.length).toBe(9);
  });
});
