# Acceptance Tests: Expense Approval

## Test: Automatic Approval for Small Expenses

```javascript
it('should auto-approve expenses under $100', async () => {
  const expense = await submitExpense({ amount: 75, category: 'Meals' });
  const approval = await waitForApproval(expense.id, { timeout: 5000 });
  
  expect(approval.status).toBe('APPROVED');
  expect(approval.approver).toBe('SYSTEM');
  expect(approval.approvalTime).toBeLessThan(5000); // <5 seconds
});
```

*See Scenario 1 for comprehensive test examples*
