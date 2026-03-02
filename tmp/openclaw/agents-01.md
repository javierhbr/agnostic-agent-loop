# The OpenClaw Orchestrator Pattern: Configuring Nested Sub-Agent Hierarchies

To set up the orchestrator sub-agent pattern in OpenClaw, you need to enable nested sub-agents so your main agent can spawn an orchestrator, which in turn spawns its own worker sub-agents.

Here is the step-by-step process to configure and use this pattern:

Step 1: Enable Nested Sub-Agents in Configuration

By default, OpenClaw limits sub-agents to a single layer (`maxSpawnDepth: 1`), meaning a sub-agent cannot spawn its own sub-agents. To enable the orchestrator pattern, you must modify your configuration file to allow nesting:

- Set **maxSpawnDepth: 2** in your configuration. This allows the specific orchestrator hierarchy: Main Agent → Orchestrator Sub-Agent → Worker Sub-Sub-Agents.

Step 2: Understand the Depth Levels and Granted Tools

Once nesting is enabled, OpenClaw automatically changes the tool permissions based on the depth of the agent:

- **Depth 0 (Main Agent):** This is your primary session that always has full permissions.
- **Depth 1 (The Orchestrator):** Because `maxSpawnDepth` is set to 2 or higher, Depth 1 sub-agents automatically receive orchestration tools: `sessions_spawn`, `subagents`, `sessions_list`, and `sessions_history`. This gives the orchestrator the ability to manage its own children.
- **Depth 2 (Leaf Workers):** These are the sub-sub-agents doing the actual work. They are stripped of session tools (meaning `sessions_spawn` is denied) so they cannot spawn any further children.

Step 3: Spawn the Orchestrator

To kick off the workflow, use the `/subagents` slash command or the `sessions_spawn` tool from your main chat.

- Command example: `/subagents spawn <agentId> <task>`.
- Once spawned, the orchestrator (Depth 1) will begin processing your task and can use its `sessions_spawn` tool to delegate specific parts of the task to Depth 2 worker agents in the background.

Step 4: The Announce Chain (Reporting Back)

You do not need to manually check on the Depth 2 workers; results flow back up the chain automatically:

1. When a **Depth 2 worker** finishes its task, it announces its result directly to its parent (the Depth 1 orchestrator).
2. The **Depth 1 orchestrator** receives all the announcements from its workers, synthesizes the final results, and finishes its run.
3. The orchestrator then announces the final synthesized result back to your **Main Agent**, which delivers it to you in the user chat.

Safety Limits and Controls

To prevent runaway compute costs or infinite fan-out loops, OpenClaw enforces strict limits on orchestrators:

- **Concurrency Limits:** Each agent session (including the orchestrator) is limited by `maxChildrenPerAgent`, which defaults to 5 active children at a time.
- **Cascade Stopping:** If an orchestrator goes off track, you can stop it using `/subagents kill <id>` (or `/stop` in the main chat). This triggers a "cascade stop," which will instantly abort the orchestrator and all of its active depth-2 children.

# The App Factory: Eleven Agents of Autonomous Development

specialized AI agents. To avoid context bloat, the workload is distributed so that each agent handles a specific phase of the app development lifecycle.

Here are the 11 AI agents and their specific roles in building and shipping apps:

**1. Sheldon (The Orchestrator)** Sheldon is the central manager of the factory. Instead of doing the heavy lifting, Sheldon uses less than 5% of its context window to sit on top of the system, read the project state files, and direct traffic. It spawns the right sub-agents for the appropriate phase and passes tasks down the chain.

**2. Shan (The Researcher)** Running on a 5-minute cron job, Shan autonomously scans platforms like Reddit and X to identify user pain points and broken workflows. It cross-references these issues with App Store categories to find high-demand, low-competition opportunities and generates a "one-pager" pitch for a new app.

**3. The Validation Agent** Once the research is complete, Sheldon passes the one-pager to the Validation Agent. This agent checks the research to ensure the idea is technically feasible and validates the app concept before any code is written.

