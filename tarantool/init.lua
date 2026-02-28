local log = require("log")
local runner = require("runner")
local create_user_impl = require("create_user")

box.cfg{ listen = 3301, memtx_memory = 256 * 1024 * 1024, log_level = 5 }

box.schema.user.create("app", { password = "secret", if_not_exists = true })
box.schema.user.grant("app", "read,write,execute", "universe", nil, { if_not_exists = true })

runner.run({
    require("001_init_spaces"),
    require("002_seed_roles"),
})

_G.create_user = create_user_impl
log.info("Tarantool started")
