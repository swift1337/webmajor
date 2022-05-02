# Webmajor 🕵🏻‍♂️
Simple tool for inspecting and proxying http data of your app


## Usage
For example, if `localhost:8080` is your development application, each request
to `http://localhost:9999/...` will be proxied to `http://localhost:8080/...` and saved for later inspection.

```bash
./webmajor -source http://localhost:8080 -service-port 9999
May  2 22:37:10.000 INF starting server
May  2 22:37:10.000 INF source base: http://localhost:8080
May  2 22:37:10.000 INF to visit dashboard, open http://localhost:9999/__webmajor
```


## Installation
You can check repo releases and download binary for your platform or build service manually from sources


#### Why "Webmajor"?
In Russia, there's a popular meme about *[Tovarishch](https://en.wikipedia.org/wiki/Tovarishch) Major*, a military guy
who always spies on ordinary people on the Internet (like man-in-the-middle)

<img src="https://user-images.githubusercontent.com/11892559/166313917-c258fa10-7398-4847-b655-d6ac2ffd0f79.jpeg" style="width: 200px">
