#!/usr/bin/env node

import fs from "node:fs";
import path from "node:path";

function parseArgs(argv) {
  const args = {
    configPath: "./openclaw.json",
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
  --config          Path to openclaw.json
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

function ensureIdentityFile(workspace, role) {
  const identityPath = path.join(workspace, "IDENTITY.md");
  if (fs.existsSync(identityPath)) {
    return;
  }
  const content = `# ${role.name}

Role id: ${role.id}

This workspace belongs to the ${role.name} agent.
Focus on the responsibilities implied by the role name and the task received from the main agent.
`;
  fs.writeFileSync(identityPath, content, "utf8");
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
    const mainAgent = {
      id: mainAgentId,
      workspace: path.join(workspaceRoot, `workspace-${mainAgentId}`),
      subagents: { allowAgents: [] },
    };
    if (mainAgentId === "main") {
      mainAgent.default = true;
    }
    nextAgents.push(mainAgent);
    mainIndex = nextAgents.length - 1;
  }

  return { agents: nextAgents, mainIndex };
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
  const agents = ensured.agents;
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
      agents.push({
        id: role.id,
        workspace,
      });
    } else if (!agents[roleIndex]?.workspace) {
      agents[roleIndex] = {
        ...agents[roleIndex],
        workspace,
      };
    }

    allowAgents.add(role.id);
    ensureDir(workspace);
    ensureIdentityFile(workspace, role);
    copyMainSkillsToWorkspace(mainWorkspace, workspace);
  }

  mainAgent.subagents.allowAgents = [...allowAgents].sort();
  config.agents.list = agents;

  writeJson(configPath, config);

  console.log(`Installed workflow sub-agents into ${configPath}`);
  console.log(`Main agent: ${args.mainAgentId}`);
  console.log(`Installed roles: ${args.roles.map((role) => role.id).join(", ")}`);
  console.log(`Allowed agents: ${mainAgent.subagents.allowAgents.join(", ")}`);
}

main();
