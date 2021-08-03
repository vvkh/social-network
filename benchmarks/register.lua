wrk.method = "POST"
wrk.path = "/register/"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

function read_file(file)
    lines = {}
    for line in io.lines(file) do
        lines[#lines + 1] = line
    end
    return lines
end
requests = read_file("benchmarks/requests/register.txt")

local thread_id_counter = 0
local threads = {}
function setup(thread)
   thread:set("request_id", 1)
   thread:set("thread_id", thread_id_counter)
   table.insert(threads, thread)
   for _, t in ipairs(threads) do
      t:set("thread_count", thread_id_counter + 1)
   end
   thread_id_counter = thread_id_counter + 1
end

request = function()
    requests_per_thread = math.floor(#requests / thread_count)
    min_request = requests_per_thread * thread_id
    max_request = min_request + requests_per_thread - 1
    counter = min_request + request_id
    request_id = request_id + 1

    if counter > max_request then
        wrk.thread:stop()
    end
    body = requests[counter]
    if body == nil then
    	print("counter = ", counter)
    end
    return wrk.format(
        wrk.method,
        wrk.path,
        wrk.headers,
        body
    )
end
