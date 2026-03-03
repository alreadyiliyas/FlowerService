local log = require("log")
local runner = require("migrations.runner")
local auth = require("modules.auth")

runner.run({
    require("migrations.001_init_spaces"),
    require("migrations.002_seed_roles"),
})

_G.create_user = auth.create_user
_G.verify_account = auth.verify_account
_G.set_password = auth.set_password
_G.update_password = auth.update_password

log.info("Tarantool app loaded")

