
# Error Propagation Principles in Go

## 1. What Happened?
This part of the error contains information about **what happened**, e.g.,
- "Disk full"
- "Socket closed"
- "Credentials expired"

This information is likely generated implicitly by whatever caused the error, but it can be decorated with additional context to help the user understand it better.

## 2. When and Where It Occurred?
Errors should always contain:
- A **complete stack trace**, starting from how the call was initiated to where the error was instantiated.
- Stack traces should **not** be included in the error message but should be **easily accessible** when handling the error up the stack.
- Additional **contextual information**, such as:
  - Machine ID (in a distributed system, this helps identify where the error occurred).
  - **Timestamp** in UTC, showing when the error was instantiated.

## 3. A Friendly User-Facing Message
The error message displayed to the user should:
- Be customized for the system and its users.
- Contain **only relevant information** from the previous points.
- Be **human-centric** and easy to understand.
- Indicate whether the issue is **transitory** or permanent.
- Be concise (ideally one line of text).

## 4. How the User Can Get More Information
At some point, someone will need detailed error information. Errors should provide:
- An **Error ID**, which can be cross-referenced with logs for full details.
- Logs containing:
  - The **exact time** the error occurred (not when it was logged).
  - The **stack trace**.
  - Any **additional context** from when the error was created.
- Optionally, a **hash of the stack trace** can be included to group similar issues in bug trackers.

