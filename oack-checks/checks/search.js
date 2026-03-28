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

  await context.step("Verify default catalog", async () => {
    const info = await page.textContent('[data-testid="results-info"]');
    if (!info || !info.includes("25")) {
      throw new Error("Expected 25 Pokémon, got: " + info);
    }
  });

  await context.step("Search by name", async () => {
    await page.fill('[data-testid="search-input"]', "pikachu");
    // Wait for debounce + response
    await page.waitForFunction(
      () => {
        const el = document.querySelector('[data-testid="results-info"]');
        return el && el.textContent.includes("1");
      },
      { timeout: 5000 },
    );
    const name = await page.textContent(
      '[data-testid="pokemon-grid"] [data-testid="pokemon-name"]',
    );
    if (!name || !name.includes("Pikachu")) {
      throw new Error("Expected Pikachu in results, got: " + name);
    }
  });

  await context.screenshot("search-pikachu");

  await context.step("Filter by type", async () => {
    // Clear search first
    await page.fill('[data-testid="search-input"]', "");
    await page.click('[data-testid="filter-fire"]');
    await page.waitForFunction(
      () => {
        const el = document.querySelector('[data-testid="pokemon-grid"]');
        return el && el.textContent.includes("Charmander");
      },
      { timeout: 5000 },
    );
  });

  await context.screenshot("filter-fire");
};
