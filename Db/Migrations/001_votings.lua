local username =  os.getenv('TARANTOOL_USER') or 'user'
local password = os.getenv('TARANTOOL_PASS') or 'password'


-- Проверяем, существует ли пользователь. Если нет, создаем его.
if box.schema.user.exists(username) == false then
    box.schema.user.create(username, {password = password})
    box.schema.user.grant(username, 'read,write,execute', 'universe')
end

-- Для работы с индексами
local voting_id_seq = box.schema.sequence.create('voting_id_seq', { if_not_exists = true, min=0, start=0 })

local counters = box.schema.create_space('counters', { if_not_exists = true })

counters:format({
    { name = 'counter_name', type = 'string' },
    { name = 'last_id', type = 'unsigned' }
})

counters:create_index('primary', { parts = { 'counter_name', }, if_not_exists = true })

function get_next_voting_id()
    return voting_id_seq:next()
end

function get_next_options_id()
    local counter_name = 'options_id'
    local result = counters:select({counter_name})

    if #result > 0 then
        local updated = counters:update(
            {counter_name},
            {{'+', 2, 1}}
        )
        return updated[2]
    else
        counters:insert({counter_name, 1})
        return 1
    end
end

-- Спэйс для хранения голосований
local votings = box.schema.create_space('votings', { if_not_exists = true })

votings:format({
    { name = 'id', type = 'unsigned' },
    { name = 'creator_id', type = 'unsigned' },
    { name = 'question', type = 'string' },
    { name = 'options_id', type = 'unsigned' },
    { name = 'is_active', type = 'boolean' }
})

votings:create_index('primary', { sequence = 'voting_id_seq', if_not_exists = true })

-- Спэйс для хранения вариантов голосований
local options = box.schema.create_space('options', { if_not_exists = true })

options:format({
    { name = 'options_id', type = 'unsigned' },
    { name = 'option_id', type = 'unsigned' },
    {name = 'text', type = 'string'},
    {name = 'count', type = 'integer'}
    
})
options:create_index('options_id_idx', { parts = { 'options_id', 'option_id', }, if_not_exists = true })


