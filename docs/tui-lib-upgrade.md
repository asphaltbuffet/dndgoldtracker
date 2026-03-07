# TUI Library Upgrade: Charmbracelet v1 → v2

## Technical Summary

### Module Path Changes (`go.mod`, `go.sum`, `gomod2nix.toml`)

The three core Charmbracelet libraries were replaced with their v2 major
versions. The v2 modules declare a new canonical module path under `charm.land`
rather than `github.com/charmbracelet`, so both the module require directives
and every import path in source files changed:

| Before | After |
|---|---|
| `github.com/charmbracelet/bubbletea` | `charm.land/bubbletea/v2` |
| `github.com/charmbracelet/bubbles` | `charm.land/bubbles/v2` |
| `github.com/charmbracelet/lipgloss` | `charm.land/lipgloss/v2` |

The Go toolchain version was bumped from `1.23.5` to `1.25.6`.

The indirect dependency graph changed substantially:

**Removed**:
- `go-osc52`
- `erikgeiser/coninput`
- `mattn/go-isatty`
- `mattn/go-localereader`
- `muesli/ansi`
- `muesli/termenv`
- `golang.org/x/text`

**Added**:
- `charmbracelet/colorprofile`
- `charmbracelet/ultraviolet`
- `charmbracelet/x/termios`
- `charmbracelet/x/windows`
- `clipperhouse/displaywidth`
- `clipperhouse/uax29/v2`
- `xo/terminfo`

---

### BubbleTea v2 API Changes (`ui/ui.go`, `ui/updates.go`)

**Key message type renamed:**

All five `Update` handlers switched from `tea.KeyMsg` to `tea.KeyPressMsg`. BubbleTea v2
distinguishes press and release events; `tea.KeyPressMsg` is the direct replacement for the v1
`tea.KeyMsg`. All `msg.String()` comparisons (`"ctrl+c"`, `"tab"`, `"enter"`, etc.) are unchanged.

**`View()` return type changed:**

`View()` now returns `tea.View` instead of `string`. All return sites wrap their string output in
`tea.NewView(...)`. Sub-view functions (`choicesView`, `moneyView`, etc.) still return `string` and
are composed inside the top-level `View()`.

---

### Bubbles v2 API Changes (`ui/utils.go`, `ui/ui.go`)

**`cursor` subpackage eliminated:**

The `github.com/charmbracelet/bubbles/cursor` import is gone entirely. The `cursor.Mode` enum
(three states: `CursorBlink`, `CursorStatic`, `CursorHide`) is replaced by a single `bool` on
`textinput.Model`: `VirtualCursor` (`true` = software blinking cursor, `false` = real terminal
cursor). Consequently:

- `model.cursorMode cursor.Mode` → `model.virtualCursor bool` (initialized to `true`)
- `changeCursorMode` went from returning `[]tea.Cmd` (one per input, to propagate the mode change
  asynchronously) to returning nothing — `SetVirtualCursor(bool)` is synchronous and no commands
  are needed
- The three-state cycle (`ctrl+r` previously stepped through blink → static → hide) is now a
  simple toggle between virtual and real
- The `cursorStyle` package-level variable (which applied the focused lipgloss style to the cursor
  block) was removed; `CursorStyle` in v2 takes `image/color.Color`, not a lipgloss style

**TextInput style API changed from direct field assignment to `SetStyles`:**

Styles are no longer set via exported struct fields. The pattern is now:

```go
// Before
t.PromptStyle = focusedStyle
t.TextStyle = focusedStyle

// After
s := textinput.DefaultStyles(false)
s.Focused.Prompt = focusedStyle
s.Focused.Text = focusedStyle
t.SetStyles(s)
```

`updateFocusIndex` uses `inputs[i].Styles()` to read current styles before mutating and writing
back via `SetStyles(s)`, since the styles struct is now owned by the model rather than being
separately addressable fields.

---

## Architectural Impact

**Cursor state simplified.** The previous design stored a three-state enum in the model and sent
batched async commands to propagate cursor mode changes to all inputs. The v2 design is a
synchronous bool toggle — `changeCursorMode` became a void function with no return value, and its
callers in `updateMoney`, `updateExperience`, and `updateAddMember` now return `nil` instead of
`tea.Batch(cmds...)`. This reduces async plumbing.

**TextInput style encapsulation tightened.** Styles are now encapsulated inside the textinput
model rather than being exposed as mutable exported fields. The read-modify-write pattern
(`Styles()` → mutate → `SetStyles()`) is the only supported interface, which prevents accidental
partial updates and aligns with the immutable-update philosophy of the BubbleTea architecture.

**`View()` return type opens declarative rendering.** Returning `tea.View` (a struct) instead of
`string` allows the runtime to read metadata fields (`AltScreen`, `MouseMode`, etc.) from the view
declaratively rather than requiring program options or imperative commands. This codebase does not
yet use any of those fields, but the door is open.

**Dependency graph modernized.** Several v1-era compatibility shims were removed
(`muesli/termenv`, `go-osc52`, `coninput`, `mattn/go-localereader`). Their functionality is now
provided by the newer Charmbracelet platform packages (`charmbracelet/ultraviolet`,
`charmbracelet/colorprofile`, `x/termios`, `x/windows`), which form the new low-level I/O and
color-profile foundation.

---

## Next Steps

1. **Restore cursor color styling.** The `cursorStyle = focusedStyle` variable was removed because
   `textinput.CursorStyle.Color` is `image/color.Color`, not a lipgloss style. If the pink cursor
   (ANSI 205) is still desired, convert `lipgloss.Color("205")` to an `image/color.Color` and set
   `s.Cursor.Color` in `configureInputs`.

2. **Review `ctrl+r` behavior change.** The cursor toggle went from a three-state cycle
   (blink → static → hide) to a two-state toggle (virtual ↔ real). The help text still reads
   "cursor mode is virtual/real (ctrl+r to change style)" — verify this matches the intended UX, or
   simplify/remove the hint if the distinction is no longer meaningful.

3. **Fix `updateInputs` mutation bug** (pre-existing, now more visible). `ui/utils.go` updates
   `inputs` slice elements in-place but the caller never propagates the updated slice back to the
   model. In v2's immutable-update pattern this is a latent bug: refactor `updateInputs` to return
   `([]textinput.Model, tea.Cmd)` and assign the result back to the appropriate model field
   (`m.coinInputs`, `m.xpInputs`, `m.memberInputs`).

4. **Consider adopting `tea.View` metadata fields.** Now that `View()` returns `tea.View`, fields
   like `v.AltScreen = true` or `v.MouseMode = tea.MouseModeCellMotion` can be set declaratively
   instead of being passed as `tea.NewProgram` options or sent as commands. Evaluate whether any
   current or future program options should move into `View()`.

5. **Update Nix derivation** for any NixOS/home-manager packaging. `gomod2nix.toml` has been
   regenerated with v2 hashes, but the Nix flake's `buildGoModule` or `gomod2nix` call may need
   the `go` attribute bumped to match the new `go 1.24.2` toolchain requirement.
