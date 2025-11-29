# Test-driven development
This instruction describes the TDD loop you and I will follow to complete a task.

1. Write a failing test, one test at a time:
    - You will help me create the boilerplate test, including a good test name in a suitable test suite.
    - We will then write the **Arrange-Act-Assert (AAA)** test code.
    - You will run the build to confirm it fails.
    - Whenever you try to modify an existing testcase due to a new feature, disscuss with me before modifying.
2. Write code to pass the test:
    - You will ask me before generating any function prototype without implementation
    - You or me will write the minimal implementation code to make the test pass.
    - You will build and run to confirm all tests pass.
3. Review whether we need more tests. If yes, go back to step 1.       
4. Refactor:
    - work with me to decide if we need to refactor.
    - If we do, you or me will modify the code.
    - You will then run the tests to ensure they still pass.
5. Commit:
    - You will mark the task as complete in ./FEATURES.md.
    - You will propose a **short commit message including the task ID** and, after my approval, commit the code.
6. Back to Implementation:
    - You will revisit ./IMPLEMENTATION.md to complete the step.

## Test Case Coding style
**write tests like this:**
``` go
func TestNewNote_ValidCreation_WithInjectedID(t *testing.T) {
	id := "test-id"
	title := "Test Note"
	ownerID := "owner-1"

	note, err := NewNoteWithVersion(id, title, ownerID, 0)
	if err != nil {
		t.Fatalf("Failed to create a valid note: %v", err)
	}

	if note.ID != id {
		t.Errorf("Expected ID to be '%s', but got '%s'", id, note.ID)
	}

	if note.Version != 0 {
		t.Errorf("Expected Version to be 0, but got %d", note.Version)
	}

	if note.ContentIDs == nil {
		t.Fatalf("Expected ContentIDs to be an empty slice, but it was nil")
	}

	if len(note.ContentIDs) != 0 {
		t.Errorf("Expected ContentIDs to be empty, but got %d elements", len(note.ContentIDs))
	}
}
```

**don't write test like this:**
```go
unc TestUserAccessibleNotes(t *testing.T) {
	user, _ := user.NewUser("", "testuser", "testhash")

	t.Run("should have no accessible notes initially", func(t *testing.T) {
		assert.Empty(t, user.AccessibleNoteIDs())
	})

	t.Run("should add accessible note IDs", func(t *testing.T) {
		noteID1 := "note1"
		noteID2 := "note2"

		user.AddAccessibleNoteID(noteID1)
		user.AddAccessibleNoteID(noteID2)

		expected := []string{noteID1, noteID2}
		assert.ElementsMatch(t, expected, user.AccessibleNoteIDs())
	})

	t.Run("should not add duplicate accessible note IDs", func(t *testing.T) {
		noteID1 := "note1"

		initialCount := len(user.AccessibleNoteIDs())
		user.AddAccessibleNoteID(noteID1)

		assert.Len(t, user.AccessibleNoteIDs(), initialCount)
	})

	t.Run("should remove accessible note IDs", func(t *testing.T) {
		noteIDToRemove := "note1"

		user.RemoveAccessibleNoteID(noteIDToRemove)

		assert.NotContains(t, user.AccessibleNoteIDs(), noteIDToRemove)
		assert.Len(t, user.AccessibleNoteIDs(), 1) // Only note2 should remain
	})

	t.Run("should not change anything if removing a non-existent note ID", func(t *testing.T) {
		noteIDNonExistent := "nonexistent"
		initialCount := len(user.AccessibleNoteIDs())

		user.RemoveAccessibleNoteID(noteIDNonExistent)

		assert.Len(t, user.AccessibleNoteIDs(), initialCount)
	})

	t.Run("modifying returned slice should not affect internal state", func(t *testing.T) {
		ids := user.AccessibleNoteIDs()
		ids = append(ids, "newID")

		assert.NotContains(t, user.AccessibleNoteIDs(), "newID")
	})
}
```

## Commit Message Example
write commit message as a one-liner in this style:

```bash
git commit -m "feat(backend) T7.5: Implement UserRepository and UserPO"
```

don't write commit message that are too verbal like this:

```bash
git commit -m "feat(backend): Implement UserRepository and UserPO (T15.2)\\Implement the `UserRepository` interface and `InMemoryUserRepository` for user data persistence.\\Define `UserPO` for user data.\\Add error `ErrUsernameInUse` for duplicate username scenarios.\\Add tests for `InMemoryUserRepository`.\\Update `docs/FEATURES.md` to mark T15.2 as complete and add T15.3."
```
