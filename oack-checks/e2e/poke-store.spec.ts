import { test, expect, type Page } from "@playwright/test";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const DEMO_USER = {
  email: "ash@pokemon.com",
  password: "pikachu123",
  name: "Ash Ketchum",
};

async function login(page: Page) {
  await page.goto("/login");
  await page.getByTestId("email-input").fill(DEMO_USER.email);
  await page.getByTestId("password-input").fill(DEMO_USER.password);
  await page.getByTestId("login-submit").click();
  await page.waitForURL(/\/store/);
}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

test.describe("Login", () => {
  test("shows login form with demo account buttons", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByTestId("login-form")).toBeVisible();
    await expect(page.getByTestId("demo-ash")).toBeVisible();
    await expect(page.getByTestId("demo-misty")).toBeVisible();
    await expect(page.getByTestId("demo-brock")).toBeVisible();
  });

  test("rejects invalid credentials", async ({ page }) => {
    await page.goto("/login");
    await page.getByTestId("email-input").fill("bad@example.com");
    await page.getByTestId("password-input").fill("wrong");
    await page.getByTestId("login-submit").click();
    await expect(page.getByTestId("login-error")).toBeVisible();
  });

  test("logs in with valid credentials and redirects to store", async ({
    page,
  }) => {
    await login(page);
    await expect(page.getByTestId("user-name")).toHaveText(DEMO_USER.name);
  });
});

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

test.describe("Item search", () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test("displays all pokemon by default", async ({ page }) => {
    await expect(page.getByTestId("results-info")).toContainText("25");
  });

  test("searches by name", async ({ page }) => {
    await page.getByTestId("search-input").fill("pikachu");
    await expect(page.getByTestId("results-info")).toContainText("1");
    await expect(
      page.getByTestId("pokemon-grid").getByTestId("pokemon-name"),
    ).toHaveText("Pikachu");
  });

  test("filters by type", async ({ page }) => {
    await page.getByTestId("filter-fire").click();
    const names = page
      .getByTestId("pokemon-grid")
      .getByTestId("pokemon-name");
    await expect(names.first()).toBeVisible();
    // Fire type should include Charmander
    await expect(page.getByTestId("pokemon-grid")).toContainText("Charmander");
  });
});

// ---------------------------------------------------------------------------
// Cart — add & remove
// ---------------------------------------------------------------------------

test.describe("Cart", () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test("adds item to cart and shows badge", async ({ page }) => {
    // Add Bulbasaur (id=1)
    await page.getByTestId("add-to-cart-1").click();
    await expect(page.getByTestId("add-to-cart-1")).toContainText("Added!");
    await expect(page.locator("#cart-badge")).toBeVisible();
  });

  test("shows item in cart page", async ({ page }) => {
    await page.getByTestId("add-to-cart-1").click();
    await page.getByTestId("nav-cart").click();
    await page.waitForURL(/\/cart/);
    await expect(page.getByTestId("cart-item-1")).toBeVisible();
    await expect(
      page.getByTestId("cart-item-1").getByTestId("cart-item-name"),
    ).toContainText("Bulbasaur");
  });

  test("removes item from cart", async ({ page }) => {
    await page.getByTestId("add-to-cart-1").click();
    await page.getByTestId("nav-cart").click();
    await page.waitForURL(/\/cart/);
    await page.getByTestId("remove-1").click();
    await expect(page.getByTestId("empty-cart")).toBeVisible();
  });
});

// ---------------------------------------------------------------------------
// Checkout
// ---------------------------------------------------------------------------

test.describe("Checkout", () => {
  test("completes full checkout flow", async ({ page }) => {
    await login(page);

    // Add an item
    await page.getByTestId("add-to-cart-4").click(); // Charmander
    await page.getByTestId("nav-cart").click();
    await page.waitForURL(/\/cart/);
    await page.getByTestId("checkout-btn").click();
    await page.waitForURL(/\/checkout/);

    // Fill checkout form
    await page.getByTestId("delivery-address").fill("123 Pallet Town");
    await page.getByTestId("delivery-region").selectOption("kanto");
    await page.getByTestId("poke-code").fill("00100");
    await page.getByTestId("card-number").fill("4242424242424242");
    await page.getByTestId("card-expiry").fill("12/30");
    await page.getByTestId("card-cvv").fill("123");
    await page.getByTestId("place-order-btn").click();

    // Verify order confirmation
    await expect(page.getByTestId("order-success")).toBeVisible();
    await expect(page.getByTestId("order-id")).toContainText("ORD-");
    await expect(page.getByTestId("order-status")).toContainText("confirmed");
  });
});

// ---------------------------------------------------------------------------
// Logout
// ---------------------------------------------------------------------------

test.describe("Logout", () => {
  test("logs out and redirects to login", async ({ page }) => {
    await login(page);
    await page.getByTestId("logout-btn").click();
    await page.waitForURL(/\/login/);
    await expect(page.getByTestId("login-form")).toBeVisible();
  });

  test("cannot access store after logout", async ({ page }) => {
    await login(page);
    await page.getByTestId("logout-btn").click();
    await page.waitForURL(/\/login/);
    await page.goto("/store");
    // Should redirect back to login
    await page.waitForURL(/\/login/);
  });
});
