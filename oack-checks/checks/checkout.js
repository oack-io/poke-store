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
    await page.click('[data-testid="add-to-cart-4"]');
    await page.waitForFunction(
      () => {
        const btn = document.querySelector('[data-testid="add-to-cart-4"]');
        return btn && btn.textContent.includes("Added!");
      },
      { timeout: 5000 },
    );
  });

  await context.step("Navigate to checkout", async () => {
    await page.click('[data-testid="nav-cart"]');
    await page.waitForURL("**/cart", { timeout: 10000 });
    await page.click('[data-testid="checkout-btn"]');
    await page.waitForURL("**/checkout", { timeout: 10000 });
  });

  await context.step("Fill checkout form", async () => {
    await page.fill('[data-testid="delivery-address"]', "123 Pallet Town");
    await page.selectOption('[data-testid="delivery-region"]', "kanto");
    await page.fill('[data-testid="poke-code"]', "00100");
    await page.fill('[data-testid="card-number"]', "4242424242424242");
    await page.fill('[data-testid="card-expiry"]', "12/30");
    await page.fill('[data-testid="card-cvv"]', "123");
  });

  await context.screenshot("checkout-filled");

  await context.step("Place order", async () => {
    await page.click('[data-testid="place-order-btn"]');
    const success = page.locator('[data-testid="order-success"]');
    await success.waitFor({ state: "visible", timeout: 10000 });
  });

  await context.step("Verify order confirmation", async () => {
    const orderId = await page.textContent('[data-testid="order-id"]');
    if (!orderId || !orderId.includes("ORD-")) {
      throw new Error("Expected order ID starting with ORD-, got: " + orderId);
    }
    const status = await page.textContent('[data-testid="order-status"]');
    if (!status || !status.includes("confirmed")) {
      throw new Error("Expected status 'confirmed', got: " + status);
    }
  });

  await context.screenshot("order-confirmed");
};
