# Acceptance Tests: Mobile Push Notifications

## Overview

These are example acceptance tests for **User Story 2: Task Assignment Notifications**. This demonstrates the ATDD (Acceptance Test-Driven Development) approach where tests are written BEFORE implementation.

---

## User Story Reference

**US-2: Task Assignment Notifications**

**As a** Project Manager  
**I want to** receive a push notification when a task is assigned to me  
**So that** I can prioritize new work immediately

---

## Acceptance Criteria → Test Mapping

| Criterion | Test(s) |
|-----------|---------|
| Notification appears within 30 seconds of task assignment | Test 1, Test 2 |
| Notification shows task title and assigner name | Test 3 |
| Clicking notification opens the task detail page | Test 4 |
| Notification badge shows unread count | Test 5 |
| Works on desktop and mobile browsers | Test 6, Test 7 |

---

## RED Phase: Write Failing Tests

### Test 1: Notification Delivered Within 30 Seconds

```javascript
describe('Task Assignment Notification - Delivery Timing', () => {
  it('should deliver notification within 30 seconds of task assignment', async () => {
    // ARRANGE
    const projectManager = await createUser({ role: 'PM' });
    const assigner = await createUser({ role: 'TEAM_MEMBER' });
    await enableNotifications(projectManager);
    
    const startTime = Date.now();
    
    // ACT
    const task = await assignTask({
      title: 'Update landing page',
      assignee: projectManager.id,
      assigner: assigner.id
    });
    
    const notification = await waitForNotification(projectManager.id, {
      timeout: 30000 // 30 seconds
    });
    
    const deliveryTime = Date.now() - startTime;
    
    // ASSERT
    expect(notification).toBeDefined();
    expect(deliveryTime).toBeLessThan(30000); // Less than 30 seconds
    expect(notification.type).toBe('TASK_ASSIGNED');
    expect(notification.taskId).toBe(task.id);
  });
});
```

**Expected Result (RED):** ❌ Test fails - notification delivery system not implemented

---

### Test 2: No Notification if User Disabled Task Assignments

```javascript
it('should not send notification if user disabled task assignments in preferences', async () => {
  // ARRANGE
  const projectManager = await createUser({ role: 'PM' });
  await enableNotifications(projectManager);
  await setNotificationPreference(projectManager.id, {
    taskAssignments: false // User disabled this type
  });
  
  // ACT
  await assignTask({
    assignee: projectManager.id,
    assigner: someUser.id
  });
  
  await wait(5000); // Wait 5 seconds
  
  const notifications = await getNotifications(projectManager.id);
  
  // ASSERT
  expect(notifications).toHaveLength(0); // No notifications sent
});
```

**Expected Result (RED):** ❌ Test fails - preference enforcement not implemented

---

### Test 3: Notification Contains Task Title and Assigner Name

```javascript
it('should include task title and assigner name in notification body', async () => {
  // ARRANGE
  const projectManager = await createUser({ 
    role: 'PM',
    name: 'Paula Chen'
  });
  const assigner = await createUser({ 
    name: 'Marcus Williams'
  });
  await enableNotifications(projectManager);
  
  // ACT
  await assignTask({
    title: 'Fix critical bug in checkout',
    assignee: projectManager.id,
    assigner: assigner.id
  });
  
  const notification = await waitForNotification(projectManager.id);
  
  // ASSERT
  expect(notification.title).toContain('Task Assigned');
  expect(notification.body).toContain('Fix critical bug in checkout');
  expect(notification.body).toContain('Marcus Williams');
  
  // Example format: "Marcus Williams assigned you: Fix critical bug in checkout"
});
```

**Expected Result (RED):** ❌ Test fails - notification formatting not implemented

---

### Test 4: Clicking Notification Opens Task Detail Page

```javascript
it('should open task detail page when notification clicked', async () => {
  // ARRANGE
  const projectManager = await createUser({ role: 'PM' });
  await enableNotifications(projectManager);
  
  const task = await assignTask({
    assignee: projectManager.id,
    assigner: someUser.id
  });
  
  const notification = await waitForNotification(projectManager.id);
  
  // ACT
  const clickResult = await simulateNotificationClick(notification);
  
  // ASSERT
  expect(clickResult.url).toBe(`/tasks/${task.id}`);
  expect(clickResult.focused).toBe(true); // Window/tab focused
  expect(notification.read).toBe(true); // Marked as read
});
```

**Expected Result (RED):** ❌ Test fails - click handling not implemented

---

### Test 5: Notification Badge Shows Unread Count

