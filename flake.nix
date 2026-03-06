{
  description = "dndgoldtracker";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable-small";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
    nur.url = "github:nix-community/NUR";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gomod2nix,
    nur,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [gomod2nix.overlays.default nur.overlays.default];
        };
        lib = pkgs.lib;
        version =
          if (self ? shortRev)
          then self.shortRev
          else "dev";
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            glow
            go
            gomod2nix.packages.${system}.default
            gopls
            goreleaser # may be optional
            imagemagick # used with `vhs`
            mise # task runner
            nixd # nix lsp for flake editing
            nodejs
            upx
            # uv
            vhs # for making demos
          ];

          shellHook = ''
            mise trust --all
          '';

          # necessary for `-race` in tests
          CGO_ENABLED = "1";
        };
      }
    )
    // {
      overlays.default = final: prev: {};
    };
}
