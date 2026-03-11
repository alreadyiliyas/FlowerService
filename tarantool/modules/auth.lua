local E = require("errors.errors")

local M = {}

function M.create_user(first_name, last_name, phone_number, role_name)
    role_name = role_name or "user"
    local now = os.time()

    return box.atomic(function()
        if box.space.auths.index.phone_number:get{phone_number} then
            error(E.PHONE_ALREADY_EXISTS)
        end

        local role = box.space.roles.index.name:get{role_name}
        if role == nil then
            error(E.ROLE_NOT_FOUND)
        end

        local user = box.space.users:auto_increment{
            first_name, last_name, role.id, 1, false, now, now, box.NULL
        }

        box.space.auths:auto_increment{
            user.id, phone_number, box.NULL, "phone", box.NULL, box.NULL, now, box.NULL
        }

        return {user.id, user.first_name, user.last_name, role.name, now}
    end)
end

function M.verify_account(phone_number)
    return box.atomic(function()
        local account = box.space.auths.index.phone_number:get{phone_number}
        if account == nil then
            error(E.ACCOUNT_NOT_FOUND)
        end

        local user = box.space.users:get{account.user_id}
        if user == nil then
            error(E.USER_NOT_FOUND)
        end

        if user.is_active then
            error(E.ALREADY_ACTIVE)
        end

        local now = os.time()
        box.space.users:update(user.id, {
            {"=", 6, true},
            {"=", 8, now},
            {"=", 5, user.version + 1},
        })

        return true
    end)
end

function M.set_password(phone_number, password_hash)
    return box.atomic(function()
        local account = box.space.auths.index.phone_number:get{phone_number}
        if account == nil then
            error(E.ACCOUNT_NOT_FOUND)
        end

        local user = box.space.users:get{account.user_id}
        if user == nil then
            error(E.USER_NOT_FOUND)
        end

        if user.is_active == false then
            error(E.ALREADY_NOT_ACTIVE)
        end

        if account.password_hash ~= nil then
            error(E.USER_SET_PASSWORD)
        end

        local now = os.time()
        box.space.auths:update(account.id, {
            {"=", 7, password_hash},
            {"=", 8, now},
        })

        return true
    end)
end

function M.update_password(phone_number, password_hash)
    return box.atomic(function()
        local account = box.space.auths.index.phone_number:get{phone_number}
        if account == nil then
            error(E.ACCOUNT_NOT_FOUND)
        end

        local user = box.space.users:get{account.user_id}
        if user == nil then
            error(E.USER_NOT_FOUND)
        end

        if user.is_active == false then
            error(E.ALREADY_NOT_ACTIVE)
        end

        local now = os.time()
        box.space.auths:update(account.id, {
            {"=", 7, password_hash},
            {"=", 8, now},
        })

        return true
    end)
end

function M.get_account_by_phone_number(phone_number)
    return box.atomic(function()
        local account = box.space.auths.index.phone_number:get{phone_number}
        if account == nil then
            error(E.ACCOUNT_NOT_FOUND)
        end

        local user = box.space.users:get{account.user_id}
        if user == nil then
            error(E.USER_NOT_FOUND)
        end
        if user.is_active == false then
            error(E.ALREADY_NOT_ACTIVE)
        end

        local role = box.space.roles:get{user.role_id}
        if role == nil then
            error(E.ROLE_NOT_FOUND)
        end
        
        return {user.id, user.first_name, user.last_name, role.name, account.password_hash}
    end)
end

return M
