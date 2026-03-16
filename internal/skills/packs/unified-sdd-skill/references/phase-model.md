# Phase Model

The default workflow is:

```text
[Platform] -> [Route] -> [Specify] -> [Plan] -> [Deliver]
```

Iteration rollout:

```text
Iteration 1
[Platform] -> [Route] -> [Specify]

Iteration 2
[Plan] -> [Deliver]
```

Deliver flow:

```text
[Deliver]
  Build -> Create PR -> Review PR -> Verify -> Deploy -> Archive
```

Future option:

```text
[Plan] -> [Build] -> [Deploy]
```

Only split Build and Deploy when release coordination becomes large enough to
deserve its own phase.
