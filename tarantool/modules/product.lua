local E = require("errors.errors")

local P = {}

function P.create_product(first_name, last_name, phone_number, role_name)
    local now = os.time()

    return box.atomic(function()
        
    end)
end


return P