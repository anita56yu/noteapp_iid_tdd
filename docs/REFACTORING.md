# Refactoring Process

This document outlines the steps for refactoring code to improve its internal structure without changing its external behavior. Refactoring is typically identified during the **Review** phase of IID.

## 1. Define the Goal
-   Clearly state the "code smell" or design issue that needs to be addressed (e.g., "The use case layer is leaking repository errors").
-   Define the desired outcome (e.g., "The use case will translate repository errors into its own error types").
-   Create a new task in `./FEATURES.md` for this refactoring effort.

## 2. Verify the Safety Net
-   Before making any changes, run the entire test suite (`go test ./...`) to ensure all existing tests are passing. This confirms you are starting from a stable, green state.

## 3. Plan the Refactoring
-   Break the refactoring down into the smallest possible, incremental changes. Each change should be a self-contained piece of the overall refactoring goal.
-   **Example:** For translating errors, the plan would be to handle `ErrNoteNotFound`, `ErrNilNote`, and `ErrEmptyTitle` one by one.

## 4. Execute Using TDD
-   For each incremental change identified in the plan, follow the complete TDD loop as described in **`./TDD.md`**. This includes:
    1.  **Write a failing test:** Modify an existing test or create a new one that fails until the refactoring is implemented.
    2.  **Write code to pass the test:** Make the minimal code change to make the test pass.
    3.  **Review and Refactor:** Review the change and refactor further if needed, ensuring all tests continue to pass.

## 5. Commit
-   Once all incremental changes are complete and all tests are passing, follow the **Commit** step in **`./TDD.md`**. This includes marking the task as complete and committing the changes with a `refactor` type commit message.
