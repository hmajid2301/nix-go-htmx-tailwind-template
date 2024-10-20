
# Nix, Go, htmx & Tailwind CSS Template

**Currently a WIP may not work**

A template for creating new Go "full-stack" web services. Powered by **Nix** package manager. We will leverage
the reproducibility of Nix to build our Go web-service as a Nix package and then use Nix to build our Docker image.

## Background

The aim of this template is to focus mostly on using Go and needing to write as little frontend as needed especially
JS/TS. Hence, even playwright tests for this project are written in Go. I am a back-end software developer hence I want
to keep the front-end stack as simple as possible, i.e. not need to use node to install dependencies if possible.

## Stack

  - Backend
    - Go
      - Standard Library HTTP Server
    - SQLite DB
  - Frontend
    - htmx
    - Tailwind CSS

## Features

  - Powered by Go
  - htmx from the server
  - templ as the templating engine
  - Live reloading powered by air
  - Nix for reproducibility
      - gomod2nix to build go binary with Nix
      - Development Shells
      - Pre Commit Hooks
      - Build Docker images
        - CI Image
        - Service
  - Gitlab for CI/CD pipeline
  - Taskfiles as the task runner
  - Copier template management
  - Testing
      - `gotestsum` as the test runner
        - Better output
      - Output for JUnit and Cobertura reports (code coverage)
      - E2E with playwright
  - Deploy to K8s
  - Renovate for automatic dependency updates
  - Database
      - code generated using **sqlc**
      - Migrations with Goose


## Usage

```bash
nix-shell -p copier
copier copy https://gitlab.com/hmajid2301/nix-go-htmx-tailwind-template banterbus
```
