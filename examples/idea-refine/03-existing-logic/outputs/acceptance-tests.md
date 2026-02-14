# Acceptance Tests: Ticket Routing

## Test: VIP Customer Priority Routing

```javascript
it('should route VIP customer tickets to senior agents', async () => {
  const ticket = await createTicket({ 
    customer: vipCustomer, 
    priority: 'HIGH' 
  });
  
  const assignment = await waitForRouting(ticket.id);
  
  expect(assignment.agent.level).toBe('SENIOR');
  expect(assignment.queueTime).toBeLessThan(60000); // <1 min
});
```
