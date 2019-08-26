

# TODO

- impl programmatical API

    - extract cli
    
- refactor cli to 0 required args:
    - [--watch WATCHED_DIR ...]
    - [--poll-interval POLL_INTERVAL]
        - (also rename PollEvery to PollInterval, 
        list a watched dir every `PollInterval/len(WatchedDirs)` 
        to smoothen the fs stress)
        https://golang.org/pkg/testing/
        https://www.toptal.com/go/your-introductory-course-to-testing-with-go
    - [--nginx NGINX_CLI_ARG ...]

- unit tests
