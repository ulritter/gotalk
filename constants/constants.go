package constants

const RAWFILE = "https://raw.githubusercontent.com/ulritter/gotalk-app/main/language.json"
const LANGFILE = "../language.json"

// actions for client <-> server communication
const ACTION_CHANGENICK = "changenick"
const ACTION_SENDMESSAGE = "message"
const ACTION_LISTUSERS = "listusers"
const ACTION_REVISION = "revision"
const ACTION_SENDSTATUS = "status"
const ACTION_EXIT = "exit"
const ACTION_INIT = "init"

// end user commands on ui
const CMD_PREFIX = '/'
const CMD_EXIT1 = "exit"
const CMD_EXIT2 = "quit"
const CMD_EXIT3 = "q"
const CMD_CHANGENICK = "nick"
const CMD_LISTUSERS = "list"
const CMD_HELP = "help"
const CMD_HELP1 = "?"

const BUFSIZE = 4096

const REVISION = "0.8.3"
