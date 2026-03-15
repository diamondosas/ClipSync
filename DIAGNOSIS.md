# ClipSync Diagnosis & Recommendations

Based on my analysis of your codebase, here are the primary reasons why your application is not functioning as expected, along with simple explanations and steps to fix them.

## 1. Port Mismatch (The Biggest Bug)
In your code, the "Listener" and the "Connector" are looking at different doors.
- **Listener:** In `connect.go`, your `Listen` function uses port `9999` (from `globals.PORT`).
- **Connector:** In `connect.go`, your `Connect` function is hardcoded to use port `9000`.
- **Result:** Even if devices find each other, they can't talk because one is shouting at door 9000 while the other is listening at door 9999.

**Fix:** Ensure `Connect` uses `globals.PORT` instead of the hardcoded `9000`.

## 2. The "One-Shot" Discovery
In `discover.go`, your `entry` function only looks for a device **once**.
- **The Problem:** It reads one result from the channel and then stops. If it finds *your own* computer first (which it usually does), it hits the `return` statement and never looks for anyone else.
- **Result:** Your app stops searching as soon as it sees itself.

**Fix:** Change the `entry` function to use a `for ... range` loop so it keeps processing every device it finds on the network.

## 3. The "One-Shot" Receiver
In `main.go`, your receiving logic only happens once.
- **The Problem:** Inside your `main` function, the goroutine that calls `network.RecieveClipboard()` does it once, writes it to the clipboard, and then the goroutine finishes.
- **Result:** You can only receive one clipboard update per session.

**Fix:** Wrap the logic inside that goroutine in a `for` loop so it stays alive and keeps waiting for new data.

## 4. Global Variable Collision
Both `Listen()` and `Connect()` are overwriting the same global variable `network.Conn`.
- **The Problem:** If `Connect` runs after `Listen`, it replaces the listening connection. This makes the networking state very unpredictable.
- **Result:** One part of your app might "steal" the connection from another part.

**Fix:** You should separate the connection used for sending from the connection used for receiving, or use a single connection that handles both.

## 5. The Infinite Loop (Clipboard Feedback)
- **The Problem:** When you receive a clipboard update from another device, you write it to your local clipboard. Your app is *also* watching your local clipboard for changes.
- **Result:** Your app sees the change you just wrote, thinks *you* made it, and sends it back to the other device. This creates an infinite loop of messages.

**Fix:** You need a way to tell your "Watcher" to ignore changes that were just written by the "Receiver."

## 6. Instant Shutdown of Services
In `discover.go`, your `RegisterDevice` function has a `defer server.Shutdown()` at the end.
- **The Problem:** The function starts the server and then immediately reaches the end of the function. Because of the `defer`, it shuts the server down instantly.
- **Result:** Your computer tells the network "I'm here!" and then immediately says "Goodbye!" before anyone can see it.

**Fix:** The server needs to stay alive as long as your program is running. You should move the shutdown logic so it only happens when the program actually stops.

## 7. WaitGroup Counting
In `main.go`, you use `globals.WG.Add(5)`.
- **The Problem:** You only call `Done()` inside `RunWithContext`, but you only use `RunWithContext` three times. The other two background tasks never call `Done()`.
- **Result:** Your program will hang forever at the end (`WG.Wait()`) and never shut down cleanly because it's waiting for 2 tasks that already finished (or never started correctly).

## 8. Premature Connection Closing
In both `Connect` and `Listen`, you have `defer Conn.Close()`.
- **The Problem:** Because these functions run and finish very quickly (they don't wait for messages), the `defer` command closes your connection right after it is opened.
- **Result:** By the time you try to send or receive data, the connection is already closed.

---

### Suggested Next Steps
1. **Unify Ports:** Make sure every part of the app uses `globals.PORT`.
2. **Loop Everything:** Use `for` loops in `entry` (discovery) and your receiver in `main.go`.
3. **Keep Services and Connections Alive:** Remove `defer Conn.Close()` or `defer server.Shutdown()` from functions that finish quickly. These should stay open until the entire app closes.
4. **Fix the Counters:** Make sure your `Add()` calls match your `Done()` calls exactly.
5. **Debug Logging:** Add more `log.Println` statements to see exactly where the code stops (e.g., "Checking discovered device...", "Actually connecting to...").
