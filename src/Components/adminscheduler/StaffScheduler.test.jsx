import React from 'react';
import { render, screen } from '@testing-library/react';
import StaffScheduler from './adminscheduler';

test('renders the main heading', () => {
  render(<StaffScheduler />);
  expect(screen.getByText(/Hotel Admin Scheduler/i)).toBeTruthy();
});

test('renders Room Availability section', () => {
  render(<StaffScheduler />);
  expect(screen.getByText(/Room Availability/i)).toBeTruthy();
});



test('renders Weekly Schedule section', () => {
  render(<StaffScheduler />);
  expect(screen.getByText(/Weekly Schedule/i)).toBeTruthy();
});
