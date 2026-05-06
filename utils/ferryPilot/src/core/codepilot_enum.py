FILE_MAP_PATH = "../config/file_map.json"
FILE_MAP_FROZEN_PATH = "config/file_map.json"  # PyInstaller 打包后 _MEIPASS 下的路径
RemoteKey = "remote"
SourceSectionKey = "source"
TargetSectionKey = "target"
PathKey = "path"
ExcludeKey = "exclude"

type ConfigType = str
ConfigGlobal: ConfigType = "global"
ConfigProject: ConfigType = "project"
ConfigTypeList: list[ConfigType] = [ConfigProject, ConfigGlobal]

type TargetAgent = str
SourceAgent: TargetAgent = "source"
CodexAgent: TargetAgent = "codex"
GeminiAgent: TargetAgent = "gemini"
CopilotAgent: TargetAgent = "copilot"
CursorAgent: TargetAgent = "cursor"
ClaudeAgent: TargetAgent = "claude"
TargetAgentList: list[TargetAgent] = [CodexAgent, GeminiAgent, CopilotAgent, CursorAgent, ClaudeAgent]

type DirKey = str
SubAgentsKey: DirKey = "sub_agents"
