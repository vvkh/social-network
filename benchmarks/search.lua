wrk.method = "GET"
wrk.path = "/profiles/?"

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
    if not user_created then
        return wrk.format(
            "POST",
            "/register/",
            {["Content-Type"] = "application/x-www-form-urlencoded"},
            "age=18&first_name=wrk&last_name=wrk&password=wrk&sex=male&username=wrk"
        )
    end
    if not authenticated then
        return wrk.format(
            "POST",
            "/login/",
            {["Content-Type"] = "application/x-www-form-urlencoded"},
            "password=wrk&username=wrk"
        )
    end

    counter = counter + 1
    if counter > #requests then
        counter = 0
    end
    queryParams = requests[counter]
    return wrk.format(
        wrk.method,
        wrk.path .. queryParams,
        wrk.headers,
        ""
    )
end

response = function(status, headers, body)
    if not user_created then
        print("user wrk created")
        user_created = true
        return
    end
    if not authenticated then
        cookie = headers["Set-Cookie"]
        wrk.headers["Cookie"] = cookie
        print("authenticated wrk with cookie ", cookie)
        authenticated = true
        return
    end
end