**4. The Builder (Claude Opus 4.6)** Tasked with the actual coding, the Builder takes the validated idea and writes the entire Swift and SwiftUI codebase from a single prompt. It incorporates pre-configured templates for features like Apple StoreKit (for payments) and Gemini Flash AI wrappers to power the app's internal AI functions.

**5. The Reviewer (Codex 5.3)** Because an AI model can have a bias and cut corners if it reviews its own code, the Reviewer agent independently reads every file created by the Builder. It specifically checks for crash risks, missing features, permission bugs, and overall code quality.

**6. The Quality Control Agent** This agent runs the app through six automated quality checks, acting as a strict quality gate. The app must score at least an 8/10 to pass to the next phase. If the app fails three times, this agent flags it for manual human review.

**7. The Monetization Agent** This agent analyzes the initial validation document and user demographic to pick the best payment strategy (such as a 7-day free trial or a premium paywall). It then autonomously hooks up Apple StoreKit to implement the chosen strategy.

**8. The App Store Packaging Agent** This agent prepares the app for the public. It generates the App Store listing, writes the description, and selects the best keywords. It also navigates the finished app to take automated screenshots and uses a tool called Nano Banana Pro to design a unique app icon.

**9. The Onboarding Agent** Because user retention is critical, this agent generates 3 to 5 custom onboarding screens for the app. These screens are designed to explain the app's core features to new users as soon as they open it.

**10. The Promo Video Agent** Before the app hits the store, this agent uses a tool called Votion (integrated inside Claude Code) to automatically generate a high-quality demo video showcasing the app's features.

**11. Larry (The Social Media Marketer)** Larry is an open-source OpenClaw agent (created by Oliver Henry) dedicated strictly to distribution. Larry manages multiple niche TikTok and Instagram theme pages. When an app ships, Larry autonomously generates hooks, images, and slideshows to natively promote the app. It uses self-learning loops to analyze which hooks get views and optimize future marketing content.

_Note: Once all the agents finish their tasks, the system queues the app to the Apple Developer portal. Pressing "submit" is the only manual step required by a human in this entire pipeline__._


# The 11-Agent Autonomous App Factory Blueprint

To clone and set up your own 11-agent autonomous app factory running 24/7, you need to combine the high-volume production pipeline of an "App Factory" with OpenClaw's official **Multi-Agent Routing**, **Sub-Agents** system, and robust coordination tools like **Agent Mail** or **Network-AI**.

Here is the step-by-step blueprint to build, coordinate, and launch your multi-agent team:

Step 1: Set Up the Dedicated Gateway & Environment

Do not run this 24/7 system on your personal daily-driver computer or using your personal AI subscription, as it poses security risks and can violate terms of service.

- **Hardware:** Use a dedicated Mac Mini or a Cloud VPS to run the OpenClaw Gateway continuously.
- **API Management:** Use a service like OpenRouter to manage your API tokens. This allows you to route specific tasks to different models and track costs separately from your personal accounts.
- **Enable Multi-Agent Mode:** OpenClaw can host multiple isolated agents side-by-side on one Gateway. Use the command `openclaw agent add <agentId>` to create a unique workspace, `agentDir`, and session store for each of your 11 agents.

Step 2: Clone the 11 Agents and Assign Roles

Create separate identities for your agents. You can bind each agent to its own Slack bot or Discord bot using `accountId` bindings in your `openclaw.json` file. While you can give them unique personalities in their `IDENTITY.md` or `SOUL.md` files, remember the core rule of agent swarms: **agents are fundamentally fungible generalists**. Simply naming an agent a "frontend developer" doesn't make it better at frontend code; its capability comes from the model you assign it.

Assign the following models to your 11 agents based on their pipeline roles:

- **The Builder (Claude Opus 4.6):** Writes the Swift/SwiftUI code.
- **The Reviewer (Codex 5.3):** Independently verifies code for bugs and crashes. _Never use the same model to build and review_.
- **The Marketer/Router (Sonnet 4.6):** Handles fast routing, metadata logging, and marketing.
- **The Wrapper (Gemini Flash):** Embedded inside the generated apps to power the user-facing AI features affordably.

