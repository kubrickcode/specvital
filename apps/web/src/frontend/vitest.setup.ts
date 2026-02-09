import "@testing-library/jest-dom/vitest";

// Mock ResizeObserver for cmdk and other components
class ResizeObserverMock {
  observe() {}
  unobserve() {}
  disconnect() {}
}

globalThis.ResizeObserver = ResizeObserverMock;

// Mock scrollIntoView for cmdk
Element.prototype.scrollIntoView = () => {};
