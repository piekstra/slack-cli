# Integration Tests

Manual integration tests for verifying slack-cli against a live Slack workspace. Tests are organized from safe (read-only) to destructive, so you can stop at any section.

---

## Part 1: Setup

### Step 1: Create Slack App

1. Go to https://api.slack.com/apps
2. Click "Create New App" → "From scratch"
3. Name it (e.g., "slack-cli-test") and select your workspace

### Step 2: Configure Bot Token Scopes

In **OAuth & Permissions**, add these scopes:

| Scope | Purpose | Required For |
|-------|---------|--------------|
| `channels:read` | List public channels, get channel info | Part 2 |
| `channels:history` | Read message history | Part 2 |
| `groups:read` | List private channels | Part 2 |
| `groups:history` | Read private channel history | Part 2 |
| `users:read` | List users, get user info | Part 2 |
| `users:read.email` | See user email addresses | Part 2 |
| `team:read` | Get workspace info | Part 2 |
| `chat:write` | Send, update, delete messages | Part 3 |
| `reactions:write` | Add/remove reactions | Part 3 |
| `channels:manage` | Create, archive, set topic/purpose, invite | Parts 4 & 5 |
| `groups:write` | Topic/purpose/invite for private channels | Part 4 |

**Note:** `channels:manage` is a superset that includes `channels:write.topic` and `channels:write.invites`. You can use the granular scopes instead if you want more limited permissions.

### Step 3: Install App & Configure CLI

```bash
# Install app to workspace (in Slack app settings)
# Copy the "Bot User OAuth Token" (starts with xoxb-)

# Build the CLI
make build

# Store the token
./bin/slack-cli config set-token
# Paste your xoxb-... token when prompted

# Verify it works
./bin/slack-cli workspace info
```

### Step 4: Discover Test Inputs

Use the CLI to find the IDs you need (Slack IDs are opaque and not easily visible in the UI):

```bash
# Find your test channel ID (bot must already be in the channel)
# Look for your test channel name in the output
slack-cli channels list

# Example output:
# ID            NAME              MEMBERS
# C08UR9H3YHU   testing           5
# C07ABC123DE   general           42

# Find a user ID (optional, for invite tests)
slack-cli users list --limit 10
```

Set these for easy reference during testing:

```bash
export TEST_CHANNEL_ID="C..."      # From channels list output
export TEST_USER_ID="U..."         # Optional: from users list output
```

| Variable | Description | Required |
|----------|-------------|----------|
| `TEST_CHANNEL_ID` | Channel the bot is already in | Yes |
| `TEST_USER_ID` | A user to invite to channels | Optional (Part 4 only) |

**Prerequisite:** The bot must already be invited to `TEST_CHANNEL_ID`. Use `/invite @your-bot-name` in Slack if needed.

---

## Part 2: Read-Only Tests

**Scopes required:** `team:read`, `users:read`, `channels:read`, `channels:history`

These tests don't modify anything. Safe to run anytime.

### 2.1 Workspace Info

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli workspace info` | Shows workspace ID, name, domain |
| 2 | `slack-cli workspace info -o json` | Valid JSON with `id`, `name`, `domain` |
| 3 | `slack-cli workspace info -o table` | Formatted table output |

### 2.2 Users

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli users list` | Table with ID, USERNAME, REAL NAME |
| 2 | `slack-cli users list --limit 3` | Exactly 3 users |
| 3 | `slack-cli users list -o json` | Valid JSON array |
| 4 | `slack-cli users get $TEST_USER_ID` | User details (ID, name, email, status) |
| 5 | `slack-cli users get $TEST_USER_ID -o json` | Full user object with nested profile |
| 6 | `slack-cli users get UINVALID999` | Error: `user_not_found` |

### 2.3 Channels

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels list` | Table with ID, NAME, MEMBERS |
| 2 | `slack-cli channels list --limit 5` | Exactly 5 channels |
| 3 | `slack-cli channels list --types public_channel` | Only public channels |
| 4 | `slack-cli channels list --exclude-archived=false` | Includes archived channels |
| 5 | `slack-cli channels list -o json` | Valid JSON array |
| 6 | `slack-cli channels get $TEST_CHANNEL_ID` | Channel details (ID, name, topic, purpose, members) |
| 7 | `slack-cli channels get $TEST_CHANNEL_ID -o json` | Full channel object |
| 8 | `slack-cli channels get CINVALID999` | Error: `channel_not_found` |

### 2.4 Message History

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli messages history $TEST_CHANNEL_ID` | Table with timestamp, user, text |
| 2 | `slack-cli messages history $TEST_CHANNEL_ID --limit 5` | Exactly 5 messages |
| 3 | `slack-cli messages history $TEST_CHANNEL_ID -o json` | Valid JSON array of messages |