```javascript
it('should increment notification badge count for unread notifications', async () => {
  // ARRANGE
  const projectManager = await createUser({ role: 'PM' });
  await enableNotifications(projectManager);
  
  const initialBadge = await getBadgeCount(projectManager.id);
  expect(initialBadge).toBe(0);
  
  // ACT
  await assignTask({ assignee: projectManager.id, assigner: user1.id });
  await assignTask({ assignee: projectManager.id, assigner: user2.id });
  await assignTask({ assignee: projectManager.id, assigner: user3.id });
  
  await wait(2000); // Wait for delivery
  
  const updatedBadge = await getBadgeCount(projectManager.id);
  
  // ASSERT
  expect(updatedBadge).toBe(3); // 3 unread notifications
});
```

**Expected Result (RED):** ❌ Test fails - badge count tracking not implemented

---

### Test 6: Works on Desktop Chrome

```javascript
it('should deliver notification on desktop Chrome browser', async () => {
  // ARRANGE
  const browser = await launchBrowser('chrome', 'desktop');
  const page = await browser.newPage();
  await page.goto('https://notifyhub.local');
  
  const projectManager = await createUser({ role: 'PM' });
  await page.login(projectManager);
  await page.grantNotificationPermission();
  
  // ACT
  await assignTask({
    assignee: projectManager.id,
    assigner: someUser.id
  });
  
  const notification = await page.waitForNotification({ timeout: 30000 });
  
  // ASSERT
  expect(notification).toBeDefined();
  expect(notification.visible).toBe(true);
  
  await browser.close();
});
```

**Expected Result (RED):** ❌ Test fails - Chrome desktop not implemented

---

### Test 7: Works on Mobile Safari

```javascript
it('should deliver notification on mobile Safari (iOS 16.4+)', async () => {
  // ARRANGE
  const browser = await launchBrowser('safari', 'mobile');
  const page = await browser.newPage();
  await page.goto('https://notifyhub.local');
  
  const projectManager = await createUser({ role: 'PM' });
  await page.login(projectManager);
  await page.grantNotificationPermission();
  
  // ACT
  await assignTask({
    assignee: projectManager.id,
    assigner: someUser.id
  });
  
  const notification = await page.waitForNotification({ timeout: 30000 });
  
  // ASSERT
  expect(notification).toBeDefined();
  expect(notification.visible).toBe(true);
  
  await browser.close();
});
```

**Expected Result (RED):** ❌ Test fails - Mobile Safari not implemented

---

## GREEN Phase: Implement Until Tests Pass

After writing these failing tests:

1. **Implement notification trigger detection** - Detect when task assigned
2. **Build delivery queue** - Queue notification for delivery
3. **Integrate web-push** - Send push to user's subscription
4. **Add click handling** - Deep link to task page
5. **Implement badge counting** - Track unread notifications
6. **Test cross-browser** - Ensure works on Chrome, Safari, Firefox, Edge

**Goal:** All 7 tests turn green ✅

---

## REFACTOR Phase: Improve Without Breaking Tests

After tests pass, refactor for:

- Performance (delivery <30s → <10s)
- Code clarity (extract notification formatting logic)
- Error handling (retry failed deliveries)
- Monitoring (log delivery metrics)

**Rule:** Tests must stay green throughout refactoring.

---

## Running the Tests

```bash
# Run all acceptance tests for this user story
npm test -- acceptance/task-assignment-notifications.spec.js

# Run specific test
npm test -- acceptance/task-assignment-notifications.spec.js -t "should deliver notification within 30 seconds"

# Run in watch mode during development
npm test -- --watch acceptance/task-assignment-notifications.spec.js
```

---

## Non-Technical Explanation

### What These Tests Verify

1. **Speed:** Notifications arrive quickly (under 30 seconds)
2. **Content:** Users see the right information (task title, who assigned it)
3. **Navigation:** Clicking goes to the right place
4. **Preferences:** Users can turn off notifications they don't want
5. **Compatibility:** Works on different devices and browsers

### Why We Write Tests First (ATDD)

**Traditional Approach:**
```
Build feature → Hope it works → Find bugs → Fix bugs → Repeat
```

**ATDD Approach:**
```
Define success criteria → Write tests → Tests fail → Build feature → Tests pass → Done
```

**Benefits:**
- **No surprises:** We know exactly what "done" looks like
- **Fewer bugs:** Tests catch issues before users see them
- **Living documentation:** Tests show how features actually work
- **Confidence:** Refactoring is safe because tests prevent regressions

### The RED-GREEN-REFACTOR Cycle

1. **RED:** Write a test that fails (because feature doesn't exist yet)
2. **GREEN:** Write just enough code to make the test pass
3. **REFACTOR:** Improve the code while keeping tests green

This cycle ensures we build exactly what was specified in the acceptance criteria - nothing more, nothing less.

---

[← Back to Dev Plan](dev-plan.md) | [View Learnings →](../learnings.md)