Step 3: Coordinate via Orchestrator & Shared State

If 11 agents work on the same app simultaneously without coordination, they will suffer from split-brain failures and overwrite each other's files.

- **Set Up the Orchestrator (Sub-Agents):** Do not let one agent build the whole app. Use an orchestrator agent (e.g., "Sheldon"). In OpenClaw, configure `maxSpawnDepth: 2` to enable the orchestrator pattern. The main orchestrator can use the `/subagents spawn <agentId> <task>` command to delegate specific phases to your sub-agents in the background.
- **Prevent Collisions with Atomic State:** Integrate a coordination layer like **Network-AI** or **Agent Mail**. These systems provide a "Shared Blackboard" with file-system mutexes (atomic commits).
- **File Reservations:** Before an agent edits a file, it must register with Agent Mail, reserve the file path, announce its work via a thread ID, and release the reservation when finished so no other agent conflicts with it.

Step 4: Execute the 9-Step App Factory Loop

Once your orchestrator and state management are running, execute this continuous pipeline:

1. **Research:** A research agent runs on a 5-minute cron job to scan X, Reddit, and the App Store for pain points, generating a one-page pitch.
2. **Validation:** A validation sub-agent reviews the pitch for feasibility.
3. **Building:** The Builder agent (Opus 4.6) codes the app using templates pre-configured with Apple StoreKit.
4. **Reviewing:** The Reviewer agent (Codex 5.3) verifies the code.
5. **Quality Gates:** The app goes through 6 automated checks. It must score an 8/10. If it fails 3 times, it is flagged for manual human review.
6. **Payments:** The system selects a monetization strategy (e.g., 7-day free trial or premium paywall) and hooks up Apple StoreKit.
7. **Packaging:** Agents generate the app store listing, description, and keywords. Integrate **Nano Banana Pro** to auto-generate the app icon, while the AI navigates the app to take screenshots.
8. **Onboarding:** The factory generates 3 to 5 onboarding screens.
9. **Submission:** The app is queued to the Apple Developer portal. _Pressing "submit" is the only manual step you should perform_.

Step 5: Enforce Quality Loops & Memory Management

To keep the code clean and the agents focused, enforce strict quality loops:

- **Run Until Clean:** Instruct your agents to perform Self-Reviews, Cross-Reviews of other agents' code, and Random Code Exploration. Keep running these loops until they consistently come back with zero bugs.
- **Handle Context Compaction:** When an agent's memory compresses, it can forget its current workflow. Immediately feed it a "Post-Compaction" prompt to re-establish its rules, tool knowledge, and current tasks.

Step 6: Deploy Autonomous Marketing

An app won't sell without marketing. Dedicate agents specifically to distribution:

- **Promo Videos:** Use **Votion** inside Claude Code to automatically generate a demo video for the app before it hits the store.
- **Social Media Theme Pages:** Create TikTok and Instagram accounts for different niches. Integrate the open-source **Larry skill** (an OpenClaw agent) to autonomously generate hooks, images, and slideshows that natively mention your newly shipped apps. It uses self-learning loops to optimize its hooks based on view metrics.

# Building the 11-Agent Autonomous App Factory

To build a fully autonomous 11-agent app factory that runs 24/7, you must combine the high-volume production pipeline of Dubibubii's "App Factory" with the secure, multi-agent team architecture developed by Brian Casel.

Here is the step-by-step guide to setting up, staffing, and coordinating your AI development team using OpenClaw:

Step 1: Secure Dedicated Hardware and Access

Never run a fully autonomous agent team on your personal, daily-driver computer, as it requires 24/7 uptime and poses a security risk to your personal files.

- **Hardware:** Set up a dedicated machine (like a Mac Mini) or a Cloud VPS (starting at $5/month) to run the OpenClaw Gateway continuously.
- **Security & Permissions:** Treat your agents like human employees. Create a dedicated email address, a separate GitHub account, and an isolated Dropbox account strictly for the agents. Only share specific folders between your main computer and the agents' Dropbox to keep the rest of your files walled off.
- **API Management:** Do not use your personal Claude Max plan, as autonomous tasks can violate terms of service or drain your limits. Instead, use OpenRouter to manage API tokens, allowing you to carefully optimize which models your agents use and track your costs.

