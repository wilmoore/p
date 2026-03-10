import { execFile } from "node:child_process";
import { promises as fs } from "node:fs";
import path from "node:path";
import process from "node:process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

const REPO_ROOT = process.cwd();
const HANDOFF_DIR = path.join(REPO_ROOT, "doc", ".plan", "session-handoff");
const INDEX_PATH = path.join(HANDOFF_DIR, "index.json");
const LEDGER_PATH = path.join(REPO_ROOT, "doc", ".plan", "session-handoff.md");

function isoNow() {
  return new Date().toISOString();
}

function die(message, exitCode = 1) {
  process.stderr.write(`${message}\n`);
  process.exit(exitCode);
}

async function ensureDirs() {
  await fs.mkdir(path.join(HANDOFF_DIR, "sessions"), { recursive: true });
  await fs.mkdir(path.join(HANDOFF_DIR, "archive"), { recursive: true });
}

async function loadIndex() {
  try {
    const raw = await fs.readFile(INDEX_PATH, "utf8");
    const parsed = JSON.parse(raw);
    if (!parsed || typeof parsed !== "object") {
      throw new Error("index.json is not an object");
    }
    if (!Array.isArray(parsed.sessions)) {
      parsed.sessions = [];
    }
    if (typeof parsed.version !== "number") {
      parsed.version = 1;
    }
    return parsed;
  } catch (err) {
    if (err && err.code === "ENOENT") {
      return { version: 1, sessions: [] };
    }
    throw err;
  }
}

async function saveIndex(index) {
  const json = JSON.stringify(index, null, 2) + "\n";
  await fs.writeFile(INDEX_PATH, json, "utf8");
}

function relFromRepo(p) {
  return path.relative(REPO_ROOT, p).split(path.sep).join("/");
}

function ledgerMarkdown(index) {
  const updated = isoNow();
  const pending = index.sessions.filter((s) => s.status === "pending");
  const current = pending.length > 0 ? pending[0].id : "";

  const lines = [];
  lines.push("# Session Handoff Ledger");
  lines.push("");
  lines.push(`Updated: ${updated}`);
  lines.push(`Current session: ${current}`);
  lines.push("");
  lines.push(`## Outstanding Snapshots (${pending.length})`);
  lines.push("");

  if (pending.length === 0) {
    lines.push("- None");
  } else {
    pending.forEach((s, i) => {
      const rel = s.relativePath ? `doc/.plan/session-handoff/${s.relativePath}` : "";
      lines.push(`${i + 1}. [pending] ${s.id} — ${s.branch || "(unknown branch)"} (${s.workingTree || "unknown"})`);
      if (rel) lines.push(`   File: ${rel}`);
      if (s.updatedAt) lines.push(`   Updated: ${s.updatedAt}`);
    });
  }

  lines.push("");
  lines.push("## Recent Activity");
  lines.push("");

  const recent = index.sessions
    .filter((s) => s.status !== "pending")
    .slice()
    .sort((a, b) => (b.updatedAt || "").localeCompare(a.updatedAt || ""))
    .slice(0, 5);

  if (recent.length === 0) {
    lines.push("- None");
  } else {
    recent.forEach((s) => {
      const note = s.note ? ` — ${s.note}` : "";
      const reason = s.reason ? ` — ${s.reason}` : "";
      lines.push(`- ${s.status}: ${s.id}${note}${reason}`);
    });
  }

  lines.push("");
  lines.push("## Commands");
  lines.push("");
  lines.push("- `node bin/session-handoff.mjs list` — show pending snapshots");
  lines.push("- `node bin/session-handoff.mjs ack <id> [--note \"done\"]` — mark complete");
  lines.push("- `node bin/session-handoff.mjs dismiss <id> --reason \"why\"` — abandon work");
  lines.push("- `node bin/session-handoff.mjs write --trigger \"/pro:session.handoff\"` — capture a fresh snapshot");
  lines.push("");
  lines.push("All snapshots live under `doc/.plan/session-handoff/sessions/`. Review each file before acknowledging or dismissing it.");
  lines.push("");
  return lines.join("\n");
}

