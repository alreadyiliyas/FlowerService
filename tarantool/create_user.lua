-- /c:/Users/ILYAS/GolandProjects/Flower/tarantool/create_user.lua
local E = require("errors")

return function(first_name, last_name, phone_number, role_name)
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
