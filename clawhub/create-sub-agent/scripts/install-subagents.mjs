#!/usr/bin/env node

import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { execFileSync } from "node:child_process";

const defaultConfigPath = path.join(os.homedir(), ".openclaw", "openclaw.json");

function parseArgs(argv) {
  const args = {
    configPath: defaultConfigPath,
    mainAgentId: "main",
    workspaceRoot: ".",
    roles: [],
  };

  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    if (arg === "--config" && argv[i + 1]) {
      args.configPath = argv[i + 1];
      i += 1;
      continue;
    }
    if (arg === "--main-agent" && argv[i + 1]) {
      args.mainAgentId = argv[i + 1];
      i += 1;
      continue;
    }
    if (arg === "--workspace-root" && argv[i + 1]) {
      args.workspaceRoot = argv[i + 1];
      i += 1;
      continue;
    }
    if (arg === "--roles" && argv[i + 1]) {
      args.roles = parseRoles(argv[i + 1]);
      i += 1;
      continue;
    }
    if (arg === "--help" || arg === "-h") {
      printHelp();
      process.exit(0);
    }
  }

  return args;
}

function parseRoles(raw) {
  if (!raw?.trim()) {
    throw new Error(
      "--roles is required. Pass a comma-separated role list such as pm,architect,reviewer or id:name entries such as pm:Product Manager,reviewer:Code Reviewer.",
    );
  }

  const entries = raw
    .split(",")
    .map((value) => value.trim())
    .filter(Boolean);
  if (entries.length === 0) {
    throw new Error("--roles is required. Pass at least one role.");
  }

  return entries.map((entry) => {
    const [rawId, ...nameParts] = entry.split(":");
    const id = rawId?.trim();
    if (!id) {
      throw new Error(`Invalid role entry "${entry}".`);
    }
    const name = nameParts.join(":").trim() || id;
    return { id, name };
  });
}

function printHelp() {
  console.log(`Usage:
  node .agents/skills/create-sub-agent/scripts/install-subagents.mjs [--config <path>] [--main-agent <id>] [--workspace-root <dir>] [--roles <ids>]

Options:
  --config          Path to openclaw.json (default: ~/.openclaw/openclaw.json)
  --main-agent      Main agent id that will be allowed to spawn workflow sub-agents
  --workspace-root  Root directory where workspace-* folders will be created
  --roles           Required. Comma-separated role ids or id:name entries to install
`);
}

function readJson(filePath) {
  if (!fs.existsSync(filePath)) {
    throw new Error(`Config file not found: ${filePath}`);
  }
  const raw = fs.readFileSync(filePath, "utf8");
  return JSON.parse(raw);
}

function writeJson(filePath, value) {
  fs.writeFileSync(filePath, `${JSON.stringify(value, null, 2)}\n`, "utf8");
}

function findAgentIndex(list, agentId) {
  return list.findIndex((entry) => entry?.id === agentId);
}

function ensureDir(dirPath) {
  fs.mkdirSync(dirPath, { recursive: true });
}

const IDENTITY_BLOCK_START = "<!-- OPENCLAW:ROLE-IDENTITY:START -->";
const IDENTITY_BLOCK_END = "<!-- OPENCLAW:ROLE-IDENTITY:END -->";

function buildManagedIdentityBlock(role) {
  return `${IDENTITY_BLOCK_START}
# ${role.name}

Role id: ${role.id}

This workspace belongs to the ${role.name} agent.
Focus on the responsibilities implied by the role name and the task received from the main agent.
${IDENTITY_BLOCK_END}`;
}

function ensureIdentityFile(workspace, role) {
  const identityPath = path.join(workspace, "IDENTITY.md");
  const managedBlock = buildManagedIdentityBlock(role);

  if (!fs.existsSync(identityPath)) {
    fs.writeFileSync(identityPath, `${managedBlock}\n`, "utf8");
    return;
  }

  const existing = fs.readFileSync(identityPath, "utf8");
  const start = existing.indexOf(IDENTITY_BLOCK_START);
  const end = existing.indexOf(IDENTITY_BLOCK_END);

  if (start !== -1 && end !== -1 && end > start) {
    const before = existing.slice(0, start).replace(/\s*$/, "");
    const after = existing.slice(end + IDENTITY_BLOCK_END.length).replace(/^\s*/, "");
    const parts = [before, managedBlock, after].filter(Boolean);
    fs.writeFileSync(identityPath, `${parts.join("\n\n")}\n`, "utf8");
    return;
  }

  const separator = existing.trim().length === 0 ? "" : "\n\n";
  fs.writeFileSync(identityPath, `${existing.replace(/\s*$/, "")}${separator}${managedBlock}\n`, "utf8");
}

