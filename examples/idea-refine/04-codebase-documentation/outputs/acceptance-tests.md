# Acceptance Tests: Payment Refactoring

## Test: Existing Credit Card Processing Unchanged

```javascript
it('should process credit cards exactly as before refactoring', async () => {
  const payment = await processPayment({
    method: 'CREDIT_CARD',
    amount: 9999,
    card: testCard
  });
  
  expect(payment.status).toBe('SUCCESS');
  expect(payment.provider).toBe('STRIPE');
  expect(payment.processingTime).toBeLessThan(2000);
  
  // Verify Stripe API called with same parameters as legacy code
  expect(stripeAPI).toHaveBeenCalledWith(expectedLegacyFormat);
});
```
