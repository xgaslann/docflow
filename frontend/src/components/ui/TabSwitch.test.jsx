import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { TabSwitch } from '../TabSwitch';
import { FileText, FileOutput } from 'lucide-react';

const mockTabs = [
  { id: 'tab1', label: 'Tab One', icon: FileText },
  { id: 'tab2', label: 'Tab Two', icon: FileOutput },
];

describe('TabSwitch', () => {
  it('renders all tabs', () => {
    render(
      <TabSwitch
        tabs={mockTabs}
        activeTab="tab1"
        onChange={vi.fn()}
      />
    );

    expect(screen.getByText('Tab One')).toBeInTheDocument();
    expect(screen.getByText('Tab Two')).toBeInTheDocument();
  });

  it('calls onChange when tab is clicked', () => {
    const onChange = vi.fn();
    render(
      <TabSwitch
        tabs={mockTabs}
        activeTab="tab1"
        onChange={onChange}
      />
    );

    fireEvent.click(screen.getByText('Tab Two'));
    expect(onChange).toHaveBeenCalledWith('tab2');
  });

  it('highlights active tab', () => {
    render(
      <TabSwitch
        tabs={mockTabs}
        activeTab="tab1"
        onChange={vi.fn()}
      />
    );

    const activeButton = screen.getByText('Tab One').closest('button');
    expect(activeButton.className).toContain('bg-indigo-500');
  });

  it('does not highlight inactive tab', () => {
    render(
      <TabSwitch
        tabs={mockTabs}
        activeTab="tab1"
        onChange={vi.fn()}
      />
    );

    const inactiveButton = screen.getByText('Tab Two').closest('button');
    expect(inactiveButton.className).not.toContain('bg-indigo-500');
  });

  it('renders tabs without icons', () => {
    const tabsWithoutIcons = [
      { id: 'tab1', label: 'Tab One' },
      { id: 'tab2', label: 'Tab Two' },
    ];

    render(
      <TabSwitch
        tabs={tabsWithoutIcons}
        activeTab="tab1"
        onChange={vi.fn()}
      />
    );

    expect(screen.getByText('Tab One')).toBeInTheDocument();
    expect(screen.getByText('Tab Two')).toBeInTheDocument();
  });

  it('updates when activeTab prop changes', () => {
    const { rerender } = render(
      <TabSwitch
        tabs={mockTabs}
        activeTab="tab1"
        onChange={vi.fn()}
      />
    );

    let activeButton = screen.getByText('Tab One').closest('button');
    expect(activeButton.className).toContain('bg-indigo-500');

    rerender(
      <TabSwitch
        tabs={mockTabs}
        activeTab="tab2"
        onChange={vi.fn()}
      />
    );

    activeButton = screen.getByText('Tab Two').closest('button');
    expect(activeButton.className).toContain('bg-indigo-500');
  });
});
