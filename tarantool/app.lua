local log = require("log")
local runner = require("migrations.runner")
local auth = require("modules.auth")
local user = require("modules.user")
local category = require("modules.category")
local product = require("modules.product")

runner.run({
    require("migrations.001_init_spaces"),
    require("migrations.002_seed_roles"),
    require("migrations.003_add_user_avatar"),
    require("migrations.004_add_catalog"),
})

-- Account
_G.create_user = auth.create_user
_G.verify_account = auth.verify_account
_G.set_password = auth.set_password
_G.update_password = auth.update_password
_G.get_account_by_phone_number = auth.get_account_by_phone_number

-- User 
_G.get_user_info_by_phone_number = user.get_user_info_by_phone_number
_G.update_user_info_by_phone_number = user.update_user_info_by_phone_number
_G.delete_user = user.delete_user


-- Categories
_G.create_category = category.create_category
_G.list_categories = category.list_categories
_G.get_category = category.get_category
_G.update_category = category.update_category
_G.delete_category = category.delete_category

-- Products
_G.create_product = product.create_product
_G.list_products = product.list_products
_G.get_product = product.get_product
_G.update_product = product.update_product
_G.delete_product = product.delete_product

log.info("Tarantool app loaded")