### 2.5 Output Formats

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels list -o text` | Same as default (human-readable) |
| 2 | `slack-cli channels list -o json \| jq '.[0].id'` | Works with jq |
| 3 | `slack-cli channels list --no-color` | No ANSI escape codes in output |

---

## Part 3: Messaging Tests

**Scopes required:** `chat:write`, `reactions:write`

These tests create messages, then clean them up at the end.

### 3.1 Send & Verify Message

| Step | Command | Expected | Capture |
|------|---------|----------|---------|
| 1 | `slack-cli messages send $TEST_CHANNEL_ID "Integration test message"` | "Message sent (ts: X)" | **Save TS₁** |
| 2 | `slack-cli messages history $TEST_CHANNEL_ID --limit 1` | Shows your message |
| 3 | `slack-cli messages send $TEST_CHANNEL_ID "JSON test" -o json` | JSON with `ts` field | (verify only) |
| 4 | `slack-cli messages send $TEST_CHANNEL_ID "Plain text" --simple` | Message without Block Kit formatting |

### 3.2 Multiline Message (stdin)

| Step | Command | Expected |
|------|---------|----------|
| 1 | `echo -e "Line 1\nLine 2\nLine 3" \| slack-cli messages send $TEST_CHANNEL_ID -` | Multiline message appears in Slack |

### 3.3 Reactions

Using **TS₁** from step 3.1:

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli messages react $TEST_CHANNEL_ID <TS₁> thumbsup` | "Added :thumbsup: reaction" |
| 2 | `slack-cli messages react $TEST_CHANNEL_ID <TS₁> :heart:` | "Added :heart: reaction" (colons stripped) |
| 3 | `slack-cli messages react $TEST_CHANNEL_ID <TS₁> thumbsup` | Error: `already_reacted` |
| 4 | `slack-cli messages unreact $TEST_CHANNEL_ID <TS₁> thumbsup` | "Removed :thumbsup: reaction" |
| 5 | `slack-cli messages unreact $TEST_CHANNEL_ID <TS₁> heart` | "Removed :heart: reaction" |

### 3.4 Threading

Using **TS₁** from step 3.1:

| Step | Command | Expected | Capture |
|------|---------|----------|---------|
| 1 | `slack-cli messages send $TEST_CHANNEL_ID "Thread reply" --thread <TS₁>` | "Message sent" as thread reply | **Save TS₂** |
| 2 | `slack-cli messages thread $TEST_CHANNEL_ID <TS₁>` | Shows parent + reply |
| 3 | `slack-cli messages thread $TEST_CHANNEL_ID <TS₁> -o json` | JSON array of thread messages |

### 3.5 Update Message

Using **TS₁** from step 3.1:

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli messages update $TEST_CHANNEL_ID <TS₁> "Updated message text"` | "Message updated" |
| 2 | `slack-cli messages history $TEST_CHANNEL_ID --limit 1` | Shows updated text |

### 3.6 Cleanup: Delete Messages

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli messages delete $TEST_CHANNEL_ID <TS₂> --force` | "Message deleted" (thread reply) |
| 2 | `slack-cli messages delete $TEST_CHANNEL_ID <TS₁> --force` | "Message deleted" (parent) |

**Verify:** `slack-cli messages history $TEST_CHANNEL_ID --limit 3` should not show the deleted messages.

---

## Part 4: Channel Metadata Tests

**Scopes required:** `channels:write`

These tests modify channel metadata but restore original values afterward.

### 4.1 Save Original State

| Step | Command | Capture |
|------|---------|---------|
| 1 | `slack-cli channels get $TEST_CHANNEL_ID -o json` | **Save original TOPIC and PURPOSE** |

### 4.2 Modify Topic & Purpose

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels set-topic $TEST_CHANNEL_ID "Integration test topic"` | "Set topic for channel" |
| 2 | `slack-cli channels get $TEST_CHANNEL_ID` | Topic shows "Integration test topic" |
| 3 | `slack-cli channels set-purpose $TEST_CHANNEL_ID "Integration test purpose"` | "Set purpose for channel" |
| 4 | `slack-cli channels get $TEST_CHANNEL_ID` | Purpose shows "Integration test purpose" |

### 4.3 Invite User (Optional)

Skip if you didn't set `TEST_USER_ID`.

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels invite $TEST_CHANNEL_ID $TEST_USER_ID` | "Invited 1 user(s)" or "already_in_channel" |