Step 2: Establish the Orchestrator and Dashboard

If one agent tries to manage an entire app build, its context window will bloat and it will fail. You must build a multi-agent hierarchy.

- **The Orchestrator:** Create a central manager agent (e.g., "Sheldon") that sits at the top of the hierarchy. It should use less than 5% of its context window and act purely to direct traffic and delegate tasks to specialized sub-agents.
- **The Dashboard:** Because OpenClaw's built-in cron system can make it difficult to assign scheduled tasks to specific agents, build a custom dashboard (using Claude Code or Rails) to track active projects, shipped apps, queue lengths, and quality scores.

Step 3: Configure the 11-Agent Team Profiles

Instead of using Telegram, set up your agents as individual Slack bots. Slack provides better Markdown support and allows you to use threaded replies to manage multiple agents simultaneously.

You can place your agents in a shared workspace so they access the same "brain" folder and `AGENTS.md` rules. Use the `IDENTITY.md` file to assign unique personalities and roles to each of your 11 agents. Assign specific models to agents based on their roles:

- **The Builder (Claude Opus 4.6):** Handles heavy reasoning, writing the entire Swift and SwiftUI codebase from a single prompt.
- **The Reviewer (Codex 5.3):** Independently verifies every file for crash risks and permission bugs. _Never use the same model to build and review, or it will cut corners_.
- **The Marketer/Assistant (Sonnet 4.6):** Handles fast routing, logging metadata, and marketing tasks where speed and efficiency matter more than deep reasoning.
- **The App Integration (Gemini Flash):** A cheap, fast wrapper used _inside_ the generated apps to power the user-facing AI features.

Step 4: Coordinate via File-Based State Management

To coordinate 11 agents without them overwriting each other's work or losing track of the project, **do not rely on conversation history for state**.

- **The Project State File:** Have the orchestrator read a shared markdown file to determine what phase a project is in. It then spawns a sub-agent session with the appropriate prompt, the sub-agent does the work, updates the state file, and exits.
- **File Reservations:** To prevent conflicts, agents should register with an "Agent Mail" system or shared blackboard. They must "reserve" specific files before editing them, announce their work via thread IDs, and release the reservations when finished.

Step 5: Execute the 9-Step App Pipeline

Once coordinated, the team runs this continuous loop:

1. **Research:** A research agent (e.g., "Shan") runs on a 5-minute cron job, scanning Reddit, X, and the App Store for pain points to generate a one-page pitch.
2. **Validation:** A validation agent reviews the research to ensure the idea is viable.
3. **Building:** The builder agent codes the app using templates pre-configured with Apple StoreKit.
4. **Reviewing:** The reviewer agent checks the code for errors.
5. **Quality Gates:** The app runs through 6 automated checks. It must score an 8/10. If it fails 3 times, it is flagged for manual human review.
6. **Payments:** The factory analyzes the validation document to choose a monetization strategy (e.g., a 7-day free trial or premium paywall) and hooks up Apple StoreKit.
7. **Packaging:** Agents generate the app store listing, use **Nano Banana Pro** to design a unique app icon, and automatically navigate the app to take screenshots.
8. **Onboarding:** The factory generates 3 to 5 onboarding screens to teach new users how to use the app.
9. **Submission:** The app is queued to the Apple Developer portal. Pressing "submit" is the only manual step you perform.

Step 6: Deploy Autonomous Marketing

Because building the app is only half the battle, the factory must also market it autonomously.

- **Promo Videos:** Use **Votion** inside Claude Code to automatically generate demo videos for the apps.
- **Social Media Distribution:** Set up multiple niche TikTok and Instagram accounts. Integrate the open-source **Larry skill** (an OpenClaw agent) to autonomously generate hooks, images, and slideshows. The Larry agent natively mentions your new apps and uses self-learning loops to optimize its content based on views and conversions.