async function writeLedger(index) {
  await fs.writeFile(LEDGER_PATH, ledgerMarkdown(index), "utf8");
}

function parseFlags(args) {
  const out = { _: [] };
  for (let i = 0; i < args.length; i++) {
    const a = args[i];
    if (!a.startsWith("--")) {
      out._.push(a);
      continue;
    }
    const eq = a.indexOf("=");
    if (eq !== -1) {
      out[a.slice(2, eq)] = a.slice(eq + 1);
      continue;
    }
    const key = a.slice(2);
    const next = args[i + 1];
    if (next && !next.startsWith("--")) {
      out[key] = next;
      i++;
    } else {
      out[key] = true;
    }
  }
  return out;
}

async function updateSnapshotFile(absPath, fields) {
  let body;
  try {
    body = await fs.readFile(absPath, "utf8");
  } catch (err) {
    if (err && err.code === "ENOENT") return;
    throw err;
  }

  const replaceLine = (label, value) => {
    const re = new RegExp(`^(${label}:)\\s.*$`, "m");
    if (re.test(body)) {
      body = body.replace(re, `$1 ${value}`);
    }
  };

  if (fields.status) replaceLine("Status", fields.status);
  if (fields.updatedAt) replaceLine("Updated", fields.updatedAt);

  if (fields.note) {
    if (!body.includes("\n\n## Notes\n")) {
      body = body.trimEnd() + "\n\n## Notes\n\n";
    }
    body = body.trimEnd() + `\n- ${fields.note}\n`;
  }

  await fs.writeFile(absPath, body, "utf8");
}

async function moveToArchiveIfPresent(session) {
  if (!session.relativePath) return session;

  const abs = path.join(HANDOFF_DIR, session.relativePath);
  const filename = path.basename(abs);
  const archiveRel = path.join("archive", filename).split(path.sep).join("/");
  const archiveAbs = path.join(HANDOFF_DIR, "archive", filename);

  try {
    await fs.rename(abs, archiveAbs);
  } catch (err) {
    if (err && err.code === "ENOENT") {
      // If the file isn't present, just leave paths alone.
      return session;
    }
    throw err;
  }

  return { ...session, relativePath: archiveRel };
}

async function listCmd() {
  await ensureDirs();
  const index = await loadIndex();
  const pending = index.sessions.filter((s) => s.status === "pending");
  if (pending.length === 0) {
    process.stdout.write("No pending session handoff snapshots.\n");
    return;
  }

  pending.forEach((s) => {
    const rel = s.relativePath ? `doc/.plan/session-handoff/${s.relativePath}` : "";
    process.stdout.write(`${s.id}\t${s.branch || ""}\t${rel}\n`);
  });
}

async function ackCmd(id, note) {
  if (!id) die("usage: node bin/session-handoff.mjs ack <id> [--note \"done\"]");

  await ensureDirs();
  const index = await loadIndex();
  const i = index.sessions.findIndex((s) => s.id === id);
  if (i === -1) die(`unknown session id: ${id}`);

  const updatedAt = isoNow();
  let session = { ...index.sessions[i] };
  session.status = "acked";
  session.updatedAt = updatedAt;
  if (note) session.note = note;

  // Update snapshot file before moving.
  if (session.relativePath) {
    await updateSnapshotFile(path.join(HANDOFF_DIR, session.relativePath), {
      status: "acked",
      updatedAt,
      note: note ? `Ack: ${note}` : "Acked",
    });
  }

  session = await moveToArchiveIfPresent(session);
  index.sessions[i] = session;
  await saveIndex(index);
  await writeLedger(index);

  process.stdout.write(`Acked ${id}\n`);
}