### 4.4 Restore Original State

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels set-topic $TEST_CHANNEL_ID "<original topic>"` | Restored |
| 2 | `slack-cli channels set-purpose $TEST_CHANNEL_ID "<original purpose>"` | Restored |

---

## Part 5: Destructive Tests ⚠️

**Scopes required:** `channels:manage`

**Warning:** These tests create and archive channels. They require elevated permissions and will leave artifacts in your workspace if interrupted.

### 5.1 Create Channels

| Step | Command | Expected | Capture |
|------|---------|----------|---------|
| 1 | `slack-cli channels create test-integ-$(date +%s)` | "Created channel: test-integ-X (C...)" | **Save NEW_CHANNEL_ID** |
| 2 | `slack-cli channels create test-private-$(date +%s) --private` | "Created channel" (private) | **Save PRIVATE_CHANNEL_ID** |
| 3 | `slack-cli channels get <NEW_CHANNEL_ID>` | Shows new channel details |

### 5.2 Channel Creation Errors

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels create general` | Error: `name_taken` |
| 2 | `slack-cli channels create "has spaces"` | Error: `invalid_name_specials` |

### 5.3 Archive & Unarchive

Using **NEW_CHANNEL_ID** from step 5.1:

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels archive <NEW_CHANNEL_ID> --force` | "Archived channel" |
| 2 | `slack-cli channels archive <NEW_CHANNEL_ID>` | Error: `already_archived` |
| 3 | `slack-cli channels unarchive <NEW_CHANNEL_ID>` | "Unarchived channel" |
| 4 | `slack-cli channels unarchive <NEW_CHANNEL_ID>` | Error: `not_archived` |

### 5.4 Cleanup: Archive Test Channels

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels archive <NEW_CHANNEL_ID> --force` | "Archived channel" |
| 2 | `slack-cli channels archive <PRIVATE_CHANNEL_ID> --force` | "Archived channel" |

---

## Part 6: Config Command Tests ⚠️

**Warning:** These tests manipulate your stored token. Run last and be prepared to re-authenticate.

### 6.1 Config Show

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli config show` | Shows masked token and storage location |

### 6.2 Config Test (Optional)

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli config test` | "Token is valid" + workspace info |

### 6.3 Token Manipulation (Careful!)

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli config delete-token` | "API token deleted" |
| 2 | `slack-cli workspace info` | Error: no token configured |
| 3 | `slack-cli config set-token` | Prompts for token, stores it |
| 4 | `slack-cli workspace info` | Works again |

---

## Part 7: Error Handling & Edge Cases

These can be run at any time to verify error handling.

### 7.1 Invalid Input Errors

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli channels get INVALID` | Error with helpful hint |
| 2 | `slack-cli users get INVALID` | Error with helpful hint |
| 3 | `slack-cli messages react $TEST_CHANNEL_ID badts emoji` | Validation error |
| 4 | `slack-cli messages send` | Usage help shown |
| 5 | `slack-cli channels unknown` | Unknown command error |

### 7.2 Permission Errors

| Step | Command | Expected |
|------|---------|----------|
| 1 | `SLACK_API_TOKEN=invalid slack-cli workspace info` | Error: `invalid_auth` |
| 2 | (Use token without `channels:manage`) `slack-cli channels create test` | Error describing missing scope |

### 7.3 Edge Cases

| Step | Command | Expected |
|------|---------|----------|
| 1 | `slack-cli messages send $TEST_CHANNEL_ID "Hello 👋 世界"` | Unicode preserved |
| 2 | `slack-cli messages send $TEST_CHANNEL_ID 'Test <>&"'"'"' chars'` | Special chars escaped |
| 3 | `slack-cli channels create my-test-with-hyphens` | Works (clean up after) |

---

## Troubleshooting

| Error | Cause | Solution |
|-------|-------|----------|
| `invalid_auth` | Token invalid or expired | Regenerate token, run `config set-token` |
| `not_in_channel` | Bot not in channel | `/invite @botname` in Slack |
| `channel_not_found` | Wrong ID or no access | Verify ID, check bot permissions |
| `missing_scope` | Token lacks required scope | Add scope in Slack app settings, reinstall |
| `ratelimited` | Too many requests | Wait and retry |
| `already_archived` | Channel already archived | Use `channels unarchive` first |
| `name_taken` | Channel name exists | Choose different name |

---

## Quick Reference

```bash
# Show all commands
slack-cli --help

# Command-specific help
slack-cli channels --help
slack-cli messages send --help
```
