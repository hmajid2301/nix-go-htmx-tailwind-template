{ pkgs ? (let
  inherit (builtins) fetchTree fromJSON readFile;
  inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
in import (fetchTree nixpkgs.locked) {
  overlays = [ (import "${fetchTree gomod2nix.locked}/overlay.nix") ];
}), mkGoEnv ? pkgs.mkGoEnv, gomod2nix ? pkgs.gomod2nix, pre-commit-hooks
, myPackages, ... }:
let
  goEnv = mkGoEnv { pwd = ./.; };
  pre-commit-check = pre-commit-hooks.lib.${pkgs.system}.run {
    src = ./.;
    hooks = {
      golangci-lint.enable = true;
      gotest.enable = true;
      gofmt.enable = true;
      generate = {
        enable = true;
        name = "Generate code such as tailwindcss, templ components and sqlc";
        entry = "task generate";
        pass_filenames = false;
      };

    };
  };
in pkgs.mkShell {
  hardeningDisable = [ "all" ];
  shellHook = ''
    export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
    export PLAYWRIGHT_NODEJS_PATH="${pkgs.nodejs}/bin/node"
    export PLAYWRIGHT_BROWSERS_PATH="${pkgs.playwright-driver.browsers}"
    export GOOSE_MIGRATION_DIR="./db/migrations"
    ${pre-commit-check.shellHook}
  '';
  buildInputs = pre-commit-check.enabledPackages;
  packages = myPackages ++ [ goEnv gomod2nix ];
}
