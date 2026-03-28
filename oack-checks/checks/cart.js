module.exports = async function (page, context) {
  const email = context.LOGIN_EMAIL;
  const password = context.LOGIN_PASSWORD;
  const baseURL = context.BASE_URL || "https://poke-store.oack.io";

  await context.step("Login", async () => {
    await page.goto(baseURL + "/login");
    await page.fill('[data-testid="email-input"]', email);
    await page.fill('[data-testid="password-input"]', password);
    await page.click('[data-testid="login-submit"]');
    await page.waitForURL("**/store", { timeout: 10000 });
  });

  await context.step("Add item to cart", async () => {
    await page.click('[data-testid="add-to-cart-1"]');
    // Wait for button to change to "Added!"
    await page.waitForFunction(
      () => {
        const btn = document.querySelector('[data-testid="add-to-cart-1"]');
        return btn && btn.textContent.includes("Added!");
      },
      { timeout: 5000 },
    );
    // Verify cart badge appears
    const badge = page.locator('[data-testid="cart-badge"]');
    await badge.waitFor({ state: "visible", timeout: 5000 });
  });

  await context.step("View cart", async () => {
    await page.click('[data-testid="nav-cart"]');
    await page.waitForURL("**/cart", { timeout: 10000 });
    const item = page.locator('[data-testid="cart-item-1"]');
    await item.waitFor({ state: "visible", timeout: 5000 });
    const name = await page.textContent(
      '[data-testid="cart-item-1"] [data-testid="cart-item-name"]',
    );
    if (!name || !name.includes("Bulbasaur")) {
      throw new Error("Expected Bulbasaur in cart, got: " + name);
    }
  });

  await context.screenshot("cart-with-item");

  await context.step("Remove item from cart", async () => {
    await page.click('[data-testid="remove-1"]');
    const empty = page.locator('[data-testid="empty-cart"]');
    await empty.waitFor({ state: "visible", timeout: 5000 });
  });

  await context.screenshot("cart-empty");
};