async function dismissCmd(id, reason) {
  if (!id) die("usage: node bin/session-handoff.mjs dismiss <id> --reason \"why\"");
  if (!reason || reason === true) die("dismiss requires --reason");

  await ensureDirs();
  const index = await loadIndex();
  const i = index.sessions.findIndex((s) => s.id === id);
  if (i === -1) die(`unknown session id: ${id}`);

  const updatedAt = isoNow();
  let session = { ...index.sessions[i] };
  session.status = "dismissed";
  session.updatedAt = updatedAt;
  session.reason = reason;

  if (session.relativePath) {
    await updateSnapshotFile(path.join(HANDOFF_DIR, session.relativePath), {
      status: "dismissed",
      updatedAt,
      note: `Dismissed: ${reason}`,
    });
  }

  session = await moveToArchiveIfPresent(session);
  index.sessions[i] = session;
  await saveIndex(index);
  await writeLedger(index);

  process.stdout.write(`Dismissed ${id}\n`);
}

async function gitSafe(args) {
  try {
    const { stdout } = await execFileAsync("git", args, { cwd: REPO_ROOT });
    return stdout.trim();
  } catch {
    return "";
  }
}

function makeSessionId() {
  // Keep same shape as existing ids.
  const ts = new Date().toISOString().replace(/[:.]/g, "-");
  const rand = Math.random().toString(16).slice(2, 10);
  return `session-${ts}-${rand}`;
}

async function writeCmd(trigger) {
  await ensureDirs();

  const index = await loadIndex();
  const id = makeSessionId();
  const now = isoNow();
  const branch = await gitSafe(["branch", "--show-current"]);
  const statusPorcelain = await gitSafe(["status", "--porcelain=v1"]);
  const workingTree = statusPorcelain ? "dirty" : "clean";
  const recentCommit = await gitSafe(["log", "-1", "--oneline"]);

  const filename = `${id}.md`;
  const relativePath = `sessions/${filename}`;
  const absPath = path.join(HANDOFF_DIR, "sessions", filename);

  const snapshot = [
    "# Session Handoff Snapshot",
    "",
    `ID: ${id}`,
    "Status: pending",
    `Created: ${now}`,
    `Updated: ${now}`,
    `Trigger: ${trigger || "manual"}`,
    "",
    "## Current State",
    "",
    `- Branch: \`${branch || "(unknown)"}\``,
    `- Working tree: ${workingTree}`,
    `- Last commit: \`${recentCommit || "(unknown)"}\``,
    "- In-progress items: unknown",
    "",
    "## Next Steps",
    "",
    "1. Review the outstanding checklist stored in this file.",
    "2. Once complete, acknowledge the snapshot:",
    "   ```bash",
    `   node bin/session-handoff.mjs ack ${id}`,
    "   ```",
    "3. If the work is obsolete, dismiss it instead:",
    "   ```bash",
    `   node bin/session-handoff.mjs dismiss ${id} --reason \"why\"`,
    "   ```",
    "",
  ].join("\n");

  await fs.writeFile(absPath, snapshot, "utf8");
  index.sessions.push({
    id,
    filename,
    relativePath,
    status: "pending",
    createdAt: now,
    updatedAt: now,
    trigger: trigger || "manual",
    branch,
    workingTree,
    recentCommit,
    backlogSummary: "- In-progress items: unknown",
  });
  await saveIndex(index);
  await writeLedger(index);

  process.stdout.write(`${id}\n`);
}

async function main() {
  const flags = parseFlags(process.argv.slice(2));
  const cmd = flags._[0];
  const id = flags._[1];

  try {
    if (cmd === "list") {
      await listCmd();
      return;
    }
    if (cmd === "ack") {
      await ackCmd(id, flags.note && flags.note !== true ? flags.note : "");
      return;
    }
    if (cmd === "dismiss") {
      await dismissCmd(id, flags.reason);
      return;
    }
    if (cmd === "write") {
      await writeCmd(flags.trigger && flags.trigger !== true ? flags.trigger : "");
      return;
    }

    die(
      [
        "usage:",
        "  node bin/session-handoff.mjs list",
        "  node bin/session-handoff.mjs ack <id> [--note \"done\"]",
        "  node bin/session-handoff.mjs dismiss <id> --reason \"why\"",
        "  node bin/session-handoff.mjs write --trigger \"/pro:session.handoff\"",
      ].join("\n"),
      2,
    );
  } catch (err) {
    die(err && err.stack ? err.stack : String(err));
  }
}

await main();
