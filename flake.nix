{
  description = "a terminal-based notes and tasks organizer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };

        pname = "gobrain";
        version = "1.0.4";
        author = "haykh";

        src = pkgs.fetchFromGitHub {
          inherit pname version;
          owner = author;
          repo = pname;
          rev = "master";
          hash = "sha256-V2ToDRnrOq4hul/h8iN5VS9YZLM0mlOcEw907m8WOvg=";
        };
      in
      {
        packages.default = pkgs.buildGoModule {
          inherit pname version src;

          vendorHash = "sha256-AfdgJceYwgZB0lVpVRWYoGkbBQWWCV+FHx33Ru7lCUM=";

          meta = with pkgs.lib; {
            homepage = "https://github.com/${author}/gobrain";
            license = licenses.unlicense;
            maintainers = [ author ];
          };
        };

        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.gopls
            pkgs.gotools
          ];

          shellHook = ''
            echo "Welcome to gobrain development shell!"
          '';
        };
      }
    );
}
