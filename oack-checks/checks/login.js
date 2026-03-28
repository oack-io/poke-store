module.exports = async function (page, context) {
  const email = context.LOGIN_EMAIL;
  const password = context.LOGIN_PASSWORD;
  const baseURL = context.BASE_URL || "https://poke-store.oack.io";

  await context.step("Navigate to login", async () => {
    await page.goto(baseURL + "/login");
    const form = page.locator('[data-testid="login-form"]');
    await form.waitFor({ state: "visible", timeout: 10000 });
  });

  await context.step("Submit credentials", async () => {
    await page.fill('[data-testid="email-input"]', email);
    await page.fill('[data-testid="password-input"]', password);
    await page.click('[data-testid="login-submit"]');
    await page.waitForURL("**/store", { timeout: 10000 });
  });

  await context.step("Verify logged in", async () => {
    const userName = await page.textContent('[data-testid="user-name"]');
    if (!userName || !userName.includes("Ash")) {
      throw new Error("Expected user name to contain 'Ash', got: " + userName);
    }
  });

  await context.screenshot("login-success");
};
