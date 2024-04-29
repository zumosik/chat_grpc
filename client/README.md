# Go + Templ + Tailwind + HTMX
> GoTTH - A simple, modern stack for building fast web applications.

This project is template with all needed to be installed and configured for apps using **GoTTH** stack.
* Go - Backend
* Tailwind - CSS
* Templ - Templating
* HTMX - Interactivity

## How to run?
1. Install tailwind css binary  
    https://tailwindcss.com/blog/standalone-cli  
  
    For linux you can use this command:  
    ```bash
    make install-tailwind-linux
    ```
2. Run air for hot reload  
    ```bash
    make dev
    ```
3. Run tailwind watch in another terminal  
    ```bash
    make tailwind-watch
    ```
4. Open browser and go to port localhost:3032 **(port can be changed in configs/local.yml)**

## Makefile commands
```bash
make tailwind-watch
```
This target watches the ./static/css/input.css file and automatically rebuilds the Tailwind CSS styles whenever changes are detected.
(tailwindcss bin required)
```bash
make tailwind-build
```
This target minifies the Tailwind CSS styles by running the tailwindcss command.
(tailwindcss bin required)

```bash
make templ-watch
```
This target watches for changes to *.templ files and automatically generates them.
(templ required)

```bash
make templ-generate
```
This target generates templates using the templ command.
(templ required)

```bash
make dev
```
This target runs the development server using Air, which helps in hot-reloading your Go application during development.
(Air, templ and tailwindcss bin required)
```bash
make build
```
This target orchestrates the building process by executing the tailwind-build, templ-generate, and go build commands sequentially. It creates the binary output in the ./bin/ directory.
(templ and tailwindcss bin required)
```bash
make test
```
This target runs all Go tests for the application.