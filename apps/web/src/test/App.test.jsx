import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import App from "../App.jsx";

// 渲染 /auth 路由，验证登录页正常渲染与表单交互
describe("App 冒烟测试", () => {
  beforeEach(() => {
    window.history.pushState({}, "", "/auth");
  });

  it("在 /auth 路由渲染登录页", () => {
    render(<App />);
    expect(screen.getByText("登录 PulseFeed")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("请输入账号")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("输入密码")).toBeInTheDocument();
  });

  it("可在登录与注册模式间切换", async () => {
    const user = userEvent.setup();
    render(<App />);

    await user.click(screen.getByRole("button", { name: "注册" }));
    expect(screen.getByPlaceholderText("输入昵称")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /注册并登录/ })).toBeInTheDocument();
  });
});
