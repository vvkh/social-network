wrk.method = "POST"
wrk.path = "/register/"
wrk.headers["Content-type"] = "application/x-www-form-urlencoded"
wrk.body = "age=18&first_name=wrk&last_name=wrk&password=wrk&sex=male&username=wrk"

function read_file(file)
    lines = {}
    for line in io.lines(file) do
        lines[#lines + 1] = line
    end
    return lines
end

requests = read_file("benchmarks/requests/search.txt")

counter = 0
user_created = nil
authenticated = nil

request = function()
    return wrk.format(
        wrk.method,
        wrk.path,
        wrk.headers,
        wrk.body
    )
end


response = function(status, headers, body)
    if not user_created then
        user_created = true
        wrk.path = "/login/"
        wrk.body = "password=wrk&username=wrk"
        return
    end
    if not authenticated and status == 302 then
        authenticated = true
        cookie = headers["Set-Cookie"]
        wrk.headers["Cookie"] = cookie

        wrk.method = "GET"
        wrk.headers["Content-type"] = ""
        wrk.body = ""
    end
    if not authenticated then
        return
    end
    counter = counter + 1
    if counter >= #requests then
        counter = 1
    end
    queryParams = requests[counter]
    wrk.path = "/profiles/?" .. queryParams
end