function resolveWorkspacePath(workspace, configDir, fallbackPath) {
  if (!workspace) {
    return fallbackPath;
  }
  return path.isAbsolute(workspace) ? workspace : path.resolve(configDir, workspace);
}

function copyMainSkillsToWorkspace(mainWorkspace, roleWorkspace) {
  const sourceSkillsDir = path.join(mainWorkspace, "skills");
  if (!fs.existsSync(sourceSkillsDir)) {
    return;
  }

  const targetSkillsDir = path.join(roleWorkspace, "skills");
  fs.mkdirSync(targetSkillsDir, { recursive: true });
  fs.cpSync(sourceSkillsDir, targetSkillsDir, {
    recursive: true,
    force: false,
    errorOnExist: false,
  });
}

function ensureMainAgent(agents, mainAgentId, workspaceRoot) {
  const nextAgents = [...agents];
  let mainIndex = findAgentIndex(nextAgents, mainAgentId);

  if (mainIndex === -1) {
    throw new Error(
      `Main agent "${mainAgentId}" not found in config. Create/configure the main agent first before installing sub-agents.`,
    );
  }

  return { agents: nextAgents, mainIndex };
}

function ensureAgentExists(role, workspace, cwd) {
  execFileSync(
    "openclaw",
    ["agents", "add", role.id, "--workspace", workspace, "--non-interactive"],
    {
      cwd,
      stdio: "inherit",
    },
  );
}

function main() {
  const args = parseArgs(process.argv.slice(2));
  if (args.roles.length === 0) {
    throw new Error("--roles is required.");
  }

  const configPath = path.resolve(args.configPath);
  const configDir = path.dirname(configPath);
  const workspaceRoot = path.resolve(args.workspaceRoot);
  const config = readJson(configPath);

  config.agents ??= {};
  config.agents.list ??= [];

  const ensured = ensureMainAgent(config.agents.list, args.mainAgentId, workspaceRoot);
  let agents = ensured.agents;
  const mainAgent = agents[ensured.mainIndex];
  const mainWorkspace = resolveWorkspacePath(
    mainAgent.workspace,
    configDir,
    path.join(workspaceRoot, `workspace-${args.mainAgentId}`),
  );

  mainAgent.subagents ??= {};
  const allowAgents = new Set(mainAgent.subagents.allowAgents ?? []);

  for (const role of args.roles) {
    const workspace = path.join(workspaceRoot, `workspace-${role.id}`);
    const roleIndex = findAgentIndex(agents, role.id);

    if (roleIndex === -1) {
      ensureAgentExists(role, workspace, workspaceRoot);
      const refreshedConfig = readJson(configPath);
      refreshedConfig.agents ??= {};
      refreshedConfig.agents.list ??= [];
      agents = refreshedConfig.agents.list;
    } else if (!agents[roleIndex]?.workspace) {
      agents[roleIndex] = { ...agents[roleIndex], workspace };
    }

    allowAgents.add(role.id);
    ensureDir(workspace);
    ensureIdentityFile(workspace, role);
    copyMainSkillsToWorkspace(mainWorkspace, workspace);
  }

  const mainAgentIndex = findAgentIndex(agents, args.mainAgentId);
  const nextMainAgent = agents[mainAgentIndex];
  nextMainAgent.subagents ??= {};
  nextMainAgent.subagents.allowAgents = [...allowAgents].sort();
  config.agents.list = agents;
  config.agents.list[mainAgentIndex] = nextMainAgent;

  writeJson(configPath, config);

  console.log(`Installed workflow sub-agents into ${configPath}`);
  console.log(`Main agent: ${args.mainAgentId}`);
  console.log(`Installed roles: ${args.roles.map((role) => role.id).join(", ")}`);
  console.log(`Allowed agents: ${mainAgent.subagents.allowAgents.join(", ")}`);
}

main();
