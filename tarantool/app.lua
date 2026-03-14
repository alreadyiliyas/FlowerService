local log = require("log")
local runner = require("migrations.runner")
local auth = require("modules.auth")
local user = require("modules.user")

runner.run({
    require("migrations.001_init_spaces"),
    require("migrations.002_seed_roles"),
})

-- Account
_G.create_user = auth.create_user
_G.verify_account = auth.verify_account
_G.set_password = auth.set_password
_G.update_password = auth.update_password
_G.get_account_by_phone_number = auth.get_account_by_phone_number

-- User 
_G.get_user_info_by_phone_number = user.get_user_info_by_phone_number

log.info("Tarantool app loaded")

