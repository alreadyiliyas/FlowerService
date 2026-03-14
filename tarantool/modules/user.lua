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
        
        return {user.id, user.first_name, user.last_name, role.name, user.is_active, user.created_at, user.updated_at}
    end)
end

return U