local E = require("errors.errors")

local U = {}

function U.get_user_info_by_phone_number(phone_number)
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
        
        return {user.id, user.first_name, user.last_name, role.name, user.is_active, user.created_at, user.updated_at, user.avatar_url, account.phone_number}
    end)
end

function U.update_user_info_by_phone_number(phone_number, new_phone_number, first_name, last_name, avatar_url)
    return box.atomic(function()
        local account = box.space.auths.index.phone_number:get{phone_number}
        if account == nil then
            error(E.ACCOUNT_NOT_FOUND)
        end

        if new_phone_number == "" then
            new_phone_number = nil
        end

        if new_phone_number ~= nil and new_phone_number ~= phone_number then
            if box.space.auths.index.phone_number:get{new_phone_number} then
                error(E.PHONE_ALREADY_EXISTS)
            end
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

        local now = os.time()
        box.space.users:update(user.id, {
            {"=", 2, first_name},
            {"=", 3, last_name},
            {"=", 8, now},
            {"=", 9, avatar_url},
            {"=", 5, user.version + 1},
        })

        if new_phone_number ~= nil and new_phone_number ~= phone_number then
            box.space.auths:update(account.id, {
                {"=", 3, new_phone_number},
                {"=", 8, now},
            })
        end

        local actual_phone = new_phone_number or phone_number

        return {user.id, first_name, last_name, role.name, user.is_active, user.created_at, now, avatar_url, actual_phone}
    end)
end

return U
