import { test, expect } from '@playwright/test';

// Helper to log in as Ash
async function loginAsAsh(page) {
  await page.goto('/login');
  await page.getByTestId('email-input').fill('ash@pokemon.com');
  await page.getByTestId('password-input').fill('pikachu123');
  await page.getByTestId('login-submit').click();
  await page.waitForURL('/store');
}

test.describe('PokéStore', () => {
  test.describe('Home Page', () => {
    test('should display the landing page', async ({ page }) => {
      await page.goto('/');
      await expect(page.locator('.hero-title')).toBeVisible();
      await expect(page.getByTestId('browse-btn')).toBeVisible();
      await expect(page.getByTestId('login-btn')).toBeVisible();
    });
  });

  test.describe('Authentication', () => {
    test('should show login page with demo accounts', async ({ page }) => {
      await page.goto('/login');
      await expect(page.getByTestId('login-form')).toBeVisible();
      await expect(page.getByTestId('demo-ash')).toBeVisible();
      await expect(page.getByTestId('demo-misty')).toBeVisible();
      await expect(page.getByTestId('demo-brock')).toBeVisible();
    });

    test('should fill credentials from demo account button', async ({ page }) => {
      await page.goto('/login');
      await page.getByTestId('demo-ash').click();
      await expect(page.getByTestId('email-input')).toHaveValue('ash@pokemon.com');
      await expect(page.getByTestId('password-input')).toHaveValue('pikachu123');
    });

    test('should show error on invalid credentials', async ({ page }) => {
      await page.goto('/login');
      await page.getByTestId('email-input').fill('wrong@email.com');
      await page.getByTestId('password-input').fill('wrongpass');
      await page.getByTestId('login-submit').click();
      await expect(page.getByTestId('login-error')).toBeVisible();
    });

    test('should log in successfully and redirect to store', async ({ page }) => {
      await loginAsAsh(page);
      await expect(page.getByTestId('user-name')).toHaveText('Ash Ketchum');
    });
  });

  test.describe('Store', () => {
    test.beforeEach(async ({ page }) => {
      await loginAsAsh(page);
    });

    test('should display Pokémon cards', async ({ page }) => {
      await expect(page.getByTestId('pokemon-grid')).toBeVisible();
      await expect(page.getByTestId('results-info')).toContainText('25 Pokémon found');
    });

    test('should search Pokémon by name', async ({ page }) => {
      await page.getByTestId('search-input').fill('pikachu');
      await page.waitForTimeout(500); // debounce
      await expect(page.getByTestId('results-info')).toContainText('1 Pokémon found');
      await expect(page.getByTestId('pokemon-card-25')).toBeVisible();
    });

    test('should filter by type', async ({ page }) => {
      await page.getByTestId('filter-fire').click();
      await expect(page.getByTestId('pokemon-card-4')).toBeVisible(); // Charmander
    });

    test('should add Pokémon to cart', async ({ page }) => {
      await page.getByTestId('add-to-cart-25').click(); // Add Pikachu
      await expect(page.getByTestId('add-to-cart-25')).toHaveText('Added!');
    });
  });

  test.describe('Cart', () => {
    test.beforeEach(async ({ page }) => {
      await loginAsAsh(page);
      // Add Pikachu to cart
      await page.getByTestId('add-to-cart-25').click();
      await page.waitForTimeout(500);
    });

    test('should show items in cart', async ({ page }) => {
      await page.goto('/cart');
      await expect(page.getByTestId('cart-item-25')).toBeVisible();
      await expect(page.getByTestId('cart-item-name')).toHaveText('Pikachu');
    });

    test('should remove item from cart', async ({ page }) => {
      await page.goto('/cart');
      await page.getByTestId('remove-25').click();
      await expect(page.getByTestId('empty-cart')).toBeVisible();
    });

    test('should navigate to checkout', async ({ page }) => {
      await page.goto('/cart');
      await page.getByTestId('checkout-btn').click();
      await page.waitForURL('/checkout');
      await expect(page.getByTestId('checkout-form')).toBeVisible();
    });
  });

  test.describe('Checkout', () => {
    test.beforeEach(async ({ page }) => {
      await loginAsAsh(page);
      await page.getByTestId('add-to-cart-25').click();
      await page.waitForTimeout(500);
    });

    test('should complete checkout successfully', async ({ page }) => {
      await page.goto('/checkout');
      await page.getByTestId('delivery-address').fill('Pallet Town, Route 1');
      await page.getByTestId('card-number').fill('4242 4242 4242 4242');
      await page.getByTestId('card-expiry').fill('12/28');
      await page.getByTestId('card-cvv').fill('123');
      await page.getByTestId('place-order-btn').click();

      await expect(page.getByTestId('order-success')).toBeVisible();
      await expect(page.getByTestId('order-id')).toBeVisible();
      await expect(page.getByTestId('order-status')).toHaveText('confirmed');
    });
  });
});
