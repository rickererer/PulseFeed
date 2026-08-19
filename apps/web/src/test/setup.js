import "@testing-library/jest-dom";

// jsdom 不实现 matchMedia，Material Symbols / 组件库依赖时兜底
if (!window.matchMedia) {
  window.matchMedia = () => ({
    matches: false,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {}
  });
}

// 清理 localStorage，保证每个用例独立
beforeEach(() => {
  window.localStorage.clear();
  window.history.pushState({}, "", "/");
});
