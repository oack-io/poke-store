// Full end-to-end flow: login -> search -> add to cart -> checkout -> logout
module.exports = async function (page, context) {
  const email = context.LOGIN_EMAIL;
  const password = context.LOGIN_PASSWORD;
  const baseURL = context.BASE_URL || "https://poke-store.oack.io";

  // ---- Login ----
  await context.step("Login", async () => {
    await page.goto(baseURL + "/login");
    await page.fill('[data-testid="email-input"]', email);
    await page.fill('[data-testid="password-input"]', password);
    await page.click('[data-testid="login-submit"]');
    await page.waitForURL("**/store", { timeout: 10000 });
    const userName = await page.textContent('[data-testid="user-name"]');
    if (!userName || !userName.includes("Ash")) {
      throw new Error("Expected user name containing 'Ash', got: " + userName);
    }
  });

  // ---- Search ----
  await context.step("Search for Pikachu", async () => {
    await page.fill('[data-testid="search-input"]', "pikachu");
    await page.waitForFunction(
      () => {
        const el = document.querySelector('[data-testid="results-info"]');
        return el && el.textContent.includes("1");
      },
      { timeout: 5000 },
    );
  });

  await context.screenshot("search-results");

  // ---- Add to cart ----
  await context.step("Add Pikachu to cart", async () => {
    await page.click('[data-testid="add-to-cart-25"]');
    await page.waitForFunction(
      () => {
        const btn = document.querySelector('[data-testid="add-to-cart-25"]');
        return btn && btn.textContent.includes("Added!");
      },
      { timeout: 5000 },
    );
  });

  // ---- Cart ----
  await context.step("View cart", async () => {
    await page.click('[data-testid="nav-cart"]');
    await page.waitForURL("**/cart", { timeout: 10000 });
    const item = page.locator('[data-testid="cart-item-25"]');
    await item.waitFor({ state: "visible", timeout: 5000 });
  });

  await context.screenshot("cart-with-pikachu");

  // ---- Checkout ----
  await context.step("Checkout", async () => {
    await page.click('[data-testid="checkout-btn"]');
    await page.waitForURL("**/checkout", { timeout: 10000 });
    await page.fill('[data-testid="delivery-address"]', "123 Pallet Town");
    await page.selectOption('[data-testid="delivery-region"]', "kanto");
    await page.fill('[data-testid="poke-code"]', "00100");
    await page.fill('[data-testid="card-number"]', "4242424242424242");
    await page.fill('[data-testid="card-expiry"]', "12/30");
    await page.fill('[data-testid="card-cvv"]', "123");
    await page.click('[data-testid="place-order-btn"]');
    const success = page.locator('[data-testid="order-success"]');
    await success.waitFor({ state: "visible", timeout: 10000 });
  });

  await context.step("Verify order", async () => {
    const orderId = await page.textContent('[data-testid="order-id"]');
    if (!orderId || !orderId.includes("ORD-")) {
      throw new Error("Expected order ID with ORD- prefix, got: " + orderId);
    }
  });

  await context.screenshot("order-confirmed");

  // ---- Logout ----
  await context.step("Logout", async () => {
    await page.click('[data-testid="back-to-store"]');
    await page.waitForURL("**/store", { timeout: 10000 });
    await page.click('[data-testid="logout-btn"]');
    await page.waitForURL("**/login", { timeout: 10000 });
    const form = page.locator('[data-testid="login-form"]');
    await form.waitFor({ state: "visible", timeout: 5000 });
  });

  await context.screenshot("logged-out");
};
